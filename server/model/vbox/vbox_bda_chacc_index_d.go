// 自动生成模板VboxBdaChaccIndexD
package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VboxBdaChaccIndexD 结构体
type VboxBdaChaccIndexD struct {
	global.GVA_MODEL

	Uid             *int     `json:"uid" form:"uid" gorm:"column:uid;comment:用户id;size:10;"`
	UserName        string   `json:"username" form:"username" gorm:"column:username;comment:用户名;size:64;"`
	PAccount        string   `json:"p_account" form:"p_account" gorm:"column:p_account;comment:付方账户名;size:50;"`
	ChannelCode     *int     `json:"channelCode" form:"channelCode" gorm:"column:channel_code;comment:通道code;"`
	ProductId       *int     `json:"productId" form:"productId" gorm:"column:product_id;comment:通道ID;"`
	ProductName     string   `json:"productName" form:"productName" gorm:"column:product_name;comment:通道名;size:128;"`
	OrderQuantify   *int     `json:"orderQuantify" form:"orderQuantify" gorm:"column:order_quantify;comment:订单量;"`
	OkOrderQuantify *int     `json:"okOrderQuantify" form:"okOrderQuantify" gorm:"column:ok_order_quantify;comment:成功订单量;"`
	Ratio           *float64 `json:"ratio" form:"ratio" gorm:"column:ratio;comment:成交率;"`
	Income          *int     `json:"income" form:"income" gorm:"column:income;comment:成交金额;"`
	Dt              string   `json:"dt" form:"dt" gorm:"column:dt;comment:时间yyyy-MM-dd;size:32;"`
	CreatedBy       uint     `gorm:"column:created_by;comment:创建者"`
	UpdatedBy       uint     `gorm:"column:updated_by;comment:更新者"`
	DeletedBy       uint     `gorm:"column:deleted_by;comment:删除者"`
}

// TableName VboxBdaChaccIndexD 表名
func (VboxBdaChaccIndexD) TableName() string {
	return "vbox_bda_chacc_index_d"
}
