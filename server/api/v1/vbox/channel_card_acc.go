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

type ChannelCardAccApi struct {
}

var channelCardAccService = service.ServiceGroupApp.VboxServiceGroup.ChannelCardAccService

// SwitchEnableChannelCardAcc 开启或关闭ChannelCardAcc
// @Tags ChannelCardAcc
// @Summary 开启或关闭ChannelCardAcc
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body payermng.VboxPayAccount true "开启或关闭ChannelCardAcc"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vca/ChannelCardAcc [put]
func (cardAccApi *ChannelCardAccApi) SwitchEnableChannelCardAcc(c *gin.Context) {
	var channelCardAcc vboxReq.ChannelCardAccUpd
	err := c.ShouldBindJSON(&channelCardAcc)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	verify := utils.Rules{
		"status": {utils.NotEmpty()},
	}
	if err := utils.Verify(channelCardAcc, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	channelCardAcc.UpdatedBy = utils.GetUserID(c)

	if err := channelCardAccService.SwitchEnableChannelCardAcc(channelCardAcc, c); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// SwitchEnableChannelCardAccByIds 批量开关通道账号
// @Tags ChannelCardAcc
// @Summary 批量开关通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量开关通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量操作成功"}"
// @Router /vca/deleteChannelCardAccByIds [delete]
func (cardAccApi *ChannelCardAccApi) SwitchEnableChannelCardAccByIds(c *gin.Context) {
	var upd vboxReq.ChannelCardAccUpd
	err := c.ShouldBindJSON(&upd)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	updatedBy := utils.GetUserID(c)
	if err := channelCardAccService.SwitchEnableChannelCardAccByIds(upd, updatedBy, c); err != nil {
		global.GVA_LOG.Error("批量操作失败!", zap.Error(err))
		response.FailWithMessage("批量操作失败", c)
	} else {
		response.OkWithMessage("批量操作成功", c)
	}
}

// CreateChannelCardAcc 创建通道账号
// @Tags ChannelCardAcc
// @Summary 创建通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelCardAcc true "创建通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /vca/createChannelCardAcc [post]
func (cardAccApi *ChannelCardAccApi) CreateChannelCardAcc(c *gin.Context) {
	var vca vbox.ChannelCardAcc
	err := c.ShouldBindJSON(&vca)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vca.CreatedBy = utils.GetUserID(c)
	verify := utils.Rules{
		"Token": {utils.NotEmpty()},
		"Cid":   {utils.NotEmpty()},
	}
	if err := utils.Verify(vca, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := channelCardAccService.CreateChannelCardAcc(&vca, c); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteChannelCardAcc 删除通道账号
// @Tags ChannelCardAcc
// @Summary 删除通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelCardAcc true "删除通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vca/ChannelCardAcc [delete]
func (cardAccApi *ChannelCardAccApi) DeleteChannelCardAcc(c *gin.Context) {
	var vca vbox.ChannelCardAcc
	err := c.ShouldBindJSON(&vca)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vca.DeletedBy = utils.GetUserID(c)
	if err := channelCardAccService.DeleteChannelCardAcc(vca, c); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteChannelCardAccByIds 批量删除通道账号
// @Tags ChannelCardAcc
// @Summary 批量删除通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /vca/deleteChannelCardAccByIds [delete]
func (cardAccApi *ChannelCardAccApi) DeleteChannelCardAccByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := channelCardAccService.DeleteChannelCardAccByIds(IDS, c, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateChannelCardAcc 更新通道账号
// @Tags ChannelCardAcc
// @Summary 更新通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelCardAcc true "更新通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vca/updateCardAcc [put]
func (cardAccApi *ChannelCardAccApi) UpdateChannelCardAcc(c *gin.Context) {
	var vca vbox.ChannelCardAcc
	err := c.ShouldBindJSON(&vca)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vca.UpdatedBy = utils.GetUserID(c)
	verify := utils.Rules{
		"AcAccount": {utils.NotEmpty()},
		"Cid":       {utils.NotEmpty()},
	}
	if err := utils.Verify(vca, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := channelCardAccService.UpdateChannelCardAcc(vca); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindChannelCardAcc 用id查询通道账号
// @Tags ChannelCardAcc
// @Summary 用id查询通道账号
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.ChannelCardAcc true "用id查询通道账号"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vca/findChannelCardAcc [get]
func (cardAccApi *ChannelCardAccApi) FindChannelCardAcc(c *gin.Context) {
	var vca vbox.ChannelCardAcc
	err := c.ShouldBindQuery(&vca)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if vca.ID != 0 {
		if ret, err := channelCardAccService.GetChannelCardAcc(vca.ID); err != nil {
			global.GVA_LOG.Error("查询失败!", zap.Error(err))
			response.FailWithMessage("查询失败", c)
		} else {
			response.OkWithData(gin.H{"ret": ret}, c)
		}
	}
	if vca.AcId != "" {
		if ret, err := channelCardAccService.GetChannelCardAccByAcId(vca.AcId); err != nil {
			global.GVA_LOG.Error("查询失败!", zap.Error(err))
			response.FailWithMessage("查询失败", c)
		} else {
			response.OkWithData(gin.H{"ret": ret}, c)
		}
	}

}

// GetChannelCardAccList 分页获取通道账号列表
// @Tags ChannelCardAcc
// @Summary 分页获取通道账号列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.ChannelCardAccSearch true "分页获取通道账号列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vca/getChannelCardAccList [get]
func (cardAccApi *ChannelCardAccApi) GetChannelCardAccList(c *gin.Context) {
	var pageInfo vboxReq.ChannelCardAccSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := utils2.GetUserIDS(c)
	if list, total, err := channelCardAccService.GetChannelCardAccInfoList(pageInfo, ids); err != nil {
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
