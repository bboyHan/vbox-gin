package request

type ChannelGuideImgTask struct {
	ChannelId string `json:"channelId" form:"channelId" gorm:"-"`
}
