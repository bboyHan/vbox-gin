package response

type ChannelPayCodeStatistics struct {
	Order    uint    `json:"order" form:"order"`       // 排序
	Location string  `json:"location" form:"location"` // 地点
	CodeNums uint    `json:"codeNums" form:"codeNums"` // 产码数
	Ratio    float64 `json:"ratio" form:"ratio"`       // 占比
}
