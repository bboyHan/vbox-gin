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
	"github.com/flipped-aurora/gin-vue-admin/server/utils/captcha"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChannelPayCodeApi struct {
}

var channelPayCodeService = service.ServiceGroupApp.VboxServiceGroup.ChannelPayCodeService

// GetPayCodeOverview 获取预产统计情况
// @Tags ChannelPayCode
// @Summary 获取预产统计情况
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelPayCode true "获取预产统计情况"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vboxChannelPayCode/getPayCodeOverview [post]
func (channelPayCodeApi *ChannelPayCodeApi) GetPayCodeOverview(c *gin.Context) {
	var pageInfo vboxReq.ChannelPayCodeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pageInfo.CreatedBy = utils.GetUserID(c)
	if list, err := channelPayCodeService.GetPayCodeOverview(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(gin.H{"list": list}, c)
	}
}

// CreateChannelPayCode 创建通道账户付款二维码
// @Tags ChannelPayCode
// @Summary 创建通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelPayCode true "创建通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /vboxChannelPayCode/createVboxChannelPayCode [post]
func (channelPayCodeApi *ChannelPayCodeApi) CreateChannelPayCode(c *gin.Context) {
	var channelPayCode vbox.ChannelPayCode
	err := c.ShouldBindJSON(&channelPayCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 验证图片二维码合法性
	imgB64 := channelPayCode.ImgBaseStr
	content, err := captcha.ParseQrCodeImageFromBase64(imgB64)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	global.GVA_LOG.Info("图片解析内容 ----", zap.Any("content", content))

	channelPayCode.ImgContent = content
	channelPayCode.CreatedBy = utils.GetUserID(c)
	if err = channelPayCodeService.CreateChannelPayCode(&channelPayCode); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteChannelPayCode 删除通道账户付款二维码
// @Tags ChannelPayCode
// @Summary 删除通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelPayCode true "删除通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vboxChannelPayCode/deleteVboxChannelPayCode [delete]
func (channelPayCodeApi *ChannelPayCodeApi) DeleteChannelPayCode(c *gin.Context) {
	var channelPayCode vbox.ChannelPayCode
	err := c.ShouldBindJSON(&channelPayCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	channelPayCode.DeletedBy = utils.GetUserID(c)
	if err := channelPayCodeService.DeleteChannelPayCode(channelPayCode); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteChannelPayCodeByIds 批量删除通道账户付款二维码
// @Tags ChannelPayCode
// @Summary 批量删除通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /vboxChannelPayCode/deleteVboxChannelPayCodeByIds [delete]
func (channelPayCodeApi *ChannelPayCodeApi) DeleteChannelPayCodeByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := channelPayCodeService.DeleteChannelPayCodeByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxChannelPayCode 更新通道账户付款二维码
// @Tags ChannelPayCode
// @Summary 更新通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelPayCode true "更新通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vboxChannelPayCode/updateVboxChannelPayCode [put]
func (channelPayCodeApi *ChannelPayCodeApi) UpdateVboxChannelPayCode(c *gin.Context) {
	var channelPayCode vbox.ChannelPayCode
	err := c.ShouldBindJSON(&channelPayCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	channelPayCode.UpdatedBy = utils.GetUserID(c)
	if err := channelPayCodeService.UpdateChannelPayCode(channelPayCode); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindChannelPayCode 用id查询通道账户付款二维码
// @Tags ChannelPayCode
// @Summary 用id查询通道账户付款二维码
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.ChannelPayCode true "用id查询通道账户付款二维码"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vboxChannelPayCode/findVboxChannelPayCode [get]
func (channelPayCodeApi *ChannelPayCodeApi) FindChannelPayCode(c *gin.Context) {
	var vboxChannelPayCode vbox.ChannelPayCode
	err := c.ShouldBindQuery(&vboxChannelPayCode)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if reChannelPayCode, err := channelPayCodeService.GetChannelPayCode(vboxChannelPayCode.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rechannelPayCode": reChannelPayCode}, c)
	}
}

// GetChannelPayCodeList 分页获取通道账户付款二维码列表
// @Tags ChannelPayCode
// @Summary 分页获取通道账户付款二维码列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.ChannelPayCodeSearch true "分页获取通道账户付款二维码列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vboxChannelPayCode/getVboxChannelPayCodeList [get]
func (channelPayCodeApi *ChannelPayCodeApi) GetChannelPayCodeList(c *gin.Context) {
	var pageInfo vboxReq.ChannelPayCodeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := utils2.GetUserIDS(c)
	if list, total, err := channelPayCodeService.GetChannelPayCodeInfoList(pageInfo, ids); err != nil {
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

// GetChannelPayCodeStatisByLocation 分页获取二维码统计排名列表
// @Tags ChannelPayCode
// @Summary 分页获取通道账户付款二维码列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.getChannelPayCodeStatisByLocation true "分页获取二维码统计排名列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vboxChannelPayCode/getChannelPayCodeStatisByLocation [get]
func (channelPayCodeApi *ChannelPayCodeApi) GetChannelPayCodeStatisByLocation(c *gin.Context) {
	var pageInfo vboxReq.ChannelPayCodeSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := utils2.GetUserIDS(c)
	if list, total, err := channelPayCodeService.GetChannelPayCodeNumsByLocation(pageInfo, ids); err != nil {
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
