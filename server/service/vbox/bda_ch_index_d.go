package vbox

import (
	"fmt"
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

type BdaChIndexDService struct {
}

// CreateBdaChIndexD 创建用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChIndexDService *BdaChIndexDService) CreateBdaChIndexD(bdaChIndexD *vbox.BdaChIndexD) (err error) {
	err = global.GVA_DB.Create(bdaChIndexD).Error
	return err
}

// DeleteBdaChIndexD 删除用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChIndexDService *BdaChIndexDService) DeleteBdaChIndexD(bdaChIndexD vbox.BdaChIndexD) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChIndexD{}).Where("id = ?", bdaChIndexD.ID).Update("deleted_by", bdaChIndexD.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&bdaChIndexD).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteBdaChIndexDByIds 批量删除用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChIndexDService *BdaChIndexDService) DeleteBdaChIndexDByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChIndexD{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.BdaChIndexD{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateBdaChIndexD 更新用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChIndexDService *BdaChIndexDService) UpdateBdaChIndexD(bdaChIndexD vbox.BdaChIndexD) (err error) {
	err = global.GVA_DB.Save(&bdaChIndexD).Error
	return err
}

// GetBdaChIndexD 根据id获取用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChIndexDService *BdaChIndexDService) GetBdaChIndexD(id uint) (bdaChIndexD vbox.BdaChIndexD, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&bdaChIndexD).Error
	return
}

// GetBdaChIndexDInfoList 分页获取用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChIndexDService *BdaChIndexDService) GetBdaChIndexDInfoList(info vboxReq.BdaChIndexDSearch) (list []vbox.BdaChIndexD, total int64, err error) {
	fmt.Println("统计开始")
	bdaChIndexDService.CronVboxBdaChIndexD()
	fmt.Println("统计结束")

	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.BdaChIndexD{})
	var bdaChIndexDs []vbox.BdaChIndexD
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Uid != nil {
		db = db.Where("uid = ?", info.Uid)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&bdaChIndexDs).Error

	return bdaChIndexDs, total, err
}

// Author yoga
func (bdaChIndexDService *BdaChIndexDService) CronVboxBdaChIndexDByHand(dt string) (err error) {

	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ? ", dt)
	var uids []int
	err = db.Select("distinct created_by").Pluck("created_by", &uids).Error
	fmt.Println("uids=", uids)
	if err != nil {
		return
	}
	var chIds []string
	err = db.Select("distinct channel_code").Pluck("channel_code", &chIds).Error
	if err != nil {
		return
	}
	for _, uid := range uids {
		for _, chId := range chIds {
			//yInfoList, yOrderTotal, err := GetUsersVboxPayOrderInfoList(uid, chId, 0)
			yInfoList, yOrderTotal, err := GetUsersVboxPayOrderInfoListByHand(uid, chId, dt)

			fmt.Println(yInfoList)
			fmt.Println(yOrderTotal)
			if err != nil {
				return err
			}
			yGroupedCounts := make(map[string]int16)
			yOkGroupedCounts := make(map[string]int16)
			yOkGroupedCosts := make(map[string]int)

			for _, order := range yInfoList {
				uid := strconv.Itoa(int(order.CreatedBy)) + "-" + order.ChannelCode
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
			key := strconv.Itoa(uid) + "-" + chId
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

			chCode := vcp.ChannelCode

			entity := vbox.BdaChIndexD{
				Uid:             &uid,
				Username:        userInfo.Username,
				ChannelCode:     chCode,
				ProductId:       chId,
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
			log.Println("统计 ch 结果=", entity)
			err = global.GVA_DB.Save(&entity).Error
		}
	}

	return err
}

// Author yoga
func (bdaChIndexDService *BdaChIndexDService) CronVboxBdaChIndexD() (err error) {

	dt := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ? ", dt)
	var uids []int
	err = db.Select("distinct created_by").Pluck("created_by", &uids).Error
	fmt.Println("uids=", uids)
	if err != nil {
		return
	}
	var chIds []string
	err = db.Select("distinct channel_code").Pluck("channel_code", &chIds).Error
	if err != nil {
		return
	}
	for _, uid := range uids {
		for _, chId := range chIds {
			//yInfoList, yOrderTotal, err := GetUsersVboxPayOrderInfoList(uid, chId, 0)
			yInfoList, yOrderTotal, err := GetUsersVboxPayOrderInfoListByHand(uid, chId, dt)

			fmt.Println(yInfoList)
			fmt.Println(yOrderTotal)
			if err != nil {
				return err
			}
			yGroupedCounts := make(map[string]int16)
			yOkGroupedCounts := make(map[string]int16)
			yOkGroupedCosts := make(map[string]int)

			for _, order := range yInfoList {
				uid := strconv.Itoa(int(order.CreatedBy)) + "-" + order.ChannelCode
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
			key := strconv.Itoa(uid) + "-" + chId
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

			chCode := vcp.ChannelCode

			entity := vbox.BdaChIndexD{
				Uid:             &uid,
				Username:        userInfo.Username,
				ChannelCode:     chCode,
				ProductId:       chId,
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
			log.Println("统计 ch 结果=", entity)
			err = global.GVA_DB.Save(&entity).Error
		}
	}

	return err
}

func GetUsersVboxPayOrderInfoList(id int, chId string, num int) (list []vbox.PayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var vpos []vbox.PayOrder
	fmt.Println("id=", id, "chid=", chId, "num=", num)
	err = db.Where("created_by = ? and channel_code = ? and DATE(created_at) = (CURDATE() - INTERVAL ? DAY)", id, chId, num).Find(&vpos).Error
	total = int64(len(vpos))
	return vpos, total, err
}

func GetUsersVboxPayOrderInfoListByHand(id int, chId string, dt string) (list []vbox.PayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var vpos []vbox.PayOrder
	fmt.Println("id=", id, "chid=", chId, "dt=", dt)
	err = db.Where("created_by = ? and channel_code = ? and DATE(created_at) = ?", id, chId, dt).Find(&vpos).Error
	total = int64(len(vpos))
	return vpos, total, err
}
