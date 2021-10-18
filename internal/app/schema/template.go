package schema

import "time"

// Template 邮件模板管理对象
type Template struct {
	ID        uint      `json:"id"`                                   // 唯一标识
	Name      string    `json:"name" binding:"required"`              // 名称
	From      string    `json:"from,omitempty" binding:"email"`       // 发件人邮箱
	FromName  string    `json:"from_name,omitempty"`                  // 联系人名称
	ReplyTo   string    `json:"reply_to,omitempty"`                   // 回复邮箱地址
	Subject   string    `json:"subject,omitempty" binding:"required"` // 邮件主题
	Content   string    `json:"content,omitempty" binding:"required"` // 邮件正文模板
	Status    int       `json:"status,omitempty"`                     // 状态(1:启用 2:禁用)
	CreatedAt time.Time `json:"created_at,omitempty"`                 // 创建时间
	UpdatedAt time.Time `json:"updated_at,omitempty"`                 // 更新时间
}

// TemplateQueryParam 查询条件
type TemplateQueryParam struct {
	PaginationParam
	Lite   bool   `form:"lite"`
	Name   string `form:"name"`
	Status int    `form:"status"`
}

// TemplateQueryOptions 查询可选参数项
type TemplateQueryOptions struct {
	OrderFields  []*OrderField // 排序字段
	SelectFields []string      // 查询字段
}

// TemplateQueryResult 查询结果
type TemplateQueryResult struct {
	Data       Templates
	PageResult *PaginationResult
}

// Templates 邮件模板管理列表
type Templates []*Template
