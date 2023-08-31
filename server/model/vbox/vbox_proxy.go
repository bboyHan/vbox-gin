// 自动生成模板VboxProxy
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Proxy 结构体
type Proxy struct {
	global.GVA_MODEL
	Chan   string `json:"chan" form:"chan" gorm:"column:chan;comment:;size:50;"`
	Type   *int   `json:"type" form:"type" gorm:"column:type;comment:;size:10;"`
	Url    string `json:"url" form:"url" gorm:"column:url;comment:;size:200;"`
	Remark string `json:"remark" form:"remark" gorm:"column:remark;comment:;size:50;"`
}

// TableName VboxProxy 表名
func (Proxy) TableName() string {
	return "vbox_proxy"
}
