package service

import (
	"context"

	"github.com/google/wire"

	"email-send-manager/internal/app/dao"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
	"email-send-manager/pkg/util/snowflake"
)

var TemplateSet = wire.NewSet(wire.Struct(new(TemplateSrv), "*"))

type TemplateSrv struct {
	TransRepo    *dao.TransRepo
	TemplateRepo *dao.TemplateRepo
}

func (a *TemplateSrv) Query(ctx context.Context, params schema.TemplateQueryParam, opts ...schema.TemplateQueryOptions) (*schema.TemplateQueryResult, error) {
	result, err := a.TemplateRepo.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *TemplateSrv) Get(ctx context.Context, id uint64, opts ...schema.TemplateQueryOptions) (*schema.Template, error) {
	item, err := a.TemplateRepo.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *TemplateSrv) Create(ctx context.Context, item schema.Template) (*schema.IDResult, error) {
	item.ID = snowflake.MustID()

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.TemplateRepo.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *TemplateSrv) Update(ctx context.Context, id uint64, item schema.Template) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	item.ID = oldItem.ID
	item.CreatedAt = oldItem.CreatedAt

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.TemplateRepo.Update(ctx, id, item)
	})
}

func (a *TemplateSrv) Delete(ctx context.Context, id uint64) error {
	oldItem, err := a.TemplateRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.TemplateRepo.Delete(ctx, id)
	})
}

func (a *TemplateSrv) UpdateStatus(ctx context.Context, id uint64, status int) error {
	oldItem, err := a.TemplateRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Status == status {
		return nil
	}

	return a.TemplateRepo.UpdateStatus(ctx, id, status)
}
