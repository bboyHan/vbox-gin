// 自动生成模板VboxProxy
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 信道 结构体  Proxy
type Proxy struct {
	global.GVA_MODEL
	Chan   string `json:"chan" form:"chan" gorm:"column:chan;comment:渠道;size:32;"`        //渠道
	Type   *int   `json:"type" form:"type" gorm:"column:type;comment:类型;size:4;"`         //类型
	Status *int   `json:"status" form:"status" gorm:"column:status;comment:状态开关;size:4;"` //状态开关
	Url    string `json:"url" form:"url" gorm:"column:url;comment:访问地址;type:text;"`       //访问地址
	Remark string `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:64;"`  //备注
}

// TableName 信道 VboxProxy自定义表名 vbox_proxy
func (Proxy) TableName() string {
	return "vbox_proxy"
}
