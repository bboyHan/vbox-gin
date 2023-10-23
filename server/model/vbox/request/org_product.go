package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
)

type OrgProductSearch struct {
	vbox.OrgProduct
	ProductName string `json:"productName" form:"productName"`
	request.PageInfo
}
