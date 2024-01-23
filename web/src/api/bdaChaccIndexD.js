import service from '@/utils/request'

// @Tags BdaChaccIndexD
// @Summary 创建用户通道账户粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChaccIndexD true "创建用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bdaChaccIndexD/createBdaChaccIndexD [post]
export const createBdaChaccIndexD = (data) => {
  return service({
    url: '/bdaChaccIndexD/createBdaChaccIndexD',
    method: 'post',
    data
  })
}

// @Tags BdaChaccIndexD
// @Summary 删除用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChaccIndexD true "删除用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChaccIndexD/deleteBdaChaccIndexD [delete]
export const deleteBdaChaccIndexD = (data) => {
  return service({
    url: '/bdaChaccIndexD/deleteBdaChaccIndexD',
    method: 'delete',
    data
  })
}

// @Tags BdaChaccIndexD
// @Summary 批量删除用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChaccIndexD/deleteBdaChaccIndexD [delete]
export const deleteBdaChaccIndexDByIds = (data) => {
  return service({
    url: '/bdaChaccIndexD/deleteBdaChaccIndexDByIds',
    method: 'delete',
    data
  })
}

// @Tags BdaChaccIndexD
// @Summary 更新用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChaccIndexD true "更新用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChaccIndexD/updateBdaChaccIndexD [put]
export const updateBdaChaccIndexD = (data) => {
  return service({
    url: '/bdaChaccIndexD/updateBdaChaccIndexD',
    method: 'put',
    data
  })
}

// @Tags BdaChaccIndexD
// @Summary 用id查询用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.BdaChaccIndexD true "用id查询用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChaccIndexD/findBdaChaccIndexD [get]
export const findBdaChaccIndexD = (params) => {
  return service({
    url: '/bdaChaccIndexD/findBdaChaccIndexD',
    method: 'get',
    params
  })
}

// @Tags BdaChaccIndexD
// @Summary 分页获取用户通道粒度成率统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取用户通道粒度成率统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChaccIndexD/getBdaChaccIndexDList [get]
export const getBdaChaccIndexDList = (params) => {
  return service({
    url: '/bdaChaccIndexD/getBdaChaccIndexDList',
    method: 'get',
    params
  })
}


// @Tags BdaChaccIndexD
// @Summary 分页获取用户通道粒度成率统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取用户通道粒度成率统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChaccIndexD/getBdaChaccIndexDUesrOverview [get]
export const getBdaChaccIndexDUesrOverview = (params) => {
  return service({
    url: '/bdaChaccIndexD/getBdaChaccIndexDUesrOverview',
    method: 'get',
    params
  })
}
