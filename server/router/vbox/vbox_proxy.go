package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type ProxyRouter struct {
}

// InitVboxProxyRouter 初始化 信道 路由信息
func (s *ProxyRouter) InitVboxProxyRouter(Router *gin.RouterGroup) {
	vboxProxyRouter := Router.Group("vboxProxy").Use(middleware.OperationRecord())
	vboxProxyRouterWithoutRecord := Router.Group("vboxProxy")
	var vboxProxyApi = v1.ApiGroupApp.VboxApiGroup.VboxProxyApi
	{
		vboxProxyRouter.POST("createVboxProxy", vboxProxyApi.CreateVboxProxy)             // 新建信道
		vboxProxyRouter.DELETE("deleteVboxProxy", vboxProxyApi.DeleteVboxProxy)           // 删除信道
		vboxProxyRouter.DELETE("deleteVboxProxyByIds", vboxProxyApi.DeleteVboxProxyByIds) // 批量删除信道
		vboxProxyRouter.PUT("updateVboxProxy", vboxProxyApi.UpdateVboxProxy)              // 更新信道
	}
	{
		vboxProxyRouterWithoutRecord.GET("findVboxProxy", vboxProxyApi.FindVboxProxy)       // 根据ID获取信道
		vboxProxyRouterWithoutRecord.GET("getVboxProxyList", vboxProxyApi.GetVboxProxyList) // 获取信道列表
	}
}
