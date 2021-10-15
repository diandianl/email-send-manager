package setting

import (
	"context"
	"gorm.io/datatypes"

	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/util/structure"
)

// GetSettingDB Get Setting db model
func GetSettingDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Setting))
}

// SchemaSetting Setting schema
type SchemaSetting schema.Setting

// ToSetting Convert to Setting entity
func (a SchemaSetting) ToSetting() *Setting {
	item := new(Setting)
	structure.Copy(a, item)
	return item
}

// Setting Setting entity
type Setting struct {
	Key   string         `gorm:"primaryKey;"` // 配置Key
	Value datatypes.JSON `gorm:""`            // 配置值
}

// ToSchemaSetting Convert to Setting schema
func (a Setting) ToSchemaSetting() *schema.Setting {
	item := new(schema.Setting)
	structure.Copy(a, item)
	return item
}

// Settings Setting entity list
type Settings []*Setting

// ToSchemaSettings Convert to Setting schema list
func (a Settings) ToSchemaSettings() []*schema.Setting {
	list := make([]*schema.Setting, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaSetting()
	}
	return list
}
