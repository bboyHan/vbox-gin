package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"gorm.io/gorm"
)

type VboxUserWalletService struct {
}

// CreateVboxUserWallet 创建VboxUserWallet记录
// Author [piexlmax](https://github.com/piexlmax)
func (vuwService *VboxUserWalletService) CreateVboxUserWallet(vuw *vbox.VboxUserWallet) (err error) {
	err = global.GVA_DB.Create(vuw).Error
	return err
}

// DeleteVboxUserWallet 删除VboxUserWallet记录
// Author [piexlmax](https://github.com/piexlmax)
func (vuwService *VboxUserWalletService) DeleteVboxUserWallet(vuw vbox.VboxUserWallet) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxUserWallet{}).Where("id = ?", vuw.ID).Update("deleted_by", vuw.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&vuw).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteVboxUserWalletByIds 批量删除VboxUserWallet记录
// Author [piexlmax](https://github.com/piexlmax)
func (vuwService *VboxUserWalletService) DeleteVboxUserWalletByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxUserWallet{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxUserWallet{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateVboxUserWallet 更新VboxUserWallet记录
// Author [piexlmax](https://github.com/piexlmax)
func (vuwService *VboxUserWalletService) UpdateVboxUserWallet(vuw vbox.VboxUserWallet) (err error) {
	err = global.GVA_DB.Save(&vuw).Error
	return err
}

// GetVboxUserWallet 根据id获取VboxUserWallet记录
// Author [piexlmax](https://github.com/piexlmax)
func (vuwService *VboxUserWalletService) GetVboxUserWallet(id uint) (vuw vbox.VboxUserWallet, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vuw).Error
	return
}

// GetVboxUserWalletInfoList 分页获取VboxUserWallet记录
// Author [piexlmax](https://github.com/piexlmax)
func (vuwService *VboxUserWalletService) GetVboxUserWalletInfoList(info vboxReq.VboxUserWalletSearch, ids []int) (list []vbox.VboxUserWallet, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxUserWallet{})
	var vuws []vbox.VboxUserWallet
	db = db.Where("uid in (?)", ids).Order("create_time desc").Find(&vuws)

	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.UserName != "" {
		db = db.Where("username = ?", info.UserName)
	}

	if info.Remark != "" {
		db = db.Where("remark = ?", info.Remark)
	}
	if info.CreateTime != nil {
		db = db.Where("create_time > ?", info.CreateTime)
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&vuws).Error
	return vuws, total, err
}

func (vuwService *VboxUserWalletService) GetVboxUserWalletAvailablePoints(uid uint) (rechargeData int, err error) {
	queryA := `
        SELECT coalesce(sum(recharge),0) as recharge
		FROM vbox_user_wallet
		WHERE  uid = ? ;
    `

	queryB := `
        SELECT coalesce(sum(cost),0) as recharge
		from (
		   SELECT (money * coalesce(unit_price,0)) as cost
		FROM vbox_pay_order
		WHERE  uid = ? and order_status = ?
		)t;
    `

	resultA := 0
	err = global.GVA_DB.Model(&vbox.VboxUserWallet{}).Raw(queryA, uid).Scan(&resultA).Error
	if err != nil {
		return
	}

	resultB := 0
	err = global.GVA_DB.Model(&vbox.VboxPayOrder{}).Raw(queryB, uid, 1).Scan(&resultB).Error
	if err != nil {
		return
	}
	rechargeData = resultA - resultB
	return rechargeData, err

}
