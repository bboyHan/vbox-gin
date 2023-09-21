// 自动生成模板VboxTeams
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VboxTeams 结构体
type VboxTeams struct {
	global.GVA_MODEL
	Name      string `json:"name" form:"name" gorm:"column:name;comment:团队名;size:191;"`
	Parent_id int    `json:"parentId" form:"parentId" gorm:"column:parent_id;comment:父节点ID;"`
	CreatedBy uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName VboxTeams 表名
func (VboxTeams) TableName() string {
	return "vbox_teams"
}
