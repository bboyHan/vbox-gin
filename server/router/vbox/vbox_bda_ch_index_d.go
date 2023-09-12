package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VboxBdaChIndexDRouter struct {
}

// InitVboxBdaChIndexDRouter 初始化 VboxBdaChIndexD 路由信息
func (s *VboxBdaChIndexDRouter) InitVboxBdaChIndexDRouter(Router *gin.RouterGroup) {
	bdaChDRouter := Router.Group("bdaChD").Use(middleware.OperationRecord())
	bdaChDRouterWithoutRecord := Router.Group("bdaChD")
	var bdaChDApi = v1.ApiGroupApp.Vbox.VboxBdaChIndexDApi
	{
		bdaChDRouter.POST("createVboxBdaChIndexD", bdaChDApi.CreateVboxBdaChIndexD)             // 新建VboxBdaChIndexD
		bdaChDRouter.DELETE("deleteVboxBdaChIndexD", bdaChDApi.DeleteVboxBdaChIndexD)           // 删除VboxBdaChIndexD
		bdaChDRouter.DELETE("deleteVboxBdaChIndexDByIds", bdaChDApi.DeleteVboxBdaChIndexDByIds) // 批量删除VboxBdaChIndexD
		bdaChDRouter.PUT("updateVboxBdaChIndexD", bdaChDApi.UpdateVboxBdaChIndexD)              // 更新VboxBdaChIndexD
	}
	{
		bdaChDRouterWithoutRecord.GET("findVboxBdaChIndexD", bdaChDApi.FindVboxBdaChIndexD)       // 根据ID获取VboxBdaChIndexD
		bdaChDRouterWithoutRecord.GET("getVboxBdaChIndexDList", bdaChDApi.GetVboxBdaChIndexDList) // 获取VboxBdaChIndexD列表
	}
}
