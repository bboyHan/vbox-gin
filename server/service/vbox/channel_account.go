package vbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	j3 "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxResp "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	http2 "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ChannelAccountService struct {
}

/*
    创建布隆过滤器
	bloomFilterKey := "myFilter"
	bloomFilterErrorRate := 0.01
	bloomFilterCapacity := 10000
	global.GVA_REDIS.Do(context.Background(), "BF.RESERVE", bloomFilterKey, bloomFilterErrorRate, bloomFilterCapacity)
	if err != nil {
		fmt.Println("创建布隆过滤器失败:", err)
		return
	}

	// 将元素添加到布隆过滤器中
	err = global.GVA_REDIS.Do(context.Background(), "BF.ADD", "myFilter", "hello").Err()
	if err != nil {
		panic(err)
	}

	// 检查元素是否存在于布隆过滤器中
	exists, err := global.GVA_REDIS.Do(context.Background(), "BF.EXISTS", "myFilter", "hello").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(exists) // 输出 true

	if boolValue, ok := exists.(bool); ok {
		// 如果接口类型的值是bool类型，那么boolValue就是该值的bool表示
		// 在这里你可以使用boolValue
		fmt.Println("The interface value is a bool:", boolValue)
	} else {
		// 如果接口类型的值不是bool类型，这里的代码将会执行
		fmt.Println("The interface value is not a bool")
	}
	// 尝试检查不存在的元素
	exists, err = global.GVA_REDIS.Do(context.Background(), "BF.EXISTS", "myFilter", "world").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(exists) // 输出 false

	// 清除布隆过滤器
	err = global.GVA_REDIS.Do(context.Background(), "DEL", "myFilter").Err()
	if err != nil {
		panic(err)
	}
*/

// QueryOrgAccAvailable 查询通道账号的官方记录
func (vcaService *ChannelAccountService) QueryOrgAccAvailable(vca *vbox.ChannelAccount) (res interface{}, err error) {

	// filterKey设置

	// 当前用户所属部门
	allOrgIDs := utils2.GetAllOrgID()

	for _, orgID := range allOrgIDs {
		var accList = make(map[string][]interface{})

		// 当前部门拥有的产品
		chanCodeIDs := utils2.GetChannelCodeByOrgID(orgID)

		// 当前部门的所有用户
		userIDs := utils2.GetUsersByOrgIds([]uint{orgID})

		fmt.Printf("org id : %v  ", orgID)
		fmt.Printf("code ids : %v  ", chanCodeIDs)
		fmt.Printf("user ids : %v\n", userIDs)
		var accTempList []vbox.ChannelAccount

		err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("created_by in ?", userIDs).
			Where("cid in ?", chanCodeIDs).Where("status = ? and sys_status = ?", 1, 1).Find(&accTempList).Error
		if err != nil {
			global.GVA_LOG.Error("查询可用通道账号失败:", zap.Error(err))
			return
		}

		// 根据 chan code 分组
		for _, accT := range accTempList {
			if _, ok := accList[accT.Cid]; ok {
				accList[accT.Cid] = append(accList[accT.Cid], accT.AcAccount)
			} else {
				accList[accT.Cid] = []interface{}{accT.AcAccount}
			}
		}

		// 遍历分组后的acc
		for k, accL := range accList {
			fmt.Printf("key : %s", k)
			fmt.Printf("accL : %v\n", accL)

			blKey := fmt.Sprintf(global.ChanOrgAccFilter, strconv.FormatUint(uint64(orgID), 10), k)
			global.GVA_REDIS.Do(context.Background(), "BF.RESERVE", blKey, global.BloomFilterErrorRate, global.BloomFilterCapacity)

			err = global.GVA_REDIS.Do(context.Background(), "BF.MADD", blKey, "accL", "a", "b").Err()
			if err != nil {
				global.GVA_LOG.Error("查询可用通道账号添加到BL失败:", zap.Error(err))
				return
			}

			var count interface{}
			count, err = global.GVA_REDIS.Do(context.Background(), "BF.COUNT", blKey).Result()
			if err != nil {
				global.GVA_LOG.Error("查询可用通道账号添加到BL失败:", zap.Error(err))
				return
			}

			fmt.Printf("accL count : %v\n", count)

		}

	}

	// deepUserIDs := utils2.GetDeepUserIDs(userId)

	//

	//global.GVA_REDIS.Do(context.Background(), "BF.RESERVE", , global.BloomFilterErrorRate, global.BloomFilterCapacity)
	//if err != nil {
	//	global.GVA_LOG.Error("创建布隆过滤器失败:", zap.Error(err))
	//	return
	//}

	return nil, err
}

