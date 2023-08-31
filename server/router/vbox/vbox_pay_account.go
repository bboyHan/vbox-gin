package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type PayAccountRouter struct {
}

// InitVboxPayAccountRouter 初始化 VboxPayAccount 路由信息
func (s *PayAccountRouter) InitVboxPayAccountRouter(Router *gin.RouterGroup) {
	vpaRouter := Router.Group("vpa").Use(middleware.OperationRecord())
	vpaRouterWithoutRecord := Router.Group("vpa")
	var vpaApi = v1.ApiGroupApp.Vbox.PayAccountApi
	{
		vpaRouter.POST("createVboxPayAccount", vpaApi.CreateVboxPayAccount)             // 新建VboxPayAccount
		vpaRouter.DELETE("deleteVboxPayAccount", vpaApi.DeleteVboxPayAccount)           // 删除VboxPayAccount
		vpaRouter.DELETE("deleteVboxPayAccountByIds", vpaApi.DeleteVboxPayAccountByIds) // 批量删除VboxPayAccount
		vpaRouter.PUT("updateVboxPayAccount", vpaApi.UpdateVboxPayAccount)              // 更新VboxPayAccount
		vpaRouter.PUT("switchEnable", vpaApi.SwitchEnableVboxPayAccount)                // 开关VboxPayAccount
	}
	{
		vpaRouterWithoutRecord.GET("findVboxPayAccount", vpaApi.FindVboxPayAccount)       // 根据ID获取VboxPayAccount
		vpaRouterWithoutRecord.GET("getVboxPayAccountList", vpaApi.GetVboxPayAccountList) // 获取VboxPayAccount列表
	}
}
