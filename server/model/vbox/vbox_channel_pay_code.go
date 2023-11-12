// 自动生成模板VboxChannelPayCode
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VboxChannelPayCode 结构体  通道账户付款二维码
type VboxChannelPayCode struct {
	global.GVA_MODEL
	Cid          string `json:"cid" form:"cid" gorm:"column:cid;comment:通道id;size:10;"`                                    //通道id
	Ac_account   string `json:"acAccount" form:"acAccount" gorm:"column:ac_account;comment:通道账户名;size:128;"`               //通道账户名
	Time_limit   string `json:"timeLimit" form:"timeLimit" gorm:"column:time_limit;comment:过期时间;size:64;"`                 //过期时间
	Operator     string `json:"operator" form:"operator" gorm:"column:operator;comment:运营商;size:64;"`                      //运营商
	Location     string `json:"location" form:"location" gorm:"column:location;comment:省市;size:128;"`                      //省市
	Img_base_str string `json:"imgBaseStr" form:"imgBaseStr" gorm:"column:img_base_str;comment:图片base64编码;type:longtext;"` //图片base64编码
	Uid          *int   `json:"uid" form:"uid" gorm:"column:uid;comment:用户id;size:20;"`                                    //用户id
	CreatedBy    uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy    uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy    uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 通道账户付款二维码 VboxChannelPayCode自定义表名 vbox_channel_pay_code
func (VboxChannelPayCode) TableName() string {
	return "vbox_channel_pay_code"
}
