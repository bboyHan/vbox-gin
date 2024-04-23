package vbox

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/gin-gonic/gin"
)

type OrderRouter struct{}

func (s *OrderRouter) InitPubAccessRouter(Router *gin.RouterGroup) {
	orderRouter := Router.Group("order")
	orderApi := v1.ApiGroupApp.VboxApiGroup.PayOrderApi

	{
		orderRouter.POST("create", orderApi.CreateOrder2PayAcc)
		orderRouter.POST("query", orderApi.QueryOrder2PayAcc)
		orderRouter.GET("detail", orderApi.QueryOrderSimple)
		orderRouter.GET("qry", orderApi.QueryOrderNoAuth)
		orderRouter.POST("cbExt", orderApi.CallbackOrderExt)
		orderRouter.POST("cb", orderApi.CallbackTestSimple)
	}

	vcaRouter := Router.Group("vca")
	var vcaApi = v1.ApiGroupApp.VboxApiGroup.ChannelAccountApi

	{
		vcaRouter.GET("queryOrgAccAvailable", vcaApi.QueryOrgAccAvailable) // 查询账户所有的通道账号可用情况
	}

}
