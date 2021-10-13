package service

import (
	"context"

	"github.com/google/wire"

	"email-send-manager/internal/app/dao"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
	"email-send-manager/pkg/util/snowflake"
)

var SendBatchSet = wire.NewSet(wire.Struct(new(SendBatchSrv), "*"))

type SendBatchSrv struct {
	TransRepo     *dao.TransRepo
	SendBatchRepo *dao.SendBatchRepo
}

func (a *SendBatchSrv) Query(ctx context.Context, params schema.SendBatchQueryParam, opts ...schema.SendBatchQueryOptions) (*schema.SendBatchQueryResult, error) {
	result, err := a.SendBatchRepo.Query(ctx, params, opts...)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (a *SendBatchSrv) Get(ctx context.Context, id uint64, opts ...schema.SendBatchQueryOptions) (*schema.SendBatch, error) {
	item, err := a.SendBatchRepo.Get(ctx, id, opts...)
	if err != nil {
		return nil, err
	} else if item == nil {
		return nil, errors.ErrNotFound
	}

	return item, nil
}

func (a *SendBatchSrv) Create(ctx context.Context, item schema.SendBatch) (*schema.IDResult, error) {
	item.ID = snowflake.MustID()

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.SendBatchRepo.Create(ctx, item)
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(item.ID), nil
}

func (a *SendBatchSrv) Update(ctx context.Context, id uint64, item schema.SendBatch) error {
	oldItem, err := a.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	item.ID = oldItem.ID
	item.Creator = oldItem.Creator
	item.CreatedAt = oldItem.CreatedAt

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.SendBatchRepo.Update(ctx, id, item)
	})
}

func (a *SendBatchSrv) Delete(ctx context.Context, id uint64) error {
	oldItem, err := a.SendBatchRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	}

	return a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		return a.SendBatchRepo.Delete(ctx, id)
	})
}

func (a *SendBatchSrv) UpdateStatus(ctx context.Context, id uint64, status int) error {
	oldItem, err := a.SendBatchRepo.Get(ctx, id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.ErrNotFound
	} else if oldItem.Status == status {
		return nil
	}

	return a.SendBatchRepo.UpdateStatus(ctx, id, status)
}
