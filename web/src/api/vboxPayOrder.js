import service from '@/utils/request'

// @Tags VboxPayOrder
// @Summary 用order_id查询VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxPayOrder true "用id查询VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vpo/findVboxPayOrder [get]
export const queryOrderSimple = (params) => {
  return service({
    url: '/order/detail',
    method: 'get',
    params
  })
}

// @Tags VboxPayOrder
// @Summary 用id查询VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query model.VboxPayOrder true "用id查询VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vpo/findVboxPayOrder [get]
export const findVboxPayOrder = (params) => {
  return service({
    url: '/vpo/findVboxPayOrder',
    method: 'get',
    params
  })
}

// @Tags VboxPayOrder
// @Summary 分页获取VboxPayOrder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "分页获取VboxPayOrder列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getVboxPayOrderList [get]
export const getVboxPayOrderList = (params) => {
  return service({
    url: '/vpo/getVboxPayOrderList',
    method: 'get',
    params
  })
}

// @Tags getSelectUserPayOrderAnalysis
// @Summary 获取单个用户分析展示数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取单个用户分析展示数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getSelectUserPayOrderAnalysis [get]
export const getSelectUserPayOrderAnalysis = (params) => {
  return service({
    url: '/vpo/getSelectUserPayOrderAnalysis',
    method: 'get',
    params
  })
}

// @Tags GetSelectUserPayOrderAnalysisH
// @Summary 获取单个用户分析展示近一两个小时的数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取单个用户分析展示数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/GetSelectUserPayOrderAnalysisH [get]
export const getSelectUserPayOrderAnalysisH = () => {
  return service({
    url: '/vpo/getSelectUserPayOrderAnalysisH',
    method: 'get'
  })
}


// @Tags getHomePagePayOrderAnalysis
// @Summary 获取首页单个用户分析展示数据
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取单个用户分析展示数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getHomePagePayOrderAnalysis [get]
export const getHomePagePayOrderAnalysis = () => {
  return service({
    url: '/vpo/getHomePagePayOrderAnalysis',
    method: 'get'
    
  })
}


// @Tags getSelectPayOrderAnalysisQuantifyCharts
// @Summary 获取单个用户各个通道下每天的成单数据图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取单个用户各个通道下每天的成单数据图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getSelectPayOrderAnalysisQuantifyCharts [get]
export const getSelectPayOrderAnalysisQuantifyCharts = (params) => {
  return service({
    url: '/vpo/getSelectPayOrderAnalysisQuantifyCharts',
    method: 'get',
    params
  })
}


// @Tags getSelectPayOrderAnalysisChannelIncomeCharts
// @Summary 获取单个用户各个通道下每天的成单数据图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取单个用户各个通道下每天的成单数据图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getSelectPayOrderAnalysisChannelIncomeCharts [get]
export const getSelectPayOrderAnalysisChannelIncomeCharts = (params) => {
  return service({
    url: '/vpo/getSelectPayOrderAnalysisChannelIncomeCharts',
    method: 'get',
    params
  })
}



// @Tags getSelectPayOrderAnalysisIncomeBarCharts
// @Summary 获取单个用户每天的收入数据图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取单个用户各个通道下每天的成单数据图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getSelectPayOrderAnalysisIncomeBarCharts [get]
export const getSelectPayOrderAnalysisIncomeBarCharts = (params) => {
  return service({
    url: '/vpo/getSelectPayOrderAnalysisIncomeBarCharts',
    method: 'get',
    params
  })
}




// @Tags getVboxUserPayOrderAnalysis
// @Summary 获取用户订单看板
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取用户订单看板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getVboxUserPayOrderAnalysis [get]
export const getVboxUserPayOrderAnalysis = () => {
  return service({
    url: '/vpo/getVboxUserPayOrderAnalysis',
    method: 'get'
  })
}



// @Tags getVboxUserPayOrderAnalysisIncomeCharts
// @Summary 获取用户订单看板收入图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取用户订单看板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getVboxUserPayOrderAnalysisIncomeCharts [get]
export const getVboxUserPayOrderAnalysisIncomeCharts = () => {
  return service({
    url: '/vpo/getVboxUserPayOrderAnalysisIncomeCharts',
    method: 'get'
  })
}

