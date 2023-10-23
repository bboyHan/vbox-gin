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

type ChannelProductApi struct {
}

var vcpService = service.ServiceGroupApp.VboxServiceGroup.ChannelProductService

// CreateChannelProduct 创建通道产品
// @Tags ChannelProduct
// @Summary 创建通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelProduct true "创建通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /vcp/createChannelProduct [post]
func (vcpApi *ChannelProductApi) CreateChannelProduct(c *gin.Context) {
	var vcp vbox.ChannelProduct
	err := c.ShouldBindJSON(&vcp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vcp.CreatedBy = utils.GetUserID(c)
	if err := vcpService.CreateChannelProduct(&vcp); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteChannelProduct 删除通道产品
// @Tags ChannelProduct
// @Summary 删除通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelProduct true "删除通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vcp/deleteChannelProduct [delete]
func (vcpApi *ChannelProductApi) DeleteChannelProduct(c *gin.Context) {
	var vcp vbox.ChannelProduct
	err := c.ShouldBindJSON(&vcp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vcp.DeletedBy = utils.GetUserID(c)
	if err := vcpService.DeleteChannelProduct(vcp); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteChannelProductByIds 批量删除通道产品
// @Tags ChannelProduct
// @Summary 批量删除通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /vcp/deleteChannelProductByIds [delete]
func (vcpApi *ChannelProductApi) DeleteChannelProductByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := vcpService.DeleteChannelProductByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateChannelProduct 更新通道产品
// @Tags ChannelProduct
// @Summary 更新通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelProduct true "更新通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vcp/updateChannelProduct [put]
func (vcpApi *ChannelProductApi) UpdateChannelProduct(c *gin.Context) {
	var vcp vbox.ChannelProduct
	err := c.ShouldBindJSON(&vcp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	vcp.UpdatedBy = utils.GetUserID(c)
	if err := vcpService.UpdateChannelProduct(vcp); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindChannelProduct 用id查询通道产品
// @Tags ChannelProduct
// @Summary 用id查询通道产品
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.ChannelProduct true "用id查询通道产品"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vcp/findChannelProduct [get]
func (vcpApi *ChannelProductApi) FindChannelProduct(c *gin.Context) {
	var vcp vbox.ChannelProduct
	err := c.ShouldBindQuery(&vcp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if revcp, err := vcpService.GetChannelProduct(vcp.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"revcp": revcp}, c)
	}
}

// GetChannelProductList 分页获取通道产品列表
// @Tags ChannelProduct
// @Summary 分页获取通道产品列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.ChannelProductSearch true "分页获取通道产品列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getChannelProductList [get]
func (vcpApi *ChannelProductApi) GetChannelProductList(c *gin.Context) {
	var pageInfo vboxReq.ChannelProductSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := vcpService.GetChannelProductInfoList(pageInfo); err != nil {
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

// GetChannelProductAll 获取ChannelProduct所有列表
// @Tags VboxChannelProduct
// @Summary 获取ChannelProduct所有列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxChannelProductSearch true "获取ChannelProduct所有列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getVboxChannelProductList [get]
func (vcpApi *ChannelProductApi) GetChannelProductAll(c *gin.Context) {
	var pageInfo vboxReq.ChannelProductSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, err := vcpService.GetChannelProductAll(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(gin.H{"list": list}, c)
	}
}

// GetChannelProductSelf 获取ChannelProduct所有列表(当前用户组织下的)
// @Tags VboxChannelProduct
// @Summary 获取ChannelProduct所有列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.ChannelProductSearch true "获取ChannelProduct所有列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getChannelProductList [get]
func (vcpApi *ChannelProductApi) GetChannelProductSelf(c *gin.Context) {
	orgIds := utils2.GetOrgIDS(c)
	if list, err := vcpService.GetOrgProductList(orgIds); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(gin.H{"list": list}, c)
	}
}
