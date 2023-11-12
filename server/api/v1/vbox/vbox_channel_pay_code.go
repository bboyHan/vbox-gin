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

type VboxChannelPayCodeApi struct {
}

var vboxChannelPayCodeService = service.ServiceGroupApp.VboxServiceGroup.VboxChannelPayCodeService

// CreateVboxChannelPayCode 创建通道账户付款二维码
// @Tags VboxChannelPayCode
// @Summary 创建通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxChannelPayCode true "创建通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /vboxChannelPayCode/createVboxChannelPayCode [post]
func (vboxChannelPayCodeApi *VboxChannelPayCodeApi) CreateVboxChannelPayCode(c *gin.Context) {
	var vboxChannelPayCode vbox.VboxChannelPayCode
	err := c.ShouldBindJSON(&vboxChannelPayCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vboxChannelPayCode.CreatedBy = utils.GetUserID(c)
	if err := vboxChannelPayCodeService.CreateVboxChannelPayCode(&vboxChannelPayCode); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxChannelPayCode 删除通道账户付款二维码
// @Tags VboxChannelPayCode
// @Summary 删除通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxChannelPayCode true "删除通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vboxChannelPayCode/deleteVboxChannelPayCode [delete]
func (vboxChannelPayCodeApi *VboxChannelPayCodeApi) DeleteVboxChannelPayCode(c *gin.Context) {
	var vboxChannelPayCode vbox.VboxChannelPayCode
	err := c.ShouldBindJSON(&vboxChannelPayCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vboxChannelPayCode.DeletedBy = utils.GetUserID(c)
	if err := vboxChannelPayCodeService.DeleteVboxChannelPayCode(vboxChannelPayCode); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxChannelPayCodeByIds 批量删除通道账户付款二维码
// @Tags VboxChannelPayCode
// @Summary 批量删除通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /vboxChannelPayCode/deleteVboxChannelPayCodeByIds [delete]
func (vboxChannelPayCodeApi *VboxChannelPayCodeApi) DeleteVboxChannelPayCodeByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := vboxChannelPayCodeService.DeleteVboxChannelPayCodeByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxChannelPayCode 更新通道账户付款二维码
// @Tags VboxChannelPayCode
// @Summary 更新通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxChannelPayCode true "更新通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vboxChannelPayCode/updateVboxChannelPayCode [put]
func (vboxChannelPayCodeApi *VboxChannelPayCodeApi) UpdateVboxChannelPayCode(c *gin.Context) {
	var vboxChannelPayCode vbox.VboxChannelPayCode
	err := c.ShouldBindJSON(&vboxChannelPayCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vboxChannelPayCode.UpdatedBy = utils.GetUserID(c)
	if err := vboxChannelPayCodeService.UpdateVboxChannelPayCode(vboxChannelPayCode); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxChannelPayCode 用id查询通道账户付款二维码
// @Tags VboxChannelPayCode
// @Summary 用id查询通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxChannelPayCode true "用id查询通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vboxChannelPayCode/findVboxChannelPayCode [get]
func (vboxChannelPayCodeApi *VboxChannelPayCodeApi) FindVboxChannelPayCode(c *gin.Context) {
	var vboxChannelPayCode vbox.VboxChannelPayCode
	err := c.ShouldBindQuery(&vboxChannelPayCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if revboxChannelPayCode, err := vboxChannelPayCodeService.GetVboxChannelPayCode(vboxChannelPayCode.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"revboxChannelPayCode": revboxChannelPayCode}, c)
	}
}

// GetVboxChannelPayCodeList 分页获取通道账户付款二维码列表
// @Tags VboxChannelPayCode
// @Summary 分页获取通道账户付款二维码列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxChannelPayCodeSearch true "分页获取通道账户付款二维码列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vboxChannelPayCode/getVboxChannelPayCodeList [get]
func (vboxChannelPayCodeApi *VboxChannelPayCodeApi) GetVboxChannelPayCodeList(c *gin.Context) {
	var pageInfo vboxReq.VboxChannelPayCodeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := vboxChannelPayCodeService.GetVboxChannelPayCodeInfoList(pageInfo); err != nil {
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
