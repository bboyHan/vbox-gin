package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChannelShopRouter struct {
}

// InitChannelShopRouter 初始化 引导商铺 路由信息
func (s *ChannelShopRouter) InitChannelShopRouter(Router *gin.RouterGroup) {
	channelShopRouter := Router.Group("channelShop").Use(middleware.OperationRecord())
	channelShopRouterWithoutRecord := Router.Group("channelShop")
	var channelShopApi = v1.ApiGroupApp.VboxApiGroup.ChannelShopApi
	{
		channelShopRouter.POST("createChannelShop", channelShopApi.CreateChannelShop)             // 新建引导商铺
		channelShopRouter.DELETE("deleteChannelShop", channelShopApi.DeleteChannelShop)           // 删除引导商铺
		channelShopRouter.DELETE("deleteChannelShopByIds", channelShopApi.DeleteChannelShopByIds) // 批量删除引导商铺
		channelShopRouter.PUT("updateChannelShop", channelShopApi.UpdateChannelShop)              // 更新引导商铺
	}
	{
		channelShopRouterWithoutRecord.GET("findChannelShop", channelShopApi.FindChannelShop)                       // 根据ID获取引导商铺
		channelShopRouterWithoutRecord.GET("findChannelShopByProductId", channelShopApi.FindChannelShopByProductId) // 根据productID获取引导商铺
		channelShopRouterWithoutRecord.GET("getChannelShopList", channelShopApi.GetChannelShopList)                 // 获取引导商铺列表
	}
}
