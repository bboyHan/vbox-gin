// 自动生成模板PayAccount
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// 付方 结构体  PayAccount
type PayAccount struct {
	global.GVA_MODEL
	Uid       uint   `json:"uid" form:"uid" gorm:"column:uid;comment:UID;size:16;"`                   //UID
	PAccount  string `json:"pAccount" form:"pAccount" gorm:"column:p_account;comment:付方账户名;size:64;"` //付方账户名
	PKey      string `json:"pKey" form:"pKey" gorm:"column:p_key;comment:付方Key;size:64;"`             //付方Key
	PRemark   string `json:"pRemark" form:"pRemark" gorm:"column:p_remark;comment:付方备注;size:64;"`     //付方备注
	Status    int    `json:"status" form:"status" gorm:"column:status;comment:状态开关;size:4;"`          //状态开关
	CreatedBy uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 付方 PayAccount自定义表名 vbox_pay_account
func (PayAccount) TableName() string {
	return "vbox_pay_account"
}
