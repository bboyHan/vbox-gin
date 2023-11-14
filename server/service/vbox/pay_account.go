package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"gorm.io/gorm"
	"strings"
)

type PayAccountService struct {
}

// CreatePayAccount 创建付方记录
// Author [piexlmax](https://github.com/piexlmax)
func (paccService *PayAccountService) CreatePayAccount(pacc *vbox.PayAccount) (err error) {
	pacc.PAccount = rand_string.RandomInt(8)
	pacc.PKey = strings.ToLower(rand_string.RandomLetter(32))
	pacc.Status = 1
	err = global.GVA_DB.Create(pacc).Error
	return err
}

// DeletePayAccount 删除付方记录
// Author [piexlmax](https://github.com/piexlmax)
func (paccService *PayAccountService) DeletePayAccount(pacc vbox.PayAccount) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.PayAccount{}).Where("id = ?", pacc.ID).Update("deleted_by", pacc.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&pacc).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeletePayAccountByIds 批量删除付方记录
// Author [piexlmax](https://github.com/piexlmax)
func (paccService *PayAccountService) DeletePayAccountByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.PayAccount{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.PayAccount{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdatePayAccount 更新付方记录
// Author [piexlmax](https://github.com/piexlmax)
func (paccService *PayAccountService) UpdatePayAccount(pacc vbox.PayAccount) (err error) {
	err = global.GVA_DB.Save(&pacc).Error
	return err
}

// GetPayAccount 根据id获取付方记录
// Author [piexlmax](https://github.com/piexlmax)
func (paccService *PayAccountService) GetPayAccount(id uint) (pacc vbox.PayAccount, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&pacc).Error
	return
}

// GetPayAccountInfoList 分页获取付方记录
// Author [piexlmax](https://github.com/piexlmax)
func (paccService *PayAccountService) GetPayAccountInfoList(info vboxReq.PayAccountSearch, ids []uint) (list []vbox.PayAccount, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayAccount{})
	var paccs []vbox.PayAccount
	db.Where("created_by in (?)", ids)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.PRemark != "" {
		db = db.Where("p_remark LIKE ?", "%"+info.PRemark+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&paccs).Error
	return paccs, total, err
}

// GetPAccGateway 根据id获取信道记录
func (paccService *PayAccountService) GetPAccGateway(req vboxReq.VboxProxySearch) (vboxProxy vbox.Proxy, err error) {
	Chan := req.Chan
	err = global.GVA_DB.Where("chan = ?", Chan).First(&vboxProxy).Error
	return
}
