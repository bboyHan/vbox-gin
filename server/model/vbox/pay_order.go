// 自动生成模板PayOrder
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// 订单 结构体  PayOrder
type PayOrder struct {
	global.GVA_MODEL
	OrderId        string     `json:"orderId" form:"orderId" gorm:"column:order_id;comment:订单ID;size:128;"`                    //订单ID
	PAccount       string     `json:"pAccount" form:"pAccount" gorm:"column:p_account;comment:付方ID;size:64;"`                  //付方ID
	Money          int        `json:"money" form:"money" gorm:"column:money;comment:金额;size:16;"`                              //金额
	UnitPrice      int        `json:"unitPrice" form:"unitPrice" gorm:"column:unit_price;comment:单价积分;size:16;"`               //单价积分
	Uid            int        `json:"uid" form:"uid" gorm:"column:uid;comment:用户ID;size:16;"`                                  //用户ID
	AcId           string     `json:"acId" form:"acId" gorm:"column:ac_id;comment:账号ID;size:32;"`                              //账号ID
	ChannelCode    string     `json:"channelCode" form:"channelCode" gorm:"column:channel_code;comment:通道编码;size:32;"`         //通道编码
	PlatformOid    string     `json:"platformOid" form:"platformOid" gorm:"column:platform_oid;comment:平台id;size:256;"`        //平台id
	PayIp          string     `json:"payIp" form:"payIp" gorm:"column:pay_ip;comment:客户ip;size:128;"`                          //客户ip
	PayRegion      string     `json:"payRegion" form:"payRegion" gorm:"column:pay_region;comment:区域;size:128;"`                //区域
	PayDevice      string     `json:"payDevice" form:"payDevice" gorm:"column:pay_device;comment:客户端设备;size:200;"`             //客户端设备
	ResourceUrl    string     `json:"resourceUrl" form:"resourceUrl" gorm:"column:resource_url;comment:支付链接;type:text;"`       //支付链接
	NotifyUrl      string     `json:"notifyUrl" form:"notifyUrl" gorm:"column:notify_url;comment:回调地址;size:200;"`              //回调地址
	OrderStatus    int        `json:"orderStatus" form:"orderStatus" gorm:"column:order_status;comment:订单状态;size:2;"`          //订单状态
	CallbackStatus int        `json:"callbackStatus" form:"callbackStatus" gorm:"column:callback_status;comment:回调状态;size:2;"` //回调状态
	CodeUseStatus  int        `json:"codeUseStatus" form:"codeUseStatus" gorm:"column:code_use_status;comment:取码状态;size:2;"`   //取码状态
	AsyncTime      *time.Time `json:"asyncTime" form:"asyncTime" gorm:"column:async_time;comment:异步执行时间;"`                     //异步执行时间
	CallTime       *time.Time `json:"callTime" form:"callTime" gorm:"column:call_time;comment:回调时间;"`                          //回调时间
	CreatedBy      uint       `gorm:"column:created_by;comment:创建者"`
	UpdatedBy      uint       `gorm:"column:updated_by;comment:更新者"`
	DeletedBy      uint       `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 订单 PayOrder自定义表名 vbox_pay_order
func (PayOrder) TableName() string {
	return "vbox_pay_order"
}
