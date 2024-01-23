package response

type ChaAccUserCardResp struct {
	Uid        int    `json:"uid" form:"uid"`               //
	AcidCnt    uint   `json:"acidCnt" form:"acidCnt"`       //
	ChannelCnt uint   `json:"channelCnt" form:"channelCnt"` //
	OkOrderCnt uint   `json:"okOrderCnt" form:"okOrderCnt"` //
	OkIncome   uint   `json:"okIncome" form:"okIncome"`     //
	Dt         string `json:"dt" form:"dt"`                 //
}
