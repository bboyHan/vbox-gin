package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type VboxTeamsRouter struct {
}

// InitVboxTeamsRouter 初始化 VboxTeams 路由信息
func (s *VboxTeamsRouter) InitVboxTeamsRouter(Router *gin.RouterGroup) {
	teamsRouter := Router.Group("teams").Use(middleware.OperationRecord())
	teamsRouterWithoutRecord := Router.Group("teams")
	var teamsApi = v1.ApiGroupApp.Vbox.VboxTeamsApi
	{
		teamsRouter.POST("createVboxTeams", teamsApi.CreateVboxTeams)             // 新建VboxTeams
		teamsRouter.DELETE("deleteVboxTeams", teamsApi.DeleteVboxTeams)           // 删除VboxTeams
		teamsRouter.DELETE("deleteVboxTeamsByIds", teamsApi.DeleteVboxTeamsByIds) // 批量删除VboxTeams
		teamsRouter.PUT("updateVboxTeams", teamsApi.UpdateVboxTeams)              // 更新VboxTeams
	}
	{
		teamsRouterWithoutRecord.GET("findVboxTeams", teamsApi.FindVboxTeams)       // 根据ID获取VboxTeams
		teamsRouterWithoutRecord.GET("getVboxTeamsList", teamsApi.GetVboxTeamsList) // 获取VboxTeams列表
	}
}
