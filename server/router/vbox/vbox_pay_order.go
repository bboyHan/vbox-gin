package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VboxPayOrderRouter struct {
}

// InitVboxPayOrderRouter 初始化 VboxPayOrder 路由信息
func (s *VboxPayOrderRouter) InitVboxPayOrderRouter(Router *gin.RouterGroup) {
	vpoRouter := Router.Group("vpo").Use(middleware.OperationRecord())
	vpoRouterWithoutRecord := Router.Group("vpo")
	var vpoApi = v1.ApiGroupApp.Vbox.VboxPayOrderApi
	{
		vpoRouter.POST("createVboxPayOrder", vpoApi.CreateVboxPayOrder)             // 新建VboxPayOrder
		vpoRouter.DELETE("deleteVboxPayOrder", vpoApi.DeleteVboxPayOrder)           // 删除VboxPayOrder
		vpoRouter.DELETE("deleteVboxPayOrderByIds", vpoApi.DeleteVboxPayOrderByIds) // 批量删除VboxPayOrder
		vpoRouter.PUT("updateVboxPayOrder", vpoApi.UpdateVboxPayOrder)              // 更新VboxPayOrder
	}
	{
		vpoRouterWithoutRecord.GET("findVboxPayOrder", vpoApi.FindVboxPayOrder)                                                         // 根据ID获取VboxPayOrder
		vpoRouterWithoutRecord.GET("getVboxPayOrderList", vpoApi.GetVboxPayOrderList)                                                   // 获取VboxPayOrder列表
		vpoRouterWithoutRecord.GET("getVboxUserPayOrderAnalysis", vpoApi.GetVboxUserPayOrderAnalysis)                                   // 获取用户订单看板
		vpoRouterWithoutRecord.GET("getSelectUserPayOrderAnalysis", vpoApi.GetSelectUserPayOrderAnalysis)                               // 获取用户订单看板
		vpoRouterWithoutRecord.GET("getSelectPayOrderAnalysisQuantifyCharts", vpoApi.GetSelectPayOrderAnalysisQuantifyCharts)           // 获取用户订单看板
		vpoRouterWithoutRecord.GET("getSelectPayOrderAnalysisChannelIncomeCharts", vpoApi.GetSelectPayOrderAnalysisChannelIncomeCharts) // 获取用户订单看板
		vpoRouterWithoutRecord.GET("getVboxUserPayOrderAnalysisIncomeCharts", vpoApi.GetVboxUserPayOrderAnalysisIncomeCharts)           // 获取用户订单看板
		vpoRouterWithoutRecord.GET("getSelectPayOrderAnalysisIncomeBarCharts", vpoApi.GetSelectPayOrderAnalysisIncomeBarCharts)         // 获取用户订单看板
	}
}
