import service from '@/utils/request'

// @Tags VboxPayOrder
// @Summary 创建VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxPayOrder true "创建VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/createVboxPayOrder [post]
export const createVboxPayOrder = (data) => {
  return service({
    url: '/vpo/createVboxPayOrder',
    method: 'post',
    data
  })
}

// @Tags VboxPayOrder
// @Summary 删除VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxPayOrder true "删除VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vpo/deleteVboxPayOrder [delete]
export const deleteVboxPayOrder = (data) => {
  return service({
    url: '/vpo/deleteVboxPayOrder',
    method: 'delete',
    data
  })
}

// @Tags VboxPayOrder
// @Summary 删除VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vpo/deleteVboxPayOrder [delete]
export const deleteVboxPayOrderByIds = (data) => {
  return service({
    url: '/vpo/deleteVboxPayOrderByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxPayOrder
// @Summary 更新VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxPayOrder true "更新VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vpo/updateVboxPayOrder [put]
export const updateVboxPayOrder = (data) => {
  return service({
    url: '/vpo/updateVboxPayOrder',
    method: 'put',
    data
  })
}

// @Tags VboxPayOrder
// @Summary 用id查询VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxPayOrder true "用id查询VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vpo/findVboxPayOrder [get]
export const findVboxPayOrder = (params) => {
  return service({
    url: '/vpo/findVboxPayOrder',
    method: 'get',
    params
  })
}

// @Tags VboxPayOrder
// @Summary 分页获取VboxPayOrder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxPayOrder列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getVboxPayOrderList [get]
export const getVboxPayOrderList = (params) => {
  return service({
    url: '/vpo/getVboxPayOrderList',
    method: 'get',
    params
  })
}
