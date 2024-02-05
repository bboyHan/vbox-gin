package response

type ChOrgTableResp struct {
	OrganizationId   int     `json:"organizationId" form:"organizationId"`     //团队id
	OrganizationName string  `json:"organizationName" form:"organizationName"` //团队名
	ChannelCode      string  `json:"channelCode" form:"channelCode"`
	StepTime         string  `json:"stepTime" form:"stepTime"`
	ProductId        string  `json:"productId" form:"ProductId"`
	ProductName      string  `json:"productName" form:"productName"`
	OrderQuantify    int     `json:"orderQuantify" form:"orderQuantify"`
	OkOrderQuantify  int     `json:"OkOrderQuantify" form:"OkOrderQuantify"`
	Ratio            float64 `json:"ratio" form:"ratio"`
	OkIncome         int     `json:"okIncome" form:"okIncome"`
}

type ChaOrgUserCardResp struct {
	Uid        int    `json:"uid" form:"uid"`               //
	ChannelCnt uint   `json:"channelCnt" form:"channelCnt"` //
	OkOrderCnt uint   `json:"okOrderCnt" form:"okOrderCnt"` //
	OkIncome   uint   `json:"okIncome" form:"okIncome"`     //
	Dt         string `json:"dt" form:"dt"`                 //
}

type UserDayChOrgIncomeLineChart struct {
	OrganizationName string `json:"organizationName" form:"organizationName"` //团队名
	OkIncome         uint   `json:"okIncome" form:"okIncome"`
	StepTime         string `json:"stepTime" form:"stepTime"`
}
