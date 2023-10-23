package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"gorm.io/gorm"
)

type ChannelAccountService struct {
}

// CreateChannelAccount 创建通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) CreateChannelAccount(vca *vbox.ChannelAccount) (err error) {
	err = global.GVA_DB.Create(vca).Error
	return err
}

// DeleteChannelAccount 删除通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) DeleteChannelAccount(vca vbox.ChannelAccount) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).Update("deleted_by", vca.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&vca).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteChannelAccountByIds 批量删除通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) DeleteChannelAccountByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.ChannelAccount{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.ChannelAccount{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// SwitchEnableChannelAccount 开关通道账号
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) SwitchEnableChannelAccount(vca vboxReq.ChannelAccountUpd) (err error) {
	err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).Update("status", vca.Status).Error
	return err
}

// SwitchEnableChannelAccountByIds 批量开关通道账号记录
// Author [bboyhan](https://github.com/bboyhan)
func (vcaService *ChannelAccountService) SwitchEnableChannelAccountByIds(upd vboxReq.ChannelAccountUpd, updatedBy uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.ChannelAccount{}).Where("id in ?", upd.Ids).Update("status", upd.Status).Update("updated_by", updatedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", upd.Ids).Updates(&vbox.ChannelAccount{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannelAccount 更新通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) UpdateChannelAccount(vca vbox.ChannelAccount) (err error) {
	err = global.GVA_DB.Save(&vca).Error
	return err
}

// GetChannelAccount 根据id获取通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) GetChannelAccount(id uint) (vca vbox.ChannelAccount, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vca).Error
	return
}

// GetChannelAccountInfoList 分页获取通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) GetChannelAccountInfoList(info vboxReq.ChannelAccountSearch, ids []uint) (list []vbox.ChannelAccount, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelAccount{})
	var vcas []vbox.ChannelAccount
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.AcRemark != "" {
		db = db.Where("ac_remark LIKE ?", "%"+info.AcRemark+"%")
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account LIKE ?", "%"+info.AcAccount+"%")
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Where("created_by in (?)", ids).Find(&vcas).Error
	return vcas, total, err
}
