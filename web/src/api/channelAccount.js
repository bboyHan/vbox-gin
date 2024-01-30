import service from '@/utils/request'

// @Tags ChannelAccount
// @Summary 获取通道账号充值记录
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "获取通道账号充值记录"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /vca/queryAccOrderHis [get]
export const queryAccOrderHis = (data) => {
  return service({
    url: '/vca/queryAccOrderHis',
    method: 'post',
    data: data
  })
}

// @Tags ChannelAccount
// @Summary 获取通道账号可用数
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "获取通道账号可用数"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /vca/countAcc [get]
export const countAcc = (data) => {
  return service({
    url: '/vca/countAcc',
    method: 'get',
    params: data
  })
}

// @Tags ChannelAccount
// @Summary 开启/关闭通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "设置通道账号开关"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /vca/switchEnableCA [put]
export const switchEnableCA = (data) => {
  return service({
    url: '/vca/switchEnable',
    method: 'put',
    data: data
  })
}

// @Tags ChannelAccount
// @Summary 开启/关闭通道账号（批量）
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "设置通道账号开关"（批量）
// @Success 200 {string} string "{"success":true,"data":{},"msg":"操作成功"}"
// @Router /vca/switchEnableCA [put]
export const switchEnableCAByIds = (data) => {
  return service({
    url: '/vca/switchEnableByIds',
    method: 'put',
    data: data
  })
}

// @Tags ChannelAccount
// @Summary 创建通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "创建通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /vca/createChannelAccount [post]
export const createChannelAccount = (data) => {
  return service({
    url: '/vca/createChannelAccount',
    method: 'post',
    data
  })
}

// @Tags ChannelAccount
// @Summary 删除通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "删除通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vca/deleteChannelAccount [delete]
export const deleteChannelAccount = (data) => {
  return service({
    url: '/vca/deleteChannelAccount',
    method: 'delete',
    data
  })
}

// @Tags ChannelAccount
// @Summary 批量删除通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vca/deleteChannelAccount [delete]
export const deleteChannelAccountByIds = (data) => {
  return service({
    url: '/vca/deleteChannelAccountByIds',
    method: 'delete',
    data
  })
}

// @Tags ChannelAccount
// @Summary 更新通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "更新通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vca/updateChannelAccount [put]
export const updateChannelAccount = (data) => {
  return service({
    url: '/vca/updateChannelAccount',
    method: 'put',
    data
  })
}

// @Tags ChannelAccount
// @Summary 通道转移
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "通道转移"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vca/updateChannelAccount [put]
export const transferChannelForAcc = (data) => {
  return service({
    url: '/vca/transferChannelForAcc',
    method: 'post',
    data
  })
}

// @Tags ChannelAccount
// @Summary 用id查询通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ChannelAccount true "用id查询通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vca/findChannelAccount [get]
export const findChannelAccount = (params) => {
  return service({
    url: '/vca/findChannelAccount',
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
export const getChannelAccountList = (params) => {
  return service({
    url: '/vca/getChannelAccountList',
    method: 'get',
    params
  })
}

