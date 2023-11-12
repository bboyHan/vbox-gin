package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VboxChannelPayCodeRouter struct {
}

// InitVboxChannelPayCodeRouter 初始化 通道账户付款二维码 路由信息
func (s *VboxChannelPayCodeRouter) InitVboxChannelPayCodeRouter(Router *gin.RouterGroup) {
	vboxChannelPayCodeRouter := Router.Group("vboxChannelPayCode").Use(middleware.OperationRecord())
	vboxChannelPayCodeRouterWithoutRecord := Router.Group("vboxChannelPayCode")
	var vboxChannelPayCodeApi = v1.ApiGroupApp.VboxApiGroup.VboxChannelPayCodeApi
	{
		vboxChannelPayCodeRouter.POST("createVboxChannelPayCode", vboxChannelPayCodeApi.CreateVboxChannelPayCode)             // 新建通道账户付款二维码
		vboxChannelPayCodeRouter.DELETE("deleteVboxChannelPayCode", vboxChannelPayCodeApi.DeleteVboxChannelPayCode)           // 删除通道账户付款二维码
		vboxChannelPayCodeRouter.DELETE("deleteVboxChannelPayCodeByIds", vboxChannelPayCodeApi.DeleteVboxChannelPayCodeByIds) // 批量删除通道账户付款二维码
		vboxChannelPayCodeRouter.PUT("updateVboxChannelPayCode", vboxChannelPayCodeApi.UpdateVboxChannelPayCode)              // 更新通道账户付款二维码
	}
	{
		vboxChannelPayCodeRouterWithoutRecord.GET("findVboxChannelPayCode", vboxChannelPayCodeApi.FindVboxChannelPayCode)       // 根据ID获取通道账户付款二维码
		vboxChannelPayCodeRouterWithoutRecord.GET("getVboxChannelPayCodeList", vboxChannelPayCodeApi.GetVboxChannelPayCodeList) // 获取通道账户付款二维码列表
	}
}
