package vbox

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxResp "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"strconv"
	"strings"
	"time"
)

type ChannelPayCodeService struct {
}

// GetPayCodeOverviewByChanAcc 获取指定通道账号的预产统计情况(根据acc)
func (channelPayCodeService *ChannelPayCodeService) GetPayCodeOverviewByChanAcc(info vboxReq.ChannelPayCodeSearch, ids []uint) (ret []vboxResp.DataSExtOverView, err error) {

	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelPayCode{}).Table("vbox_channel_pay_code")
	db.Where("created_by in ?", ids)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account like ?", "%"+info.AcAccount+"%")
	}
	if info.AcId != "" {
		db = db.Where("ac_id = ?", info.AcId)
	}
	if info.Location != "" {
		db = db.Where("location = ?", info.Location)
	}
	if info.Money != 0 {
		db = db.Where("money = ?", info.Money)
	}
	if info.Operator != "" {
		db = db.Where("operator = ?", info.Operator)
	}
	if info.Mid != "" {
		db = db.Where("mid = ?", info.Mid)
	}
	if info.CodeStatus != 0 {
		db = db.Where("code_status = ?", info.CodeStatus)
	}

	// x1 money x2 operator x3 location x4 count
	err = db.Debug().Select("money as x1,operator as x2,location as x3, count(1) as x4").Group("x1,x2,x3").Find(&ret).Error

	return ret, err
}

// GetPayCodeOverview 获取预产统计情况
func (channelPayCodeService *ChannelPayCodeService) GetPayCodeOverview(info vboxReq.ChannelPayCodeSearch, ids []uint) (ret []vboxResp.DataSExtOverView, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelPayCode{}).Table("vbox_channel_pay_code")
	db.Where("created_by in ?", ids)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account like ?", "%"+info.AcAccount+"%")
	}
	if info.AcId != "" {
		db = db.Where("ac_id = ?", info.AcId)
	}
	if info.Location != "" {
		db = db.Where("location = ?", info.Location)
	}
	if info.Money != 0 {
		db = db.Where("money = ?", info.Money)
	}
	if info.Operator != "" {
		db = db.Where("operator = ?", info.Operator)
	}
	if info.Mid != "" {
		db = db.Where("mid = ?", info.Mid)
	}
	if info.CodeStatus != 0 {
		db = db.Where("code_status = ?", info.CodeStatus)
	}

	// x1 money x2 operator x3 location x4 count
	err = db.Debug().Select("money as x1,operator as x2,location as x3, count(1) as x4").Group("x1,x2,x3").Find(&ret).Error

	return ret, nil
}

