package service

import (
	"context"

	"github.com/google/wire"

	"email-send-manager/internal/app/dao"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
)

var RecordSet = wire.NewSet(wire.Struct(new(RecordSrv), "*"))

type RecordSrv struct {
	TransRepo  *dao.TransRepo
	RecordRepo *dao.RecordRepo
}

func (a *RecordSrv) Query(ctx context.Context, params schema.RecordQueryParam, opts ...schema.RecordQueryOptions) (*schema.RecordQueryResult, error) {
	result, err := a.RecordRepo.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *RecordSrv) Get(ctx context.Context, id uint, opts ...schema.RecordQueryOptions) (*schema.Record, error) {
	item, err := a.RecordRepo.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *RecordSrv) Create(ctx context.Context, item schema.Record) (*schema.IDResult, error) {

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.RecordRepo.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *RecordSrv) Update(ctx context.Context, id uint, item schema.Record) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	item.ID = oldItem.ID
	item.CreatedAt = oldItem.CreatedAt

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.RecordRepo.Update(ctx, id, item)
	})
}

func (a *RecordSrv) Delete(ctx context.Context, id uint) error {
	oldItem, err := a.RecordRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.RecordRepo.Delete(ctx, id)
	})
}

func (a *RecordSrv) UpdateStatus(ctx context.Context, id uint, status int) error {
	oldItem, err := a.RecordRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Status == status {
		return nil
	}

	return a.RecordRepo.UpdateStatus(ctx, id, status)
}
