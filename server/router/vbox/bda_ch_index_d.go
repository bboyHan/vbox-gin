package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BdaChIndexDRouter struct {
}

// InitBdaChIndexDRouter 初始化 用户通道粒度成率统计-天更新 路由信息
func (s *BdaChIndexDRouter) InitBdaChIndexDRouter(Router *gin.RouterGroup) {
	bdaChIndexDRouter := Router.Group("bdaChIndexD").Use(middleware.OperationRecord())
	bdaChIndexDRouterWithoutRecord := Router.Group("bdaChIndexD")
	var bdaChIndexDApi = v1.ApiGroupApp.VboxApiGroup.BdaChIndexDApi
	{
		bdaChIndexDRouter.POST("createBdaChIndexD", bdaChIndexDApi.CreateBdaChIndexD)             // 新建用户通道粒度成率统计-天更新
		bdaChIndexDRouter.DELETE("deleteBdaChIndexD", bdaChIndexDApi.DeleteBdaChIndexD)           // 删除用户通道粒度成率统计-天更新
		bdaChIndexDRouter.DELETE("deleteBdaChIndexDByIds", bdaChIndexDApi.DeleteBdaChIndexDByIds) // 批量删除用户通道粒度成率统计-天更新
		bdaChIndexDRouter.PUT("updateBdaChIndexD", bdaChIndexDApi.UpdateBdaChIndexD)              // 更新用户通道粒度成率统计-天更新
	}
	{
		bdaChIndexDRouterWithoutRecord.GET("findBdaChIndexD", bdaChIndexDApi.FindBdaChIndexD)       // 根据ID获取用户通道粒度成率统计-天更新
		bdaChIndexDRouterWithoutRecord.GET("getBdaChIndexDList", bdaChIndexDApi.GetBdaChIndexDList) // 获取用户通道粒度成率统计-天更新列表
	}
}
