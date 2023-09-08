package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChannelAccountRouter struct {
}

// InitChannelAccountRouter 初始化 ChannelAccount 路由信息
func (s *ChannelAccountRouter) InitChannelAccountRouter(Router *gin.RouterGroup) {
	vcaRouter := Router.Group("vca").Use(middleware.OperationRecord())
	vcaRouterWithoutRecord := Router.Group("vca")
	var vcaApi = v1.ApiGroupApp.Vbox.ChannelAccountApi
	{
		vcaRouter.POST("createChannelAccount", vcaApi.CreateChannelAccount)             // 新建ChannelAccount
		vcaRouter.DELETE("deleteChannelAccount", vcaApi.DeleteChannelAccount)           // 删除ChannelAccount
		vcaRouter.DELETE("deleteChannelAccountByIds", vcaApi.DeleteChannelAccountByIds) // 批量删除ChannelAccount
		vcaRouter.PUT("updateChannelAccount", vcaApi.UpdateChannelAccount)              // 更新ChannelAccount
		vcaRouter.PUT("switchEnable", vcaApi.SwitchEnableChannelAccount)                // 开关VboxPayAccount
	}
	{
		vcaRouterWithoutRecord.GET("findChannelAccount", vcaApi.FindChannelAccount)       // 根据ID获取ChannelAccount
		vcaRouterWithoutRecord.GET("queryCAHisRecords", vcaApi.QueryCAHisRecords)         // 根据ID获取ChannelAccount
		vcaRouterWithoutRecord.GET("getChannelAccountList", vcaApi.GetChannelAccountList) // 获取ChannelAccount列表
	}
}
