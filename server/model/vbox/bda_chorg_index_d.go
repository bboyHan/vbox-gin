package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// VboxBdaChorgIndexD 结构体  通道团队统计-天更新
type BdaChorgIndexD struct {
	global.GVA_MODEL
	OrganizationId   int     `json:"organizationId" form:"organizationId" gorm:"column:organization_id;comment:团队id;size:64;"`      //团队id
	OrganizationName string  `json:"organizationName" form:"organizationName" gorm:"column:organization_name;comment:团队名;size:64;"` //团队名
	Cid              string  `json:"cid" form:"cid" gorm:"column:cid;comment:通道ID;size:32;"`                                        //通道ID
	ChannelCode      string  `json:"channelCode" form:"channelCode" gorm:"column:channel_code;comment:通道code;size:20;"`             //通道code
	ProductId        string  `json:"productId" form:"productId" gorm:"column:product_id;comment:产品ID;size:128;"`                    //产品ID
	ProductName      string  `json:"productName" form:"productName" gorm:"column:product_name;comment:产品名称;size:128;"`              //产品名称
	OrderQuantify    int     `json:"orderQuantify" form:"orderQuantify" gorm:"column:order_quantify;comment:成功订单量;size:64;"`        //订单量
	OkOrderQuantify  int     `json:"okOrderQuantify" form:"okOrderQuantify" gorm:"column:ok_order_quantify;comment:成功订单量;"`         //成功订单量
	Ratio            float64 `json:"ratio" form:"ratio" gorm:"column:ratio;comment:成交率;size:64;"`                                   //成交率
	Income           int     `json:"income" form:"income" gorm:"column:income;comment:成交金额;size:64;"`                               //成交金额
	Dt               string  `json:"dt" form:"dt" gorm:"column:dt;comment:天;size:32;"`                                              //天
	CreatedBy        uint    `gorm:"column:created_by;comment:创建者"`
	UpdatedBy        uint    `gorm:"column:updated_by;comment:更新者"`
	DeletedBy        uint    `gorm:"column:deleted_by;comment:删除者"`
}

// TableName 通道团队统计-天更新 VboxBdaChorgIndexD自定义表名 vbox_bda_chorg_index_d
func (BdaChorgIndexD) TableName() string {
	return "vbox_bda_chorg_index_d"
}
