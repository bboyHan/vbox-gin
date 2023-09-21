package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VboxTeamsUserRouter struct {
}

// InitVboxTeamsUserRouter 初始化 VboxTeamsUser 路由信息
func (s *VboxTeamsUserRouter) InitVboxTeamsUserRouter(Router *gin.RouterGroup) {
	tUsersRouter := Router.Group("tUsers").Use(middleware.OperationRecord())
	tUsersRouterWithoutRecord := Router.Group("tUsers")
	var tUsersApi = v1.ApiGroupApp.Vbox.VboxTeamsUserApi
	{
		tUsersRouter.POST("createVboxTeamsUser", tUsersApi.CreateVboxTeamsUser)             // 新建VboxTeamsUser
		tUsersRouter.DELETE("deleteVboxTeamsUser", tUsersApi.DeleteVboxTeamsUser)           // 删除VboxTeamsUser
		tUsersRouter.DELETE("deleteVboxTeamsUserByIds", tUsersApi.DeleteVboxTeamsUserByIds) // 批量删除VboxTeamsUser
		tUsersRouter.PUT("updateVboxTeamsUser", tUsersApi.UpdateVboxTeamsUser)              // 更新VboxTeamsUser
	}
	{
		tUsersRouterWithoutRecord.GET("findVboxTeamsUser", tUsersApi.FindVboxTeamsUser)       // 根据ID获取VboxTeamsUser
		tUsersRouterWithoutRecord.GET("getVboxTeamsUserList", tUsersApi.GetVboxTeamsUserList) // 获取VboxTeamsUser列表
		tUsersRouterWithoutRecord.GET("findTeamUserAll", tUsersApi.FindTeamUserAll)           // 获取VboxTeamsUser列表
		tUsersRouterWithoutRecord.PUT("transferTeamUser", tUsersApi.TransferTeamUser)         // 获取VboxTeamsUser列表
	}
}
