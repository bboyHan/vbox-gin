package student

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/student"
	studentReq "github.com/flipped-aurora/gin-vue-admin/server/model/student/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type StudentApi struct {
}

var stdService = service.ServiceGroupApp.StudentServiceGroup.StudentService

// CreateStudent 创建Student
// @Tags Student
// @Summary 创建Student
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body student.Student true "创建Student"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /std/createStudent [post]
func (stdApi *StudentApi) CreateStudent(c *gin.Context) {

	//fmt.Println("c: %v", *c)
	var std student.Student
	err := c.ShouldBindJSON(&std)
	fmt.Println(*std.Score)
	fmt.Println(std.Name)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	std.CreatedBy = utils.GetUserID(c)
	if err := stdService.CreateStudent(&std); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// DeleteStudent 删除Student
// @Tags Student
// @Summary 删除Student
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body student.Student true "删除Student"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /std/deleteStudent [delete]
func (stdApi *StudentApi) DeleteStudent(c *gin.Context) {
	var std student.Student
	err := c.ShouldBindJSON(&std)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	std.DeletedBy = utils.GetUserID(c)
	if err := stdService.DeleteStudent(std); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// DeleteStudentByIds 批量删除Student
// @Tags Student
// @Summary 批量删除Student
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "批量删除Student"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"批量删除成功"}"
// @Router /std/deleteStudentByIds [delete]
func (stdApi *StudentApi) DeleteStudentByIds(c *gin.Context) {
	var IDS request.IdsReq
	err := c.ShouldBindJSON(&IDS)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	deletedBy := utils.GetUserID(c)
	if err := stdService.DeleteStudentByIds(IDS, deletedBy); err != nil {
		global.GVA_LOG.Error("批量删除失败!", zap.Error(err))
		response.FailWithMessage("批量删除失败", c)
	} else {
		response.OkWithMessage("批量删除成功", c)
	}
}

// UpdateStudent 更新Student
// @Tags Student
// @Summary 更新Student
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body student.Student true "更新Student"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /std/updateStudent [put]
func (stdApi *StudentApi) UpdateStudent(c *gin.Context) {
	var std student.Student
	err := c.ShouldBindJSON(&std)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	std.UpdatedBy = utils.GetUserID(c)
	if err := stdService.UpdateStudent(std); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// FindStudent 用id查询Student
// @Tags Student
// @Summary 用id查询Student
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query student.Student true "用id查询Student"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"查询成功"}"
// @Router /std/findStudent [get]
func (stdApi *StudentApi) FindStudent(c *gin.Context) {
	var std student.Student
	err := c.ShouldBindQuery(&std)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if restd, err := stdService.GetStudent(std.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithData(gin.H{"restd": restd}, c)
	}
}

// GetStudentList 分页获取Student列表
// @Tags Student
// @Summary 分页获取Student列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query studentReq.StudentSearch true "分页获取Student列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /std/getStudentList [get]
func (stdApi *StudentApi) GetStudentList(c *gin.Context) {
	var pageInfo studentReq.StudentSearch
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if list, total, err := stdService.GetStudentInfoList(pageInfo); err != nil {
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
