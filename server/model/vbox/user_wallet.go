// 自动生成模板UserWallet
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// UserWallet 结构体  用户钱包
type UserWallet struct {
	global.GVA_MODEL
	Uid       uint   `json:"uid" form:"uid" gorm:"column:uid;comment:用户ID;size:16;"`               //用户ID
	Username  string `json:"username" form:"username" gorm:"column:username;comment:用户名;size:64;"` //用户名
	Recharge  int    `json:"recharge" form:"recharge" gorm:"column:recharge;comment:积分;size:32;"`  //积分
	EventId   string `json:"eventId" form:"eventId" gorm:"column:event_id;comment:事件ID;size:64;"`  //事件ID
	Type      int    `json:"type" form:"type" gorm:"column:type;comment:事件类型,0-充值,1-消费;size:4;"`   //事件类型
	Remark    string `json:"remark" form:"remark" gorm:"column:remark;comment:附加参数;size:256;"`     //说明
	CreatedBy uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 用户钱包 UserWallet自定义表名 vbox_user_wallet
func (UserWallet) TableName() string {
	return "vbox_user_wallet"
}
