package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BdaChIndexDApi struct {
}

var bdaChIndexDService = service.ServiceGroupApp.VboxServiceGroup.BdaChIndexDService

// CreateBdaChIndexD 创建用户通道粒度成率统计-天更新
// @Tags BdaChIndexD
// @Summary 创建用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChIndexD true "创建用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bdaChIndexD/createBdaChIndexD [post]
func (bdaChIndexDApi *BdaChIndexDApi) CreateBdaChIndexD(c *gin.Context) {
	var bdaChIndexD vbox.BdaChIndexD
	err := c.ShouldBindJSON(&bdaChIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChIndexD.CreatedBy = utils.GetUserID(c)
	if err := bdaChIndexDService.CreateBdaChIndexD(&bdaChIndexD); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteBdaChIndexD 删除用户通道粒度成率统计-天更新
// @Tags BdaChIndexD
// @Summary 删除用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChIndexD true "删除用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChIndexD/deleteBdaChIndexD [delete]
func (bdaChIndexDApi *BdaChIndexDApi) DeleteBdaChIndexD(c *gin.Context) {
	var bdaChIndexD vbox.BdaChIndexD
	err := c.ShouldBindJSON(&bdaChIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChIndexD.DeletedBy = utils.GetUserID(c)
	if err := bdaChIndexDService.DeleteBdaChIndexD(bdaChIndexD); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteBdaChIndexDByIds 批量删除用户通道粒度成率统计-天更新
// @Tags BdaChIndexD
// @Summary 批量删除用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /bdaChIndexD/deleteBdaChIndexDByIds [delete]
func (bdaChIndexDApi *BdaChIndexDApi) DeleteBdaChIndexDByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := bdaChIndexDService.DeleteBdaChIndexDByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateBdaChIndexD 更新用户通道粒度成率统计-天更新
// @Tags BdaChIndexD
// @Summary 更新用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChIndexD true "更新用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChIndexD/updateBdaChIndexD [put]
func (bdaChIndexDApi *BdaChIndexDApi) UpdateBdaChIndexD(c *gin.Context) {
	var bdaChIndexD vbox.BdaChIndexD
	err := c.ShouldBindJSON(&bdaChIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChIndexD.UpdatedBy = utils.GetUserID(c)
	if err := bdaChIndexDService.UpdateBdaChIndexD(bdaChIndexD); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindBdaChIndexD 用id查询用户通道粒度成率统计-天更新
// @Tags BdaChIndexD
// @Summary 用id查询用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.BdaChIndexD true "用id查询用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChIndexD/findBdaChIndexD [get]
func (bdaChIndexDApi *BdaChIndexDApi) FindBdaChIndexD(c *gin.Context) {
	var bdaChIndexD vbox.BdaChIndexD
	err := c.ShouldBindQuery(&bdaChIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rebdaChIndexD, err := bdaChIndexDService.GetBdaChIndexD(bdaChIndexD.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rebdaChIndexD": rebdaChIndexD}, c)
	}
}

// GetBdaChIndexDList 分页获取用户通道粒度成率统计-天更新列表
// @Tags BdaChIndexD
// @Summary 分页获取用户通道粒度成率统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.BdaChIndexDSearch true "分页获取用户通道粒度成率统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChIndexD/getBdaChIndexDList [get]
func (bdaChIndexDApi *BdaChIndexDApi) GetBdaChIndexDList(c *gin.Context) {
	var pageInfo vboxReq.BdaChIndexDSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := bdaChIndexDService.GetBdaChIndexDInfoList(pageInfo); err != nil {
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

// CronVboxBdaChIndexDByHand 分页获取用户通道粒度成率统计-天更新列表
// @Tags BdaChIndexD
// @Summary 分页获取用户通道粒度成率统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.BdaChIndexDSearch true "分页获取用户通道粒度成率统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChIndexD/cronVboxBdaChIndexDByHand [get]
func (bdaChIndexDApi *BdaChIndexDApi) CronVboxBdaChIndexDByHand(c *gin.Context) {
	var pageInfo vboxReq.BdaChIndexDSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bdaChIndexDService.CronVboxBdaChIndexDByHand(pageInfo.Dt); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithMessage("调度成功", c)
	}
}
