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

type ChaOrgRealCardResp struct {
	Title           string  `json:"title" form:"title"` //
	OrderQuantify   int     `json:"orderQuantify" form:"orderQuantify"`
	OkOrderQuantify int     `json:"okOrderQuantify" form:"okOrderQuantify"`
	Ratio           float64 `json:"ratio" form:"ratio"`
	Income          int     `json:"income" form:"income"` //
	Dt              string  `json:"dt" form:"dt"`         //
}

type OrgRealStatisicsResp struct {
	OrganizationId   int    `json:"organizationId" form:"organizationId"`     //
	OrganizationName string `json:"organizationName" form:"organizationName"` //
	StepTime         string `json:"stepTime" form:"stepTime"`
	OrderQuantify    int    `json:"orderQuantify" form:"orderQuantify"`
	OkOrderQuantify  int    `json:"okOrderQuantify" form:"okOrderQuantify"`
	OkIncome         int    `json:"okIncome" form:"okIncome"` //
}

type CidRealStatisicsResp struct {
	ChannelCode     string `json:"channelCode" form:"channelCode"`
	StepTime        string `json:"stepTime" form:"stepTime"`
	ProductId       string `json:"productId" form:"ProductId"`
	ProductName     string `json:"productName" form:"productName"`
	OrderQuantify   int    `json:"orderQuantify" form:"orderQuantify"`
	OkOrderQuantify int    `json:"okOrderQuantify" form:"okOrderQuantify"`
	OkIncome        int    `json:"okIncome" form:"okIncome"` //
}

type UidRealStatisicsResp struct {
	Uid             string `json:"uid" form:"uid"`
	StepTime        string `json:"stepTime" form:"stepTime"`
	UserName        string `json:"userName" form:"userName"`
	OrderQuantify   int    `json:"orderQuantify" form:"orderQuantify"`
	OkOrderQuantify int    `json:"okOrderQuantify" form:"okOrderQuantify"`
	OkIncome        int    `json:"okIncome" form:"okIncome"` //
}

type PaccRealStatisicsResp struct {
	PAccount        string `json:"pAccount" form:"pAccount"`
	StepTime        string `json:"stepTime" form:"stepTime"`
	UserName        string `json:"userName" form:"userName"`
	OrderQuantify   int    `json:"orderQuantify" form:"orderQuantify"`
	OkOrderQuantify int    `json:"okOrderQuantify" form:"okOrderQuantify"`
	OkIncome        int    `json:"okIncome" form:"okIncome"` //
}
