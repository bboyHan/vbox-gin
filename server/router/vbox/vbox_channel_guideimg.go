package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ChannelGuideImgRouter struct {
}

// InitChannel_guideimgRouter 初始化 Channel_guideimg 路由信息
func (s *ChannelGuideImgRouter) InitChannelGuideImgRouter(Router *gin.RouterGroup) {
	chGuideImgRouter := Router.Group("chGuideImg").Use(middleware.OperationRecord())
	chGuideImgRouterWithoutRecord := Router.Group("chGuideImg")
	var chGuideImgApi = v1.ApiGroupApp.Vbox.ChannelGuideImgApi
	{
		chGuideImgRouter.POST("createChannel_guideimg", chGuideImgApi.CreateChannelGuideImg)             // 新建Channel_guideimg
		chGuideImgRouter.DELETE("deleteChannel_guideimg", chGuideImgApi.DeleteChannelGuideImg)           // 删除Channel_guideimg
		chGuideImgRouter.DELETE("deleteChannel_guideimgByIds", chGuideImgApi.DeleteChannelGuideImgByIds) // 批量删除Channel_guideimg
		chGuideImgRouter.PUT("updateChannel_guideimg", chGuideImgApi.UpdateChannelGuideImg)              // 更新Channel_guideimg
	}
	{
		chGuideImgRouterWithoutRecord.GET("findChannel_guideimg", chGuideImgApi.FindChannelGuideImg)       // 根据ID获取Channel_guideimg
		chGuideImgRouterWithoutRecord.GET("getChannel_guideimgList", chGuideImgApi.GetChannelGuideImgList) // 获取Channel_guideimg列表
	}
}
