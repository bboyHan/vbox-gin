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
	UserName string `json:"userName" form:"userName"`
	request.PageInfo
}
