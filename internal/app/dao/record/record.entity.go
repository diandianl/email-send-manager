package record

import (
	"context"

	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/util/structure"
)

// GetRecordDB Get Record db model
func GetRecordDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Record))
}

// SchemaRecord Record schema
type SchemaRecord schema.Record

// ToRecord Convert to Record entity
func (a SchemaRecord) ToRecord() *Record {
	item := new(Record)
	structure.Copy(a, item)
	return item
}

// Record Record entity
type Record struct {
	util.Model
	TemplateID uint64 `gorm:""`                              // 模板ID
	CustomerID uint64 `gorm:""`                              // 客户ID
	Status     int    `gorm:"type:tinyint;index;default:0;"` // 状态(0:成功 1:失败)
	Reason     string `gorm:""`                              // 失败原因
}

// ToSchemaRecord Convert to Record schema
func (a Record) ToSchemaRecord() *schema.Record {
	item := new(schema.Record)
	structure.Copy(a, item)
	return item
}

// Records Record entity list
type Records []*Record

// ToSchemaRecords Convert to Record schema list
func (a Records) ToSchemaRecords() []*schema.Record {
	list := make([]*schema.Record, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaRecord()
	}
	return list
}
