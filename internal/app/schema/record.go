package schema

import "time"

// Record 邮件发送记录对象
type Record struct {
	ID         uint64    `json:"id"`          // 唯一标识
	TemplateID uint64    `json:"template_id"` // 模板ID
	CustomerID uint64    `json:"customer_id"` // 客户ID
	Status     int       `json:"status"`      // 结果状态(0:成功 1:失败)
	Reason     string    `json:"reason"`      // 失败原因
	CreatedAt  time.Time `json:"created_at"`  // 创建时间
	UpdatedAt  time.Time `json:"updated_at"`  // 更新时间
}

// RecordQueryParam 查询条件
type RecordQueryParam struct {
	PaginationParam
}

// RecordQueryOptions 查询可选参数项
type RecordQueryOptions struct {
	OrderFields  []*OrderField // 排序字段
	SelectFields []string      // 查询字段
}

// RecordQueryResult 查询结果
type RecordQueryResult struct {
	Data       Records
	PageResult *PaginationResult
}

// Records 邮件发送记录列表
type Records []*Record
