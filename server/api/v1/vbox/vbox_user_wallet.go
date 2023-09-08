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

type VboxUserWalletApi struct {
}

var vuwService = service.ServiceGroupApp.VboxServiceGroup.VboxUserWalletService


// CreateVboxUserWallet 创建VboxUserWallet
// @Tags VboxUserWallet
// @Summary 创建VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxUserWallet true "创建VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vuw/createVboxUserWallet [post]
func (vuwApi *VboxUserWalletApi) CreateVboxUserWallet(c *gin.Context) {
	var vuw vbox.VboxUserWallet
	err := c.ShouldBindJSON(&vuw)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    vuw.CreatedBy = utils.GetUserID(c)
	if err := vuwService.CreateVboxUserWallet(&vuw); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxUserWallet 删除VboxUserWallet
// @Tags VboxUserWallet
// @Summary 删除VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxUserWallet true "删除VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vuw/deleteVboxUserWallet [delete]
func (vuwApi *VboxUserWalletApi) DeleteVboxUserWallet(c *gin.Context) {
	var vuw vbox.VboxUserWallet
	err := c.ShouldBindJSON(&vuw)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    vuw.DeletedBy = utils.GetUserID(c)
	if err := vuwService.DeleteVboxUserWallet(vuw); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxUserWalletByIds 批量删除VboxUserWallet
// @Tags VboxUserWallet
// @Summary 批量删除VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /vuw/deleteVboxUserWalletByIds [delete]
func (vuwApi *VboxUserWalletApi) DeleteVboxUserWalletByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    deletedBy := utils.GetUserID(c)
	if err := vuwService.DeleteVboxUserWalletByIds(IDS,deletedBy); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxUserWallet 更新VboxUserWallet
// @Tags VboxUserWallet
// @Summary 更新VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxUserWallet true "更新VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vuw/updateVboxUserWallet [put]
func (vuwApi *VboxUserWalletApi) UpdateVboxUserWallet(c *gin.Context) {
	var vuw vbox.VboxUserWallet
	err := c.ShouldBindJSON(&vuw)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    vuw.UpdatedBy = utils.GetUserID(c)
	if err := vuwService.UpdateVboxUserWallet(vuw); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxUserWallet 用id查询VboxUserWallet
// @Tags VboxUserWallet
// @Summary 用id查询VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxUserWallet true "用id查询VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vuw/findVboxUserWallet [get]
func (vuwApi *VboxUserWalletApi) FindVboxUserWallet(c *gin.Context) {
	var vuw vbox.VboxUserWallet
	err := c.ShouldBindQuery(&vuw)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if revuw, err := vuwService.GetVboxUserWallet(vuw.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"revuw": revuw}, c)
	}
}

// GetVboxUserWalletList 分页获取VboxUserWallet列表
// @Tags VboxUserWallet
// @Summary 分页获取VboxUserWallet列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxUserWalletSearch true "分页获取VboxUserWallet列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vuw/getVboxUserWalletList [get]
func (vuwApi *VboxUserWalletApi) GetVboxUserWalletList(c *gin.Context) {
	var pageInfo vboxReq.VboxUserWalletSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := vuwService.GetVboxUserWalletInfoList(pageInfo); err != nil {
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
