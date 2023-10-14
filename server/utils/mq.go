package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"log"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var (
	orderWaitExchange = "vbox.order.waiting_exchange"
	orderWaitQueue    = "vbox.order.waiting_queue"
	orderWaitKey      = "vbox.order.waiting"

	orderConfirmDelayedExchange   = "vbox.order.confirm_delayed_exchange"
	orderConfirmDelayedRoutingKey = "vbox.order.confirm_delayed_routing_key"
	orderConfirmDelayedQueue      = "vbox.order.confirm_delayed_queue"
	orderConfirmDeadRoutingKey    = "vbox.order.confirm_dead_routing_key"
	orderConfirmDeadExchange      = "vbox.order.confirm_dead_exchange"
	orderConfirmDeadQueue         = "vbox.order.confirm_dead_queue"
)
var delay = global.GVA_CONFIG.RabbitMQ.RetryReConnDelay
var ConnPool *MqConnPool

type MqConnPool struct {
	pool  chan *Connection
	mutex sync.Mutex
}

func Init() (err error) {
	mqCfg := global.GVA_CONFIG.RabbitMQ

	ConnPool = &MqConnPool{
		pool: make(chan *Connection, mqCfg.PoolSize),
	}

	for i := 0; i < mqCfg.PoolSize; i++ {
		conn, err := Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
			mqCfg.Username,
			mqCfg.Password,
			mqCfg.Addr,
			mqCfg.Port))
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}
		ConnPool.pool <- conn
	}
	return err
}

func (cp *MqConnPool) GetConnection() (*Connection, error) {
	cp.mutex.Lock()
	defer cp.mutex.Unlock()

	select {
	case conn := <-cp.pool: // 从连接池中取出连接
		return conn, nil
	default:
		// 如果连接池已满，则等待可用连接
		return <-cp.pool, nil
	}
}

func (cp *MqConnPool) ReturnConnection(conn *Connection) {
	cp.pool <- conn // 归还连接到连接池
}

func MqOrderWaitingTask() {
	// 示例：发送消息
	conn, err := ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer ConnPool.ReturnConnection(conn)

	// ------------- 创建 订单初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(orderWaitExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 111:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(orderWaitQueue); err != nil {
		global.GVA_LOG.Error("create queue err 111:", zap.Any("err", err))
	}
	if err := ch.QueueBind(orderWaitQueue, orderWaitKey, orderWaitExchange); err != nil {
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
			deliveries, err := ch.Consume(orderWaitQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", orderWaitQueue))
			}

			for msg := range deliveries {
				v := &vbox.VboxPayOrder{}
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
				userList, tot, err := GetOwnerUserIdsList(vpa.Uid)
				var idList []int
				for _, user := range userList {
					idList = append(idList, int(user.ID))
				}
				if err != nil || tot == 0 {
					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}
				db := global.GVA_DB.Model(&vbox.ChannelAccount{}).Table("vbox_channel_account").
					Where("uid in (?)", idList).Count(&total)

				limit, offset := RandSize2DB(int(total), 20)
				var vcas []vbox.ChannelAccount
				err = db.Where("status = ? and sys_status = ?", 1, 1).Where("cid = ?", v.ChannelCode).
					Where("uid in (?)", idList).Limit(limit).Offset(offset).
					Find(&vcas).Error
				if err != nil || len(vcas) == 0 {
					if len(vcas) == 0 {
						err = errors.New("库存不足！ 请联系对接人。")
					}
					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}

				vca := vcas[rand.Intn(len(vcas))]
				marshal, err := json.Marshal(v)
				if err := global.GVA_DB.Model(&vbox.VboxPayOrder{}).Where("id = ?", v.ID).
					Update("uid", vca.Uid).Update("ac_id", vca.AcId).
					Error; err != nil {
					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}

				//3. 匹配账号后，更新订单信息（账号信息，订单支付链接处理）
				err = ch.PublishWithDelay(orderConfirmDelayedExchange, orderConfirmDelayedRoutingKey, marshal, 1*time.Minute)

				global.GVA_LOG.Info("匹配到账号了，发一个准备查单的消息 : ", zap.Any("对应单号", v.OrderId))

				if err != nil {
					_ = msg.Reject(true)
					continue
				}
				_ = msg.Ack(false)
			}
			wg.Done()

			/*// 说明：执行账号匹配
			if err := NewConsumer(orderWaitQueue, func(body []byte) error {
				//s := string(body)
				v := &vbox.VboxPayOrder{}
				err := json.Unmarshal(body, v)
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
				userList, tot, err := GetOwnerUserIdsList(vpa.Uid)
				var idList []int
				for _, user := range userList {
					idList = append(idList, int(user.ID))
				}
				if err != nil || tot == 0 {
					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}
				db := global.GVA_DB.Model(&vbox.ChannelAccount{}).Table("vbox_channel_account").
					Where("uid in (?)", idList).Count(&total)

				limit, offset := RandSize2DB(int(total), 20)
				var vcas []vbox.ChannelAccount
				err = db.Where("status = ? and sys_status = ?", 1, 1).Where("cid = ?", v.ChannelCode).
					Where("uid in (?)", idList).Limit(limit).Offset(offset).
					Find(&vcas).Error
				if err != nil || len(vcas) == 0 {
					if len(vcas) == 0 {
						err = errors.New("库存不足！ 请联系对接人。")
					}
					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}

				vca := vcas[rand.Intn(len(vcas))]
				marshal, err := json.Marshal(v)
				if err := global.GVA_DB.Model(&vbox.VboxPayOrder{}).Where("id = ?", v.ID).
					Update("uid", vca.Uid).Update("ac_id", vca.AcId).
					Error; err != nil {
					global.GVA_LOG.Info("MqOrderWaitingTask..." + err.Error())
				}

				//3. 匹配账号后，更新订单信息（账号信息，订单支付链接处理）
				err = ch.PublishWithDelay(orderConfirmDelayedExchange, orderConfirmDelayedRoutingKey, marshal, 1*time.Minute)

				log.Printf("匹配到账号了，发一个准备查单的消息 : %v", v.OrderId)
				return nil
			}); err != nil {
				global.GVA_LOG.Info("orderWaitQueue", zap.Any("consume err", err))
			}*/
		}(i + 1)
	}
	global.GVA_LOG.Info("MqOrderWaitingTask 初始化搞定")
	// 等待所有消费者完成处理
	wg.Wait()
}

