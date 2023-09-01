package request

type ChannelShopBatchSub struct {
	Address string `json:"address" form:"address" gorm:"column:address;comment:店地址;size:500;"`
	Money   string `json:"money" form:"money" gorm:"column:money;comment:金额;size:11;"`
}
