package vbox

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type VboxUserWalletApi struct {
}

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
	// 保存被分配积分的用户数据
	// 创建db
	db := global.GVA_DB.Model(&system.SysUser{})
	var user system.SysUser
	db.Where("username = ?", vuw.UserName).Find(&user)
	vuw.Uid = user.ID
	output := fmt.Sprintf("划转至【%d】,积分：【%d】", vuw.UserName, vuw.Recharge)
	vuw.Remark = output
	fmt.Println("id = ", vuw.ID)
	if err := vuwService.CreateVboxUserWallet(&vuw); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	}

	// 保存分配积分的用户数据
	userId := utils.GetUserID(c)
	var vuwUser vbox.VboxUserWallet
	vuwUser = vuw

	vuwUser.Uid = userId
	// 创建db
	createDb := global.GVA_DB.Model(&system.SysUser{})
	var createUser system.SysUser
	createDb.Where("id = ?", userId).Find(&createUser)
	vuwUser.UserName = createUser.Username
	output2 := fmt.Sprintf("划转至:【%s】,积分:【%d】", vuwUser.UserName, vuwUser.Recharge)
	vuwUser.Remark = output2
	vuwUser.Recharge = -1 * vuw.Recharge
	vuwUser.ID = 0
	if err := vuwService.CreateVboxUserWallet(&vuwUser); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	}

	response.OkWithMessage("创建成功", c)

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
	if err := vuwService.DeleteVboxUserWalletByIds(IDS, deletedBy); err != nil {
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
	userId := uint(utils.GetUserID(c))
	var pageInfo vboxReq.VboxUserWalletSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userList, tot, err := userService.GetOwnerUserIdsList(userId)
	var idList []int
	for _, user := range userList {
		idList = append(idList, int(user.ID))
	}
	if err != nil || tot == 0 {
		return
	}

	if list, total, err := vuwService.GetVboxUserWalletInfoList(pageInfo, idList); err != nil {
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

// GetVboxUserWalletAvailablePoints 用id查询VboxUserWallet
// @Tags VboxUserWallet
// @Summary 用id查询VboxUserWallet
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxUserWallet true "用id查询VboxUserWallet"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /vuw/GetVboxUserWalletAvailablePoints [get]
func (vuwApi *VboxUserWalletApi) GetVboxUserWalletAvailablePoints(c *gin.Context) {
	userId := uint(utils.GetUserID(c))
	var vuw vbox.VboxUserWallet
	err := c.ShouldBindQuery(&vuw)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userList, tot, err := userService.GetOwnerUserIdsList(userId)
	var idList []int
	for _, user := range userList {
		idList = append(idList, int(user.ID))
	}
	if err != nil || tot == 0 {
		return
	}

	if rechargeData, err := vuwService.GetVboxUserWalletAvailablePoints(userId, idList); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rechargeData": rechargeData}, c)
	}
}
