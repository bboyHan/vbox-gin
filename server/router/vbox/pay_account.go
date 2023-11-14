package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PayAccountRouter struct {
}

// InitPayAccountRouter 初始化 付方 路由信息
func (s *PayAccountRouter) InitPayAccountRouter(Router *gin.RouterGroup) {
	paccRouter := Router.Group("pacc").Use(middleware.OperationRecord())
	paccRouterWithoutRecord := Router.Group("pacc")
	var paccApi = v1.ApiGroupApp.VboxApiGroup.PayAccountApi
	{
		paccRouter.POST("createPayAccount", paccApi.CreatePayAccount)             // 新建付方
		paccRouter.DELETE("deletePayAccount", paccApi.DeletePayAccount)           // 删除付方
		paccRouter.DELETE("deletePayAccountByIds", paccApi.DeletePayAccountByIds) // 批量删除付方
		paccRouter.PUT("updatePayAccount", paccApi.UpdatePayAccount)              // 更新付方
		paccRouter.PUT("switchEnable", paccApi.SwitchEnablePayAccount)            // 开关VboxPayAccount
	}
	{
		paccRouterWithoutRecord.GET("getPAccGateway", paccApi.GetPAccGateway)       // 获取付方网关
		paccRouterWithoutRecord.GET("findPayAccount", paccApi.FindPayAccount)       // 根据ID获取付方
		paccRouterWithoutRecord.GET("getPayAccountList", paccApi.GetPayAccountList) // 获取付方列表
	}
}
