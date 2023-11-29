package vbox

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"gorm.io/gorm"
	"sort"
	"time"
)

type ChannelShopService struct {
}

// CreateChannelShop 创建引导商铺记录
func (channelShopService *ChannelShopService) CreateChannelShop(channelShop *vboxReq.ChannelShop) (err error) {
	length := len(channelShop.ChannelShopList)
	if length < 1 {
		err = fmt.Errorf("创建的店铺地址至少包含1个")
		return err
	}
	pid := channelShop.ProductId
	if pid == "" {
		pid = time.Now().Format("20060102150405") + rand_string.RandomInt(3)
	}
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		for _, c := range channelShop.ChannelShopList {
			if c.ID == 0 { //走创建
				chNew := &vbox.ChannelShop{
					Cid:        channelShop.Cid,
					ProductId:  pid,
					ShopRemark: channelShop.ShopRemark,
					Address:    c.Address,
					Money:      c.Money,
					Status:     c.Status,
					CreatedBy:  channelShop.CreatedBy,
				}
				if err := tx.Model(&vbox.ChannelShop{}).Create(&chNew).Error; err != nil {
					return err
				}
			} else { //走更新
				chNew := &vbox.ChannelShop{
					Cid:        channelShop.Cid,
					ProductId:  pid,
					ShopRemark: channelShop.ShopRemark,
					Address:    c.Address,
					Money:      c.Money,
					Status:     c.Status,
					UpdatedBy:  channelShop.CreatedBy,
				}
				if err := tx.Model(&vbox.ChannelShop{}).Where("id = ?", c.ID).Updates(&chNew).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	return err
}

// DeleteChannelShop 删除引导商铺记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelShopService *ChannelShopService) DeleteChannelShop(channelShop vbox.ChannelShop) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.ID).Update("deleted_by", channelShop.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&channelShop).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteChannelShopByIds 批量删除引导商铺记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelShopService *ChannelShopService) DeleteChannelShopByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.ChannelShop{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.ChannelShop{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannelShop 更新引导商铺记录（ type: 1-更新店名 2-开关单条 3-开关整个店 ）
func (channelShopService *ChannelShopService) UpdateChannelShop(channelShop vboxReq.ChannelShopReq) (err error) {
	// 1-更新店名 2-开关单条 3-开关整个店
	if channelShop.Type == 1 {
		err = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("product_id = ?", channelShop.ProductId).
			Update("shop_remark", channelShop.ShopRemark).
			Update("updated_by", channelShop.UpdatedBy).Error
	} else if channelShop.Type == 2 {
		err = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.Id).
			Update("status", channelShop.Status).
			Update("updated_by", channelShop.UpdatedBy).Error

	} else if channelShop.Type == 3 {
		err = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("product_id = ?", channelShop.ProductId).
			Update("status", channelShop.Status).
			Update("updated_by", channelShop.UpdatedBy).Error
	} else {
		err = fmt.Errorf("不支持的操作，type检查")
	}

	return err
}

// GetChannelShop 根据id获取引导商铺记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelShopService *ChannelShopService) GetChannelShop(id uint) (channelShop vbox.ChannelShop, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&channelShop).Error
	return
}

// GetChannelShopByProductId 根据ProductId获取引导商铺记录
func (channelShopService *ChannelShopService) GetChannelShopByProductId(productId string) (res vboxReq.ChannelShop, err error) {

	var csDB []vbox.ChannelShop
	err = global.GVA_DB.Where("product_id = ?", productId).Find(&csDB).Error

	if len(csDB) == 0 { //无记录
		return vboxReq.ChannelShop{}, nil
	}

	res = vboxReq.ChannelShop{
		ProductId:       productId,
		ShopRemark:      csDB[0].ShopRemark,
		Cid:             csDB[0].Cid,
		ChannelShopList: []vboxReq.ChannelShopSub{}}

	for _, record := range csDB {
		res.ChannelShopList = append(res.ChannelShopList, vboxReq.ChannelShopSub{
			ID:      record.ID,
			Address: record.Address,
			Money:   record.Money,
			Status:  record.Status,
		})
	}

	return res, nil
}

// GetChannelShopInfoList 分页获取引导商铺记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelShopService *ChannelShopService) GetChannelShopInfoList(info vboxReq.ChannelShopSearch, ids []uint) (res []vboxReq.ChannelShop, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelShop{})
	var channelShops []vbox.ChannelShop

	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.ShopRemark != "" {
		db = db.Where("shop_remark LIKE ?", "%"+info.ShopRemark+"%")
	}

	if info.ProductId != "" {
		db = db.Where("product_id = ?", info.ProductId)
	}

	err = db.Where("created_by in (?)", ids).Find(&channelShops).Error

	shopMap := make(map[string][]vbox.ChannelShop)

	for _, shop := range channelShops {
		if _, ok := shopMap[shop.ProductId]; !ok { // 如果map中没有productID这个key，则添加一个空切片
			shopMap[shop.ProductId] = []vbox.ChannelShop{}
		}
		shopMap[shop.ProductId] = append(shopMap[shop.ProductId], shop)
	}

	for productID, v := range shopMap {
		if len(v) == 0 {
			continue
		}
		csNew := vboxReq.ChannelShop{
			ProductId:       productID,
			ShopRemark:      v[0].ShopRemark,
			Cid:             v[0].Cid,
			ChannelShopList: []vboxReq.ChannelShopSub{}}

		for _, record := range v {
			csNew.ChannelShopList = append(csNew.ChannelShopList, vboxReq.ChannelShopSub{
				ID:      record.ID,
				Address: record.Address,
				Money:   record.Money,
				Status:  record.Status,
			})
		}

		res = append(res, csNew)
	}

	sort.Sort(vboxReq.ChannelShopList(res))

	return res, err
}
