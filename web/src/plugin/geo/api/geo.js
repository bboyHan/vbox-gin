import service from '@/utils/request'

// @Tags System
// @Summary 发送邮件
// @Security ApiKeyAuth
// @Produce  application/json
// @Param data body email_response.Email true "发送邮件必须的参数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /get/getGeos [post]
export const getGeos = (params) => {
  return service({
    url: '/geo/getGeos',
    method: 'get',
    params
  })
}

export const editGeo = (data) => {
  return service({
    url: '/geo/editGeo',
    method: 'put',
    data
  })
}

export const createGeo = (data) => {
  return service({
    url: '/geo/createGeo',
    method: 'post',
    data
  })
}

export const getGeo = (params) => {
  return service({
    url: '/geo/getGeo',
    method: 'get',
    params
  })
}

export const deleteGeo = (params) => {
  return service({
    url: '/geo/deleteGeo',
    method: 'delete',
    params
  })
}

