package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	organization "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/model"
)

type OrganizationSearch struct {
	organization.Organization
	request.PageInfo
}

type OrgUserSearch struct {
	organization.OrgUser
	Username string `json:"username" form:"username"`
	OrgIds   []uint `json:"orgIds" form:"orgIds"`
	request.PageInfo
}
