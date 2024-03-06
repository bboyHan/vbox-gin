package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"time"
)

type BdaChorgIndexDSearch struct {
	vbox.BdaChorgIndexD
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

type OrgSelectForm struct {
	SysUserID      *uint  `json:"sysUserID" form:"sysUserID"`
	OrganizationID int    `json:"organizationID" form:"organizationID"`
	Cid            string `json:"cid" form:"cid"`
	PAccount       string `json:"pAccount" form:"pAccount"`
	Uid            uint   `json:"uid" form:"uid"`
	Dt             string `json:"dt" form:"dt"`
}
