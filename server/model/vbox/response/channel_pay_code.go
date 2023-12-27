package response

type ChannelPayCodeStatistics struct {
	Order    uint    `json:"order" form:"order"`       // 排序
	Location string  `json:"location" form:"location"` // 地点
	CodeNums uint    `json:"codeNums" form:"codeNums"` // 产码数
	Ratio    float64 `json:"ratio" form:"ratio"`       // 占比
}

type ChannelPayCodeStatisticsResult struct {
	Location string `json:"location" form:"location"` // 地点
	CodeNums uint   `json:"codeNums" form:"codeNums"` // 产码数
}

type PayCodeVO struct {
	CID        string `json:"cid" form:"cid"`               // 通道ID
	Location   string `json:"location" form:"location"`     // 地区编码
	Name       string `json:"name" form:"name"`             // 名称
	Money      string `json:"money" form:"money"`           // 金额
	WaitCnt    int64  `json:"waitCnt" form:"waitCnt"`       // 等待产码数
	PendingCnt int64  `json:"pendingCnt" form:"pendingCnt"` // 待产码数
}

type PCMoney struct {
	Money  string       `json:"money" form:"money"`
	OPList []PCOperator `json:"op" form:"op"`
}

type PCLoc struct {
	Location string  `json:"loc" form:"loc"`
	Cnt      []PCCnt `json:"cnt" form:"cnt"`
}

type PCCnt struct {
	WaitCnt    int64 `json:"waitCnt" form:"waitCnt"`       // 等待产码数
	PendingCnt int64 `json:"pendingCnt" form:"pendingCnt"` // 待产码数
}

type PCOperator struct {
	Operator string  `json:"operator" form:"operator"`
	LocList  []PCLoc `json:"locList" form:"locList"`
}
