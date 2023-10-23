package vbox

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type OrgProductRouter struct{}

func (s *OrgProductRouter) InitOrgProductRouter(Router *gin.RouterGroup) {
	opRouter := Router.Group("org").Use(middleware.OperationRecord())
	opRouterWithoutRecord := Router.Group("org")
	var opApi = v1.ApiGroupApp.VboxApiGroup.OrgProductApi
	{
		opRouter.POST("createOrgProduct", opApi.CreateOrgProduct)                // 产品入组织
		opRouterWithoutRecord.DELETE("deleteOrgProduct", opApi.DeleteOrgProduct) // 删除当前组织下选中产品
	}
	{
		opRouterWithoutRecord.GET("findOrgProductAll", opApi.FindOrgProductAll)   // 获取当前组织下所有可用产品
		opRouterWithoutRecord.GET("findOrgProductList", opApi.FindOrgProductList) // 获取当前组织下所有可用产品(分页)
	}
}
