package vbox

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ChannelShopApi
	ChannelApi
	ChannelGuideImgApi
}

var (
	chShopService     = service.ServiceGroupApp.VboxServiceGroup.ChannelShopService
	chService         = service.ServiceGroupApp.VboxServiceGroup.ChannelService
	chGuideImgService = service.ServiceGroupApp.VboxServiceGroup.ChannelGuideImgService
)
