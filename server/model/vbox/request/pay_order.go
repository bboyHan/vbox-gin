package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"time"
)

type PayOrderSearch struct {
	vbox.PayOrder
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

// CreateOrder2PayAccount 结构体
type CreateOrder2PayAccount struct {
	Account     string `json:"account" form:"account" url:"account"`
	Key         string `json:"key" form:"key" url:"key"`
	OrderId     string `json:"order_id" form:"order_id" url:"order_id"`
	Money       int    `json:"money" form:"money" url:"money"`
	ChannelCode string `json:"channel_code" form:"channel_code" url:"channel_code"`
	NotifyUrl   string `json:"notify_url" form:"notify_url" url:"notify_url"`
	Sign        string `json:"sign" form:"sign" url:"sign"`
}

// QueryOrder2PayAccount 结构体
type QueryOrder2PayAccount struct {
	Account string `json:"account" form:"account" url:"account"`
	OrderId string `json:"order_id" form:"order_id" url:"order_id"`
	Key     string `json:"key" form:"key" url:"key"`
	Sign    string `json:"sign" form:"sign" url:"sign"`
}

// QueryOrderSimple 结构体
type QueryOrderSimple struct {
	OrderId   string `json:"order_id" form:"order_id" url:"order_id"`
	UserAgent string `json:"user_agent" form:"user_agent" url:"user_agent"`
	PayIp     string `json:"pay_ip" form:"pay_ip" url:"pay_ip"`
	PayRegion string `json:"pay_region" form:"pay_region" url:"pay_region"`
	PayDevice string `json:"pay_device" form:"pay_device" url:"pay_device"`
}

// CreateOrderTest 结构体
type CreateOrderTest struct {
	Money       int    `json:"money" form:"money" url:"money"`
	ChannelCode string `json:"channel_code" form:"channel_code" url:"channel_code"`
	AuthCaptcha string `json:"auth_captcha" form:"auth_captcha" url:"auth_captcha"`
	Username    string `json:"username" form:"username" url:"username"`
	UserId      uint   `json:"userId" form:"userId" url:"userId"`
	NotifyUrl   string `json:"notify_url" form:"notify_url" url:"notify_url"`
}

type OrdersDtData struct {
	ChannelCode string `json:"channelCode" form:"channelCode"`
	Dt          string `json:"dt" form:"dt"`
	request.PageInfo
}
