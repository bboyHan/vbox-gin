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
	consumerCount := 5
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

				global.GVA_LOG.Info("收到一条需要处理查询冷却状态的账号", zap.Any("info", v))

				split := strings.Split(v, ":")
				splitAcc := strings.Split(split[1], "_")
				acID := splitAcc[0]
				pcIDTmp := splitAcc[1]
				pcIDs := strings.Split(pcIDTmp, ",")

				global.GVA_LOG.Info("开始处理查询冷却状态的账号", zap.Any("acID", acID), zap.Any("pcIDs", pcIDs))

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
						if pcDB.ExpTime.Before(time.Now()) {
							global.GVA_REDIS.ZAdd(context.Background(), pcKey, redis.Z{Score: 0, Member: pcMem})
							global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id =?", pcDB.ID).Update("code_status", 2)
						} else { // 过期了，直接删除redis，并且状态置为失效 3
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
