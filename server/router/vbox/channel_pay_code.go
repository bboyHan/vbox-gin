package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VboxChannelPayCodeRouter struct {
}

// InitChannelPayCodeRouter 初始化 通道账户付款二维码 路由信息
func (s *VboxChannelPayCodeRouter) InitChannelPayCodeRouter(Router *gin.RouterGroup) {
	vboxChannelPayCodeRouter := Router.Group("channelPayCode").Use(middleware.OperationRecord())
	vboxChannelPayCodeRouterWithoutRecord := Router.Group("channelPayCode")
	var vboxChannelPayCodeApi = v1.ApiGroupApp.VboxApiGroup.ChannelPayCodeApi
	{
		vboxChannelPayCodeRouter.POST("createChannelPayCode", vboxChannelPayCodeApi.CreateChannelPayCode)             // 新建通道账户付款二维码
		vboxChannelPayCodeRouter.DELETE("deleteChannelPayCode", vboxChannelPayCodeApi.DeleteChannelPayCode)           // 删除通道账户付款二维码
		vboxChannelPayCodeRouter.DELETE("deleteChannelPayCodeByIds", vboxChannelPayCodeApi.DeleteChannelPayCodeByIds) // 批量删除通道账户付款二维码
		vboxChannelPayCodeRouter.PUT("updateChannelPayCode", vboxChannelPayCodeApi.UpdateVboxChannelPayCode)          // 更新通道账户付款二维码
	}
	{
		vboxChannelPayCodeRouterWithoutRecord.GET("getPayCodeOverview", vboxChannelPayCodeApi.GetPayCodeOverview)                               // 
		vboxChannelPayCodeRouterWithoutRecord.GET("getPayCodeOverviewByChanAcc", vboxChannelPayCodeApi.GetPayCodeOverviewByChanAcc)             //
		vboxChannelPayCodeRouterWithoutRecord.GET("findChannelPayCode", vboxChannelPayCodeApi.FindChannelPayCode)                               // 根据ID获取通道账户付款二维码
		vboxChannelPayCodeRouterWithoutRecord.GET("getChannelPayCodeList", vboxChannelPayCodeApi.GetChannelPayCodeList)                         // 获取通道账户付款二维码列表
		vboxChannelPayCodeRouterWithoutRecord.GET("getChannelPayCodeStatisByLocation", vboxChannelPayCodeApi.GetChannelPayCodeStatisByLocation) // 分页获取二维码统计排名列表
	}
}
