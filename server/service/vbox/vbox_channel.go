package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/channel"
	channelReq "github.com/flipped-aurora/gin-vue-admin/server/model/channel/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"gorm.io/gorm"
)

type ChannelService struct {
}

// CreateChannel 创建Channel记录
// Author [piexlmax](https://github.com/piexlmax)
func (chService *ChannelService) CreateChannel(ch *channel.Channel) (err error) {
	err = global.GVA_DB.Create(ch).Error
	return err
}

// DeleteChannel 删除Channel记录
// Author [piexlmax](https://github.com/piexlmax)
func (chService *ChannelService) DeleteChannel(ch channel.Channel) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&channel.Channel{}).Where("id = ?", ch.ID).Update("deleted_by", ch.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&ch).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteChannelByIds 批量删除Channel记录
// Author [piexlmax](https://github.com/piexlmax)
func (chService *ChannelService) DeleteChannelByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&channel.Channel{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&channel.Channel{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannel 更新Channel记录
// Author [piexlmax](https://github.com/piexlmax)
func (chService *ChannelService) UpdateChannel(ch channel.Channel) (err error) {
	err = global.GVA_DB.Save(&ch).Error
	return err
}

// GetChannel 根据id获取Channel记录
// Author [piexlmax](https://github.com/piexlmax)
func (chService *ChannelService) GetChannel(id uint) (ch channel.Channel, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&ch).Error
	return
}

// GetChannelByChannel 根据channel获取Channel记录
// Author yoga
func (chService *ChannelService) GetChannelByChannel(channel string) (ch channel.Channel, err error) {
	err = global.GVA_DB.Where("c_channel_id = ?", channel).First(&ch).Error
	return
}

// GetChannelInfoList 分页获取Channel记录
// Author [piexlmax](https://github.com/piexlmax)
func (chService *ChannelService) GetChannelInfoList(info channelReq.ChannelSearch) (list []channel.Channel, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&channel.Channel{})
	var chs []channel.Channel
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&chs).Error
	return chs, total, err
}
