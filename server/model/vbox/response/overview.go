package response

type DataOverView struct {
	X interface{} `json:"x" form:"x" url:"x"`
	Y interface{} `json:"y" form:"y" url:"y"`
}

type DataRateOverView struct {
	X0 int `json:"x0" form:"x0" url:"x0"`
	X1 int `json:"x1" form:"x1" url:"x1"`
	X2 int `json:"x2" form:"x2" url:"x2"`
	X3 int `json:"x3" form:"x3" url:"x3"`
	X4 int `json:"x4" form:"x4" url:"x4"`
}

type DataSOverView struct {
	X string `json:"x" form:"x" url:"x"`
	Y string `json:"y" form:"y" url:"y"`
}

type DataSExtOverView struct {
	X0 string `json:"x0" form:"x0" url:"x0"`
	X1 string `json:"x1" form:"x1" url:"x1"`
	X2 string `json:"x2" form:"x2" url:"x2"`
	X3 string `json:"x3" form:"x3" url:"x3"`
	X4 string `json:"x4" form:"x4" url:"x4"`
}

// DataWalletOverView
// 1、我的团队页，指标按顺序为：前日积分收入、前日积分支出、昨日积分收入、昨日积分支出、今日积分收入、今日积分支出、总积分收入、总积分支出、总余额（收入+支出）
// 2、查询消费区分，指标按顺序为：3日前积分消费、2日前积分消费、1日前积分消费、今日积分消费、昨日积分收入、昨日积分支出、今日积分收入、今日积分支出、总余额（收入+支出）
type DataWalletOverView struct {
	X0 int `json:"x0" form:"x0" url:"x0"`
	X1 int `json:"x1" form:"x1" url:"x1"`
	X2 int `json:"x2" form:"x2" url:"x2"`
	X3 int `json:"x3" form:"x3" url:"x3"`
	X4 int `json:"x4" form:"x4" url:"x4"`
	X5 int `json:"x5" form:"x5" url:"x5"`
	X6 int `json:"x6" form:"x6" url:"x6"`
	X7 int `json:"x7" form:"x7" url:"x7"`
	X8 int `json:"x8" form:"x8" url:"x8"`
	X9 int `json:"x9" form:"x9" url:"x9"`
}
