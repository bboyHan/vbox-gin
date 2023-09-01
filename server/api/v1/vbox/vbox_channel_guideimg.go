package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox_channel_guideimg"
	vbox_channel_guideimgReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox_channel_guideimg/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChannelGuideImgApi struct {
}

// CreateChannelGuideImg 创建Channel_guideimg
// @Tags Channel_guideimg
// @Summary 创建Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox_channel_guideimg.Channel_guideimg true "创建Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chGuideImg/createChannelGuideImg [post]
func (chGuideImgApi *ChannelGuideImgApi) CreateChannelGuideImg(c *gin.Context) {
	var chGuideImg vbox_channel_guideimg.ChannelGuideImg
	err := c.ShouldBindJSON(&chGuideImg)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chGuideImg.CreatedBy = utils.GetUserID(c)
	if err := chGuideImgService.CreateChannelGuideimg(&chGuideImg); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteChannelGuideImg 删除Channel_guideimg
// @Tags Channel_guideimg
// @Summary 删除Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox_channel_guideimg.Channel_guideimg true "删除Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chGuideImg/deleteChannelGuideImg [delete]
func (chGuideImgApi *ChannelGuideImgApi) DeleteChannelGuideImg(c *gin.Context) {
	var chGuideImg vbox_channel_guideimg.ChannelGuideImg
	err := c.ShouldBindJSON(&chGuideImg)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chGuideImg.DeletedBy = utils.GetUserID(c)
	if err := chGuideImgService.DeleteChannelGuideimg(chGuideImg); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteChannelGuideImgByIds 批量删除Channel_guideimg
// @Tags Channel_guideimg
// @Summary 批量删除Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /chGuideImg/deleteChannelGuideImgByIds [delete]
func (chGuideImgApi *ChannelGuideImgApi) DeleteChannelGuideImgByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := chGuideImgService.DeleteChannelGuideimgByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateChannelGuideImg 更新Channel_guideimg
// @Tags Channel_guideimg
// @Summary 更新Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox_channel_guideimg.Channel_guideimg true "更新Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /chGuideImg/updateChannelGuideImg [put]
func (chGuideImgApi *ChannelGuideImgApi) UpdateChannelGuideImg(c *gin.Context) {
	var chGuideImg vbox_channel_guideimg.ChannelGuideImg
	err := c.ShouldBindJSON(&chGuideImg)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chGuideImg.UpdatedBy = utils.GetUserID(c)
	if err := chGuideImgService.UpdateChannelGuideimg(chGuideImg); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindChannelGuideImg 用id查询Channel_guideimg
// @Tags Channel_guideimg
// @Summary 用id查询Channel_guideimg
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox_channel_guideimg.Channel_guideimg true "用id查询Channel_guideimg"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /chGuideImg/findChannelGuideImg [get]
func (chGuideImgApi *ChannelGuideImgApi) FindChannelGuideImg(c *gin.Context) {
	var chGuideImg vbox_channel_guideimg.ChannelGuideImg
	err := c.ShouldBindQuery(&chGuideImg)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rechGuideImg, err := chGuideImgService.GetChannelGuideimg(chGuideImg.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rechGuideImg": rechGuideImg}, c)
	}
}

// GetChannelGuideImgList 分页获取Channel_guideimg列表
// @Tags Channel_guideimg
// @Summary 分页获取Channel_guideimg列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox_channel_guideimgReq.Channel_guideimgSearch true "分页获取Channel_guideimg列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chGuideImg/getChannelGuideImgList [get]
func (chGuideImgApi *ChannelGuideImgApi) GetChannelGuideImgList(c *gin.Context) {
	var pageInfo vbox_channel_guideimgReq.ChannelGuideImgSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := chGuideImgService.GetChannelGuideimgInfoList(pageInfo); err != nil {
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

// getChannelGuideImgTaskList 获取通道引导图片
// @Tags Channel_guideimg
// @Summary 获取通道引导图片
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox_channel_guideimgReq.Channel_guideimgSearch true "分页获取Channel_guideimg列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chGuideImg/getChannelGuideImgTaskList [get]
func (chGuideImgApi *ChannelGuideImgApi) GetChannelGuideImgTaskList(c *gin.Context) {
	var pageInfo vbox_channel_guideimgReq.ChannelGuideImgTask
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := chGuideImgService.GetChannelGuideImgTaskList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "获取成功", c)
	}
}
