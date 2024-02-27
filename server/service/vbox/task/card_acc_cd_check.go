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
)

// 订单等待查询
const (
	CardAccCDCheckDelayedExchange   = "vbox.order.pool_card_acc_cd_delayed_exchange"
	CardAccCDCheckDelayedRoutingKey = "vbox.order.pool_card_acc_cd_delayed_routing_key"
	CardAccCDCheckDelayedQueue      = "vbox.order.pool_card_acc_cd_delayed_queue"
	CardAccCDCheckDeadRoutingKey    = "vbox.order.pool_card_acc_cd_dead_routing_key"
	CardAccCDCheckDeadExchange      = "vbox.order.pool_card_acc_cd_dead_exchange"
	CardAccCDCheckDeadQueue         = "vbox.order.pool_card_acc_cd_dead_queue"
)

// CardAccCDCheckTask 卡密-查单池账号-冷却检查
func CardAccCDCheckTask() {

	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		//log.Fatalf("Failed to get connection from pool: %v", err)
		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
	ch, _ := conn.Channel()
	if err = ch.ExchangeDeclare(CardAccCDCheckDeadExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclare(CardAccCDCheckDeadQueue); err != nil {
		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueBind(CardAccCDCheckDeadQueue, CardAccCDCheckDeadRoutingKey, CardAccCDCheckDeadExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclareWithDelay(CardAccCDCheckDelayedQueue, CardAccCDCheckDeadExchange, CardAccCDCheckDeadRoutingKey); err != nil {
		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
	}
	if err = ch.ExchangeDeclare(CardAccCDCheckDelayedExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
	}
	if err = ch.QueueBind(CardAccCDCheckDelayedQueue, CardAccCDCheckDelayedRoutingKey, CardAccCDCheckDelayedExchange); err != nil {
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

			connX, errX := mq.MQ.ConnPool.GetConnection()
			if errX != nil {
				//log.Fatalf("Failed to get connection from pool: %v", err)
				global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(errX))
			}
			defer mq.MQ.ConnPool.ReturnConnection(connX)
			chX, _ := connX.Channel()

			// 说明：执行查单回调处理
			deliveries, errC := chX.Consume(CardAccCDCheckDeadQueue, "", false, false, false, false, nil)
			if errC != nil {
				global.GVA_LOG.Error("mq 消费者异常， err", zap.Error(errC), zap.Any("queue", CardAccCDCheckDeadQueue))
			}

			for msg := range deliveries {
				//err = handler(msg.Body)
				v := string(msg.Body)

				split := strings.Split(v, "-")
				waitAccYdKey := split[0]
				waitAccMem := split[1]
				var accInfo []string
				if strings.Contains(waitAccMem, ",") {
					accInfo = strings.Split(waitAccMem, ",")
				} else {
					accInfo = strings.Split(waitAccMem, "_")
				}
				ID := accInfo[0]
				acID := accInfo[1]
				acAccount := accInfo[2]
				var money string
				if len(accInfo) == 4 {
					money = accInfo[3]
				}

				global.GVA_LOG.Info("【卡密类】收到一条需要处理查单池冷却状态的账号", zap.Any("acID", acID), zap.Any("acAccount", acAccount), zap.Any("money", money), zap.Any("waitAccYdKey", waitAccYdKey))

				var accDB vbox.ChannelCardAcc
				if errQ := global.GVA_DB.Debug().Unscoped().Where("id = ?", ID).Find(&accDB).Error; errQ != nil {
					global.GVA_LOG.Error("查找异常", zap.Error(errQ), zap.Any("ID", ID), zap.Any("acID", acID), zap.Any("acAccount", acAccount), zap.Any("money", money))
					_ = msg.Reject(false)
					continue
				}

				ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()

				cid := accDB.Cid
				if ttl <= 0 { //冷却结束，直接置为已用

					orgTmp := utils2.GetSelfOrg(accDB.CreatedBy)

					if global.ECContains(cid) { // 引导
						accKey := fmt.Sprintf(global.ChanOrgECPoolAccZSet, orgTmp[0], cid)

						if accDB.Status == 0 || accDB.SysStatus == 0 { // 表示超限了，删掉处理
							_ = global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("关闭账号，删掉处理", zap.Any("cardAccKey", accKey), zap.Any("card waitAccMem", waitAccMem), zap.Any("acc.Status", accDB.Status), zap.Any("acc.SysStatus", accDB.SysStatus))
						} else {
							global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
							global.GVA_LOG.Info("CD置为可用", zap.Any("card accKey", accKey), zap.Any("card waitAccMem", waitAccMem))
						}

					} else { // wx qb
						global.GVA_LOG.Info("非卡密类，无需处理", zap.Any("cid", cid))
					}

					_ = msg.Ack(true)
					continue
				} else { //仍然处于冷却状态，重新丢回ck check mq
					if accDB.Status == 0 || accDB.SysStatus == 0 { // 表示当前账号已经关闭，删掉处理

						orgTmp := utils2.GetSelfOrg(accDB.CreatedBy)

						if global.ECContains(cid) { // 引导
							accKey := fmt.Sprintf(global.ChanOrgECPoolAccZSet, orgTmp[0], cid)
							_ = global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("card pool 关闭账号，删掉处理", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))

						} else { // wx qb
							global.GVA_LOG.Info("非卡密类，无需处理", zap.Any("cid", cid))
						}

					} else {

						waitMsg := v

						ttlN := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
						err = chX.PublishWithDelay(CardAccCDCheckDelayedExchange, CardAccCDCheckDelayedRoutingKey, []byte(waitMsg), ttlN)
						global.GVA_LOG.Info("card pool acc 还在冷却中，重新放回cd check mq", zap.Any("msg", waitMsg), zap.Any("ttlN", ttlN))
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
	global.GVA_LOG.Info("Mq card pool acc cd Task 初始化搞定")

}
