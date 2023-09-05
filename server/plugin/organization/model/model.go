package model

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// Organization 结构体
// 如果含有time.Time 请自行import time包
type Organization struct {
	global.GVA_MODEL
	Name     string `json:"name" form:"name" gorm:"column:name;comment:;"`
	ParentID uint   `json:"parentID" form:"parentID" gorm:"column:parent_id;comment:父节点ID;"`
}

// TableName Organization 表名
func (Organization) TableName() string {
	return "organization"
}

type OrgUser struct {
	Organization   Organization   `json:"organization"`
	OrganizationID uint           `json:"organizationID,omitempty" form:"organizationID" `
	SysUserID      uint           `json:"sysUserID,omitempty" form:"sysUserID"`
	IsAdmin        bool           `json:"isAdmin" form:"isAdmin"`
	SysUser        system.SysUser `json:"sysUser"`
}

type DataAuthority struct {
	AuthorityID   uint                `json:"authorityID" gorm:"column:authority_id;comment:角色ID;"`
	Authority     system.SysAuthority `json:"authority"`
	AuthorityType int                 `json:"authorityType" gorm:"column:authority_type;comment:角色权限标记;"`
}

type OrgUserReq struct {
	OrganizationID   uint   `json:"organizationID,omitempty"`
	ToOrganizationID uint   `json:"toOrganizationID,omitempty"`
	SysUserIDS       []uint `json:"sysUserIDS,omitempty"`
}

// TableName Organization 表名
func (OrgUser) TableName() string {
	return "org_user"
}
