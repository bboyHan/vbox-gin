// 自动生成模板VboxPayAccount
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// PayAccount 结构体
type PayAccount struct {
	global.GVA_MODEL
	Uid      uint   `json:"uid" form:"uid" gorm:"column:uid;comment:UID;"`
	PAccount string `json:"pAccount" form:"pAccount" gorm:"column:p_account;comment:付方账户名;size:50;"`
	PKey     string `json:"pKey" form:"pKey" gorm:"column:p_key;comment:付方Key;size:50;"`
	PRemark  string `json:"pRemark" form:"pRemark" gorm:"column:p_remark;comment:付方备注;size:50;"`
	Status   int    `json:"status" form:"status" gorm:"column:status;comment:状态开关;"`
}

// TableName VboxPayAccount 表名
func (PayAccount) TableName() string {
	return "vbox_pay_account"
}
