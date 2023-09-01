package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox_channel_guideimg"
	vbox_channel_guideimgReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox_channel_guideimg/request"
	"gorm.io/gorm"
	"sort"
)

type ChannelGuideImgService struct {
}

// CreateChannelGuideimg 创建Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) CreateChannelGuideimg(chGuideImg *vbox_channel_guideimg.ChannelGuideImg) (err error) {
	err = global.GVA_DB.Create(chGuideImg).Error
	return err
}

// DeleteChannelGuideimg 删除Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) DeleteChannelGuideimg(chGuideImg vbox_channel_guideimg.ChannelGuideImg) (err error) {
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

// DeleteChannelGuideimgByIds 批量删除Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) DeleteChannelGuideimgByIds(ids request.IdsReq, deleted_by uint) (err error) {
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

// UpdateChannelGuideimg 更新Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) UpdateChannelGuideimg(chGuideImg vbox_channel_guideimg.ChannelGuideImg) (err error) {
	err = global.GVA_DB.Save(&chGuideImg).Error
	return err
}

// GetChannelGuideimg 根据id获取Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) GetChannelGuideimg(id uint) (chGuideImg vbox_channel_guideimg.ChannelGuideImg, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&chGuideImg).Error
	return
}

// GetChannelGuideimgInfoList 分页获取Channel_guideimg记录
// Author [piexlmax](https://github.com/piexlmax)
func (chGuideImgService *ChannelGuideImgService) GetChannelGuideimgInfoList(info vbox_channel_guideimgReq.ChannelGuideImgSearch) (list []vbox_channel_guideimg.ChannelGuideImg, total int64, err error) {
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

func (chGuideImgService *ChannelGuideImgService) GetChannelGuideImgTaskList(info vbox_channel_guideimgReq.ChannelGuideImgTask) (list []vbox_channel_guideimg.ChannelGuideImg, total int64, err error) {
	//limit := info.PageSize
	//offset := info.PageSize * (info.Page - 1)
	//fmt.Println(info.ChannelId)

	global.GVA_LOG.Info(info.ChannelId)
	// 创建db
	db := global.GVA_DB.Model(&vbox_channel_guideimg.ChannelGuideImg{})
	var chGuideImgs []vbox_channel_guideimg.ChannelGuideImg
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.ChannelId != "" {
		db = db.Where("c_channel_id = ?", info.ChannelId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Find(&chGuideImgs).Error
	sortChannelGuideImgs(chGuideImgs)
	return chGuideImgs, total, err
}

// 定义自定义排序函数
func sortByNum(a, b *vbox_channel_guideimg.ChannelGuideImg) bool {
	return *a.ImgNum < *b.ImgNum
}

// 对切片进行排序
func sortChannelGuideImgs(chGuideImgs []vbox_channel_guideimg.ChannelGuideImg) {
	sort.Slice(chGuideImgs, func(i, j int) bool {
		return sortByNum(&chGuideImgs[i], &chGuideImgs[j])
	})
}
