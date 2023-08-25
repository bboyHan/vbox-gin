import service from '@/utils/request'

// @Tags Channel
// @Summary 创建Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Channel true "创建Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ch/createChannel [post]
export const createChannel = (data) => {
  return service({
    url: '/ch/createChannel',
    method: 'post',
    data
  })
}

// @Tags Channel
// @Summary 删除Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Channel true "删除Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ch/deleteChannel [delete]
export const deleteChannel = (data) => {
  return service({
    url: '/ch/deleteChannel',
    method: 'delete',
    data
  })
}

// @Tags Channel
// @Summary 删除Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ch/deleteChannel [delete]
export const deleteChannelByIds = (data) => {
  return service({
    url: '/ch/deleteChannelByIds',
    method: 'delete',
    data
  })
}

// @Tags Channel
// @Summary 更新Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Channel true "更新Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ch/updateChannel [put]
export const updateChannel = (data) => {
  return service({
    url: '/ch/updateChannel',
    method: 'put',
    data
  })
}

// @Tags Channel
// @Summary 用id查询Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Channel true "用id查询Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /ch/findChannel [get]
export const findChannel = (params) => {
  return service({
    url: '/ch/findChannel',
    method: 'get',
    params
  })
}

// @Tags Channel
// @Summary 分页获取Channel列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Channel列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ch/getChannelList [get]
export const getChannelList = (params) => {
  return service({
    url: '/ch/getChannelList',
    method: 'get',
    params
  })
}
