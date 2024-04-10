package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	prod "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/flipped-aurora/gin-vue-admin/server/vbUtil"
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
	var operationRecordService system.OperationRecordService

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
			now := time.Now()
			connX, errX := mq.MQ.ConnPool.GetConnection()
			if errX != nil {
				//log.Fatalf("Failed to get connection from pool: %v", err)
				global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(errX))
			}
			defer mq.MQ.ConnPool.ReturnConnection(connX)
			if connX == nil {
				global.GVA_LOG.Error("connX == nil")
				return
			}
			chX, _ := connX.Channel()

			// 说明：执行查单回调处理
			deliveries, errC := chX.Consume(OrderConfirmDeadQueue, "", false, false, false, false, nil)
			if errC != nil {
				global.GVA_LOG.Error("mq 消费者异常， err", zap.Error(errC), zap.Any("queue", OrderConfirmDeadQueue))
			}

			for msg := range deliveries {

				v := request.PayOrderAndCtx{}
				errJ := json.Unmarshal(msg.Body, &v)
				nowTime := time.Now()
				if errJ != nil {
					global.GVA_LOG.Error("MqOrderConfirmTask...", zap.Error(errJ), zap.Any("错误消息体", msg.Body))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}
				orderId := v.Obj.OrderId
				global.GVA_LOG.Info("收到需要查询付款状态的订单消息", zap.Any("orderId", orderId))

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
					global.GVA_LOG.Info("该订单已经补单过，跳过", zap.Any("orderId", orderId))

					_ = msg.Ack(true)
					continue
				}

				//1. 筛选匹配是哪个产品
				chanID := v.Obj.ChannelCode
				accID := v.Obj.AcId
				var vca vbox.ChannelAccount
				errQ = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("ac_id = ?", accID).First(&vca).Error
				if errQ != nil {
					global.GVA_LOG.Error("订单匹配消息查库失败", zap.Error(errQ))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}

				orgTmp := utils2.GetSelfOrg(vca.CreatedBy)

				//2. 查询订单（账号）的充值情况
				expTime := v.Obj.ExpTime
				duration := expTime.Sub(nowTime)
				if duration > 0 {
					global.GVA_LOG.Info("该订单还没过期，继续查", zap.Any("orderID", orderId))
				} else {
					global.GVA_LOG.Info("该订单已经过期", zap.Any("orderID", orderId))
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
					odKey := fmt.Sprintf(global.PayOrderKey, orderId)
					jsonString, _ := json.Marshal(v.Obj)
					global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)

					record := sysModel.SysOperationRecord{
						Ip:      v.Ctx.ClientIP,
						Method:  v.Ctx.Method,
						Path:    v.Ctx.UrlPath,
						Agent:   v.Ctx.UserAgent,
						MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
						Type:    global.OrderType,
						Body:    "",
						Status:  200,
						Latency: time.Since(now),
						Resp:    fmt.Sprintf(global.OrderTimeoutMsg),
						UserID:  v.Ctx.UserID,
					}
					err = operationRecordService.CreateSysOperationRecord(record)
					continue
				}

				money := v.Obj.Money
				odCreatedTime := *odDB.CreatedAt

				prodInfo, errJ := vbUtil.GetProductByCode(chanID)
				t, errJ := vbUtil.GetCDTimeByCode(chanID)
				if errJ != nil {
					global.GVA_LOG.Error("获取产品信息异常", zap.Error(errJ))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}

				startTime := odCreatedTime.Add(-t.PreTime)
				endTime := *expTime
				global.GVA_LOG.Info("传入的查询时间范围", zap.Any("startTime", startTime), zap.Any("endTime", endTime), zap.Any("orderID", orderId))
				var qryURL, platID string
				var flag bool
				switch {
				case strings.Contains(prodInfo.ProductId, "qb") || strings.Contains(prodInfo.ProductId, "dnf"):
					var waterList []prod.Payment
					if global.PcContains(vca.Cid) {
						global.GVA_LOG.Info("pc 核对订单", zap.Any("od plat", v.Obj.PlatId))
						var recordList *prod.Records
						recordList, errX = product.QryQQRecordsByID(vca, v.Obj.PlatId)
						if errX != nil {
							global.GVA_LOG.Error("获取订单信息异常", zap.Error(errX))
							break
						}
						waterList = recordList.WaterList
					} else {
						waterList, errX, qryURL = product.QryQQRecordsBetween(vca, startTime, endTime)
						if errX != nil {
							global.GVA_LOG.Error("获取订单信息异常", zap.Error(errX))
							break
						}
					}

					if len(waterList) <= 0 {
						global.GVA_LOG.Info("waterList 0,还没有QB的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", orderId))
					} else {
						rdMap := product.Classifier(waterList)
						if vm, ok := rdMap["Q币"]; !ok {
							global.GVA_LOG.Info("map[qb] 0,还没有QB的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", orderId))
						} else {
							if rd, ok2 := vm[strconv.FormatInt(int64(money*100), 10)]; !ok2 {
								global.GVA_LOG.Info("map[qb] 1,还没有QB的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", orderId))

							} else { // 证明这种金额的，充上了
								for _, tg := range rd {
									if strings.Contains(tg, vca.AcAccount) {
										flag = true
										platID = strings.Split(tg, ",")[1]
										break
									}
								}
								if flag {
									global.GVA_LOG.Info("查单链接", zap.Any("订单ID", orderId), zap.Any("url", qryURL))
									global.GVA_LOG.Info("核对官方订单", zap.Any("platID", platID), zap.Any("订单ID", orderId), zap.Any("查单起始时间", startTime), zap.Any("结束时间", endTime))
									v.Obj.PlatId = platID
								}
							}
						}
					}

				case strings.Contains(prodInfo.ProductId, "j3"):
					var j3Record *prod.J3BalanceData
					j3Record, errX = product.QryJ3Record(vca)
					if j3Record != nil {
						nowBalance := j3Record.LeftCoins
						J3AccBalanceKey := fmt.Sprintf(global.J3AccBalanceZSet, vca.AcId)
						nowTimeUnix := odDB.CreatedAt.Unix()
						//对账时间
						checkTime := time.Now().Unix()
						valMembers, errZ := global.GVA_REDIS.ZRangeByScore(context.Background(), J3AccBalanceKey, &redis.ZRangeBy{
							Min:    strconv.FormatInt(nowTimeUnix, 10),
							Max:    strconv.FormatInt(nowTimeUnix, 10),
							Offset: 0,
							Count:  1,
						}).Result()
						if errZ != nil {
							global.GVA_LOG.Error("redis err", zap.Error(errZ))
						}
						if len(valMembers) > 0 {
							mem := valMembers[0]
							split := strings.Split(mem, ",")
							hisBalance, _ := strconv.Atoi(split[4])
							if hisBalance+money*100 == nowBalance { // 充值成功的情况
								qryURL = vca.Token

								keyMem := fmt.Sprintf("%s,%s,%v,%d,%d,%d,%d", orderId, vca.AcAccount, money, nowTimeUnix, hisBalance, checkTime, nowBalance)
								delMem := fmt.Sprintf("%s,%s,%v,%d,%d,%d,%d", orderId, vca.AcAccount, money, nowTimeUnix, hisBalance, 0, 0)

								global.GVA_LOG.Info("查单链接", zap.Any("订单ID", orderId), zap.Any("url", qryURL))
								global.GVA_LOG.Info("核对官方订单", zap.Any("核准", keyMem))

								global.GVA_REDIS.ZAdd(context.Background(), J3AccBalanceKey, redis.Z{
									Score:  float64(nowTimeUnix),
									Member: keyMem,
								})
								global.GVA_REDIS.ZRem(context.Background(), J3AccBalanceKey, delMem)

								v.Obj.PlatId = keyMem

								flag = true
							}
						}
					}

				case strings.Contains(prodInfo.ProductId, "wy"):
					var wyRecord *prod.WYBalanceData
					wyRecord, errX = product.QryWYRecord(vca.Token)
					if wyRecord != nil {
						nowBalance := wyRecord.JsBalance
						WYAccBalanceKey := fmt.Sprintf(global.WYAccBalanceZSet, vca.AcId)
						nowTimeUnix := odDB.CreatedAt.Unix()
						//对账时间
						checkTime := time.Now().Unix()
						valMembers, errZ := global.GVA_REDIS.ZRangeByScore(context.Background(), WYAccBalanceKey, &redis.ZRangeBy{
							Min:    strconv.FormatInt(nowTimeUnix, 10),
							Max:    strconv.FormatInt(nowTimeUnix, 10),
							Offset: 0,
							Count:  1,
						}).Result()
						if errZ != nil {
							global.GVA_LOG.Error("redis err", zap.Error(errZ))
						}
						if len(valMembers) > 0 {
							mem := valMembers[0]
							split := strings.Split(mem, ",")
							hisBalance, _ := strconv.Atoi(split[4])
							if hisBalance+money*10 == nowBalance { // 充值成功的情况
								qryURL = vca.Token

								keyMem := fmt.Sprintf("%s,%s,%v,%d,%d,%d,%d", orderId, vca.AcAccount, money, nowTimeUnix, hisBalance, checkTime, nowBalance)
								delMem := fmt.Sprintf("%s,%s,%v,%d,%d,%d,%d", orderId, vca.AcAccount, money, nowTimeUnix, hisBalance, 0, 0)

								global.GVA_LOG.Info("核对官方订单", zap.Any("订单ID", orderId), zap.Any("核准", keyMem), zap.Any("money", money))

								global.GVA_REDIS.ZAdd(context.Background(), WYAccBalanceKey, redis.Z{
									Score:  float64(nowTimeUnix),
									Member: keyMem,
								})
								global.GVA_REDIS.ZRem(context.Background(), WYAccBalanceKey, delMem)

								v.Obj.PlatId = keyMem

								flag = true
							}
						}
					}
				case strings.Contains(prodInfo.ProductId, "db"):
					var dyRecord *prod.DyWalletInfoRecord
					dyRecord, errX = product.QryDyRecord(vca.Token)
					if dyRecord != nil {
						nowBalance := dyRecord.Diamond
						DyAccBalanceKey := fmt.Sprintf(global.DyAccBalanceZSet, vca.AcId)
						nowTimeUnix := odDB.CreatedAt.Unix()
						//对账时间
						checkTime := time.Now().Unix()
						valMembers, errZ := global.GVA_REDIS.ZRangeByScore(context.Background(), DyAccBalanceKey, &redis.ZRangeBy{
							Min:    strconv.FormatInt(nowTimeUnix, 10),
							Max:    strconv.FormatInt(nowTimeUnix, 10),
							Offset: 0,
							Count:  1,
						}).Result()
						if errZ != nil {
							global.GVA_LOG.Error("redis err", zap.Error(errZ))
						}
						if len(valMembers) > 0 {
							mem := valMembers[0]
							split := strings.Split(mem, ",")
							hisBalance, _ := strconv.Atoi(split[4])
							global.GVA_LOG.Info("查单金额", zap.Any("应为金额", hisBalance+money*10), zap.Any("实际金额", nowBalance))
							if hisBalance+money*10 == nowBalance { // 充值成功的情况
								qryURL = vca.Token

								keyMem := fmt.Sprintf("%s,%s,%v,%d,%d,%d,%d", orderId, vca.AcAccount, money, nowTimeUnix, hisBalance, checkTime, nowBalance)
								delMem := fmt.Sprintf("%s,%s,%v,%d,%d,%d,%d", orderId, vca.AcAccount, money, nowTimeUnix, hisBalance, 0, 0)

								//global.GVA_LOG.Info("查单链接", zap.Any("订单ID", orderId), zap.Any("url", qryURL))
								global.GVA_LOG.Info("核对官方订单", zap.Any("订单ID", orderId), zap.Any("核准", keyMem), zap.Any("money", money))

								global.GVA_REDIS.ZAdd(context.Background(), DyAccBalanceKey, redis.Z{
									Score:  float64(nowTimeUnix),
									Member: keyMem,
								})
								global.GVA_REDIS.ZRem(context.Background(), DyAccBalanceKey, delMem)

								v.Obj.PlatId = keyMem

								flag = true
							}
						}
					}
				case strings.Contains(prodInfo.ProductId, "dnf"):
					var waterList []prod.Payment
					waterList, errX, qryURL = product.QryQQRecordsBetween(vca, startTime, endTime)
					if len(waterList) <= 0 {
						global.GVA_LOG.Info("waterList 0,还没有DNF的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", orderId))
					} else {
						rdMap := product.Classifier(waterList)
						if vm, ok := rdMap["DNF"]; !ok {
							global.GVA_LOG.Info("map[DNF] 0,还没有DNF的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", orderId))
						} else {
							if rd, ok2 := vm[strconv.FormatInt(int64(money*100), 10)]; !ok2 {
								global.GVA_LOG.Info("map[DNF] 1,还没有DNF的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", orderId))

							} else { // 证明这种金额的，充上了

								for _, tg := range rd {
									if strings.Contains(tg, vca.AcAccount) {
										flag = true
										platID = strings.Split(tg, ",")[1]
										break
									}
								}
								if flag {
									global.GVA_LOG.Info("查单链接", zap.Any("订单ID", orderId), zap.Any("url", qryURL))
									global.GVA_LOG.Info("核对官方订单", zap.Any("platID", platID), zap.Any("订单ID", orderId), zap.Any("查单起始时间", startTime), zap.Any("结束时间", endTime))
									v.Obj.PlatId = platID
								}
							}
						}
					}
				case strings.Contains(prodInfo.ProductId, "sdo"):
					var waterList []prod.SdoDaoYuOrderRecord
					waterList, errX = product.QrySdoDaoYuRecordBetween(vca, startTime, endTime)
					if len(waterList) <= 0 {
						global.GVA_LOG.Info("waterList 0,还没有Sdo的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", orderId))
					} else {
						rdMap := product.Classifier(waterList)
						global.GVA_LOG.Info("筛后的map", zap.Any("map", rdMap), zap.Any("waterList", waterList))
						if vm, ok := rdMap["彩虹岛"]; !ok {
							global.GVA_LOG.Info("还没有Sdo的充值记录，空空如也")
						} else {
							runStr := string(rune(money))
							formatMoney := strconv.FormatInt(int64(money), 10)
							//打印两个 money转化
							global.GVA_LOG.Info("money", zap.Any("runStr", runStr), zap.Any("formatMoney", formatMoney))
							if rd, ok2 := vm[formatMoney]; !ok2 {
								global.GVA_LOG.Info("还没有Sdo的充值记录")
							} else { // 证明这种金额的，充上了
								for _, tg := range rd {
									if strings.Contains(tg, vca.AcAccount) {
										flag = true
										platID = strings.Split(tg, "_")[1]
										v.Obj.PlatId = platID
										break
									}
								}
							}
						}
					}
				case strings.Contains(prodInfo.ProductId, "qn"):
					var waterList []prod.QNRecord
					waterList, errX = product.QryQNRecords(vca, startTime, endTime, v.Obj.PlatId)
					if len(waterList) <= 0 {
						global.GVA_LOG.Info("waterList 0,还没有qn shop的充值记录", zap.Any("查询对应账号ID", accID), zap.Any("查询对应订单", orderId))
					} else {
						global.GVA_LOG.Info("查询sku Title", zap.Any("water list len", len(waterList)), zap.Any("markID", v.Obj.PlatId))
						shopInfo := waterList[0]
						if shopInfo.Money == strconv.FormatInt(int64(money), 10) {
							if flag {
								global.GVA_LOG.Info("核对官方订单", zap.Any("markID", shopInfo.SkuTitle), zap.Any("订单ID", orderId),
									zap.Any("plat money", shopInfo.Money), zap.Any("od money", money), zap.Any("查单起始时间", startTime), zap.Any("结束时间", endTime))
								v.Obj.PlatId = platID
							}
						}
					}
				case strings.Contains(prodInfo.ProductId, "card"):
					//库记录的platID是否有值
					if odDB.PlatId == "" {
						global.GVA_LOG.Info("卡密尚未提交，继续查询", zap.Any("money", money), zap.Any("orderId", orderId))
						// 重新丢回去 下一个20s再查一次
						marshal, _ := json.Marshal(v)
						err = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 15*time.Second)
						_ = msg.Ack(true)
						continue
					} else {
						v.Obj.PlatId = odDB.PlatId

						//获取查卡池ck
						cardAcc, errCK := vbUtil.GetOrgPoolCK(orgTmp[0], chanID)
						if errCK != nil {
							global.GVA_LOG.Error("获取查单池卡密ck异常", zap.Error(errCK))
						}

						//a.拿账户的ck查一遍卡密，如果卡密ok则执行绑卡
						errE := product.ECardQuery(odDB.PlatId, cardAcc.Token)
						if errE != nil {
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
								Type:    global.OrderType,
								Status:  500,
								Latency: time.Since(time.Now()),
								Resp:    fmt.Sprintf(global.OrderConfirmBindQryRetMsg, v.Obj.PlatId, errE.Error()),
								UserID:  v.Ctx.UserID,
							}
							err = operationRecordService.CreateSysOperationRecord(record)

							global.GVA_LOG.Error("匹配查单池账号异常,查询卡密合法性错误", zap.Error(errE))
							if errX = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
								global.GVA_LOG.Error("更新订单异常", zap.Error(errX))
								_ = msg.Reject(false)
								continue
							}
							_ = msg.Ack(true)
							continue
						}

						// 释放查单池ck账号
						cdTime := 10 * time.Second
						accWaitYdKey := fmt.Sprintf(global.YdECPoolAccWaiting, cardAcc.AcId)
						accInfoVal := fmt.Sprintf("%d,%s,%s", cardAcc.ID, cardAcc.AcId, cardAcc.AcAccount)

						global.GVA_REDIS.Set(context.Background(), accWaitYdKey, accInfoVal, cdTime)

						waitMsg := strings.Join([]string{accWaitYdKey, accInfoVal}, "-")

						err = chX.PublishWithDelay(CardAccCDCheckDelayedExchange, CardAccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)

						ecRecord, errW := product.ECardBind(odDB.PlatId, vca.Token)
						if errW != nil {
							// 查单有问题，直接订单要置为超时，消息置为处理完毕
							global.GVA_LOG.Error("查询充值记录异常", zap.Error(errW))

							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
								Type:    global.OrderType,
								Status:  500,
								Latency: time.Since(time.Now()),
								Resp:    fmt.Sprintf(global.OrderConfirmBindRetMsg, v.Obj.PlatId, errW.Error()),
								UserID:  v.Ctx.UserID,
							}
							err = operationRecordService.CreateSysOperationRecord(record)

							if errX = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
								global.GVA_LOG.Error("更新订单异常", zap.Error(errX))
								_ = msg.Reject(false)
								continue
							}
						}
						if ecRecord != nil { //绑卡成功
							global.GVA_LOG.Info("绑卡成功，核对卡密金额信息", zap.Any("orderId", odDB.OrderId), zap.Any("money", money), zap.Any("ecRecord", ecRecord))
							amount := ecRecord["OrderMoney"].(float64)
							if amount == float64(money) {
								record := sysModel.SysOperationRecord{
									Ip:      v.Ctx.ClientIP,
									Method:  v.Ctx.Method,
									Path:    v.Ctx.UrlPath,
									Agent:   v.Ctx.UserAgent,
									MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
									Type:    global.OrderType,
									Status:  200,
									Latency: time.Since(time.Now()),
									Resp:    fmt.Sprintf(global.OrderConfirmBindMsg, v.Obj.PlatId),
									UserID:  v.Ctx.UserID,
								}
								err = operationRecordService.CreateSysOperationRecord(record)

								flag = true
							} else {
								record := sysModel.SysOperationRecord{
									Ip:      v.Ctx.ClientIP,
									Method:  v.Ctx.Method,
									Path:    v.Ctx.UrlPath,
									Agent:   v.Ctx.UserAgent,
									MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
									Type:    global.OrderType,
									Status:  200,
									Latency: time.Since(time.Now()),
									Resp:    fmt.Sprintf(global.OrderConfirmBindOtherMoneyMsg, v.Obj.PlatId),
									UserID:  v.Ctx.UserID,
								}

								err = operationRecordService.CreateSysOperationRecord(record)
							}
						}
					}
				default:
					global.GVA_LOG.Error("未知产品类型", zap.Any("订单ID", orderId), zap.Any("产品类型", prodInfo.ProductId))
					errX = fmt.Errorf("未知产品类型 %v", prodInfo.ProductId)
				}

				if errX != nil {
					// 查单有问题，直接订单要置为超时，消息置为处理完毕
					global.GVA_LOG.Error("查询充值记录异常", zap.Error(errX), zap.Any("orderID", orderId))
					es := errX.Error()
					if strings.Contains(es, "timed out") {
						global.GVA_LOG.Warn("查官方记录超时，重试一次...")
						// 重新丢回去 下一个20s再查一次
						marshal, _ := json.Marshal(v)
						_ = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 1*time.Second)
						_ = msg.Ack(true)
						continue
					}

					if errX = global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; err != nil {
						global.GVA_LOG.Error("更新订单异常", zap.Error(errX))
						_ = msg.Reject(false)
						continue
					}

					//入库操作记录
					record := sysModel.SysOperationRecord{
						Ip:      v.Ctx.ClientIP,
						Method:  v.Ctx.Method,
						Path:    v.Ctx.UrlPath,
						Agent:   v.Ctx.UserAgent,
						MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
						Type:    global.OrderType,
						Status:  500,
						Latency: time.Since(time.Now()),
						Resp:    fmt.Sprintf(global.OrderConfirmErrMsg, vca.AcAccount, vca.AcId),
						UserID:  v.Ctx.UserID,
					}

					err = operationRecordService.CreateSysOperationRecord(record)
					if err != nil {
						global.GVA_LOG.Error("当前账号回调查单，官方查单异常，record 入库失败..." + err.Error())
					}

					err = global.GVA_DB.Debug().Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", vca.ID).
						Update("sys_status", 0).Error

					vca.Status = 0
					oc := vboxReq.ChanAccAndCtx{
						Obj: vca,
						Ctx: vboxReq.Context{
							Body:      fmt.Sprintf(global.AccQryEx, vca.AcId, vca.AcAccount),
							ClientIP:  "127.0.0.1",
							Method:    "POST",
							UrlPath:   "/vca/switchEnable",
							UserAgent: "",
							UserID:    int(vca.CreatedBy),
						},
					}
					marshal, _ := json.Marshal(oc)

					_ = chX.Publish(ChanAccEnableCheckExchange, ChanAccEnableCheckKey, marshal)

					global.GVA_LOG.Error("处理订单为失败，发起一条关号清理资源", zap.Any("orderID", orderId), zap.Any("ac_id", vca.AcId), zap.Any("ac_account", vca.AcAccount))

					_ = msg.Ack(true)
					continue
				}

				if flag {
					//3. 查询充值成功后，更新订单信息（订单状态，订单支付链接处理）
					if err := global.GVA_DB.Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 1).Error; err != nil {
						global.GVA_LOG.Error("更新订单异常", zap.Error(err))
						_ = msg.Reject(false)
						continue
					}
					global.GVA_LOG.Info("订单查到已支付并确认消息消费，更新订单状态", zap.Any("orderId", orderId))

					// 同时把订单 redis信息也设置一下缓存信息
					v.Obj.OrderStatus = 1
					odKey := fmt.Sprintf(global.PayOrderKey, orderId)
					jsonString, _ := json.Marshal(v.Obj)
					global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)

					//支付成功后，释放本号
					//accWaitYdKey := fmt.Sprintf(global.YdProdMoneyAccWaiting, prodInfo.ProductId, accID, money)
					//accInfoVal := fmt.Sprintf("%d,%s,%s,%v", vca.ID, accID, vca.AcAccount, money)

					var accWaitYdKey, accInfoVal string
					switch {
					case strings.Contains(prodInfo.Ext, "money"):
						accWaitYdKey = fmt.Sprintf(global.ProdAccMoneyWaiting, prodInfo.ProductId, accID, money)
						accInfoVal = fmt.Sprintf("%v,%s,%s,%v", vca.ID, accID, vca.AcAccount, money)
					default:
						accWaitYdKey = fmt.Sprintf(global.ProdAccWaiting, prodInfo.ProductId, accID)
						accInfoVal = fmt.Sprintf("%v,%s,%s", vca.ID, accID, vca.AcAccount)
					}

					// 设置一个释放时间
					cdTime := t.PreTime
					global.GVA_LOG.Info(fmt.Sprintf("订单中的账号%v秒后发起释放通知", cdTime), zap.Any("acAccount", vca.AcAccount), zap.Any("acId", accID), zap.Any("order ID", orderId))
					global.GVA_REDIS.Set(context.Background(), accWaitYdKey, accInfoVal, cdTime)

					waitMsg := strings.Join([]string{accWaitYdKey, accInfoVal}, "-")
					err = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)

					record := sysModel.SysOperationRecord{
						Ip:      v.Ctx.ClientIP,
						Method:  v.Ctx.Method,
						Path:    v.Ctx.UrlPath,
						Agent:   v.Ctx.UserAgent,
						MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
						Type:    global.OrderType,
						Status:  200,
						Latency: time.Since(time.Now()),
						Resp:    fmt.Sprintf(global.OrderConfirmMsg, v.Obj.PlatId),
						UserID:  v.Ctx.UserID,
					}

					err = operationRecordService.CreateSysOperationRecord(record)

					// 并且发起一个回调通知的消息
					marshal, _ := json.Marshal(v)
					err = chX.Publish(OrderCallbackExchange, OrderCallbackKey, marshal)
					global.GVA_LOG.Info("【系统自动】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", orderId))

					_ = msg.Ack(true)
					continue
				}

				if err != nil {
					_ = msg.Reject(false)
					continue
				}

				// 重新丢回去 下一个20s再查一次
				marshal, _ := json.Marshal(v)
				if global.TxContains(v.Obj.ChannelCode) {
					_ = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 25*time.Second)
				} else if global.J3Contains(v.Obj.ChannelCode) {
					_ = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 15*time.Second)
				} else if global.DyContains(v.Obj.ChannelCode) {
					_ = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 15*time.Second)
				} else {
					_ = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 25*time.Second)
				}
			}
			wg.Done()
		}(i + 1)
	}

	// 等待所有消费者完成处理
	wg.Wait()
	//time.Sleep(time.Minute)
	global.GVA_LOG.Info("MqOrderConfirmTask 初始化搞定")

}
