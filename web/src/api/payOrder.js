import service from '@/utils/request'

// @Tags PayOrder
// @Summary 创建订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PayOrder true "创建订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /payOrder/createPayOrder [post]
export const createPayOrder = (data) => {
  return service({
    url: '/payOrder/createPayOrder',
    method: 'post',
    data
  })
}

// @Tags PayOrder
// @Summary 删除订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PayOrder true "删除订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /payOrder/deletePayOrder [delete]
export const deletePayOrder = (data) => {
  return service({
    url: '/payOrder/deletePayOrder',
    method: 'delete',
    data
  })
}

// @Tags PayOrder
// @Summary 批量删除订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /payOrder/deletePayOrder [delete]
export const deletePayOrderByIds = (data) => {
  return service({
    url: '/payOrder/deletePayOrderByIds',
    method: 'delete',
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
