package mq

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	MQ = new(Mq)
)

type (
	Connection struct {
		*amqp.Connection
	}

	Channel struct {
		*amqp.Channel
		closed int32
	}

	ConnPool struct {
		pool  chan *Connection
		mutex sync.Mutex
	}

	Mq struct {
		ConnPool *ConnPool
		Channel  *Channel
	}
)

func (m *Mq) Init() (err error) {
	mqCfg := global.GVA_CONFIG.RabbitMQ
	m.ConnPool = &ConnPool{
		pool: make(chan *Connection, mqCfg.PoolSize),
	}

	for i := 0; i < mqCfg.PoolSize; i++ {
		url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
			mqCfg.Username,
			mqCfg.Password,
			mqCfg.Addr,
			mqCfg.Port)
		conn, err := MQ.Dial(url)
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ: %v", err)
		}
		m.ConnPool.pool <- conn
	}
	return err
}

// Dial wrap amqp.Dial, dial and get reconnect connection
func (m *Mq) Dial(url string) (*Connection, error) {
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
				time.Sleep(global.GVA_CONFIG.RabbitMQ.RetryReConnDelay * time.Second)

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
	global.GVA_LOG.Debug(fmt.Sprintf("publish with delay + exp : %d", timer.Milliseconds()), zap.Any("exchange", exchange), zap.Any("key", key), zap.Any("timer", timer))
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

// Consume wrap amqp.Channel.Consume, the returned delivery will end only when channel closed by developer
func (ch *Channel) Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	deliveries := make(chan amqp.Delivery)

	go func() {
		for {
			d, err := ch.Channel.Consume(queue, consumer, autoAck, exclusive, noLocal, noWait, args)
			if err != nil {
				log.Printf("consume failed, err: %v", err)
				time.Sleep(global.GVA_CONFIG.RabbitMQ.RetryReConnDelay * time.Second)
				continue
			}

			for msg := range d {
				deliveries <- msg
			}

			// sleep before IsClose call. closed flag may not set before sleep.
			time.Sleep(global.GVA_CONFIG.RabbitMQ.RetryReConnDelay * time.Second)

			if ch.IsClosed() {
				break
			}
		}
	}()

	return deliveries, nil
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
			zap.L().Info(fmt.Sprintf("channel closed, reason: %v", reason))

			// reconnect if not closed by developer
			for {
				// wait 1s for connection reconnect
				time.Sleep(global.GVA_CONFIG.RabbitMQ.RetryReConnDelay * time.Second)

				ch, err := c.Connection.Channel()
				if err == nil {
					zap.L().Error(fmt.Sprintf("channel recreate success"))
					channel.Channel = ch
					break
				}

				zap.L().Info(fmt.Sprintf("channel recreate failed, err: %v", err))
			}
		}

	}()

	return channel, nil
}

func (cp *ConnPool) GetConnection() (*Connection, error) {
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

func (cp *ConnPool) ReturnConnection(conn *Connection) {
	cp.pool <- conn // 归还连接到连接池
}
