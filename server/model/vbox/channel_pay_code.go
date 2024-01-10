package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// ChannelPayCode 结构体  通道账户付款二维码
type ChannelPayCode struct {
	global.GVA_MODEL
	Cid        string     `json:"cid" form:"cid" gorm:"column:cid;comment:通道id;size:10;"`                                                      //通道id
	AcAccount  string     `json:"acAccount" form:"acAccount" gorm:"column:ac_account;comment:通道账户名;size:128;"`                                 //通道账户名
	AcId       string     `json:"acId" form:"acId" gorm:"column:ac_id;comment:账户ID;size:50;"`                                                  //账户ID
	AcRemark   string     `json:"acRemark" form:"acRemark" gorm:"column:ac_remark;comment:账户备注;size:128;"`                                     //账户备注
	Money      int        `json:"money" form:"money" gorm:"column:money;comment:金额;size:16;"`                                                  //金额
	ExpTime    *time.Time `json:"expTime" form:"expTime" gorm:"column:exp_time;comment:过期时间;"`                                                 //过期时间
	Operator   string     `json:"operator" form:"operator" gorm:"column:operator;comment:运营商;size:64;"`                                        //运营商
	Location   string     `json:"location" form:"location" gorm:"column:location;comment:省市;size:128;"`                                        //省市
	ImgBaseStr string     `json:"imgBaseStr" form:"imgBaseStr" gorm:"column:img_base_str;comment:图片base64编码;type:longtext;"`                   //图片base64编码
	ImgContent string     `json:"imgContent" form:"imgContent" gorm:"column:img_content;comment:图片解析内容;type:longtext;"`                        //图片base64编码
	Mid        string     `json:"mid" form:"mid" gorm:"column:mid;comment:标识id;size:20;"`                                                      //标识id
	PlatId     string     `json:"platId" form:"platId" gorm:"column:plat_id;comment:平台id;size:256;"`                                           //平台id
	CodeStatus uint       `json:"codeStatus" form:"codeStatus" gorm:"default:2;column:code_status;comment:状态,1-已使用,2-待使用,3-已失效，4-等候中;size:2;"` //产码状态
	CreatedBy  uint       `gorm:"column:created_by;comment:创建者"`
	UpdatedBy  uint       `gorm:"column:updated_by;comment:更新者"`
	DeletedBy  uint       `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 通道账户付款二维码 VboxChannelPayCode自定义表名 vbox_channel_pay_code
func (ChannelPayCode) TableName() string {
	return "vbox_channel_pay_code"
}
