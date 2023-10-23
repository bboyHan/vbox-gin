import service from '@/utils/request'

// @Tags ChannelShop
// @Summary 创建引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelShop true "创建引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /channelShop/createChannelShop [post]
export const createChannelShop = (data) => {
  return service({
    url: '/channelShop/createChannelShop',
    method: 'post',
    data
  })
}

// @Tags ChannelShop
// @Summary 删除引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelShop true "删除引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /channelShop/deleteChannelShop [delete]
export const deleteChannelShop = (data) => {
  return service({
    url: '/channelShop/deleteChannelShop',
    method: 'delete',
    data
  })
}

// @Tags ChannelShop
// @Summary 批量删除引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /channelShop/deleteChannelShop [delete]
export const deleteChannelShopByIds = (data) => {
  return service({
    url: '/channelShop/deleteChannelShopByIds',
    method: 'delete',
    data
  })
}

// @Tags ChannelShop
// @Summary 更新引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelShop true "更新引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /channelShop/updateChannelShop [put]
export const updateChannelShop = (data) => {
  return service({
    url: '/channelShop/updateChannelShop',
    method: 'put',
    data
  })
}

// @Tags ChannelShop
// @Summary 用id查询引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ChannelShop true "用id查询引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /channelShop/findChannelShop [get]
export const findChannelShop = (params) => {
  return service({
    url: '/channelShop/findChannelShop',
    method: 'get',
    params
  })
}

// @Tags ChannelShop
// @Summary 分页获取引导商铺列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取引导商铺列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /channelShop/getChannelShopList [get]
export const getChannelShopList = (params) => {
  return service({
    url: '/channelShop/getChannelShopList',
    method: 'get',
    params
  })
}

// @Tags req.ChannelShop
// @Summary 用product_id查询引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ChannelShop true "用id查询引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /channelShop/findChannelShop [get]
export const findChannelShopByProductID = (params) => {
  return service({
    url: '/channelShop/findChannelShopByProductID',
    method: 'get',
    params
  })
}
