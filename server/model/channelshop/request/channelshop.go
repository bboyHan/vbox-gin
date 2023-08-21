package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/channelshop"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ChannelShopSearch struct {
	channelshop.ChannelShop
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}
