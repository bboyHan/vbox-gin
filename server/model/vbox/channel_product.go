package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ChannelProduct 通道产品 结构体
type ChannelProduct struct {
	global.GVA_MODEL
	ParentId    string           `json:"parentId" form:"parentId" gorm:"column:parent_id;comment:产品父ID;size:10;"`            //产品父ID
	ChannelCode string           `json:"channelCode" form:"channelCode" gorm:"column:channel_code;comment:通道编码ID;size:10;"`  //通道编码ID
	ProductId   string           `json:"productId" form:"productId" gorm:"column:product_id;comment:产品ID;size:128;"`         //产品ID
	ProductName string           `json:"productName" form:"productName" gorm:"column:product_name;comment:产品名称;size:128;"`   //产品名称
	Ext         string           `json:"ext" form:"ext" gorm:"column:ext;comment:附加参数;size:256;"`                            //附加参数
	Type        int              `json:"type" form:"type" gorm:"column:type;comment:通道方式，0-原生,1-引导,2-预产,3-卡密,4-商品;size:4;"`  //通道方式，0-原生、1-引导、2-预产
	PayType     string           `json:"payType" form:"payType" gorm:"column:pay_type;comment:支付方式，weixin、alipay等;size:20;"` //支付方式，weixin、alipay等
	Children    []ChannelProduct `json:"children" gorm:"-"`
	CreatedBy   uint             `gorm:"column:created_by;comment:创建者"`
	UpdatedBy   uint             `gorm:"column:updated_by;comment:更新者"`
	DeletedBy   uint             `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 通道产品 ChannelProduct自定义表名 vbox_channel_product
func (ChannelProduct) TableName() string {
	return "vbox_channel_product"
}