func MqOrderWaitingTask2() {
	// 示例：发送消息
	conn, err := ConnPool.GetConnection()
	if err != nil {
		//log.Fatalf("Failed to get connection from pool: %v", err)
		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
	}
	defer ConnPool.ReturnConnection(conn)

	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(orderConfirmDeadExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(orderConfirmDeadQueue); err != nil {
		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
	}
	if err := ch.QueueBind(orderConfirmDeadQueue, orderConfirmDeadRoutingKey, orderConfirmDeadExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
	}
	if err := ch.QueueDeclareWithDelay(orderConfirmDelayedQueue, orderConfirmDeadExchange, orderConfirmDeadRoutingKey); err != nil {
		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
	}
	if err := ch.ExchangeDeclare(orderConfirmDelayedExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
	}
	if err := ch.QueueBind(orderConfirmDelayedQueue, orderConfirmDelayedRoutingKey, orderConfirmDelayedExchange); err != nil {
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
			deliveries, err := ch.Consume(orderConfirmDeadQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", orderConfirmDeadQueue))
			}

			for msg := range deliveries {
				//err = handler(msg.Body)
				v := &vbox.VboxPayOrder{}
				err := json.Unmarshal(msg.Body, v)
				if err != nil {
					global.GVA_LOG.Error("MqOrderConfirmTask...", zap.Error(err))
				}
				global.GVA_LOG.Info("我收到延迟消息，consume msg :", zap.Any("orderId", v.OrderId))

				//1. 筛选匹配是哪个产品

				//2. 查询订单（账号）的充值情况

				//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
				if err := global.GVA_DB.Model(&vbox.VboxPayOrder{}).Where("id = ?", v.ID).Update("order_status", 1).Error; err != nil {
					global.GVA_LOG.Error("", zap.Error(err))
				}
				global.GVA_LOG.Info("订单支付了，我更新一下状态 ", zap.Any("orderId", v.OrderId))

				if err != nil {
					_ = msg.Reject(true)
					continue
				}
				_ = msg.Ack(false)
			}
			wg.Done()
			/*if err := NewConsumer(orderConfirmDeadQueue, func(body []byte) error {
				//s := string(body)
				v := &vbox.VboxPayOrder{}
				err := json.Unmarshal(body, v)
				if err != nil {
					global.GVA_LOG.Error("MqOrderConfirmTask...", zap.Error(err))
				}
				global.GVA_LOG.Info("我收到延迟消息，consume msg :", zap.Any("orderId", v.OrderId))

				//1. 筛选匹配是哪个产品

				//2. 查询订单（账号）的充值情况

				//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
				if err := global.GVA_DB.Model(&vbox.VboxPayOrder{}).Where("id = ?", v.ID).Update("order_status", 1).Error; err != nil {
					return err
				}
				global.GVA_LOG.Info("订单支付了，我更新一下状态 ", zap.Any("orderId", v.OrderId))

				return nil
			}); err != nil {
				global.GVA_LOG.Error("", zap.Error(err))

			}*/
		}(i + 1)
	}

	// 等待所有消费者完成处理
	wg.Wait()
	//time.Sleep(time.Minute)
	global.GVA_LOG.Info("MqOrderWaitingTask2 初始化搞定")

}

