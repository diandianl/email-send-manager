import request from '@/utils/request'

// 查询配置列表
export function querySettings(key) {
  return request({
    url: '/api/v1/settings/' + key,
    method: 'get'
  })
}

// upsert settings
export function upsertSettings(data) {
  return request({
    url: '/api/v1/settings',
    method: 'post',
    data: data
  })
}
