package service

import (
	"context"

	"github.com/google/wire"

	"email-send-manager/internal/app/dao"
	"email-send-manager/internal/app/schema"
)

var SettingSet = wire.NewSet(wire.Struct(new(SettingSrv), "*"))

type SettingSrv struct {
	TransRepo   *dao.TransRepo
	SettingRepo *dao.SettingRepo
}

func (a *SettingSrv) Get(ctx context.Context, key string) (*schema.Setting, error) {
	result, err := a.SettingRepo.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if result == nil {
		result = &schema.Setting{Key: key}
	}
	return result, nil
}

func (a *SettingSrv) Upsert(ctx context.Context, item schema.Setting) error {
	return a.SettingRepo.Create(ctx, item)
}
