package vbox

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxRep "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"gorm.io/gorm"
	"time"
)

type VboxPayOrderService struct {
}

//var userService = service.ServiceGroupApp.SystemServiceGroup.UserService

// CreateVboxPayOrder 创建VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) CreateVboxPayOrder(vpo *vbox.VboxPayOrder) (err error) {
	err = global.GVA_DB.Create(vpo).Error
	return err
}

// DeleteVboxPayOrder 删除VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) DeleteVboxPayOrder(vpo vbox.VboxPayOrder) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxPayOrder{}).Where("id = ?", vpo.ID).Update("deleted_by", vpo.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&vpo).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteVboxPayOrderByIds 批量删除VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) DeleteVboxPayOrderByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxPayOrder{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxPayOrder{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateVboxPayOrder 更新VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) UpdateVboxPayOrder(vpo vbox.VboxPayOrder) (err error) {
	err = global.GVA_DB.Save(&vpo).Error
	return err
}

// GetVboxPayOrder 根据id获取VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) GetVboxPayOrder(id uint) (vpo vbox.VboxPayOrder, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vpo).Error
	return
}

// GetVboxPayOrderInfoList 分页获取VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) GetVboxPayOrderInfoList(info vboxReq.VboxPayOrderSearch, ids []int) (list []vbox.VboxPayOrder, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxPayOrder{})
	var vpos []vbox.VboxPayOrder

	db = db.Where("uid in (?)", ids).Find(&vpos)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.OrderId != "" {
		db = db.Where("order_id = ?", info.OrderId)
	}
	if info.PAccount != "" {
		db = db.Where("p_account = ?", info.PAccount)
	}
	if info.Cost != nil {
		db = db.Where("cost = ?", info.Cost)
	}
	if info.Uid != nil {
		db = db.Where("uid = ?", info.Uid)
	}
	if info.AcId != "" {
		db = db.Where("ac_id = ?", info.AcId)
	}
	if info.CChannelId != "" {
		db = db.Where("c_channel_id = ?", info.CChannelId)
	}
	if info.PlatformOid != "" {
		db = db.Where("platform_oid = ?", info.PlatformOid)
	}
	if info.PayIp != "" {
		db = db.Where("pay_ip = ?", info.PayIp)
	}
	if info.PayRegion != "" {
		db = db.Where("pay_region = ?", info.PayRegion)
	}
	if info.ResourceUrl != "" {
		db = db.Where("resource_url = ?", info.ResourceUrl)
	}
	if info.NotifyUrl != "" {
		db = db.Where("notify_url = ?", info.NotifyUrl)
	}
	if info.OrderStatus != nil {
		db = db.Where("order_status = ?", info.OrderStatus)
	}
	if info.CallbackStatus != nil {
		db = db.Where("callback_status = ?", info.CallbackStatus)
	}
	if info.CodeUseStatus != nil {
		db = db.Where("code_use_status = ?", info.CodeUseStatus)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&vpos).Error
	return vpos, total, err
}

// GetUsersVboxPayOrderInfoList 分页获取VboxPayOrder记录
// Author youga
func (vpoService *VboxPayOrderService) GetUsersVboxPayOrderInfoList(ids []int, num int) (list []vbox.VboxPayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxPayOrder{})
	var vpos []vbox.VboxPayOrder
	err = db.Where("uid in (?) and DATE(create_time) = (CURDATE() - INTERVAL ? DAY)", ids, num).Find(&vpos).Error
	total = int64(len(vpos))
	return vpos, total, err
}

