package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"time"
)

type ChannelShopSearch struct {
	vbox.ChannelShop
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

type ChannelShop struct {
	Cid             string           `json:"cid" form:"cid"`               //通道ID
	ProductId       string           `json:"productId" form:"productId"`   //产品ID
	ShopRemark      string           `json:"shopRemark" form:"shopRemark"` //店铺备注
	ChannelShopList []ChannelShopSub `json:"list" form:"list"`
	CreatedBy       uint             `json:"createdAt" form:"createdAt"`
}

type ChannelShopSub struct {
	ID      uint   `json:"id" form:"id"`
	Address string `json:"address" form:"address"` //店地址
	Money   int    `json:"money" form:"money"`     //金额
	Status  int    `json:"status" form:"status"`   //开关
}
