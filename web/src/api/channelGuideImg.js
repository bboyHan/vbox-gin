import service from '@/utils/request'

// @Tags Channel_guideimg
// @Summary 创建Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Channel_guideimg true "创建Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chGuideImg/createChannel_guideimg [post]
export const createChannel_guideimg = (data) => {
  return service({
    url: '/chGuideImg/createChannel_guideimg',
    method: 'post',
    data
  })
}

// @Tags Channel_guideimg
// @Summary 删除Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Channel_guideimg true "删除Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chGuideImg/deleteChannel_guideimg [delete]
export const deleteChannel_guideimg = (data) => {
  return service({
    url: '/chGuideImg/deleteChannel_guideimg',
    method: 'delete',
    data
  })
}

// @Tags Channel_guideimg
// @Summary 删除Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chGuideImg/deleteChannel_guideimg [delete]
export const deleteChannel_guideimgByIds = (data) => {
  return service({
    url: '/chGuideImg/deleteChannel_guideimgByIds',
    method: 'delete',
    data
  })
}

// @Tags Channel_guideimg
// @Summary 更新Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.Channel_guideimg true "更新Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /chGuideImg/updateChannel_guideimg [put]
export const updateChannel_guideimg = (data) => {
  return service({
    url: '/chGuideImg/updateChannel_guideimg',
    method: 'put',
    data
  })
}

// @Tags Channel_guideimg
// @Summary 用id查询Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.Channel_guideimg true "用id查询Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /chGuideImg/findChannel_guideimg [get]
export const findChannel_guideimg = (params) => {
  return service({
    url: '/chGuideImg/findChannel_guideimg',
    method: 'get',
    params
  })
}

// @Tags Channel_guideimg
// @Summary 分页获取Channel_guideimg列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取Channel_guideimg列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chGuideImg/getChannel_guideimgList [get]
export const getChannel_guideimgList = (params) => {
  return service({
    url: '/chGuideImg/getChannel_guideimgList',
    method: 'get',
    params
  })
}