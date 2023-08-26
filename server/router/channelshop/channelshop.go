package channelshop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChannelShopRouter struct {
}

// InitChannelShopRouter 初始化 ChannelShop 路由信息
func (s *ChannelShopRouter) InitChannelShopRouter(Router *gin.RouterGroup) {
	chShopRouter := Router.Group("chShop").Use(middleware.OperationRecord())
	chShopRouterWithoutRecord := Router.Group("chShop")
	var chShopApi = v1.ApiGroupApp.ChannelshopApiGroup.ChannelShopApi
	{
		chShopRouter.POST("createChannelShop", chShopApi.CreateChannelShop)             // 新建ChannelShop
		chShopRouter.DELETE("deleteChannelShop", chShopApi.DeleteChannelShop)           // 删除ChannelShop
		chShopRouter.DELETE("deleteChannelShopByIds", chShopApi.DeleteChannelShopByIds) // 批量删除ChannelShop
		chShopRouter.PUT("updateChannelShop", chShopApi.UpdateChannelShop)              // 更新ChannelShop
	}
	{
		chShopRouterWithoutRecord.GET("findChannelShop", chShopApi.FindChannelShop)                                   // 根据ID获取ChannelShop
		chShopRouterWithoutRecord.GET("getShopMarkList", chShopApi.GetShopMarkList)                                   // 根据ID获取ChannelShop
		chShopRouterWithoutRecord.GET("getChannelShopList", chShopApi.GetChannelShopList)                             // 获取ChannelShop列表
		chShopRouterWithoutRecord.GET("getChannelShopListByChanelRemark", chShopApi.GetChannelShopListByChanelRemark) // 获取同一通道下同一店铺的ChannelShop列表
		chShopRouterWithoutRecord.POST("batchUpdateChannelShopStatus", chShopApi.BatchUpdateChannelShopStatus)        // 批量控制ChannelShop列表状态
		chShopRouterWithoutRecord.POST("batchCreateChannelShop", chShopApi.BatchCreateChannelShop)                    // 批量控制ChannelShop列表状态
	}
}
