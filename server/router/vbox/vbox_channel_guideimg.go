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
		chGuideImgRouter.POST("createChannelGuideimg", chGuideImgApi.CreateChannelGuideImg)             // 新建Channel_guideimg
		chGuideImgRouter.DELETE("deleteChannelGuideimg", chGuideImgApi.DeleteChannelGuideImg)           // 删除Channel_guideimg
		chGuideImgRouter.DELETE("deleteChannelGuideimgByIds", chGuideImgApi.DeleteChannelGuideImgByIds) // 批量删除Channel_guideimg
		chGuideImgRouter.PUT("updateChannelGuideimg", chGuideImgApi.UpdateChannelGuideImg)              // 更新Channel_guideimg
	}
	{
		chGuideImgRouterWithoutRecord.GET("findChannelGuideimg", chGuideImgApi.FindChannelGuideImg)       // 根据ID获取Channel_guideimg
		chGuideImgRouterWithoutRecord.GET("getChannelGuideimgList", chGuideImgApi.GetChannelGuideImgList) // 获取Channel_guideimg列表
	}
}
