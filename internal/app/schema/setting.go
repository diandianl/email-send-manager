package schema

import "encoding/json"

// Setting 系统配置对象
type Setting struct {
	Key   string          `json:"key"`   // 配置Key
	Value json.RawMessage `json:"value"` // 配置值
}

// SettingQueryParam 查询条件
type SettingQueryParam struct {
	PaginationParam
}

// SettingQueryOptions 查询可选参数项
type SettingQueryOptions struct {
	OrderFields  []*OrderField // 排序字段
	SelectFields []string      // 查询字段
}

// SettingQueryResult 查询结果
type SettingQueryResult struct {
	Data       Settings
}

// Settings 系统配置列表
type Settings []*Setting
