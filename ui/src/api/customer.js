import request from '@/utils/request'

// 查询客户列表
export function queryCustomer(query) {
  return request({
    url: '/api/v1/customers',
    method: 'get',
    params: query
  })
}

// 获取客户详细
export function getCustomer(id) {
  return request({
    url: '/api/v1/customers/' + id,
    method: 'get'
  })
}

// 新增客户
export function addCustomer(data) {
  return request({
    url: '/api/v1/customers',
    method: 'post',
    data: data
  })
}

// 修改岗位
export function updateCustomer(data, id) {
  return request({
    url: '/api/v1/customers/' + id,
    method: 'put',
    data: data
  })
}

// 删除岗位
export function delCustomer(customerId) {
  return request({
    url: '/api/v1/customers/' + customerId,
    method: 'delete'
  })
}

export const importUrl = process.env.VUE_APP_BASE_API + '/api/v1/customers/import'
