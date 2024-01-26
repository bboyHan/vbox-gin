package response

type ChannelAccountUnused struct {
	Cid   string     `json:"cid" form:"cid"`     // 通道
	Total uint       `json:"total" form:"total"` // 数量
	List  []AccQueue `json:"list" form:"list" gorm:"-"`
}

type AccQueue struct {
	Money  string `json:"money" form:"money"`
	Unused int64  `json:"unused" form:"unused"`
}
