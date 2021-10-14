import request from '@/utils/request'

// 查询客户列表
export function queryTemplate(query) {
  return request({
    url: '/api/v1/templates',
    method: 'get',
    params: query
  })
}

// 获取客户详细
export function getTemplate(id) {
  return request({
    url: '/api/v1/templates/' + id,
    method: 'get'
  })
}

// 新增客户
export function addTemplate(data) {
  return request({
    url: '/api/v1/templates',
    method: 'post',
    data: data
  })
}

// 修改岗位
export function updateTemplate(data, id) {
  return request({
    url: '/api/v1/templates/' + id,
    method: 'put',
    data: data
  })
}

// 删除岗位
export function delTemplate(templateId) {
  return request({
    url: '/api/v1/templates/' + templateId,
    method: 'delete'
  })
}

