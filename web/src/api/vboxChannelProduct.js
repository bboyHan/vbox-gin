import service from '@/utils/request'

// @Tags VboxChannelProduct
// @Summary 创建VboxChannelProduct
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelProduct true "创建VboxChannelProduct"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/createVboxChannelProduct [post]
export const createVboxChannelProduct = (data) => {
  return service({
    url: '/vcp/createVboxChannelProduct',
    method: 'post',
    data
  })
}

// @Tags VboxChannelProduct
// @Summary 删除VboxChannelProduct
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelProduct true "删除VboxChannelProduct"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vcp/deleteVboxChannelProduct [delete]
export const deleteVboxChannelProduct = (data) => {
  return service({
    url: '/vcp/deleteVboxChannelProduct',
    method: 'delete',
    data
  })
}

// @Tags VboxChannelProduct
// @Summary 删除VboxChannelProduct
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxChannelProduct"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vcp/deleteVboxChannelProduct [delete]
export const deleteVboxChannelProductByIds = (data) => {
  return service({
    url: '/vcp/deleteVboxChannelProductByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxChannelProduct
// @Summary 更新VboxChannelProduct
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelProduct true "更新VboxChannelProduct"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vcp/updateVboxChannelProduct [put]
export const updateVboxChannelProduct = (data) => {
  return service({
    url: '/vcp/updateVboxChannelProduct',
    method: 'put',
    data
  })
}

// @Tags VboxChannelProduct
// @Summary 用id查询VboxChannelProduct
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxChannelProduct true "用id查询VboxChannelProduct"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vcp/findVboxChannelProduct [get]
export const findVboxChannelProduct = (params) => {
  return service({
    url: '/vcp/findVboxChannelProduct',
    method: 'get',
    params
  })
}

// @Tags VboxChannelProduct
// @Summary 分页获取VboxChannelProduct列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxChannelProduct列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getVboxChannelProductList [get]
export const getVboxChannelProductList = (params) => {
  return service({
    url: '/vcp/getVboxChannelProductList',
    method: 'get',
    params
  })
}

// @Tags VboxChannelProduct
// @Summary 获取VboxChannelProduct列表所有
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxChannelProduct列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getVboxChannelProductList [get]
export const getVboxChannelProductAll = (params) => {
  return service({
    url: '/channelProduct/all',
    method: 'get',
    params
  })
}