func (channelPayCodeService *ChannelPayCodeService) CreateChannelPayCode(vboxChannelPayCode *vbox.ChannelPayCode) (err error) {

	mid := time.Now().Format("20060102150405") + rand_string.RandomInt(3)
	vboxChannelPayCode.Mid = mid

	// 先查一下库中记录
	// 查一下数据库中预产对应的acAccount、money一致，并且code_status=2(取用池已经有待取用的预产),则将当前记录放入等候池

	// 组织
	orgTmp := utils2.GetSelfOrg(vboxChannelPayCode.CreatedBy)

	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)
	ch, err := conn.Channel()
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
	}

	// 入取用池
	pcKey := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0],
		vboxChannelPayCode.Cid, vboxChannelPayCode.Money, vboxChannelPayCode.Operator, vboxChannelPayCode.Location)

	pattern := fmt.Sprintf(global.ChanOrgPayCodePrefix, orgTmp[0], vboxChannelPayCode.Cid, vboxChannelPayCode.Money)
	keys := global.GVA_REDIS.Keys(context.Background(), pattern).Val()

	var flag bool
	for _, key := range keys {
		waitMembers := global.GVA_REDIS.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{Min: "4", Max: "4", Offset: 0, Count: -1}).Val()

		for _, member := range waitMembers {
			if strings.Contains(member, vboxChannelPayCode.AcAccount) {
				flag = true
				break
			}
		}
	}
	if flag {
		global.GVA_LOG.Info("当前添加的账号正在冷却中（有预产正在处理中）", zap.Any("acc", vboxChannelPayCode.AcAccount))
		vboxChannelPayCode.CodeStatus = 4
		err = global.GVA_DB.Create(vboxChannelPayCode).Error

		waitAccPcKey := fmt.Sprintf(global.AccWaiting, vboxChannelPayCode.AcId)

		// 设置一个冷却时间
		var cdTime time.Duration
		ttl := global.GVA_REDIS.TTL(context.Background(), waitAccPcKey).Val()
		if ttl > 0 {
			global.GVA_LOG.Info("当前添加的账号正在冷却中（有预产正在处理中）", zap.Any("ttl", ttl))
			cdTime = ttl
		} else {
			duration, _ := HandleExpTime2Product(vboxChannelPayCode.Cid)
			cdTime = duration + 60*time.Second
		}

		// 把当前acAccount下所有的预产等待队列置为冷却状态
		waitIDsTmp := strings.Join([]string{fmt.Sprintf("%d", vboxChannelPayCode.ID)}, ",")
		global.GVA_REDIS.Set(context.Background(), waitAccPcKey, waitIDsTmp, cdTime)

		waitMsg := strings.Join([]string{waitAccPcKey, waitIDsTmp}, "_")
		err = ch.PublishWithDelay(task.PayCodeCDCheckDelayedExchange, task.PayCodeCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)

		pcMem := fmt.Sprintf("%d", vboxChannelPayCode.ID) + "_" + vboxChannelPayCode.Mid + "_" + vboxChannelPayCode.AcAccount + "_" + vboxChannelPayCode.ImgContent
		global.GVA_REDIS.ZAdd(context.Background(), pcKey, redis.Z{Score: 4, Member: pcMem})
	} else {
		global.GVA_LOG.Info("当前添加的账号没有冷却中（没有预产正在处理中）")
		vboxChannelPayCode.CodeStatus = 2

		err = global.GVA_DB.Create(vboxChannelPayCode).Error

		pcMem := fmt.Sprintf("%d", vboxChannelPayCode.ID) + "_" + vboxChannelPayCode.Mid + "_" + vboxChannelPayCode.AcAccount + "_" + vboxChannelPayCode.ImgContent
		global.GVA_REDIS.ZAdd(context.Background(), pcKey, redis.Z{Score: 0, Member: pcMem})
	}

	//根据expTime 处理到期的消息校验，放到PayCodeDelayedRoutingKey
	if vboxChannelPayCode.ExpTime.Unix() > 0 {
		// 过期时间
		expTime := vboxChannelPayCode.ExpTime
		// 过期时间差
		expTimeDiff := expTime.Sub(time.Now())
		global.GVA_LOG.Info("过期时间差", zap.Any("expTimeDiff", expTimeDiff))
		marshal, _ := json.Marshal(vboxChannelPayCode)

		err = ch.PublishWithDelay(task.PayCodeDelayedExchange, task.PayCodeDelayedRoutingKey, marshal, expTimeDiff)
		global.GVA_LOG.Info("消息发完了", zap.Any("expTimeDiff", expTimeDiff))
	}

	return err
}

// DeleteChannelPayCode 删除通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) DeleteChannelPayCode(vboxChannelPayCode vbox.ChannelPayCode) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

		var pcDB vbox.ChannelPayCode
		// 先查数据是否存在
		err = tx.Model(&vbox.ChannelPayCode{}).Where("id = ?", vboxChannelPayCode.ID).First(&pcDB).Error
		if err != nil {
			return err
		} else {
			//	处理掉待用池子中的付款码
			orgTmp := utils2.GetSelfOrg(pcDB.CreatedBy)

			// 删待取池中数据
			key := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0],
				pcDB.Cid, pcDB.Money, pcDB.Operator, pcDB.Location)
			pcMem := fmt.Sprintf("%d", pcDB.ID) + "_" + pcDB.Mid + "_" + pcDB.AcAccount + "_" + pcDB.ImgContent
			global.GVA_REDIS.ZRem(context.Background(), key, pcMem)

		}

		if err = tx.Model(&vbox.ChannelPayCode{}).Where("id = ?", vboxChannelPayCode.ID).Update("deleted_by", vboxChannelPayCode.DeletedBy).Error; err != nil {
			return err
		}

		if err = tx.Delete(&vboxChannelPayCode).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteChannelPayCodeByIds 批量删除通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) DeleteChannelPayCodeByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var pcDBList []vbox.ChannelPayCode
		if err := tx.Model(&vbox.ChannelPayCode{}).Where("id in ?", ids.Ids).Find(&pcDBList).Error; err != nil {
			return err
		}
		for _, pcDB := range pcDBList {
			// 处理掉待用池子中的付款码
			orgTmp := utils2.GetSelfOrg(pcDB.CreatedBy)

			// 删待取池中数据
			key := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0],
				pcDB.Cid, pcDB.Money, pcDB.Operator, pcDB.Location)
			pcMem := fmt.Sprintf("%d", pcDB.ID) + "_" + pcDB.Mid + "_" + pcDB.AcAccount + "_" + pcDB.ImgContent
			global.GVA_REDIS.ZRem(context.Background(), key, pcMem)

		}

		if err = tx.Model(&vbox.ChannelPayCode{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.ChannelPayCode{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannelPayCode 更新通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) UpdateChannelPayCode(vboxChannelPayCode vbox.ChannelPayCode) (err error) {
	err = global.GVA_DB.Save(&vboxChannelPayCode).Error
	return err
}

// GetChannelPayCode 根据id获取通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) GetChannelPayCode(id uint) (vboxChannelPayCode vbox.ChannelPayCode, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vboxChannelPayCode).Error
	return
}

// GetChannelPayCodeInfoList 分页获取通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) GetChannelPayCodeInfoList(info vboxReq.ChannelPayCodeSearch, ids []uint) (list []vbox.ChannelPayCode, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelPayCode{})
	var vboxChannelPayCodes []vbox.ChannelPayCode
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account like ?", "%"+info.AcAccount+"%")
	}
	if info.Location != "" {
		db = db.Where("location = ?", info.Location)
	}
	if info.Operator != "" {
		db = db.Where("operator = ?", info.Operator)
	}
	if info.Mid != "" {
		db = db.Where("mid = ?", info.Mid)
	}
	if info.CodeStatus != 0 {
		db = db.Where("code_status = ?", info.CodeStatus)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Where("created_by in ?", ids).Order("id desc").Find(&vboxChannelPayCodes).Error
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	return vboxChannelPayCodes, total, err
}

