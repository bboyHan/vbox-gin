import service from '@/utils/request'

// @Tags PayOrder
// @Summary 补单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PayOrder true "补单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /payOrder/callback2Pa [post]
export const callback2Pa = (data) => {
  return service({
    url: '/payOrder/callback2Pa',
    method: 'post',
    data
  })
}

// @Tags Ip
// @Summary 查询ip归属
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PayOrder true "查询ip归属"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /payOrder/queryIpRegion [get]
export const queryIpRegion = (params) => {
  return service({
    url: '/payOrder/queryIpRegion',
    method: 'get',
    params
  })
}

// @Tags PayOrder
// @Summary 用id查询订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.PayOrder true "用id查询订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /payOrder/findPayOrder [get]
export const findPayOrder = (params) => {
  return service({
    url: '/payOrder/findPayOrder',
    method: 'get',
    params
  })
}

// @Tags PayOrder
// @Summary 分页获取订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取订单列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /payOrder/getPayOrderList [get]
export const getPayOrderList = (params) => {
  return service({
    url: '/payOrder/getPayOrderList',
    method: 'get',
    params
  })
}

// @Tags PayOrder
// @Summary 获取订单acc统计展示数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取订单acc统计展示数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /payOrder/getPayOrderList [get]
export const getOrderAccOverview = (params) => {
  return service({
    url: '/payOrder/getOrderAccOverview',
    method: 'get',
    params
  })
}

// @Tags PayOrder
// @Summary getPayOrderRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "getPayOrderRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /payOrder/getPayOrderRate [get]
export const getPayOrderRate = (params) => {
  return service({
    url: '/payOrder/getPayOrderRate',
    method: 'get',
    params
  })
}

// @Tags PayOrder
// @Summary getPayOrderOverview
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "getPayOrderOverview"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /payOrder/getPayOrderOverview [get]
export const getPayOrderOverview = (params) => {
  return service({
    url: '/payOrder/getPayOrderOverview',
    method: 'get',
    params
  })
}

// @Tags PayOrder
// @Summary 分页获取订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取订单列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /payOrder/getPayOrderListBydt [get]
export const getPayOrderListByDt = (params) => {
  return service({
    url: '/payOrder/getPayOrderListByDt',
    method: 'get',
    params
  })
}


// @Tags PayOrder
// @Summary 分页获取订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取订单列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /payOrder/getPayOrderListLatestHour [get]
export const getPayOrderListLatestHour = (params) => {
  return service({
    url: '/payOrder/getPayOrderListLatestHour',
    method: 'get',
    params
  })
}


// @Tags PayOrder
// @Summary 创建VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxPayOrder true "创建VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/orderTest [post]
export const createOrderTest = (data) => {
  return service({
    url: '/payOrder/orderTest',
    method: 'post',
    data
  })
}

// @Tags PayOrder
// @Summary 用order_id查询VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxPayOrder true "用id查询VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vpo/findVboxPayOrder [get]
export const queryOrderSimple = (params) => {
  return service({
    url: '/order/detail',
    method: 'get',
    params
  })
}

// @Tags PayOrder
// @Summary 客户端回补信息
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxPayOrder true "客户端回补信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vpo/cbExt [get]
export const cbExt = (data) => {
  return service({
    url: '/order/cbExt',
    method: 'post',
    data
  })
}
