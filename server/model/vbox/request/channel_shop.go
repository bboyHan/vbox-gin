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

type ChannelShopReq struct {
	Id         uint   `json:"id" form:"id"`                 //ID
	Type       uint   `json:"type" form:"type"`             //操作类型 1-更新店名 2-开关单条 3-开关整个店
	Cid        string `json:"cid" form:"cid"`               //通道ID
	ProductId  string `json:"productId" form:"productId"`   //产品ID
	Status     int    `json:"status" form:"status"`         //开关
	ShopRemark string `json:"shopRemark" form:"shopRemark"` //店铺备注
	UpdatedBy  uint   `json:"updatedBy" form:"createdBy"`
}

type ChannelShop struct {
	Cid             string           `json:"cid" form:"cid"`               //通道ID
	ProductId       string           `json:"productId" form:"productId"`   //产品ID
	ShopRemark      string           `json:"shopRemark" form:"shopRemark"` //店铺备注
	ChannelShopList []ChannelShopSub `json:"list" form:"list"`
	CreatedBy       uint             `json:"createdBy" form:"createdBy"`
}

type ChannelShopSub struct {
	ID      uint   `json:"id" form:"id"`
	Address string `json:"address" form:"address"` //店地址
	Money   int    `json:"money" form:"money"`     //金额
	Status  int    `json:"status" form:"status"`   //开关
}

type ChannelShopList []ChannelShop

func (c ChannelShopList) Len() int {
	return len(c)
}

func (c ChannelShopList) Less(i, j int) bool {
	return c[i].ProductId < c[j].ProductId
}

func (c ChannelShopList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
