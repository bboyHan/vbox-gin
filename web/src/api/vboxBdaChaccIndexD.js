import service from '@/utils/request'

// @Tags VboxBdaChaccIndexD
// @Summary 创建VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxBdaChaccIndexD true "创建VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChaccD/createVboxBdaChaccIndexD [post]
export const createVboxBdaChaccIndexD = (data) => {
  return service({
    url: '/bdaChaccD/createVboxBdaChaccIndexD',
    method: 'post',
    data
  })
}

// @Tags VboxBdaChaccIndexD
// @Summary 删除VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxBdaChaccIndexD true "删除VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChaccD/deleteVboxBdaChaccIndexD [delete]
export const deleteVboxBdaChaccIndexD = (data) => {
  return service({
    url: '/bdaChaccD/deleteVboxBdaChaccIndexD',
    method: 'delete',
    data
  })
}

// @Tags VboxBdaChaccIndexD
// @Summary 删除VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChaccD/deleteVboxBdaChaccIndexD [delete]
export const deleteVboxBdaChaccIndexDByIds = (data) => {
  return service({
    url: '/bdaChaccD/deleteVboxBdaChaccIndexDByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxBdaChaccIndexD
// @Summary 更新VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxBdaChaccIndexD true "更新VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChaccD/updateVboxBdaChaccIndexD [put]
export const updateVboxBdaChaccIndexD = (data) => {
  return service({
    url: '/bdaChaccD/updateVboxBdaChaccIndexD',
    method: 'put',
    data
  })
}

// @Tags VboxBdaChaccIndexD
// @Summary 用id查询VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxBdaChaccIndexD true "用id查询VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChaccD/findVboxBdaChaccIndexD [get]
export const findVboxBdaChaccIndexD = (params) => {
  return service({
    url: '/bdaChaccD/findVboxBdaChaccIndexD',
    method: 'get',
    params
  })
}

// @Tags VboxBdaChaccIndexD
// @Summary 分页获取VboxBdaChaccIndexD列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxBdaChaccIndexD列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChaccD/getVboxBdaChaccIndexDList [get]
export const getVboxBdaChaccIndexDList = (params) => {
  return service({
    url: '/bdaChaccD/getVboxBdaChaccIndexDList',
    method: 'get',
    params
  })
}
