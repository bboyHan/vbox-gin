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
    url: '/pacc/switchEnable',
    method: 'put',
    data: data
  })
}

// @Tags PayAccount
// @Summary 创建付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PayAccount true "创建付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /pacc/createPayAccount [post]
export const createPayAccount = (data) => {
  return service({
    url: '/pacc/createPayAccount',
    method: 'post',
    data
  })
}

// @Tags PayAccount
// @Summary 删除付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PayAccount true "删除付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pacc/deletePayAccount [delete]
export const deletePayAccount = (data) => {
  return service({
    url: '/pacc/deletePayAccount',
    method: 'delete',
    data
  })
}

// @Tags PayAccount
// @Summary 批量删除付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pacc/deletePayAccount [delete]
export const deletePayAccountByIds = (data) => {
  return service({
    url: '/pacc/deletePayAccountByIds',
    method: 'delete',
    data
  })
}

// @Tags PayAccount
// @Summary 更新付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.PayAccount true "更新付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /pacc/updatePayAccount [put]
export const updatePayAccount = (data) => {
  return service({
    url: '/pacc/updatePayAccount',
    method: 'put',
    data
  })
}

// @Tags PayAccount
// @Summary 用id查询付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.PayAccount true "用id查询付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /pacc/findPayAccount [get]
export const findPayAccount = (params) => {
  return service({
    url: '/pacc/findPayAccount',
    method: 'get',
    params
  })
}

// @Tags PayAccount
// @Summary 获取服务网关
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.PayAccount true "获取服务网关"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /pacc/findPayAccount [get]
export const getPAccGateway = (params) => {
  return service({
    url: '/pacc/getPAccGateway',
    method: 'get',
    params
  })
}

// @Tags PayAccount
// @Summary 分页获取付方列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取付方列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pacc/getPayAccountList [get]
export const getPayAccountList = (params) => {
  return service({
    url: '/pacc/getPayAccountList',
    method: 'get',
    params
  })
}
