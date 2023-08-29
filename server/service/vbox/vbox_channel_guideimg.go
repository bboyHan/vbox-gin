package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox_channel_guideimg"
	vbox_channel_guideimgReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox_channel_guideimg/request"
	"gorm.io/gorm"
)

type ChannelGuideImgService struct {
}

// CreateChannel_guideimg 创建Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) CreateChannel_guideimg(chGuideImg *vbox_channel_guideimg.ChannelGuideImg) (err error) {
	err = global.GVA_DB.Create(chGuideImg).Error
	return err
}

// DeleteChannel_guideimg 删除Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) DeleteChannel_guideimg(chGuideImg vbox_channel_guideimg.ChannelGuideImg) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox_channel_guideimg.ChannelGuideImg{}).Where("id = ?", chGuideImg.ID).Update("deleted_by", chGuideImg.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&chGuideImg).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteChannel_guideimgByIds 批量删除Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) DeleteChannel_guideimgByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox_channel_guideimg.ChannelGuideImg{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox_channel_guideimg.ChannelGuideImg{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannel_guideimg 更新Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) UpdateChannel_guideimg(chGuideImg vbox_channel_guideimg.ChannelGuideImg) (err error) {
	err = global.GVA_DB.Save(&chGuideImg).Error
	return err
}

// GetChannel_guideimg 根据id获取Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) GetChannel_guideimg(id uint) (chGuideImg vbox_channel_guideimg.ChannelGuideImg, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&chGuideImg).Error
	return
}

// GetChannel_guideimgInfoList 分页获取Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) GetChannel_guideimgInfoList(info vbox_channel_guideimgReq.ChannelGuideImgSearch) (list []vbox_channel_guideimg.ChannelGuideImg, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox_channel_guideimg.ChannelGuideImg{})
	var chGuideImgs []vbox_channel_guideimg.ChannelGuideImg
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&chGuideImgs).Error
	return chGuideImgs, total, err
}
