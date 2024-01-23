package product

type J3BalanceRecord struct {
	Data    J3BalanceData `json:"data" form:"data" url:"url"`
	Message string        `json:"message" form:"message" url:"message"`
	Code    int           `json:"code" form:"code" url:"code"`
}

type J3BalanceData struct {
	GameName   string `json:"gameName" form:"gameName" url:"gameName"`
	Account    string `json:"account" form:"account" url:"account"`
	ZoneName   string `json:"zoneName" form:"zoneName" url:"zoneName"`
	LeftCoins  int    `json:"leftCoins" form:"leftCoins" url:"leftCoins"`
	LeftSecond string `json:"leftSecond" form:"leftSecond" url:"leftSecond"`
	LeftDate   string `json:"leftDate" form:"leftDate" url:"leftDate"`
	Integral   int    `json:"integral" form:"integral" url:"integral"`
}

type J3Records struct {
	J3BalanceData J3BalanceData     `json:"info" form:"info" url:"info"`
	List          []J3AccountRecord `json:"list" form:"list" url:"list"`
}

type J3AccountRecord struct {
	OrderID    string `json:"orderId" form:"orderId" url:"orderId"`          // 订单号
	AcAccount  string `json:"acAccount" form:"acAccount" url:"acAccount"`    // 账户名
	Money      int    `json:"money" form:"money" url:"money"`                // 金额
	NowTime    string `json:"nowTime" form:"nowTime" url:"nowTime"`          // 当前时间
	HisBalance string `json:"hisBalance" form:"hisBalance" url:"hisBalance"` // 历史余额
	CheckTime  string `json:"checkTime" form:"checkTime" url:"checkTime"`    // 核对时间
	NowBalance string `json:"nowBalance" form:"nowBalance" url:"nowBalance"` // 当前余额
}
