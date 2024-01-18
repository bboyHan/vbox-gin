package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"go.uber.org/zap"
	"sync"
	"time"
)

// 订单等待查询
// - 订单确认查询
const (
	PayCodeDelayedExchange   = "vbox.order.pay_code_delayed_exchange"
	PayCodeDelayedRoutingKey = "vbox.order.pay_code_delayed_routing_key"
	PayCodeDelayedQueue      = "vbox.order.pay_code_delayed_queue"

	PayCodeDeadExchange   = "vbox.order.pay_code_dead_exchange"
	PayCodeDeadRoutingKey = "vbox.order.pay_code_dead_routing_key"
	PayCodeDeadQueue      = "vbox.order.pay_code_dead_queue"
)

// PayCodeExpCheck 预产过期检查
func PayCodeExpCheck() {
	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(PayCodeDeadExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(PayCodeDeadQueue); err != nil {
		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
	}
	if err := ch.QueueBind(PayCodeDeadQueue, PayCodeDeadRoutingKey, PayCodeDeadExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
	}
	if err := ch.QueueDeclareWithDelay(PayCodeDelayedQueue, PayCodeDeadExchange, PayCodeDeadRoutingKey); err != nil {
		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
	}
	if err := ch.ExchangeDeclare(PayCodeDelayedExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
	}
	if err := ch.QueueBind(PayCodeDelayedQueue, PayCodeDelayedRoutingKey, PayCodeDelayedExchange); err != nil {
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
			deliveries, err := ch.Consume(PayCodeDeadQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", PayCodeDeadQueue))
			}

			for msg := range deliveries {
				//err = handler(msg.Body)
				v := vbox.ChannelPayCode{}
				err := json.Unmarshal(msg.Body, &v)
				nowTime := time.Now()
				if err != nil {
					global.GVA_LOG.Error("Mq Pay Code Task...", zap.Error(err), zap.Any("错误消息体", msg.Body))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}
				global.GVA_LOG.Info("处理需要核验失效状态的预产码消息", zap.Any("预产码MID", v.Mid))

				//如果 v.expTime 比当前时间早，则表示已经过期了，先查一下数据库记录，如果记录的code_status仍为2，并把该记录在数据中更新为过期状态，code_status置为3
				// 先查一下数据库记录
				var pcDB vbox.ChannelPayCode
				err = global.GVA_DB.Unscoped().Where("mid =?", v.Mid).First(&pcDB).Error
				if err != nil {
					global.GVA_LOG.Error("查询预产码失败", zap.Error(err), zap.Any("预产码MID", v.Mid))
				}
				if pcDB.CodeStatus == 2 || pcDB.CodeStatus == 4 {
					global.GVA_LOG.Info("预产码过期了并且状态为2，需要更新", zap.Any("预产码MID", v.Mid), zap.Any("过期时间", v.ExpTime), zap.Any("当前时间", nowTime))
					v.CodeStatus = 3
					err = global.GVA_DB.Save(&v).Error
					if err != nil {
						global.GVA_LOG.Error("更新预产码状态失败", zap.Error(err), zap.Any("预产码", v))
					}

					// 把redis预产池里的预产码也删除掉
					orgTmp := utils2.GetSelfOrg(pcDB.CreatedBy)
					pcKey := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0],
						pcDB.Cid, pcDB.Money, pcDB.Operator, pcDB.Location)
					pcMem := fmt.Sprintf("%d", pcDB.ID) + "_" + pcDB.Mid + "_" + pcDB.AcAccount + "_" + pcDB.ImgContent
					global.GVA_REDIS.ZRem(context.Background(), pcKey, pcMem)

					global.GVA_LOG.Info("预产码已经过期或失效了，处理掉", zap.Any("预产码MID", v.Mid))
					_ = msg.Ack(true)
					continue
				} else if pcDB.CodeStatus == 1 {
					//如果 code_status = 1，则表示已经使用了，则直接跳过
					global.GVA_LOG.Info("预产码已经使用了，不需要处理", zap.Any("预产码MID", v.Mid))
					_ = msg.Ack(true)

					// 把redis预产池里的预产码也删除掉
					orgTmp := utils2.GetSelfOrg(pcDB.CreatedBy)
					pcKey := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0],
						pcDB.Cid, pcDB.Money, pcDB.Operator, pcDB.Location)
					pcMem := fmt.Sprintf("%d", pcDB.ID) + "_" + pcDB.Mid + "_" + pcDB.AcAccount + "_" + pcDB.ImgContent
					global.GVA_REDIS.ZRem(context.Background(), pcKey, pcMem)

					continue
				} else {
					//	其它状态，删redis
					// 把redis预产池里的预产码也删除掉
					orgTmp := utils2.GetSelfOrg(pcDB.CreatedBy)
					pcKey := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0],
						pcDB.Cid, pcDB.Money, pcDB.Operator, pcDB.Location)
					pcMem := fmt.Sprintf("%d", pcDB.ID) + "_" + pcDB.Mid + "_" + pcDB.AcAccount + "_" + pcDB.ImgContent
					global.GVA_REDIS.ZRem(context.Background(), pcKey, pcMem)

					global.GVA_LOG.Info("预产码已经过期或失效了，处理掉", zap.Any("预产码MID", v.Mid))
					_ = msg.Ack(true)
					continue
				}
				//
				//if err != nil {
				//	global.GVA_LOG.Error("Mq Pay Code Task...", zap.Error(err))
				//	_ = msg.Reject(false)
				//	continue
				//}
				//
				//_ = msg.Ack(true)
				//global.GVA_LOG.Info("核验完成", zap.Any("对应预产MID", v.Mid))
			}
			wg.Done()
		}(i + 1)
	}

	// 等待所有消费者完成处理
	wg.Wait()
	//time.Sleep(time.Minute)
	global.GVA_LOG.Info("MqOrderConfirmTask 初始化搞定")

}
