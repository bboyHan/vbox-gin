package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
)

type OrgProductService struct {
}

func (opService *OrgProductService) CreateOrgUser(orgProduct vbox.OrgProductReq) error {
	var Products []vbox.OrgProduct
	var CProducts []vbox.OrgProduct
	err := global.GVA_DB.Find(&Products, "organization_id = ?", orgProduct.OrganizationID).Error
	if err != nil {
		return err
	}
	var ProductMap = make(map[uint]bool)
	for i := range Products {
		ProductMap[Products[i].ChannelProductID] = true
	}

	for i := range orgProduct.ChannelProductIDS {
		if !ProductMap[orgProduct.ChannelProductIDS[i]] {
			CProducts = append(CProducts, vbox.OrgProduct{ChannelProductID: orgProduct.ChannelProductIDS[i], OrganizationID: orgProduct.OrganizationID})
		}
	}
	err = global.GVA_DB.Create(&CProducts).Error
	return err
}

func (opService *OrgProductService) DeleteOrgProduct(ids []uint, orgID uint) (err error) {
	return global.GVA_DB.Where("channel_product_id in (?) and organization_id = ?", ids, orgID).Delete(&[]vbox.OrgProduct{}).Error
}

func (opService *OrgProductService) FindOrgProductAll(orgID string) ([]uint, error) {
	var Products []vbox.OrgProduct
	var ids []uint
	err := global.GVA_DB.Find(&Products, "organization_id = ?", orgID).Error
	if err != nil {
		return ids, err
	}
	for i := range Products {
		ids = append(ids, Products[i].ChannelProductID)
	}
	return ids, err
}

// GetOrgProductList 分页获取当前组织下可用产品
func (opService *OrgProductService) GetOrgProductList(info vboxReq.OrgProductSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Table("organization").Select("p.id cp_id,p.product_name,p.channel_code,organization.id oid,organization.name").
		Joins("inner join org_product op on op.organization_id = organization.id").
		Joins("inner join vbox_channel_product p on op.channel_product_id = p.id")
	var orgs []vbox.OrgProductRes
	// 如果有条件搜索 下方会自动创建搜索语句
	db = db.Where("op.organization_id = ?", info.OrganizationID)
	if info.ProductName != "" {
		db = db.Where("p.product_name LIKE ?", "%"+info.ProductName+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Debug().Limit(limit).Offset(offset).Scan(&orgs).Error
	return orgs, total, err
}
