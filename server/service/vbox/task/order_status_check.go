package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	"go.uber.org/zap"
	"strings"
	"sync"
)

// 订单等待查询
const (
	OrderStatusCheckDelayedExchange   = "vbox.order.status_delayed_exchange"
	OrderStatusCheckDelayedRoutingKey = "vbox.order.status_delayed_routing_key"
	OrderStatusCheckDelayedQueue      = "vbox.order.status_delayed_queue"
	OrderStatusCheckDeadRoutingKey    = "vbox.order.status_dead_routing_key"
	OrderStatusCheckDeadExchange      = "vbox.order.status_dead_exchange"
	OrderStatusCheckDeadQueue         = "vbox.order.status_dead_queue"
)

// OrderStatusCheckTask 处理订单10分钟后仍未匹配，切订单状态仍处于待支付的订单
func OrderStatusCheckTask() {

	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		//log.Fatalf("Failed to get connection from pool: %v", err)
		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
	ch, _ := conn.Channel()
	if err = ch.ExchangeDeclare(OrderStatusCheckDeadExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclare(OrderStatusCheckDeadQueue); err != nil {
		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueBind(OrderStatusCheckDeadQueue, OrderStatusCheckDeadRoutingKey, OrderStatusCheckDeadExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclareWithDelay(OrderStatusCheckDelayedQueue, OrderStatusCheckDeadExchange, OrderStatusCheckDeadRoutingKey); err != nil {
		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
	}
	if err = ch.ExchangeDeclare(OrderStatusCheckDelayedExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
	}
	if err = ch.QueueBind(OrderStatusCheckDelayedQueue, OrderStatusCheckDelayedRoutingKey, OrderStatusCheckDelayedExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 333:", zap.Any("err", err))
	}

	// 设置初始消费者数量
	consumerCount := 15
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
			deliveries, err := chX.Consume(OrderStatusCheckDeadQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("mq 消费者异常， err", zap.Error(err), zap.Any("queue", OrderStatusCheckDeadQueue))
			}

			for msg := range deliveries {
				//err = handler(msg.Body)
				v := string(msg.Body)

				split := strings.Split(v, "-")
				orderID := split[0]
				ID := split[1]

				global.GVA_LOG.Info("收到一条需要处理查询的订单，查看是否为长时间仍未匹配的订单", zap.Any("orderID", orderID), zap.Any("ID", ID))

				/*msgID := fmt.Sprintf(global.MsgFilterMem, msg.MessageId, orderID)
				// 检查消息是否已经被处理过
				exists, errR := global.GVA_REDIS.SIsMember(context.Background(), global.MsgFilterKey, msgID).Result()
				if errR != nil {
					global.GVA_LOG.Error("redis ex", zap.Error(errR))
				}

				if exists {
					// 消息已经被处理过，直接返回
					global.GVA_LOG.Info("消息已经被处理过", zap.Any("msgID", msgID))
					// 消息已经处理过，不再处理
					_ = msg.Ack(false)
					continue
				}
				// 将消息ID添加到已处理集合
				errR = global.GVA_REDIS.SAdd(context.Background(), global.MsgFilterKey, msgID).Err()
				if errR != nil {
					global.GVA_LOG.Error("redis ex", zap.Error(errR))
				}
				global.GVA_LOG.Info("消息首次被处理", zap.Any("msgID", msgID))*/

				var orderDB vbox.PayOrder
				if errQ := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).First(&orderDB, ID).Error; errQ != nil {
					global.GVA_LOG.Error("查库异常记录", zap.Any("err", errQ), zap.Any("orderID", orderID))
					_ = msg.Reject(false)
					continue
				}

				if orderDB.OrderStatus == 2 && orderDB.AcId == "" {
					if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", ID).Update("order_status", 0).Error; errDB != nil {
						global.GVA_LOG.Info("订单查库异常.", zap.Error(errDB))
						_ = msg.Reject(false)
						continue
					}
					global.GVA_LOG.Info("该订单长时间为匹配，且无账号匹配，置订单状态为匹配失败", zap.Any("orderID", orderID), zap.Any("ID", ID))
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
	global.GVA_LOG.Info("Mq order status check Task 初始化搞定")

}
