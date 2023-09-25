package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"time"
)

type VboxChannelRateSearch struct {
	vbox.VboxChannelRate
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

type VboxChannelRateReq struct {
	TeamId int  `json:"teamId" form:"teamId"`
	Uid    uint `json:"uid" form:"uid"`
	request.PageInfo
}

//type RateReq struct {
//	TeamId        int        `json:"teamId"`
//	TeamName      string     `json:"teamName"`
//	Uid           uint       `json:"uid"`
//	UserName      string     `json:"userName"`
//	AuthorityName string     `json:"authorityName"`
//	LeaderId      string     `json:"leaderId"`  // 这里为了兼容空字符串和 null 值，使用了一个接口类型
//	DeletedAt     *time.Time `json:"deletedAt"` // 同上
//}
