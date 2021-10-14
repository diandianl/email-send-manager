package schema

import "time"

// SendBatch 邮件发送批次对象
type SendBatch struct {
	TemplateID       uint    `json:"template_id" binding:"required"` // 模板ID
	ReverseSelection bool      `json:"reverse_selection"`              // 是否是反选客户
	CustomerIDs      []uint  `json:"customer_ids"`                   // 客户列表
	CreatedAt        time.Time `json:"created_at"`                     // 创建时间
}

// SendBatchQueryParam 查询条件
type SendBatchQueryParam struct {
	PaginationParam
}

// SendBatchQueryOptions 查询可选参数项
type SendBatchQueryOptions struct {
	OrderFields  []*OrderField // 排序字段
	SelectFields []string      // 查询字段
}

// SendBatchQueryResult 查询结果
type SendBatchQueryResult struct {
	Data       SendBatches
	PageResult *PaginationResult
}

// SendBatches 邮件发送批次列表
type SendBatches []*SendBatch
