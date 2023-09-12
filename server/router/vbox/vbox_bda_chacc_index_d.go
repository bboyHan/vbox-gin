package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VboxBdaChaccIndexDRouter struct {
}

// InitVboxBdaChaccIndexDRouter 初始化 VboxBdaChaccIndexD 路由信息
func (s *VboxBdaChaccIndexDRouter) InitVboxBdaChaccIndexDRouter(Router *gin.RouterGroup) {
	bdaChaccDRouter := Router.Group("bdaChaccD").Use(middleware.OperationRecord())
	bdaChaccDRouterWithoutRecord := Router.Group("bdaChaccD")
	var bdaChaccDApi = v1.ApiGroupApp.Vbox.VboxBdaChaccIndexDApi
	{
		bdaChaccDRouter.POST("createVboxBdaChaccIndexD", bdaChaccDApi.CreateVboxBdaChaccIndexD)             // 新建VboxBdaChaccIndexD
		bdaChaccDRouter.DELETE("deleteVboxBdaChaccIndexD", bdaChaccDApi.DeleteVboxBdaChaccIndexD)           // 删除VboxBdaChaccIndexD
		bdaChaccDRouter.DELETE("deleteVboxBdaChaccIndexDByIds", bdaChaccDApi.DeleteVboxBdaChaccIndexDByIds) // 批量删除VboxBdaChaccIndexD
		bdaChaccDRouter.PUT("updateVboxBdaChaccIndexD", bdaChaccDApi.UpdateVboxBdaChaccIndexD)              // 更新VboxBdaChaccIndexD
	}
	{
		bdaChaccDRouterWithoutRecord.GET("findVboxBdaChaccIndexD", bdaChaccDApi.FindVboxBdaChaccIndexD)       // 根据ID获取VboxBdaChaccIndexD
		bdaChaccDRouterWithoutRecord.GET("getVboxBdaChaccIndexDList", bdaChaccDApi.GetVboxBdaChaccIndexDList) // 获取VboxBdaChaccIndexD列表
	}
}
