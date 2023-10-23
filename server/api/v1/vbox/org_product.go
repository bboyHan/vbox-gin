package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrgProductApi struct {
}

var opService = service.ServiceGroupApp.VboxServiceGroup.OrgProductService

func (opApi *OrgProductApi) CreateOrgProduct(c *gin.Context) {
	var orgProduct vbox.OrgProductReq
	c.ShouldBindJSON(&orgProduct)
	if err := opService.CreateOrgUser(orgProduct); err != nil {
		global.GVA_LOG.Error("变更失败!", zap.Error(err))
		response.FailWithMessage("变更失败", c)
	} else {
		response.OkWithMessage("变更成功", c)
	}
}

func (opApi *OrgProductApi) DeleteOrgProduct(c *gin.Context) {
	var orgProduct vbox.OrgProductReq
	c.ShouldBindJSON(&orgProduct)
	if err := opService.DeleteOrgProduct(orgProduct.ChannelProductIDS, orgProduct.OrganizationID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

func (opApi *OrgProductApi) FindOrgProductAll(c *gin.Context) {
	org := c.Query("organizationID")
	if ProductIds, err := opService.FindOrgProductAll(org); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(ProductIds, c)
	}
}

func (opApi *OrgProductApi) FindOrgProductList(c *gin.Context) {
	var pageInfo vboxReq.OrgProductSearch
	c.ShouldBindQuery(&pageInfo)
	if list, total, err := opService.GetOrgProductList(pageInfo); err != nil {
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
