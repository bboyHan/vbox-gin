import service from '@/utils/request'

// @Tags PAccount
// @Summary 开启关闭付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.SysUser true "设置付方信息"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"修改成功"}"
// @Router /user/switchEnablePA [put]
export const switchEnablePA = (data) => {
  return service({
    url: '/vpa/switchEnable',
    method: 'put',
    data: data
  })
}

// @Tags VboxPayAccount
// @Summary 创建VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxPayAccount true "创建VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpa/createVboxPayAccount [post]
export const createVboxPayAccount = (data) => {
  return service({
    url: '/vpa/createVboxPayAccount',
    method: 'post',
    data
  })
}

// @Tags VboxPayAccount
// @Summary 删除VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxPayAccount true "删除VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vpa/deleteVboxPayAccount [delete]
export const deleteVboxPayAccount = (data) => {
  return service({
    url: '/vpa/deleteVboxPayAccount',
    method: 'delete',
    data
  })
}

// @Tags VboxPayAccount
// @Summary 删除VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vpa/deleteVboxPayAccount [delete]
export const deleteVboxPayAccountByIds = (data) => {
  return service({
    url: '/vpa/deleteVboxPayAccountByIds',
    method: 'delete',
    data
  })
}

// @Tags VboxPayAccount
// @Summary 更新VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.VboxPayAccount true "更新VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vpa/updateVboxPayAccount [put]
export const updateVboxPayAccount = (data) => {
  return service({
    url: '/vpa/updateVboxPayAccount',
    method: 'put',
    data
  })
}

// @Tags VboxPayAccount
// @Summary 用id查询VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxPayAccount true "用id查询VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vpa/findVboxPayAccount [get]
export const findVboxPayAccount = (params) => {
  return service({
    url: '/vpa/findVboxPayAccount',
    method: 'get',
    params
  })
}

// @Tags VboxPayAccount
// @Summary 分页获取VboxPayAccount列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxPayAccount列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpa/getVboxPayAccountList [get]
export const getVboxPayAccountList = (params) => {
  return service({
    url: '/vpa/getVboxPayAccountList',
    method: 'get',
    params
  })
}
