import service from '@/utils/request'

// @Tags ChannelProduct
// @Summary 创建通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelProduct true "创建通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /vcp/createChannelProduct [post]
export const createChannelProduct = (data) => {
  return service({
    url: '/vcp/createChannelProduct',
    method: 'post',
    data
  })
}

// @Tags ChannelProduct
// @Summary 删除通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelProduct true "删除通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vcp/deleteChannelProduct [delete]
export const deleteChannelProduct = (data) => {
  return service({
    url: '/vcp/deleteChannelProduct',
    method: 'delete',
    data
  })
}

// @Tags ChannelProduct
// @Summary 批量删除通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vcp/deleteChannelProduct [delete]
export const deleteChannelProductByIds = (data) => {
  return service({
    url: '/vcp/deleteChannelProductByIds',
    method: 'delete',
    data
  })
}

// @Tags ChannelProduct
// @Summary 更新通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelProduct true "更新通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vcp/updateChannelProduct [put]
export const updateChannelProduct = (data) => {
  return service({
    url: '/vcp/updateChannelProduct',
    method: 'put',
    data
  })
}

// @Tags ChannelProduct
// @Summary 用id查询通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ChannelProduct true "用id查询通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vcp/findChannelProduct [get]
export const findChannelProduct = (params) => {
  return service({
    url: '/vcp/findChannelProduct',
    method: 'get',
    params
  })
}

// @Tags ChannelProduct
// @Summary 分页获取通道产品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取通道产品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getChannelProductList [get]
export const getChannelProductList = (params) => {
  return service({
    url: '/vcp/getChannelProductList',
    method: 'get',
    params
  })
}

// @Tags ChannelProduct
// @Summary 获取ChannelProduct所有列表(当前用户组织下的)
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取ChannelProduct所有列表(当前用户组织下的)"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getChannelProductList [get]
export const getChannelProductSelf = (params) => {
  return service({
    url: '/vcp/getChannelProductSelf',
    method: 'get',
    params
  })
}

// @Tags ChannelProduct
// @Summary 获取所有通道产品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取所有通道产品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getChannelProductList [get]
export const getChannelProductAll = (params) => {
  return service({
    url: '/vcp/getChannelProductAll',
    method: 'get',
    params
  })
}
