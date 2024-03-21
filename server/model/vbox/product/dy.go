package product

type DyWalletInfoResponse struct {
	Data struct {
		Money             int    `json:"money"`
		UserID            int64  `json:"user_id"`
		Diamond           int    `json:"diamond"`
		PercentWord       int    `json:"percent_word"`
		GoldCoins         int    `json:"gold_coins"`
		Coupon            string `json:"coupon"`
		CouponDescription string `json:"coupon_description"`
		CouponDetail      []int  `json:"coupon_detail"`
		GoldenBeans       int    `json:"golden_beans"`
		CashMoney         int    `json:"cash_money"`
		RemunerationMoney int    `json:"remuneration_money"`
		ContingentMoney   int    `json:"contingent_money"`
		HasLivePermission bool   `json:"has_live_permission"`
		Message           string `json:"message"`
	} `json:"data"`
	Extra struct {
		Now int64 `json:"now"`
	} `json:"extra"`
	StatusCode int `json:"status_code"`
}

type DyWalletInfoRecord struct {
	Money   int   `json:"money"`
	UserID  int64 `json:"user_id"`
	Diamond int   `json:"diamond"`
}

type DyAccountRecord struct {
	OrderID    string `json:"orderId" form:"orderId" url:"orderId"`          // 订单号
	AcAccount  string `json:"acAccount" form:"acAccount" url:"acAccount"`    // 账户名
	Money      int    `json:"money" form:"money" url:"money"`                // 金额
	NowTime    string `json:"nowTime" form:"nowTime" url:"nowTime"`          // 当前时间
	HisBalance string `json:"hisBalance" form:"hisBalance" url:"hisBalance"` // 历史余额
	CheckTime  string `json:"checkTime" form:"checkTime" url:"checkTime"`    // 核对时间
	NowBalance string `json:"nowBalance" form:"nowBalance" url:"nowBalance"` // 当前余额
}
