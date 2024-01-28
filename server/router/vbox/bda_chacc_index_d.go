package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type BdaChaccIndexDRouter struct {
}

// InitBdaChaccIndexDRouter 初始化 用户通道粒度成率统计-天更新 路由信息
func (s *BdaChaccIndexDRouter) InitBdaChaccIndexDRouter(Router *gin.RouterGroup) {
	bdaChaccIndexDRouter := Router.Group("bdaChaccIndexD").Use(middleware.OperationRecord())
	bdaChaccIndexDRouterWithoutRecord := Router.Group("bdaChaccIndexD")
	var bdaChaccIndexDApi = v1.ApiGroupApp.VboxApiGroup.BdaChaccIndexDApi
	{
		bdaChaccIndexDRouter.POST("createBdaChaccIndexD", bdaChaccIndexDApi.CreateBdaChaccIndexD)             // 新建用户通道粒度成率统计-天更新
		bdaChaccIndexDRouter.DELETE("deleteBdaChaccIndexD", bdaChaccIndexDApi.DeleteBdaChaccIndexD)           // 删除用户通道粒度成率统计-天更新
		bdaChaccIndexDRouter.DELETE("deleteBdaChaccIndexDByIds", bdaChaccIndexDApi.DeleteBdaChaccIndexDByIds) // 批量删除用户通道粒度成率统计-天更新
		bdaChaccIndexDRouter.PUT("updateBdaChaccIndexD", bdaChaccIndexDApi.UpdateBdaChaccIndexD)              // 更新用户通道粒度成率统计-天更新
	}
	{
		bdaChaccIndexDRouterWithoutRecord.GET("findBdaChaccIndexD", bdaChaccIndexDApi.FindBdaChaccIndexD)                       // 根据ID获取用户通道粒度成率统计-天更新
		bdaChaccIndexDRouterWithoutRecord.GET("getBdaChaccIndexDList", bdaChaccIndexDApi.GetBdaChaccIndexDList)                 // 获取用户通道粒度成率统计-天更新列表
		bdaChaccIndexDRouterWithoutRecord.GET("getBdaChaccIndexDUesrOverview", bdaChaccIndexDApi.GetBdaChaccIndexDUesrOverview) // 获取用户通道粒度成率统计-天更新列表
		bdaChaccIndexDRouterWithoutRecord.GET("getBdaChaccIndexToDayIncome", bdaChaccIndexDApi.GetBdaChaccIndexToDayIncome)     // 获取用户通道粒度成率统计-天更新列表
		bdaChaccIndexDRouterWithoutRecord.GET("getBdaChaccIndexToDayInOkCnt", bdaChaccIndexDApi.GetBdaChaccIndexToDayInOkCnt)   // 获取用户通道粒度成率统计-天更新列表
		bdaChaccIndexDRouterWithoutRecord.GET("getBdaChaccIndexToWeekIncome", bdaChaccIndexDApi.GetBdaChaccIndexToWeekIncome)   // 获取用户通道粒度成率统计-天更新列表
		bdaChaccIndexDRouterWithoutRecord.GET("getBdaChaccIndexToWeekInOkCnt", bdaChaccIndexDApi.GetBdaChaccIndexToWeekInOkCnt) // 获取用户通道粒度成率统计-天更新列表
		bdaChaccIndexDRouterWithoutRecord.GET("getBdaChaccIndexDListWeek", bdaChaccIndexDApi.GetBdaChaccIndexDListWeek)         // 获取用户通道粒度成率统计-天更新列表
	}
}
