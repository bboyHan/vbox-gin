import service from '@/utils/request'

// @Tags VboxProxy
// @Summary 创建VboxProxy
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxProxy true "创建VboxProxy"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vboxProxy/createVboxProxy [post]
export const createVboxProxy = (data) => {
  return service({
    url: '/vboxProxy/createVboxProxy',
    method: 'post',
    data
  })
}

// @Tags VboxProxy
// @Summary 删除VboxProxy
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxProxy true "删除VboxProxy"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vboxProxy/deleteVboxProxy [delete]
export const deleteVboxProxy = (data) => {
  return service({
    url: '/vboxProxy/deleteVboxProxy',
    method: 'delete',
    data
  })
}

// @Tags VboxProxy
// @Summary 删除VboxProxy
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxProxy"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vboxProxy/deleteVboxProxy [delete]
export const deleteVboxProxyByIds = (data) => {
  return service({
    url: '/vboxProxy/deleteVboxProxyByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxProxy
// @Summary 更新VboxProxy
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxProxy true "更新VboxProxy"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vboxProxy/updateVboxProxy [put]
export const updateVboxProxy = (data) => {
  return service({
    url: '/vboxProxy/updateVboxProxy',
    method: 'put',
    data
  })
}

// @Tags VboxProxy
// @Summary 用id查询VboxProxy
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxProxy true "用id查询VboxProxy"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vboxProxy/findVboxProxy [get]
export const findVboxProxy = (params) => {
  return service({
    url: '/vboxProxy/findVboxProxy',
    method: 'get',
    params
  })
}

// @Tags VboxProxy
// @Summary 分页获取VboxProxy列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxProxy列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vboxProxy/getVboxProxyList [get]
export const getVboxProxyList = (params) => {
  return service({
    url: '/vboxProxy/getVboxProxyList',
    method: 'get',
    params
  })
}
