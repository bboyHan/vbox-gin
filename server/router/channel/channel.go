package channel

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChannelRouter struct {
}

// InitChannelRouter 初始化 Channel 路由信息
func (s *ChannelRouter) InitChannelRouter(Router *gin.RouterGroup) {
	chRouter := Router.Group("ch").Use(middleware.OperationRecord())
	chRouterWithoutRecord := Router.Group("ch")
	var chApi = v1.ApiGroupApp.ChannelApiGroup.ChannelApi
	{
		chRouter.POST("createChannel", chApi.CreateChannel)             // 新建Channel
		chRouter.DELETE("deleteChannel", chApi.DeleteChannel)           // 删除Channel
		chRouter.DELETE("deleteChannelByIds", chApi.DeleteChannelByIds) // 批量删除Channel
		chRouter.PUT("updateChannel", chApi.UpdateChannel)              // 更新Channel
	}
	{
		chRouterWithoutRecord.GET("findChannel", chApi.FindChannel)       // 根据ID获取Channel
		chRouterWithoutRecord.GET("getChannelList", chApi.GetChannelList) // 获取Channel列表
	}
}
