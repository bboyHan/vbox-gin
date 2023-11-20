package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BdaChShopIndexDRouter struct {
}

// InitBdaChShopIndexDRouter 初始化 用户通道店铺成率统计-天更新 路由信息
func (s *BdaChShopIndexDRouter) InitBdaChShopIndexDRouter(Router *gin.RouterGroup) {
	bdaChshopIndexDRouter := Router.Group("bdaChshopIndexD").Use(middleware.OperationRecord())
	bdaChshopIndexDRouterWithoutRecord := Router.Group("bdaChshopIndexD")
	var bdaChshopIndexDApi = v1.ApiGroupApp.VboxApiGroup.BdaChShopIndexDApi
	{
		bdaChshopIndexDRouter.POST("createBdaChShopIndexD", bdaChshopIndexDApi.CreateBdaChShopIndexD)             // 新建用户通道店铺成率统计-天更新
		bdaChshopIndexDRouter.DELETE("deleteBdaChShopIndexD", bdaChshopIndexDApi.DeleteBdaChShopIndexD)           // 删除用户通道店铺成率统计-天更新
		bdaChshopIndexDRouter.DELETE("deleteBdaChShopIndexDByIds", bdaChshopIndexDApi.DeleteBdaChShopIndexDByIds) // 批量删除用户通道店铺成率统计-天更新
		bdaChshopIndexDRouter.PUT("updateBdaChShopIndexD", bdaChshopIndexDApi.UpdateBdaChShopIndexD)              // 更新用户通道店铺成率统计-天更新
	}
	{
		bdaChshopIndexDRouterWithoutRecord.GET("findBdaChShopIndexD", bdaChshopIndexDApi.FindBdaChShopIndexD)       // 根据ID获取用户通道店铺成率统计-天更新
		bdaChshopIndexDRouterWithoutRecord.GET("getBdaChShopIndexDList", bdaChshopIndexDApi.GetBdaChShopIndexDList) // 获取用户通道店铺成率统计-天更新列表
	}
}
