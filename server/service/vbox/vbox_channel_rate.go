package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"gorm.io/gorm"
)

type VboxChannelRateService struct {
}

// CreateVboxChannelRate 创建VboxChannelRate记录
// Author [piexlmax](https://github.com/piexlmax)
func (chRateService *VboxChannelRateService) CreateVboxChannelRate(chRate *vbox.VboxChannelRate) (err error) {
	err = global.GVA_DB.Create(chRate).Error
	return err
}

// DeleteVboxChannelRate 删除VboxChannelRate记录
// Author [piexlmax](https://github.com/piexlmax)
func (chRateService *VboxChannelRateService) DeleteVboxChannelRate(chRate vbox.VboxChannelRate) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxChannelRate{}).Where("id = ?", chRate.ID).Update("deleted_by", chRate.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&chRate).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteVboxChannelRateByIds 批量删除VboxChannelRate记录
// Author [piexlmax](https://github.com/piexlmax)
func (chRateService *VboxChannelRateService) DeleteVboxChannelRateByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxChannelRate{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxChannelRate{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateVboxChannelRate 更新VboxChannelRate记录
// Author [piexlmax](https://github.com/piexlmax)
func (chRateService *VboxChannelRateService) UpdateVboxChannelRate(chRate vbox.VboxChannelRate) (err error) {
	err = global.GVA_DB.Save(&chRate).Error
	return err
}

// GetVboxChannelRate 根据id获取VboxChannelRate记录
// Author [piexlmax](https://github.com/piexlmax)
func (chRateService *VboxChannelRateService) GetVboxChannelRate(id uint) (chRate vbox.VboxChannelRate, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&chRate).Error
	return
}

// GetVboxChannelRateInfoList 分页获取VboxChannelRate记录
// Author [piexlmax](https://github.com/piexlmax)
func (chRateService *VboxChannelRateService) GetVboxChannelRateInfoList(info vboxReq.VboxChannelRateSearch) (list []vbox.VboxChannelRate, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxChannelRate{})
	var chRates []vbox.VboxChannelRate
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&chRates).Error
	return chRates, total, err
}

func (chRateService *VboxChannelRateService) GetVboxTeamUserChannelRateList(info vboxReq.VboxChannelRateSearch) (list []vbox.UserChannelProductRate, total int64, err error) {
	channelProducts, err := GetVboxChannelProductInfoList()
	rateChannelProducts := make([]vbox.UserChannelProductRate, 0)

	// todo
	channelCodeToRate := map[string]float64{
		"6000": 0.5,
		"6001": 0.8,
		"6002": 0.9,
		"6003": 0.7,
		"6006": 0.3,
	}

	processRateChannelProducts(channelProducts, channelCodeToRate, &rateChannelProducts)
	return rateChannelProducts, int64(len(rateChannelProducts)), err
}

func GetVboxChannelProductInfoList() (list []vbox.ChannelProduct, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelProduct{})
	if err = db.Where("parent_id = ?", "0").Error; err != nil {
		return
	}
	var vcps []vbox.ChannelProduct
	err = db.Where("parent_id = ?", "0").Find(&vcps).Error
	for k := range vcps {
		err = findChildrenChannelProduct(&vcps[k])
	}
	return vcps, err
}

func findChildrenChannelProduct(vcp *vbox.ChannelProduct) (err error) {
	err = global.GVA_DB.Where("parent_id = ?", vcp.ChannelCode).Find(&vcp.Children).Error
	if len(vcp.Children) > 0 {
		for k := range vcp.Children {
			err = findChildrenChannelProduct(&vcp.Children[k])
		}
	}
	return err
}

func processRateChannelProducts(channelProducts []vbox.ChannelProduct, channelCodeToRate map[string]float64, result *[]vbox.UserChannelProductRate) {
	for _, cp := range channelProducts {
		rate, exists := channelCodeToRate[cp.ChannelCode]
		if !exists {
			rate = 0.0 // 如果没有设置对应的 rate，默认为 0.0
		}

		rcp := vbox.UserChannelProductRate{
			ChannelProduct: cp,
			Rate:           rate,
		}

		*result = append(*result, rcp)

		if len(cp.Children) > 0 {
			processRateChannelProducts(cp.Children, channelCodeToRate, result)
		}
	}
}
