import service from '@/utils/request'

// @Tags VboxBdaChIndexD
// @Summary 创建VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxBdaChIndexD true "创建VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChD/createVboxBdaChIndexD [post]
export const createVboxBdaChIndexD = (data) => {
  return service({
    url: '/bdaChD/createVboxBdaChIndexD',
    method: 'post',
    data
  })
}

// @Tags VboxBdaChIndexD
// @Summary 删除VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxBdaChIndexD true "删除VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChD/deleteVboxBdaChIndexD [delete]
export const deleteVboxBdaChIndexD = (data) => {
  return service({
    url: '/bdaChD/deleteVboxBdaChIndexD',
    method: 'delete',
    data
  })
}

// @Tags VboxBdaChIndexD
// @Summary 删除VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChD/deleteVboxBdaChIndexD [delete]
export const deleteVboxBdaChIndexDByIds = (data) => {
  return service({
    url: '/bdaChD/deleteVboxBdaChIndexDByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxBdaChIndexD
// @Summary 更新VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxBdaChIndexD true "更新VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChD/updateVboxBdaChIndexD [put]
export const updateVboxBdaChIndexD = (data) => {
  return service({
    url: '/bdaChD/updateVboxBdaChIndexD',
    method: 'put',
    data
  })
}

// @Tags VboxBdaChIndexD
// @Summary 用id查询VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxBdaChIndexD true "用id查询VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChD/findVboxBdaChIndexD [get]
export const findVboxBdaChIndexD = (params) => {
  return service({
    url: '/bdaChD/findVboxBdaChIndexD',
    method: 'get',
    params
  })
}

// @Tags VboxBdaChIndexD
// @Summary 分页获取VboxBdaChIndexD列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxBdaChIndexD列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChD/getVboxBdaChIndexDList [get]
export const getVboxBdaChIndexDList = (params) => {
  return service({
    url: '/bdaChD/getVboxBdaChIndexDList',
    method: 'get',
    params
  })
}
