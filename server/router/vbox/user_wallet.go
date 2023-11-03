package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type UserWalletRouter struct {
}

// InitUserWalletRouter 初始化 用户钱包 路由信息
func (s *UserWalletRouter) InitUserWalletRouter(Router *gin.RouterGroup) {
	userWalletRouter := Router.Group("userWallet").Use(middleware.OperationRecord())
	userWalletRouterWithoutRecord := Router.Group("userWallet")
	var userWalletApi = v1.ApiGroupApp.VboxApiGroup.UserWalletApi
	{
		userWalletRouter.POST("transfer", userWalletApi.TransferUserWallet)                   // 划转、充值
		userWalletRouter.POST("createUserWallet", userWalletApi.CreateUserWallet)             // 新建用户钱包
		userWalletRouter.DELETE("deleteUserWallet", userWalletApi.DeleteUserWallet)           // 删除用户钱包
		userWalletRouter.DELETE("deleteUserWalletByIds", userWalletApi.DeleteUserWalletByIds) // 批量删除用户钱包
		userWalletRouter.PUT("updateUserWallet", userWalletApi.UpdateUserWallet)              // 更新用户钱包
	}
	{
		userWalletRouterWithoutRecord.GET("findUserWallet", userWalletApi.FindUserWallet)       // 根据ID获取用户钱包
		userWalletRouterWithoutRecord.GET("getUserWalletSelf", userWalletApi.GetUserWalletSelf) // 根据ID获取用户钱包
		userWalletRouterWithoutRecord.GET("getUserWalletList", userWalletApi.GetUserWalletList) // 获取用户钱包列表
	}
}
