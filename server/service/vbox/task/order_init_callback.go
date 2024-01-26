package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxRep "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	OrderCallbackExchange = "vbox.order.callback_exchange"
	OrderCallbackQueue    = "vbox.order.callback_queue"
	OrderCallbackKey      = "vbox.order.callback"
)

// OrderCallbackTask 订单回调
func OrderCallbackTask() {

	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 创建 订单初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err = ch.ExchangeDeclare(OrderCallbackExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 111:", zap.Any("err", err))
	}
	if err = ch.QueueDeclare(OrderCallbackQueue); err != nil {
		global.GVA_LOG.Error("create queue err 111:", zap.Any("err", err))
	}
	if err = ch.QueueBind(OrderCallbackQueue, OrderCallbackKey, OrderCallbackExchange); err != nil {
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
			deliveries, errC := ch.Consume(OrderCallbackQueue, "", false, false, false, false, nil)
			if errC != nil {
				global.GVA_LOG.Error("err", zap.Error(errC), zap.Any("queue", OrderCallbackQueue))
			}

			for msg := range deliveries {
				//v1 := &map[string]interface{}{}
				//err := json.Unmarshal(msg.Body, v1)
				//global.GVA_LOG.Info(fmt.Sprintf("%v", v1))
				now := time.Now()
				var operationRecordService system.OperationRecordService

				v := request.PayOrderAndCtx{}
				err = json.Unmarshal(msg.Body, &v)
				if err != nil {
					global.GVA_LOG.Info("MqOrder Callback Task..." + err.Error())
				}
				orderId := v.Obj.OrderId
				global.GVA_LOG.Info("收到一条需要进行发起回调的订单消息", zap.Any("order ID", orderId))

				/*msgID := fmt.Sprintf(global.MsgFilterMem, msg.MessageId, v.Obj.OrderId)
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

				//1. 先直接设置订单状态
				key := fmt.Sprintf(global.PayOrderKey, orderId)
				var order vbox.PayOrder
				err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).First(&order).Error

				global.GVA_LOG.Info("最新订单信息", zap.Any("orderId", orderId), zap.Any("订单状态", order.OrderStatus))
				//查出来了，设置一下redis
				order.OrderStatus = 1
				jsonString, _ := json.Marshal(order)
				global.GVA_REDIS.Set(context.Background(), key, jsonString, 10*time.Second)

				if err != nil {
					global.GVA_LOG.Error("订单匹配异常，消息丢弃", zap.Any("对应单号", orderId), zap.Error(err))

					_ = msg.Reject(false)
					continue
				}

				//1.0 核查商户
				var vpa vbox.PayAccount
				if strings.Contains(v.Obj.PAccount, "TEST") {
					global.GVA_LOG.Info("测试单，商户检测跳过", zap.Any("入参商户", v.Obj.PAccount))
					vpa = vbox.PayAccount{
						PAccount: v.Obj.PAccount,
						Uid:      v.Obj.CreatedBy,
					}
				} else {
					var count int64
					count, err = global.GVA_REDIS.Exists(context.Background(), global.PayAccPrefix+v.Obj.PAccount).Result()
					if count == 0 {
						if err != nil {
							global.GVA_LOG.Error("当前缓存池无此商户，redis err", zap.Error(err))
						}
						global.GVA_LOG.Info("当前缓存池无此商户，查一下库。。。", zap.Any("入参商户ID", v.Obj.PAccount))

						err = global.GVA_DB.Table("vbox_pay_account").
							Where("p_account = ?", v.Obj.PAccount).First(&vpa).Error
						jsonStr, _ := json.Marshal(vpa)
						global.GVA_REDIS.Set(context.Background(), global.PayAccPrefix+v.Obj.PAccount, jsonStr, 10*time.Minute)
					} else {
						jsonStr, _ := global.GVA_REDIS.Get(context.Background(), global.PayAccPrefix+v.Obj.PAccount).Bytes()
						err = json.Unmarshal(jsonStr, &vpa)
					}
				}

				notifyUrl := v.Obj.NotifyUrl
				client := vbHttp.NewHTTPClient()
				var headers = map[string]string{
					"Content-Type": "application/json",
				}
				var payUrl string
				payUrl, err = HandlePayUrl2PAcc(orderId)

				signBody := &vboxRep.Order2PayAccountRes{
					OrderId:   orderId,
					Money:     v.Obj.Money,
					Status:    1,
					NotifyUrl: notifyUrl,
					PayUrl:    payUrl,
					Key:       vpa.PKey,
				}
				//global.GVA_LOG.Info("初始body", zap.Any("body", signBody))
				sign := utils.CalSign(signBody)
				signBody.Sign = sign

				notifyBody := vboxRep.OrderSign2PayAccountRes{
					OrderId:   signBody.OrderId,
					Money:     signBody.Money,
					Status:    signBody.Status,
					NotifyUrl: signBody.NotifyUrl,
					PayUrl:    signBody.PayUrl,
					Sign:      signBody.Sign,
				}

				global.GVA_LOG.Info("请求地址", zap.Any("notifyUrl", notifyUrl))
				global.GVA_LOG.Info("请求body", zap.Any("notifyBody", notifyBody))

				var options = &vbHttp.RequestOptions{
					Headers:      headers,
					MaxRedirects: 3,
					PayloadType:  "json",
					Payload:      notifyBody,
				}

				response, errH := client.Post(notifyUrl, options)
				if errH != nil {
					global.GVA_LOG.Error("回调异常", zap.Error(errH))
					_ = msg.Reject(false)

					//入库操作记录
					record := sysModel.SysOperationRecord{
						Ip:      v.Ctx.ClientIP,
						Method:  v.Ctx.Method,
						Path:    v.Ctx.UrlPath,
						Agent:   v.Ctx.UserAgent,
						Status:  500,
						Latency: time.Since(now),
						Resp:    fmt.Sprintf(global.NotifyEx, errH.Error(), response),
						UserID:  v.Ctx.UserID,
					}
					err = operationRecordService.CreateSysOperationRecord(record)
					if err != nil {
						global.GVA_LOG.Error("record 入库失败..." + err.Error())
					}

					continue
				}
				global.GVA_LOG.Info("回调响应消息", zap.Any("状态码", response.StatusCode), zap.Any("响应内容", string(response.Body)))

				nowTime := time.Now()
				if v.Obj.HandStatus == 3 {
					//3. 更新回调成功的状态
					if errD := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).
						Update("order_status", 1).Update("cb_status", 1).Update("hand_status", 1).Update("cb_time", nowTime).Error; errD != nil {
						global.GVA_LOG.Error("更新订单异常", zap.Error(errD))
						_ = msg.Reject(false)
						continue
					}

					//	补单单独记录一下日志
					record := sysModel.SysOperationRecord{
						Ip:      v.Ctx.ClientIP,
						Method:  v.Ctx.Method,
						Path:    v.Ctx.UrlPath,
						Agent:   v.Ctx.UserAgent,
						Status:  200,
						Latency: time.Since(now),
						Resp:    fmt.Sprintf(global.NotifyHandSuccess, response.StatusCode, string(response.Body)),
						UserID:  v.Ctx.UserID,
					}
					err = operationRecordService.CreateSysOperationRecord(record)
					if err != nil {
						global.GVA_LOG.Error("record 入库失败..." + err.Error())
					}
				} else {
					//3. 更新回调成功的状态
					if errD := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).
						Update("cb_status", 1).Update("cb_time", nowTime).Error; errD != nil {
						global.GVA_LOG.Error("更新订单异常", zap.Error(errD))
						_ = msg.Reject(false)
						continue
					}
				}

				// 4.入库wallet
				var c int64
				global.GVA_DB.Model(&vbox.UserWallet{}).Where("event_id = ?", orderId).Count(&c)

				if c == 0 {
					wallet := vbox.UserWallet{
						Uid:       v.Obj.CreatedBy,
						CreatedBy: v.Obj.CreatedBy,
						Type:      global.WalletOrderType,
						EventId:   orderId,
						Recharge:  -v.Obj.Money,
						Remark:    fmt.Sprintf(global.WalletEventOrderCost, v.Obj.Money, v.Obj.ChannelCode, orderId),
					}

					global.GVA_DB.Model(&vbox.UserWallet{}).Save(&wallet)
				}

				_ = msg.Ack(true)
				global.GVA_LOG.Info("订单完成，回调完成", zap.Any("对应单号", orderId))

			}
			wg.Done()
		}(i + 1)
	}
	global.GVA_LOG.Info("Vbox OrderCallbackTask 初始化搞定")
	// 等待所有消费者完成处理
	wg.Wait()
}
