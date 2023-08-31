package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProxyRouter struct {
}

// InitVboxProxyRouter 初始化 VboxProxy 路由信息
func (s *ProxyRouter) InitVboxProxyRouter(Router *gin.RouterGroup) {
	vboxProxyRouter := Router.Group("vboxProxy").Use(middleware.OperationRecord())
	vboxProxyRouterWithoutRecord := Router.Group("vboxProxy")
	var vboxProxyApi = v1.ApiGroupApp.Vbox.ProxyApi
	{
		vboxProxyRouter.POST("createVboxProxy", vboxProxyApi.CreateVboxProxy)             // 新建VboxProxy
		vboxProxyRouter.DELETE("deleteVboxProxy", vboxProxyApi.DeleteVboxProxy)           // 删除VboxProxy
		vboxProxyRouter.DELETE("deleteVboxProxyByIds", vboxProxyApi.DeleteVboxProxyByIds) // 批量删除VboxProxy
		vboxProxyRouter.PUT("updateVboxProxy", vboxProxyApi.UpdateVboxProxy)              // 更新VboxProxy
	}
	{
		vboxProxyRouterWithoutRecord.GET("findVboxProxy", vboxProxyApi.FindVboxProxy)       // 根据ID获取VboxProxy
		vboxProxyRouterWithoutRecord.GET("getVboxProxyList", vboxProxyApi.GetVboxProxyList) // 获取VboxProxy列表
	}
}
