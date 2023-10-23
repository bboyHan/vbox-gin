package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VboxProxyApi struct {
}

var vboxProxyService = service.ServiceGroupApp.VboxServiceGroup.ProxyService

// CreateVboxProxy 创建信道
// @Tags Proxy
// @Summary 创建信道
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.Proxy true "创建信道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /vboxProxy/createVboxProxy [post]
func (vboxProxyApi *VboxProxyApi) CreateVboxProxy(c *gin.Context) {
	var vboxProxy vbox.Proxy
	err := c.ShouldBindJSON(&vboxProxy)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vboxProxyService.CreateVboxProxy(&vboxProxy); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxProxy 删除信道
// @Tags Proxy
// @Summary 删除信道
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.Proxy true "删除信道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vboxProxy/deleteVboxProxy [delete]
func (vboxProxyApi *VboxProxyApi) DeleteVboxProxy(c *gin.Context) {
	var vboxProxy vbox.Proxy
	err := c.ShouldBindJSON(&vboxProxy)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vboxProxyService.DeleteVboxProxy(vboxProxy); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxProxyByIds 批量删除信道
// @Tags Proxy
// @Summary 批量删除信道
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除信道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /vboxProxy/deleteVboxProxyByIds [delete]
func (vboxProxyApi *VboxProxyApi) DeleteVboxProxyByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vboxProxyService.DeleteVboxProxyByIds(IDS); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxProxy 更新信道
// @Tags Proxy
// @Summary 更新信道
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.Proxy true "更新信道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vboxProxy/updateVboxProxy [put]
func (vboxProxyApi *VboxProxyApi) UpdateVboxProxy(c *gin.Context) {
	var vboxProxy vbox.Proxy
	err := c.ShouldBindJSON(&vboxProxy)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := vboxProxyService.UpdateVboxProxy(vboxProxy); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxProxy 用id查询信道
// @Tags Proxy
// @Summary 用id查询信道
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.Proxy true "用id查询信道"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vboxProxy/findVboxProxy [get]
func (vboxProxyApi *VboxProxyApi) FindVboxProxy(c *gin.Context) {
	var vboxProxy vbox.Proxy
	err := c.ShouldBindQuery(&vboxProxy)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if revboxProxy, err := vboxProxyService.GetVboxProxy(vboxProxy.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"revboxProxy": revboxProxy}, c)
	}
}

// GetVboxProxyList 分页获取信道列表
// @Tags Proxy
// @Summary 分页获取信道列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxProxySearch true "分页获取信道列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vboxProxy/getVboxProxyList [get]
func (vboxProxyApi *VboxProxyApi) GetVboxProxyList(c *gin.Context) {
	var pageInfo vboxReq.VboxProxySearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := vboxProxyService.GetVboxProxyInfoList(pageInfo); err != nil {
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
