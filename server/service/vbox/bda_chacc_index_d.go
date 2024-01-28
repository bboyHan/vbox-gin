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
	//bdaChaccIndexDService.CronVboxBdaChaccIndexD()
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
	if info.AcAccount != "" {
		db = db.Where("ac_account = ?", info.AcAccount)
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

// GetBdaChaccIndexDInfoList 分页获取用户通道粒度成率统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChaccIndexDService *BdaChaccIndexDService) GetBdaChaccIndexDInfoListWeek(info vboxReq.BdaChaccIndexDSearch, ids []uint) (list []vbox.BdaChaccIndexD, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	now := time.Now()
	sTime := now.AddDate(0, 0, -6).Format("2006-01-02")
	// 创建db
	db := global.GVA_DB.Model(&vbox.BdaChaccIndexD{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') >= ? and created_by in ? ", sTime, ids).Order("dt desc")
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
	if info.AcAccount != "" {
		db = db.Where("ac_account = ?", info.AcAccount)
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
	// 删除符合条件的数据
	result := global.GVA_DB.Where("dt = ?", dt).Delete(&vbox.BdaChaccIndexD{})
	// 检查删除操作的错误
	if result.Error != nil {
		// 处理错误
	} else {
		// 输出受影响的行数
		fmt.Println(result.RowsAffected)
	}

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

func (bdaChaccIndexDService *BdaChaccIndexDService) GetBdaChaccIndexDUesrOverview(res vboxReq.BdaChaccIndexDSearch) (list []vboxResp.ChaAccUserCardResp, total int64, err error) {
	querySql := `
		SELECT
			t1.uid,
			COALESCE(acid_cnt,0) as acidCnt,
			COALESCE(channel_cnt,0) as channelCnt,
			COALESCE(ok_order_cnt,0) as okOrderCnt,
			COALESCE(ok_income,0) as okIncome
		FROM
		(
			SELECT 
				created_by as uid,
				count(DISTINCT ac_id) as acid_cnt,
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

	item, err := getChaAccUserCardResp(querySql, dt, *res.Uid)
	list = append(list, item)
	item_1, err := getChaAccUserCardResp(querySql, dt_1, *res.Uid)
	list = append(list, item_1)
	item_2, err := getChaAccUserCardResp(querySql, dt_2, *res.Uid)
	list = append(list, item_2)

	return list, 3, err
}

func (bdaChaccIndexDService *BdaChaccIndexDService) GetBdaChaccIndexToDayIncome(res vboxReq.BdaChaccIndexDSearch) (data vboxResp.ChartData, err error) {
	querySql := `
			SELECT
			 uid,
			 ac_id as acId,
			 sum(if(order_status =1,money,0)) as okIncome,
			 step_time as stepTime
			FROM(
					SELECT
							created_by as uid,
							ac_id,
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
			GROUP BY uid,ac_id,step_time
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
	err = db.Select("distinct ac_id").Where("order_status=1").Pluck("ac_id", &accIds).Error
	fmt.Println("accIds=", accIds)
	if err != nil {
		return
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
		var item vboxResp.UserDayIncomeLineChart
		err := rows.Scan(&item.Uid, &item.AcID, &item.OkIncome, &item.StepTime)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		key := item.StepTime + "-" + item.AcID
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
			Name:   accId,
			Type:   "line",
			Stack:  "Total",
			Smooth: true,
			Data:   incomes,
		}
		seriesData = append(seriesData, incomesEntity)
	}

	entity := vboxResp.ChartData{
		LegendData: accIds,
		XAxisData:  timeData,
		SeriesData: seriesData,
	}
	//fmt.Printf(" entity: %v\n", entity)
	return entity, err
}

func (bdaChaccIndexDService *BdaChaccIndexDService) GetBdaChaccIndexToDayOkCnt(res vboxReq.BdaChaccIndexDSearch) (data vboxResp.ChartData, err error) {
	querySql := `
			SELECT
			 uid,
			 ac_id as acId,
			 count(if(order_status =1,1,0)) as okIncome,
			 step_time as stepTime
			FROM(
					SELECT
							created_by as uid,
							ac_id,
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
			GROUP BY uid,ac_id,step_time
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
	err = db.Select("distinct ac_id").Where("order_status=1").Pluck("ac_id", &accIds).Error
	fmt.Println("accIds=", accIds)
	if err != nil {
		return
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
		var item vboxResp.UserDayIncomeLineChart
		err := rows.Scan(&item.Uid, &item.AcID, &item.OkIncome, &item.StepTime)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		key := item.StepTime + "-" + item.AcID
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
			Name:   accId,
			Type:   "line",
			Stack:  "Total",
			Smooth: true,
			Data:   incomes,
		}
		seriesData = append(seriesData, incomesEntity)
	}

	entity := vboxResp.ChartData{
		LegendData: accIds,
		XAxisData:  timeData,
		SeriesData: seriesData,
	}
	//fmt.Printf(" entity: %v\n", entity)
	return entity, err
}

func (bdaChaccIndexDService *BdaChaccIndexDService) GetBdaChaccIndexToWeekIncome(res vboxReq.BdaChaccIndexDSearch) (data vboxResp.ChartData, err error) {
	querySql := `
			SELECT
			 uid,
			 ac_id as acId,
			 sum(if(order_status =1,money,0)) as okIncome,
			 step_time as stepTime
			FROM(
					SELECT
							created_by as uid,
							ac_id,
							order_status,
							money,
						    DATE_FORMAT(created_at, '%Y-%m-%d') AS step_time	
					from vbox_pay_order 
					where DATE_FORMAT(created_at, '%Y-%m-%d') >= ?
					and created_by = ? 
			)t 
			GROUP BY uid,ac_id,step_time
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
	err = db.Select("distinct ac_id").Where("order_status=1").Pluck("ac_id", &accIds).Error
	fmt.Println("accIds=", accIds)
	if err != nil {
		return
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
		var item vboxResp.UserDayIncomeLineChart
		err := rows.Scan(&item.Uid, &item.AcID, &item.OkIncome, &item.StepTime)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		key := item.StepTime + "-" + item.AcID
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
			Name:   accId,
			Type:   "line",
			Stack:  "Total",
			Smooth: true,
			Data:   incomes,
		}
		seriesData = append(seriesData, incomesEntity)
	}

	entity := vboxResp.ChartData{
		LegendData: accIds,
		XAxisData:  timeData,
		SeriesData: seriesData,
	}
	//fmt.Printf(" entity: %v\n", entity)
	return entity, err
}

func (bdaChaccIndexDService *BdaChaccIndexDService) GetBdaChaccIndexToWeekOkCnt(res vboxReq.BdaChaccIndexDSearch) (data vboxResp.ChartData, err error) {

	querySql := `
			SELECT
			 uid,
			 ac_id as acId,
			 sum(if(order_status =1,1,0)) as okIncome,
			 step_time as stepTime
			FROM(
					SELECT
							created_by as uid,
							ac_id,
							order_status,
							money,
						    DATE_FORMAT(created_at, '%Y-%m-%d') AS step_time	
					from vbox_pay_order 
					where DATE_FORMAT(created_at, '%Y-%m-%d') >= ?
					and created_by = ? 
			)t 
			GROUP BY uid,ac_id,step_time
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
	err = db.Select("distinct ac_id").Where("order_status=1").Pluck("ac_id", &accIds).Error
	fmt.Println("accIds=", accIds)
	if err != nil {
		return
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
		var item vboxResp.UserDayIncomeLineChart
		err := rows.Scan(&item.Uid, &item.AcID, &item.OkIncome, &item.StepTime)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		key := item.StepTime + "-" + item.AcID
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
			Name:   accId,
			Type:   "line",
			Stack:  "Total",
			Smooth: true,
			Data:   incomes,
		}
		seriesData = append(seriesData, incomesEntity)
	}

	entity := vboxResp.ChartData{
		LegendData: accIds,
		XAxisData:  timeData,
		SeriesData: seriesData,
	}
	//fmt.Printf(" entity: %v\n", entity)
	return entity, err
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

func getChaAccUserCardResp(querySql string, dt string, uid int) (res vboxResp.ChaAccUserCardResp, err error) {
	//fmt.Println("dt-->", dt, "uid-->", uid, "querySql-->", querySql)
	db := global.GVA_DB.Model(&vboxResp.ChaAccUserCardResp{})
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
		var item vboxResp.ChaAccUserCardResp
		err := rows.Scan(&item.Uid, &item.AcidCnt, &item.ChannelCnt, &item.OkOrderCnt, &item.OkIncome)
		if err != nil {
			panic(err)
		}
		item.Dt = dt
		//fmt.Println("--->", item)
		return item, err

	}

	defaultItem := vboxResp.ChaAccUserCardResp{
		Uid:        uid,
		AcidCnt:    0,
		ChannelCnt: 0,
		OkOrderCnt: 0,
		OkIncome:   0,
		Dt:         dt,
	}
	return defaultItem, err

	// 打印查询结果
	//fmt.Println(list)
}
