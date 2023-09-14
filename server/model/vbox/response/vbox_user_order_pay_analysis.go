package response

type VboxUserOrderPayAnalysis struct {
	Uid              uint   `json:"uid" form:"uid" gorm:"column:uid;comment:父角色ID;size:10;"`
	Username         string `json:"username" form:"username" gorm:"-"`
	Balance          int    `json:"balance" form:"balance" gorm:"-"`
	ChIdCnt          int    `json:"chIdCnt" form:"chIdCnt" gorm:"-"`
	OpenChId         int    `json:"openChId" form:"openChId" gorm:"-"`
	YOrderQuantify   int    `json:"yOrderQuantify" form:"yOrderQuantify" gorm:"-"`
	YOkOrderQuantify int    `json:"yOkOrderQuantify" form:"yOkOrderQuantify" gorm:"-"`
	YOkRate          int    `json:"yOkRate" form:"yOkRate" gorm:"-"`
	YInCome          int    `json:"yInCome" form:"yInCome" gorm:"-"`

	TOrderQuantify   int `json:"tOrderQuantify" form:"tOrderQuantify" gorm:"-"`
	TOkOrderQuantify int `json:"tOkOrderQuantify" form:"tOkOrderQuantify" gorm:"-"`
	TOkRate          int `json:"tOkRate" form:"tOkRate" gorm:"-"`
	TInCome          int `json:"tInCome" form:"tInCome" gorm:"-"`
}

type VboxUserOrderPayAnalysisH struct {
	Uid              int    `json:"uid" form:"uid" gorm:"column:uid;comment:父角色ID;size:10;"`
	Username         string `json:"username" form:"username" gorm:"-"`
	Balance          int    `json:"balance" form:"balance" gorm:"-"`
	ChIdCnt          int    `json:"chIdCnt" form:"chIdCnt" gorm:"-"`
	OpenChId         int    `json:"openChId" form:"openChId" gorm:"-"`
	NewOpenChId      int    `json:"newOpenChId" form:"openChId" gorm:"-"`
	YOrderQuantify   int    `json:"yOrderQuantify" form:"yOrderQuantify" gorm:"-"`
	YOkOrderQuantify int    `json:"yOkOrderQuantify" form:"yOkOrderQuantify" gorm:"-"`
	YOkRate          int    `json:"yOkRate" form:"yOkRate" gorm:"-"`
	YInCome          int    `json:"yInCome" form:"yInCome" gorm:"-"`

	TOrderQuantify   int `json:"tOrderQuantify" form:"tOrderQuantify" gorm:"-"`
	TOkOrderQuantify int `json:"tOkOrderQuantify" form:"tOkOrderQuantify" gorm:"-"`
	TOkRate          int `json:"tOkRate" form:"tOkRate" gorm:"-"`
	TInCome          int `json:"tInCome" form:"tInCome" gorm:"-"`
}
type LineChartData struct {
	LegendData []string               `json:"legendData"`
	XData      []string               `json:"xData"`
	Lists      []LineChartDataYSeries `json:"lists"`
}

type LineChartDataYSeries struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Smooth bool   `json:"smooth"`
	Data   []int  `json:"data"`
}

type CustomBarChartData struct {
	XData []string `json:"xData"`
	Lists []int    `json:"lists"`
}
