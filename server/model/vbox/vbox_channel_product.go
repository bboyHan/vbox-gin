// 自动生成模板VboxChannelProduct
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ChannelProduct 结构体
type ChannelProduct struct {
	global.GVA_MODEL
	ParentId    *uint            `json:"parentId" form:"parentId" gorm:"column:parent_id;comment:产品父ID;size:20;"`
	ChannelCode string           `json:"channelCode" form:"channelCode" gorm:"column:channel_code;comment:通道编码ID;size:20;"`
	ProductName string           `json:"productName" form:"productName" gorm:"column:product_name;comment:产品名称;size:128;"`
	ProductId   string           `json:"productId" form:"productId" gorm:"column:product_id;comment:产品ID;size:128;"`
	Ext         string           `json:"ext" form:"ext" gorm:"column:ext;comment:附加参数;size:256;"`
	Type        uint             `json:"type" form:"type" gorm:"column:type;comment:通道方式，0-原生、1-引导、2-预产;size:10;"`
	PayType     string           `json:"payType" form:"payType" gorm:"column:pay_type;comment:支付方式，weixin、alipay等;size:20;"`
	Children    []ChannelProduct `json:"children" gorm:"-"`
}

type UserChannelProductRate struct {
	ID          uint                     `gorm:"ID"` // 主键ID
	ParentId    *uint                    `json:"parentId"`
	ChannelCode string                   `json:"channelCode"`
	ProductName string                   `json:"productName"`
	ProductId   string                   `json:"productId"`
	Ext         string                   `json:"ext"`
	Type        uint                     `json:"type"`
	PayType     string                   `json:"payType"`
	Children    []UserChannelProductRate `json:"children" gorm:"-"`
	Rate        float64                  `json:"rate"`
}

// TableName VboxChannelProduct 表名
func (ChannelProduct) TableName() string {
	return "vbox_channel_product"
}
