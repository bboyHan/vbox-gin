package response

type VboxUserOrderPayAnalysis struct {
	Uid              *int   `json:"uid" form:"uid" gorm:"column:uid;comment:父角色ID;size:10;"`
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
