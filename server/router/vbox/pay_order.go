package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type PayOrderRouter struct {
}

// InitPayOrderRouter 初始化 订单 路由信息
func (s *PayOrderRouter) InitPayOrderRouter(Router *gin.RouterGroup) {
	payOrderRouterWithoutRecord := Router.Group("payOrder")
	var payOrderApi = v1.ApiGroupApp.VboxApiGroup.PayOrderApi
	{
		payOrderRouterWithoutRecord.GET("findPayOrder", payOrderApi.FindPayOrder)                           // 根据ID获取订单
		payOrderRouterWithoutRecord.GET("getPayOrderList", payOrderApi.GetPayOrderList)                     // 获取订单列表
		payOrderRouterWithoutRecord.GET("getPayOrderOverview", payOrderApi.GetPayOrderOverview)             // 获取订单统计展示数据
		payOrderRouterWithoutRecord.GET("getPayOrderListByDt", payOrderApi.GetPayOrderListByDt)             // 获取某天订单列表
		payOrderRouterWithoutRecord.GET("getPayOrderListLatestHour", payOrderApi.GetPayOrderListLatestHour) // 获取最近一小时订单列表
		payOrderRouterWithoutRecord.GET("queryIpRegion", payOrderApi.QueryIpRegion)                         // 查询IP区域
		payOrderRouterWithoutRecord.POST("orderTest", payOrderApi.CreateOrderTest)                          // 测试订单
	}
}
