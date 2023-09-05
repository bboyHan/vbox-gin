package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
    "github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
    "github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
    "github.com/flipped-aurora/gin-vue-admin/server/service"
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "github.com/flipped-aurora/gin-vue-admin/server/utils"
)

type ChannelAccountApi struct {
}

var vcaService = service.ServiceGroupApp.VboxServiceGroup.ChannelAccountService


// CreateChannelAccount 创建ChannelAccount
// @Tags ChannelAccount
// @Summary 创建ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelAccount true "创建ChannelAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vca/createChannelAccount [post]
func (vcaApi *ChannelAccountApi) CreateChannelAccount(c *gin.Context) {
	var vca vbox.ChannelAccount
	err := c.ShouldBindJSON(&vca)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    vca.CreatedBy = utils.GetUserID(c)
	if err := vcaService.CreateChannelAccount(&vca); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteChannelAccount 删除ChannelAccount
// @Tags ChannelAccount
// @Summary 删除ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelAccount true "删除ChannelAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vca/deleteChannelAccount [delete]
func (vcaApi *ChannelAccountApi) DeleteChannelAccount(c *gin.Context) {
	var vca vbox.ChannelAccount
	err := c.ShouldBindJSON(&vca)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    vca.DeletedBy = utils.GetUserID(c)
	if err := vcaService.DeleteChannelAccount(vca); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteChannelAccountByIds 批量删除ChannelAccount
// @Tags ChannelAccount
// @Summary 批量删除ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除ChannelAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /vca/deleteChannelAccountByIds [delete]
func (vcaApi *ChannelAccountApi) DeleteChannelAccountByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    deletedBy := utils.GetUserID(c)
	if err := vcaService.DeleteChannelAccountByIds(IDS,deletedBy); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateChannelAccount 更新ChannelAccount
// @Tags ChannelAccount
// @Summary 更新ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.ChannelAccount true "更新ChannelAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vca/updateChannelAccount [put]
func (vcaApi *ChannelAccountApi) UpdateChannelAccount(c *gin.Context) {
	var vca vbox.ChannelAccount
	err := c.ShouldBindJSON(&vca)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    vca.UpdatedBy = utils.GetUserID(c)
	if err := vcaService.UpdateChannelAccount(vca); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindChannelAccount 用id查询ChannelAccount
// @Tags ChannelAccount
// @Summary 用id查询ChannelAccount
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.ChannelAccount true "用id查询ChannelAccount"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vca/findChannelAccount [get]
func (vcaApi *ChannelAccountApi) FindChannelAccount(c *gin.Context) {
	var vca vbox.ChannelAccount
	err := c.ShouldBindQuery(&vca)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if revca, err := vcaService.GetChannelAccount(vca.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"revca": revca}, c)
	}
}

// GetChannelAccountList 分页获取ChannelAccount列表
// @Tags ChannelAccount
// @Summary 分页获取ChannelAccount列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.ChannelAccountSearch true "分页获取ChannelAccount列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vca/getChannelAccountList [get]
func (vcaApi *ChannelAccountApi) GetChannelAccountList(c *gin.Context) {
	var pageInfo vboxReq.ChannelAccountSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := vcaService.GetChannelAccountInfoList(pageInfo); err != nil {
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
