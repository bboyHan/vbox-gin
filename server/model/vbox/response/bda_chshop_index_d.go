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

type ChaShopUserCardResp struct {
	Uid        int    `json:"uid" form:"uid"`               //
	ShopCnt    uint   `json:"shopCnt" form:"ShopCnt"`       //
	ChannelCnt uint   `json:"channelCnt" form:"channelCnt"` //
	OkOrderCnt uint   `json:"okOrderCnt" form:"okOrderCnt"` //
	OkIncome   uint   `json:"okIncome" form:"okIncome"`     //
	Dt         string `json:"dt" form:"dt"`                 //
}

type UserDayChShopIncomeLineChart struct {
	Uid      int    `json:"uid" form:"uid"`
	ShopId   string `json:"shopId" form:"shopId"`
	OkIncome uint   `json:"okIncome" form:"okIncome"`
	StepTime string `json:"stepTime" form:"stepTime"`
}
