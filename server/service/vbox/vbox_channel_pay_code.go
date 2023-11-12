package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"gorm.io/gorm"
)

type VboxChannelPayCodeService struct {
}

// CreateVboxChannelPayCode 创建通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxChannelPayCodeService *VboxChannelPayCodeService) CreateVboxChannelPayCode(vboxChannelPayCode *vbox.VboxChannelPayCode) (err error) {
	err = global.GVA_DB.Create(vboxChannelPayCode).Error
	return err
}

// DeleteVboxChannelPayCode 删除通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxChannelPayCodeService *VboxChannelPayCodeService) DeleteVboxChannelPayCode(vboxChannelPayCode vbox.VboxChannelPayCode) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxChannelPayCode{}).Where("id = ?", vboxChannelPayCode.ID).Update("deleted_by", vboxChannelPayCode.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&vboxChannelPayCode).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteVboxChannelPayCodeByIds 批量删除通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxChannelPayCodeService *VboxChannelPayCodeService) DeleteVboxChannelPayCodeByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxChannelPayCode{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxChannelPayCode{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateVboxChannelPayCode 更新通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxChannelPayCodeService *VboxChannelPayCodeService) UpdateVboxChannelPayCode(vboxChannelPayCode vbox.VboxChannelPayCode) (err error) {
	err = global.GVA_DB.Save(&vboxChannelPayCode).Error
	return err
}

// GetVboxChannelPayCode 根据id获取通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxChannelPayCodeService *VboxChannelPayCodeService) GetVboxChannelPayCode(id uint) (vboxChannelPayCode vbox.VboxChannelPayCode, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vboxChannelPayCode).Error
	return
}

// GetVboxChannelPayCodeInfoList 分页获取通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxChannelPayCodeService *VboxChannelPayCodeService) GetVboxChannelPayCodeInfoList(info vboxReq.VboxChannelPayCodeSearch) (list []vbox.VboxChannelPayCode, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxChannelPayCode{})
	var vboxChannelPayCodes []vbox.VboxChannelPayCode
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.Ac_account != "" {
		db = db.Where("ac_account = ?", info.Ac_account)
	}
	if info.Location != "" {
		db = db.Where("location = ?", info.Location)
	}
	if info.Uid != nil {
		db = db.Where("uid = ?", info.Uid)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&vboxChannelPayCodes).Error
	return vboxChannelPayCodes, total, err
}
