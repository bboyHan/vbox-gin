package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChannelProductRouter struct {
}

// InitChannelProductRouter 初始化 通道产品 路由信息
func (s *ChannelProductRouter) InitChannelProductRouter(Router *gin.RouterGroup) {
	vcpRouter := Router.Group("vcp").Use(middleware.OperationRecord())
	vcpRouterWithoutRecord := Router.Group("vcp")
	var vcpApi = v1.ApiGroupApp.VboxApiGroup.ChannelProductApi
	{
		vcpRouter.POST("createChannelProduct", vcpApi.CreateChannelProduct)             // 新建通道产品
		vcpRouter.DELETE("deleteChannelProduct", vcpApi.DeleteChannelProduct)           // 删除通道产品
		vcpRouter.DELETE("deleteChannelProductByIds", vcpApi.DeleteChannelProductByIds) // 批量删除通道产品
		vcpRouter.PUT("updateChannelProduct", vcpApi.UpdateChannelProduct)              // 更新通道产品
	}
	{
		vcpRouterWithoutRecord.GET("findChannelProduct", vcpApi.FindChannelProduct)       // 根据ID获取通道产品
		vcpRouterWithoutRecord.GET("getChannelProductList", vcpApi.GetChannelProductList) // 获取通道产品列表
		vcpRouterWithoutRecord.GET("getChannelProductAll", vcpApi.GetChannelProductAll)   // 获取通道产品列表
		vcpRouterWithoutRecord.GET("getChannelProductSelf", vcpApi.GetChannelProductSelf) // 获取通道产品列表
	}
}
