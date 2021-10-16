import request from '@/utils/request'

// 查询当前发送进度
export function queryCurrent() {
  return request({
    url: '/api/v1/send-batches/current',
    method: 'get'
  })
}

// 删除发送记录
export function doSendBatch(batch) {
  return request({
    url: '/api/v1/send-batches',
    method: 'post',
    data: batch
  })
}

export function cancelSend() {
  return request({
    url: '/api/v1/send-batches',
    method: 'delete'
  })
}

