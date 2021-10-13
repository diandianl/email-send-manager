package send_batch

import (
	"context"

	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/util/structure"
)

// GetSendBatchDB Get SendBatch db model
func GetSendBatchDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(SendBatch))
}

// SchemaSendBatch SendBatch schema
type SchemaSendBatch schema.SendBatch

// ToSendBatch Convert to SendBatch entity
func (a SchemaSendBatch) ToSendBatch() *SendBatch {
	item := new(SendBatch)
	structure.Copy(a, item)
	return item
}

// SendBatch SendBatch entity
type SendBatch struct {
	util.Model
	TemplateID       uint64   `gorm:""`                              // 模板ID
	ReverseSelection bool     `gorm:""`                              // 是否是反选客户
	CustomerIDs      []uint64 `gorm:""`                              // 客户列表
}

// ToSchemaSendBatch Convert to SendBatch schema
func (a SendBatch) ToSchemaSendBatch() *schema.SendBatch {
	item := new(schema.SendBatch)
	structure.Copy(a, item)
	return item
}

// SendBatches SendBatch entity list
type SendBatches []*SendBatch

// ToSchemaSendBatches Convert to SendBatch schema list
func (a SendBatches) ToSchemaSendBatches() []*schema.SendBatch {
	list := make([]*schema.SendBatch, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaSendBatch()
	}
	return list
}
