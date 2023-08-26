package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ChannelShopBatch struct {
	Uid            *int                  `json:"uid" form:"uid" gorm:"column:uid;comment:用户id;size:11;"`
	Cid            string                `json:"cid" form:"cid" gorm:"column:cid;comment:通道ID;size:50;"`
	Channel        string                `json:"channel" form:"channel" gorm:"column:channel;comment:通道;size:50;"`
	Shop_remark    string                `json:"shop_remark" form:"shop_remark" gorm:"column:shop_remark;comment:店铺备注;size:50;"`
	Address        string                `json:"address" form:"address" gorm:"column:address;comment:店地址;size:500;"`
	Money          *int                  `json:"money" form:"money" gorm:"column:money;comment:金额;size:11;"`
	ShopMarkList   []ChannelShopBatchSub `json:"shopMarkList" gorm:"-"`
	StartCreatedAt *time.Time            `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time            `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}
