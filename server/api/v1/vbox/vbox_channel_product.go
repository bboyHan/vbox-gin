package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/captcha"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

var redisStore = captcha.NewCustomRedisStore("channel_product:", time.Hour*12)

type ChannelProductApi struct {
}

// CreateVboxChannelProduct 创建VboxChannelProduct
// @Tags VboxChannelProduct
// @Summary 创建VboxChannelProduct
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxChannelProduct true "创建VboxChannelProduct"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/createVboxChannelProduct [post]
func (vcpApi *ChannelProductApi) CreateVboxChannelProduct(c *gin.Context) {
	var vcp vbox.ChannelProduct
	err := c.ShouldBindJSON(&vcp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vcpService.CreateVboxChannelProduct(&vcp); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxChannelProduct 删除VboxChannelProduct
// @Tags VboxChannelProduct
// @Summary 删除VboxChannelProduct
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxChannelProduct true "删除VboxChannelProduct"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vcp/deleteVboxChannelProduct [delete]
func (vcpApi *ChannelProductApi) DeleteVboxChannelProduct(c *gin.Context) {
	var vcp vbox.ChannelProduct
	err := c.ShouldBindJSON(&vcp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vcpService.DeleteVboxChannelProduct(vcp); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// UpdateVboxChannelProduct 更新VboxChannelProduct
// @Tags VboxChannelProduct
// @Summary 更新VboxChannelProduct
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxChannelProduct true "更新VboxChannelProduct"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vcp/updateVboxChannelProduct [put]
func (vcpApi *ChannelProductApi) UpdateVboxChannelProduct(c *gin.Context) {
	var vcp vbox.ChannelProduct
	err := c.ShouldBindJSON(&vcp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vcpService.UpdateVboxChannelProduct(vcp); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxChannelProduct 用id查询VboxChannelProduct
// @Tags VboxChannelProduct
// @Summary 用id查询VboxChannelProduct
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxChannelProduct true "用id查询VboxChannelProduct"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vcp/findVboxChannelProduct [get]
func (vcpApi *ChannelProductApi) FindVboxChannelProduct(c *gin.Context) {
	var vcp vbox.ChannelProduct
	err := c.ShouldBindQuery(&vcp)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if revcp, err := vcpService.GetVboxChannelProduct(vcp.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"revcp": revcp}, c)
	}
}

// GetVboxChannelProductList 分页获取VboxChannelProduct列表
// @Tags VboxChannelProduct
// @Summary 分页获取VboxChannelProduct列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxChannelProductSearch true "分页获取VboxChannelProduct列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getVboxChannelProductList [get]
func (vcpApi *ChannelProductApi) GetVboxChannelProductList(c *gin.Context) {
	var pageInfo vboxReq.VboxChannelProductSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := vcpService.GetVboxChannelProductInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		//jsonData, _ := json.Marshal(list)
		//redisStore.Set("list", string(jsonData))
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}

// GetVboxChannelProductAll 分页获取VboxChannelProduct列表
// @Tags VboxChannelProduct
// @Summary 分页获取VboxChannelProduct列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxChannelProductSearch true "分页获取VboxChannelProduct列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vcp/getVboxChannelProductList [get]
func (vcpApi *ChannelProductApi) GetVboxChannelProductAll(c *gin.Context) {
	var pageInfo vboxReq.VboxChannelProductSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, err := vcpService.GetVboxChannelProductInfoAll(); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithData(gin.H{"list": list}, c)
	}
}
