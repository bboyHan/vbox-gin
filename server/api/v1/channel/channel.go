package channel

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/channel"
	channelReq "github.com/flipped-aurora/gin-vue-admin/server/model/channel/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChannelApi struct {
}

var chService = service.ServiceGroupApp.ChannelServiceGroup.ChannelService

// CreateChannel 创建Channel
// @Tags Channel
// @Summary 创建Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body channel.Channel true "创建Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ch/createChannel [post]
func (chApi *ChannelApi) CreateChannel(c *gin.Context) {
	var ch channel.Channel
	err := c.ShouldBindJSON(&ch)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ch.CreatedBy = utils.GetUserID(c)
	verify := utils.Rules{
		"Type": {utils.NotEmpty()},
	}
	if err := utils.Verify(ch, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := chService.CreateChannel(&ch); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteChannel 删除Channel
// @Tags Channel
// @Summary 删除Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body channel.Channel true "删除Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ch/deleteChannel [delete]
func (chApi *ChannelApi) DeleteChannel(c *gin.Context) {
	var ch channel.Channel
	err := c.ShouldBindJSON(&ch)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ch.DeletedBy = utils.GetUserID(c)
	if err := chService.DeleteChannel(ch); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteChannelByIds 批量删除Channel
// @Tags Channel
// @Summary 批量删除Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /ch/deleteChannelByIds [delete]
func (chApi *ChannelApi) DeleteChannelByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := chService.DeleteChannelByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateChannel 更新Channel
// @Tags Channel
// @Summary 更新Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body channel.Channel true "更新Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ch/updateChannel [put]
func (chApi *ChannelApi) UpdateChannel(c *gin.Context) {
	var ch channel.Channel
	err := c.ShouldBindJSON(&ch)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ch.UpdatedBy = utils.GetUserID(c)
	verify := utils.Rules{
		"Type": {utils.NotEmpty()},
	}
	if err := utils.Verify(ch, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := chService.UpdateChannel(ch); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindChannel 用id查询Channel
// @Tags Channel
// @Summary 用id查询Channel
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query channel.Channel true "用id查询Channel"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /ch/findChannel [get]
func (chApi *ChannelApi) FindChannel(c *gin.Context) {
	var ch channel.Channel
	err := c.ShouldBindQuery(&ch)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rech, err := chService.GetChannel(ch.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rech": rech}, c)
	}
}

// GetChannelList 分页获取Channel列表
// @Tags Channel
// @Summary 分页获取Channel列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query channelReq.ChannelSearch true "分页获取Channel列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ch/getChannelList [get]
func (chApi *ChannelApi) GetChannelList(c *gin.Context) {
	var pageInfo channelReq.ChannelSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := chService.GetChannelInfoList(pageInfo); err != nil {
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
