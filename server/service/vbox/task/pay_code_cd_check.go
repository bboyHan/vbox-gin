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
	PayCodeCDCheckDelayedExchange   = "vbox.order.pay_code_cd_delayed_exchange"
	PayCodeCDCheckDelayedRoutingKey = "vbox.order.pay_code_cd_delayed_routing_key"
	PayCodeCDCheckDelayedQueue      = "vbox.order.pay_code_cd_delayed_queue"
	PayCodeCDCheckDeadRoutingKey    = "vbox.order.pay_code_cd_dead_routing_key"
	PayCodeCDCheckDeadExchange      = "vbox.order.pay_code_cd_dead_exchange"
	PayCodeCDCheckDeadQueue         = "vbox.order.pay_code_cd_dead_queue"
)

// PayCodeCDCheckTask 预产冷却检查
func PayCodeCDCheckTask() {

	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		//log.Fatalf("Failed to get connection from pool: %v", err)
		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
	ch, _ := conn.Channel()
	if err = ch.ExchangeDeclare(PayCodeCDCheckDeadExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclare(PayCodeCDCheckDeadQueue); err != nil {
		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueBind(PayCodeCDCheckDeadQueue, PayCodeCDCheckDeadRoutingKey, PayCodeCDCheckDeadExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclareWithDelay(PayCodeCDCheckDelayedQueue, PayCodeCDCheckDeadExchange, PayCodeCDCheckDeadRoutingKey); err != nil {
		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
	}
	if err = ch.ExchangeDeclare(PayCodeCDCheckDelayedExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
	}
	if err = ch.QueueBind(PayCodeCDCheckDelayedQueue, PayCodeCDCheckDelayedRoutingKey, PayCodeCDCheckDelayedExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 333:", zap.Any("err", err))
	}

	// 设置初始消费者数量
	consumerCount := 10
	// 使用 WaitGroup 来等待所有消费者完成处理
	var wg sync.WaitGroup
	wg.Add(consumerCount)
	// 启动多个消费者
	for i := 0; i < consumerCount; i++ {
		go func(consumerID int) {
			// 说明：执行查单回调处理
			deliveries, err := ch.Consume(PayCodeCDCheckDeadQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("mq 消费者异常， err", zap.Error(err), zap.Any("queue", PayCodeCDCheckDeadQueue))
			}

			for msg := range deliveries {
				//err = handler(msg.Body)
				v := string(msg.Body)

				global.GVA_LOG.Info("【预产类】收到一条需要处理查询冷却状态的账号", zap.Any("info", v))

				split := strings.Split(v, "-")
				//waitAccPcKey - "vb_acc_waiting_pc:acid_%s"
				waitAccPcKey := split[0]
				acID := strings.Split(strings.Split(waitAccPcKey, ":")[1], "_")[1]
				pcIDTmp := split[1]
				pcIDs := strings.Split(pcIDTmp, ",")

				global.GVA_LOG.Info("开始处理查询冷却状态的账号", zap.Any("acID", acID), zap.Any("pcIDs", pcIDs))

				var accDB vbox.ChannelAccount
				global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("ac_id = ?", acID).First(&accDB)

				var pcDBList []vbox.ChannelPayCode
				global.GVA_DB.Debug().Model(&vbox.ChannelPayCode{}).Where("id in ? ", pcIDs).Find(&pcDBList)

				for _, pcDB := range pcDBList {
					orgTmp := utils2.GetSelfOrg(pcDB.CreatedBy)
					pcKey := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0], pcDB.Cid, pcDB.Money, pcDB.Operator, pcDB.Location)
					pcMem := fmt.Sprintf("%d", pcDB.ID) + "_" + pcDB.Mid + "_" + pcDB.AcAccount + "_" + pcDB.ImgContent

					if pcDB.CodeStatus == 1 {
						global.GVA_LOG.Info("已使用，置为已用", zap.Any("pcID", pcDB.ID))
						global.GVA_REDIS.ZAdd(context.Background(), pcKey, redis.Z{Score: 1, Member: pcMem})
					} else if pcDB.CodeStatus == 4 {
						global.GVA_LOG.Info("冷却结束，开始处理恢复", zap.Any("pcID", pcDB.ID))

						//查一下订单是否超出账户限制
						var flag bool
						// 1. 查询该用户的余额是否充足
						var balance int
						err = global.GVA_DB.Model(&vbox.UserWallet{}).Select("IFNULL(sum(recharge), 0) as balance").
							Where("uid = ?", pcDB.CreatedBy).Scan(&balance).Error

						if balance <= 0 { //余额不足，则 log 一条
							//入库操作记录
							flag = true

							msgX := fmt.Sprintf(global.BalanceNotEnough, pcDB.AcId, pcDB.AcAccount)

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

								msg := fmt.Sprintf(global.AccDailyLimitNotEnough, accDB.AcId, accDB.AcAccount)
								global.GVA_LOG.Error("当前账号日消耗已经超限...", zap.Any("msg", msg))
								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
									Update("sys_status", 0).Error

							}
						}
						// 2.2 总限制
						if accDB.TotalLimit > 0 {

							var totalSum int

							err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as totalSum").
								Where("ac_id = ?", accDB.AcId).
								Where("channel_code = ?", accDB.Cid).
								Where("order_status = ?", 1).Scan(&totalSum).Error

							if err != nil {
								global.GVA_LOG.Error("当前账号计算总消耗查mysql错误，直接丢了..." + err.Error())
							}

							if totalSum > accDB.TotalLimit { // 如果总消费已经超了，不允许开启了，直接结束
								flag = true

								//入库操作记录
								msgX := fmt.Sprintf(global.AccTotalLimitNotEnough, accDB.AcId, accDB.AcAccount)
								global.GVA_LOG.Error("当前账号总消耗已经超限...", zap.Any("msg", msgX))

								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
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
								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", accDB.ID).
									Update("sys_status", 0).Error
								global.GVA_LOG.Warn("当前账号笔数消耗已经超限额了，结束...", zap.Any("ac info", accDB))
							}
						}

						if pcDB.ExpTime.After(time.Now()) && !flag { // 设置的过期时间比当前时间晚，表示还可使用
							global.GVA_REDIS.ZAdd(context.Background(), pcKey, redis.Z{Score: 0, Member: pcMem})
							global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id =?", pcDB.ID).Update("code_status", 2)
						} else { // 超限额，直接结束，并且状态置为失效 3
							global.GVA_REDIS.ZRem(context.Background(), pcKey, pcMem)
							global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id =?", pcDB.ID).Update("code_status", 3)
						}

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
