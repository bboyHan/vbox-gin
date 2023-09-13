package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChannelProductRouter struct {
}

// InitVboxChannelProductRouter 初始化 VboxChannelProduct 路由信息
func (s *ChannelProductRouter) InitVboxChannelProductRouter(Router *gin.RouterGroup) {
	vcpRouter := Router.Group("vcp").Use(middleware.OperationRecord())
	vcpRouterWithoutRecord := Router.Group("vcp")
	var vcpApi = v1.ApiGroupApp.Vbox.ChannelProductApi
	{
		vcpRouter.POST("createVboxChannelProduct", vcpApi.CreateVboxChannelProduct)   // 新建VboxChannelProduct
		vcpRouter.DELETE("deleteVboxChannelProduct", vcpApi.DeleteVboxChannelProduct) // 删除VboxChannelProduct
		vcpRouter.PUT("updateVboxChannelProduct", vcpApi.UpdateVboxChannelProduct)    // 更新VboxChannelProduct
	}
	{
		vcpRouterWithoutRecord.GET("findVboxChannelProduct", vcpApi.FindVboxChannelProduct)       // 根据ID获取VboxChannelProduct
		vcpRouterWithoutRecord.GET("getVboxChannelProductList", vcpApi.GetVboxChannelProductList) // 获取VboxChannelProduct列表
		vcpRouterWithoutRecord.GET("getVboxChannelProductAll", vcpApi.GetVboxChannelProductAll)   // 获取VboxChannelProduct列表所有
	}
}
