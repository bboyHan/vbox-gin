import service from '@/utils/request'

// @Tags BdaChorgIndexD
// @Summary 创建通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChorgIndexD true "创建通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bdaChorg/createBdaChorgIndexD [post]
export const createBdaChorgIndexD = (data) => {
  return service({
    url: '/bdaChorg/createBdaChorgIndexD',
    method: 'post',
    data
  })
}

// @Tags BdaChorgIndexD
// @Summary 删除通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChorgIndexD true "删除通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChorg/deleteBdaChorgIndexD [delete]
export const deleteBdaChorgIndexD = (data) => {
  return service({
    url: '/bdaChorg/deleteBdaChorgIndexD',
    method: 'delete',
    data
  })
}

// @Tags BdaChorgIndexD
// @Summary 批量删除通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChorg/deleteBdaChorgIndexD [delete]
export const deleteBdaChorgIndexDByIds = (data) => {
  return service({
    url: '/bdaChorg/deleteBdaChorgIndexDByIds',
    method: 'delete',
    data
  })
}

// @Tags BdaChorgIndexD
// @Summary 更新通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChorgIndexD true "更新通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChorg/updateBdaChorgIndexD [put]
export const updateBdaChorgIndexD = (data) => {
  return service({
    url: '/bdaChorg/updateBdaChorgIndexD',
    method: 'put',
    data
  })
}

// @Tags BdaChorgIndexD
// @Summary 用id查询通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.BdaChorgIndexD true "用id查询通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChorg/findBdaChorgIndexD [get]
export const findBdaChorgIndexD = (params) => {
  return service({
    url: '/bdaChorg/findBdaChorgIndexD',
    method: 'get',
    params
  })
}

// @Tags BdaChorgIndexD
// @Summary 分页获取通道团队统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取通道团队统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChorg/getBdaChorgIndexDList [get]
export const getBdaChorgIndexDList = (params) => {
  return service({
    url: '/bdaChorg/getBdaChorgIndexDList',
    method: 'get',
    params
  })
}
