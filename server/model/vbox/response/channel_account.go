package response

type ChannelAccountUnused struct {
	Cid   string `json:"cid" form:"cid"`     // 通道
	Total uint   `json:"total" form:"total"` // 数量
}
