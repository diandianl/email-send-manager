package setting

import (
	"context"
	"gorm.io/gorm/clause"

	"github.com/google/wire"
	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/errors"
)

// SettingSet Injection wire
var SettingSet = wire.NewSet(wire.Struct(new(SettingRepo), "*"))

type SettingRepo struct {
	DB *gorm.DB
}

func (a *SettingRepo) getQueryOption(opts ...schema.SettingQueryOptions) schema.SettingQueryOptions {
	var opt schema.SettingQueryOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func (a *SettingRepo) Query(ctx context.Context) (*schema.SettingQueryResult, error) {
	db := GetSettingDB(ctx, a.DB)

	var list Settings
	err := db.Find(&list).Error
	if err != nil {
		return nil, errors.WithStack(err)
	}

	qr := &schema.SettingQueryResult{
		Data:       list.ToSchemaSettings(),
	}
	return qr, nil
}

func (a *SettingRepo) Get(ctx context.Context, key string) (*schema.Setting, error) {
	var item Setting
	ok, err := util.FindOne(ctx, GetSettingDB(ctx, a.DB).Where("key=?", key), &item)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !ok {
		return nil, nil
	}

	return item.ToSchemaSetting(), nil
}

func (a *SettingRepo) Create(ctx context.Context, item schema.Setting) error {
	eitem := SchemaSetting(item).ToSetting()
	result := GetSettingDB(ctx, a.DB).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value"}),
	}).Create(eitem)
	return errors.WithStack(result.Error)
}
