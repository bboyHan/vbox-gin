package vbox

import "github.com/flipped-aurora/gin-vue-admin/server/service"

type ApiGroup struct {
	ChannelShopApi
	ChannelApi
	ChannelGuideImgApi
	ProxyApi
	PayAccountApi
	VboxPayOrderApi
	ChannelProductApi
	ChannelAccountApi
	VboxUserWalletApi
	VboxBdaChIndexDApi
	VboxBdaChaccIndexDApi
	VboxTeamsUserApi
	VboxTeamsApi
	VboxChannelRateApi
}

var (
	chShopService     = service.ServiceGroupApp.VboxServiceGroup.ChannelShopService
	chService         = service.ServiceGroupApp.VboxServiceGroup.ChannelService
	chGuideImgService = service.ServiceGroupApp.VboxServiceGroup.ChannelGuideImgService
	vpoService        = service.ServiceGroupApp.VboxServiceGroup.VboxPayOrderService
	vpaService        = service.ServiceGroupApp.VboxServiceGroup.PayAccountService
	vboxProxyService  = service.ServiceGroupApp.VboxServiceGroup.ProxyService
	vcpService        = service.ServiceGroupApp.VboxServiceGroup.ChannelProductService
	vcaService        = service.ServiceGroupApp.VboxServiceGroup.ChannelAccountService
	vuwService        = service.ServiceGroupApp.VboxServiceGroup.VboxUserWalletService
	bdaChDService     = service.ServiceGroupApp.VboxServiceGroup.VboxBdaChIndexDService
	bdaChaccDService  = service.ServiceGroupApp.VboxServiceGroup.VboxBdaChaccIndexDService
	userService       = service.ServiceGroupApp.SystemServiceGroup.UserService
)
