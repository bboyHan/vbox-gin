// 自动生成模板Student
package student

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Student 结构体
type Student struct {
	global.GVA_MODEL
	Name      string `json:"name" form:"name" gorm:"column:name;comment:名;size:32;"`
	Score     *int   `json:"score" form:"score" gorm:"column:score;comment:;size:2;"`
	CreatedBy uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName Student 表名
func (Student) TableName() string {
	return "student"
}
