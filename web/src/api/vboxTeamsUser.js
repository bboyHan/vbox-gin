import service from '@/utils/request'

// @Tags VboxTeamsUser
// @Summary 创建VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxTeamsUser true "创建VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /tUsers/createVboxTeamsUser [post]
export const createVboxTeamsUser = (data) => {
  return service({
    url: '/tUsers/createVboxTeamsUser',
    method: 'post',
    data
  })
}

// @Tags VboxTeamsUser
// @Summary 删除VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxTeamsUser true "删除VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /tUsers/deleteVboxTeamsUser [delete]
export const deleteVboxTeamsUser = (data) => {
  return service({
    url: '/tUsers/deleteVboxTeamsUser',
    method: 'delete',
    data
  })
}

// @Tags VboxTeamsUser
// @Summary 删除VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /tUsers/deleteVboxTeamsUser [delete]
export const deleteVboxTeamsUserByIds = (data) => {
  return service({
    url: '/tUsers/deleteVboxTeamsUserByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxTeamsUser
// @Summary 更新VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxTeamsUser true "更新VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /tUsers/updateVboxTeamsUser [put]
export const updateVboxTeamsUser = (data) => {
  return service({
    url: '/tUsers/updateVboxTeamsUser',
    method: 'put',
    data
  })
}

// @Tags VboxTeamsUser
// @Summary 用id查询VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxTeamsUser true "用id查询VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /tUsers/findVboxTeamsUser [get]
export const findVboxTeamsUser = (params) => {
  return service({
    url: '/tUsers/findVboxTeamsUser',
    method: 'get',
    params
  })
}

// @Tags VboxTeamsUser
// @Summary 分页获取VboxTeamsUser列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxTeamsUser列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /tUsers/getVboxTeamsUserList [get]
export const getVboxTeamsUserList = (params) => {
  return service({
    url: '/tUsers/getVboxTeamsUserList',
    method: 'get',
    params
  })
}

export const findTeamUserAll = (params) => {
  return service({
    url: '/tUsers/findTeamUserAll',
    method: 'get',
    params
  })
}


export const transferTeamUserApi = (data) => {
  return service({
    url: '/tUsers/transferTeamUser',
    method: 'put',
    data
  })
}
