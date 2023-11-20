import service from '@/utils/request'

// @Tags BdaChShopIndexD
// @Summary 创建用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChShopIndexD true "创建用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bdaChshopIndexD/createBdaChShopIndexD [post]
export const createBdaChShopIndexD = (data) => {
  return service({
    url: '/bdaChshopIndexD/createBdaChShopIndexD',
    method: 'post',
    data
  })
}

// @Tags BdaChShopIndexD
// @Summary 删除用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChShopIndexD true "删除用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChshopIndexD/deleteBdaChShopIndexD [delete]
export const deleteBdaChShopIndexD = (data) => {
  return service({
    url: '/bdaChshopIndexD/deleteBdaChShopIndexD',
    method: 'delete',
    data
  })
}

// @Tags BdaChShopIndexD
// @Summary 批量删除用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChshopIndexD/deleteBdaChShopIndexD [delete]
export const deleteBdaChShopIndexDByIds = (data) => {
  return service({
    url: '/bdaChshopIndexD/deleteBdaChShopIndexDByIds',
    method: 'delete',
    data
  })
}

// @Tags BdaChShopIndexD
// @Summary 更新用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.BdaChShopIndexD true "更新用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChshopIndexD/updateBdaChShopIndexD [put]
export const updateBdaChShopIndexD = (data) => {
  return service({
    url: '/bdaChshopIndexD/updateBdaChShopIndexD',
    method: 'put',
    data
  })
}

// @Tags BdaChShopIndexD
// @Summary 用id查询用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.BdaChShopIndexD true "用id查询用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChshopIndexD/findBdaChShopIndexD [get]
export const findBdaChShopIndexD = (params) => {
  return service({
    url: '/bdaChshopIndexD/findBdaChShopIndexD',
    method: 'get',
    params
  })
}

// @Tags BdaChShopIndexD
// @Summary 分页获取用户通道店铺成率统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取用户通道店铺成率统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChshopIndexD/getBdaChShopIndexDList [get]
export const getBdaChShopIndexDList = (params) => {
  return service({
    url: '/bdaChshopIndexD/getBdaChShopIndexDList',
    method: 'get',
    params
  })
}
