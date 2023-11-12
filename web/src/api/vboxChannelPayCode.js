import service from '@/utils/request'

// @Tags VboxChannelPayCode
// @Summary 创建通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelPayCode true "创建通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /vboxChannelPayCode/createVboxChannelPayCode [post]
export const createVboxChannelPayCode = (data) => {
  return service({
    url: '/vboxChannelPayCode/createVboxChannelPayCode',
    method: 'post',
    data
  })
}

// @Tags VboxChannelPayCode
// @Summary 删除通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelPayCode true "删除通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vboxChannelPayCode/deleteVboxChannelPayCode [delete]
export const deleteVboxChannelPayCode = (data) => {
  return service({
    url: '/vboxChannelPayCode/deleteVboxChannelPayCode',
    method: 'delete',
    data
  })
}

// @Tags VboxChannelPayCode
// @Summary 批量删除通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vboxChannelPayCode/deleteVboxChannelPayCode [delete]
export const deleteVboxChannelPayCodeByIds = (data) => {
  return service({
    url: '/vboxChannelPayCode/deleteVboxChannelPayCodeByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxChannelPayCode
// @Summary 更新通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelPayCode true "更新通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vboxChannelPayCode/updateVboxChannelPayCode [put]
export const updateVboxChannelPayCode = (data) => {
  return service({
    url: '/vboxChannelPayCode/updateVboxChannelPayCode',
    method: 'put',
    data
  })
}

// @Tags VboxChannelPayCode
// @Summary 用id查询通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxChannelPayCode true "用id查询通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vboxChannelPayCode/findVboxChannelPayCode [get]
export const findVboxChannelPayCode = (params) => {
  return service({
    url: '/vboxChannelPayCode/findVboxChannelPayCode',
    method: 'get',
    params
  })
}

// @Tags VboxChannelPayCode
// @Summary 分页获取通道账户付款二维码列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取通道账户付款二维码列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vboxChannelPayCode/getVboxChannelPayCodeList [get]
export const getVboxChannelPayCodeList = (params) => {
  return service({
    url: '/vboxChannelPayCode/getVboxChannelPayCodeList',
    method: 'get',
    params
  })
}