// QueryAccOrderHis 查询通道账号的官方记录
func (vcaService *ChannelAccountService) QueryAccOrderHis(vca *vbox.ChannelAccount) (res interface{}, err error) {

	var url string

	c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordQBPrefix).Result()
	if c == 0 {
		var channelCode string
		if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) { // tx系
			channelCode = "qb_proxy"
		} else if global.J3Contains(vca.Cid) {
			channelCode = "j3_proxy"
		}

		err = global.GVA_DB.Debug().Model(&vbox.Proxy{}).Select("url").Where("status = ? and type = ? and chan=?", 1, 1, channelCode).
			First(&url).Error

		if err != nil {
			return nil, errors.New("该信道无资源配置")
		}

		global.GVA_REDIS.Set(context.Background(), global.ProductRecordQBPrefix, url, 10*time.Minute)

	} else {
		url, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordQBPrefix).Result()
	}

	if global.TxContains(vca.Cid) { // tx系

		openID, openKey, err := product.Secret(vca.Token)
		if err != nil {
			return nil, err
		}
		records := product.Records(url, openID, openKey, 24*30*time.Hour)

		if records.Ret != 0 {
			return nil, fmt.Errorf("该账号ck存在异常，请核查")
		}
		//classifier := product.Classifier(records.WaterList)
		return records, nil
	} else if global.J3Contains(vca.Cid) {

		record, err := product.QryJ3Record(*vca)
		if err != nil {
			return nil, err
		}

		timeMax := strconv.FormatInt(time.Now().Unix(), 10)
		global.GVA_LOG.Info("", zap.Any("timeMax", timeMax))
		accKey := fmt.Sprintf(global.J3AccBalanceZSet, vca.AcId)

		list := global.GVA_REDIS.ZRevRangeByScore(context.Background(), accKey, &redis.ZRangeBy{
			Min:    "0",
			Max:    timeMax,
			Offset: 0,
			Count:  100,
		}).Val()
		var accRecords []j3.J3AccountRecord
		for _, mem := range list {
			// 原格式 keyMem := fmt.Sprintf("%s_%s_%v_%d_%d_%d_%d", v.Obj.OrderId, vca.AcAccount, money, nowTimeUnix, hisBalance, checkTime, nowBalance)
			// 映射到J3AccountRecord 中
			keyMem := strings.Split(mem, "_")
			money, _ := strconv.Atoi(keyMem[2])
			//nowTimeStr := keyMem[3]
			//checkTimeStr := keyMem[5]
			//var nowTime, checkTime time.Time
			//nowUnix, _ := strconv.ParseInt(nowTimeStr, 10, 64)
			//if checkTimeStr != "0" {
			//	checkUnix, _ := strconv.ParseInt(checkTimeStr, 10, 64)
			//	checkTime = time.Unix(checkUnix, 0)
			//}
			//nowTime = time.Unix(nowUnix, 0)

			accRecord := j3.J3AccountRecord{
				OrderID:    keyMem[0],
				AcAccount:  keyMem[1],
				Money:      money,
				NowTime:    keyMem[3],
				HisBalance: keyMem[4],
				CheckTime:  keyMem[5],
				NowBalance: keyMem[6],
			}

			accRecords = append(accRecords, accRecord)
		}
		ret := &j3.J3Records{
			J3BalanceData: *record,
			List:          accRecords,
		}
		return ret, nil
	} else if global.PcContains(vca.Cid) {
		openID, openKey, err := product.Secret(vca.Token)
		if err != nil {
			return nil, err
		}
		records := product.Records(url, openID, openKey, 24*30*time.Hour)

		if records.Ret != 0 {
			return nil, fmt.Errorf("该账号ck存在异常，请核查")
		}
		//classifier := product.Classifier(records.WaterList)
		return records, nil
	}

	return res, err
}

// CountAcc 查询可用通道的 当前等待取用的账号个数
func (vcaService *ChannelAccountService) CountAcc(ids []uint) (res []vboxResp.ChannelAccountUnused, err error) {
	err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Select("count(1) as total, cid").Where("status = ? and sys_status = ? and created_by in (?)", 1, 1, ids).
		Group("cid").Order("id desc").Find(&res).Error
	return res, err
}

