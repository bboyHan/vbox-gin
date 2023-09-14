// 自动生成模板VboxPayOrder
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// VboxPayOrder 结构体
type VboxPayOrder struct {
	global.GVA_MODEL
	OrderId        string     `json:"order_id" form:"order_id" gorm:"column:order_id;comment:订单ID;size:100;"`
	PAccount       string     `json:"p_account" form:"p_account" gorm:"column:p_account;comment:付方ID;size:50;"`
	Money          int        `json:"money" form:"money" gorm:"column:money;comment:金额;"`
	Uid            int        `json:"uid" form:"uid" gorm:"column:uid;comment:用户ID;"`
	AcId           string     `json:"ac_id" form:"ac_id" gorm:"column:ac_id;comment:账号ID;size:50;"`
	ChannelCode    string     `json:"channel_code" form:"channel_code" gorm:"column:channel_code;comment:通道编码;size:50;"`
	PlatformOid    string     `json:"platform_oid" form:"platform_oid" gorm:"column:platform_oid;comment:平台id;size:500;"`
	PayDevice      string     `json:"pay_device" form:"pay_device" gorm:"column:pay_device;comment:客户端设备;size:200;"`
	PayIp          string     `json:"pay_ip" form:"pay_ip" gorm:"column:pay_ip;comment:客户ip;size:50;"`
	PayRegion      string     `json:"pay_region" form:"pay_region" gorm:"column:pay_region;comment:区域;size:50;"`
	ResourceUrl    string     `json:"resource_url" form:"resource_url" gorm:"column:resource_url;comment:支付链接;type:text;"`
	NotifyUrl      string     `json:"notify_url" form:"notify_url" gorm:"column:notify_url;comment:回调地址;size:200;"`
	OrderStatus    int        `json:"order_status" form:"order_status" gorm:"column:order_status;comment:订单状态;size:2;"`
	CallbackStatus int        `json:"callback_status" form:"callback_status" gorm:"column:callback_status;comment:回调状态;size:2;"`
	CodeUseStatus  int        `json:"code_use_status" form:"code_use_status" gorm:"column:code_use_status;comment:取码状态;size:2;"`
	AsyncTime      *time.Time `json:"async_time" form:"async_time" gorm:"column:async_time;comment:异步执行时间;"`
	CallTime       *time.Time `json:"call_time" form:"call_time" gorm:"column:call_time;comment:回调时间;"`
}

// TableName VboxPayOrder 表名
func (VboxPayOrder) TableName() string {
	return "vbox_pay_order"
}