// Connection amqp.Connection wrapper
type Connection struct {
	*amqp.Connection
}

// Channel amqp.Channel wapper
type Channel struct {
	*amqp.Channel
	closed int32
}

// IsClosed indicate closed by developer
func (ch *Channel) IsClosed() bool {
	return atomic.LoadInt32(&ch.closed) == 1
}

// Close ensure closed flag set
func (ch *Channel) Close() error {
	if ch.IsClosed() {
		return amqp.ErrClosed
	}

	atomic.StoreInt32(&ch.closed, 1)

	return ch.Channel.Close()
}

// ExchangeDeclare 创建交换机.
func (ch *Channel) ExchangeDeclare(name string, kind string) (err error) {
	return ch.Channel.ExchangeDeclare(name, kind, true, false, false, false, nil)
}

// Publish 发布消息.
func (ch *Channel) Publish(exchange, key string, body []byte) (err error) {
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(context.Background(), exchange, key, false, false,
		//amqp.Publishing{ContentType: "text/plain", Body: body})
		amqp.Publishing{ContentType: "application/json", Body: body})
	return err
}

// PublishWithDelay 发布延迟消息.
func (ch *Channel) PublishWithDelay(exchange, key string, body []byte, timer time.Duration) (err error) {
	_, err = ch.Channel.PublishWithDeferredConfirmWithContext(context.Background(), exchange, key, false, false,
		amqp.Publishing{ContentType: "application/json", Body: body, Expiration: fmt.Sprintf("%d", timer.Milliseconds())})
	//amqp.Publishing{ContentType: "text/plain", Body: body, Expiration: fmt.Sprintf("%d", timer.Milliseconds())})
	return err
}

// QueueDeclare 创建队列.
func (ch *Channel) QueueDeclare(name string) (err error) {
	_, err = ch.Channel.QueueDeclare(name, true, false, false, false, nil)
	return
}

// QueueDeclareWithDelay 创建延迟队列.
func (ch *Channel) QueueDeclareWithDelay(name, exchange, key string) (err error) {
	_, err = ch.Channel.QueueDeclare(name, true, false, false, false, amqp.Table{
		"x-dead-letter-exchange":    exchange,
		"x-dead-letter-routing-key": key,
	})
	return
}

// QueueBind 绑定队列.
func (ch *Channel) QueueBind(name, key, exchange string) (err error) {
	return ch.Channel.QueueBind(name, key, exchange, false, nil)
}

// NewConsumer 实例化一个消费者, 会单独用一个channel.
func NewConsumer(queue string, handler func([]byte) error) error {
	conn, err := ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer ConnPool.ReturnConnection(conn)

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("new mq channel err: %v", err)
	}

	deliveries, err := ch.Consume(queue, "", false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("consume err: %v, queue: %s", err, queue)
	}

	for msg := range deliveries {
		err = handler(msg.Body)
		if err != nil {
			_ = msg.Reject(true)
			continue
		}
		_ = msg.Ack(false)
	}

	return nil
}

// Channel wrap amqp.Connection.Channel, get a auto reconnect channel
func (c *Connection) Channel() (*Channel, error) {
	ch, err := c.Connection.Channel()
	if err != nil {
		return nil, err
	}

	channel := &Channel{
		Channel: ch,
	}

	go func() {
		for {
			reason, ok := <-channel.Channel.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok || channel.IsClosed() {
				log.Println("channel closed")
				_ = channel.Close() // close again, ensure closed flag set when connection closed
				break
			}
			log.Printf("channel closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(delay * time.Second)

				ch, err := c.Connection.Channel()
				if err == nil {
					log.Println("channel recreate success")
					channel.Channel = ch
					break
				}

				log.Printf("channel recreate failed, err: %v", err)
			}
		}

	}()

	return channel, nil
}

// Dial wrap amqp.Dial, dial and get reconnect connection
func Dial(url string) (*Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, err
	}

	connection := &Connection{
		Connection: conn,
	}

	go func() {
		for {
			reason, ok := <-connection.Connection.NotifyClose(make(chan *amqp.Error))
			// exit this goroutine if closed by developer
			if !ok {
				log.Println("connection closed")
				break
			}
			log.Printf("connection closed, reason: %v", reason)

			// reconnect if not closed by developer
			for {
				// wait 1s for reconnect
				time.Sleep(delay * time.Second)

				conn, err := amqp.Dial(url)
				if err == nil {
					connection.Connection = conn
					log.Println("reconnect success")
					break
				}

				log.Printf("reconnect failed, err: %v", err)
			}
		}
	}()

	return connection, nil
}

// Consume wrap amqp.Channel.Consume, the returned delivery will end only when channel closed by developer
func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := ch.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				log.Printf("consume failed, err: %v", err)
				time.Sleep(delay * time.Second)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(delay * time.Second)

			if ch.IsClosed() {
				break
			}
		}
	}()

	return deliveries, nil
}
