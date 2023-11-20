package vbox

type RouterGroup struct {
	ChannelAccountRouter
	ChannelProductRouter
	PayOrderRouter
	PayAccountRouter
	ProxyRouter
	ChannelShopRouter
	OrderRouter
	OrgProductRouter
	UserWalletRouter
	VboxChannelPayCodeRouter
	BdaChIndexDRouter
	BdaChaccIndexDRouter
	BdaChShopIndexDRouter
}
