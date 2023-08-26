import service from '@/utils/request'

// @Tags ChannelShop
// @Summary 创建ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelShop true "创建ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chShop/createChannelShop [post]
export const createChannelShop = (data) => {
  return service({
    url: '/chShop/createChannelShop',
    method: 'post',
    data
  })
}


// @Tags ChannelShop
// @Summary 批量创建ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelShop true "批量创建ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chShop/batchCreateChannelShop [post]
export const batchCreateChannelShop = (data) => {
  return service({
    url: '/chShop/batchCreateChannelShop',
    method: 'post',
    data
  })
}


// @Tags ChannelShop
// @Summary 批量更新ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelShop true "批量更新ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chShop/batchUpdateChannelShopStatus [post]
export const batchUpdateChannelShopStatus = (data) => {
  return service({
    url: '/chShop/batchUpdateChannelShopStatus',
    method: 'post',
    data
  })
}



// @Tags ChannelShop
// @Summary 删除ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelShop true "删除ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chShop/deleteChannelShop [delete]
export const deleteChannelShop = (data) => {
  return service({
    url: '/chShop/deleteChannelShop',
    method: 'delete',
    data
  })
}

// @Tags ChannelShop
// @Summary 删除ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chShop/deleteChannelShop [delete]
export const deleteChannelShopByIds = (data) => {
  return service({
    url: '/chShop/deleteChannelShopByIds',
    method: 'delete',
    data
  })
}

// @Tags ChannelShop
// @Summary 更新ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelShop true "更新ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /chShop/updateChannelShop [put]
export const updateChannelShop = (data) => {
  return service({
    url: '/chShop/updateChannelShop',
    method: 'put',
    data
  })
}

// @Tags ChannelShop
// @Summary 用id查询ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ChannelShop true "用id查询ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /chShop/findChannelShop [get]
export const findChannelShop = (params) => {
  return service({
    url: '/chShop/findChannelShop',
    method: 'get',
    params
  })
}

// @Tags ChannelShop
// @Summary 用id查询ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ChannelShop true "用id查询ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /chShop/findChannelShop [get]
export const getShopMarkList = () => {
  return service({
    url: '/chShop/getShopMarkList',
    method: 'get'
  })
}


// @Tags ChannelShop
// @Summary 分页获取ChannelShop列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取ChannelShop列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chShop/getChannelShopList [get]
export const getChannelShopList = (params) => {
  return service({
    url: '/chShop/getChannelShopList',
    method: 'get',
    params
  })
}


// @Tags ChannelShop
// @Summary 分页获取ChannelShop列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取ChannelShop列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chShop/getChannelShopListByChanelRemark [get]
export const getChannelShopListByChanelRemark = (params) => {
  return service({
    url: '/chShop/getChannelShopListByChanelRemark',
    method: 'get',
    params
  })
}
