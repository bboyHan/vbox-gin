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

type ChannelShopApi struct {
}

var channelShopService = service.ServiceGroupApp.VboxServiceGroup.ChannelShopService

// CreateChannelShop 创建引导商铺
// @Tags ChannelShop
// @Summary 创建引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelShop true "创建引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /channelShop/createChannelShop [post]
func (channelShopApi *ChannelShopApi) CreateChannelShop(c *gin.Context) {
	var channelShop vboxReq.ChannelShop
	err := c.ShouldBindJSON(&channelShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	channelShop.CreatedBy = utils.GetUserID(c)
	verify := utils.Rules{
		"Cid":             {utils.NotEmpty()},
		"ShopRemark":      {utils.NotEmpty()},
		"ChannelShopList": {utils.NotEmpty()},
	}
	if err := utils.Verify(channelShop, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := channelShopService.CreateChannelShop(&channelShop); err != nil {
		global.GVA_LOG.Error("创建/更新失败!", zap.Error(err))
		response.FailWithMessage("创建/更新失败", c)
	} else {
		response.OkWithMessage("创建/更新成功", c)
	}
}

// DeleteChannelShop 删除引导商铺
// @Tags ChannelShop
// @Summary 删除引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelShop true "删除引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /channelShop/deleteChannelShop [delete]
func (channelShopApi *ChannelShopApi) DeleteChannelShop(c *gin.Context) {
	var channelShop vbox.ChannelShop
	err := c.ShouldBindJSON(&channelShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	channelShop.DeletedBy = utils.GetUserID(c)
	if err := channelShopService.DeleteChannelShop(channelShop); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteChannelShopByIds 批量删除引导商铺
// @Tags ChannelShop
// @Summary 批量删除引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /channelShop/deleteChannelShopByIds [delete]
func (channelShopApi *ChannelShopApi) DeleteChannelShopByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := channelShopService.DeleteChannelShopByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateChannelShop 更新引导商铺
// @Tags ChannelShop
// @Summary 更新引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelShop true "更新引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /channelShop/updateChannelShop [put]
func (channelShopApi *ChannelShopApi) UpdateChannelShop(c *gin.Context) {
	var channelShop vbox.ChannelShop
	err := c.ShouldBindJSON(&channelShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	channelShop.UpdatedBy = utils.GetUserID(c)
	verify := utils.Rules{
		"ShopRemark": {utils.NotEmpty()},
	}
	if err := utils.Verify(channelShop, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := channelShopService.UpdateChannelShop(channelShop); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindChannelShop 用id查询引导商铺
// @Tags ChannelShop
// @Summary 用id查询引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.ChannelShop true "用id查询引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /channelShop/findChannelShop [get]
func (channelShopApi *ChannelShopApi) FindChannelShop(c *gin.Context) {
	var channelShop vbox.ChannelShop
	err := c.ShouldBindQuery(&channelShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rechannelShop, err := channelShopService.GetChannelShop(channelShop.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rechannelShop": rechannelShop}, c)
	}
}

// FindChannelShopByProductId 用ProductId查询引导商铺
// @Tags ChannelShop
// @Summary 用id查询引导商铺
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.ChannelShop true "用id查询引导商铺"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /channelShop/findChannelShop [get]
func (channelShopApi *ChannelShopApi) FindChannelShopByProductId(c *gin.Context) {
	var channelShop vbox.ChannelShop
	err := c.ShouldBindQuery(&channelShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rechannelShop, err := channelShopService.GetChannelShopByProductId(channelShop.ProductId); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rechannelShop": rechannelShop}, c)
	}
}

// GetChannelShopList 分页获取引导商铺列表
// @Tags ChannelShop
// @Summary 分页获取引导商铺列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.ChannelShopSearch true "分页获取引导商铺列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /channelShop/getChannelShopList [get]
func (channelShopApi *ChannelShopApi) GetChannelShopList(c *gin.Context) {
	var pageInfo vboxReq.ChannelShopSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	ids := utils2.GetUserIDS(c)
	if list, err := channelShopService.GetChannelShopInfoList(pageInfo, ids); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(gin.H{"list": list}, "获取成功", c)
	}
}
