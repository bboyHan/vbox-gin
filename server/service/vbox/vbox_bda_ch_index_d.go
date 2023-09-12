package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
    "gorm.io/gorm"
)

type VboxBdaChIndexDService struct {
}

// CreateVboxBdaChIndexD 创建VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService) CreateVboxBdaChIndexD(bdaChD *vbox.VboxBdaChIndexD) (err error) {
	err = global.GVA_DB.Create(bdaChD).Error
	return err
}

// DeleteVboxBdaChIndexD 删除VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService)DeleteVboxBdaChIndexD(bdaChD vbox.VboxBdaChIndexD) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&vbox.VboxBdaChIndexD{}).Where("id = ?", bdaChD.ID).Update("deleted_by", bdaChD.DeletedBy).Error; err != nil {
              return err
        }
        if err = tx.Delete(&bdaChD).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteVboxBdaChIndexDByIds 批量删除VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService)DeleteVboxBdaChIndexDByIds(ids request.IdsReq,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&vbox.VboxBdaChIndexD{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxBdaChIndexD{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateVboxBdaChIndexD 更新VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService)UpdateVboxBdaChIndexD(bdaChD vbox.VboxBdaChIndexD) (err error) {
	err = global.GVA_DB.Save(&bdaChD).Error
	return err
}

// GetVboxBdaChIndexD 根据id获取VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService)GetVboxBdaChIndexD(id uint) (bdaChD vbox.VboxBdaChIndexD, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&bdaChD).Error
	return
}

// GetVboxBdaChIndexDInfoList 分页获取VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService)GetVboxBdaChIndexDInfoList(info vboxReq.VboxBdaChIndexDSearch) (list []vbox.VboxBdaChIndexD, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&vbox.VboxBdaChIndexD{})
    var bdaChDs []vbox.VboxBdaChIndexD
    // 如果有条件搜索 下方会自动创建搜索语句
    if info.StartCreatedAt !=nil && info.EndCreatedAt !=nil {
     db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
    }
    if info.Uid != nil {
        db = db.Where("uid = ?",info.Uid)
    }
    if info.UserName != "" {
        db = db.Where("username = ?",info.UserName)
    }
    if info.ChannelCode != nil {
        db = db.Where("channel_code = ?",info.ChannelCode)
    }
    if info.ProductId != nil {
        db = db.Where("product_id = ?",info.ProductId)
    }
    if info.ProductName != "" {
        db = db.Where("product_name = ?",info.ProductName)
    }
    if info.Dt != "" {
        db = db.Where("dt = ?",info.Dt)
    }
	err = db.Count(&total).Error
	if err!=nil {
    	return
    }

	err = db.Limit(limit).Offset(offset).Find(&bdaChDs).Error
	return  bdaChDs, total, err
}
