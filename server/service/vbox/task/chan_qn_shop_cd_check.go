package task

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
)

// 订单等待查询
const (
	QNShopCDCheckDelayedExchange   = "vbox.order.qn_shop_cd_delayed_exchange"
	QNShopCDCheckDelayedRoutingKey = "vbox.order.qn_shop_cd_delayed_routing_key"
	QNShopCDCheckDelayedQueue      = "vbox.order.qn_shop_cd_delayed_queue"
	QNShopCDCheckDeadRoutingKey    = "vbox.order.qn_shop_cd_dead_routing_key"
	QNShopCDCheckDeadExchange      = "vbox.order.qn_shop_cd_dead_exchange"
	QNShopCDCheckDeadQueue         = "vbox.order.qn_shop_cd_dead_queue"
)

// QNShopCDCheckTask 引导类-QNShop-冷却检查
func QNShopCDCheckTask() {

	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		//log.Fatalf("Failed to get connection from pool: %v", err)
		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
	ch, _ := conn.Channel()
	if err = ch.ExchangeDeclare(QNShopCDCheckDeadExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclare(QNShopCDCheckDeadQueue); err != nil {
		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueBind(QNShopCDCheckDeadQueue, QNShopCDCheckDeadRoutingKey, QNShopCDCheckDeadExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclareWithDelay(QNShopCDCheckDelayedQueue, QNShopCDCheckDeadExchange, QNShopCDCheckDeadRoutingKey); err != nil {
		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
	}
	if err = ch.ExchangeDeclare(QNShopCDCheckDelayedExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
	}
	if err = ch.QueueBind(QNShopCDCheckDelayedQueue, QNShopCDCheckDelayedRoutingKey, QNShopCDCheckDelayedExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 333:", zap.Any("err", err))
	}

	// 设置初始消费者数量
	consumerCount := 20
	// 使用 WaitGroup 来等待所有消费者完成处理
	var wg sync.WaitGroup
	wg.Add(consumerCount)
	// 启动多个消费者
	for i := 0; i < consumerCount; i++ {
		go func(consumerID int) {

			connX, errX := mq.MQ.ConnPool.GetConnection()
			if errX != nil {
				//log.Fatalf("Failed to get connection from pool: %v", err)
				global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(errX))
			}
			defer mq.MQ.ConnPool.ReturnConnection(connX)
			if connX == nil {
				global.GVA_LOG.Error("connX == nil", zap.Any("err", errX))
				return
			}
			chX, _ := connX.Channel()

			// 说明：执行查单回调处理
			deliveries, errC := chX.Consume(QNShopCDCheckDeadQueue, "", false, false, false, false, nil)
			if errC != nil {
				global.GVA_LOG.Error("mq 消费者异常， err", zap.Error(errC), zap.Any("queue", QNShopCDCheckDeadQueue))
			}

			for msg := range deliveries {
				//err = handler(msg.Body)
				v := string(msg.Body)

				split := strings.Split(v, "-")
				waitAccYdKey := split[0]
				waitAccMem := split[1]
				var accInfo []string
				accInfo = strings.Split(waitAccMem, ",")

				ID := accInfo[0]
				MID := accInfo[1]
				markID := accInfo[2]
				money := accInfo[3]
				uid := accInfo[4]
				var vcaList []vbox.ChannelAccount
				global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("created_by =?", uid).Scan(&vcaList)

				if len(vcaList) == 0 {
					global.GVA_LOG.Error("len(vcaList) == 0")
					_ = msg.Ack(true)
					continue
				}

				accDB := vcaList[0]

				//string 转uint
				UID, _ := strconv.ParseUint(uid, 10, 64)
				orgTmp := utils2.GetSelfOrg(uint(UID))

				global.GVA_LOG.Info("【引导类】收到一条需要处理查询冷却状态的QN shop", zap.Any("ID", ID), zap.Any("MID", MID), zap.Any("money", money),
					zap.Any("markID", markID), zap.Any("waitAccYdKey", waitAccYdKey))

				//查一下订单是否超出账户限制
				var flag bool
				// 1. 查询该用户的余额是否充足
				var balance int
				err = global.GVA_DB.Model(&vbox.UserWallet{}).Select("IFNULL(sum(recharge), 0) as balance").
					Where("uid = ?", accDB.CreatedBy).Scan(&balance).Error

				//余额不足，则 log 一条
				if balance <= 0 {

					//入库操作记录
					flag = true

					msgX := fmt.Sprintf(global.BalanceNotEnough, accDB.AcId, accDB.AcAccount)

					global.GVA_LOG.Error("余额不足...", zap.Any("msg", msgX))
					err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
						Update("sys_status", 0).Error
				}

				// 2. 查询账号是否有超 金额限制，或者笔数限制

				// 2.1 日限制
				if accDB.DailyLimit > 0 {
					var dailySum int
					//// 获取今天的时间范围
					//startOfDay := time.Now().UTC().Truncate(24 * time.Hour)
					//endOfDay := startOfDay.Add(24 * time.Hour)
					//// 获取本地时区
					//loc, _ := time.LoadLocation("Asia/Shanghai") // 请替换为你实际使用的时区
					//startOfDay = startOfDay.In(loc)
					//endOfDay = endOfDay.In(loc)

					err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as dailySum").
						Where("ac_id = ?", accDB.AcId).
						Where("channel_code = ?", accDB.Cid).
						Where("order_status = ? AND created_at BETWEEN CURDATE() AND CURDATE() + INTERVAL 1 DAY - INTERVAL 1 SECOND", 1).Scan(&dailySum).Error

					if err != nil {
						global.GVA_LOG.Error("当前账号计算日消耗查mysql错误，直接丢了..." + err.Error())
						_ = msg.Reject(false)
						err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
							Update("sys_status", 0).Error
						continue
					}

					if dailySum > accDB.DailyLimit { // 如果日消费已经超了，不允许开启了，直接结束
						flag = true

						msgX := fmt.Sprintf(global.AccDailyLimitNotEnough, accDB.AcId, accDB.AcAccount, dailySum, accDB.DailyLimit)
						global.GVA_LOG.Error("当前账号日消耗已经超限...", zap.Any("msg", msgX), zap.Any("daily Sum", dailySum), zap.Any("daily limit", accDB.DailyLimit))
						err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
							Update("sys_status", 0).Error

					}
				}
				// 2.2 总限制
				if accDB.TotalLimit > 0 {

					var totalSum int

					err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as totalSum").
						Where("ac_id = ?", accDB.AcId).
						Where("order_status = ?", 1).Scan(&totalSum).Error

					if err != nil {
						global.GVA_LOG.Error("当前账号计算总消耗查mysql错误，直接丢了..." + err.Error())
					}

					if totalSum > accDB.TotalLimit { // 如果总消费已经超了，不允许开启了，直接结束
						flag = true

						//入库操作记录
						msgX := fmt.Sprintf(global.AccTotalLimitNotEnough, accDB.AcId, accDB.AcAccount, totalSum, accDB.TotalLimit)
						global.GVA_LOG.Error("当前账号总消耗已经超限...", zap.Any("msg", msgX), zap.Any("total Sum", totalSum), zap.Any("total limit", accDB.TotalLimit))

						err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
							Update("sys_status", 0).Error

						global.GVA_LOG.Info("当前账号总消耗已经超限额了，结束...", zap.Any("ac info", accDB))
					}
				}
				// 2.3 进单限制
				if accDB.InCntLimit > 0 {

					var count int64

					err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("ac_id = ? and order_status = ?", accDB.AcId, 1).Count(&count).Error

					if err != nil {
						global.GVA_LOG.Error("当前账号笔数消耗查mysql错误，直接丢了..." + err.Error())
					}

					if int(count) >= accDB.InCntLimit { // 如果笔数消费已经超了，不允许开启了，直接结束

						flag = true
						msgX := fmt.Sprintf(global.AccInCntLimitNotEnough, accDB.AcId, accDB.AcAccount, count, accDB.InCntLimit)

						global.GVA_LOG.Error("当前账号笔数消耗已经超限额...", zap.Any("msg", msgX), zap.Any("cnt", count), zap.Any("limit", accDB.InCntLimit))
						err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
							Update("sys_status", 0).Error
						global.GVA_LOG.Warn("当前账号笔数消耗已经超限额了，结束...", zap.Any("ac info", accDB))
					}
				}
				// 2.4 拉单限制
				if accDB.CountLimit > 0 {

					var count int64

					err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("ac_id = ?", accDB.AcId).Count(&count).Error

					if err != nil {
						global.GVA_LOG.Error("当前账号笔数消耗查mysql错误，直接丢了..." + err.Error())
					}

					if int(count) >= accDB.CountLimit { // 如果笔数消费已经超了，不允许开启了，直接结束

						flag = true
						msgX := fmt.Sprintf(global.AccCountLimitNotEnough, accDB.AcId, accDB.AcAccount, count, accDB.CountLimit)

						global.GVA_LOG.Error("当前账号笔数消耗已经超限额...", zap.Any("msg", msgX), zap.Any("cnt", count), zap.Any("limit", accDB.CountLimit))
						err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
							Update("sys_status", 0).Error
						global.GVA_LOG.Warn("当前账号笔数消耗已经超限额了，结束...", zap.Any("ac info", accDB))
					}
				}

				ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()

				if ttl <= 0 { //冷却结束，直接置为已用

					accKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgTmp[0], "5001", money)

					if flag || accDB.Status == 0 || accDB.SysStatus == 0 { // 表示超限了，删掉处理
						_ = global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
						global.GVA_LOG.Info("CD超限或关闭QN shop，删掉处理", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem), zap.Any("acc.Status", accDB.Status), zap.Any("acc.SysStatus", accDB.SysStatus))
					} else {
						global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
						global.GVA_LOG.Info("QN shop CD置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
					}

					_ = msg.Ack(true)
					continue
				} else { //仍然处于冷却状态，重新丢回ck check mq
					if flag || accDB.Status == 0 || accDB.SysStatus == 0 { // 表示超限或者当前账号已经关闭，删掉处理

						accKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgTmp[0], "5001", money)
						_ = global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
						global.GVA_LOG.Info("qn shop CD超限或关闭账号，删掉处理", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))

					} else {
						waitMsg := v

						ttlN := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
						err = chX.PublishWithDelay(QNShopCDCheckDelayedExchange, QNShopCDCheckDelayedRoutingKey, []byte(waitMsg), ttlN)
						global.GVA_LOG.Info("qn shop 还在冷却中，重新放回cd check mq", zap.Any("msg", waitMsg), zap.Any("ttlN", ttlN))
					}
				}

				if err != nil {
					_ = msg.Reject(false)
					continue
				}

				_ = msg.Ack(true)

			}
			wg.Done()
		}(i + 1)
	}

	// 等待所有消费者完成处理
	wg.Wait()
	//time.Sleep(time.Minute)
	global.GVA_LOG.Info("Mq qn shop cd Task 初始化搞定")

}