/*
*
 */
func (channelPayCodeService *ChannelPayCodeService) GetChannelPayCodeNumsByLocation(info vboxReq.ChannelPayCodeSearch, ids []uint) (list []vboxResp.ChannelPayCodeStatistics, total int64, err error) {

	query := `
		    SELECT
			 code as location,count(mid) as codeNums
			FROM(
			
				SELECT  
					t.mid,
					a.name as code
				from (
					SELECT location, mid
					from vbox_channel_pay_code
					where  location !='' and created_by in ?
					and operator in ? and code_status in ?
				) t 
				join geo_provinces a 
				on a.code = SUBSTRING(t.location,1,?)
			
			)b
			GROUP BY code
		    ORDER BY codeNums desc
		;
		`

	querySubF := `
			SELECT
			    location,
			   codeNums
			from (
			SELECT code AS location,
				count( mid ) AS codeNums 
			FROM
				(
				SELECT
					t.mid,
					a.name AS code 
				FROM
					( SELECT 
					      location,
					      mid 
					  FROM vbox_channel_pay_code 
					  WHERE location != '' 
						   AND created_by IN ? 
						   AND SUBSTRING( location, 1,? ) = ?
					  	and operator in ? and code_status in ?
					  ) t
					JOIN geo_cities a ON a.code = SUBSTRING( t.location, 1,? ) 
					AND LENGTH( location ) >= ? 
				) b 
			GROUP BY code
			UNION ALL
			SELECT
				'无具体省市区' AS location,
				count( mid ) AS codeNums 
			FROM
				( SELECT location, mid FROM vbox_channel_pay_code
			   WHERE location != '' AND created_by IN ? AND location = ?
			   and operator in ? and code_status in ?
			   ) c 
			) d 
		where codeNums > 0
		ORDER BY
				codeNums DESC
		;
		`
	querySubS := `
			SELECT
				location,
				codeNums 
			FROM
				(
				SELECT code AS location,
					count(mid) AS codeNums 
				FROM
					(
					SELECT
						t.mid,
						a.name AS code 
					FROM
						(
						SELECT
							location,
							mid 
						FROM
							vbox_channel_pay_code 
						WHERE location != '' 
							AND created_by IN ? 
							AND SUBSTRING( location, 1,? ) = ? 
							AND LENGTH( location ) = ? 
						and operator in ? and code_status in ?
						) t
						JOIN geo_areas a ON a.code = SUBSTRING( t.location, 1,? ) 
					) b 
				GROUP BY code
				UNION ALL
				SELECT
					'无具体省市区' AS location,
					count( mid ) AS codeNums 
				FROM
					( SELECT location, mid FROM vbox_channel_pay_code
				   WHERE location != '' AND created_by IN ? AND location = ? 
				   and operator in ? and code_status in ?
				   ) c 
				) d 
			where codeNums > 0
			ORDER BY
				codeNums DESC
		;
		`

	querySubCity := `
		SELECT
			a.name AS location ,
			codeNums
		FROM
			(
			SELECT
				location,
				count(mid) as codeNums 
			FROM
				vbox_channel_pay_code 
			WHERE location != '' 
				AND created_by IN ? 
				AND location = ?
			 and operator in ? and code_status in ?
			group by location
			) t
			JOIN geo_areas a
			ON a.code =  t.location 
;
		`

	fmt.Println("CodeStatus=", info.CodeStatus, "operators", info.Operator)
	codeStatus := []uint{1, 2, 3}
	if info.CodeStatus != 0 {
		codeStatus = []uint{info.CodeStatus}
	}
	operators := []string{"yidong", "liantong", "dianxin"}
	if info.Operator != "" {
		operators = []string{info.Operator}
	}

	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelPayCode{})

	var totalGroup int64 = 0

	//fmt.Println("info.Location:", info.Location)

	var codeStatisResultList []vboxResp.ChannelPayCodeStatisticsResult
	// 全国各省
	if info.Location == "" {
		//fmt.Println("0 >>>>>")
		rows, err := db.Raw(query, ids, operators, codeStatus, 2).Rows()
		if err != nil {
			// 处理错误
		}

		defer rows.Close()
		for rows.Next() {
			var result vboxResp.ChannelPayCodeStatisticsResult
			err := rows.Scan(&result.Location, &result.CodeNums)
			if err != nil {
				// 处理错误
			}
			codeStatisResultList = append(codeStatisResultList, result)
			totalGroup += int64(result.CodeNums)
		}

	}

	if len(info.Location) == 2 {
		//fmt.Println("2 >>>>>")
		rows, err := db.Raw(querySubF, ids, 2, info.Location, operators, codeStatus, 4, 4, ids, info.Location, operators, codeStatus).Rows()
		if err != nil {
			// 处理错误
		}
		defer rows.Close()
		for rows.Next() {
			var result vboxResp.ChannelPayCodeStatisticsResult
			err := rows.Scan(&result.Location, &result.CodeNums)
			if err != nil {
				// 处理错误
			}
			codeStatisResultList = append(codeStatisResultList, result)
			totalGroup += int64(result.CodeNums)
		}
	}

	if len(info.Location) == 4 {
		//fmt.Println("4 >>>>>")
		rows, err := db.Raw(querySubS, ids, 4, info.Location, 6, operators, codeStatus, 6, ids, info.Location, operators, codeStatus).Rows()
		if err != nil {
			// 处理错误
		}
		defer rows.Close()
		for rows.Next() {
			var result vboxResp.ChannelPayCodeStatisticsResult
			err := rows.Scan(&result.Location, &result.CodeNums)
			if err != nil {
				// 处理错误
			}
			codeStatisResultList = append(codeStatisResultList, result)
			totalGroup += int64(result.CodeNums)
		}

	}
	if len(info.Location) == 6 {
		//fmt.Println("6 >>>>>")
		rows, err := db.Raw(querySubCity, ids, info.Location, operators, codeStatus).Rows()
		if err != nil {
			// 处理错误
		}
		defer rows.Close()
		for rows.Next() {
			var result vboxResp.ChannelPayCodeStatisticsResult
			// scan中有严格的字段顺序 要和sql中的一致
			err := rows.Scan(&result.Location, &result.CodeNums)
			if err != nil {
				// 处理错误
			}
			codeStatisResultList = append(codeStatisResultList, result)
			totalGroup += int64(result.CodeNums)
		}

	}
	for i, statis := range codeStatisResultList {
		//fmt.Println("A num: ", i, "code: ", statis.Location, "total", totalGroup, "codeNums", statis.CodeNums)
		ratio := math.Round(float64(statis.CodeNums)/float64(totalGroup)*10000) / 10000
		entity := vboxResp.ChannelPayCodeStatistics{
			Order:    uint(i + 1),
			Location: statis.Location,
			CodeNums: statis.CodeNums,
			Ratio:    ratio,
		}
		list = append(list, entity)
	}
	total = int64(len(list))
	fmt.Println("total:", total)
	if total == 0 {
		entity := vboxResp.ChannelPayCodeStatistics{
			Order:    1,
			Location: "无",
			CodeNums: 0,
			Ratio:    1,
		}
		list = append(list, entity)
	}
	return list, total, err

}

