package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"strings"
)

type PayAccountService struct {
}

// CreateVboxPayAccount 创建VboxPayAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpaService *PayAccountService) CreateVboxPayAccount(vpa *vbox.PayAccount) (err error) {
	vpa.PAccount = rand_string.RandomInt(8)
	vpa.PKey = strings.ToLower(rand_string.RandomLetter(32))
	vpa.Status = 1
	err = global.GVA_DB.Create(vpa).Error
	return err
}

// DeleteVboxPayAccount 删除VboxPayAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpaService *PayAccountService) DeleteVboxPayAccount(vpa vbox.PayAccount) (err error) {
	err = global.GVA_DB.Delete(&vpa).Error
	return err
}

// DeleteVboxPayAccountByIds 批量删除VboxPayAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpaService *PayAccountService) DeleteVboxPayAccountByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]vbox.PayAccount{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateVboxPayAccount 更新VboxPayAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpaService *PayAccountService) UpdateVboxPayAccount(vpa vbox.PayAccount) (err error) {
	err = global.GVA_DB.Save(&vpa).Error
	return err
}

// GetVboxPayAccount 根据id获取VboxPayAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpaService *PayAccountService) GetVboxPayAccount(id uint) (vpa vbox.PayAccount, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vpa).Error
	return
}

// GetVboxPayAccountInfoList 分页获取VboxPayAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpaService *PayAccountService) GetVboxPayAccountInfoList(info vboxReq.VboxPayAccountSearch) (list []vbox.PayAccount, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayAccount{})
	var vpas []vbox.PayAccount
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&vpas).Error
	return vpas, total, err
}
