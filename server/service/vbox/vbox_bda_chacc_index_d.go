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
func (bdaChaccDService *VboxBdaChaccIndexDService) DeleteVboxBdaChaccIndexD(bdaChaccD vbox.VboxBdaChaccIndexD) (err error) {
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
func (bdaChaccDService *VboxBdaChaccIndexDService) DeleteVboxBdaChaccIndexDByIds(ids request.IdsReq, deleted_by uint) (err error) {
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
func (bdaChaccDService *VboxBdaChaccIndexDService) UpdateVboxBdaChaccIndexD(bdaChaccD vbox.VboxBdaChaccIndexD) (err error) {
	err = global.GVA_DB.Save(&bdaChaccD).Error
	return err
}

// GetVboxBdaChaccIndexD 根据id获取VboxBdaChaccIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccDService *VboxBdaChaccIndexDService) GetVboxBdaChaccIndexD(id uint) (bdaChaccD vbox.VboxBdaChaccIndexD, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&bdaChaccD).Error
	return
}

// GetVboxBdaChaccIndexDInfoList 分页获取VboxBdaChaccIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccDService *VboxBdaChaccIndexDService) GetVboxBdaChaccIndexDInfoList(info vboxReq.VboxBdaChaccIndexDSearch) (list []vbox.VboxBdaChaccIndexD, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxBdaChaccIndexD{})
	var bdaChaccDs []vbox.VboxBdaChaccIndexD
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Uid != nil {
		db = db.Where("uid = ?", info.Uid)
	}
	if info.UserName != "" {
		db = db.Where("username = ?", info.UserName)
	}
	if info.PAccount != "" {
		db = db.Where("p_account = ?", info.PAccount)
	}
	if info.ChannelCode != nil {
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

	err = db.Limit(limit).Offset(offset).Find(&bdaChaccDs).Error
	return bdaChaccDs, total, err
}

// Author yoga
func (bdaChaccDService *VboxBdaChaccIndexDService) CronVboxBdaChaccIndexD() (err error) {
	queryB := `
		    SELECT distinct c_channel_id  as chId
			FROM vbox_pay_order
			WHERE  uid = ? and DATE(create_time) = (CURDATE() - INTERVAL ? DAY)
		;
    `

	queryC := `
		    SELECT distinct p_account  as accId
			FROM vbox_pay_order
			WHERE  uid = ? and c_channel_id = ? and DATE(create_time) = (CURDATE() - INTERVAL ? DAY)
		;
    `

	dt := time.Now().Format("2006-01-02")
	//fmt.Println(dt)
	//var bdaChDs []vbox.VboxBdaChIndexD
	userList, tot, err := GetOwnerUserIdsList(1)
	if err != nil || tot == 0 {
		return
	}
	var idList []int
	for _, user := range userList {
		idList = append(idList, int(user.ID))
	}
	fmt.Println(idList)
	for _, uid := range idList {

		rowDt, errA := global.GVA_DB.Raw(queryB, uid, 1).Rows()
		var chIds []string
		if errA != nil {
			return err
		}
		for rowDt.Next() {
			var chId string
			scanErr := rowDt.Scan(&chId)
			if scanErr != nil {
				return err
			}
			chIds = append(chIds, chId)
		}
		//fmt.Println("chIds=", chIds)
		for _, chId := range chIds {

			rowDt, errA := global.GVA_DB.Raw(queryC, uid, chId, 1).Rows()
			var accIds []string
			if errA != nil {
				return err
			}
			for rowDt.Next() {
				var accId string
				scanErr := rowDt.Scan(&accId)
				if scanErr != nil {
					return err
				}
				accIds = append(accIds, accId)
			}
			fmt.Println("chIds=", accIds)
			for _, accId := range accIds {
				yInfoList, yOrderTotal, err := GetUsersAccChVboxPayOrderInfoList(uid, chId, accId, 1)
				if err != nil {
					return err
				}
				yGroupedCounts := make(map[string]int16)
				yOkGroupedCounts := make(map[string]int16)
				yOkGroupedCosts := make(map[string]int)

				for _, order := range yInfoList {
					uid := strconv.Itoa(order.Uid) + "-" + order.CChannelId + "-" + order.PAccount
					yGroupedCounts[uid]++
					if order.OrderStatus == 1 {
						yOkGroupedCounts[uid]++
						yOkGroupedCosts[uid] += order.Cost
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
				err = global.GVA_DB.Where("product_id = ?", chId).First(&vcp).Error
				if err != nil {
					return err
				}

				chCode := int(vcp.ChannelCode)

				entity := vbox.VboxBdaChaccIndexD{
					Uid:             &uid,
					UserName:        userInfo.Username,
					PAccount:        accId,
					ChannelCode:     &chCode,
					ProductId:       chId,
					ProductName:     vcp.ProductName,
					OrderQuantify:   int(yOrderQuantify),
					OkOrderQuantify: yOkOrderQuantify,
					Ratio:           float64(yOkRate),
					Income:          yInCome,
					Dt:              dt,
					CreatedBy:       1,
					UpdatedBy:       1,
					DeletedBy:       1,
				}
				//fmt.Println(entity.Dt)
				log.Println("统计 acc 结果=", entity)
				err = global.GVA_DB.Save(&entity).Error
			}
		}

	}

	return err
}

func GetUsersAccChVboxPayOrderInfoList(id int, chId string, accId string, num int) (list []vbox.VboxPayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxPayOrder{})
	var vpos []vbox.VboxPayOrder
	err = db.Where("uid = ? and p_account = ? and c_channel_id = ? and DATE(create_time) = (CURDATE() - INTERVAL ? DAY)", id, accId, chId, num).Find(&vpos).Error
	total = int64(len(vpos))
	return vpos, total, err
}
