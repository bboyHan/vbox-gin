package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"time"
)

type ChannelAccountSearch struct {
	vbox.ChannelAccount
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	Status         *int       `json:"status" form:"status"`
	SysStatus      *int       `json:"sysStatus" form:"sysStatus"`
	request.PageInfo
}

type ChannelAccountUpd struct {
	Ids    []uint `json:"ids" form:"ids" url:"ids"`
	ID     uint   `json:"id" form:"id" url:"id"`
	Status uint   `json:"status" form:"status" url:"status"`
}
