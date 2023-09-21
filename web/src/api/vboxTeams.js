import service from '@/utils/request'

// @Tags VboxTeams
// @Summary 创建VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxTeams true "创建VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /teams/createVboxTeams [post]
export const createVboxTeams = (data) => {
  return service({
    url: '/teams/createVboxTeams',
    method: 'post',
    data
  })
}

// @Tags VboxTeams
// @Summary 删除VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxTeams true "删除VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /teams/deleteVboxTeams [delete]
export const deleteVboxTeams = (data) => {
  return service({
    url: '/teams/deleteVboxTeams',
    method: 'delete',
    data
  })
}

// @Tags VboxTeams
// @Summary 删除VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /teams/deleteVboxTeams [delete]
export const deleteVboxTeamsByIds = (data) => {
  return service({
    url: '/teams/deleteVboxTeamsByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxTeams
// @Summary 更新VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxTeams true "更新VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /teams/updateVboxTeams [put]
export const updateVboxTeams = (data) => {
  return service({
    url: '/teams/updateVboxTeams',
    method: 'put',
    data
  })
}

// @Tags VboxTeams
// @Summary 用id查询VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxTeams true "用id查询VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /teams/findVboxTeams [get]
export const findVboxTeams = (params) => {
  return service({
    url: '/teams/findVboxTeams',
    method: 'get',
    params
  })
}

// @Tags VboxTeams
// @Summary 分页获取VboxTeams列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxTeams列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /teams/getVboxTeamsList [get]
export const getVboxTeamsList = (params) => {
  return service({
    url: '/teams/getVboxTeamsList',
    method: 'get',
    params
  })
}


