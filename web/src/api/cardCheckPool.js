import service from '@/utils/request'

// @Tags CardAcc
// @Summary 开启/关闭通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "设置通道账号开关"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /vca/switchEnableCA [put]
export const switchEnableCA = (data) => {
  return service({
    url: '/cardAcc/switchEnable',
    method: 'put',
    data: data
  })
}

// @Tags CardAcc
// @Summary 开启/关闭通道账号（批量）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "设置通道账号开关"（批量）
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /vca/switchEnableCA [put]
export const switchEnableCAByIds = (data) => {
  return service({
    url: '/cardAcc/switchEnableByIds',
    method: 'put',
    data: data
  })
}

// @Tags CardAcc
// @Summary 创建通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "创建通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /vca/createCardAcc [post]
export const createCardAcc = (data) => {
  return service({
    url: '/cardAcc/createChannelCardAcc',
    method: 'post',
    data
  })
}

// @Tags CardAcc
// @Summary 删除通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "删除通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vca/deleteCardAcc [delete]
export const deleteCardAcc = (data) => {
  return service({
    url: '/cardAcc/deleteChannelCardAcc',
    method: 'delete',
    data
  })
}

// @Tags CardAcc
// @Summary 批量删除通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vca/deleteCardAccByIds [delete]
export const deleteCardAccByIds = (data) => {
  return service({
    url: '/cardAcc/deleteChannelCardAccByIds',
    method: 'delete',
    data
  })
}

// @Tags CardAcc
// @Summary 更新通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "更新通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vca/updateChannelAccount [put]
export const updateCardAcc = (data) => {
  return service({
    url: '/cardAcc/updateChannelCardAcc',
    method: 'put',
    data
  })
}

// @Tags CardAcc
// @Summary 用id查询通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ChannelAccount true "用id查询通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vca/findChannelAccount [get]
export const findCardAcc = (params) => {
  return service({
    url: '/cardAcc/getChannelCardAcc',
    method: 'get',
    params
  })
}

// @Tags ChannelAccount
// @Summary 分页获取通道账号列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取通道账号列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vca/getChannelAccountList [get]
export const getCardAccList = (params) => {
  return service({
    url: '/cardAcc/getChannelCardAccList',
    method: 'get',
    params
  })
}

