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

type VboxBdaChaccIndexDApi struct {
}

var bdaChaccDService = service.ServiceGroupApp.VboxServiceGroup.VboxBdaChaccIndexDService


// CreateVboxBdaChaccIndexD 创建VboxBdaChaccIndexD
// @Tags VboxBdaChaccIndexD
// @Summary 创建VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxBdaChaccIndexD true "创建VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChaccD/createVboxBdaChaccIndexD [post]
func (bdaChaccDApi *VboxBdaChaccIndexDApi) CreateVboxBdaChaccIndexD(c *gin.Context) {
	var bdaChaccD vbox.VboxBdaChaccIndexD
	err := c.ShouldBindJSON(&bdaChaccD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    bdaChaccD.CreatedBy = utils.GetUserID(c)
    verify := utils.Rules{
        "Dt":{utils.NotEmpty()},
    }
	if err := utils.Verify(bdaChaccD, verify); err != nil {
    		response.FailWithMessage(err.Error(), c)
    		return
    	}
	if err := bdaChaccDService.CreateVboxBdaChaccIndexD(&bdaChaccD); err != nil {
        global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteVboxBdaChaccIndexD 删除VboxBdaChaccIndexD
// @Tags VboxBdaChaccIndexD
// @Summary 删除VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxBdaChaccIndexD true "删除VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /bdaChaccD/deleteVboxBdaChaccIndexD [delete]
func (bdaChaccDApi *VboxBdaChaccIndexDApi) DeleteVboxBdaChaccIndexD(c *gin.Context) {
	var bdaChaccD vbox.VboxBdaChaccIndexD
	err := c.ShouldBindJSON(&bdaChaccD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    bdaChaccD.DeletedBy = utils.GetUserID(c)
	if err := bdaChaccDService.DeleteVboxBdaChaccIndexD(bdaChaccD); err != nil {
        global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteVboxBdaChaccIndexDByIds 批量删除VboxBdaChaccIndexD
// @Tags VboxBdaChaccIndexD
// @Summary 批量删除VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /bdaChaccD/deleteVboxBdaChaccIndexDByIds [delete]
func (bdaChaccDApi *VboxBdaChaccIndexDApi) DeleteVboxBdaChaccIndexDByIds(c *gin.Context) {
	var IDS request.IdsReq
    err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    deletedBy := utils.GetUserID(c)
	if err := bdaChaccDService.DeleteVboxBdaChaccIndexDByIds(IDS,deletedBy); err != nil {
        global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateVboxBdaChaccIndexD 更新VboxBdaChaccIndexD
// @Tags VboxBdaChaccIndexD
// @Summary 更新VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body vbox.VboxBdaChaccIndexD true "更新VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /bdaChaccD/updateVboxBdaChaccIndexD [put]
func (bdaChaccDApi *VboxBdaChaccIndexDApi) UpdateVboxBdaChaccIndexD(c *gin.Context) {
	var bdaChaccD vbox.VboxBdaChaccIndexD
	err := c.ShouldBindJSON(&bdaChaccD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
    bdaChaccD.UpdatedBy = utils.GetUserID(c)
      verify := utils.Rules{
          "Dt":{utils.NotEmpty()},
      }
    if err := utils.Verify(bdaChaccD, verify); err != nil {
      	response.FailWithMessage(err.Error(), c)
      	return
     }
	if err := bdaChaccDService.UpdateVboxBdaChaccIndexD(bdaChaccD); err != nil {
        global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindVboxBdaChaccIndexD 用id查询VboxBdaChaccIndexD
// @Tags VboxBdaChaccIndexD
// @Summary 用id查询VboxBdaChaccIndexD
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vbox.VboxBdaChaccIndexD true "用id查询VboxBdaChaccIndexD"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /bdaChaccD/findVboxBdaChaccIndexD [get]
func (bdaChaccDApi *VboxBdaChaccIndexDApi) FindVboxBdaChaccIndexD(c *gin.Context) {
	var bdaChaccD vbox.VboxBdaChaccIndexD
	err := c.ShouldBindQuery(&bdaChaccD)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if rebdaChaccD, err := bdaChaccDService.GetVboxBdaChaccIndexD(bdaChaccD.ID); err != nil {
        global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"rebdaChaccD": rebdaChaccD}, c)
	}
}

// GetVboxBdaChaccIndexDList 分页获取VboxBdaChaccIndexD列表
// @Tags VboxBdaChaccIndexD
// @Summary 分页获取VboxBdaChaccIndexD列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query vboxReq.VboxBdaChaccIndexDSearch true "分页获取VboxBdaChaccIndexD列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /bdaChaccD/getVboxBdaChaccIndexDList [get]
func (bdaChaccDApi *VboxBdaChaccIndexDApi) GetVboxBdaChaccIndexDList(c *gin.Context) {
	var pageInfo vboxReq.VboxBdaChaccIndexDSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := bdaChaccDService.GetVboxBdaChaccIndexDInfoList(pageInfo); err != nil {
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
