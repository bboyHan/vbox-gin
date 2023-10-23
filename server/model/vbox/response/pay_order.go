package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"time"
)

// PayOrderRes 结构体
type PayOrderRes struct {
	//vbox.ChannelAccount `json:"ca" gorm:"embedded"`
	//vbox.VboxPayOrder   `json:"pa" gorm:"embedded"`
	vbox.ChannelAccount `gorm:"embedded"`
	vbox.PayOrder       `gorm:"embedded"`
	UserId              string    `json:"user_id" gorm:"embedded"`
	CreatedTime         time.Time `json:"created_time" gorm:"embedded"`
}

// Order2PayAccountRes 结构体
type Order2PayAccountRes struct {
	OrderId   string `json:"order_id" form:"order_id" url:"order_id"`
	PayUrl    string `json:"pay_url" form:"pay_url" url:"pay_url"`
	NotifyUrl string `json:"notify_url" form:"notify_url" url:"notify_url"`
	Money     int    `json:"money" form:"money" url:"money"`
	Status    int    `json:"status" form:"status" url:"status"`
}

// OrderSimpleRes 结构体
type OrderSimpleRes struct {
	OrderId     string `json:"order_id" form:"order_id" url:"order_id"`
	Account     string `json:"account" form:"account" url:"account"`
	ResourceUrl string `json:"resource_url" form:"resource_url" url:"resource_url"`
	ChannelCode string `json:"channel_code" form:"channel_code" url:"channel_code"`
	Money       int    `json:"money" form:"money" url:"money"`
	Status      int    `json:"status" form:"status" url:"status"`
}
