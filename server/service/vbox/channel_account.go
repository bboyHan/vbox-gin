package vbox

import (
	"context"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxResp "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"gorm.io/gorm"
	"time"
)

type ChannelAccountService struct {
}

// QueryAccOrderHis 查询通道账号的官方记录
func (vcaService *ChannelAccountService) QueryAccOrderHis(vca *vbox.ChannelAccount) (res interface{}, err error) {

	rdConn := global.GVA_REDIS.Conn()
	defer rdConn.Close()
	var url string

	c, err := rdConn.Exists(context.Background(), global.ProductRecordQBPrefix).Result()
	if c == 0 {
		var channelCode string
		if global.TxContains(vca.Cid) { // tx系
			channelCode = "qb_proxy"
		}

		err = global.GVA_DB.Debug().Model(&vbox.Proxy{}).Select("url").Where("status = ? and type = ? and chan=?", 1, 1, channelCode).
			First(&url).Error

		if err != nil {
			return nil, errors.New("该信道无资源配置")
		}

		rdConn.Set(context.Background(), global.ProductRecordQBPrefix, url, 10*time.Minute)

	} else {
		url, _ = rdConn.Get(context.Background(), global.ProductRecordQBPrefix).Result()
	}

	if global.TxContains(vca.Cid) { // tx系

		openID, openKey, err := product.Secret(vca.Token)
		if err != nil {
			return nil, err
		}
		records := product.Records(url, openID, openKey, 24*30*time.Hour)
		//classifier := product.Classifier(records.WaterList)
		return records, nil
	}

	return res, err
}

// CountAcc 查询可用通道的 当前等待取用的账号个数
func (vcaService *ChannelAccountService) CountAcc(ids []uint) (res []vboxResp.ChannelAccountUnused, err error) {
	err = global.GVA_DB.Debug().Model(&vbox.ChannelAccount{}).Select("count(1) as total, cid").Where("status = ? and sys_status = ? and created_by in (?)", 1, 1, ids).
		Group("cid").Order("id desc").Find(&res).Error
	return res, err
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
// Author [bboyhan](https://github.com/bboyhan)
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

// GetChannelAccountByAcId 根据AcId获取通道账号记录
func (vcaService *ChannelAccountService) GetChannelAccountByAcId(acId string) (vca vbox.ChannelAccount, err error) {
	err = global.GVA_DB.Where("ac_id = ?", acId).First(&vca).Error
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

	err = db.Where("created_by in (?)", ids).Order("id desc").Find(&vcas).Error
	return vcas, total, err
}
