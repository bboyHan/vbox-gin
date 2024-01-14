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
	"strings"
	"sync"
	"time"
)

// 订单等待查询
const (
	AccCDCheckDelayedExchange   = "vbox.order.acc_cd_delayed_exchange"
	AccCDCheckDelayedRoutingKey = "vbox.order.acc_cd_delayed_routing_key"
	AccCDCheckDelayedQueue      = "vbox.order.acc_cd_delayed_queue"
	AccCDCheckDeadRoutingKey    = "vbox.order.acc_cd_dead_routing_key"
	AccCDCheckDeadExchange      = "vbox.order.acc_cd_dead_exchange"
	AccCDCheckDeadQueue         = "vbox.order.acc_cd_dead_queue"
)

// AccCDCheckTask 引导类-通道账号-冷却检查
func AccCDCheckTask() {

	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		//log.Fatalf("Failed to get connection from pool: %v", err)
		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
	ch, _ := conn.Channel()
	if err = ch.ExchangeDeclare(AccCDCheckDeadExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclare(AccCDCheckDeadQueue); err != nil {
		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueBind(AccCDCheckDeadQueue, AccCDCheckDeadRoutingKey, AccCDCheckDeadExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclareWithDelay(AccCDCheckDelayedQueue, AccCDCheckDeadExchange, AccCDCheckDeadRoutingKey); err != nil {
		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
	}
	if err = ch.ExchangeDeclare(AccCDCheckDelayedExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
	}
	if err = ch.QueueBind(AccCDCheckDelayedQueue, AccCDCheckDelayedRoutingKey, AccCDCheckDelayedExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 333:", zap.Any("err", err))
	}

	// 设置初始消费者数量
	consumerCount := 5
	// 使用 WaitGroup 来等待所有消费者完成处理
	var wg sync.WaitGroup
	wg.Add(consumerCount)
	// 启动多个消费者
	for i := 0; i < consumerCount; i++ {
		go func(consumerID int) {
			// 说明：执行查单回调处理
			deliveries, err := ch.Consume(AccCDCheckDeadQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("mq 消费者异常， err", zap.Error(err), zap.Any("queue", AccCDCheckDeadQueue))
			}

			for msg := range deliveries {
				//err = handler(msg.Body)
				v := string(msg.Body)

				global.GVA_LOG.Info("【预产类】收到一条需要处理查询冷却状态的账号", zap.Any("info", v))

				split := strings.Split(v, "-")
				//waitAccYdKey - "vb_acc_waiting_yd:acid_%s:money_%v"
				waitAccYdKey := split[0]
				waitAccMem := split[1]
				accInfo := strings.Split(waitAccMem, "_")
				ID := accInfo[0]
				acID := accInfo[1]
				acAccount := accInfo[2]
				money := accInfo[3]

				global.GVA_LOG.Info("【引导类】收到一条需要处理查询冷却状态的账号", zap.Any("info", v))

				var accDB vbox.ChannelAccount
				if errQ := global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", ID).First(&accDB); errQ != nil {
					global.GVA_LOG.Error("没找到记录", zap.Any("err", errQ), zap.Any("ID", ID), zap.Any("acID", acID), zap.Any("acAccount", acAccount), zap.Any("money", money))
					_ = msg.Reject(false)
					continue
				}

				//检查限额情况
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
					err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
						Update("sys_status", 0).Error
				}

				// 2. 查询账号是否有超 金额限制，或者笔数限制

				// 2.1 日限制
				if accDB.DailyLimit > 0 {
					var dailySum int
					// 获取今天的时间范围
					startOfDay := time.Now().UTC().Truncate(24 * time.Hour)
					endOfDay := startOfDay.Add(24 * time.Hour)

					err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("sum(money) as dailySum").
						Where("ac_id = ?", accDB.AcId).
						Where("order_status = ? AND created_at BETWEEN ? AND ?", 1, startOfDay, endOfDay).Scan(&dailySum).Error

					if err != nil {
						global.GVA_LOG.Error("当前账号计算日消耗查mysql错误，直接丢了..." + err.Error())
						_ = msg.Reject(false)
						err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
							Update("sys_status", 0).Error
						continue
					}

					if dailySum > accDB.DailyLimit { // 如果日消费已经超了，不允许开启了，直接结束
						flag = true

						msg := fmt.Sprintf(global.AccDailyLimitNotEnough, accDB.AcId, accDB.AcAccount)
						global.GVA_LOG.Error("当前账号日消耗已经超限...", zap.Any("msg", msg))
						err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
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
						msgX := fmt.Sprintf(global.AccTotalLimitNotEnough, accDB.AcId, accDB.AcAccount)
						global.GVA_LOG.Error("当前账号总消耗已经超限...", zap.Any("msg", msgX))

						err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
							Update("sys_status", 0).Error

						global.GVA_LOG.Info("当前账号总消耗已经超限额了，结束...", zap.Any("ac info", accDB))
					}
				}
				// 2.3 笔数限制
				if accDB.CountLimit > 0 {

					var count int64

					err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("ac_id = ? and order_status = ?", accDB.AcId, 1).Count(&count).Error

					if err != nil {
						global.GVA_LOG.Error("当前账号笔数消耗查mysql错误，直接丢了..." + err.Error())
					}

					if int(count) >= accDB.CountLimit { // 如果笔数消费已经超了，不允许开启了，直接结束

						flag = true
						msgX := fmt.Sprintf(global.AccCountLimitNotEnough, accDB.AcId, accDB.AcAccount)

						global.GVA_LOG.Error("当前账号笔数消耗已经超限额...", zap.Any("msg", msgX))
						err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
							Update("sys_status", 0).Error
						global.GVA_LOG.Warn("当前账号笔数消耗已经超限额了，结束...", zap.Any("ac info", accDB))
					}
				}

				ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()

				if ttl <= 0 { //没进入冷却状态，直接置为已用
					orgTmp := utils2.GetSelfOrg(accDB.CreatedBy)

					accKey := fmt.Sprintf(global.ChanOrgAccZSet, orgTmp[0], accDB.Cid, money)

					if flag { // 表示超限了，删掉处理
						_ = global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
						global.GVA_LOG.Info("超限了，删掉处理", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
					} else {
						global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
						global.GVA_LOG.Info("置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
					}
					_ = msg.Ack(true)
					continue
				} else { //仍然处于冷却状态，重新丢回ck check mq
					if flag { // 表示超限了，删掉处理
						orgTmp := utils2.GetSelfOrg(accDB.CreatedBy)
						accKey := fmt.Sprintf(global.ChanOrgAccZSet, orgTmp[0], accDB.Cid, money)
						_ = global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
						global.GVA_LOG.Info("超限了，删掉处理", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))

					} else {
						waitMsg := v
						err = ch.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), ttl)
						global.GVA_LOG.Info("还在冷却中，重新放回ck check mq", zap.Any("ttl", ttl))
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
	global.GVA_LOG.Info("Mq PayCodeCDCheck Task 初始化搞定")

}