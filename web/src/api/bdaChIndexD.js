import service from '@/utils/request'

// @Tags BdaChIndexD
// @Summary 创建用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChIndexD true "创建用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bdaChIndexD/createBdaChIndexD [post]
export const createBdaChIndexD = (data) => {
  return service({
    url: '/bdaChIndexD/createBdaChIndexD',
    method: 'post',
    data
  })
}

// @Tags BdaChIndexD
// @Summary 删除用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChIndexD true "删除用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChIndexD/deleteBdaChIndexD [delete]
export const deleteBdaChIndexD = (data) => {
  return service({
    url: '/bdaChIndexD/deleteBdaChIndexD',
    method: 'delete',
    data
  })
}

// @Tags BdaChIndexD
// @Summary 批量删除用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChIndexD/deleteBdaChIndexD [delete]
export const deleteBdaChIndexDByIds = (data) => {
  return service({
    url: '/bdaChIndexD/deleteBdaChIndexDByIds',
    method: 'delete',
    data
  })
}

// @Tags BdaChIndexD
// @Summary 更新用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChIndexD true "更新用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChIndexD/updateBdaChIndexD [put]
export const updateBdaChIndexD = (data) => {
  return service({
    url: '/bdaChIndexD/updateBdaChIndexD',
    method: 'put',
    data
  })
}

// @Tags BdaChIndexD
// @Summary 用id查询用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.BdaChIndexD true "用id查询用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChIndexD/findBdaChIndexD [get]
export const findBdaChIndexD = (params) => {
  return service({
    url: '/bdaChIndexD/findBdaChIndexD',
    method: 'get',
    params
  })
}

// @Tags BdaChIndexD
// @Summary 分页获取用户通道粒度成率统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取用户通道粒度成率统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChIndexD/getBdaChIndexDList [get]
export const getBdaChIndexDList = (params) => {
  return service({
    url: '/bdaChIndexD/getBdaChIndexDList',
    method: 'get',
    params
  })
}