func (vpoService *VboxPayOrderService) GetVboxUserPayOrderAnalysis(id uint, idList []int) (list []vboxRep.VboxUserOrderPayAnalysis, total int64, err error) {
	tInfoList, tOrderTotal, err := vpoService.GetUsersVboxPayOrderInfoList(idList, 0)
	if err != nil {
		return
	}
	yInfoList, yOrderTotal, err := vpoService.GetUsersVboxPayOrderInfoList(idList, 1)
	if err != nil {
		return
	}

	tGroupedCounts := make(map[int]int16)
	tOkGroupedCounts := make(map[int]int16)
	tOkGroupedCosts := make(map[int]int)

	yGroupedCounts := make(map[int]int16)
	yOkGroupedCounts := make(map[int]int16)
	yOkGroupedCosts := make(map[int]int)

	for _, order := range tInfoList {
		uid := *order.Uid
		tGroupedCounts[uid]++
		if *order.OrderStatus == 1 {
			tOkGroupedCounts[uid]++
			tOkGroupedCosts[uid] += *order.Cost
		}
	}

	for _, order := range yInfoList {
		uid := *order.Uid
		yGroupedCounts[uid]++
		if *order.OrderStatus == 1 {
			yOkGroupedCounts[uid]++
			yOkGroupedCosts[uid] += *order.Cost
		}
	}

	fmt.Println("Uid分组统计OrderStatus=1的个数和Cost总和：")

	for _, uid := range idList {
		tOrderQuantify := tOrderTotal
		tOkOrderQuantify := 0
		tOkRate := 0
		tInCome := 0
		yOrderQuantify := yOrderTotal
		yOkOrderQuantify := 0
		yOkRate := 0
		yInCome := 0
		// 判断 tGroupedCounts 中是否包含指定的 uid 键
		_, tContainsUID := tGroupedCounts[uid]
		_, tOkContainsUID := tOkGroupedCounts[uid]
		_, yContainsUID := yGroupedCounts[uid]
		_, yOkContainsUID := yOkGroupedCounts[uid]

		if tContainsUID {
			tOrderQuantify = int64(tGroupedCounts[uid])
			if tOkContainsUID {
				tOkOrderQuantify = int(tOkGroupedCounts[uid])
			}

			if tOrderQuantify > 0 {
				result := float64(tOkOrderQuantify) / float64(tOrderQuantify)
				tOkRate = int(result * 100)
				tInCome = tOkGroupedCosts[uid]
			}

		}

		if yContainsUID {

			yOrderQuantify = int64(yGroupedCounts[uid])
			if yOkContainsUID {
				yOkOrderQuantify = int(yOkGroupedCounts[uid])
			}

			if yOrderQuantify > 0 {
				result := float64(yOkOrderQuantify) / float64(yOrderQuantify)
				yOkRate = int(result * 100)
				yInCome = yOkGroupedCosts[uid]
			}

		}

		var userInfo system.SysUser
		err = global.GVA_DB.Where("`id` = ?", uid).First(&userInfo).Error
		if err != nil {
			return
		}

		entity := vboxRep.VboxUserOrderPayAnalysis{
			Uid:              &uid,
			Username:         userInfo.Username,
			Balance:          99999, //todo
			ChIdCnt:          2,     //todo
			OpenChId:         1,     //todo
			YOrderQuantify:   int(yOrderQuantify),
			YOkOrderQuantify: yOkOrderQuantify,
			YOkRate:          yOkRate,
			YInCome:          yInCome,
			TOrderQuantify:   int(tOrderQuantify),
			TOkOrderQuantify: tOkOrderQuantify,
			TOkRate:          tOkRate,
			TInCome:          tInCome,
		}
		list = append(list, entity)

	}
	return list, int64(len(idList)), err

}

// 判断是否是今天
func isToday(createTime *time.Time, now time.Time) bool {
	year, month, day := now.Date()
	createYear, createMonth, createDay := createTime.Date()

	return year == createYear && month == createMonth && day == createDay
}

// 判断是否是昨天
func isYesterday(createTime *time.Time, now time.Time) bool {
	yesterday := now.AddDate(0, 0, -1)
	year, month, day := yesterday.Date()
	createYear, createMonth, createDay := createTime.Date()

	return year == createYear && month == createMonth && day == createDay
}
