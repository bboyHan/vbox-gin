package task

//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/flipped-aurora/gin-vue-admin/server/global"
//	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
//	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
//	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
//	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
//	"github.com/flipped-aurora/gin-vue-admin/server/mq"
//	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
//	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
//	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
//	"github.com/flipped-aurora/gin-vue-admin/server/utils"
//	"github.com/redis/go-redis/v9"
//	"go.uber.org/zap"
//	"strconv"
//	"strings"
//	"sync"
//	"time"
//)
//
//const (
//	OrderConfirmDelayedExchange   = "vbox.order.confirm_delayed_exchange"
//	OrderConfirmDelayedRoutingKey = "vbox.order.confirm_delayed_routing_key"
//	OrderConfirmDelayedQueue      = "vbox.order.confirm_delayed_queue"
//	OrderConfirmDeadRoutingKey    = "vbox.order.confirm_dead_routing_key"
//	OrderConfirmDeadExchange      = "vbox.order.confirm_dead_exchange"
//	OrderConfirmDeadQueue         = "vbox.order.confirm_dead_queue"
//)
//
//// OrderConfirmTask 处理查单（如果单子支付成功，发起回调通知的消息）
//func OrderConfirmTask() {
//	var operationRecordService system.OperationRecordService
//
//	// 示例：发送消息
//	conn, err := mq.MQ.ConnPool.GetConnection()
//	if err != nil {
//		//log.Fatalf("Failed to get connection from pool: %v", err)
//		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
//	}
//	defer mq.MQ.ConnPool.ReturnConnection(conn)
//
//	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
//	ch, _ := conn.Channel()
//	if err = ch.ExchangeDeclare(OrderConfirmDeadExchange, "direct"); err != nil {
//		global.GVA_LOG.Error("create exchange err 222:", zap.Any("err", err))
//	}
//	if err = ch.QueueDeclare(OrderConfirmDeadQueue); err != nil {
//		global.GVA_LOG.Error("create queue err 222:", zap.Any("err", err))
//	}
//	if err = ch.QueueBind(OrderConfirmDeadQueue, OrderConfirmDeadRoutingKey, OrderConfirmDeadExchange); err != nil {
//		global.GVA_LOG.Error("bind queue err 222:", zap.Any("err", err))
//	}
//	if err = ch.QueueDeclareWithDelay(OrderConfirmDelayedQueue, OrderConfirmDeadExchange, OrderConfirmDeadRoutingKey); err != nil {
//		global.GVA_LOG.Error("create queue err 333:", zap.Any("err", err))
//	}
//	if err = ch.ExchangeDeclare(OrderConfirmDelayedExchange, "direct"); err != nil {
//		global.GVA_LOG.Error("create exchange err 333:", zap.Any("err", err))
//	}
//	if err = ch.QueueBind(OrderConfirmDelayedQueue, OrderConfirmDelayedRoutingKey, OrderConfirmDelayedExchange); err != nil {
//		global.GVA_LOG.Error("bind queue err 333:", zap.Any("err", err))
//	}
//
//	// 设置初始消费者数量
//	consumerCount := 30
//	// 使用 WaitGroup 来等待所有消费者完成处理
//	var wg sync.WaitGroup
//	wg.Add(consumerCount)
//	// 启动多个消费者
//	for i := 0; i < consumerCount; i++ {
//		go func(consumerID int) {
//			connX, errX := mq.MQ.ConnPool.GetConnection()
//			if errX != nil {
//				//log.Fatalf("Failed to get connection from pool: %v", err)
//				global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(errX))
//			}
//			defer mq.MQ.ConnPool.ReturnConnection(connX)
//			if connX == nil {
//				global.GVA_LOG.Error("connX == nil")
//				return
//			}
//			chX, _ := connX.Channel()
//
//			// 说明：执行查单回调处理
//			deliveries, errC := chX.Consume(OrderConfirmDeadQueue, "", false, false, false, false, nil)
//			if errC != nil {
//				global.GVA_LOG.Error("mq 消费者异常， err", zap.Error(errC), zap.Any("queue", OrderConfirmDeadQueue))
//			}
//
//			for msg := range deliveries {
//
//				v := request.PayOrderAndCtx{}
//				err := json.Unmarshal(msg.Body, &v)
//				nowTime := time.Now()
//				if err != nil {
//					global.GVA_LOG.Error("MqOrderConfirmTask...", zap.Error(err), zap.Any("错误消息体", msg.Body))
//					// 如果解析消息失败，则直接丢弃消息
//					_ = msg.Reject(false)
//					continue
//				}
//				global.GVA_LOG.Info("收到需要查询付款状态的订单消息", zap.Any("orderId", v.Obj.OrderId))
//
//				// 直接先查一下单，如果已经补单过，则直接跳过
//				var odDB vbox.PayOrder
//				errQ := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id =?", v.Obj.ID).First(&odDB).Error
//				if errQ != nil {
//					global.GVA_LOG.Error("订单匹配消息查库失败", zap.Error(errQ))
//					// 如果解析消息失败，则直接丢弃消息
//					_ = msg.Reject(false)
//					continue
//				}
//				if odDB.HandStatus == 1 {
//					global.GVA_LOG.Info("该订单已经补单过，跳过", zap.Any("orderId", v.Obj.OrderId))
//					_ = msg.Ack(true)
//					continue
//				}
//
//				//1. 筛选匹配是哪个产品
//				chanID := v.Obj.ChannelCode
//				accID := v.Obj.AcId
//				var vca vbox.ChannelAccount
//				err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("ac_id = ?", accID).First(&vca).Error
//				if err != nil {
//					global.GVA_LOG.Error("订单匹配消息查库失败", zap.Any("err", err.Error()))
//					// 如果解析消息失败，则直接丢弃消息
//					_ = msg.Reject(false)
//					continue
//				}
//
//				orgTmp := utils2.GetSelfOrg(vca.CreatedBy)
//
//				//2. 查询订单（账号）的充值情况
//				expTime := v.Obj.ExpTime
//				duration := expTime.Sub(nowTime)
//				if duration > 0 {
//					global.GVA_LOG.Info("该订单还没过期，继续查", zap.Any("orderID", v.Obj.OrderId))
//				} else {
//					global.GVA_LOG.Info("该订单已经过期", zap.Any("orderID", v.Obj.OrderId))
//					//过期了， 更新成过期状态，消息丢掉
//					if err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 3).Error; err != nil {
//						global.GVA_LOG.Error("更新订单异常", zap.Error(err))
//						// 如果解析消息失败，则直接丢弃消息
//						_ = msg.Reject(false)
//						continue
//					}
//					_ = msg.Ack(true)
//
//					// 同时把订单 redis信息也设置一下缓存信息
//					v.Obj.OrderStatus = 3
//					odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
//					jsonString, _ := json.Marshal(v.Obj)
//					global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)
//
//					continue
//				}
//
//				money := v.Obj.Money
//				odCreatedTime := *odDB.CreatedAt
//
//				var startTime, endTime time.Time
//				var errQryRecord error
//				// 处理核对订单的查询时间范围
//				switch chanID {
//				case "1001":
//					//处理查询的起始时间（提前3分钟）
//					startTime = odCreatedTime.Add(-180 * time.Second)
//					endTime = *expTime
//					records, errX, qryURL := product.QryQQRecordsBetween(vca, startTime, *expTime)
//
//				}
//
//				if global.TxContains(chanID) {
//
//					//global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", odCreatedTime), zap.Any("传入的过期时间", *expTime))
//					//处理查询的起始时间（提前3分钟）
//					startTime := odCreatedTime.Add(-180 * time.Second)
//					//global.GVA_LOG.Info("查询的起始时间", zap.Any("查询的起始时间", startTime))
//					records, errX, qryURL := product.QryQQRecordsBetween(vca, startTime, *expTime)
//
//					if errX != nil {
//						// 查单有问题，直接订单要置为超时，消息置为处理完毕
//						global.GVA_LOG.Error("查询充值记录异常", zap.Error(errX))
//
//						if errX = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
//							global.GVA_LOG.Error("更新订单异常", zap.Error(errX))
//							_ = msg.Reject(false)
//							continue
//						}
//
//						//入库操作记录
//						record := sysModel.SysOperationRecord{
//							Ip:      v.Ctx.ClientIP,
//							Method:  v.Ctx.Method,
//							Path:    v.Ctx.UrlPath,
//							Agent:   v.Ctx.UserAgent,
//							Status:  500,
//							Latency: time.Since(time.Now()),
//							Resp:    fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//							UserID:  v.Ctx.UserID,
//						}
//
//						err = operationRecordService.CreateSysOperationRecord(record)
//						if err != nil {
//							global.GVA_LOG.Error("当前账号回调查单，官方查单异常，record 入库失败..." + err.Error())
//						}
//
//						err = global.GVA_DB.Debug().Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).
//							Update("sys_status", 0).Error
//
//						vca.Status = 0
//						oc := vboxReq.ChanAccAndCtx{
//							Obj: vca,
//							Ctx: vboxReq.Context{
//								Body:      fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//								ClientIP:  "127.0.0.1",
//								Method:    "POST",
//								UrlPath:   "/vca/switchEnable",
//								UserAgent: "",
//								UserID:    int(vca.CreatedBy),
//							},
//						}
//						marshal, _ := json.Marshal(oc)
//
//						_ = chX.Publish(ChanAccEnableCheckExchange, ChanAccEnableCheckKey, marshal)
//
//						global.GVA_LOG.Error("处理订单为失败，发起一条关号清理资源", zap.Any("orderID", v.Obj.OrderId), zap.Any("ac_id", vca.AcId), zap.Any("ac_account", vca.AcAccount))
//
//						//// 重新丢回去 下一个20s再查一次
//						//marshal, _ := json.Marshal(v)
//						//err = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 29*time.Second)
//						_ = msg.Ack(true)
//						continue
//					}
//					rdMap := product.Classifier(records.WaterList)
//					if vm, ok := rdMap["Q币"]; !ok {
//						global.GVA_LOG.Info("还没有QB的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", v.Obj.OrderId))
//					} else {
//						if rd, ok2 := vm[strconv.FormatInt(int64(money*100), 10)]; !ok2 {
//							global.GVA_LOG.Info("还没有QB的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", v.Obj.OrderId))
//
//						} else { // 证明这种金额的，充上了
//
//							var flag bool
//							var platID string
//							for _, tg := range rd {
//								if strings.Contains(tg, vca.AcAccount) {
//									flag = true
//									platID = strings.Split(tg, ",")[1]
//									break
//								}
//							}
//							if flag {
//								global.GVA_LOG.Info("查单链接", zap.Any("订单ID", v.Obj.OrderId), zap.Any("url", qryURL))
//								global.GVA_LOG.Info("核对官方订单", zap.Any("platID", platID), zap.Any("订单ID", v.Obj.OrderId), zap.Any("查单起始时间", startTime), zap.Any("结束时间", *expTime))
//
//								//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
//								if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 1).Error; err != nil {
//									global.GVA_LOG.Error("更新订单异常", zap.Error(err))
//									_ = msg.Reject(false)
//									continue
//								}
//								_ = msg.Ack(true)
//								global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("orderId", v.Obj.OrderId))
//
//								// 同时把订单 redis信息也设置一下缓存信息
//								v.Obj.OrderStatus = 1
//								odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
//								jsonString, _ := json.Marshal(v.Obj)
//								global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)
//
//								// 并且发起一个回调通知的消息
//								marshal, _ := json.Marshal(v)
//								err = chX.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
//								global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))
//
//								//支付成功后，释放本号
//								accWaitYdKey := fmt.Sprintf(global.YdQBAccWaiting, accID, money)
//								accInfoVal := fmt.Sprintf("%d_%s_%s_%v", vca.ID, accID, vca.AcAccount, money)
//
//								cdTime := 210 * time.Second
//								// 设置一个冷却时间
//								global.GVA_LOG.Info("当前添加的账号新一轮冷却", zap.Any("accWaitYdKey", accWaitYdKey), zap.Any("ttl", cdTime), zap.Any("order ID", v.Obj.OrderId))
//								global.GVA_REDIS.Set(context.Background(), accWaitYdKey, accInfoVal, cdTime)
//
//								waitMsg := strings.Join([]string{accWaitYdKey, accInfoVal}, "-")
//								err = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), 0)
//								global.GVA_LOG.Info("3分半后，释放本账号key", zap.Any("acAccount", vca.AcAccount), zap.Any("acId", accID), zap.Any("orderId", odDB.OrderId))
//
//								continue
//							}
//
//						}
//					}
//				} else if global.DnfContains(chanID) {
//					//global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", odCreatedTime), zap.Any("传入的过期时间", *expTime))
//					//处理查询的起始时间（提前2分半钟）
//					startTime := odCreatedTime.Add(-150 * time.Second)
//					//global.GVA_LOG.Info("查询的起始时间", zap.Any("查询的起始时间", startTime))
//					records, errX, qryURL := product.QryQQRecordsBetween(vca, startTime, *expTime)
//					if errX != nil {
//						// 查单有问题，直接订单要置为超时，消息置为处理完毕
//						global.GVA_LOG.Error("查询充值记录异常", zap.Error(errX))
//
//						if errX = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
//							global.GVA_LOG.Error("更新订单异常", zap.Error(errX))
//							_ = msg.Reject(false)
//							continue
//						}
//
//						//入库操作记录
//						record := sysModel.SysOperationRecord{
//							Ip:      v.Ctx.ClientIP,
//							Method:  v.Ctx.Method,
//							Path:    v.Ctx.UrlPath,
//							Agent:   v.Ctx.UserAgent,
//							Status:  500,
//							Latency: time.Since(time.Now()),
//							Resp:    fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//							UserID:  v.Ctx.UserID,
//						}
//
//						err = operationRecordService.CreateSysOperationRecord(record)
//						if err != nil {
//							global.GVA_LOG.Error("当前账号回调查单，官方查单异常，record 入库失败..." + err.Error())
//						}
//
//						err = global.GVA_DB.Debug().Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).
//							Update("sys_status", 0).Error
//
//						vca.Status = 0
//						oc := vboxReq.ChanAccAndCtx{
//							Obj: vca,
//							Ctx: vboxReq.Context{
//								Body:      fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//								ClientIP:  "127.0.0.1",
//								Method:    "POST",
//								UrlPath:   "/vca/switchEnable",
//								UserAgent: "",
//								UserID:    int(vca.CreatedBy),
//							},
//						}
//						marshal, _ := json.Marshal(oc)
//
//						_ = chX.Publish(ChanAccEnableCheckExchange, ChanAccEnableCheckKey, marshal)
//
//						global.GVA_LOG.Error("处理订单为失败，发起一条关号清理资源", zap.Any("orderID", v.Obj.OrderId), zap.Any("ac_id", vca.AcId), zap.Any("ac_account", vca.AcAccount))
//
//						_ = msg.Ack(true)
//						continue
//					}
//					rdMap := product.Classifier(records.WaterList)
//					if vm, ok := rdMap["DNF"]; !ok {
//						global.GVA_LOG.Info("还没有Dnf的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", v.Obj.OrderId))
//					} else {
//						if rd, ok2 := vm[strconv.FormatInt(int64(money*100), 10)]; !ok2 {
//							global.GVA_LOG.Info("还没有Dnf的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", v.Obj.OrderId))
//						} else { // 证明这种金额的，充上了
//							var flag bool
//							var platID string
//							for _, tg := range rd {
//								if strings.Contains(tg, vca.AcAccount) {
//									flag = true
//									platID = strings.Split(tg, ",")[1]
//									break
//								}
//							}
//							if flag {
//								global.GVA_LOG.Info("查单链接", zap.Any("订单ID", v.Obj.OrderId), zap.Any("url", qryURL))
//								global.GVA_LOG.Info("核对官方订单", zap.Any("platID", platID), zap.Any("订单ID", v.Obj.OrderId), zap.Any("查单起始时间", startTime), zap.Any("结束时间", *expTime))
//
//								//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
//								if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 1).Error; err != nil {
//									global.GVA_LOG.Error("更新订单异常", zap.Error(err))
//									_ = msg.Reject(false)
//									continue
//								}
//								_ = msg.Ack(true)
//								global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("orderId", v.Obj.OrderId))
//
//								// 同时把订单 redis信息也设置一下缓存信息
//								v.Obj.OrderStatus = 1
//								odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
//								jsonString, _ := json.Marshal(v.Obj)
//								global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)
//
//								// 并且发起一个回调通知的消息
//								marshal, _ := json.Marshal(v)
//								err = chX.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
//								global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))
//
//								continue
//							}
//						}
//					}
//				} else if global.SdoContains(chanID) {
//
//					global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", odCreatedTime), zap.Any("传入的过期时间", *expTime))
//					records, errX := product.QrySdoDaoYuRecordBetween(vca, odCreatedTime, *expTime)
//					if errX != nil {
//						// 查单有问题，直接订单要置为超时，消息置为处理完毕
//						global.GVA_LOG.Error("查询充值记录异常", zap.Error(errX))
//
//						if errX = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
//							global.GVA_LOG.Error("更新订单异常", zap.Error(errX))
//							_ = msg.Reject(false)
//							continue
//						}
//
//						//入库操作记录
//						record := sysModel.SysOperationRecord{
//							Ip:      v.Ctx.ClientIP,
//							Method:  v.Ctx.Method,
//							Path:    v.Ctx.UrlPath,
//							Agent:   v.Ctx.UserAgent,
//							Status:  500,
//							Latency: time.Since(time.Now()),
//							Resp:    fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//							UserID:  v.Ctx.UserID,
//						}
//
//						err = operationRecordService.CreateSysOperationRecord(record)
//						if err != nil {
//							global.GVA_LOG.Error("当前账号回调查单，官方查单异常，record 入库失败..." + err.Error())
//						}
//
//						err = global.GVA_DB.Debug().Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).
//							Update("sys_status", 0).Error
//
//						vca.Status = 0
//						oc := vboxReq.ChanAccAndCtx{
//							Obj: vca,
//							Ctx: vboxReq.Context{
//								Body:      fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//								ClientIP:  "127.0.0.1",
//								Method:    "POST",
//								UrlPath:   "/vca/switchEnable",
//								UserAgent: "",
//								UserID:    int(vca.CreatedBy),
//							},
//						}
//						marshal, _ := json.Marshal(oc)
//
//						_ = chX.Publish(ChanAccEnableCheckExchange, ChanAccEnableCheckKey, marshal)
//
//						global.GVA_LOG.Error("处理订单为失败，发起一条关号清理资源", zap.Any("orderID", v.Obj.OrderId), zap.Any("ac_id", vca.AcId), zap.Any("ac_account", vca.AcAccount))
//
//						_ = msg.Ack(true)
//						continue
//					}
//					rdMap := product.Classifier(records)
//					global.GVA_LOG.Info("筛后的map", zap.Any("map", rdMap), zap.Any("records", records))
//					if vm, ok := rdMap["彩虹岛"]; !ok {
//						global.GVA_LOG.Info("还没有Sdo的充值记录，空空如也")
//					} else {
//						runStr := string(rune(money))
//						formatMoney := strconv.FormatInt(int64(money), 10)
//						//打印两个 money转化
//						global.GVA_LOG.Info("money", zap.Any("runStr", runStr), zap.Any("formatMoney", formatMoney))
//						if rd, ok2 := vm[formatMoney]; !ok2 {
//							global.GVA_LOG.Info("还没有Sdo的充值记录")
//						} else { // 证明这种金额的，充上了
//							var flag bool
//							var platID string
//							for _, tg := range rd {
//								if strings.Contains(tg, vca.AcAccount) {
//									flag = true
//									platID = strings.Split(tg, "_")[1]
//									break
//								}
//							}
//							if flag {
//								//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
//								if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).
//									Update("order_status", 1).Update("plat_id", platID).Error; err != nil {
//									global.GVA_LOG.Error("更新订单异常", zap.Error(err))
//									_ = msg.Reject(false)
//									continue
//								}
//								_ = msg.Ack(true)
//								global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("orderId", v.Obj.OrderId))
//
//								// 同时把订单 redis信息也设置一下缓存信息
//								v.Obj.OrderStatus = 1
//								odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
//								jsonString, _ := json.Marshal(v.Obj)
//								global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)
//
//								// 并且发起一个回调通知的消息
//								marshal, _ := json.Marshal(v)
//								err = chX.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
//								global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))
//
//								accJucKey := fmt.Sprintf(global.PayAccKey, accID)
//								//支付成功后，释放本号
//								waitAccYdKey := fmt.Sprintf(global.YdSdoAccWaiting, accID, money)
//
//								waitAccMem := fmt.Sprintf("%v,%s,%s", vca.ID, vca.AcId, vca.AcAccount)
//								waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
//								global.GVA_REDIS.Del(context.Background(), accJucKey)
//								global.GVA_REDIS.Del(context.Background(), waitAccYdKey)
//
//								_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), 0)
//								global.GVA_LOG.Info("释放本账号", zap.Any("acAccount", vca.AcAccount), zap.Any("acId", accID), zap.Any("orderId", odDB.OrderId))
//
//								continue
//							}
//						}
//					}
//				} else if global.ECContains(chanID) {
//					global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", odCreatedTime), zap.Any("传入的过期时间", *expTime))
//					//库记录的platID是否有值
//					if odDB.PlatId == "" {
//						global.GVA_LOG.Info("卡密尚未提交，继续查询", zap.Any("money", money), zap.Any("orderId", v.Obj.OrderId))
//						// 重新丢回去 下一个20s再查一次
//						marshal, _ := json.Marshal(v)
//						err = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 15*time.Second)
//						_ = msg.Ack(true)
//						continue
//					} else {
//						v.Obj.PlatId = odDB.PlatId
//
//						accPoolKey := fmt.Sprintf(global.ChanOrgECPoolAccZSet, orgTmp[0], chanID)
//
//						var resPoolList []string
//						resPoolList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), accPoolKey, &redis.ZRangeBy{
//							Min:    "0",
//							Max:    "0",
//							Offset: 0,
//							Count:  -1,
//						}).Result()
//
//						if err != nil {
//							global.GVA_LOG.Error("引导类匹配异常, redis err", zap.Error(err))
//							_ = msg.Ack(true)
//							if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; errDB != nil {
//								global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(errDB))
//							}
//							continue
//						}
//						var cardID, cardAccID, cardAcAccount string
//
//						var cardAcc vbox.ChannelCardAcc
//
//						if resPoolList != nil && len(resPoolList) > 0 {
//							//accTmp := resList[0]
//							accPoolTmp := utils.RandomElement(resPoolList)
//							// 2.1 把账号设置为已用
//							global.GVA_REDIS.ZAdd(context.Background(), accPoolKey, redis.Z{
//								Score:  1,
//								Member: accPoolTmp,
//							})
//							split := strings.Split(accPoolTmp, ",")
//							cardID = split[0]
//							cardAccID = split[1]
//							cardAcAccount = split[2]
//							err = global.GVA_DB.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", cardID).First(&cardAcc).Error
//							if err != nil {
//								//return nil, fmt.Errorf("匹配通道账号不存在！ 请核查：%v", err.Error())
//
//								global.GVA_LOG.Error("匹配查单池账号异常", zap.Error(err))
//								_ = msg.Ack(true)
//								continue
//							}
//							global.GVA_LOG.Info("匹配查单池账号ck", zap.Any("cardID", cardID), zap.Any("cardAccID", cardAccID), zap.Any("cardAcAccount", cardAcAccount))
//
//							//	2.3 取用的账号进入CD 冷却
//							cdTime := 10 * time.Second
//							accWaitYdKey := fmt.Sprintf(global.YdECPoolAccWaiting, cardAccID)
//							accInfoVal := fmt.Sprintf("%s,%s,%s", cardID, cardAccID, cardAcAccount)
//
//							// 设置一个冷却时间
//							ttl := global.GVA_REDIS.TTL(context.Background(), accWaitYdKey).Val()
//							if ttl > 0 {
//								global.GVA_LOG.Info("当前添加的查单池账号正在冷却中（有预产正在处理中）", zap.Any("ttl", ttl))
//								cdTime = ttl
//							}
//							global.GVA_REDIS.Set(context.Background(), accWaitYdKey, accInfoVal, cdTime)
//
//							waitMsg := strings.Join([]string{accWaitYdKey, accInfoVal}, "-")
//							err = chX.PublishWithDelay(CardAccCDCheckDelayedExchange, CardAccCDCheckDelayedRoutingKey, []byte(waitMsg), 0)
//						} else {
//							global.GVA_LOG.Error("查单池匹配账号CK不足, list size zero", zap.Error(err), zap.Any("channelCode", v.Obj.ChannelCode), zap.Any("money", money))
//							_ = msg.Ack(true)
//							if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; errDB != nil {
//								global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(errDB))
//							}
//							record := sysModel.SysOperationRecord{
//								Ip:      v.Ctx.ClientIP,
//								Method:  v.Ctx.Method,
//								Path:    v.Ctx.UrlPath,
//								Agent:   v.Ctx.UserAgent,
//								Status:  500,
//								Latency: time.Since(time.Now()),
//								Resp:    fmt.Sprintf(global.ResourceAccNotEnough, v.Obj.ChannelCode, money),
//								UserID:  v.Ctx.UserID,
//							}
//							err = operationRecordService.CreateSysOperationRecord(record)
//							continue
//						}
//
//						//a.拿账户的ck查一遍卡密，如果卡密ok则执行绑卡
//						errE := product.ECardQuery(odDB.PlatId, cardAcc.Token)
//						if errE != nil {
//
//							global.GVA_LOG.Error("匹配查单池账号异常,查询卡密合法性错误", zap.Error(errE))
//							if errX = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
//								global.GVA_LOG.Error("更新订单异常", zap.Error(errX))
//								_ = msg.Reject(false)
//								continue
//							}
//
//							_ = msg.Ack(true)
//							continue
//						}
//
//						//b.卡密不ok，直接订单要置为超时，消息置为处理完毕
//
//						ecRecord, errW := product.ECardBind(odDB.PlatId, vca.Token)
//						if errW != nil {
//							// 查单有问题，直接订单要置为超时，消息置为处理完毕
//							global.GVA_LOG.Error("查询充值记录异常", zap.Error(errW))
//
//							if errX = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
//								global.GVA_LOG.Error("更新订单异常", zap.Error(errX))
//								_ = msg.Reject(false)
//								continue
//							}
//
//							//入库操作记录
//							record := sysModel.SysOperationRecord{
//								Ip:      v.Ctx.ClientIP,
//								Method:  v.Ctx.Method,
//								Path:    v.Ctx.UrlPath,
//								Agent:   v.Ctx.UserAgent,
//								Status:  500,
//								Latency: time.Since(time.Now()),
//								Resp:    fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//								UserID:  v.Ctx.UserID,
//							}
//
//							err = operationRecordService.CreateSysOperationRecord(record)
//							if err != nil {
//								global.GVA_LOG.Error("当前账号回调查单，官方查单异常，record 入库失败..." + err.Error())
//							}
//
//							err = global.GVA_DB.Debug().Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).
//								Update("sys_status", 0).Error
//
//							vca.Status = 0
//							oc := vboxReq.ChanAccAndCtx{
//								Obj: vca,
//								Ctx: vboxReq.Context{
//									Body:      fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//									ClientIP:  "127.0.0.1",
//									Method:    "POST",
//									UrlPath:   "/vca/switchEnable",
//									UserAgent: "",
//									UserID:    int(vca.CreatedBy),
//								},
//							}
//							marshal, _ := json.Marshal(oc)
//
//							_ = chX.Publish(ChanAccEnableCheckExchange, ChanAccEnableCheckKey, marshal)
//
//							global.GVA_LOG.Error("处理订单为失败，发起一条关号清理资源", zap.Any("orderID", v.Obj.OrderId), zap.Any("ac_id", vca.AcId), zap.Any("ac_account", vca.AcAccount))
//
//							_ = msg.Ack(true)
//							continue
//						}
//						if ecRecord != nil { //绑卡成功
//							global.GVA_LOG.Info("绑卡成功，核对卡密金额信息", zap.Any("orderId", odDB.OrderId), zap.Any("ecRecord", ecRecord))
//							amount := ecRecord["OrderMoney"].(float64)
//							if amount == float64(money) {
//								//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
//								if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 1).Error; err != nil {
//									global.GVA_LOG.Error("更新订单异常", zap.Error(err))
//									_ = msg.Reject(false)
//									continue
//								}
//
//								// 4.入库wallet
//								var count int64
//								global.GVA_DB.Model(&vbox.UserWallet{}).Where("event_id = ?", v.Obj.OrderId).Count(&count)
//
//								if count == 0 {
//									wallet := vbox.UserWallet{
//										Uid:       v.Obj.CreatedBy,
//										CreatedBy: v.Obj.CreatedBy,
//										Type:      global.WalletOrderType,
//										EventId:   v.Obj.OrderId,
//										Recharge:  -money,
//										Remark:    fmt.Sprintf(global.WalletEventOrderCost, money, v.Obj.ChannelCode, v.Obj.OrderId),
//									}
//
//									global.GVA_DB.Model(&vbox.UserWallet{}).Save(&wallet)
//								}
//
//								_ = msg.Ack(true)
//								global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("pa", v.Obj.PAccount), zap.Any("orderId", v.Obj.OrderId))
//
//								// 同时把订单 redis信息也设置一下缓存信息
//								v.Obj.OrderStatus = 1
//								odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
//								jsonString, _ := json.Marshal(v.Obj)
//								global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)
//
//								// 并且发起一个回调通知的消息
//								marshal, _ := json.Marshal(v)
//								err = chX.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
//								global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))
//
//								accJucKey := fmt.Sprintf(global.PayAccKey, accID)
//								//支付成功后，释放本号
//								waitAccYdKey := fmt.Sprintf(global.YdECAccWaiting, accID)
//
//								waitAccMem := fmt.Sprintf("%v,%s,%s", vca.ID, vca.AcId, vca.AcAccount)
//								waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
//								global.GVA_REDIS.Del(context.Background(), accJucKey)
//								global.GVA_REDIS.Del(context.Background(), waitAccYdKey)
//
//								_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), 0)
//								global.GVA_LOG.Info("释放本账号", zap.Any("acAccount", vca.AcAccount), zap.Any("acId", accID), zap.Any("orderId", odDB.OrderId))
//							}
//						}
//					}
//
//				} else if global.J3Contains(chanID) {
//					global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", odCreatedTime), zap.Any("传入的过期时间", *expTime))
//					j3Record, errX := product.QryJ3Record(vca)
//					if errX != nil {
//						// 查单有问题，直接订单要置为超时，消息置为处理完毕
//						global.GVA_LOG.Error("查询充值记录异常", zap.Error(errX))
//
//						if errX = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
//							global.GVA_LOG.Error("更新订单异常", zap.Error(errX))
//							_ = msg.Reject(false)
//							continue
//						}
//
//						//入库操作记录
//						record := sysModel.SysOperationRecord{
//							Ip:      v.Ctx.ClientIP,
//							Method:  v.Ctx.Method,
//							Path:    v.Ctx.UrlPath,
//							Agent:   v.Ctx.UserAgent,
//							Status:  500,
//							Latency: time.Since(time.Now()),
//							Resp:    fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//							UserID:  v.Ctx.UserID,
//						}
//
//						err = operationRecordService.CreateSysOperationRecord(record)
//						if err != nil {
//							global.GVA_LOG.Error("当前账号回调查单，官方查单异常，record 入库失败..." + err.Error())
//						}
//
//						err = global.GVA_DB.Debug().Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).
//							Update("sys_status", 0).Error
//
//						vca.Status = 0
//						oc := vboxReq.ChanAccAndCtx{
//							Obj: vca,
//							Ctx: vboxReq.Context{
//								Body:      fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//								ClientIP:  "127.0.0.1",
//								Method:    "POST",
//								UrlPath:   "/vca/switchEnable",
//								UserAgent: "",
//								UserID:    int(vca.CreatedBy),
//							},
//						}
//						marshal, _ := json.Marshal(oc)
//
//						_ = chX.Publish(ChanAccEnableCheckExchange, ChanAccEnableCheckKey, marshal)
//
//						global.GVA_LOG.Error("处理订单为失败，发起一条关号清理资源", zap.Any("orderID", v.Obj.OrderId), zap.Any("ac_id", vca.AcId), zap.Any("ac_account", vca.AcAccount))
//
//						_ = msg.Ack(true)
//						continue
//					}
//
//					nowBalance := j3Record.LeftCoins
//					J3AccBalanceKey := fmt.Sprintf(global.J3AccBalanceZSet, vca.AcId)
//					nowTimeUnix := odDB.CreatedAt.Unix()
//					//对账时间
//					checkTime := time.Now().Unix()
//					valMembers, errZ := global.GVA_REDIS.ZRangeByScore(context.Background(), J3AccBalanceKey, &redis.ZRangeBy{
//						Min:    strconv.FormatInt(nowTimeUnix, 10),
//						Max:    strconv.FormatInt(nowTimeUnix, 10),
//						Offset: 0,
//						Count:  1,
//					}).Result()
//					if len(valMembers) > 0 {
//						mem := valMembers[0]
//						split := strings.Split(mem, ",")
//						hisBalance, _ := strconv.Atoi(split[4])
//						if hisBalance+money*100 == nowBalance { // 充值成功的情况
//							qryURL := vca.Token
//
//							keyMem := fmt.Sprintf("%s,%s,%v,%d,%d,%d,%d", v.Obj.OrderId, vca.AcAccount, money, nowTimeUnix, hisBalance, checkTime, nowBalance)
//							delMem := fmt.Sprintf("%s,%s,%v,%d,%d,%d,%d", v.Obj.OrderId, vca.AcAccount, money, nowTimeUnix, hisBalance, 0, 0)
//
//							global.GVA_LOG.Info("查单链接", zap.Any("订单ID", v.Obj.OrderId), zap.Any("url", qryURL))
//							global.GVA_LOG.Info("核对官方订单", zap.Any("核准", keyMem))
//
//							global.GVA_REDIS.ZAdd(context.Background(), J3AccBalanceKey, redis.Z{
//								Score:  float64(nowTimeUnix),
//								Member: keyMem,
//							})
//							global.GVA_REDIS.ZRem(context.Background(), J3AccBalanceKey, delMem)
//
//							v.Obj.PlatId = keyMem
//
//							//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
//							if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 1).
//								Update("plat_id", keyMem).Error; err != nil {
//								global.GVA_LOG.Error("更新订单异常", zap.Error(err))
//								_ = msg.Reject(false)
//								continue
//							}
//
//							// 4.入库wallet
//							var count int64
//							global.GVA_DB.Model(&vbox.UserWallet{}).Where("event_id = ?", v.Obj.OrderId).Count(&count)
//
//							if count == 0 {
//								wallet := vbox.UserWallet{
//									Uid:       v.Obj.CreatedBy,
//									CreatedBy: v.Obj.CreatedBy,
//									Type:      global.WalletOrderType,
//									EventId:   v.Obj.OrderId,
//									Recharge:  -money,
//									Remark:    fmt.Sprintf(global.WalletEventOrderCost, money, v.Obj.ChannelCode, v.Obj.OrderId),
//								}
//
//								global.GVA_DB.Model(&vbox.UserWallet{}).Save(&wallet)
//							}
//
//							_ = msg.Ack(true)
//							global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("pa", v.Obj.PAccount), zap.Any("orderId", v.Obj.OrderId))
//
//							// 同时把订单 redis信息也设置一下缓存信息
//							v.Obj.OrderStatus = 1
//							odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
//							jsonString, _ := json.Marshal(v.Obj)
//							global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)
//
//							// 并且发起一个回调通知的消息
//							marshal, _ := json.Marshal(v)
//							err = chX.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
//							global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))
//
//							accJucKey := fmt.Sprintf(global.PayAccKey, accID)
//							//支付成功后，释放本号
//							waitAccYdKey := fmt.Sprintf(global.YdJ3AccWaiting, accID)
//
//							waitAccMem := fmt.Sprintf("%v,%s,%s", vca.ID, vca.AcId, vca.AcAccount)
//							waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
//							global.GVA_REDIS.Del(context.Background(), accJucKey)
//							global.GVA_REDIS.Del(context.Background(), waitAccYdKey)
//
//							_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), 0)
//							global.GVA_LOG.Info("释放本账号", zap.Any("acAccount", vca.AcAccount), zap.Any("acId", accID), zap.Any("orderId", odDB.OrderId))
//
//							_ = msg.Ack(true)
//							continue
//						} else {
//							global.GVA_LOG.Info("未对账成功，当前余额为", zap.Any("nowBalance", nowBalance), zap.Any("hisBalance", hisBalance), zap.Any("money", money), zap.Any("orderId", v.Obj.OrderId))
//							// 重新丢回去 下一个20s再查一次
//							marshal, _ := json.Marshal(v)
//							err = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 15*time.Second)
//							_ = msg.Ack(true)
//							continue
//						}
//
//					} else {
//						// 异常情况处理
//						global.GVA_LOG.Error("订单匹配消息查redis数据异常", zap.Error(errZ), zap.Any("valMembers", valMembers))
//						// 如果解析消息失败，则直接丢弃消息
//						_ = msg.Reject(false)
//						continue
//					}
//
//				} else if global.QnContains(chanID) {
//				} else if global.PcContains(chanID) {
//					global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", odCreatedTime), zap.Any("传入的过期时间", *expTime))
//
//					records, errQy := product.QryQQRecordsByID(vca, odDB.PlatId)
//					if errQy != nil {
//						// 查单有问题，直接订单要置为超时，消息置为处理完毕
//						global.GVA_LOG.Error("查询充值记录异常", zap.Error(errQy))
//
//						if errQy = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
//							global.GVA_LOG.Error("更新订单异常", zap.Error(errQy))
//							_ = msg.Reject(false)
//							continue
//						}
//
//						//入库操作记录
//						record := sysModel.SysOperationRecord{
//							Ip:      v.Ctx.ClientIP,
//							Method:  v.Ctx.Method,
//							Path:    v.Ctx.UrlPath,
//							Agent:   v.Ctx.UserAgent,
//							Status:  500,
//							Latency: time.Since(time.Now()),
//							Body:    fmt.Sprintf(global.AccRecordEx, vca.AcId),
//							Resp:    fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//							UserID:  v.Ctx.UserID,
//						}
//
//						err = operationRecordService.CreateSysOperationRecord(record)
//						if err != nil {
//							global.GVA_LOG.Error("当前账号回调查单，官方查单异常，record 入库失败..." + err.Error())
//						}
//
//						err = global.GVA_DB.Debug().Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).
//							Update("sys_status", 0).Error
//
//						vca.Status = 0
//						oc := vboxReq.ChanAccAndCtx{
//							Obj: vca,
//							Ctx: vboxReq.Context{
//								Body:      fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
//								ClientIP:  "127.0.0.1",
//								Method:    "POST",
//								UrlPath:   "/vca/switchEnable",
//								UserAgent: "",
//								UserID:    int(vca.CreatedBy),
//							},
//						}
//						marshal, _ := json.Marshal(oc)
//
//						_ = chX.Publish(ChanAccEnableCheckExchange, ChanAccEnableCheckKey, marshal)
//
//						global.GVA_LOG.Error("处理订单为失败，发起一条关号清理资源", zap.Any("orderID", v.Obj.OrderId), zap.Any("ac_id", vca.AcId), zap.Any("ac_account", vca.AcAccount))
//
//						_ = msg.Ack(true)
//						continue
//					}
//					rdMap := product.Classifier(records.WaterList)
//					if vm, ok := rdMap["Q币"]; !ok {
//						global.GVA_LOG.Info("还没有QB的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", v.Obj.OrderId))
//					} else {
//						if rd, ok2 := vm[strconv.FormatInt(int64(money*100), 10)]; !ok2 {
//							global.GVA_LOG.Info("还没有QB的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", v.Obj.OrderId))
//
//						} else { // 证明这种金额的，充上了
//							if utils.Contains(rd, vca.AcAccount) {
//								//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
//								if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 1).Error; err != nil {
//									global.GVA_LOG.Error("更新订单异常", zap.Error(err))
//									_ = msg.Reject(false)
//									continue
//								}
//
//								// 4.入库wallet
//								var count int64
//								global.GVA_DB.Model(&vbox.UserWallet{}).Where("event_id = ?", v.Obj.OrderId).Count(&count)
//
//								if count == 0 {
//									wallet := vbox.UserWallet{
//										Uid:       v.Obj.CreatedBy,
//										CreatedBy: v.Obj.CreatedBy,
//										Type:      global.WalletOrderType,
//										EventId:   v.Obj.OrderId,
//										Recharge:  -money,
//										Remark:    fmt.Sprintf(global.WalletEventOrderCost, money, v.Obj.ChannelCode, v.Obj.OrderId),
//									}
//
//									global.GVA_DB.Model(&vbox.UserWallet{}).Save(&wallet)
//								}
//
//								_ = msg.Ack(true)
//								global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("orderId", v.Obj.OrderId))
//
//								// 同时把订单 redis信息也设置一下缓存信息
//								v.Obj.OrderStatus = 1
//								odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
//								jsonString, _ := json.Marshal(v.Obj)
//								global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)
//
//								// 并且发起一个回调通知的消息
//								marshal, _ := json.Marshal(v)
//								err = chX.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
//								global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))
//
//								_ = msg.Ack(true)
//								continue
//							}
//						}
//					}
//				}
//
//				if err != nil {
//					_ = msg.Reject(false)
//					continue
//				}
//
//				// 重新丢回去 下一个20s再查一次
//				marshal, err := json.Marshal(v)
//				if global.TxContains(v.Obj.ChannelCode) {
//					err = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 25*time.Second)
//				} else if global.J3Contains(v.Obj.ChannelCode) {
//					err = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 15*time.Second)
//				} else {
//					err = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 25*time.Second)
//				}
//			}
//			wg.Done()
//		}(i + 1)
//	}
//
//	// 等待所有消费者完成处理
//	wg.Wait()
//	//time.Sleep(time.Minute)
//	global.GVA_LOG.Info("MqOrderConfirmTask 初始化搞定")
//
//}
