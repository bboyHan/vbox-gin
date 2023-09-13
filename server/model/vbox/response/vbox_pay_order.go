package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/vbox"

// VboxPayOrderRes 结构体
type VboxPayOrderRes struct {
	vbox.ChannelAccount `gorm:"embedded"`
	vbox.VboxPayOrder   `gorm:"embedded"`
	UserId              string `json:"user_id" gorm:"embedded"`
}

//ID             uint       `gorm:"-"`          // 主键ID
//CreatedAt      time.Time  `json:"-" gorm:"-"` // 创建时间
//UpdatedAt      time.Time  `json:"-" gorm:"-"` // 更新时间
//DeletedAt      time.Time  `json:"-" gorm:"-"` // 删除时间
//OrderId        string     `json:"orderId" form:"order_id" gorm:"-"`
//PAccount       string     `json:"p_account" form:"p_account" gorm:"-"`
//Cost           *int       `json:"cost" form:"cost" gorm:"-"`
//Uid            *int       `json:"uid" form:"uid" gorm:"-"`
//AcId           string     `json:"acId" form:"ac_id" gorm:"-"`
//AcAccount      string     `json:"acAccount" form:"ac_account" gorm:"-""`
//AcRemark       string     `json:"ac_remark" form:"ac_remark" gorm:"-"`
//CChannelId     string     `json:"c_channel_id" form:"c_channel_id" gorm:"-"`
//PlatformOid    string     `json:"platform_oid" form:"platform_oid" gorm:"-"`
//PayDevice      string     `json:"pay_device" form:"pay_device" gorm:"-"`
//PayIp          string     `json:"pay_ip" form:"pay_ip" gorm:"-"`
//PayRegion      string     `json:"pay_region" form:"pay_region" gorm:"-"`
//ResourceUrl    string     `json:"resource_url" form:"resource_url" gorm:"-"`
//NotifyUrl      string     `json:"notify_url" form:"notify_url" gorm:"-"`
//OrderStatus    *int       `json:"order_status" form:"order_status" gorm:"-"`
//CallbackStatus *int       `json:"callback_status" form:"callback_status" gorm:"-"`
//CodeUseStatus  *int       `json:"code_use_status" form:"code_use_status" gorm:"-"`
//CreateTime     *time.Time `json:"create_time" form:"create_time" gorm:"-"`
//AsyncTime      *time.Time `json:"async_time" form:"async_time" gorm:"-"`
//CallTime       *time.Time `json:"call_time" form:"call_time" gorm:"-"`
//vbox.ChannelAccount `gorm:"embedded" json:"acc"`
//vbox.VboxPayOrder   `gorm:"embedded" json:"order"`
