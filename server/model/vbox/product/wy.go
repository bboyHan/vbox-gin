package product

type WYBalanceData struct {
	Account   string `json:"account" form:"account" url:"account"`
	TyBalance int    `json:"tyBalance" form:"tyBalance" url:"tyBalance"`
	JsBalance int    `json:"jsBalance" form:"jsBalance" url:"jsBalance"`
}

type WYRecords struct {
	WYBalanceData WYBalanceData     `json:"info" form:"info" url:"info"`
	List          []WYAccountRecord `json:"list" form:"list" url:"list"`
}

type WYAccountRecord struct {
	OrderID    string `json:"orderId" form:"orderId" url:"orderId"`          // 订单号
	AcAccount  string `json:"acAccount" form:"acAccount" url:"acAccount"`    // 账户名
	Money      int    `json:"money" form:"money" url:"money"`                // 金额
	NowTime    string `json:"nowTime" form:"nowTime" url:"nowTime"`          // 当前时间
	HisBalance string `json:"hisBalance" form:"hisBalance" url:"hisBalance"` // 历史余额
	CheckTime  string `json:"checkTime" form:"checkTime" url:"checkTime"`    // 核对时间
	NowBalance string `json:"nowBalance" form:"nowBalance" url:"nowBalance"` // 当前余额
}
