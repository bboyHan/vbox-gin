package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VboxUserWalletRouter struct {
}

// InitVboxUserWalletRouter 初始化 VboxUserWallet 路由信息
func (s *VboxUserWalletRouter) InitVboxUserWalletRouter(Router *gin.RouterGroup) {
	vuwRouter := Router.Group("vuw").Use(middleware.OperationRecord())
	vuwRouterWithoutRecord := Router.Group("vuw")
	var vuwApi = v1.ApiGroupApp.VboxApiGroup.VboxUserWalletApi
	{
		vuwRouter.POST("createVboxUserWallet", vuwApi.CreateVboxUserWallet)   // 新建VboxUserWallet
		vuwRouter.DELETE("deleteVboxUserWallet", vuwApi.DeleteVboxUserWallet) // 删除VboxUserWallet
		vuwRouter.DELETE("deleteVboxUserWalletByIds", vuwApi.DeleteVboxUserWalletByIds) // 批量删除VboxUserWallet
		vuwRouter.PUT("updateVboxUserWallet", vuwApi.UpdateVboxUserWallet)    // 更新VboxUserWallet
	}
	{
		vuwRouterWithoutRecord.GET("findVboxUserWallet", vuwApi.FindVboxUserWallet)        // 根据ID获取VboxUserWallet
		vuwRouterWithoutRecord.GET("getVboxUserWalletList", vuwApi.GetVboxUserWalletList)  // 获取VboxUserWallet列表
	}
}
