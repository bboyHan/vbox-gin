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

type VboxChannelRateApi struct {
}

var chRateService = service.ServiceGroupApp.VboxServiceGroup.VboxChannelRateService

// CreateVboxChannelRate 创建VboxChannelRate
// @Tags VboxChannelRate
// @Summary 创建VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxChannelRate true "创建VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chRate/createVboxChannelRate [post]
func (chRateApi *VboxChannelRateApi) CreateVboxChannelRate(c *gin.Context) {
	var chRate vbox.VboxChannelRate
	err := c.ShouldBindJSON(&chRate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chRate.CreatedBy = utils.GetUserID(c)
	if err := chRateService.CreateVboxChannelRate(&chRate); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxChannelRate 删除VboxChannelRate
// @Tags VboxChannelRate
// @Summary 删除VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxChannelRate true "删除VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chRate/deleteVboxChannelRate [delete]
func (chRateApi *VboxChannelRateApi) DeleteVboxChannelRate(c *gin.Context) {
	var chRate vbox.VboxChannelRate
	err := c.ShouldBindJSON(&chRate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chRate.DeletedBy = utils.GetUserID(c)
	if err := chRateService.DeleteVboxChannelRate(chRate); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxChannelRateByIds 批量删除VboxChannelRate
// @Tags VboxChannelRate
// @Summary 批量删除VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /chRate/deleteVboxChannelRateByIds [delete]
func (chRateApi *VboxChannelRateApi) DeleteVboxChannelRateByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := chRateService.DeleteVboxChannelRateByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxChannelRate 更新VboxChannelRate
// @Tags VboxChannelRate
// @Summary 更新VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxChannelRate true "更新VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /chRate/updateVboxChannelRate [put]
func (chRateApi *VboxChannelRateApi) UpdateVboxChannelRate(c *gin.Context) {
	var chRate vbox.VboxChannelRate
	err := c.ShouldBindJSON(&chRate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chRate.UpdatedBy = utils.GetUserID(c)
	if err := chRateService.UpdateVboxChannelRate(chRate); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxChannelRate 用id查询VboxChannelRate
// @Tags VboxChannelRate
// @Summary 用id查询VboxChannelRate
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxChannelRate true "用id查询VboxChannelRate"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /chRate/findVboxChannelRate [get]
func (chRateApi *VboxChannelRateApi) FindVboxChannelRate(c *gin.Context) {
	var chRate vbox.VboxChannelRate
	err := c.ShouldBindQuery(&chRate)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rechRate, err := chRateService.GetVboxChannelRate(chRate.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rechRate": rechRate}, c)
	}
}

// GetVboxChannelRateList 分页获取VboxChannelRate列表
// @Tags VboxChannelRate
// @Summary 分页获取VboxChannelRate列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxChannelRateSearch true "分页获取VboxChannelRate列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chRate/getVboxChannelRateList [get]
func (chRateApi *VboxChannelRateApi) GetVboxChannelRateList(c *gin.Context) {
	var pageInfo vboxReq.VboxChannelRateSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := chRateService.GetVboxChannelRateInfoList(pageInfo); err != nil {
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

// GetVboxTeamUserChannelRateList 分页获取VboxChannelRate列表
// @Tags VboxChannelRate
// @Summary 分页获取VboxChannelRate列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxChannelRateSearch true "分页获取VboxChannelRate列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chRate/getVboxTeamUserChannelRateList [get]
func (chRateApi *VboxChannelRateApi) GetVboxTeamUserChannelRateList(c *gin.Context) {
	var pageInfo vboxReq.VboxChannelRateSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := chRateService.GetVboxTeamUserChannelRateList(pageInfo); err != nil {
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
