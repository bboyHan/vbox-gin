package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChannelCardAccRouter struct {
}

// InitChannelCardAccRouter 初始化 通道账号 路由信息
func (s *ChannelAccountRouter) InitChannelCardAccRouter(Router *gin.RouterGroup) {
	cardAccRouter := Router.Group("cardAcc").Use(middleware.OperationRecord())
	cardAccRouterWithoutRecord := Router.Group("cardAcc")
	var vcaApi = v1.ApiGroupApp.VboxApiGroup.ChannelCardAccApi
	{
		cardAccRouter.POST("createChannelCardAcc", vcaApi.CreateChannelCardAcc)             // 新建通道账号
		cardAccRouter.DELETE("deleteChannelCardAcc", vcaApi.DeleteChannelCardAcc)           // 删除通道账号
		cardAccRouter.DELETE("deleteChannelCardAccByIds", vcaApi.DeleteChannelCardAccByIds) // 批量删除通道账号
		cardAccRouter.PUT("updateChannelCardAcc", vcaApi.UpdateChannelCardAcc)              // 更新通道账号
		cardAccRouter.PUT("switchEnable", vcaApi.SwitchEnableChannelCardAcc)                // 开关VboxPayAccount
		cardAccRouter.PUT("switchEnableByIds", vcaApi.SwitchEnableChannelCardAccByIds)      // 批量开关VboxPayAccount
	}
	{
		cardAccRouterWithoutRecord.GET("getChannelCardAcc", vcaApi.FindChannelCardAcc)        // 根据ID获取通道账号
		cardAccRouterWithoutRecord.GET("getChannelCardAccList", vcaApi.GetChannelCardAccList) // 获取通道账号列表
	}
}
