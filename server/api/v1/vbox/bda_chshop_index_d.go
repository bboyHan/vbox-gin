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

type BdaChShopIndexDApi struct {
}

var bdaChshopIndexDService = service.ServiceGroupApp.VboxServiceGroup.BdaChShopIndexDService

// CreateBdaChShopIndexD 创建用户通道店铺成率统计-天更新
// @Tags BdaChShopIndexD
// @Summary 创建用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChShopIndexD true "创建用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"创建成功"}"
// @Router /bdaChshopIndexD/createBdaChShopIndexD [post]
func (bdaChshopIndexDApi *BdaChShopIndexDApi) CreateBdaChShopIndexD(c *gin.Context) {
	var bdaChshopIndexD vbox.BdaChShopIndexD
	err := c.ShouldBindJSON(&bdaChshopIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChshopIndexD.CreatedBy = utils.GetUserID(c)
	if err := bdaChshopIndexDService.CreateBdaChShopIndexD(&bdaChshopIndexD); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteBdaChShopIndexD 删除用户通道店铺成率统计-天更新
// @Tags BdaChShopIndexD
// @Summary 删除用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChShopIndexD true "删除用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChshopIndexD/deleteBdaChShopIndexD [delete]
func (bdaChshopIndexDApi *BdaChShopIndexDApi) DeleteBdaChShopIndexD(c *gin.Context) {
	var bdaChshopIndexD vbox.BdaChShopIndexD
	err := c.ShouldBindJSON(&bdaChshopIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChshopIndexD.DeletedBy = utils.GetUserID(c)
	if err := bdaChshopIndexDService.DeleteBdaChShopIndexD(bdaChshopIndexD); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteBdaChShopIndexDByIds 批量删除用户通道店铺成率统计-天更新
// @Tags BdaChShopIndexD
// @Summary 批量删除用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /bdaChshopIndexD/deleteBdaChShopIndexDByIds [delete]
func (bdaChshopIndexDApi *BdaChShopIndexDApi) DeleteBdaChShopIndexDByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := bdaChshopIndexDService.DeleteBdaChShopIndexDByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateBdaChShopIndexD 更新用户通道店铺成率统计-天更新
// @Tags BdaChShopIndexD
// @Summary 更新用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.BdaChShopIndexD true "更新用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChshopIndexD/updateBdaChShopIndexD [put]
func (bdaChshopIndexDApi *BdaChShopIndexDApi) UpdateBdaChShopIndexD(c *gin.Context) {
	var bdaChshopIndexD vbox.BdaChShopIndexD
	err := c.ShouldBindJSON(&bdaChshopIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	bdaChshopIndexD.UpdatedBy = utils.GetUserID(c)
	if err := bdaChshopIndexDService.UpdateBdaChShopIndexD(bdaChshopIndexD); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindBdaChShopIndexD 用id查询用户通道店铺成率统计-天更新
// @Tags BdaChShopIndexD
// @Summary 用id查询用户通道店铺成率统计-天更新
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.BdaChShopIndexD true "用id查询用户通道店铺成率统计-天更新"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChshopIndexD/findBdaChShopIndexD [get]
func (bdaChshopIndexDApi *BdaChShopIndexDApi) FindBdaChShopIndexD(c *gin.Context) {
	var bdaChshopIndexD vbox.BdaChShopIndexD
	err := c.ShouldBindQuery(&bdaChshopIndexD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rebdaChshopIndexD, err := bdaChshopIndexDService.GetBdaChShopIndexD(bdaChshopIndexD.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rebdaChshopIndexD": rebdaChshopIndexD}, c)
	}
}

// GetBdaChShopIndexDList 分页获取用户通道店铺成率统计-天更新列表
// @Tags BdaChShopIndexD
// @Summary 分页获取用户通道店铺成率统计-天更新列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.BdaChShopIndexDSearch true "分页获取用户通道店铺成率统计-天更新列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChshopIndexD/getBdaChShopIndexDList [get]
func (bdaChshopIndexDApi *BdaChShopIndexDApi) GetBdaChShopIndexDList(c *gin.Context) {
	var pageInfo vboxReq.BdaChShopIndexDSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := bdaChshopIndexDService.GetBdaChShopIndexDInfoList(pageInfo); err != nil {
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
