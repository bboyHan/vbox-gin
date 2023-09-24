import service from '@/utils/request'

// @Tags VboxChannelRate
// @Summary 创建VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelRate true "创建VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chRate/createLatestVboxChannelRate [post]
export const createLatestVboxChannelRate = (data) => {
  return service({
    url: '/chRate/createLatestVboxChannelRate',
    method: 'post',
    data
  })
}

// @Tags VboxChannelRate
// @Summary 创建VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelRate true "创建VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chRate/createVboxChannelRate [post]
export const createVboxChannelRate = (data) => {
  return service({
    url: '/chRate/createVboxChannelRate',
    method: 'post',
    data
  })
}

// @Tags VboxChannelRate
// @Summary 删除VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelRate true "删除VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chRate/deleteVboxChannelRate [delete]
export const deleteVboxChannelRate = (data) => {
  return service({
    url: '/chRate/deleteVboxChannelRate',
    method: 'delete',
    data
  })
}

// @Tags VboxChannelRate
// @Summary 删除VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chRate/deleteVboxChannelRate [delete]
export const deleteVboxChannelRateByIds = (data) => {
  return service({
    url: '/chRate/deleteVboxChannelRateByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxChannelRate
// @Summary 更新VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxChannelRate true "更新VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /chRate/updateVboxChannelRate [put]
export const updateVboxChannelRate = (data) => {
  return service({
    url: '/chRate/updateVboxChannelRate',
    method: 'put',
    data
  })
}

// @Tags VboxChannelRate
// @Summary 用id查询VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxChannelRate true "用id查询VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /chRate/findVboxChannelRate [get]
export const findVboxChannelRate = (params) => {
  return service({
    url: '/chRate/findVboxChannelRate',
    method: 'get',
    params
  })
}

// @Tags VboxChannelRate
// @Summary 分页获取VboxChannelRate列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxChannelRate列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chRate/getVboxChannelRateList [get]
export const getVboxChannelRateList = (params) => {
  return service({
    url: '/chRate/getVboxChannelRateList',
    method: 'get',
    params
  })
}

// @Tags VboxChannelRate
// @Summary 分页获取VboxChannelRate列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxChannelRate列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chRate/getVboxTeamUserChannelRateList [get]
export const getVboxTeamUserChannelRateList = (params) => {
  return service({
    url: '/chRate/getVboxTeamUserChannelRateList',
    method: 'get',
    params
  })
}
