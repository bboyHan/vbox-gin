package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"gorm.io/gorm"
)

type ChannelProductService struct {
}

// CreateChannelProduct 创建通道产品记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) CreateChannelProduct(vcp *vbox.ChannelProduct) (err error) {
	err = global.GVA_DB.Create(vcp).Error
	return err
}

// DeleteChannelProduct 删除通道产品记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) DeleteChannelProduct(vcp vbox.ChannelProduct) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.ChannelProduct{}).Where("id = ?", vcp.ID).Update("deleted_by", vcp.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&vcp).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteChannelProductByIds 批量删除通道产品记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) DeleteChannelProductByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.ChannelProduct{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.ChannelProduct{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannelProduct 更新通道产品记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) UpdateChannelProduct(vcp vbox.ChannelProduct) (err error) {
	err = global.GVA_DB.Save(&vcp).Error
	return err
}

// GetChannelProduct 根据id获取通道产品记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) GetChannelProduct(id uint) (vcp vbox.ChannelProduct, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vcp).Error
	return
}

// GetChannelProductSelf 获取ChannelProduct所有列表(当前用户组织下的)
func (vcpService *ChannelProductService) GetChannelProductSelf(ids []uint, search vboxReq.ChannelProductSearch) (list []vbox.ChannelProduct, err error) {
	// 创建db
	var productIds []uint
	db := global.GVA_DB.Model(&vbox.OrgProduct{})
	if err = db.Distinct("channel_product_id").Select("channel_product_id").Where("organization_id in ?", ids).Find(&productIds).Error; err != nil {
		return
	}

	db = global.GVA_DB.Model(&vbox.ChannelProduct{})
	var channelProducts []vbox.ChannelProduct

	if search.Type != 0 {
		db.Where("type = ?", search.Type)
	}

	err = db.Select("vbox_channel_product.*").
		Where("id in ?", productIds).Find(&channelProducts).Error

	channelProducts = SetChildren(channelProducts)
	return channelProducts, err
}

// GetChannelProductInfoList 分页获取通道产品记录
func (vcpService *ChannelProductService) GetChannelProductInfoList(info vboxReq.ChannelProductSearch) (list []vbox.ChannelProduct, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelProduct{})
	if err = db.Where("parent_id = ?", "0").Count(&total).Error; total == 0 || err != nil {
		return
	}
	var vcps []vbox.ChannelProduct
	err = db.Limit(limit).Offset(offset).Where("parent_id = ?", "0").Find(&vcps).Error
	for k := range vcps {
		err = vcpService.findChildrenChannelProduct(&vcps[k])
	}
	return vcps, total, err
}

// GetChannelProductAll 获取VboxChannelProduct所有记录
func (vcpService *ChannelProductService) GetChannelProductAll() (list []vbox.ChannelProduct, err error) {
	db := global.GVA_DB.Model(&vbox.ChannelProduct{})
	var vcps []vbox.ChannelProduct
	err = db.Find(&vcps).Error
	return vcps, err
}

// @author: [piexlmax](https://github.com/piexlmax)
// @function: findChildrenAuthority
// @description: 查询子角色
// @param: authority *model.SysAuthority
// @return: err error
func (vcpService *ChannelProductService) findChildrenChannelProduct(vcp *vbox.ChannelProduct) (err error) {
	err = global.GVA_DB.Where("parent_id = ?", vcp.ChannelCode).Find(&vcp.Children).Error
	if len(vcp.Children) > 0 {
		for k := range vcp.Children {
			err = vcpService.findChildrenChannelProduct(&vcp.Children[k])
		}
	}
	return err
}

func SetChildren(channelProducts []vbox.ChannelProduct) []vbox.ChannelProduct {
	if len(channelProducts) == 0 {
		return nil
	}

	productMap := make(map[string]*vbox.ChannelProduct)

	var rootProducts []vbox.ChannelProduct

	for i := range channelProducts {
		product := &channelProducts[i]
		code := product.ChannelCode
		parentId := product.ParentId

		if parentId == "0" {
			rootProducts = append(rootProducts, *product)
		} else {
			parent, ok := productMap[parentId]
			if !ok {
				parent = &vbox.ChannelProduct{
					ChannelCode: parentId,
					Children:    []vbox.ChannelProduct{},
				}
				productMap[parentId] = parent
			}
			parent.Children = append(parent.Children, *product)
		}

		productMap[code] = product
	}

	for i := range rootProducts {
		rootProduct := &rootProducts[i]
		if children, exist := productMap[rootProduct.ChannelCode]; exist {
			rootProduct.Children = children.Children
		}
	}

	return rootProducts
}
