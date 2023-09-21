// 自动生成模板VboxTeamsUser
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
	_ "time"
)

// VboxTeamsUser 结构体
type VboxTeamsUser struct {
	global.GVA_MODEL
	Team_id   uint   `json:"teamId" form:"teamId" gorm:"column:team_id;comment:团队ID;size:191;"`
	Team_name string `json:"teamName" form:"teamName" gorm:"column:team_name;comment:团队名;size:191;"`
	Uid       uint   `json:"uid" form:"uid" gorm:"column:uid;comment:用户ID;size:191;"`
	Leader_id string `json:"leaderId" form:"leaderId" gorm:"column:leader_id;comment:团队领导角色名ID;size:191;"`
	CreatedBy uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint   `gorm:"column:deleted_by;comment:删除者"`
}

type VboxTeamsUserRep struct {
	Team_id       uint           `json:"teamId"`
	Team_name     string         `json:"teamName"`
	Uid           uint           `json:"uid"`
	Username      string         `json:"userName" `
	AuthorityName string         `json:"authorityName" `
	Leader_id     string         `json:"leaderId"`
	DeletedAt     gorm.DeletedAt `json:"deletedAt"` // 删除时间
}

type TeamsUserReq struct {
	TeamID     uint   `json:"teamID,omitempty"`
	ToTeamID   uint   `json:"toTeamID,omitempty"`
	SysUserIDS []uint `json:"sysUserIDS,omitempty"`
}

// TableName VboxTeamsUser 表名
func (VboxTeamsUser) TableName() string {
	return "vbox_teams_user"
}
