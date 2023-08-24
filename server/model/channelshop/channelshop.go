// 自动生成模板ChannelShop
package channelshop

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ChannelShop 结构体
type ChannelShop struct {
	global.GVA_MODEL
	Uid          *int          `json:"uid" form:"uid" gorm:"column:uid;comment:用户id;size:11;"`
	Cid          string        `json:"cid" form:"cid" gorm:"column:cid;comment:通道ID;size:50;"`
	Channel      string        `json:"channel" form:"channel" gorm:"column:channel;comment:通道;size:50;"`
	Shop_remark  string        `json:"shop_remark" form:"shop_remark" gorm:"column:shop_remark;comment:店铺备注;size:50;"`
	Address      string        `json:"address" form:"address" gorm:"column:address;comment:店地址;size:500;"`
	Money        *int          `json:"money" form:"money" gorm:"column:money;comment:金额;size:11;"`
	Status       *int          `json:"status" form:"status" gorm:"column:status;comment:开关;size:11;"`
	DisMoney     string        `json:"disMoney" gorm:"-"`
	OpenAndClose string        `json:"openAndClose" gorm:"-"`
	TreeLevel    int           `json:"treeLevel" gorm:"-"`
	Children     []ChannelShop `json:"children" gorm:"-"`
	CreatedBy    uint          `gorm:"column:created_by;comment:创建者"`
	UpdatedBy    uint          `gorm:"column:updated_by;comment:更新者"`
	DeletedBy    uint          `gorm:"column:deleted_by;comment:删除者"`
}

// TableName ChannelShop 表名
func (ChannelShop) TableName() string {
	return "vbox_channel_shop"
}
