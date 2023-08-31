package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PayAccountApi struct {
}

// SwitchEnableVboxPayAccount 开启或关闭VboxPayAccount
// @Tags VboxPayAccount
// @Summary 开启或关闭VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body payermng.VboxPayAccount true "开启或关闭VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vboxPayAccount/updateVboxPayAccount [put]
func (vboxPayAccountApi *PayAccountApi) SwitchEnableVboxPayAccount(c *gin.Context) {
	var vboxPayAccount vbox.PayAccount
	err := c.ShouldBindJSON(&vboxPayAccount)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"status": {utils.NotEmpty()},
	}
	if err := utils.Verify(vboxPayAccount, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vpaService.UpdateVboxPayAccount(vboxPayAccount); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// CreateVboxPayAccount 创建VboxPayAccount
// @Tags VboxPayAccount
// @Summary 创建VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayAccount true "创建VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpa/createVboxPayAccount [post]
func (vboxPayAccountApi *PayAccountApi) CreateVboxPayAccount(c *gin.Context) {
	var vpa vbox.PayAccount
	err := c.ShouldBindJSON(&vpa)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vpa.Uid = utils.GetUserID(c)
	if err := vpaService.CreateVboxPayAccount(&vpa); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxPayAccount 删除VboxPayAccount
// @Tags VboxPayAccount
// @Summary 删除VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayAccount true "删除VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vpa/deleteVboxPayAccount [delete]
func (vboxPayAccountApi *PayAccountApi) DeleteVboxPayAccount(c *gin.Context) {
	var vpa vbox.PayAccount
	err := c.ShouldBindJSON(&vpa)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vpaService.DeleteVboxPayAccount(vpa); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxPayAccountByIds 批量删除VboxPayAccount
// @Tags VboxPayAccount
// @Summary 批量删除VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /vpa/deleteVboxPayAccountByIds [delete]
func (vboxPayAccountApi *PayAccountApi) DeleteVboxPayAccountByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vpaService.DeleteVboxPayAccountByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxPayAccount 更新VboxPayAccount
// @Tags VboxPayAccount
// @Summary 更新VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayAccount true "更新VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vpa/updateVboxPayAccount [put]
func (vboxPayAccountApi *PayAccountApi) UpdateVboxPayAccount(c *gin.Context) {
	var vpa vbox.PayAccount
	err := c.ShouldBindJSON(&vpa)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vpaService.UpdateVboxPayAccount(vpa); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxPayAccount 用id查询VboxPayAccount
// @Tags VboxPayAccount
// @Summary 用id查询VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxPayAccount true "用id查询VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vpa/findVboxPayAccount [get]
func (vboxPayAccountApi *PayAccountApi) FindVboxPayAccount(c *gin.Context) {
	var vpa vbox.PayAccount
	err := c.ShouldBindQuery(&vpa)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if revpa, err := vpaService.GetVboxPayAccount(vpa.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"revpa": revpa}, c)
	}
}

// GetVboxPayAccountList 分页获取VboxPayAccount列表
// @Tags VboxPayAccount
// @Summary 分页获取VboxPayAccount列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxPayAccountSearch true "分页获取VboxPayAccount列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpa/getVboxPayAccountList [get]
func (vboxPayAccountApi *PayAccountApi) GetVboxPayAccountList(c *gin.Context) {
	var pageInfo vboxReq.VboxPayAccountSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := vpaService.GetVboxPayAccountInfoList(pageInfo); err != nil {
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
