package vbox

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/captcha"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PayOrderApi struct {
}

var payOrderService = service.ServiceGroupApp.VboxServiceGroup.PayOrderService

// CreateOrder2PayAcc 创建CreateOrder2PayAcc
// @Tags VboxPayOrder
// @Summary 创建Order2PayAcc
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayOrder true "创建VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /base/create [post]
func (vpoApi *PayOrderApi) CreateOrder2PayAcc(c *gin.Context) {
	var vpo vboxReq.CreateOrder2PayAccount
	err := c.ShouldBind(&vpo) // 可接收 from - json - xml
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if order, err := payOrderService.CreateOrder2PayAcc(&vpo); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(order, "创建成功", c)
	}
}

// CreateOrderTest 创建CreateOrder2Acc
// @Tags CreateOrderTest
// @Summary 创建CreateOrderTest
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayOrder true "创建CreateOrderTest"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /base/create [post]
func (vpoApi *PayOrderApi) CreateOrderTest(c *gin.Context) {
	var vpo vboxReq.CreateOrderTest
	err := c.ShouldBind(&vpo) // 可接收 from - json - xml
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userID := utils.GetUserID(c)
	var user system.SysUser

	err = global.GVA_DB.Where("id = ?", userID).First(&user).Error
	if err == nil {
		if user.EnableAuth == 1 {
			var secret string
			secret, _ = captcha.GetSecret(user.AuthCaptcha)
			if ok := captcha.ValidateCode(secret, vpo.AuthCaptcha); !ok {
				err = errors.New("双因子认证码错误")
				response.FailWithMessage(err.Error(), c)
				return
			}
		} else {
			if vpo.AuthCaptcha != "666" {
				err = errors.New("双因子认证码错误")
				response.FailWithMessage(err.Error(), c)
				return
			}
		}
	}

	vpo.Username = user.Username
	vpo.UserId = user.ID
	if order, err := payOrderService.CreateOrderTest(&vpo); err != nil {
		global.GVA_LOG.Error("测试订单创建失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(order, "创建成功", c)
	}
}

// QueryOrder2PayAcc 查询QueryOrder2PayAcc
// @Tags VboxPayOrder
// @Summary 查询QueryOrder2PayAcc
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayOrder true "查询QueryOrder2PayAcc"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /base/create [post]
func (vpoApi *PayOrderApi) QueryOrder2PayAcc(c *gin.Context) {
	var vpo vboxReq.QueryOrder2PayAccount
	err := c.ShouldBind(&vpo) // 可接收 from - json - xml
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if order, err := payOrderService.QueryOrder2PayAcc(&vpo); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(order, "查询成功", c)
	}
}

// QueryOrderSimple 查询QueryOrderSimple
// @Tags VboxPayOrder
// @Summary 查询QueryOrderSimple
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayOrder true "查询QueryOrderSimple"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /base/create [post]
func (vpoApi *PayOrderApi) QueryOrderSimple(c *gin.Context) {
	var vpo vboxReq.QueryOrderSimple
	err := c.ShouldBind(&vpo) // 可接收 from - json - xml
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	vpo.PayIp = c.ClientIP()
	vpo.UserAgent = c.Request.UserAgent()
	vpo.PayRegion, _ = utils.SearchIp2Region(vpo.PayIp)
	vpo.PayDevice = utils.GetDeviceSimpleInfo(vpo.UserAgent)

	if order, err := payOrderService.QueryOrderSimple(&vpo); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(order, "查询成功", c)
	}
}

// FindPayOrder 用id查询订单
// @Tags PayOrder
// @Summary 用id查询订单
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.PayOrder true "用id查询订单"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /payOrder/findPayOrder [get]
func (vpoApi *PayOrderApi) FindPayOrder(c *gin.Context) {
	var payOrder vbox.PayOrder
	err := c.ShouldBindQuery(&payOrder)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if repayOrder, err := payOrderService.GetPayOrder(payOrder.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"repayOrder": repayOrder}, c)
	}
}

// GetPayOrderList 分页获取订单列表
// @Tags PayOrder
// @Summary 分页获取订单列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.PayOrderSearch true "分页获取订单列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /payOrder/getPayOrderList [get]
func (vpoApi *PayOrderApi) GetPayOrderList(c *gin.Context) {
	var pageInfo vboxReq.PayOrderSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := payOrderService.GetPayOrderInfoList(pageInfo); err != nil {
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
