import service from '@/utils/request'

// @Tags ChannelAccount
// @Summary 开启关闭通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "设置通道账号信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/switchEnableCA [put]
export const switchEnableCA = (data) => {
  return service({
    url: '/vca/switchEnable',
    method: 'put',
    data: data
  })
}

// @Tags ChannelAccount
// @Summary 创建ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "创建ChannelAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vca/createChannelAccount [post]
export const createChannelAccount = (data) => {
  return service({
    url: '/vca/createChannelAccount',
    method: 'post',
    data
  })
}

// @Tags ChannelAccount
// @Summary 删除ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "删除ChannelAccount"
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
// @Summary 删除ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除ChannelAccount"
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
// @Summary 更新ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelAccount true "更新ChannelAccount"
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
// @Summary 用id查询ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ChannelAccount true "用id查询ChannelAccount"
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
// @Summary 用id查询ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.ChannelAccount true "用id查询ChannelAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vca/findChannelAccount [get]
export const queryCAHisRecords = (params) => {
  return service({
    url: '/vca/queryCAHisRecords',
    method: 'get',
    params
  })
}

// @Tags ChannelAccount
// @Summary 分页获取ChannelAccount列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取ChannelAccount列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vca/getChannelAccountList [get]
export const getChannelAccountList = (params) => {
  return service({
    url: '/vca/getChannelAccountList',
    method: 'get',
    params
  })
}
