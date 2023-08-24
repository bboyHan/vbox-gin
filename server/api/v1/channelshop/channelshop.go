package channelshop

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/channelshop"
	channelshopReq "github.com/flipped-aurora/gin-vue-admin/server/model/channelshop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ChannelShopApi struct {
}

var chShopService = service.ServiceGroupApp.ChannelshopServiceGroup.ChannelShopService

// @Tags ChannelShop
// @Summary 批量更新ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ChannelShop true "批量更新ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chShop/batchUpdateChannelShopStatus [post]
func (chShopApi *ChannelShopApi) BatchUpdateChannelShopStatus(c *gin.Context) {
	userId := utils.GetUserID(c)
	var chShop channelshop.ChannelShop
	err := c.ShouldBindJSON(&chShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chShop.CreatedBy = utils.GetUserID(c)
	if err := chShopService.BatchUpdateChannelShopStatus(chShop, userId); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// CreateChannelShop 创建ChannelShop
// @Tags ChannelShop
// @Summary 创建ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body channelshop.ChannelShop true "创建ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chShop/createChannelShop [post]
func (chShopApi *ChannelShopApi) CreateChannelShop(c *gin.Context) {
	var chShop channelshop.ChannelShop
	err := c.ShouldBindJSON(&chShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chShop.CreatedBy = utils.GetUserID(c)
	if err := chShopService.CreateChannelShop(&chShop); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteChannelShop 删除ChannelShop
// @Tags ChannelShop
// @Summary 删除ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body channelshop.ChannelShop true "删除ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /chShop/deleteChannelShop [delete]
func (chShopApi *ChannelShopApi) DeleteChannelShop(c *gin.Context) {
	var chShop channelshop.ChannelShop
	err := c.ShouldBindJSON(&chShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chShop.DeletedBy = utils.GetUserID(c)
	if err := chShopService.DeleteChannelShop(chShop); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteChannelShopByIds 批量删除ChannelShop
// @Tags ChannelShop
// @Summary 批量删除ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /chShop/deleteChannelShopByIds [delete]
func (chShopApi *ChannelShopApi) DeleteChannelShopByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := chShopService.DeleteChannelShopByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateChannelShop 更新ChannelShop
// @Tags ChannelShop
// @Summary 更新ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body channelshop.ChannelShop true "更新ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /chShop/updateChannelShop [put]
func (chShopApi *ChannelShopApi) UpdateChannelShop(c *gin.Context) {
	var chShop channelshop.ChannelShop
	err := c.ShouldBindJSON(&chShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	chShop.UpdatedBy = utils.GetUserID(c)
	if err := chShopService.UpdateChannelShop(chShop); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindChannelShop 用id查询ChannelShop
// @Tags ChannelShop
// @Summary 用id查询ChannelShop
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query channelshop.ChannelShop true "用id查询ChannelShop"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /chShop/findChannelShop [get]
func (chShopApi *ChannelShopApi) FindChannelShop(c *gin.Context) {
	var chShop channelshop.ChannelShop
	err := c.ShouldBindQuery(&chShop)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rechShop, err := chShopService.GetChannelShop(chShop.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rechShop": rechShop}, c)
	}
}

// GetChannelShopList 分页获取ChannelShop列表
// @Tags ChannelShop
// @Summary 分页获取ChannelShop列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query channelshopReq.ChannelShopSearch true "分页获取ChannelShop列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chShop/getChannelShopList [get]
func (chShopApi *ChannelShopApi) GetChannelShopList(c *gin.Context) {
	userId := utils.GetUserID(c)
	fmt.Println("GetChannelShopList userId = ", userId)
	var pageInfo channelshopReq.ChannelShopSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := chShopService.GetTreeChannelShopInfoList(pageInfo, userId); err != nil {
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

	//if list, total, err := chShopService.GetChannelShopInfoList(pageInfo, userId); err != nil {
	//	global.GVA_LOG.Error("获取失败!", zap.Error(err))
	//	response.FailWithMessage("获取失败", c)
	//} else {
	//	response.OkWithDetailed(response.PageResult{
	//		List:     list,
	//		Total:    total,
	//		Page:     pageInfo.Page,
	//		PageSize: pageInfo.PageSize,
	//	}, "获取成功", c)
	//}
}

// GetChannelShopListByChanelRemark 分页获取获取同一通道下同一店铺的ChannelShop列表
// @Tags ChannelShop
// @Summary 分页获取获取同一通道下同一店铺的ChannelShop列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query channelshopReq.ChannelShopSearch true "分页获取获取同一通道下同一店铺的ChannelShop列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /chShop/getChannelShopListByChanelRemark [get]
func (chShopApi *ChannelShopApi) GetChannelShopListByChanelRemark(c *gin.Context) {
	userId := utils.GetUserID(c)
	fmt.Println("GetChannelShopList userId = ", userId)
	var pageInfo channelshopReq.ChannelShopSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := chShopService.GetChannelShopInfoListByChanelRemark(pageInfo, userId); err != nil {
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
