package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
    vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
    "gorm.io/gorm"
)

type VboxBdaChaccIndexDService struct {
}

// CreateVboxBdaChaccIndexD 创建VboxBdaChaccIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccDService *VboxBdaChaccIndexDService) CreateVboxBdaChaccIndexD(bdaChaccD *vbox.VboxBdaChaccIndexD) (err error) {
	err = global.GVA_DB.Create(bdaChaccD).Error
	return err
}

// DeleteVboxBdaChaccIndexD 删除VboxBdaChaccIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccDService *VboxBdaChaccIndexDService)DeleteVboxBdaChaccIndexD(bdaChaccD vbox.VboxBdaChaccIndexD) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&vbox.VboxBdaChaccIndexD{}).Where("id = ?", bdaChaccD.ID).Update("deleted_by", bdaChaccD.DeletedBy).Error; err != nil {
              return err
        }
        if err = tx.Delete(&bdaChaccD).Error; err != nil {
              return err
        }
        return nil
	})
	return err
}

// DeleteVboxBdaChaccIndexDByIds 批量删除VboxBdaChaccIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccDService *VboxBdaChaccIndexDService)DeleteVboxBdaChaccIndexDByIds(ids request.IdsReq,deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
	    if err := tx.Model(&vbox.VboxBdaChaccIndexD{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
            return err
        }
        if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxBdaChaccIndexD{}).Error; err != nil {
            return err
        }
        return nil
    })
	return err
}

// UpdateVboxBdaChaccIndexD 更新VboxBdaChaccIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccDService *VboxBdaChaccIndexDService)UpdateVboxBdaChaccIndexD(bdaChaccD vbox.VboxBdaChaccIndexD) (err error) {
	err = global.GVA_DB.Save(&bdaChaccD).Error
	return err
}

// GetVboxBdaChaccIndexD 根据id获取VboxBdaChaccIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccDService *VboxBdaChaccIndexDService)GetVboxBdaChaccIndexD(id uint) (bdaChaccD vbox.VboxBdaChaccIndexD, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&bdaChaccD).Error
	return
}

// GetVboxBdaChaccIndexDInfoList 分页获取VboxBdaChaccIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccDService *VboxBdaChaccIndexDService)GetVboxBdaChaccIndexDInfoList(info vboxReq.VboxBdaChaccIndexDSearch) (list []vbox.VboxBdaChaccIndexD, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
    // 创建db
	db := global.GVA_DB.Model(&vbox.VboxBdaChaccIndexD{})
    var bdaChaccDs []vbox.VboxBdaChaccIndexD
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
    if info.PAccount != "" {
        db = db.Where("p_account = ?",info.PAccount)
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

	err = db.Limit(limit).Offset(offset).Find(&bdaChaccDs).Error
	return  bdaChaccDs, total, err
}
