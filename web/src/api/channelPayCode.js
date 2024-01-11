import service from '@/utils/request'

// @Tags VboxChannelPayCode
// @Summary 批量创建通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelPayCode true "创建通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /channelPayCode/createVboxChannelPayCode [post]
export const batchCreateChannelPayCode = (data) => {
  return service({
    url: '/channelPayCode/batchCreateChannelPayCode',
    method: 'post',
    data
  })
}

// @Tags VboxChannelPayCode
// @Summary 创建通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelPayCode true "创建通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /channelPayCode/createVboxChannelPayCode [post]
export const createChannelPayCode = (data) => {
  return service({
    url: '/channelPayCode/createChannelPayCode',
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
// @Router /channelPayCode/deleteVboxChannelPayCode [delete]
export const deleteChannelPayCode = (data) => {
  return service({
    url: '/channelPayCode/deleteChannelPayCode',
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
// @Router /channelPayCode/deleteVboxChannelPayCode [delete]
export const deleteChannelPayCodeByIds = (data) => {
  return service({
    url: '/channelPayCode/deleteChannelPayCodeByIds',
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
// @Router /channelPayCode/updateVboxChannelPayCode [put]
export const updateChannelPayCode = (data) => {
  return service({
    url: '/channelPayCode/updateChannelPayCode',
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
// @Router /channelPayCode/findVboxChannelPayCode [get]
export const findChannelPayCode = (params) => {
  return service({
    url: '/channelPayCode/findChannelPayCode',
    method: 'get',
    params
  })
}

// @Tags PayCodeVO
// @Summary 获取预产统计情况
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取预产统计情况"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /channelPayCode/getPayCodeOverview [get]
export const getPayCodeOverviewByChanAcc = (params) => {
  return service({
    url: '/channelPayCode/getPayCodeOverviewByChanAcc',
    method: 'get',
    params
  })
}

// @Tags PayCodeVO
// @Summary 获取预产统计情况
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取预产统计情况"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /channelPayCode/getPayCodeOverview [get]
export const getPayCodeOverview = (params) => {
  return service({
    url: '/channelPayCode/getPayCodeOverview',
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
// @Router /channelPayCode/getVboxChannelPayCodeList [get]
export const getChannelPayCodeList = (params) => {
  return service({
    url: '/channelPayCode/getChannelPayCodeList',
    method: 'get',
    params
  })
}

// @Tags getChannelPayCodeStatisByLocation
// @Summary 分页获取二维码统计排名列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取二维码统计排名列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /channelPayCode/getChannelPayCodeStatisByLocation [get]
export const getChannelPayCodeStatisByLocation = (params) => {
  return service({
    url: '/channelPayCode/getChannelPayCodeStatisByLocation',
    method: 'get',
    params
  })
}

