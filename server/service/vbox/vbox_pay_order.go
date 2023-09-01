package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"gorm.io/gorm"
)

type VboxPayOrderService struct {
}

// CreateVboxPayOrder 创建VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) CreateVboxPayOrder(vpo *vbox.VboxPayOrder) (err error) {
	err = global.GVA_DB.Create(vpo).Error
	return err
}

// DeleteVboxPayOrder 删除VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) DeleteVboxPayOrder(vpo vbox.VboxPayOrder) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxPayOrder{}).Where("id = ?", vpo.ID).Update("deleted_by", vpo.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&vpo).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteVboxPayOrderByIds 批量删除VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) DeleteVboxPayOrderByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxPayOrder{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxPayOrder{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateVboxPayOrder 更新VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) UpdateVboxPayOrder(vpo vbox.VboxPayOrder) (err error) {
	err = global.GVA_DB.Save(&vpo).Error
	return err
}

// GetVboxPayOrder 根据id获取VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) GetVboxPayOrder(id uint) (vpo vbox.VboxPayOrder, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vpo).Error
	return
}

// GetVboxPayOrderInfoList 分页获取VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) GetVboxPayOrderInfoList(info vboxReq.VboxPayOrderSearch, ids []int) (list []vbox.VboxPayOrder, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxPayOrder{})
	var vpos []vbox.VboxPayOrder

	db = db.Where("uid in (?)", ids).Find(&vpos)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.OrderId != "" {
		db = db.Where("order_id = ?", info.OrderId)
	}
	if info.PAccount != "" {
		db = db.Where("p_account = ?", info.PAccount)
	}
	if info.Cost != nil {
		db = db.Where("cost = ?", info.Cost)
	}
	if info.Uid != nil {
		db = db.Where("uid = ?", info.Uid)
	}
	if info.AcId != "" {
		db = db.Where("ac_id = ?", info.AcId)
	}
	if info.CChannelId != "" {
		db = db.Where("c_channel_id = ?", info.CChannelId)
	}
	if info.PlatformOid != "" {
		db = db.Where("platform_oid = ?", info.PlatformOid)
	}
	if info.PayIp != "" {
		db = db.Where("pay_ip = ?", info.PayIp)
	}
	if info.PayRegion != "" {
		db = db.Where("pay_region = ?", info.PayRegion)
	}
	if info.ResourceUrl != "" {
		db = db.Where("resource_url = ?", info.ResourceUrl)
	}
	if info.NotifyUrl != "" {
		db = db.Where("notify_url = ?", info.NotifyUrl)
	}
	if info.OrderStatus != nil {
		db = db.Where("order_status = ?", info.OrderStatus)
	}
	if info.CallbackStatus != nil {
		db = db.Where("callback_status = ?", info.CallbackStatus)
	}
	if info.CodeUseStatus != nil {
		db = db.Where("code_use_status = ?", info.CodeUseStatus)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&vpos).Error
	return vpos, total, err
}
