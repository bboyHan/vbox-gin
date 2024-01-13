package task

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

func HandleAccAvailable() (err error) {

	var accDBList []vbox.ChannelAccount
	global.GVA_DB.Model(&vbox.ChannelAccount{}).Table("vbox_channel_account").
		Where("status = ? and sys_status = ?", 1, 1).Find(&accDBList)

	//global.GVA_LOG.Info("根据开启的商户列表，开始检测可用账号情况")

	for _, accDB := range accDBList {
		accDBTmp := accDB
		go func() {
			//查一下订单是否超出账户限制
			var flag bool
			cid := accDBTmp.Cid

			// 1. 查询该用户的余额是否充足
			var balance int
			err = global.GVA_DB.Model(&vbox.UserWallet{}).Select("IFNULL(sum(recharge), 0) as balance").
				Where("uid = ?", accDBTmp.CreatedBy).Scan(&balance).Error

			if balance <= 0 { //余额不足，则 log 一条
				//入库操作记录
				flag = true

				msgX := fmt.Sprintf(global.BalanceNotEnough, accDBTmp.AcId, accDBTmp.AcAccount)

				global.GVA_LOG.Error("余额不足...", zap.Any("msg", msgX))
				err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", accDBTmp.ID).
					Update("status", 0).Update("sys_status", 0).Error
			}

			// 2.1 日限制
			if accDBTmp.DailyLimit > 0 {
				var dailySum int
				// 获取今天的时间范围
				startOfDay := time.Now().UTC().Truncate(24 * time.Hour)
				endOfDay := startOfDay.Add(24 * time.Hour)

				err = global.GVA_DB.Model(&vbox.PayOrder{}).Select("sum(money) as dailySum").
					Where("ac_id = ?", accDBTmp.AcId).
					Where("order_status = ? AND created_at BETWEEN ? AND ?", 1, startOfDay, endOfDay).Scan(&dailySum).Error

				if dailySum >= accDBTmp.DailyLimit { // 如果日消费已经超了，不允许开启了，直接结束
					flag = true

					msg := fmt.Sprintf(global.AccDailyLimitNotEnough, accDBTmp.AcId, accDBTmp.AcAccount)
					global.GVA_LOG.Error("当前账号日消耗已经超限...", zap.Any("msg", msg))
					err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", accDBTmp.ID).
						Update("status", 0).Update("sys_status", 0).Error
				}
			}
			// 2.2 总限制
			if accDBTmp.TotalLimit > 0 {

				var totalSum int

				err = global.GVA_DB.Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as totalSum").
					Where("ac_id = ?", accDBTmp.AcId).
					Where("order_status = ?", 1).Scan(&totalSum).Error

				if err != nil {
					global.GVA_LOG.Error("当前账号计算总消耗查mysql错误，直接丢了..." + err.Error())
				}

				if totalSum >= accDBTmp.TotalLimit { // 如果总消费已经超了，不允许开启了，直接结束
					flag = true

					//入库操作记录
					msgX := fmt.Sprintf(global.AccTotalLimitNotEnough, accDBTmp.AcId, accDBTmp.AcAccount)
					global.GVA_LOG.Error("当前账号总消耗已经超限...", zap.Any("msg", msgX))

					err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", accDBTmp.ID).
						Update("status", 0).Update("sys_status", 0).Error

					global.GVA_LOG.Info("当前账号总消耗已经超限额了，结束...", zap.Any("ac info", accDBTmp))
				}
			}
			// 2.3 笔数限制
			if accDBTmp.CountLimit > 0 {

				var count int64

				err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("ac_id = ? and order_status = ?", accDBTmp.AcId, 1).Count(&count).Error

				if err != nil {
					global.GVA_LOG.Error("当前账号笔数消耗查mysql错误，直接丢了..." + err.Error())
				}

				if int(count) >= accDBTmp.CountLimit { // 如果笔数消费已经超了，不允许开启了，直接结束

					flag = true
					msgX := fmt.Sprintf(global.AccCountLimitNotEnough, accDBTmp.AcId, accDBTmp.AcAccount)

					global.GVA_LOG.Error("当前账号笔数消耗已经超限额...", zap.Any("msg", msgX))
					err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", accDBTmp.ID).
						Update("status", 0).Update("sys_status", 0).Error
					global.GVA_LOG.Warn("当前账号笔数消耗已经超限额了，结束...", zap.Any("ac info", accDBTmp))
				}
			}

			if flag {
				global.GVA_LOG.Warn("当前账号已经超限额了，处理一下...", zap.Any("ac info", accDBTmp))

				if global.PcContains(cid) {
					orgTmp := utils2.GetSelfOrg(accDBTmp.CreatedBy)
					orgID := orgTmp[0]
					pattern := fmt.Sprintf(global.ChanOrgPayCodeMoneyPrefix, orgID, cid)

					var keys []string
					keys = global.GVA_REDIS.Keys(context.Background(), pattern).Val() //拿出所有该账号的码，全部处理掉

					for _, key := range keys {
						resWaitTmpList := global.GVA_REDIS.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
							Min:    "0",
							Max:    "0",
							Offset: 0,
							Count:  -1,
						}).Val()

						for _, waitMem := range resWaitTmpList {
							if strings.Contains(waitMem, accDBTmp.AcAccount) {
								//	把超限的码全部处理掉
								global.GVA_REDIS.ZRem(context.Background(), key, waitMem)

								// 把 pay code中属于该账号的码全部处理掉
								id := strings.Split(waitMem, "_")[0]
								global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id = ? ", id).Update("code_status", 4)
								global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ? ", accDBTmp.ID).
									Update("status", 0).Update("sys_status", 0)
							}
						}
					}
				}
			}
		}()

	}

	return err
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
