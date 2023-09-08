import service from '@/utils/request'

// @Tags VboxUserWallet
// @Summary 创建VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxUserWallet true "创建VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vuw/createVboxUserWallet [post]
export const createVboxUserWallet = (data) => {
  return service({
    url: '/vuw/createVboxUserWallet',
    method: 'post',
    data
  })
}

// @Tags VboxUserWallet
// @Summary 删除VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxUserWallet true "删除VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vuw/deleteVboxUserWallet [delete]
export const deleteVboxUserWallet = (data) => {
  return service({
    url: '/vuw/deleteVboxUserWallet',
    method: 'delete',
    data
  })
}

// @Tags VboxUserWallet
// @Summary 删除VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vuw/deleteVboxUserWallet [delete]
export const deleteVboxUserWalletByIds = (data) => {
  return service({
    url: '/vuw/deleteVboxUserWalletByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxUserWallet
// @Summary 更新VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxUserWallet true "更新VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vuw/updateVboxUserWallet [put]
export const updateVboxUserWallet = (data) => {
  return service({
    url: '/vuw/updateVboxUserWallet',
    method: 'put',
    data
  })
}

// @Tags VboxUserWallet
// @Summary 用id查询VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxUserWallet true "用id查询VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vuw/findVboxUserWallet [get]
export const findVboxUserWallet = (params) => {
  return service({
    url: '/vuw/findVboxUserWallet',
    method: 'get',
    params
  })
}

// @Tags VboxUserWallet
// @Summary 分页获取VboxUserWallet列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxUserWallet列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vuw/getVboxUserWalletList [get]
export const getVboxUserWalletList = (params) => {
  return service({
    url: '/vuw/getVboxUserWalletList',
    method: 'get',
    params
  })
}

// @Tags getVboxUserWalletAvailablePoints
// @Summary 获取用户可用余额
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxUserWallet列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vuw/getVboxUserWalletAvailablePoints [get]
export const getVboxUserWalletAvailablePoints = () => {
  return service({
    url: '/vuw/getVboxUserWalletAvailablePoints',
    method: 'get'
  })
}
