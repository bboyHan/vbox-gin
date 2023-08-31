package vbox

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ChannelShopApi
	ChannelApi
	ChannelGuideImgApi
	ProxyApi
	PayAccountApi
}

var (
	chShopService     = service.ServiceGroupApp.VboxServiceGroup.ChannelShopService
	chService         = service.ServiceGroupApp.VboxServiceGroup.ChannelService
	chGuideImgService = service.ServiceGroupApp.VboxServiceGroup.ChannelGuideImgService
	vpaService        = service.ServiceGroupApp.VboxServiceGroup.PayAccountService
	vboxProxyService  = service.ServiceGroupApp.VboxServiceGroup.ProxyService
)
