package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"gorm.io/gorm"
	"log"
	"strconv"
	"time"
)

type BdaChShopIndexDService struct {
}

// CreateBdaChShopIndexD 创建用户通道店铺成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChshopIndexDService *BdaChShopIndexDService) CreateBdaChShopIndexD(bdaChshopIndexD *vbox.BdaChShopIndexD) (err error) {
	err = global.GVA_DB.Create(bdaChshopIndexD).Error
	return err
}

// DeleteBdaChShopIndexD 删除用户通道店铺成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChshopIndexDService *BdaChShopIndexDService) DeleteBdaChShopIndexD(bdaChshopIndexD vbox.BdaChShopIndexD) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChShopIndexD{}).Where("id = ?", bdaChshopIndexD.ID).Update("deleted_by", bdaChshopIndexD.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&bdaChshopIndexD).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteBdaChShopIndexDByIds 批量删除用户通道店铺成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChshopIndexDService *BdaChShopIndexDService) DeleteBdaChShopIndexDByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChShopIndexD{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.BdaChShopIndexD{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateBdaChShopIndexD 更新用户通道店铺成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChshopIndexDService *BdaChShopIndexDService) UpdateBdaChShopIndexD(bdaChshopIndexD vbox.BdaChShopIndexD) (err error) {
	err = global.GVA_DB.Save(&bdaChshopIndexD).Error
	return err
}

// GetBdaChShopIndexD 根据id获取用户通道店铺成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChshopIndexDService *BdaChShopIndexDService) GetBdaChShopIndexD(id uint) (bdaChshopIndexD vbox.BdaChShopIndexD, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&bdaChshopIndexD).Error
	return
}

// GetBdaChShopIndexDInfoList 分页获取用户通道店铺成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChshopIndexDService *BdaChShopIndexDService) GetBdaChShopIndexDInfoList(info vboxReq.BdaChShopIndexDSearch) (list []vbox.BdaChShopIndexD, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.BdaChShopIndexD{})
	var bdaChshopIndexDs []vbox.BdaChShopIndexD
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Uid != nil {
		db = db.Where("uid = ?", info.Uid)
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.ShopRemark != "" {
		db = db.Where("shop_remark = ?", info.ShopRemark)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&bdaChshopIndexDs).Error
	return bdaChshopIndexDs, total, err
}

func (bdaChshopIndexDService *BdaChShopIndexDService) CronVboxBdaChShopIndexD() (err error) {
	dt := time.Now().Format("2006-01-02")

	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ? ", dt)
	var uids []int
	err = db.Select("distinct uid").Pluck("uid", &uids).Error
	if err != nil {
		return
	}
	var chIds []string
	err = db.Select("distinct channel_code").Pluck("channel_code", &chIds).Error
	if err != nil {
		return
	}

	var shops []string
	err = db.Select("distinct resource_url").Pluck("resource_url", &shops).Error
	if err != nil {
		return
	}

	for _, uid := range uids {
		for _, chId := range chIds {
			for _, shop := range shops {

				yInfoList, yOrderTotal, err := GetUsersShopChVboxPayOrderInfoList(uid, chId, shop, 1)
				if err != nil {
					return err
				}
				yGroupedCounts := make(map[string]int16)
				yOkGroupedCounts := make(map[string]int16)
				yOkGroupedCosts := make(map[string]int)

				for _, order := range yInfoList {
					uid := strconv.Itoa(order.Uid) + "-" + order.ChannelCode + "-" + order.ResourceUrl
					yGroupedCounts[uid]++
					if order.OrderStatus == 1 {
						yOkGroupedCounts[uid]++
						yOkGroupedCosts[uid] += order.Money
					}
				}

				yOrderQuantify := yOrderTotal
				yOkOrderQuantify := 0
				yOkRate := 0
				yInCome := 0
				// 判断 tGroupedCounts 中是否包含指定的 uid 键
				key := strconv.Itoa(uid) + "-" + chId + "-" + shop
				_, yContainsUID := yGroupedCounts[key]
				_, yOkContainsUID := yOkGroupedCounts[key]

				if yContainsUID {

					yOrderQuantify = int64(yGroupedCounts[key])
					if yOkContainsUID {
						yOkOrderQuantify = int(yOkGroupedCounts[key])
					}

					if yOrderQuantify > 0 {
						result := float64(yOkOrderQuantify) / float64(yOrderQuantify)
						yOkRate = int(result * 100)
						yInCome = yOkGroupedCosts[key]
					}

				}

				var userInfo system.SysUser
				err = global.GVA_DB.Where("`id` = ?", uid).First(&userInfo).Error
				if err != nil {
					return err
				}
				var vcp vbox.ChannelProduct
				err = global.GVA_DB.Where("channel_code = ?", chId).First(&vcp).Error
				if err != nil {
					return err
				}
				var vcs vbox.ChannelShop
				err = global.GVA_DB.Where("cid = ? and shop_remark = ?", chId, shop).First(&vcs).Error
				if err != nil {
					return err
				}
				chCode := vcp.ChannelCode

				entity := vbox.BdaChShopIndexD{
					Uid:             &uid,
					UserName:        userInfo.Username,
					Cid:             chId,
					ShopRemark:      vcs.ShopRemark,
					ChannelCode:     chCode,
					ProductId:       vcp.ProductId,
					ProductName:     vcp.ProductName,
					OrderQuantify:   int(yOrderQuantify),
					OkOrderQuantify: yOkOrderQuantify,
					Ratio:           float64(yOkRate),
					Income:          yInCome,
					Dt:              dt,
					CreatedBy:       uint(uid),
					UpdatedBy:       1,
					DeletedBy:       1,
				}
				//fmt.Println(entity.Dt)
				log.Println("统计 ch shop 结果=", entity)
				err = global.GVA_DB.Save(&entity).Error
			}
		}
	}
	return err
}

func GetUsersShopChVboxPayOrderInfoList(id int, chId string, shop string, num int) (list []vbox.PayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var vpos []vbox.PayOrder
	err = db.Where("uid = ? and resource_url = ? and channel_code = ? and DATE(created_at) = (CURDATE() - INTERVAL ? DAY)", id, shop, chId, num).Find(&vpos).Error
	total = int64(len(vpos))
	return vpos, total, err
}
