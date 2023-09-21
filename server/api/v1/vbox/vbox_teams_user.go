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

type VboxTeamsUserApi struct {
}

var tUsersService = service.ServiceGroupApp.VboxServiceGroup.VboxTeamsUserService

// CreateVboxTeamsUser 创建VboxTeamsUser
// @Tags VboxTeamsUser
// @Summary 创建VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxTeamsUser true "创建VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /tUsers/createVboxTeamsUser [post]
func (tUsersApi *VboxTeamsUserApi) CreateVboxTeamsUser(c *gin.Context) {
	var tUsers vbox.TeamsUserReq
	err := c.ShouldBindJSON(&tUsers)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	//tUsers.CreatedBy = utils.GetUserID(c)
	if err := tUsersService.CreateVboxTeamsUser(tUsers); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxTeamsUser 删除VboxTeamsUser
// @Tags VboxTeamsUser
// @Summary 删除VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxTeamsUser true "删除VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /tUsers/deleteVboxTeamsUser [delete]
func (tUsersApi *VboxTeamsUserApi) DeleteVboxTeamsUser(c *gin.Context) {
	var tUsers vbox.TeamsUserReq
	err := c.ShouldBindJSON(&tUsers)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	if err := tUsersService.DeleteVboxTeamsUser(tUsers.SysUserIDS, tUsers.TeamID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxTeamsUserByIds 批量删除VboxTeamsUser
// @Tags VboxTeamsUser
// @Summary 批量删除VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /tUsers/deleteVboxTeamsUserByIds [delete]
func (tUsersApi *VboxTeamsUserApi) DeleteVboxTeamsUserByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := tUsersService.DeleteVboxTeamsUserByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxTeamsUser 更新VboxTeamsUser
// @Tags VboxTeamsUser
// @Summary 更新VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxTeamsUser true "更新VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /tUsers/updateVboxTeamsUser [put]
func (tUsersApi *VboxTeamsUserApi) UpdateVboxTeamsUser(c *gin.Context) {
	var tUsers vbox.VboxTeamsUser
	err := c.ShouldBindJSON(&tUsers)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	tUsers.UpdatedBy = utils.GetUserID(c)
	if err := tUsersService.UpdateVboxTeamsUser(tUsers); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxTeamsUser 用id查询VboxTeamsUser
// @Tags VboxTeamsUser
// @Summary 用id查询VboxTeamsUser
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxTeamsUser true "用id查询VboxTeamsUser"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /tUsers/findVboxTeamsUser [get]
func (tUsersApi *VboxTeamsUserApi) FindVboxTeamsUser(c *gin.Context) {
	var tUsers vbox.VboxTeamsUser
	err := c.ShouldBindQuery(&tUsers)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if retUsers, err := tUsersService.GetVboxTeamsUser(tUsers.Team_id); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"retUsers": retUsers}, c)
	}
}

// GetVboxTeamsUserList 分页获取VboxTeamsUser列表
// @Tags VboxTeamsUser
// @Summary 分页获取VboxTeamsUser列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxTeamsUserSearch true "分页获取VboxTeamsUser列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /tUsers/getVboxTeamsUserList [get]
func (tUsersApi *VboxTeamsUserApi) GetVboxTeamsUserList(c *gin.Context) {
	var pageInfo vboxReq.VboxTeamsUserSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := tUsersService.GetVboxTeamsUserInfoList(pageInfo); err != nil {
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

func (tUsersApi *VboxTeamsUserApi) FindTeamUserAll(c *gin.Context) {
	org := c.Query("teamID")
	if UserIds, err := tUsersService.FindTeamUserAll(org); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(UserIds, c)
	}
}

//
//
//func (tUsersApi *VboxTeamsUserApi) SetOrgUserAdmin(c *gin.Context) {
//	var orgUser organization.OrgUser
//	c.ShouldBindJSON(&orgUser)
//	if err := orgService.SetOrgUserAdmin(orgUser.SysUserID, orgUser.IsAdmin); err != nil {
//		global.GVA_LOG.Error("设置失败!", zap.Error(err))
//		response.FailWithMessage("设置失败", c)
//	} else {
//		response.OkWithMessage("设置成功", c)
//	}
//}
//
//func (tUsersApi *VboxTeamsUserApi) SyncAuthority(c *gin.Context) {
//	if err := tUsersService.SyncAuthority(); err != nil {
//		global.GVA_LOG.Error("同步失败!", zap.Error(err))
//		response.FailWithMessage(err.Error(), c)
//	} else {
//		response.OkWithMessage("同步成功", c)
//	}
//}
//
//func (tUsersApi *VboxTeamsUserApi) GetAuthority(c *gin.Context) {
//	if authDataList, err := tUsersService.GetOrgAuthority(); err != nil {
//		global.GVA_LOG.Error("获取失败!", zap.Error(err))
//		response.FailWithMessage("获取失败", c)
//	} else {
//		response.OkWithData(authDataList, c)
//	}
//}

func (tUsersApi *VboxTeamsUserApi) TransferTeamUser(c *gin.Context) {
	var orgUser vbox.TeamsUserReq
	c.ShouldBindJSON(&orgUser)
	if err := tUsersService.TransferTeamUser(orgUser.SysUserIDS, orgUser.TeamID, orgUser.ToTeamID); err != nil {
		global.GVA_LOG.Error("转移失败!", zap.Error(err))
		response.FailWithMessage("转移失败", c)
	} else {
		response.OkWithMessage("转移成功", c)
	}
}
