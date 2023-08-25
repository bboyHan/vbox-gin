// 自动生成模板Channel
package channel

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Channel 结构体
type Channel struct {
	global.GVA_MODEL
	CChannelId   string `json:"c_channel_id" form:"c_channel_id" gorm:"column:c_channel_id;comment:通道id;size:50;"`
	CGame        string `json:"c_game" form:"c_game" gorm:"column:c_game;comment:game id;size:50;"`
	CGameName    string `json:"c_game_name" form:"c_game_name" gorm:"column:c_game_name;comment:game描述;size:50;"`
	CChannel     string `json:"c_channel" form:"c_channel" gorm:"column:c_channel;comment:支付通道id;size:50;"`
	CChannelName string `json:"c_channel_name" form:"c_channel_name" gorm:"column:c_channel_name;comment:支付通道描述;size:50;"`
	Type         *int   `json:"type" form:"type" gorm:"column:type;comment:自动或引导;size:11;"`
	CreatedBy    uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy    uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy    uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName Channel 表名
func (Channel) TableName() string {
	return "vbox_channel"
}
