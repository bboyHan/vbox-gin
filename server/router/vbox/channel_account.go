package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChannelAccountRouter struct {
}

// InitChannelAccountRouter 初始化 通道账号 路由信息
func (s *ChannelAccountRouter) InitChannelAccountRouter(Router *gin.RouterGroup) {
	vcaRouter := Router.Group("vca").Use(middleware.OperationRecord())
	vcaRouterWithoutRecord := Router.Group("vca")
	var vcaApi = v1.ApiGroupApp.VboxApiGroup.ChannelAccountApi
	{
		vcaRouter.POST("createChannelAccount", vcaApi.CreateChannelAccount)             // 新建通道账号
		vcaRouter.DELETE("deleteChannelAccount", vcaApi.DeleteChannelAccount)           // 删除通道账号
		vcaRouter.DELETE("deleteChannelAccountByIds", vcaApi.DeleteChannelAccountByIds) // 批量删除通道账号
		vcaRouter.PUT("updateChannelAccount", vcaApi.UpdateChannelAccount)              // 更新通道账号
		vcaRouter.PUT("switchEnable", vcaApi.SwitchEnableChannelAccount)                // 开关VboxPayAccount
		vcaRouter.PUT("switchEnableByIds", vcaApi.SwitchEnableChannelAccountByIds)      // 开关VboxPayAccount
	}
	{
		vcaRouterWithoutRecord.POST("queryAccOrderHis", vcaApi.QueryAccOrderHis)          // 查询账户所有的通道账号可用个数
		vcaRouterWithoutRecord.GET("countAcc", vcaApi.CountAcc)                           // 查询账户所有的通道账号可用个数
		vcaRouterWithoutRecord.GET("findChannelAccount", vcaApi.FindChannelAccount)       // 根据ID获取通道账号
		vcaRouterWithoutRecord.GET("getChannelAccountList", vcaApi.GetChannelAccountList) // 获取通道账号列表
	}
}
