package vbox

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ChannelShopApi
	ChannelApi
}

var (
	chShopService = service.ServiceGroupApp.VboxServiceGroup.ChannelShopService
	chService     = service.ServiceGroupApp.VboxServiceGroup.ChannelService
)
