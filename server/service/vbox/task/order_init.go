package task

import (
	"encoding/json"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"sync"
	"time"
)

// 订单等待查询
// - 订单确认查询
const (
	OrderWaitExchange = "vbox.order.waiting_exchange"
	OrderWaitQueue    = "vbox.order.waiting_queue"
	OrderWaitKey      = "vbox.order.waiting"

	OrderConfirmDelayedExchange   = "vbox.order.confirm_delayed_exchange"
	OrderConfirmDelayedRoutingKey = "vbox.order.confirm_delayed_routing_key"
	OrderConfirmDelayedQueue      = "vbox.order.confirm_delayed_queue"
	OrderConfirmDeadRoutingKey    = "vbox.order.confirm_dead_routing_key"
	OrderConfirmDeadExchange      = "vbox.order.confirm_dead_exchange"
	OrderConfirmDeadQueue         = "vbox.order.confirm_dead_queue"
)

// OrderWaitingTask 订单入库匹配
func OrderWaitingTask() {

	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 创建 订单初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(OrderWaitExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 111:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(OrderWaitQueue); err != nil {
		global.GVA_LOG.Error("create queue err 111:", zap.Any("err", err))
	}
	if err := ch.QueueBind(OrderWaitQueue, OrderWaitKey, OrderWaitExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 111:", zap.Any("err", err))
	}

	// 设置初始消费者数量
	consumerCount := 20
	// 使用 WaitGroup 来等待所有消费者完成处理
	var wg sync.WaitGroup
	wg.Add(consumerCount)

	// 启动多个消费者
	for i := 0; i < consumerCount; i++ {
		go func(consumerID int) {
			// 说明：执行账号匹配
			deliveries, err := ch.Consume(OrderWaitQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", OrderWaitQueue))
			}

			for msg := range deliveries {
				//v := &map[string]interface{}{}
				//err := json.Unmarshal(msg.Body, v)
				//global.GVA_LOG.Info(fmt.Sprintf("%v", v))

				v := &vbox.PayOrder{}
				err := json.Unmarshal(msg.Body, v)
				if err != nil {
					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}

				//1. 筛选匹配是哪个产品
				var vpa vbox.PayAccount
				err = global.GVA_DB.Table("vbox_pay_account").
					Where("p_account = ?", v.PAccount).First(&vpa).Error
				if err != nil {
					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}

				//2. 查询产品对应的账号池是否有可用账号
				var total int64 = 0
				idList := utils2.GetDeepUserIDs(vpa.Uid)

				db := global.GVA_DB.Model(&vbox.ChannelAccount{}).Table("vbox_channel_account").
					Where("uid in (?)", idList).Count(&total)

				limit, offset := utils.RandSize2DB(int(total), 20)
				var vcas []vbox.ChannelAccount
				err = db.Where("status = ? and sys_status = ?", 1, 1).Where("cid = ?", v.ChannelCode).
					Where("uid in (?)", idList).Limit(limit).Offset(offset).
					Find(&vcas).Error
				if err != nil || len(vcas) == 0 {
					if len(vcas) == 0 {
						err = errors.New("库存不足！ 请联系对接人。")
						//TODO
						// 如果库存不足，丢弃消息并且更新订单（要么重试、要么订单更新为异常单）
						global.GVA_LOG.Error("库存不足..." + err.Error())

						_ = msg.Reject(false)
						continue
					}

					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}

				vca := vcas[rand.Intn(len(vcas))]
				marshal, err := json.Marshal(v)
				if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.ID).Update("uid", vca.Uid).Update("ac_id", vca.AcId).Error; err != nil {
					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}

				//3. 匹配账号后，更新订单信息（账号信息，订单支付链接处理）
				err = ch.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 1*time.Minute)

				global.GVA_LOG.Info("匹配到账号了，发一个准备查单的消息 : ", zap.Any("对应单号", v.OrderId))

				if err != nil {
					_ = msg.Reject(true)
					continue
				}
				_ = msg.Ack(true)
			}
			wg.Done()
		}(i + 1)
	}
	global.GVA_LOG.Info("Vbox OrderWaitingTask 初始化搞定")
	// 等待所有消费者完成处理
	wg.Wait()
}

// OrderConfirmTask 处理查单回调
func OrderConfirmTask() {
	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		//log.Fatalf("Failed to get connection from pool: %v", err)
		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(OrderConfirmDeadExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(OrderConfirmDeadQueue); err != nil {
		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
	}
	if err := ch.QueueBind(OrderConfirmDeadQueue, OrderConfirmDeadRoutingKey, OrderConfirmDeadExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
	}
	if err := ch.QueueDeclareWithDelay(OrderConfirmDelayedQueue, OrderConfirmDeadExchange, OrderConfirmDeadRoutingKey); err != nil {
		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
	}
	if err := ch.ExchangeDeclare(OrderConfirmDelayedExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
	}
	if err := ch.QueueBind(OrderConfirmDelayedQueue, OrderConfirmDelayedRoutingKey, OrderConfirmDelayedExchange); err != nil {
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
			// 说明：执行查单回调处理
			deliveries, err := ch.Consume(OrderConfirmDeadQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", OrderConfirmDeadQueue))
			}

			for msg := range deliveries {
				//err = handler(msg.Body)
				v := &vbox.PayOrder{}
				err := json.Unmarshal(msg.Body, v)
				if err != nil {
					global.GVA_LOG.Error("MqOrderConfirmTask...", zap.Error(err), zap.Any("错误消息体", msg.Body))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}
				global.GVA_LOG.Info("我收到延迟消息，consume msg :", zap.Any("orderId", v.OrderId))

				//1. 筛选匹配是哪个产品

				//2. 查询订单（账号）的充值情况

				//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
				if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.ID).Update("order_status", 1).Error; err != nil {
					global.GVA_LOG.Error("更新订单", zap.Error(err))
				}
				global.GVA_LOG.Info("订单支付了，我更新一下状态 ", zap.Any("orderId", v.OrderId))

				if err != nil {
					_ = msg.Reject(true)
					continue
				}
				_ = msg.Ack(false)
			}
			wg.Done()
		}(i + 1)
	}

	// 等待所有消费者完成处理
	wg.Wait()
	//time.Sleep(time.Minute)
	global.GVA_LOG.Info("MqOrderConfirmTask 初始化搞定")

}