// CreateChannelAccount 创建通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) CreateChannelAccount(vca *vbox.ChannelAccount, c *gin.Context) (err error) {
	vca.AcId = rand_string.RandomInt(8)
	token := vca.Token
	//增加校验
	if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) {
		_, _, errX := product.Secret(token)
		if errX != nil {
			return errX
		}
	} else if global.J3Contains(vca.Cid) {
		parsedURL, errX := url.Parse(token)
		if errX != nil {
			global.GVA_LOG.Warn("无效的 URL:", zap.Error(errX))
			return errors.New("无效的URL")
		}

		query := parsedURL.Query()
		account := query.Get("account")
		zoneCode := query.Get("zoneCode")
		SN := query.Get("SN")
		sign := query.Get("sign")
		if account == "" || zoneCode == "" || SN == "" || sign == "" {
			return errors.New("账号信息不完整")
		}
		if zoneCode != "z22" {
			return errors.New("仅支持双线二区参数，请核查")
		}
		vca.AcAccount = account
	} else {
		return errors.New("该信道暂不支持创建账号")
	}

	err = global.GVA_DB.Create(vca).Error
	//vca传入的所有值 转化成 vcaDB vbox.ChannelAccount存放

	if vca.Status == 1 {
		go func() {
			var vcaDB vbox.ChannelAccount
			err = global.GVA_DB.Model(vbox.ChannelAccount{}).Where("id =?", vca.ID).First(&vcaDB).Error

			conn, err := mq.MQ.ConnPool.GetConnection()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)

			ch, err := conn.Channel()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
			}

			body := http2.DoGinContextBody(c)

			oc := vboxReq.ChanAccAndCtx{
				Obj: vcaDB,
				Ctx: vboxReq.Context{
					Body:      string(body),
					ClientIP:  c.ClientIP(),
					Method:    c.Request.Method,
					UrlPath:   c.Request.URL.Path,
					UserAgent: c.Request.UserAgent(),
					UserID:    int(vcaDB.CreatedBy),
				},
			}
			marshal, err := json.Marshal(oc)

			err = ch.Publish(task.ChanAccEnableCheckExchange, task.ChanAccEnableCheckKey, marshal)
		}()
	}

	return err
}

// DeleteChannelAccount 删除通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) DeleteChannelAccount(vca vbox.ChannelAccount, c *gin.Context) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var vcaDB vbox.ChannelAccount
		if err := tx.Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).First(&vcaDB).Error; err != nil {
			return err
		}

		conn, err := mq.MQ.ConnPool.GetConnection()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
		}
		defer mq.MQ.ConnPool.ReturnConnection(conn)

		ch, err := conn.Channel()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
		}

		body := http2.DoGinContextBody(c)

		oc := vboxReq.ChanAccAndCtx{
			Obj: vcaDB,
			Ctx: vboxReq.Context{
				Body:      string(body),
				ClientIP:  c.ClientIP(),
				Method:    c.Request.Method,
				UrlPath:   c.Request.URL.Path,
				UserAgent: c.Request.UserAgent(),
				UserID:    int(vcaDB.DeletedBy),
			},
		}
		marshal, err := json.Marshal(oc)

		err = ch.Publish(task.ChanAccDelCheckExchange, task.ChanAccDelCheckKey, marshal)

		if err := tx.Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).Update("sys_status", 2).Error; err != nil {
			return err
		}

		return nil
	})
	return err
}

// DeleteChannelAccountByIds 批量删除通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) DeleteChannelAccountByIds(ids request.IdsReq, c *gin.Context, deletedBy uint) (err error) {

	if len(ids.Ids) < 1 {
		return fmt.Errorf("传入的id为空")
	} else {
		for _, ID := range ids.Ids {
			err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
				var vcaDB vbox.ChannelAccount
				if err := tx.Model(&vbox.ChannelAccount{}).Where("id = ?", ID).First(&vcaDB).Error; err != nil {
					return err
				}

				conn, err := mq.MQ.ConnPool.GetConnection()
				if err != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
				}
				defer mq.MQ.ConnPool.ReturnConnection(conn)

				ch, err := conn.Channel()
				if err != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
				}

				body := http2.DoGinContextBody(c)

				oc := vboxReq.ChanAccAndCtx{
					Obj: vcaDB,
					Ctx: vboxReq.Context{
						Body:      string(body),
						ClientIP:  c.ClientIP(),
						Method:    c.Request.Method,
						UrlPath:   c.Request.URL.Path,
						UserAgent: c.Request.UserAgent(),
						UserID:    int(deletedBy),
					},
				}
				marshal, err := json.Marshal(oc)

				err = ch.Publish(task.ChanAccDelCheckExchange, task.ChanAccDelCheckKey, marshal)

				if err := tx.Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("sys_status", 2).Error; err != nil {
					return err
				}

				return nil
			})
		}
	}

	return err
}

