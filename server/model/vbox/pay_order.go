package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// PayOrder 订单 结构体
type PayOrder struct {
	global.GVA_MODEL
	OrderId     string     `json:"orderId" form:"orderId" gorm:"column:order_id;comment:订单ID;size:128;index;"`                                       //订单ID
	Money       int        `json:"money" form:"money" gorm:"column:money;comment:金额;size:16;"`                                                       //金额
	UnitPrice   int        `json:"unitPrice" form:"unitPrice" gorm:"column:unit_price;comment:单价积分;size:16;"`                                        //单价积分
	UnitId      int        `json:"unitId" form:"unitId" gorm:"column:unit_id;comment:积分关联ID;size:16;"`                                               //用户ID
	EventId     string     `json:"eventId" form:"eventId" gorm:"column:event_id;comment:事件关联ID;size:32;"`                                            //账号ID
	EventType   int        `json:"eventType" form:"eventType" gorm:"column:event_type;comment:事件类型（1-商铺关联，2-付码关联）;size:2;"`                          //账号ID
	AcId        string     `json:"acId" form:"acId" gorm:"column:ac_id;comment:账号ID;size:32;index;"`                                                 //账号ID
	AcAccount   string     `json:"acAccount" form:"acAccount" gorm:"column:ac_account;comment:通道账号;size:128;"`                                       //通道账号
	ChannelCode string     `json:"channelCode" form:"channelCode" gorm:"column:channel_code;comment:通道编码;size:32;"`                                  //通道编码
	PlatId      string     `json:"platId" form:"platId" gorm:"column:plat_id;comment:平台id;size:256;"`                                                //平台id
	PayIp       string     `json:"payIp" form:"payIp" gorm:"column:pay_ip;comment:客户ip;size:128;"`                                                   //客户ip
	PayRegion   string     `json:"payRegion" form:"payRegion" gorm:"column:pay_region;comment:区域;size:128;"`                                         //区域
	PayDevice   string     `json:"payDevice" form:"payDevice" gorm:"column:pay_device;comment:客户端设备;size:200;"`                                      //客户端设备
	ResourceUrl string     `json:"resourceUrl" form:"resourceUrl" gorm:"column:resource_url;comment:支付链接;type:text;"`                                //支付链接
	PAccount    string     `json:"pAccount" form:"pAccount" gorm:"column:p_account;comment:付方ID;size:64;"`                                           //付方ID
	NotifyUrl   string     `json:"notifyUrl" form:"notifyUrl" gorm:"column:notify_url;comment:回调地址;size:200;"`                                       //回调地址
	OrderStatus int        `json:"orderStatus" form:"orderStatus" gorm:"default:2;column:order_status;comment:订单状态,1-已支付,2-待支付,3-超时,0-匹配失败;size:2;"` //订单状态
	CbStatus    int        `json:"cbStatus" form:"cbStatus" gorm:"default:2;column:cb_status;comment:回调状态,1-已回调,2-待回调;size:2;"`                      //回调状态
	HandStatus  int        `json:"handStatus" form:"handStatus" gorm:"default:2;column:hand_status;comment:手动补单状态,1-已手动补单,2-默认,3-候补单;size:2;"`       //手动补单
	ExpTime     *time.Time `json:"expTime" form:"expTime" gorm:"column:exp_time;comment:订单过期时间;"`                                                    //异步执行时间
	CbTime      *time.Time `json:"cbTime" form:"cbTime" gorm:"column:cb_time;comment:回调时间;index;"`                                                   //回调时间
	CreatedBy   uint       `gorm:"column:created_by;comment:创建者"`
	UpdatedBy   uint       `gorm:"column:updated_by;comment:更新者"`
	DeletedBy   uint       `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 订单 PayOrder自定义表名 vbox_pay_order
func (PayOrder) TableName() string {
	return "vbox_pay_order"
}

type PayOrderList []PayOrder

func (c PayOrderList) Len() int {
	return len(c)
}

func (c PayOrderList) Less(i, j int) bool {
	at := c[j].CreatedAt
	u := *at
	return c[i].CreatedAt.Before(u)
}

func (c PayOrderList) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
