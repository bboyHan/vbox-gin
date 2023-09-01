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

type VboxPayOrderApi struct {
}

var vpoService = service.ServiceGroupApp.VboxServiceGroup.VboxPayOrderService


// CreateVboxPayOrder 创建VboxPayOrder
// @Tags VboxPayOrder
// @Summary 创建VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayOrder true "创建VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/createVboxPayOrder [post]
func (vpoApi *VboxPayOrderApi) CreateVboxPayOrder(c *gin.Context) {
	var vpo vbox.VboxPayOrder
	err := c.ShouldBindJSON(&vpo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    vpo.CreatedBy = utils.GetUserID(c)
	if err := vpoService.CreateVboxPayOrder(&vpo); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxPayOrder 删除VboxPayOrder
// @Tags VboxPayOrder
// @Summary 删除VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayOrder true "删除VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /vpo/deleteVboxPayOrder [delete]
func (vpoApi *VboxPayOrderApi) DeleteVboxPayOrder(c *gin.Context) {
	var vpo vbox.VboxPayOrder
	err := c.ShouldBindJSON(&vpo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    vpo.DeletedBy = utils.GetUserID(c)
	if err := vpoService.DeleteVboxPayOrder(vpo); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxPayOrderByIds 批量删除VboxPayOrder
// @Tags VboxPayOrder
// @Summary 批量删除VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /vpo/deleteVboxPayOrderByIds [delete]
func (vpoApi *VboxPayOrderApi) DeleteVboxPayOrderByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    deletedBy := utils.GetUserID(c)
	if err := vpoService.DeleteVboxPayOrderByIds(IDS,deletedBy); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxPayOrder 更新VboxPayOrder
// @Tags VboxPayOrder
// @Summary 更新VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayOrder true "更新VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /vpo/updateVboxPayOrder [put]
func (vpoApi *VboxPayOrderApi) UpdateVboxPayOrder(c *gin.Context) {
	var vpo vbox.VboxPayOrder
	err := c.ShouldBindJSON(&vpo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    vpo.UpdatedBy = utils.GetUserID(c)
	if err := vpoService.UpdateVboxPayOrder(vpo); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxPayOrder 用id查询VboxPayOrder
// @Tags VboxPayOrder
// @Summary 用id查询VboxPayOrder
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxPayOrder true "用id查询VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vpo/findVboxPayOrder [get]
func (vpoApi *VboxPayOrderApi) FindVboxPayOrder(c *gin.Context) {
	var vpo vbox.VboxPayOrder
	err := c.ShouldBindQuery(&vpo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if revpo, err := vpoService.GetVboxPayOrder(vpo.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"revpo": revpo}, c)
	}
}

// GetVboxPayOrderList 分页获取VboxPayOrder列表
// @Tags VboxPayOrder
// @Summary 分页获取VboxPayOrder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxPayOrderSearch true "分页获取VboxPayOrder列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getVboxPayOrderList [get]
func (vpoApi *VboxPayOrderApi) GetVboxPayOrderList(c *gin.Context) {
	var pageInfo vboxReq.VboxPayOrderSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := vpoService.GetVboxPayOrderInfoList(pageInfo); err != nil {
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
