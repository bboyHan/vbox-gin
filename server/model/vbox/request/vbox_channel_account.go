package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"time"
)

type ChannelAccountSearch struct{
    vbox.ChannelAccount
    StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
    EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
    StartStatus  *int  `json:"startStatus" form:"startStatus"`
    EndStatus  *int  `json:"endStatus" form:"endStatus"`
    StartSysStatus  *int  `json:"startSysStatus" form:"startSysStatus"`
    EndSysStatus  *int  `json:"endSysStatus" form:"endSysStatus"`
    request.PageInfo
}
