package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VboxChannelRateRouter struct {
}

// InitVboxChannelRateRouter 初始化 VboxChannelRate 路由信息
func (s *VboxChannelRateRouter) InitVboxChannelRateRouter(Router *gin.RouterGroup) {
	chRateRouter := Router.Group("chRate").Use(middleware.OperationRecord())
	chRateRouterWithoutRecord := Router.Group("chRate")
	var chRateApi = v1.ApiGroupApp.Vbox.VboxChannelRateApi
	{
		chRateRouter.POST("createVboxChannelRate", chRateApi.CreateVboxChannelRate)             // 新建VboxChannelRate
		chRateRouter.DELETE("deleteVboxChannelRate", chRateApi.DeleteVboxChannelRate)           // 删除VboxChannelRate
		chRateRouter.DELETE("deleteVboxChannelRateByIds", chRateApi.DeleteVboxChannelRateByIds) // 批量删除VboxChannelRate
		chRateRouter.PUT("updateVboxChannelRate", chRateApi.UpdateVboxChannelRate)              // 更新VboxChannelRate
	}
	{
		chRateRouterWithoutRecord.GET("findVboxChannelRate", chRateApi.FindVboxChannelRate)                       // 根据ID获取VboxChannelRate
		chRateRouterWithoutRecord.GET("getVboxChannelRateList", chRateApi.GetVboxChannelRateList)                 // 获取VboxChannelRate列表
		chRateRouterWithoutRecord.GET("getVboxTeamUserChannelRateList", chRateApi.GetVboxTeamUserChannelRateList) // 获取VboxChannelRate列表
	}
}
