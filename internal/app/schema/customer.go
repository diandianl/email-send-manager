package schema

import "time"

// Customer 客户管理对象
type Customer struct {
	ID        uint    `json:"id"`                      // 唯一标识
	Name      string    `json:"name" binding:"required"` // 名称
	Email     string    `json:"email"`                   // 邮箱
	Status    int       `json:"status"`                  // 状态(1:启用 2:禁用)
	CreatedAt time.Time `json:"created_at"`              // 创建时间
	UpdatedAt time.Time `json:"updated_at"`              // 更新时间
}

// CustomerQueryParam 查询条件
type CustomerQueryParam struct {
	PaginationParam
}

// CustomerQueryOptions 查询可选参数项
type CustomerQueryOptions struct {
	OrderFields  []*OrderField // 排序字段
	SelectFields []string      // 查询字段
}

// CustomerQueryResult 查询结果
type CustomerQueryResult struct {
	Data       Customers
	PageResult *PaginationResult
}

// Customers 客户管理列表
type Customers []*Customer
