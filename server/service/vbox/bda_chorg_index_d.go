package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"gorm.io/gorm"
)

type BdaChorgIndexDService struct {
}

// CreateBdaChorgIndexD 创建通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) CreateBdaChorgIndexD(bdaChorg *vbox.BdaChorgIndexD) (err error) {
	err = global.GVA_DB.Create(bdaChorg).Error
	return err
}

// DeleteBdaChorgIndexD 删除通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) DeleteBdaChorgIndexD(bdaChorg vbox.BdaChorgIndexD) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChorgIndexD{}).Where("id = ?", bdaChorg.ID).Update("deleted_by", bdaChorg.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&bdaChorg).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteBdaChorgIndexDByIds 批量删除通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) DeleteBdaChorgIndexDByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChorgIndexD{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.BdaChorgIndexD{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateBdaChorgIndexD 更新通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) UpdateBdaChorgIndexD(bdaChorg vbox.BdaChorgIndexD) (err error) {
	err = global.GVA_DB.Save(&bdaChorg).Error
	return err
}

// GetBdaChorgIndexD 根据id获取通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) GetBdaChorgIndexD(id uint) (bdaChorg vbox.BdaChorgIndexD, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&bdaChorg).Error
	return
}

// GetBdaChorgIndexDInfoList 分页获取通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) GetBdaChorgIndexDInfoList(info vboxReq.BdaChorgIndexDSearch) (list []vbox.BdaChorgIndexD, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.BdaChorgIndexD{})
	var bdaChorgs []vbox.BdaChorgIndexD
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.OrganizationId != 0 {
		db = db.Where("organization_id = ?", info.OrganizationId)
	}
	if info.OrganizationName != "" {
		db = db.Where("organization_name LIKE ?", "%"+info.OrganizationName+"%")
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.ChannelCode != "" {
		db = db.Where("channel_code = ?", info.ChannelCode)
	}
	if info.ProductId != "" {
		db = db.Where("product_id = ?", info.ProductId)
	}
	if info.ProductName != "" {
		db = db.Where("product_name LIKE ?", "%"+info.ProductName+"%")
	}
	if info.Dt != "" {
		db = db.Where("dt = ?", info.Dt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&bdaChorgs).Error
	return bdaChorgs, total, err
}
