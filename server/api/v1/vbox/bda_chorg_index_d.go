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

type BdaChorgIndexDApi struct {
}

var bdaChorgService = service.ServiceGroupApp.VboxServiceGroup.BdaChorgIndexDService

// CreateBdaChorgIndexD 创建通道团队统计-天更新
// @Tags BdaChorgIndexD
// @Summary 创建通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChorgIndexD true "创建通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bdaChorg/createBdaChorgIndexD [post]
func (bdaChorgApi *BdaChorgIndexDApi) CreateBdaChorgIndexD(c *gin.Context) {
	var bdaChorg vbox.BdaChorgIndexD
	err := c.ShouldBindJSON(&bdaChorg)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChorg.CreatedBy = utils.GetUserID(c)
	if err := bdaChorgService.CreateBdaChorgIndexD(&bdaChorg); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteBdaChorgIndexD 删除通道团队统计-天更新
// @Tags BdaChorgIndexD
// @Summary 删除通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChorgIndexD true "删除通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChorg/deleteBdaChorgIndexD [delete]
func (bdaChorgApi *BdaChorgIndexDApi) DeleteBdaChorgIndexD(c *gin.Context) {
	var bdaChorg vbox.BdaChorgIndexD
	err := c.ShouldBindJSON(&bdaChorg)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChorg.DeletedBy = utils.GetUserID(c)
	if err := bdaChorgService.DeleteBdaChorgIndexD(bdaChorg); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteBdaChorgIndexDByIds 批量删除通道团队统计-天更新
// @Tags BdaChorgIndexD
// @Summary 批量删除通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /bdaChorg/deleteBdaChorgIndexDByIds [delete]
func (bdaChorgApi *BdaChorgIndexDApi) DeleteBdaChorgIndexDByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := bdaChorgService.DeleteBdaChorgIndexDByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateBdaChorgIndexD 更新通道团队统计-天更新
// @Tags BdaChorgIndexD
// @Summary 更新通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChorgIndexD true "更新通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChorg/updateBdaChorgIndexD [put]
func (bdaChorgApi *BdaChorgIndexDApi) UpdateBdaChorgIndexD(c *gin.Context) {
	var bdaChorg vbox.BdaChorgIndexD
	err := c.ShouldBindJSON(&bdaChorg)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChorg.UpdatedBy = utils.GetUserID(c)
	if err := bdaChorgService.UpdateBdaChorgIndexD(bdaChorg); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindBdaChorgIndexD 用id查询通道团队统计-天更新
// @Tags BdaChorgIndexD
// @Summary 用id查询通道团队统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.BdaChorgIndexD true "用id查询通道团队统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChorg/findBdaChorgIndexD [get]
func (bdaChorgApi *BdaChorgIndexDApi) FindBdaChorgIndexD(c *gin.Context) {
	var bdaChorg vbox.BdaChorgIndexD
	err := c.ShouldBindQuery(&bdaChorg)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rebdaChorg, err := bdaChorgService.GetBdaChorgIndexD(bdaChorg.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rebdaChorg": rebdaChorg}, c)
	}
}

// GetBdaChorgIndexDList 分页获取通道团队统计-天更新列表
// @Tags BdaChorgIndexD
// @Summary 分页获取通道团队统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.BdaChorgIndexDSearch true "分页获取通道团队统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChorg/getBdaChorgIndexDList [get]
func (bdaChorgApi *BdaChorgIndexDApi) GetBdaChorgIndexDList(c *gin.Context) {
	var pageInfo vboxReq.BdaChorgIndexDSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := bdaChorgService.GetBdaChorgIndexDInfoList(pageInfo); err != nil {
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

// GetBdaChorgIndexRealList 分页获取通道团队统计-天更新列表
// @Tags BdaChorgIndexD
// @Summary 分页获取通道团队统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.BdaChorgIndexDSearch true "分页获取通道团队统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChorg/getBdaChorgIndexRealList [get]
func (bdaChorgApi *BdaChorgIndexDApi) GetBdaChorgIndexRealList(c *gin.Context) {
	var pageInfo vboxReq.OrgSelectForm
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pageInfo.Uid = utils.GetUserID(c)
	if list, total, err := bdaChorgService.GetBdaChorgIndexRealList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "获取成功", c)
	}
}

// GetBdaChorgIndexRealListBySelect 分页获取通道团队统计-天更新列表
// @Tags BdaChorgIndexD
// @Summary 分页获取通道团队统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.BdaChorgIndexDSearch true "分页获取通道团队统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChorg/getBdaChorgIndexRealListBySelect [get]
func (bdaChorgApi *BdaChorgIndexDApi) GetBdaChorgIndexRealListBySelect(c *gin.Context) {
	var pageInfo vboxReq.OrgSelectForm
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	pageInfo.Uid = utils.GetUserID(c)
	if list, total, err := bdaChorgService.GetBdaChorgIndexRealListBySelect(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "获取成功", c)
	}
}