// SwitchEnableChannelAccount 开关通道账号
// Author [bboyhan](https://github.com/bboyhan)
func (vcaService *ChannelAccountService) SwitchEnableChannelAccount(vca vboxReq.ChannelAccountUpd, c *gin.Context) (err error) {
	var vcaDB vbox.ChannelAccount
	err = global.GVA_DB.Where("id = ?", vca.ID).First(&vcaDB).Error
	if err != nil {
		return fmt.Errorf("不存在的账号，请核查")
	}

	// 如果是开启，则发起一条消息，去查这个账号是否能开启

	go func() {
		conn, err := mq.MQ.ConnPool.GetConnection()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
		}
		defer mq.MQ.ConnPool.ReturnConnection(conn)

		ch, err := conn.Channel()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
		}

		body := http2.DoGinContextBody(c)
		vcaDB.Status = vca.Status
		oc := vboxReq.ChanAccAndCtx{
			Obj: vcaDB,
			Ctx: vboxReq.Context{
				Body:      string(body),
				ClientIP:  c.ClientIP(),
				Method:    c.Request.Method,
				UrlPath:   c.Request.URL.Path,
				UserAgent: c.Request.UserAgent(),
				UserID:    int(vcaDB.CreatedBy),
			},
		}
		marshal, err := json.Marshal(oc)

		err = ch.Publish(task.ChanAccEnableCheckExchange, task.ChanAccEnableCheckKey, marshal)
	}()

	err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).Update("status", vca.Status).Update("updated_by", vca.UpdatedBy).Error
	return err
}

// SwitchEnableChannelAccountByIds 批量开关通道账号记录
// Author [bboyhan](https://github.com/bboyhan)
func (vcaService *ChannelAccountService) SwitchEnableChannelAccountByIds(upd vboxReq.ChannelAccountUpd, updatedBy uint, c *gin.Context) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

		// 如果是开启，则发起一条消息，去查这个账号是否能开启

		go func() {
			conn, err := mq.MQ.ConnPool.GetConnection()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)
			var vcaDBList []vbox.ChannelAccount
			err = global.GVA_DB.Model(vbox.ChannelAccount{}).Where("id in ?", upd.Ids).Find(&vcaDBList).Error

			for _, vcaDB := range vcaDBList {

				ch, err := conn.Channel()
				if err != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
					continue
				}

				body := http2.DoGinContextBody(c)
				vcaDB.Status = upd.Status

				oc := vboxReq.ChanAccAndCtx{
					Obj: vcaDB,
					Ctx: vboxReq.Context{
						Body:      string(body),
						ClientIP:  c.ClientIP(),
						Method:    c.Request.Method,
						UrlPath:   c.Request.URL.Path,
						UserAgent: c.Request.UserAgent(),
						UserID:    int(vcaDB.CreatedBy),
					},
				}
				marshal, err := json.Marshal(oc)

				err = ch.Publish(task.ChanAccEnableCheckExchange, task.ChanAccEnableCheckKey, marshal)
			}
		}()

		if err := tx.Model(&vbox.ChannelAccount{}).Where("id in ?", upd.Ids).Update("status", upd.Status).Update("updated_by", updatedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", upd.Ids).Updates(&vbox.ChannelAccount{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannelAccount 更新通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) UpdateChannelAccount(vca vbox.ChannelAccount) (err error) {
	token := vca.Token
	//增加校验
	if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) {
		_, _, errX := product.Secret(token)
		if errX != nil {
			return errX
		}
	} else if global.J3Contains(vca.Cid) {
		parsedURL, errX := url.Parse(token)
		if errX != nil {
			global.GVA_LOG.Warn("无效的 URL:", zap.Error(errX))
			return errors.New("无效的URL")
		}

		query := parsedURL.Query()
		account := query.Get("account")
		zoneCode := query.Get("zoneCode")
		SN := query.Get("SN")
		sign := query.Get("sign")
		if account == "" || zoneCode == "" || SN == "" || sign == "" {
			return errors.New("账号信息不完整")
		}
		if zoneCode != "z22" {
			return errors.New("仅支持双线二区参数，请核查")
		}
		vca.AcAccount = account
	} else {
		return errors.New("该信道暂不支持创建账号")
	}
	err = global.GVA_DB.Save(&vca).Error
	return err
}

// GetChannelAccount 根据id获取通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) GetChannelAccount(id uint) (vca vbox.ChannelAccount, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vca).Error
	return
}

// GetChannelAccountByAcId 根据AcId获取通道账号记录
func (vcaService *ChannelAccountService) GetChannelAccountByAcId(acId string) (vca vbox.ChannelAccount, err error) {
	err = global.GVA_DB.Where("ac_id = ?", acId).First(&vca).Error
	return
}

// GetChannelAccountInfoList 分页获取通道账号记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcaService *ChannelAccountService) GetChannelAccountInfoList(info vboxReq.ChannelAccountSearch, ids []uint) (list []vbox.ChannelAccount, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelAccount{})
	var vcas []vbox.ChannelAccount
	db.Where("created_by in (?)", ids)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.AcRemark != "" {
		db = db.Where("ac_remark LIKE ?", "%"+info.AcRemark+"%")
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account LIKE ?", "%"+info.AcAccount+"%")
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	if info.SysStatus != nil {
		db = db.Where("sys_status = ?", info.SysStatus)
	}
	if info.AcId != "" {
		db = db.Where("ac_id = ?", info.AcId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Order("id desc").Find(&vcas).Error
	return vcas, total, err
}
