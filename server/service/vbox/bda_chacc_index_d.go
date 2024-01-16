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

type BdaChaccIndexDService struct {
}

// CreateBdaChaccIndexD 创建用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccIndexDService *BdaChaccIndexDService) CreateBdaChaccIndexD(bdaChaccIndexD *vbox.BdaChaccIndexD) (err error) {
	err = global.GVA_DB.Create(bdaChaccIndexD).Error
	return err
}

// DeleteBdaChaccIndexD 删除用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccIndexDService *BdaChaccIndexDService) DeleteBdaChaccIndexD(bdaChaccIndexD vbox.BdaChaccIndexD) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChaccIndexD{}).Where("id = ?", bdaChaccIndexD.ID).Update("deleted_by", bdaChaccIndexD.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&bdaChaccIndexD).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteBdaChaccIndexDByIds 批量删除用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccIndexDService *BdaChaccIndexDService) DeleteBdaChaccIndexDByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChaccIndexD{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.BdaChaccIndexD{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateBdaChaccIndexD 更新用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccIndexDService *BdaChaccIndexDService) UpdateBdaChaccIndexD(bdaChaccIndexD vbox.BdaChaccIndexD) (err error) {
	err = global.GVA_DB.Save(&bdaChaccIndexD).Error
	return err
}

// GetBdaChaccIndexD 根据id获取用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccIndexDService *BdaChaccIndexDService) GetBdaChaccIndexD(id uint) (bdaChaccIndexD vbox.BdaChaccIndexD, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&bdaChaccIndexD).Error
	return
}

// GetBdaChaccIndexDInfoList 分页获取用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccIndexDService *BdaChaccIndexDService) GetBdaChaccIndexDInfoList(info vboxReq.BdaChaccIndexDSearch) (list []vbox.BdaChaccIndexD, total int64, err error) {
	fmt.Println("统计开始")
	bdaChaccIndexDService.CronVboxBdaChaccIndexD()
	fmt.Println("统计结束")
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.BdaChaccIndexD{})
	var bdaChaccIndexDs []vbox.BdaChaccIndexD
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Uid != nil {
		db = db.Where("uid = ?", info.Uid)
	}
	if info.Username != "" {
		db = db.Where("user_name = ?", info.Username)
	}
	if info.AcId != "" {
		db = db.Where("ac_id = ?", info.AcId)
	}
	if info.AcRemark != "" {
		db = db.Where("ac_remark = ?", info.AcRemark)
	}
	if info.ChannelCode != "" {
		db = db.Where("channel_code = ?", info.ChannelCode)
	}
	if info.ProductId != "" {
		db = db.Where("product_id = ?", info.ProductId)
	}
	if info.ProductName != "" {
		db = db.Where("product_name = ?", info.ProductName)
	}
	if info.Dt != "" {
		db = db.Where("dt = ?", info.Dt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&bdaChaccIndexDs).Error
	return bdaChaccIndexDs, total, err
}

func (bdaChaccIndexDService *BdaChaccIndexDService) CronVboxBdaChaccIndexDByHand(dt string) (err error) {
	//dt := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ? ", dt)
	var uids []int
	err = db.Select("distinct created_by").Pluck("created_by", &uids).Error
	if err != nil {
		return
	}
	var chIds []string
	err = db.Select("distinct channel_code").Pluck("channel_code", &chIds).Error
	if err != nil {
		return
	}

	var accIds []string
	err = db.Select("distinct ac_id").Pluck("ac_id", &chIds).Error
	if err != nil {
		return
	}

	for _, uid := range uids {
		for _, chId := range chIds {
			for _, accId := range accIds {

				yInfoList, yOrderTotal, err := GetUsersAccChVboxPayOrderInfoList(uid, chId, accId, dt)
				if err != nil {
					return err
				}
				yGroupedCounts := make(map[string]int16)
				yOkGroupedCounts := make(map[string]int16)
				yOkGroupedCosts := make(map[string]int)

				for _, order := range yInfoList {
					uid := strconv.Itoa(int(order.CreatedBy)) + "-" + order.ChannelCode + "-" + order.AcId
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
				key := strconv.Itoa(uid) + "-" + chId + "-" + accId
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
				var vca vbox.ChannelAccount
				err = global.GVA_DB.Where("ac_id = ?", accId).First(&vca).Error
				if err != nil {
					return err
				}
				chCode := vcp.ChannelCode

				entity := vbox.BdaChaccIndexD{
					Uid:             &uid,
					Username:        userInfo.Username,
					AcId:            accId,
					AcAccount:       vca.AcAccount,
					AcRemark:        vca.AcRemark,
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
				log.Println("统计 ch acc 结果=", entity)
				err = global.GVA_DB.Save(&entity).Error
			}
		}
	}
	return err
}

func (bdaChaccIndexDService *BdaChaccIndexDService) CronVboxBdaChaccIndexD() (err error) {
	dt := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	fmt.Println(dt, "统计开始")
	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ? ", dt)
	var uids []int
	err = db.Select("distinct created_by").Pluck("created_by", &uids).Error
	fmt.Println("uids=", uids)
	if err != nil {
		return
	}
	var chIds []string
	err = db.Select("distinct channel_code").Pluck("channel_code", &chIds).Error
	fmt.Println("chIds=", chIds)
	if err != nil {
		return
	}

	var accIds []string
	err = db.Select("distinct ac_id").Pluck("ac_id", &accIds).Error
	fmt.Println("accIds=", accIds)
	if err != nil {
		return
	}

	for _, uid := range uids {
		for _, chId := range chIds {
			for _, accId := range accIds {

				yInfoList, yOrderTotal, err := GetUsersAccChVboxPayOrderInfoList(uid, chId, accId, dt)
				//fmt.Println("total", yOrderTotal)
				if err != nil {
					return err
				}
				if yOrderTotal > 0 {

					yGroupedCounts := make(map[string]int16)
					yOkGroupedCounts := make(map[string]int16)
					yOkGroupedCosts := make(map[string]int)

					for _, order := range yInfoList {
						uid := strconv.Itoa(int(order.CreatedBy)) + "-" + order.ChannelCode + "-" + order.AcId
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
					key := strconv.Itoa(uid) + "-" + chId + "-" + accId
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
					var vca vbox.ChannelAccount
					fmt.Println("accId=", accId, (accId == ""))
					account := ""
					acRemark := ""
					if accId != "" {
						err = global.GVA_DB.Where("ac_id is not null and ac_id != '' and ac_id = ?", accId).First(&vca).Error
						if err != nil {
							return err
						}
						account = vca.AcAccount
						acRemark = vca.AcRemark
					}

					chCode := vcp.ChannelCode

					entity := vbox.BdaChaccIndexD{
						Uid:             &uid,
						Username:        userInfo.Username,
						AcId:            accId,
						AcAccount:       account,
						AcRemark:        acRemark,
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
					log.Println("统计 ch acc 结果=", entity)
					err = global.GVA_DB.Save(&entity).Error
				}

			}
		}
	}
	return err
}

func GetUsersAccChVboxPayOrderInfoList(id int, chId string, accId string, dt string) (list []vbox.PayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var vpos []vbox.PayOrder
	//err = db.Where("uid = ? and ac_id = ? and channel_code = ? and DATE(created_at) = (CURDATE() - INTERVAL ? DAY)", id, accId, chId, num).Find(&vpos).Error
	err = db.Where("created_by = ? and ac_id = ? and channel_code = ? and DATE(created_at) = ?", id, accId, chId, dt).Find(&vpos).Error
	total = int64(len(vpos))
	return vpos, total, err
}