func HandleExpTime2Product(chanID string) (time.Duration, error) {
	var key string

	if global.TxContains(chanID) {
		key = "1000"
	} else if global.J3Contains(chanID) {
		key = "2000"
	} else if global.PcContains(chanID) {
		key = "3000"
	}

	var expTimeStr string
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		global.GVA_LOG.Warn("当前key不存在", zap.Any("key", key))

		var proxy vbox.Proxy
		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
		err = db.Where("status = ?", 1).Where("chan = ?", key).
			First(&proxy).Error
		if err != nil || proxy.Url == "" {
			return 0, err
		}
		expTimeStr = proxy.Url
		seconds, _ := strconv.Atoi(expTimeStr)
		duration := time.Duration(seconds) * time.Second

		global.GVA_REDIS.Set(context.Background(), key, int64(duration.Seconds()), 0)
		global.GVA_LOG.Info("数据库取出该产品的有效时长", zap.Any("channel code", chanID), zap.Any("过期时间(s)", seconds))

		return duration, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
		return 0, err
	} else {
		expTimeStr, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		seconds, _ := strconv.Atoi(expTimeStr)

		duration := time.Duration(seconds) * time.Second

		//global.GVA_LOG.Info("缓存池取出该产品的有效时长", zap.Any("channel code", chanID), zap.Any("过期时间(s)", seconds))
		return duration, err
	}
}
