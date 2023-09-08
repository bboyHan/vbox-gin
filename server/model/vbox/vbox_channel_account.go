// 自动生成模板ChannelAccount
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ChannelAccount 结构体
type ChannelAccount struct {
	global.GVA_MODEL
	Uid        uint   `json:"uid" form:"uid" gorm:"column:uid;comment:用户id;size:10;"`
	AcAccount  string `json:"acAccount" form:"acAccount" gorm:"column:ac_account;comment:通道账户名;size:128;"`
	AcPwd      string `json:"acPwd" form:"acPwd" gorm:"column:ac_pwd;comment:通道账户密码;size:100;"`
	AcRemark   string `json:"acRemark" form:"acRemark" gorm:"column:ac_remark;comment:账户备注;size:128;"`
	AcId       int    `json:"acId" form:"acId" gorm:"column:ac_id;comment:账户id;size:50;"`
	Cid        *int   `json:"cid" form:"cid" gorm:"column:cid;comment:通道id;size:10;"`
	Token      string `json:"token" form:"token" gorm:"column:token;comment:ck;size:256;type:text;"`
	DailyLimit *int   `json:"dailyLimit" form:"dailyLimit" gorm:"column:daily_limit;comment:日限额;size:10;"`
	TotalLimit *int   `json:"totalLimit" form:"totalLimit" gorm:"column:total_limit;comment:总限额;size:10;"`
	CountLimit *int   `json:"countLimit" form:"countLimit" gorm:"column:count_limit;comment:笔数限制;size:10;"`
	Status     *int   `json:"status" form:"status" gorm:"column:status;comment:状态开关;size:10;"`
	SysStatus  *int   `json:"sysStatus" form:"sysStatus" gorm:"column:sys_status;comment:系统开关;size:10;"`
	CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName ChannelAccount 表名
func (ChannelAccount) TableName() string {
	return "vbox_channel_account"
}
