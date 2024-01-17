package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// ChannelAccount 通道账号 结构体
type ChannelAccount struct {
	global.GVA_MODEL
	AcId       string `json:"acId" form:"acId" gorm:"column:ac_id;comment:账户ID;size:50;"`                                 //账户ID
	AcRemark   string `json:"acRemark" form:"acRemark" gorm:"column:ac_remark;comment:账户备注;size:128;"`                    //账户备注
	AcAccount  string `json:"acAccount" form:"acAccount" gorm:"column:ac_account;comment:通道账户名;size:128;"`                //通道账户名
	AcPwd      string `json:"acPwd" form:"acPwd" gorm:"column:ac_pwd;comment:通道账户密码;size:100;"`                           //通道账户密码
	Token      string `json:"token" form:"token" gorm:"column:token;comment:ck;type:text;"`                               //ck
	Cid        string `json:"cid" form:"cid" gorm:"column:cid;comment:通道id;size:10;"`                                     //通道id
	CountLimit int    `json:"countLimit" form:"countLimit" gorm:"column:count_limit;comment:笔数限制;size:10;"`               //笔数限制
	DailyLimit int    `json:"dailyLimit" form:"dailyLimit" gorm:"column:daily_limit;comment:日限额;size:10;"`                //日限额
	TotalLimit int    `json:"totalLimit" form:"totalLimit" gorm:"column:total_limit;comment:总限额;size:10;"`                //总限额
	Status     int    `json:"status" form:"status" gorm:"column:status;comment:状态开关;size:2;"`                             //状态开关
	SysStatus  int    `json:"sysStatus" form:"sysStatus" gorm:"column:sys_status;comment:系统开关,0-关闭,1-开启,2-删除中;size:2;"`   //系统开关
	CdStatus   int    `json:"cdStatus" form:"cdStatus" gorm:"default:1;column:cd_status;comment:是否冷却,1-默认,2-冷却中;size:2;"` //冷却状态
	Uid        int    `json:"uid" form:"uid" gorm:"column:uid;comment:用户id;"`                                             //用户id
	CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 通道账号 ChannelAccount自定义表名 vbox_channel_account
func (ChannelAccount) TableName() string {
	return "vbox_channel_account"
}
