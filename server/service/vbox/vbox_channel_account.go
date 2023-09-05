package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
    "gorm.io/gorm"
)

type ChannelAccountService struct {
}

// CreateChannelAccount 创建ChannelAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) CreateChannelAccount(vca *vbox.ChannelAccount) (err error) {
	err = global.GVA_DB.Create(vca).Error
	return err
}

// DeleteChannelAccount 删除ChannelAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService)DeleteChannelAccount(vca vbox.ChannelAccount) (err error) {
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

// DeleteChannelAccountByIds 批量删除ChannelAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService)DeleteChannelAccountByIds(ids request.IdsReq,deleted_by uint) (err error) {
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

// UpdateChannelAccount 更新ChannelAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService)UpdateChannelAccount(vca vbox.ChannelAccount) (err error) {
	err = global.GVA_DB.Save(&vca).Error
	return err
}

// GetChannelAccount 根据id获取ChannelAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService)GetChannelAccount(id uint) (vca vbox.ChannelAccount, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vca).Error
	return
}

// GetChannelAccountInfoList 分页获取ChannelAccount记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService)GetChannelAccountInfoList(info vboxReq.ChannelAccountSearch) (list []vbox.ChannelAccount, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&vbox.ChannelAccount{})
    var vcas []vbox.ChannelAccount
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.AcAccount != "" {
        db = db.Where("ac_account LIKE ?","%"+ info.AcAccount+"%")
    }
    if info.AcRemark != "" {
        db = db.Where("ac_remark LIKE ?","%"+ info.AcRemark+"%")
    }
    if info.AcId != nil {
        db = db.Where("ac_id = ?",info.AcId)
    }
    if info.Cid != nil {
        db = db.Where("cid = ?",info.Cid)
    }
        if info.StartStatus != nil && info.EndStatus != nil {
            db = db.Where("status BETWEEN ? AND ? ",info.StartStatus,info.EndStatus)
        }
        if info.StartSysStatus != nil && info.EndSysStatus != nil {
            db = db.Where("sys_status BETWEEN ? AND ? ",info.StartSysStatus,info.EndSysStatus)
        }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	err = db.Limit(limit).Offset(offset).Find(&vcas).Error
	return  vcas, total, err
}
