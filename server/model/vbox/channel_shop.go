// 自动生成模板ChannelShop
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ChannelShop 引导商铺 结构体
type ChannelShop struct {
	global.GVA_MODEL
	Cid        string `json:"cid" form:"cid" gorm:"column:cid;comment:通道ID;size:32;"`                        //通道ID
	ProductId  string `json:"productId" form:"productId" gorm:"column:product_id;comment:通道;size:32;"`       //产品ID
	ShopRemark string `json:"shopRemark" form:"shopRemark" gorm:"column:shop_remark;comment:店铺备注;size:128;"` //店铺备注
	Address    string `json:"address" form:"address" gorm:"column:address;comment:店地址;type:text;"`           //店地址
	MarkId     string `json:"markId" form:"markId" gorm:"column:mark_id;comment:标识;size:128;"`
	Device     string `json:"device" form:"device" gorm:"default:default;column:device;comment:设备;size:32;"`
	Money      int    `json:"money" form:"money" gorm:"column:money;comment:金额;size:32;"`   //金额
	Status     int    `json:"status" form:"status" gorm:"column:status;comment:开关;size:4;"` //开关
	CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 引导商铺 ChannelShop自定义表名 vbox_channel_shop
func (ChannelShop) TableName() string {
	return "vbox_channel_shop"
}
