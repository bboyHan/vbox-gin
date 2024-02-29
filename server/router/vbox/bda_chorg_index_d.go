package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BdaChorgIndexDRouter struct {
}

// InitBdaChorgIndexDRouter 初始化 通道团队统计-天更新 路由信息
func (s *BdaChorgIndexDRouter) InitBdaChorgIndexDRouter(Router *gin.RouterGroup) {
	bdaChorgRouter := Router.Group("bdaChorg").Use(middleware.OperationRecord())
	bdaChorgRouterWithoutRecord := Router.Group("bdaChorg")
	var bdaChorgApi = v1.ApiGroupApp.VboxApiGroup.BdaChorgIndexDApi
	{
		bdaChorgRouter.POST("createBdaChorgIndexD", bdaChorgApi.CreateBdaChorgIndexD)             // 新建通道团队统计-天更新
		bdaChorgRouter.DELETE("deleteBdaChorgIndexD", bdaChorgApi.DeleteBdaChorgIndexD)           // 删除通道团队统计-天更新
		bdaChorgRouter.DELETE("deleteBdaChorgIndexDByIds", bdaChorgApi.DeleteBdaChorgIndexDByIds) // 批量删除通道团队统计-天更新
		bdaChorgRouter.PUT("updateBdaChorgIndexD", bdaChorgApi.UpdateBdaChorgIndexD)              // 更新通道团队统计-天更新
	}
	{
		bdaChorgRouterWithoutRecord.GET("findBdaChorgIndexD", bdaChorgApi.FindBdaChorgIndexD)             // 根据ID获取通道团队统计-天更新
		bdaChorgRouterWithoutRecord.GET("getBdaChorgIndexDList", bdaChorgApi.GetBdaChorgIndexDList)       // 获取通道团队统计-天更新列表
		bdaChorgRouterWithoutRecord.GET("getBdaChorgIndexRealList", bdaChorgApi.GetBdaChorgIndexRealList) // 获取通道团队统计-天更新列表
	}
}
