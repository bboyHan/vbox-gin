package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PayAccountApi struct {
}

var paccService = service.ServiceGroupApp.VboxServiceGroup.PayAccountService

// SwitchEnablePayAccount 开启或关闭VboxPayAccount
// @Tags VboxPayAccount
// @Summary 开启或关闭VboxPayAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body payermng.VboxPayAccount true "开启或关闭VboxPayAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vboxPayAccount/updateVboxPayAccount [put]
func (paccApi *PayAccountApi) SwitchEnablePayAccount(c *gin.Context) {
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
	if err := paccService.UpdatePayAccount(vboxPayAccount); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// CreatePayAccount 创建付方
// @Tags PayAccount
// @Summary 创建付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.PayAccount true "创建付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /pacc/createPayAccount [post]
func (paccApi *PayAccountApi) CreatePayAccount(c *gin.Context) {
	var pacc vbox.PayAccount
	err := c.ShouldBindJSON(&pacc)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pacc.CreatedBy = utils.GetUserID(c)
	pacc.Uid = utils.GetUserID(c)
	if err := paccService.CreatePayAccount(&pacc); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeletePayAccount 删除付方
// @Tags PayAccount
// @Summary 删除付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.PayAccount true "删除付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /pacc/deletePayAccount [delete]
func (paccApi *PayAccountApi) DeletePayAccount(c *gin.Context) {
	var pacc vbox.PayAccount
	err := c.ShouldBindJSON(&pacc)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pacc.DeletedBy = utils.GetUserID(c)
	if err := paccService.DeletePayAccount(pacc); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeletePayAccountByIds 批量删除付方
// @Tags PayAccount
// @Summary 批量删除付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /pacc/deletePayAccountByIds [delete]
func (paccApi *PayAccountApi) DeletePayAccountByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := paccService.DeletePayAccountByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdatePayAccount 更新付方
// @Tags PayAccount
// @Summary 更新付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.PayAccount true "更新付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /pacc/updatePayAccount [put]
func (paccApi *PayAccountApi) UpdatePayAccount(c *gin.Context) {
	var pacc vbox.PayAccount
	err := c.ShouldBindJSON(&pacc)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pacc.UpdatedBy = utils.GetUserID(c)
	if err := paccService.UpdatePayAccount(pacc); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindPayAccount 用id查询付方
// @Tags PayAccount
// @Summary 用id查询付方
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.PayAccount true "用id查询付方"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /pacc/findPayAccount [get]
func (paccApi *PayAccountApi) FindPayAccount(c *gin.Context) {
	var pacc vbox.PayAccount
	err := c.ShouldBindQuery(&pacc)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if repacc, err := paccService.GetPayAccount(pacc.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"repacc": repacc}, c)
	}
}

// GetPayAccountList 分页获取付方列表
// @Tags PayAccount
// @Summary 分页获取付方列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.PayAccountSearch true "分页获取付方列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /pacc/getPayAccountList [get]
func (paccApi *PayAccountApi) GetPayAccountList(c *gin.Context) {
	var pageInfo vboxReq.PayAccountSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := utils2.GetUserIDS(c)
	if list, total, err := paccService.GetPayAccountInfoList(pageInfo, ids); err != nil {
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

// GetPAccGateway 用付方信道地址（API）
// @Tags Proxy
// @Summary 用id查询信道
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.Proxy true "用id查询信道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /pacc/getPAccGateway [get]
func (paccApi *PayAccountApi) GetPAccGateway(c *gin.Context) {
	var req vboxReq.VboxProxySearch
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if reProxy, err := paccService.GetPAccGateway(req); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(reProxy, c)
	}
}
