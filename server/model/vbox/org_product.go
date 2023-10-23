package vbox

import "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/model"

type OrgProduct struct {
	Organization     model.Organization `json:"organization"`
	OrganizationID   uint               `json:"organizationID,omitempty" form:"organizationID" `
	ChannelProductID uint               `json:"channelProductID,omitempty" form:"channelProductID"`
	ChannelProduct   ChannelProduct     `json:"channelProduct"`
}

type OrgProductRes struct {
	Name        string `json:"name" form:"name"`
	Oid         string `json:"oid" form:"oid"`
	ProductName string `json:"productName" form:"productName"`
	ChannelCode string `json:"channelCode" form:"channelCode"`
	CpId        int    `json:"cpId" form:"cpId"`
}

type OrgProductReq struct {
	OrganizationID    uint   `json:"organizationID,omitempty"`
	ToOrganizationID  uint   `json:"toOrganizationID,omitempty"`
	ChannelProductIDS []uint `json:"channelProductIDS,omitempty"`
}

// TableName Organization 表名
func (OrgProduct) TableName() string {
	return "org_product"
}
