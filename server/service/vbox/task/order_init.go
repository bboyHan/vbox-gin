package task

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"go.uber.org/zap"
	"log"
	"sync"
)

var (
	orderWaitExchange = "vbox.order.waiting_exchange"
	orderWaitQueue    = "vbox.order.waiting_queue"
	orderWaitKey      = "vbox.order.waiting"
)

func OrderWaitingTask() {

	// 示例：发送消息
	conn, err := initialize.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer initialize.MQ.ConnPool.ReturnConnection(conn)

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
				v := &map[string]interface{}{}
				err := json.Unmarshal(msg.Body, v)
				global.GVA_LOG.Info(fmt.Sprintf("%v", v))

				/*v := &vbox.VboxPayOrder{}
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

				global.GVA_LOG.Info("匹配到账号了，发一个准备查单的消息 : ", zap.Any("对应单号", v.OrderId))*/

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
