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

type VboxBdaChIndexDService struct {
}

// CreateVboxBdaChIndexD 创建VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService) CreateVboxBdaChIndexD(bdaChD *vbox.VboxBdaChIndexD) (err error) {
	err = global.GVA_DB.Create(bdaChD).Error
	return err
}

// DeleteVboxBdaChIndexD 删除VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService) DeleteVboxBdaChIndexD(bdaChD vbox.VboxBdaChIndexD) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxBdaChIndexD{}).Where("id = ?", bdaChD.ID).Update("deleted_by", bdaChD.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&bdaChD).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteVboxBdaChIndexDByIds 批量删除VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService) DeleteVboxBdaChIndexDByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxBdaChIndexD{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxBdaChIndexD{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateVboxBdaChIndexD 更新VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService) UpdateVboxBdaChIndexD(bdaChD vbox.VboxBdaChIndexD) (err error) {
	err = global.GVA_DB.Save(&bdaChD).Error
	return err
}

// GetVboxBdaChIndexD 根据id获取VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService) GetVboxBdaChIndexD(id uint) (bdaChD vbox.VboxBdaChIndexD, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&bdaChD).Error
	return
}

// GetVboxBdaChIndexDInfoList 分页获取VboxBdaChIndexD记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChDService *VboxBdaChIndexDService) GetVboxBdaChIndexDInfoList(info vboxReq.VboxBdaChIndexDSearch, idList []int) (list []vbox.VboxBdaChIndexD, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxBdaChIndexD{})
	var bdaChDs []vbox.VboxBdaChIndexD
	db = db.Where("uid in (?)", idList)
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

	err = db.Limit(limit).Offset(offset).Find(&bdaChDs).Error
	return bdaChDs, total, err
}

// Author yoga
func (bdaChDService *VboxBdaChIndexDService) CronVboxBdaChIndexD() (err error) {
	queryB := `
		    SELECT distinct c_channel_id  as chId
			FROM vbox_pay_order
			WHERE  uid = ? and DATE(create_time) = (CURDATE() - INTERVAL ? DAY)
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

			yInfoList, yOrderTotal, err := GetUsersVboxPayOrderInfoList(uid, chId, 1)
			if err != nil {
				return err
			}
			yGroupedCounts := make(map[string]int16)
			yOkGroupedCounts := make(map[string]int16)
			yOkGroupedCosts := make(map[string]int)

			for _, order := range yInfoList {
				uid := strconv.Itoa(order.Uid) + "-" + order.CChannelId
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
			err = global.GVA_DB.Where("product_id = ?", chId).First(&vcp).Error
			if err != nil {
				return err
			}

			chCode := int(vcp.ChannelCode)

			entity := vbox.VboxBdaChIndexD{
				Uid:             &uid,
				UserName:        userInfo.Username,
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
			log.Println("统计结果=", entity)
			err = global.GVA_DB.Save(&entity).Error
		}

	}

	return err
}

func GetUsersVboxPayOrderInfoList(id int, chId string, num int) (list []vbox.VboxPayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxPayOrder{})
	var vpos []vbox.VboxPayOrder
	err = db.Where("uid = ? and c_channel_id = ? and DATE(create_time) = (CURDATE() - INTERVAL ? DAY)", id, chId, num).Find(&vpos).Error
	total = int64(len(vpos))
	return vpos, total, err
}

func GetOwnerUserIdsList(id uint) (list []system.SysUser, total int64, err error) {

	db := global.GVA_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	if id == 1 {
		err = db.Count(&total).Error
		if err != nil {
			return
		}
		err = db.Preload("Authorities").Preload("Authority").Find(&userList).Error

		return userList, total, err
	}
	if err != nil {
		return nil, 0, err
	}

	querySub := `
		WITH RECURSIVE cte (id) AS (SELECT id FROM sys_users WHERE id = ? UNION ALL SELECT sys_users.id FROM sys_users 
			JOIN cte ON sys_users.parent_id = cte.id)
		SELECT id,uuid,username,password,nick_name,side_mode,header_img,base_color,active_color,authority_id,phone,email,enable,parent_id
		FROM sys_users WHERE id IN (SELECT id FROM cte);
    `
	//var userListSub []system.SysUser
	err = db.Raw(querySub, id).Preload("Authorities").Preload("Authority").Find(&userList).Error

	total = int64(len(userList))
	return userList, total, err
}
