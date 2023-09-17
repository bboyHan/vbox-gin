package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VboxPayOrderApi struct {
}

// CreateOrder2PayAcc 创建CreateOrder2PayAcc
// @Tags VboxPayOrder
// @Summary 创建Order2PayAcc
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxPayOrder true "创建VboxPayOrder"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /base/create [post]
func (vpoApi *VboxPayOrderApi) CreateOrder2PayAcc(c *gin.Context) {
	var vpo vboxReq.CreateOrder2PayAccount
	err := c.ShouldBind(&vpo) // 可接收 from - json - xml
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if order, err := vpoService.CreateOrder2PayAcc(&vpo); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
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
func (vpoApi *VboxPayOrderApi) QueryOrder2PayAcc(c *gin.Context) {
	var vpo vboxReq.QueryOrder2PayAccount
	err := c.ShouldBind(&vpo) // 可接收 from - json - xml
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if order, err := vpoService.QueryOrder2PayAcc(&vpo); err != nil {
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
func (vpoApi *VboxPayOrderApi) QueryOrderSimple(c *gin.Context) {
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

	if order, err := vpoService.QueryOrderSimple(&vpo); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithDetailed(order, "查询成功", c)
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

	userID := utils.GetUserID(c)
	userList, tot, err := userService.GetOwnerUserIdsList(userID)
	var idList []int
	for _, user := range userList {
		idList = append(idList, int(user.ID))
	}
	if err != nil || tot == 0 {
		return
	}
	if list, total, err := vpoService.GetVboxPayOrderInfoList(pageInfo, idList); err != nil {
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

// GetVboxUserPayOrderAnalysis 获取用户订单看板
// @Tags VboxPayOrder
// @Summary 分页获取VboxPayOrder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxPayOrderSearch true "获取用户订单看板"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/GetVboxUserPayOrderAnalysis [get]
func (vpoApi *VboxPayOrderApi) GetVboxUserPayOrderAnalysis(c *gin.Context) {

	userID := utils.GetUserID(c)
	userList, tot, err := userService.GetOwnerUserIdsList(userID)
	var idList []int
	for _, user := range userList {
		idList = append(idList, int(user.ID))
	}
	if err != nil || tot == 0 {
		return
	}
	if list, total, err := vpoService.GetVboxUserPayOrderAnalysis(userID, idList); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "获取成功", c)
	}
}

// GetSelectUserPayOrderAnalysis 获取单个用户分析展示数据
// @Tags VboxPayOrder
// @Summary 分页获取VboxPayOrder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxPayOrderSearch true "获取单个用户分析展示数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getSelectUserPayOrderAnalysis [get]
func (vpoApi *VboxPayOrderApi) GetSelectUserPayOrderAnalysis(c *gin.Context) {
	userID := utils.GetUserID(c)
	var selectUser system.SysUser
	err := c.ShouldBindQuery(&selectUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 创建db
	db := global.GVA_DB.Model(&system.SysUser{})
	var user system.SysUser
	db.Where("username = ?", selectUser.Username).Find(&user)
	var idList []int
	idList = append(idList, int(user.ID))

	if list, total, err := vpoService.GetVboxUserPayOrderAnalysis(userID, idList); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:  list,
			Total: total,
		}, "获取成功", c)
	}
}

// GetSelectUserPayOrderAnalysisH 获取单个用户分析展示数据
// @Tags VboxPayOrder
// @Summary 分页获取VboxPayOrder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxPayOrderSearch true "获取单个用户分析展示数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/GetSelectUserPayOrderAnalysisH [get]
func (vpoApi *VboxPayOrderApi) GetSelectUserPayOrderAnalysisH(c *gin.Context) {
	userID := utils.GetUserID(c)
	var idList []int
	idList = append(idList, int(userID))

	if resData, err := vpoService.GetHomePagePayOrderAnalysisH(userID, idList); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {

		response.OkWithData(gin.H{"resultData": resData}, c)
	}
}

// GetHomePagePayOrderAnalysis 获取首页单个用户分析展示数据
// @Tags VboxPayOrder
// @Summary 分页获取VboxPayOrder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxPayOrderSearch true "获取首页单个用户分析展示数据"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getHomePagePayOrderAnalysis [get]
func (vpoApi *VboxPayOrderApi) GetHomePagePayOrderAnalysis(c *gin.Context) {
	userID := utils.GetUserID(c)
	var idList []int
	idList = append(idList, int(userID))

	if resData, err := vpoService.GetHomePagePayOrderAnalysis(userID, idList); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {

		response.OkWithData(gin.H{"resultData": resData}, c)
	}
}

// GetSelectPayOrderAnalysisQuantifyCharts 获取单个用户各个通道下每天的成单数据图
// @Tags VboxPayOrder
// @Summary 分页获取VboxPayOrder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxPayOrderSearch true "获取单个用户各个通道下每天的成单数据图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getSelectPayOrderAnalysisQuantifyCharts [get]
func (vpoApi *VboxPayOrderApi) GetSelectPayOrderAnalysisQuantifyCharts(c *gin.Context) {
	var selectUser system.SysUser
	err := c.ShouldBindQuery(&selectUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 创建db
	db := global.GVA_DB.Model(&system.SysUser{})
	var user system.SysUser
	db.Where("username = ?", selectUser.Username).Find(&user)

	if data, err := vpoService.GetSelectPayOrderAnalysisQuantifyCharts(int(user.ID)); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}

// GetSelectPayOrderAnalysisChannelIncomeCharts 获取单个用户各个通道下每天的收入数据图
// @Tags VboxPayOrder
// @Summary 分页获取VboxPayOrder列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxPayOrderSearch true "获取单个用户各个通道下每天的收入数据图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getSelectPayOrderAnalysisChannelIncomeCharts [get]
func (vpoApi *VboxPayOrderApi) GetSelectPayOrderAnalysisChannelIncomeCharts(c *gin.Context) {
	var selectUser system.SysUser
	err := c.ShouldBindQuery(&selectUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 创建db
	db := global.GVA_DB.Model(&system.SysUser{})
	var user system.SysUser
	db.Where("username = ?", selectUser.Username).Find(&user)

	if data, err := vpoService.GetSelectPayOrderAnalysisChannelIncomeCharts(int(user.ID)); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}

// GetSelectPayOrderAnalysisIncomeBarCharts
// @Summary 获取单个用户每天的收入数据图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query request.PageInfo true "获取单个用户各个通道下每天的成单数据图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getSelectPayOrderAnalysisIncomeBarCharts [get]
func (vpoApi *VboxPayOrderApi) GetSelectPayOrderAnalysisIncomeBarCharts(c *gin.Context) {
	var selectUser system.SysUser
	err := c.ShouldBindQuery(&selectUser)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 创建db
	db := global.GVA_DB.Model(&system.SysUser{})
	var user system.SysUser
	db.Where("username = ?", selectUser.Username).Find(&user)

	if data, err := vpoService.GetSelectPayOrderAnalysisIncomeBarCharts(int(user.ID)); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}

// GetVboxUserPayOrderAnalysisIncomeCharts 获取用户订单看板收入图
// @Tags VboxPayOrder
// @Summary 获取用户订单看板收入图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxPayOrderSearch true "获取用户订单看板收入图"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /vpo/getVboxUserPayOrderAnalysisIncomeCharts [get]
func (vpoApi *VboxPayOrderApi) GetVboxUserPayOrderAnalysisIncomeCharts(c *gin.Context) {

	userID := utils.GetUserID(c)
	userList, tot, err := userService.GetOwnerUserIdsList(userID)
	var idList []int
	for _, user := range userList {
		idList = append(idList, int(user.ID))
	}
	if err != nil || tot == 0 {
		return
	}
	if data, err := vpoService.GetVboxUserPayOrderAnalysisIncomeCharts(userID, idList); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(data, "获取成功", c)
	}
}
