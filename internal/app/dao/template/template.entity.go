package template

import (
	"context"

	"gorm.io/gorm"

	"email-send-manager/internal/app/dao/util"
	"email-send-manager/internal/app/schema"
	"email-send-manager/pkg/util/structure"
)

// GetTemplateDB Get Template db model
func GetTemplateDB(ctx context.Context, defDB *gorm.DB) *gorm.DB {
	return util.GetDBWithModel(ctx, defDB, new(Template))
}

// SchemaTemplate Template schema
type SchemaTemplate schema.Template

// ToTemplate Convert to Template entity
func (a SchemaTemplate) ToTemplate() *Template {
	item := new(Template)
	structure.Copy(a, item)
	return item
}

// Template Template entity
type Template struct {
	util.Model
	Name     string `gorm:"size:100;index;"`               // 名称
	From     string `gorm:"size:100;"`                     // 发件人邮箱
	FromName string `gorm:"size:100;"`                     // 联系人名称
	ReplyTo  string `gorm:"size:100;"`                     // 回复邮箱地址
	Subject  string `gorm:"size:100;"`                     // 邮件主题
	Content  string `gorm:"size:100;"`                     // 邮件正文模板
	Status   int    `gorm:"type:tinyint;index;default:1;"` // 状态(1:启用 0:停用)
}

// ToSchemaTemplate Convert to Template schema
func (a Template) ToSchemaTemplate() *schema.Template {
	item := new(schema.Template)
	structure.Copy(a, item)
	return item
}

// Templates Template entity list
type Templates []*Template

// ToSchemaTemplates Convert to Template schema list
func (a Templates) ToSchemaTemplates() []*schema.Template {
	list := make([]*schema.Template, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaTemplate()
	}
	return list
}
