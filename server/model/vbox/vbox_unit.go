package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

//Unit 单价积分 结构体
type Unit struct {
	global.GVA_MODEL
	ScoreId string `json:"scoreId" form:"scoreId" gorm:"column:score_id;comment:积分ID;size:32;"` //积分ID
	Chan    string `json:"chan" form:"chan" gorm:"column:chan;comment:渠道;size:32;"`             //渠道
	Type    int    `json:"type" form:"type" gorm:"column:type;comment:类型;size:4;"`              //类型
	Status  int    `json:"status" form:"status" gorm:"column:status;comment:状态开关;size:4;"`    //状态开关
	OrgID   uint   `json:"orgID" form:"orgID" gorm:"column:org_id;comment:组织ID;size:16;"`       //组织ID
	Score   int    `json:"score" form:"score" gorm:"column:score;comment:积分;size:16;"`          //积分换算比例
	Remark  string `json:"remark" form:"remark" gorm:"column:remark;comment:备注;size:64;"`       //备注
}

// TableName 单价积分 vbox_unit
func (Unit) TableName() string {
	return "vbox_unit"
}
