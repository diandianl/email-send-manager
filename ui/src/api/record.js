import request from '@/utils/request'

// 查询发送记录列表
export function queryRecord(query) {
  return request({
    url: '/api/v1/records',
    method: 'get',
    params: query
  })
}

// 删除发送记录
export function delRecord(recordId) {
  return request({
    url: '/api/v1/records/' + recordId,
    method: 'delete'
  })
}

