package vbox

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type OrderRouter struct{}

func (s *OrderRouter) InitPubAccessRouter(Router *gin.RouterGroup) {
	orderRouter := Router.Group("order")
	orderApi := v1.ApiGroupApp.Vbox.VboxPayOrderApi
	{
		orderRouter.POST("create", orderApi.CreateOrder2PayAcc)
		orderRouter.POST("query", orderApi.QueryOrder2PayAcc)
		orderRouter.GET("detail", orderApi.QueryOrderSimple)
	}

	cpRouter := Router.Group("channelProduct")
	channelProductApi := v1.ApiGroupApp.Vbox.ChannelProductApi
	{
		cpRouter.GET("all", channelProductApi.GetVboxChannelProductAll)
	}
}
