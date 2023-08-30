// 自动生成模板Channel_guideimg
package vbox_channel_guideimg

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

// Channel_guideimg 结构体
type ChannelGuideImg struct {
	global.GVA_MODEL
	ChannelId  string `json:"c_channel_id" form:"c_channel_id" gorm:"column:c_channel_id;comment:通道ID;size:50;"`
	ImgBaseStr string `json:"img_base_str" form:"img_base_str" gorm:"column:img_base_str;comment:图片base64编码;type:text;"`
	ImgNum     *int   `json:"img_num" form:"img_num" gorm:"column:img_num;comment:通道;size:2;"`
	FileName   string `json:"file_name" form:"file_name" gorm:"column:file_name;comment:店铺备注;size:200;"`
	url        string `json:"url" form:"url" gorm:"column:url;comment:店地址;size:200;"`
	tag        string `json:"tag" form:"tag" gorm:"column:tag;comment:金额;size:200;"`
	key        string `json:"key" form:"key" gorm:"column:key;comment:开关;size:200;"`
	CreatedBy  uint   `gorm:"column:created_by;comment:创建者"`
	UpdatedBy  uint   `gorm:"column:updated_by;comment:更新者"`
	DeletedBy  uint   `gorm:"column:deleted_by;comment:删除者"`
}

// TableName Channel_guideimg 表名
func (ChannelGuideImg) TableName() string {
	return "vbox_channel_guideimg"
}
