package vbox

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxResp "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
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
	fmt.Println("统计开始")
	//bdaChshopIndexDService.CronVboxBdaChShopIndexD()
	fmt.Println("统计结束")
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
	if info.Username != "" {
		db = db.Where("username LIKE ?", "%"+info.Username+"%")
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

// GetBdaChShopIndexDInfoList 分页获取用户通道店铺成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChshopIndexDService *BdaChShopIndexDService) GetBdaChShopIndexDInfoListWeek(info vboxReq.BdaChShopIndexDSearch, ids []uint) (list []vbox.BdaChShopIndexD, total int64, err error) {
	fmt.Println("统计开始")
	//bdaChshopIndexDService.CronVboxBdaChShopIndexD()
	fmt.Println("统计结束")
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	now := time.Now()
	sTime := now.AddDate(0, 0, -6).Format("2006-01-02")
	// 创建db
	//db := global.GVA_DB.Model(&vbox.BdaChShopIndexD{})
	db := global.GVA_DB.Model(&vbox.BdaChShopIndexD{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') >= ? and created_by in ? ", sTime, ids).Order("dt desc")

	var bdaChshopIndexDs []vbox.BdaChShopIndexD
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Uid != nil {
		db = db.Where("uid = ?", info.Uid)
	}
	if info.Username != "" {
		db = db.Where("username LIKE ?", "%"+info.Username+"%")
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.ShopRemark != "" {
		db = db.Where("shop_remark LIKE ?", "%"+info.ShopRemark+"%")
	}
	if info.ProductName != "" {
		db = db.Where("product_name LIKE ?", "%"+info.ProductName+"%")
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

func (bdaChshopIndexDService *BdaChShopIndexDService) CronVboxBdaChShopIndexDByHand(dt string) (err error) {

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

	var shops []string
	err = db.Select("distinct resource_url").Pluck("resource_url", &shops).Error
	if err != nil {
		return
	}

	for _, uid := range uids {
		for _, chId := range chIds {
			for _, shop := range shops {

				yInfoList, yOrderTotal, err := GetUsersShopChVboxPayOrderInfoList(uid, chId, shop, dt)
				if err != nil {
					return err
				}
				yGroupedCounts := make(map[string]int16)
				yOkGroupedCounts := make(map[string]int16)
				yOkGroupedCosts := make(map[string]int)

				for _, order := range yInfoList {
					uid := strconv.Itoa(int(order.CreatedBy)) + "-" + order.ChannelCode + "-" + order.ResourceUrl
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
					Username:        userInfo.Username,
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

func (bdaChshopIndexDService *BdaChShopIndexDService) CronVboxBdaChShopIndexD() (err error) {

	querySql := `

			SELECT
				a.uid as uid,
				a.channel_code as channelCode,
				a.shop_id as shopId,
				step_time as stepTime,
				b.shop_remark as shopRemark,
				c.product_id as productId,
				c.product_name as productName,
				d.username as username,
				a.cnt as orderQuantify,
				a.ok_cnt as okOrderQuantify,
				ROUND((100.0 * a.ok_cnt/a.cnt),2) as ratio,
				a.ok_income as okIncome
			FROM
			(
				SELECT
				 uid,
				 channel_code,
				 shop_id,
				 step_time,
				 sum(if (order_status =1,money,0)) as ok_income,
				 count(*) as cnt,
				 sum(if(order_status=1,1,0)) as ok_cnt
				FROM(
						SELECT
							created_by as uid,
							channel_code,
							order_status,
							SUBSTRING_INDEX(event_id, '_', 1) as shop_id,
							money,
							DATE_FORMAT(created_at, '%Y-%m-%d') AS step_time	
					from vbox_pay_order 
					where DATE_FORMAT(created_at, '%Y-%m-%d') = ?
					and event_type = 1
				)aa
				GROUP BY uid,channel_code,shop_id
			)a 
			left join (
			SELECT DISTINCT cid,shop_remark,product_id  as shop_id
			FROM vbox_channel_shop where address !='' ORDER BY product_id
			) b
			on a.channel_code = b.cid  
			and a.shop_id=b.shop_id
			left join (
			SELECT DISTINCT channel_code,product_id,product_name 
			FROM vbox_channel_product where channel_code !=''
			) c
			on a.channel_code = c.channel_code  
			left join (
			SELECT DISTINCT id,username
			FROM sys_users where username !=''
			) d
			on a.uid = d.id  

`
	dt := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ? ", dt)

	rows, err := db.Raw(querySql, dt).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	for rows.Next() {
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.ChShopTableResp
		err := rows.Scan(&item.UID, &item.ChannelCode, &item.ShopId, &item.StepTime, &item.ShopRemark,
			&item.ProductID, &item.ProductName, &item.Username, &item.OrderQuantify, &item.OkOrderQuantify, &item.Ratio, &item.OkIncome)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		entity := vbox.BdaChShopIndexD{
			Uid:             &item.UID,
			Username:        item.Username,
			Cid:             item.ChannelCode,
			ShopRemark:      item.ShopRemark,
			ChannelCode:     item.ChannelCode,
			ProductId:       item.ProductID,
			ProductName:     item.ProductName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           item.Ratio,
			Income:          item.OkIncome,
			Dt:              dt,
			CreatedBy:       uint(item.UID),
			UpdatedBy:       1,
			DeletedBy:       1,
		}
		//fmt.Println(entity.Dt)
		log.Println("统计 ch shop 结果=", entity)
		err = global.GVA_DB.Save(&entity).Error
		//fmt.Printf("---> UID: %d, AcID: %d, OkIncome: %f, StepTime: %s\n", item.Uid, item.AcID, item.OkIncome, item.StepTime)
	}

	//var uids []int
	//err = db.Select("distinct created_by").Pluck("created_by", &uids).Error
	//if err != nil {
	//	return
	//}
	//fmt.Println("uids=", uids)
	//var chIds []string
	//err = db.Select("distinct channel_code").Pluck("channel_code", &chIds).Error
	//if err != nil {
	//	return
	//}
	//fmt.Println("chIds=", chIds)
	//var addUrls []string
	//err = db.Select("distinct resource_url").Pluck("resource_url", &addUrls).Error
	//if err != nil {
	//	return
	//}
	//fmt.Println("addUrls=", addUrls)
	//var shops []string
	//var vcsp vbox.ChannelShop
	//vcspDb := global.GVA_DB.Model(&vcsp).Where("address is not null and address != '' and address in ? ", addUrls)
	//err = vcspDb.Select("distinct shop_remark").Pluck("shop_remark", &shops).Error
	//if err != nil {
	//	return
	//}
	//fmt.Println("shops=", shops)
	//for _, uid := range uids {
	//	for _, chId := range chIds {
	//		for _, shop := range shops {
	//
	//			yInfoList, yOrderTotal, err := GetUsersShopChVboxPayOrderInfoList(uid, chId, shop, dt)
	//			if err != nil {
	//				return err
	//			}
	//			yGroupedCounts := make(map[string]int16)
	//			yOkGroupedCounts := make(map[string]int16)
	//			yOkGroupedCosts := make(map[string]int)
	//
	//			for _, order := range yInfoList {
	//				uid := strconv.Itoa(int(order.CreatedBy)) + "-" + order.ChannelCode + "-" + order.ResourceUrl
	//				yGroupedCounts[uid]++
	//				if order.OrderStatus == 1 {
	//					yOkGroupedCounts[uid]++
	//					yOkGroupedCosts[uid] += order.Money
	//				}
	//			}
	//
	//			yOrderQuantify := yOrderTotal
	//			yOkOrderQuantify := 0
	//			yOkRate := 0
	//			yInCome := 0
	//			// 判断 tGroupedCounts 中是否包含指定的 uid 键
	//			key := strconv.Itoa(uid) + "-" + chId + "-" + shop
	//			_, yContainsUID := yGroupedCounts[key]
	//			_, yOkContainsUID := yOkGroupedCounts[key]
	//
	//			if yContainsUID {
	//
	//				yOrderQuantify = int64(yGroupedCounts[key])
	//				if yOkContainsUID {
	//					yOkOrderQuantify = int(yOkGroupedCounts[key])
	//				}
	//
	//				if yOrderQuantify > 0 {
	//					result := float64(yOkOrderQuantify) / float64(yOrderQuantify)
	//					yOkRate = int(result * 100)
	//					yInCome = yOkGroupedCosts[key]
	//				}
	//
	//			}
	//
	//			var userInfo system.SysUser
	//			err = global.GVA_DB.Where("`id` = ?", uid).First(&userInfo).Error
	//			if err != nil {
	//				return err
	//			}
	//			var vcp vbox.ChannelProduct
	//			err = global.GVA_DB.Where("channel_code = ?", chId).First(&vcp).Error
	//			if err != nil {
	//				return err
	//			}
	//			var vcs vbox.ChannelShop
	//
	//			if err != nil {
	//				return err
	//			}
	//			shopRemark := ""
	//			if shop != "" {
	//				err = global.GVA_DB.Where("cid is not null and cid != '' and cid = ? and shop_remark = ?", chId, shop).First(&vcs).Error
	//				if err != nil {
	//					return err
	//				}
	//				shopRemark = vcs.ShopRemark
	//			}
	//
	//			chCode := vcp.ChannelCode
	//
	//			entity := vbox.BdaChShopIndexD{
	//				Uid:             &uid,
	//				Username:        userInfo.Username,
	//				Cid:             chId,
	//				ShopRemark:      shopRemark,
	//				ChannelCode:     chCode,
	//				ProductId:       vcp.ProductId,
	//				ProductName:     vcp.ProductName,
	//				OrderQuantify:   int(yOrderQuantify),
	//				OkOrderQuantify: yOkOrderQuantify,
	//				Ratio:           float64(yOkRate),
	//				Income:          yInCome,
	//				Dt:              dt,
	//				CreatedBy:       uint(uid),
	//				UpdatedBy:       1,
	//				DeletedBy:       1,
	//			}
	//			//fmt.Println(entity.Dt)
	//			log.Println("统计 ch shop 结果=", entity)
	//			err = global.GVA_DB.Save(&entity).Error
	//		}
	//	}
	//}
	return err
}

func GetUsersShopChVboxPayOrderInfoList(id int, chId string, shop string, dt string) (list []vbox.PayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var vpos []vbox.PayOrder
	err = db.Where("created_by = ? and resource_url = ? and channel_code = ? and DATE(created_at) = ?", id, shop, chId, dt).Find(&vpos).Error
	total = int64(len(vpos))
	return vpos, total, err
}

func (bdaChshopIndexDService *BdaChShopIndexDService) GetBdaChShopIndexDUesrOverview(res vboxReq.BdaChShopIndexDSearch) (list []vboxResp.ChaShopUserCardResp, total int64, err error) {
	querySql := `
		SELECT
			t1.uid,
			COALESCE(shop_cnt,0) as shopCnt,
			COALESCE(channel_cnt,0) as channelCnt,
			COALESCE(ok_order_cnt,0) as okOrderCnt,
			COALESCE(ok_income,0) as okIncome
		FROM
		(
			SELECT 
				created_by as uid,
 				count(DISTINCT SUBSTRING_INDEX(event_id, '_', 1)) as shop_cnt,
				count(DISTINCT channel_code) as channel_cnt 
			from vbox_pay_order 
			where DATE_FORMAT(created_at, '%Y-%m-%d')= ?
			and created_by = ?
			GROUP BY created_by
		)t1
		LEFT JOIN
		(
		SELECT 
			created_by as uid,
			count(order_id) as ok_order_cnt,
			sum(money) as ok_income
		from vbox_pay_order 
		where DATE_FORMAT(created_at, '%Y-%m-%d')= ?
		and created_by = ?
		and order_status =1
		GROUP BY created_by
		)t2
		on t1.uid=t2.uid
`
	dt := time.Now().AddDate(0, 0, -0).Format("2006-01-02")
	dt_1 := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	dt_2 := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	fmt.Println(dt, "统计开始")

	item, err := getChaShopUserCardResp(querySql, dt, *res.Uid)
	list = append(list, item)
	item_1, err := getChaShopUserCardResp(querySql, dt_1, *res.Uid)
	list = append(list, item_1)
	item_2, err := getChaShopUserCardResp(querySql, dt_2, *res.Uid)
	list = append(list, item_2)

	return list, 3, err
}

func (bdaChshopIndexDService *BdaChShopIndexDService) GetBdaChShopIndexToDayIncome(res vboxReq.BdaChShopIndexDSearch) (data vboxResp.ChartData, err error) {
	querySql := `
			SELECT
			 uid,
			 shop_id as shopId,
			 sum(if(order_status =1,money,0)) as okIncome,
			 step_time as stepTime
			FROM(
					SELECT
							created_by as uid,
							SUBSTRING_INDEX(event_id, '_', 1) as shop_id,
							order_status,
							money,
						CONCAT(
							DATE_FORMAT(created_at, '%H'),
							':',
							LPAD(5 * FLOOR(MINUTE(created_at) / 5), 2, '0')
						) AS step_time	
					from vbox_pay_order 
					where created_at >= ?
					and created_by = ? 
			)t 
			GROUP BY uid,shop_id,step_time
		`
	uid := res.Uid
	// 生成一个包含 24 小时每 5 分钟一个时间点的切片
	var timeData []string
	// 当前时间
	now := time.Now()
	sixHoursAgo := now.Add(-6 * time.Hour)

	startTime := time.Date(sixHoursAgo.Year(), sixHoursAgo.Month(), sixHoursAgo.Day(), sixHoursAgo.Hour(), 0, 0, 0, now.Location())

	fmt.Println("startTime-->", startTime)
	//startTime := now.Add(-6 * time.Hour)
	// 生成 24 小时每 5 分钟一个时间点的时间数据
	for i := 0; i < 6*60/5; i++ {
		// 将当前时间加上 i*5 分钟
		newTime := startTime.Add(time.Minute * 5 * time.Duration(i))
		// 将时间格式化为 'HH:mm' 格式的字符串
		formattedTime := newTime.Format("15:04")
		// 将格式化后的时间字符串添加到切片中
		timeData = append(timeData, formattedTime)
	}

	//dt := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	//db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ? and created_by= ? ", dt, uid)
	sTime := now.Add(-6 * time.Hour).Format("2006-01-02 15:04:05")
	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("created_at >= ? and created_by= ? ", sTime, uid)
	var accIds []string
	err = db.Select("distinct SUBSTRING_INDEX(event_id, '_', 1) as shop_id").Where("order_status=1").Pluck("shop_id", &accIds).Error
	fmt.Println("shopIds=", accIds)

	var shopNames []vbox.ChannelShop
	err = global.GVA_DB.Model(&vbox.ChannelShop{}).Select("DISTINCT product_id, shop_remark").
		Where("product_id in ? ", accIds).Find(&shopNames).Error

	if err != nil {
		return
	}
	// 将查询结果存入 map 中
	shopMap := make(map[string]string)
	for _, shop := range shopNames {
		fmt.Println("shop.ProductId" + shop.ProductId)
		fmt.Println("shop.ShopRemark" + shop.ShopRemark)
		shopMap[shop.ProductId] = shop.ShopRemark
	}

	var shops []string
	for _, accId := range accIds {
		// 在这里处理每个账户 ID
		shops = append(shops, shopMap[accId])
	}
	//
	//db := global.GVA_DB.Model(&vboxResp.UserDayIncomeLineChart{})
	rows, err := db.Raw(querySql, sTime, uid).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 创建一个 map 用于存储结果
	incomeMap := make(map[string]uint)
	// 如果有下一行数据，继续循环
	for rows.Next() {
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.UserDayChShopIncomeLineChart
		err := rows.Scan(&item.Uid, &item.ShopId, &item.OkIncome, &item.StepTime)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		key := item.StepTime + "-" + item.ShopId
		//fmt.Println("key---->", key)
		incomeMap[key] = item.OkIncome

		//fmt.Printf("---> UID: %d, AcID: %d, OkIncome: %f, StepTime: %s\n", item.Uid, item.AcID, item.OkIncome, item.StepTime)
	}

	var seriesData []vboxResp.SeriesItem

	for _, accId := range accIds {
		var incomes []int
		// 内层循环
		for _, time := range timeData {
			//fmt.Printf(time + "-" + accId + "\n")
			// 构建 key
			key := time + "-" + accId
			// 在 incomeMap 中查找 key
			if val, ok := incomeMap[key]; ok {
				//fmt.Println("ok---->", ok)
				// 找到，将值存入数组
				incomes = append(incomes, int(val))
			} else {
				// 没找到，存入默认值 0
				incomes = append(incomes, 0)
			}
		}

		// 处理每个 accId 对应的 incomes 数组
		//fmt.Printf("AccID: %d, Incomes: %v\n", accId, incomes)
		incomesEntity := vboxResp.SeriesItem{
			Name:   shopMap[accId],
			Type:   "line",
			Stack:  "",
			Smooth: true,
			Data:   incomes,
		}
		seriesData = append(seriesData, incomesEntity)
	}

	entity := vboxResp.ChartData{
		LegendData: shops,
		XAxisData:  timeData,
		SeriesData: seriesData,
	}
	//fmt.Printf(" entity: %v\n", entity)
	return entity, err
}

func (bdaChshopIndexDService *BdaChShopIndexDService) GetBdaChShopIndexToDayOkCnt(res vboxReq.BdaChShopIndexDSearch) (data vboxResp.ChartData, err error) {
	querySql := `
			SELECT
			 uid,
			 shop_id as shopId,
			 count(if(order_status =1,1,0)) as okIncome,
			 step_time as stepTime
			FROM(
					SELECT
							created_by as uid,
							SUBSTRING_INDEX(event_id, '_', 1) as shop_id,
							order_status,
							money,
						CONCAT(
							DATE_FORMAT(created_at, '%H'),
							':',
							LPAD(5 * FLOOR(MINUTE(created_at) / 5), 2, '0')
						) AS step_time	
					from vbox_pay_order 
					where created_at >= ?
					and created_by = ? 
			)t 
			GROUP BY uid,shop_id,step_time
		`
	uid := res.Uid
	// 生成一个包含 24 小时每 5 分钟一个时间点的切片
	var timeData []string
	// 当前时间
	now := time.Now()
	sixHoursAgo := now.Add(-6 * time.Hour)

	startTime := time.Date(sixHoursAgo.Year(), sixHoursAgo.Month(), sixHoursAgo.Day(), sixHoursAgo.Hour(), 0, 0, 0, now.Location())

	fmt.Println("startTime-->", startTime)
	//startTime := now.Add(-6 * time.Hour)
	// 生成 24 小时每 5 分钟一个时间点的时间数据
	for i := 0; i < 6*60/5; i++ {
		// 将当前时间加上 i*5 分钟
		newTime := startTime.Add(time.Minute * 5 * time.Duration(i))
		// 将时间格式化为 'HH:mm' 格式的字符串
		formattedTime := newTime.Format("15:04")
		// 将格式化后的时间字符串添加到切片中
		timeData = append(timeData, formattedTime)
	}

	//dt := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	//db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ? and created_by= ? ", dt, uid)
	sTime := now.Add(-6 * time.Hour).Format("2006-01-02 15:04:05")
	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("created_at >= ? and created_by= ? ", sTime, uid)
	var accIds []string
	err = db.Select("distinct SUBSTRING_INDEX(event_id, '_', 1) as shop_id").Where("order_status=1").Pluck("shop_id", &accIds).Error
	fmt.Println("accIds=", accIds)
	if err != nil {
		return
	}
	var shopNames []vbox.ChannelShop
	err = global.GVA_DB.Model(&vbox.ChannelShop{}).Select("DISTINCT product_id, shop_remark").
		Where("product_id in ? ", accIds).Find(&shopNames).Error

	if err != nil {
		return
	}
	// 将查询结果存入 map 中
	shopMap := make(map[string]string)
	for _, shop := range shopNames {
		fmt.Println("shop.ProductId" + shop.ProductId)
		fmt.Println("shop.ShopRemark" + shop.ShopRemark)
		shopMap[shop.ProductId] = shop.ShopRemark
	}
	var shops []string
	for _, accId := range accIds {
		// 在这里处理每个账户 ID
		shops = append(shops, shopMap[accId])
	}

	//
	//db := global.GVA_DB.Model(&vboxResp.UserDayIncomeLineChart{})
	rows, err := db.Raw(querySql, sTime, uid).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 创建一个 map 用于存储结果
	incomeMap := make(map[string]uint)
	// 如果有下一行数据，继续循环
	for rows.Next() {
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.UserDayChShopIncomeLineChart
		err := rows.Scan(&item.Uid, &item.ShopId, &item.OkIncome, &item.StepTime)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		key := item.StepTime + "-" + item.ShopId
		//fmt.Println("key---->", key)
		incomeMap[key] = item.OkIncome

		//fmt.Printf("---> UID: %d, AcID: %d, OkIncome: %f, StepTime: %s\n", item.Uid, item.AcID, item.OkIncome, item.StepTime)
	}

	var seriesData []vboxResp.SeriesItem

	for _, accId := range accIds {
		var incomes []int
		// 内层循环
		for _, time := range timeData {
			//fmt.Printf(time + "-" + accId + "\n")
			// 构建 key
			key := time + "-" + accId
			// 在 incomeMap 中查找 key
			if val, ok := incomeMap[key]; ok {
				//fmt.Println("ok---->", ok)
				// 找到，将值存入数组
				incomes = append(incomes, int(val))
			} else {
				// 没找到，存入默认值 0
				incomes = append(incomes, 0)
			}
		}

		// 处理每个 accId 对应的 incomes 数组
		//fmt.Printf("AccID: %d, Incomes: %v\n", accId, incomes)
		incomesEntity := vboxResp.SeriesItem{
			Name:   shopMap[accId],
			Type:   "line",
			Stack:  "",
			Smooth: true,
			Data:   incomes,
		}
		seriesData = append(seriesData, incomesEntity)
	}

	entity := vboxResp.ChartData{
		LegendData: shops,
		XAxisData:  timeData,
		SeriesData: seriesData,
	}
	//fmt.Printf(" entity: %v\n", entity)
	return entity, err
}

func (bdaChshopIndexDService *BdaChShopIndexDService) GetBdaChShopIndexToWeekInCome(res vboxReq.BdaChShopIndexDSearch) (data vboxResp.ChartData, err error) {
	querySql := `
			SELECT
			 uid,
			 shop_id as shopId,
			 sum(if(order_status =1,money,0)) as okIncome,
			 step_time as stepTime
			FROM(
					SELECT
							created_by as uid,
							SUBSTRING_INDEX(event_id, '_', 1) as shop_id,
							order_status,
							money,
						    DATE_FORMAT(created_at, '%Y-%m-%d') AS step_time	
					from vbox_pay_order 
					where DATE_FORMAT(created_at, '%Y-%m-%d') >= ?
					and created_by = ? 
			)t 
			GROUP BY uid,shop_id,step_time
		`
	uid := res.Uid
	// 生成一个包含 24 小时每 5 分钟一个时间点的切片
	var timeData []string
	// 当前时间
	now := time.Now()
	sTime := now.AddDate(0, 0, -6).Format("2006-01-02")

	for i := 6; i >= 0; i-- {
		t := time.Now().AddDate(0, 0, -i)
		timeData = append(timeData, t.Format("2006-01-02"))
	}

	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') >= ? and created_by= ? ", sTime, uid)
	var accIds []string
	err = db.Select("distinct SUBSTRING_INDEX(event_id, '_', 1) as shop_id").Where("order_status=1").Pluck("shop_id", &accIds).Error
	fmt.Println("accIds=", accIds)
	if err != nil {
		return
	}

	var shopNames []vbox.ChannelShop
	err = global.GVA_DB.Model(&vbox.ChannelShop{}).Select("DISTINCT product_id, shop_remark").
		Where("product_id in ? ", accIds).Find(&shopNames).Error

	if err != nil {
		return
	}
	// 将查询结果存入 map 中
	shopMap := make(map[string]string)
	for _, shop := range shopNames {
		fmt.Println("shop.ProductId" + shop.ProductId)
		fmt.Println("shop.ShopRemark" + shop.ShopRemark)
		shopMap[shop.ProductId] = shop.ShopRemark
	}

	var shops []string
	for _, accId := range accIds {
		// 在这里处理每个账户 ID
		shops = append(shops, shopMap[accId])
	}
	//
	//db := global.GVA_DB.Model(&vboxResp.UserDayIncomeLineChart{})
	rows, err := db.Raw(querySql, sTime, uid).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 创建一个 map 用于存储结果
	incomeMap := make(map[string]uint)
	// 如果有下一行数据，继续循环
	for rows.Next() {
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.UserDayChShopIncomeLineChart
		err := rows.Scan(&item.Uid, &item.ShopId, &item.OkIncome, &item.StepTime)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		key := item.StepTime + "-" + item.ShopId
		//fmt.Println("key---->", key)
		incomeMap[key] = item.OkIncome

		//fmt.Printf("---> UID: %d, AcID: %d, OkIncome: %f, StepTime: %s\n", item.Uid, item.AcID, item.OkIncome, item.StepTime)
	}

	var seriesData []vboxResp.SeriesItem

	for _, accId := range accIds {
		var incomes []int
		// 内层循环
		for _, time := range timeData {
			//fmt.Printf(time + "-" + accId + "\n")
			// 构建 key
			key := time + "-" + accId
			// 在 incomeMap 中查找 key
			if val, ok := incomeMap[key]; ok {
				//fmt.Println("ok---->", int(val))
				// 找到，将值存入数组
				incomes = append(incomes, int(val))
			} else {
				// 没找到，存入默认值 0
				incomes = append(incomes, 0)
			}
		}

		// 处理每个 accId 对应的 incomes 数组
		//fmt.Printf("AccID: %d, Incomes: %v\n", accId, incomes)
		incomesEntity := vboxResp.SeriesItem{
			Name:   shopMap[accId],
			Type:   "line",
			Stack:  "",
			Smooth: true,
			Data:   incomes,
		}
		seriesData = append(seriesData, incomesEntity)
	}

	entity := vboxResp.ChartData{
		LegendData: shops,
		XAxisData:  timeData,
		SeriesData: seriesData,
	}
	//fmt.Printf(" entity: %v\n", entity)
	return entity, err
}

func (bdaChshopIndexDService *BdaChShopIndexDService) GetBdaChaShopIndexToWeekOkCnt(res vboxReq.BdaChShopIndexDSearch) (data vboxResp.ChartData, err error) {

	querySql := `
			SELECT
			 uid,
			 shop_id as shopId,
			 sum(if(order_status =1,1,0)) as okIncome,
			 step_time as stepTime
			FROM(
					SELECT
							created_by as uid,
							SUBSTRING_INDEX(event_id, '_', 1) as shop_id,
							order_status,
							money,
						    DATE_FORMAT(created_at, '%Y-%m-%d') AS step_time	
					from vbox_pay_order 
					where DATE_FORMAT(created_at, '%Y-%m-%d') >= ?
					and created_by = ? 
			)t 
			GROUP BY uid,shop_id,step_time
		`
	uid := res.Uid
	// 生成一个包含 24 小时每 5 分钟一个时间点的切片
	var timeData []string
	// 当前时间
	now := time.Now()
	sTime := now.AddDate(0, 0, -6).Format("2006-01-02")

	for i := 6; i >= 0; i-- {
		t := time.Now().AddDate(0, 0, -i)
		timeData = append(timeData, t.Format("2006-01-02"))
	}

	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') >= ? and created_by= ? ", sTime, uid)
	var accIds []string
	err = db.Select("distinct SUBSTRING_INDEX(event_id, '_', 1) as shop_id").Where("order_status=1").Pluck("shop_id", &accIds).Error
	fmt.Println("accIds=", accIds)
	if err != nil {
		return
	}
	var shopNames []vbox.ChannelShop
	err = global.GVA_DB.Model(&vbox.ChannelShop{}).Select("DISTINCT product_id, shop_remark").
		Where("product_id in ? ", accIds).Find(&shopNames).Error

	if err != nil {
		return
	}
	// 将查询结果存入 map 中
	shopMap := make(map[string]string)
	for _, shop := range shopNames {
		fmt.Println("shop.ProductId" + shop.ProductId)
		fmt.Println("shop.ShopRemark" + shop.ShopRemark)
		shopMap[shop.ProductId] = shop.ShopRemark
	}

	var shops []string
	for _, accId := range accIds {
		// 在这里处理每个账户 ID
		shops = append(shops, shopMap[accId])
	}
	//
	//db := global.GVA_DB.Model(&vboxResp.UserDayIncomeLineChart{})
	rows, err := db.Raw(querySql, sTime, uid).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 创建一个 map 用于存储结果
	incomeMap := make(map[string]uint)
	// 如果有下一行数据，继续循环
	for rows.Next() {
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.UserDayChShopIncomeLineChart
		err := rows.Scan(&item.Uid, &item.ShopId, &item.OkIncome, &item.StepTime)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		key := item.StepTime + "-" + item.ShopId
		//fmt.Println("key---->", key)
		incomeMap[key] = item.OkIncome

		//fmt.Printf("---> UID: %d, AcID: %d, OkIncome: %f, StepTime: %s\n", item.Uid, item.AcID, item.OkIncome, item.StepTime)
	}

	var seriesData []vboxResp.SeriesItem

	for _, accId := range accIds {
		var incomes []int
		// 内层循环
		for _, time := range timeData {
			//fmt.Printf(time + "-" + accId + "\n")
			// 构建 key
			key := time + "-" + accId
			// 在 incomeMap 中查找 key
			if val, ok := incomeMap[key]; ok {
				//fmt.Println("ok---->", ok)
				// 找到，将值存入数组
				incomes = append(incomes, int(val))
			} else {
				// 没找到，存入默认值 0
				incomes = append(incomes, 0)
			}
		}

		// 处理每个 accId 对应的 incomes 数组
		//fmt.Printf("AccID: %d, Incomes: %v\n", accId, incomes)
		incomesEntity := vboxResp.SeriesItem{
			Name:   shopMap[accId],
			Type:   "line",
			Stack:  "",
			Smooth: true,
			Data:   incomes,
		}
		seriesData = append(seriesData, incomesEntity)
	}

	entity := vboxResp.ChartData{
		LegendData: shops,
		XAxisData:  timeData,
		SeriesData: seriesData,
	}
	//fmt.Printf(" entity: %v\n", entity)
	return entity, err
}

func getChaShopUserCardResp(querySql string, dt string, uid int) (res vboxResp.ChaShopUserCardResp, err error) {
	//fmt.Println("dt-->", dt, "uid-->", uid, "querySql-->", querySql)
	db := global.GVA_DB.Model(&vboxResp.ChaShopUserCardResp{})
	rows, err := db.Raw(querySql, dt, uid, dt, uid).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.ChaShopUserCardResp
		err := rows.Scan(&item.Uid, &item.ShopCnt, &item.ChannelCnt, &item.OkOrderCnt, &item.OkIncome)
		if err != nil {
			panic(err)
		}
		item.Dt = dt
		//fmt.Println("--->", item)
		return item, err

	}

	defaultItem := vboxResp.ChaShopUserCardResp{
		Uid:        uid,
		ShopCnt:    0,
		ChannelCnt: 0,
		OkOrderCnt: 0,
		OkIncome:   0,
		Dt:         dt,
	}
	return defaultItem, err

	// 打印查询结果
	//fmt.Println(list)
}
