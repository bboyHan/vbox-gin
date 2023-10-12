// 自动生成模板VboxChannelRate
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VboxChannelRate 结构体
type VboxChannelRate struct {
	global.GVA_MODEL
	Uid         uint   `json:"uid" form:"uid" gorm:"column:uid;comment:uid;size:192;"`
	ChannelCode string `json:"channelCode" form:"channelCode" gorm:"column:channel_code;comment:通道编码ID;size:192;"`
	ProductName string `json:"productName" form:"productName" gorm:"column:product_name;comment:产品名称;size:192;"`
	ProductId   string `json:"productId" form:"productId" gorm:"column:product_id;comment:产品ID;size:192;"`
	//Rate        float64 `json:"rate" form:"rate" gorm:"column:rate;comment:费率;"`
	UnitPrice int  `json:"unitPrice" form:"unit_price" gorm:"column:unit_price;comment:单价积分;"`
	CreatedBy uint `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint `gorm:"column:deleted_by;comment:删除者"`
}

// TableName VboxChannelRate 表名
func (VboxChannelRate) TableName() string {
	return "vbox_channel_rate"
}
