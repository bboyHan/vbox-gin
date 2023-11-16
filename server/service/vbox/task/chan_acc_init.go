package task

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
	"log"
	"strconv"
	"sync"
)

// 账号开启查询

const (
	ChanAccEnableCheckExchange = "vbox.channel.acc_enable_check_exchange"
	ChanAccEnableCheckQueue    = "vbox.order.acc_enable_check_queue"
	ChanAccEnableCheckKey      = "vbox.order.acc_enable_check"
)

// ChanAccEnableCheckTask 订单入库匹配
func ChanAccEnableCheckTask() {
	var operationRecordService system.OperationRecordService

	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 创建 账号开启检查 初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(ChanAccEnableCheckExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(ChanAccEnableCheckQueue); err != nil {
		global.GVA_LOG.Error("create queue err:", zap.Any("err", err))
	}
	if err := ch.QueueBind(ChanAccEnableCheckQueue, ChanAccEnableCheckKey, ChanAccEnableCheckExchange); err != nil {
		global.GVA_LOG.Error("bind queue err:", zap.Any("err", err))
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
			deliveries, err := ch.Consume(ChanAccEnableCheckQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", ChanAccEnableCheckQueue))
			}

			for msg := range deliveries {
				//v := &map[string]interface{}{}
				//err := json.Unmarshal(msg.Body, v)
				global.GVA_LOG.Info(fmt.Sprintf("%v", msg.Body))

				v := vboxReq.ChanAccAndCtx{}
				err := json.Unmarshal(msg.Body, &v)
				if err != nil {
					global.GVA_LOG.Info("错了，直接丢了..." + err.Error())
					_ = msg.Reject(false)
					continue
				}
				_ = msg.Ack(true)

				// 1. 查询该用户的余额是否充足
				var balance int
				err = global.GVA_DB.Model(&vbox.UserWallet{}).Select("sum(recharge) as balance").Where("uid = ?", v.Obj.CreatedBy).First(&balance).Error
				if balance < 0 { //余额不足，则 log 一条
					//
					record := sysModel.SysOperationRecord{
						Ip:     v.Ctx.ClientIP,
						Method: v.Ctx.Method,
						Path:   v.Ctx.UrlPath,
						Agent:  v.Ctx.UserAgent,
						Body:   v.Ctx.Body,
						UserID: v.Ctx.UserID,
					}

					err = operationRecordService.CreateSysOperationRecord(record)
					if err != nil {
						return
					}
				}

				// 2. 查询账号是否有超 金额限制，或者笔数限制
				// 2.1 日限制
				if v.Obj.DailyLimit > 0 {
				}
				// 2.2 总限制
				if v.Obj.TotalLimit > 0 {

				}
				// 2.3 笔数限制
				if v.Obj.CountLimit > 0 {

				}

				// 3. 筛选匹配是哪个产品，查一下对应产品的账户是否能够正常官方使用
				chanID, err := strconv.Atoi(v.Obj.Cid)
				if chanID >= 1000 && chanID <= 1099 { //腾讯

				} else if chanID >= 2000 && chanID <= 2099 { //剑三

				}

				//2. 查询产品对应的账号池是否有可用账号

				global.GVA_LOG.Info("消息 : ", zap.Any("msg", msg.Body))
				err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).Update("sys_status", 1).Error

				if err != nil {
					_ = msg.Reject(false)
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
