package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VboxTeamsApi struct {
}

var teamsService = service.ServiceGroupApp.VboxServiceGroup.VboxTeamsService

// CreateVboxTeams 创建VboxTeams
// @Tags VboxTeams
// @Summary 创建VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxTeams true "创建VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /teams/createVboxTeams [post]
func (teamsApi *VboxTeamsApi) CreateVboxTeams(c *gin.Context) {
	var teams vbox.VboxTeams
	err := c.ShouldBindJSON(&teams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	teams.CreatedBy = utils.GetUserID(c)
	if err := teamsService.CreateVboxTeams(&teams); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxTeams 删除VboxTeams
// @Tags VboxTeams
// @Summary 删除VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxTeams true "删除VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /teams/deleteVboxTeams [delete]
func (teamsApi *VboxTeamsApi) DeleteVboxTeams(c *gin.Context) {
	var teams vbox.VboxTeams
	err := c.ShouldBindJSON(&teams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	teams.DeletedBy = utils.GetUserID(c)
	if err := teamsService.DeleteVboxTeams(teams); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxTeamsByIds 批量删除VboxTeams
// @Tags VboxTeams
// @Summary 批量删除VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /teams/deleteVboxTeamsByIds [delete]
func (teamsApi *VboxTeamsApi) DeleteVboxTeamsByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := teamsService.DeleteVboxTeamsByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxTeams 更新VboxTeams
// @Tags VboxTeams
// @Summary 更新VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxTeams true "更新VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /teams/updateVboxTeams [put]
func (teamsApi *VboxTeamsApi) UpdateVboxTeams(c *gin.Context) {
	var teams vbox.VboxTeams
	err := c.ShouldBindJSON(&teams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	teams.UpdatedBy = utils.GetUserID(c)
	if err := teamsService.UpdateVboxTeams(teams); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxTeams 用id查询VboxTeams
// @Tags VboxTeams
// @Summary 用id查询VboxTeams
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxTeams true "用id查询VboxTeams"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /teams/findVboxTeams [get]
func (teamsApi *VboxTeamsApi) FindVboxTeams(c *gin.Context) {
	var teams vbox.VboxTeams
	err := c.ShouldBindQuery(&teams)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if reteams, err := teamsService.GetVboxTeams(teams.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"reteams": reteams}, c)
	}
}

// GetVboxTeamsList 分页获取VboxTeams列表
// @Tags VboxTeams
// @Summary 分页获取VboxTeams列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxTeamsSearch true "分页获取VboxTeams列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /teams/getVboxTeamsList [get]
func (teamsApi *VboxTeamsApi) GetVboxTeamsList(c *gin.Context) {
	var pageInfo vboxReq.VboxTeamsSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := teamsService.GetVboxTeamsInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
