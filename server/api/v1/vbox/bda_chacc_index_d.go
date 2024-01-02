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

type BdaChaccIndexDApi struct {
}

var bdaChaccIndexDService = service.ServiceGroupApp.VboxServiceGroup.BdaChaccIndexDService

// CreateBdaChaccIndexD 创建用户通道粒度成率统计-天更新
// @Tags BdaChaccIndexD
// @Summary 创建用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChaccIndexD true "创建用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bdaChaccIndexD/createBdaChaccIndexD [post]
func (bdaChaccIndexDApi *BdaChaccIndexDApi) CreateBdaChaccIndexD(c *gin.Context) {
	var bdaChaccIndexD vbox.BdaChaccIndexD
	err := c.ShouldBindJSON(&bdaChaccIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChaccIndexD.CreatedBy = utils.GetUserID(c)
	if err := bdaChaccIndexDService.CreateBdaChaccIndexD(&bdaChaccIndexD); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteBdaChaccIndexD 删除用户通道粒度成率统计-天更新
// @Tags BdaChaccIndexD
// @Summary 删除用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChaccIndexD true "删除用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChaccIndexD/deleteBdaChaccIndexD [delete]
func (bdaChaccIndexDApi *BdaChaccIndexDApi) DeleteBdaChaccIndexD(c *gin.Context) {
	var bdaChaccIndexD vbox.BdaChaccIndexD
	err := c.ShouldBindJSON(&bdaChaccIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChaccIndexD.DeletedBy = utils.GetUserID(c)
	if err := bdaChaccIndexDService.DeleteBdaChaccIndexD(bdaChaccIndexD); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteBdaChaccIndexDByIds 批量删除用户通道粒度成率统计-天更新
// @Tags BdaChaccIndexD
// @Summary 批量删除用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /bdaChaccIndexD/deleteBdaChaccIndexDByIds [delete]
func (bdaChaccIndexDApi *BdaChaccIndexDApi) DeleteBdaChaccIndexDByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := bdaChaccIndexDService.DeleteBdaChaccIndexDByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateBdaChaccIndexD 更新用户通道粒度成率统计-天更新
// @Tags BdaChaccIndexD
// @Summary 更新用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChaccIndexD true "更新用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChaccIndexD/updateBdaChaccIndexD [put]
func (bdaChaccIndexDApi *BdaChaccIndexDApi) UpdateBdaChaccIndexD(c *gin.Context) {
	var bdaChaccIndexD vbox.BdaChaccIndexD
	err := c.ShouldBindJSON(&bdaChaccIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChaccIndexD.UpdatedBy = utils.GetUserID(c)
	if err := bdaChaccIndexDService.UpdateBdaChaccIndexD(bdaChaccIndexD); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindBdaChaccIndexD 用id查询用户通道粒度成率统计-天更新
// @Tags BdaChaccIndexD
// @Summary 用id查询用户通道粒度成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.BdaChaccIndexD true "用id查询用户通道粒度成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChaccIndexD/findBdaChaccIndexD [get]
func (bdaChaccIndexDApi *BdaChaccIndexDApi) FindBdaChaccIndexD(c *gin.Context) {
	var bdaChaccIndexD vbox.BdaChaccIndexD
	err := c.ShouldBindQuery(&bdaChaccIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rebdaChaccIndexD, err := bdaChaccIndexDService.GetBdaChaccIndexD(bdaChaccIndexD.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rebdaChaccIndexD": rebdaChaccIndexD}, c)
	}
}

// GetBdaChaccIndexDList 分页获取用户通道粒度成率统计-天更新列表
// @Tags BdaChaccIndexD
// @Summary 分页获取用户通道粒度成率统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.BdaChaccIndexDSearch true "分页获取用户通道粒度成率统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChaccIndexD/getBdaChaccIndexDList [get]
func (bdaChaccIndexDApi *BdaChaccIndexDApi) GetBdaChaccIndexDList(c *gin.Context) {
	var pageInfo vboxReq.BdaChaccIndexDSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := bdaChaccIndexDService.GetBdaChaccIndexDInfoList(pageInfo); err != nil {
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

// CronVboxBdaChaccIndexDByHand 分页获取用户通道粒度成率统计-天更新列表
// @Tags BdaChIndexD
// @Summary 分页获取用户通道粒度成率统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.BdaChIndexDSearch true "分页获取用户通道粒度成率统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChaccIndexD/CronVboxBdaChaccIndexDByHand [get]
func (bdaChaccIndexDApi *BdaChaccIndexDApi) CronVboxBdaChaccIndexDByHand(c *gin.Context) {
	var pageInfo vboxReq.BdaChIndexDSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := bdaChaccIndexDService.CronVboxBdaChaccIndexDByHand(pageInfo.Dt); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithMessage("调度成功", c)
	}
}
