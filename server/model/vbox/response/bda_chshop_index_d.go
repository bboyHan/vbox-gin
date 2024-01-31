package response

type ChShopTableResp struct {
	UID             int     `json:"uid" form:"uid"`
	ChannelCode     string  `json:"channelCode" form:"channelCode"`
	ShopId          string  `json:"shopId" form:"ShopId"`
	StepTime        string  `json:"stepTime" form:"stepTime"`
	ShopRemark      string  `json:"shopRemark" form:"shopRemark"`
	ProductID       string  `json:"productID" form:"ProductID"`
	ProductName     string  `json:"productName" form:"productID"`
	Username        string  `json:"username" form:"username"`
	OrderQuantify   int     `json:"orderQuantify" form:"orderQuantify"`
	OkOrderQuantify int     `json:"OkOrderQuantify" form:"OkOrderQuantify"`
	Ratio           float64 `json:"ratio" form:"ratio"`
	OkIncome        int     `json:"okIncome" form:"okIncome"`
}
