package student

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type StudentRouter struct {
}

// InitStudentRouter 初始化 Student 路由信息
func (s *StudentRouter) InitStudentRouter(Router *gin.RouterGroup) {
	stdRouter := Router.Group("std").Use(middleware.OperationRecord())
	stdRouterWithoutRecord := Router.Group("std")
	var stdApi = v1.ApiGroupApp.StudentApiGroup.StudentApi
	{
		stdRouter.POST("createStudent", stdApi.CreateStudent)             // 新建Student
		stdRouter.DELETE("deleteStudent", stdApi.DeleteStudent)           // 删除Student
		stdRouter.DELETE("deleteStudentByIds", stdApi.DeleteStudentByIds) // 批量删除Student
		stdRouter.PUT("updateStudent", stdApi.UpdateStudent)              // 更新Student
	}
	{
		stdRouterWithoutRecord.GET("findStudent", stdApi.FindStudent)       // 根据ID获取Student
		stdRouterWithoutRecord.GET("getStudentList", stdApi.GetStudentList) // 获取Student列表
	}
}
