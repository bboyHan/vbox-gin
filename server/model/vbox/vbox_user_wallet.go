// 自动生成模板VboxUserWallet
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// VboxUserWallet 结构体
type VboxUserWallet struct {
	global.GVA_MODEL
	Uid        uint       `json:"uid" form:"uid" gorm:"column:uid;comment:父角色ID;"`
	UserName   string     `json:"username" form:"username" gorm:"column:username;comment:默认菜单;size:64;"`
	Recharge   int64      `json:"recharge" form:"recharge" gorm:"column:recharge;comment:角色ID;"`
	Tariff     float64    `json:"tariff" form:"tariff" gorm:"column:tariff;comment:角色名;size:20;"`
	Remark     string     `json:"remark" form:"remark" gorm:"column:remark;comment:附加参数;size:50;"`
	CreateTime *time.Time `json:"createTime" form:"createTime" gorm:"column:create_time;comment:创建时间;"`
	VipLevel   int64      `json:"vipLevel" form:"vipLevel" gorm:"column:vip_level;comment:会员等级;"`
	CreatedBy  uint       `gorm:"column:created_by;comment:创建者"`
	UpdatedBy  uint       `gorm:"column:updated_by;comment:更新者"`
	DeletedBy  uint       `gorm:"column:deleted_by;comment:删除者"`
}

// TableName VboxUserWallet 表名
func (VboxUserWallet) TableName() string {
	return "vbox_user_wallet"
}
