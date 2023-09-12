package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VboxBdaChIndexDApi struct {
}

// CreateVboxBdaChIndexD 创建VboxBdaChIndexD
// @Tags VboxBdaChIndexD
// @Summary 创建VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxBdaChIndexD true "创建VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChD/createVboxBdaChIndexD [post]
func (bdaChDApi *VboxBdaChIndexDApi) CreateVboxBdaChIndexD(c *gin.Context) {
	var bdaChD vbox.VboxBdaChIndexD
	err := c.ShouldBindJSON(&bdaChD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChD.CreatedBy = utils.GetUserID(c)
	verify := utils.Rules{
		"Dt": {utils.NotEmpty()},
	}
	if err := utils.Verify(bdaChD, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bdaChDService.CreateVboxBdaChIndexD(&bdaChD); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxBdaChIndexD 删除VboxBdaChIndexD
// @Tags VboxBdaChIndexD
// @Summary 删除VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxBdaChIndexD true "删除VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChD/deleteVboxBdaChIndexD [delete]
func (bdaChDApi *VboxBdaChIndexDApi) DeleteVboxBdaChIndexD(c *gin.Context) {
	var bdaChD vbox.VboxBdaChIndexD
	err := c.ShouldBindJSON(&bdaChD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChD.DeletedBy = utils.GetUserID(c)
	if err := bdaChDService.DeleteVboxBdaChIndexD(bdaChD); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxBdaChIndexDByIds 批量删除VboxBdaChIndexD
// @Tags VboxBdaChIndexD
// @Summary 批量删除VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /bdaChD/deleteVboxBdaChIndexDByIds [delete]
func (bdaChDApi *VboxBdaChIndexDApi) DeleteVboxBdaChIndexDByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := bdaChDService.DeleteVboxBdaChIndexDByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxBdaChIndexD 更新VboxBdaChIndexD
// @Tags VboxBdaChIndexD
// @Summary 更新VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxBdaChIndexD true "更新VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChD/updateVboxBdaChIndexD [put]
func (bdaChDApi *VboxBdaChIndexDApi) UpdateVboxBdaChIndexD(c *gin.Context) {
	var bdaChD vbox.VboxBdaChIndexD
	err := c.ShouldBindJSON(&bdaChD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChD.UpdatedBy = utils.GetUserID(c)
	verify := utils.Rules{
		"Dt": {utils.NotEmpty()},
	}
	if err := utils.Verify(bdaChD, verify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bdaChDService.UpdateVboxBdaChIndexD(bdaChD); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxBdaChIndexD 用id查询VboxBdaChIndexD
// @Tags VboxBdaChIndexD
// @Summary 用id查询VboxBdaChIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxBdaChIndexD true "用id查询VboxBdaChIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChD/findVboxBdaChIndexD [get]
func (bdaChDApi *VboxBdaChIndexDApi) FindVboxBdaChIndexD(c *gin.Context) {
	var bdaChD vbox.VboxBdaChIndexD
	err := c.ShouldBindQuery(&bdaChD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rebdaChD, err := bdaChDService.GetVboxBdaChIndexD(bdaChD.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rebdaChD": rebdaChD}, c)
	}
}

// GetVboxBdaChIndexDList 分页获取VboxBdaChIndexD列表
// @Tags VboxBdaChIndexD
// @Summary 分页获取VboxBdaChIndexD列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxBdaChIndexDSearch true "分页获取VboxBdaChIndexD列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChD/getVboxBdaChIndexDList [get]
func (bdaChDApi *VboxBdaChIndexDApi) GetVboxBdaChIndexDList(c *gin.Context) {
	var pageInfo vboxReq.VboxBdaChIndexDSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := bdaChDService.GetVboxBdaChIndexDInfoList(pageInfo); err != nil {
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
