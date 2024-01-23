package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	OrderConfirmDelayedExchange   = "vbox.order.confirm_delayed_exchange"
	OrderConfirmDelayedRoutingKey = "vbox.order.confirm_delayed_routing_key"
	OrderConfirmDelayedQueue      = "vbox.order.confirm_delayed_queue"
	OrderConfirmDeadRoutingKey    = "vbox.order.confirm_dead_routing_key"
	OrderConfirmDeadExchange      = "vbox.order.confirm_dead_exchange"
	OrderConfirmDeadQueue         = "vbox.order.confirm_dead_queue"
)

// OrderConfirmTask 处理查单（如果单子支付成功，发起回调通知的消息）
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
	if err = ch.ExchangeDeclare(OrderConfirmDeadExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclare(OrderConfirmDeadQueue); err != nil {
		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueBind(OrderConfirmDeadQueue, OrderConfirmDeadRoutingKey, OrderConfirmDeadExchange); err != nil {
		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
	}
	if err = ch.QueueDeclareWithDelay(OrderConfirmDelayedQueue, OrderConfirmDeadExchange, OrderConfirmDeadRoutingKey); err != nil {
		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
	}
	if err = ch.ExchangeDeclare(OrderConfirmDelayedExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
	}
	if err = ch.QueueBind(OrderConfirmDelayedQueue, OrderConfirmDelayedRoutingKey, OrderConfirmDelayedExchange); err != nil {
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
				global.GVA_LOG.Error("mq 消费者异常， err", zap.Error(err), zap.Any("queue", OrderConfirmDeadQueue))
			}

			for msg := range deliveries {
				//err = handler(msg.Body)
				v := request.PayOrderAndCtx{}
				err := json.Unmarshal(msg.Body, &v)
				nowTime := time.Now()
				if err != nil {
					global.GVA_LOG.Error("MqOrderConfirmTask...", zap.Error(err), zap.Any("错误消息体", msg.Body))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}
				global.GVA_LOG.Info("收到需要查询付款状态的订单消息", zap.Any("orderId", v.Obj.OrderId))

				// 直接先查一下单，如果已经补单过，则直接跳过
				var odDB vbox.PayOrder
				errQ := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id =?", v.Obj.ID).First(&odDB).Error
				if errQ != nil {
					global.GVA_LOG.Error("订单匹配消息查库失败", zap.Error(errQ))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}
				if odDB.HandStatus == 1 {
					global.GVA_LOG.Info("该订单已经补单过，跳过", zap.Any("orderId", v.Obj.OrderId))
					_ = msg.Ack(true)
					continue
				}

				//1. 筛选匹配是哪个产品
				chanID := v.Obj.ChannelCode
				accID := v.Obj.AcId
				var vca vbox.ChannelAccount
				err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("ac_id = ?", accID).First(&vca).Error
				if err != nil {
					global.GVA_LOG.Error("订单匹配消息查库失败", zap.Any("err", err.Error()))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}
				//2. 查询订单（账号）的充值情况
				expTime := v.Obj.ExpTime
				duration := expTime.Sub(nowTime)
				if duration > 0 {
					global.GVA_LOG.Info("该订单还没过期，继续查")
				} else {
					global.GVA_LOG.Info("该订单已经过期")
					//过期了， 更新成过期状态，消息丢掉
					if err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 3).Error; err != nil {
						global.GVA_LOG.Error("更新订单异常", zap.Error(err))
						// 如果解析消息失败，则直接丢弃消息
						_ = msg.Reject(false)
						continue
					}
					_ = msg.Ack(true)

					// 同时把订单 redis信息也设置一下缓存信息
					v.Obj.OrderStatus = 3
					odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
					jsonString, _ := json.Marshal(v.Obj)
					global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)

					continue
				}

				money := v.Obj.Money
				if global.TxContains(chanID) {

					global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", *odDB.CreatedAt), zap.Any("传入的过期时间", *expTime))
					records, errQ := product.QryQQRecordsBetween(vca, *odDB.CreatedAt, *expTime)
					if errQ != nil {
						// 查单有问题，直接订单要置为超时，消息置为处理完毕
						global.GVA_LOG.Error("查询充值记录异常", zap.Error(errQ))
						// 重新丢回去 下一个20s再查一次
						marshal, _ := json.Marshal(v)
						err = ch.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 20*time.Second)
						_ = msg.Ack(true)
						continue
					}
					rdMap := product.Classifier(records.WaterList)
					if vm, ok := rdMap["Q币"]; !ok {
						global.GVA_LOG.Info("还没有QB的充值记录")
					} else {
						if rd, ok2 := vm[strconv.FormatInt(int64(money*100), 10)]; !ok2 {
							global.GVA_LOG.Info("还没有QB的充值记录")

						} else { // 证明这种金额的，充上了
							if utils.Contains(rd, vca.AcAccount) {
								//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
								if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 1).Error; err != nil {
									global.GVA_LOG.Error("更新订单异常", zap.Error(err))
									_ = msg.Reject(false)
									continue
								}
								_ = msg.Ack(true)
								global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("orderId", v.Obj.OrderId))

								// 同时把订单 redis信息也设置一下缓存信息
								v.Obj.OrderStatus = 1
								odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
								jsonString, _ := json.Marshal(v.Obj)
								global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)

								// 并且发起一个回调通知的消息
								marshal, _ := json.Marshal(v)
								err = ch.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
								global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))

								continue
							}
						}
					}
				} else if global.J3Contains(chanID) {
					global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", *odDB.CreatedAt), zap.Any("传入的过期时间", *expTime))
					j3Record, errQy := product.QryJ3Record(vca)
					if errQy != nil {
						// 查单有问题，直接订单要置为超时，消息置为处理完毕
						global.GVA_LOG.Error("查询充值记录异常", zap.Error(errQy))
						// 重新丢回去 下一个20s再查一次
						marshal, _ := json.Marshal(v)
						err = ch.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 20*time.Second)
						_ = msg.Ack(true)
						continue
					}

					nowBalance := j3Record.LeftCoins
					J3AccBalanceKey := fmt.Sprintf(global.J3AccBalanceZSet, vca.AcId)
					nowTimeUnix := odDB.CreatedAt.Unix()
					//对账时间
					checkTime := time.Now().Unix()
					valMembers, err := global.GVA_REDIS.ZRangeByScore(context.Background(), J3AccBalanceKey, &redis.ZRangeBy{
						Min:    strconv.FormatInt(nowTimeUnix, 10),
						Max:    strconv.FormatInt(nowTimeUnix, 10),
						Offset: 0,
						Count:  1,
					}).Result()
					if len(valMembers) > 0 {
						mem := valMembers[0]
						split := strings.Split(mem, "_")
						hisBalance, _ := strconv.Atoi(split[4])
						if hisBalance+money*100 == nowBalance { // 充值成功的情况

							keyMem := fmt.Sprintf("%s_%s_%v_%d_%d_%d_%d", v.Obj.OrderId, vca.AcAccount, money, nowTimeUnix, hisBalance, checkTime, nowBalance)
							delMem := fmt.Sprintf("%s_%s_%v_%d_%d_%d_%d", v.Obj.OrderId, vca.AcAccount, money, nowTimeUnix, hisBalance, 0, 0)
							global.GVA_REDIS.ZAdd(context.Background(), J3AccBalanceKey, redis.Z{
								Score:  float64(nowTimeUnix),
								Member: keyMem,
							})
							global.GVA_REDIS.ZRem(context.Background(), J3AccBalanceKey, delMem)

							v.Obj.PlatId = keyMem

							//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
							if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 1).
								Update("plat_id", keyMem).Error; err != nil {
								global.GVA_LOG.Error("更新订单异常", zap.Error(err))
								_ = msg.Reject(false)
								continue
							}

							// 4.入库wallet
							var count int64
							global.GVA_DB.Model(&vbox.UserWallet{}).Where("event_id = ?", v.Obj.OrderId).Count(&count)

							if count == 0 {
								wallet := vbox.UserWallet{
									Uid:       v.Obj.CreatedBy,
									CreatedBy: v.Obj.CreatedBy,
									Type:      global.WalletOrderType,
									EventId:   v.Obj.OrderId,
									Recharge:  -money,
									Remark:    fmt.Sprintf(global.WalletEventOrderCost, money, v.Obj.ChannelCode, v.Obj.OrderId),
								}

								global.GVA_DB.Model(&vbox.UserWallet{}).Save(&wallet)
							}

							_ = msg.Ack(true)
							global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("orderId", v.Obj.OrderId))

							// 同时把订单 redis信息也设置一下缓存信息
							v.Obj.OrderStatus = 1
							odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
							jsonString, _ := json.Marshal(v.Obj)
							global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)

							// 并且发起一个回调通知的消息
							marshal, _ := json.Marshal(v)
							err = ch.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
							global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))

							_ = msg.Ack(true)
							continue
						} else {
							global.GVA_LOG.Info("未对账成功，当前余额为", zap.Any("nowBalance", nowBalance), zap.Any("hisBalance", hisBalance), zap.Any("money", money))
							// 重新丢回去 下一个20s再查一次
							marshal, _ := json.Marshal(v)
							err = ch.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 20*time.Second)
							_ = msg.Ack(true)
							continue
						}

					} else {
						// 异常情况处理
						global.GVA_LOG.Error("订单匹配消息查redis数据异常", zap.Error(err), zap.Any("valMembers", valMembers))
						// 如果解析消息失败，则直接丢弃消息
						_ = msg.Reject(false)
						continue
					}

				} else if global.PcContains(chanID) {
					global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", *odDB.CreatedAt), zap.Any("传入的过期时间", *expTime))

					records, errQy := product.QryQQRecordsByID(vca, odDB.PlatId)
					if errQy != nil {
						// 查单有问题，直接订单要置为超时，消息置为处理完毕
						global.GVA_LOG.Error("查询充值记录异常", zap.Error(errQy))
						// 重新丢回去 下一个20s再查一次
						marshal, _ := json.Marshal(v)
						err = ch.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 20*time.Second)
						_ = msg.Ack(true)
						continue
					}
					rdMap := product.Classifier(records.WaterList)
					if vm, ok := rdMap["Q币"]; !ok {
						global.GVA_LOG.Info("还没有QB的充值记录")
					} else {
						if rd, ok2 := vm[strconv.FormatInt(int64(money*100), 10)]; !ok2 {
							global.GVA_LOG.Info("还没有QB的充值记录")

						} else { // 证明这种金额的，充上了
							if utils.Contains(rd, vca.AcAccount) {
								//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
								if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 1).Error; err != nil {
									global.GVA_LOG.Error("更新订单异常", zap.Error(err))
									_ = msg.Reject(false)
									continue
								}

								// 4.入库wallet
								var count int64
								global.GVA_DB.Model(&vbox.UserWallet{}).Where("event_id = ?", v.Obj.OrderId).Count(&count)

								if count == 0 {
									wallet := vbox.UserWallet{
										Uid:       v.Obj.CreatedBy,
										CreatedBy: v.Obj.CreatedBy,
										Type:      global.WalletOrderType,
										EventId:   v.Obj.OrderId,
										Recharge:  -money,
										Remark:    fmt.Sprintf(global.WalletEventOrderCost, money, v.Obj.ChannelCode, v.Obj.OrderId),
									}

									global.GVA_DB.Model(&vbox.UserWallet{}).Save(&wallet)
								}

								_ = msg.Ack(true)
								global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("orderId", v.Obj.OrderId))

								// 同时把订单 redis信息也设置一下缓存信息
								v.Obj.OrderStatus = 1
								odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
								jsonString, _ := json.Marshal(v.Obj)
								global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)

								// 并且发起一个回调通知的消息
								marshal, _ := json.Marshal(v)
								err = ch.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
								global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))

								_ = msg.Ack(true)
								continue
							}
						}
					}
				}

				if err != nil {
					_ = msg.Reject(false)
					continue
				}

				// 重新丢回去 下一个20s再查一次
				marshal, err := json.Marshal(v)
				err = ch.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 20*time.Second)
			}
			wg.Done()
		}(i + 1)
	}

	// 等待所有消费者完成处理
	wg.Wait()
	//time.Sleep(time.Minute)
	global.GVA_LOG.Info("MqOrderConfirmTask 初始化搞定")

}
