package service

import (
	"context"

	"github.com/google/wire"

	"email-send-manager/internal/app/dao"
	"email-send-manager/internal/app/schema"
)

var SendBatchSet = wire.NewSet(wire.Struct(new(SendBatchSrv), "TransRepo", "CustomerRepo", "TemplateRepo"))

type SendBatchSrv struct {
	current *schema.SendBatch

	TransRepo    *dao.TransRepo
	CustomerRepo *dao.CustomerRepo
	TemplateRepo *dao.TemplateRepo
}

func (a *SendBatchSrv) Current(ctx context.Context, id uint, opts ...schema.SendBatchQueryOptions) (*schema.SendBatch, error) {
	return a.current, nil
}

func (a *SendBatchSrv) Create(ctx context.Context, item schema.SendBatch) (*schema.IDResult, error) {

	err := a.TransRepo.Exec(ctx, func(ctx context.Context) error {
		tpl, err := a.TemplateRepo.Get(ctx, item.TemplateID)
		if err != nil {
			return err
		}
		if item.TemplateID == tpl.ID {
			// TODO
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return schema.NewIDResult(0), nil
}
