package task

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxRep "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/model"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"strconv"
	"strings"
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

	OrderCallbackExchange = "vbox.order.callback_exchange"
	OrderCallbackQueue    = "vbox.order.callback_queue"
	OrderCallbackKey      = "vbox.order.callback"
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
	if err = ch.ExchangeDeclare(OrderWaitExchange, "direct"); err != nil {
		global.GVA_LOG.Error("OrderWaitExchange.create exchange err 111:", zap.Any("err", err))
	}
	if err = ch.QueueDeclare(OrderWaitQueue); err != nil {
		global.GVA_LOG.Error("OrderWaitExchange.create queue err 111:", zap.Any("err", err))
	}
	if err = ch.QueueBind(OrderWaitQueue, OrderWaitKey, OrderWaitExchange); err != nil {
		global.GVA_LOG.Error("OrderWaitExchange.bind queue err 111:", zap.Any("err", err))
	}

	// 设置初始消费者数量
	consumerCount := 20
	// 使用 WaitGroup 来等待所有消费者完成处理
	var wg sync.WaitGroup
	wg.Add(consumerCount)

	// 启动多个消费者
	for i := 0; i < consumerCount; i++ {
		go func(consumerID int) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("捕获到异常:", r)
				}
			}()
			var operationRecordService system.OperationRecordService
			now := time.Now()
			// 说明：执行账号匹配
			deliveries, errQ := ch.Consume(OrderWaitQueue, "", false, false, false, false, nil)
			if errQ != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", OrderWaitQueue))
			}

			for msg := range deliveries {
				//v1 := &map[string]interface{}{}
				//err := json.Unmarshal(msg.Body, v1)

				v := request.PayOrderAndCtx{}
				err = json.Unmarshal(msg.Body, &v)
				if err != nil {
					global.GVA_LOG.Error("订单匹配消息处理失败，MqOrderWaitingTask...", zap.Any("err", err.Error()))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}
				global.GVA_LOG.Info("收到一条需要进行匹配的订单消息", zap.Any("order ID", v.Obj.OrderId), zap.Any("order", v))

				//1. 筛选匹配是哪个产品

				//1.0 核查商户
				var vpa vbox.PayAccount
				if strings.Contains(v.Obj.PAccount, "TEST") {
					global.GVA_LOG.Info("测试单，商户检测跳过", zap.Any("入参商户", v.Obj.PAccount))
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
					if err != nil {
						global.GVA_LOG.Error("订单匹配消息处理失败，MqOrderWaitingTask...", zap.Any("err", err.Error()))
						// 如果解析消息失败，则直接丢弃消息
						_ = msg.Reject(false)
						continue
					}
				}

				//2. 查供应库存账号是否充足 (优先从缓存池取，取空后查库取，如果库也空了，咋报错库存不足)

				//2.0 查一下单，如果发起初始化的时候已经匹配过账号了，就不用再匹配一次了
				var ID, accID, acAccount string
				var imgContent, MID string
				//获取当前组织
				orgTmp := utils2.GetSelfOrg(vpa.Uid)
				orgID := orgTmp[0]
				cid := v.Obj.ChannelCode

				// 设置一个冷却时间
				duration, _ := HandleExpTime2Product(cid)
				cdTime := duration + 60*time.Second

				if v.Obj.AcId != "" {
					accID = v.Obj.AcId
				} else {
					if global.TxContains(cid) {
						accKey := fmt.Sprintf(global.ChanOrgAccZSet, orgID, cid, v.Obj.Money)

						var resList []string
						resList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), accKey, &redis.ZRangeBy{
							Min:    "0",
							Max:    "0",
							Offset: 0,
							Count:  1,
						}).Result()

						if err != nil {
							global.GVA_LOG.Error("当前归属地区预产不足, redis err", zap.Error(err), zap.Any("db.region", v.Obj.PayRegion))
							// 如果解析消息失败，则直接丢弃消息
							_ = msg.Reject(false)
							continue
						}
						if resList != nil && len(resList) > 0 {
							accTmp := resList[0]

							// 2.1 把账号设置为已用
							global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{
								Score:  1,
								Member: accTmp,
							})

							// 2.2 把可用的账号给出来继续往下执行建单步骤
							split := strings.Split(accTmp, "_")
							ID = split[0]
							accID = split[1]
							acAccount = split[2]

							var vca vbox.ChannelAccount
							err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", ID).First(&vca).Error
							if err != nil {
								//return nil, fmt.Errorf("匹配通道账号不存在！ 请核查：%v", err.Error())

								global.GVA_LOG.Error("匹配账号异常", zap.Error(err))

								_ = msg.Ack(true)
								return
							}
							v.Obj.CreatedBy = vca.CreatedBy
							v.Obj.AcId = vca.AcId
							v.Ctx.UserID = int(vca.CreatedBy)

							global.GVA_LOG.Info("匹配账号", zap.Any("ID", ID), zap.Any("acID", accID), zap.Any("acAccount", vca.AcAccount))

						} else {
							global.GVA_LOG.Error("当前归属地区预产不足, list size zero", zap.Error(err), zap.Any("db.region", v.Obj.PayRegion))
							// 如果解析消息失败，则直接丢弃消息
							_ = msg.Reject(false)
							continue
						}

					} else if global.J3Contains(cid) {

					} else if global.PcContains(cid) {
						if v.Obj.PayRegion == "" {
							global.GVA_LOG.Error("当前未匹配到库中地区 err", zap.Error(err), zap.Any("db.region", v.Obj.PayRegion))
							// 如果解析消息失败，则直接丢弃消息
							_ = msg.Reject(false)
							continue
						}
						// 做一下ip访问匹配
						var province, ISP string
						if strings.Contains(v.Obj.PayRegion, "内网") {
							//	随机给一个北京的IP（电信|移动|联通）
							ISP = global.RandISP()
							province = "北京"
						} else {
							splitLoc := strings.Split(v.Obj.PayRegion, "|")
							province = splitLoc[2]
							ISP = global.ISPTranslate(splitLoc[4])
						}
						if !global.ISPContains(ISP) {
							//	随机把ISP赋值为 yidong|liantong|dianxin 其中一个
							ISP = global.RandISP()
						}
						if !global.ProvinceContains(province) {
							province = global.RandProvince()
						}
						global.GVA_LOG.Info("当前匹配到库中地区", zap.Any("db.region", v.Obj.PayRegion), zap.Any("province", province), zap.Any("ISP", ISP))

						// 查下省
						var geo model.Geo
						err = global.GVA_DB.Model(&model.Geo{}).Table("geo_provinces").
							Where("name LIKE ?", "%"+province+"%").First(&geo).Error
						if err != nil {
							global.GVA_LOG.Error("当前未匹配到库中地区 err", zap.Error(err), zap.Any("db.region", v.Obj.PayRegion))
							// 如果解析消息失败，则直接丢弃消息
							_ = msg.Reject(false)
							continue
						}

						ispPY := utils.ISP(ISP)
						provinceCode := geo.Code

						pcKey := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgID, cid, v.Obj.Money, ispPY, provinceCode)

						global.GVA_LOG.Info("pcKey", zap.Any("pcKey", pcKey))

						var resList []string
						resList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), pcKey, &redis.ZRangeBy{
							Min:    "0",
							Max:    "0",
							Offset: 0,
							Count:  1,
						}).Result()

						if err != nil {
							global.GVA_LOG.Error("当前归属地区预产不足, redis err", zap.Error(err), zap.Any("db.region", v.Obj.PayRegion))
							// 如果解析消息失败，则直接丢弃消息
							_ = msg.Reject(false)

							//入库操作记录
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.ResourceNotEnough),
								UserID:  v.Ctx.UserID,
							}
							err = operationRecordService.CreateSysOperationRecord(record)
							if err != nil {
								global.GVA_LOG.Error("record 入库失败..." + err.Error())
							}

							if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0); errDB != nil {
								global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(errDB.Error))
							}
							continue
						} else if resList != nil && len(resList) > 0 {
							pcTmp := resList[0]
							// 2.1 把预产设置为已用 ,置为1
							global.GVA_REDIS.ZAdd(context.Background(), pcKey, redis.Z{
								Score:  1,
								Member: pcTmp,
							})

							// 2.2 把可用的账号给出来继续往下执行建单步骤
							split := strings.Split(pcTmp, "_")
							ID = split[0]
							MID = split[1]
							acAccount = split[2]
							imgContent = split[3]

							var pcDB vbox.ChannelPayCode
							err = global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id = ?", ID).First(&pcDB).Error
							if err != nil {
								//return nil, fmt.Errorf("匹配通道账号不存在！ 请核查：%v", err.Error())

								global.GVA_LOG.Error("匹配预产异常", zap.Error(err))

								_ = msg.Ack(true)
								return
							}
							v.Obj.CreatedBy = pcDB.CreatedBy
							v.Obj.AcId = pcDB.AcId
							v.Obj.PlatId = pcDB.PlatId
							v.Ctx.UserID = int(pcDB.CreatedBy)

							err = global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id = ?", ID).
								Update("code_status", 1).Error

							global.GVA_LOG.Info("匹配预产", zap.Any("MID", MID), zap.Any("acAccount", acAccount), zap.Any("imgContent", imgContent))

							//匹配成功之后，要把当前acAccount相同金额的所有预产置为 等候状态，不允许被取用
							pattern := fmt.Sprintf(global.ChanOrgPayCodePrefix, orgID, cid, v.Obj.Money)
							keys, _ := global.GVA_REDIS.Keys(context.Background(), pattern).Result()

							var waitIDs []string
							for _, key := range keys {
								var resWaitTmpList []string
								resWaitTmpList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
									Min:    "0",
									Max:    "0",
									Offset: 0,
									Count:  -1,
								}).Result()

								for _, waitMem := range resWaitTmpList {
									if strings.Contains(waitMem, acAccount) {
										global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
											Score:  4,
											Member: waitMem,
										})

										//添加到 waitIDs 中
										id := strings.Split(waitMem, "_")[0]
										waitIDs = append(waitIDs, id)
									}
								}
							}

							// 把 pay code中除了本ID的，其它都让他进入冷却状态(包括对应通道账号)
							global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id in ? ", waitIDs).Update("code_status", 4)
							global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("ac_id in ? ", pcDB.AcId).
								Update("status", 4).Update("sys_status", 4)

							// 把当前acAccount下所有的预产等待队列置为冷却状态
							waitAccPcKey := fmt.Sprintf(global.AccWaiting, pcDB.AcId)
							waitIDsTmp := strings.Join(waitIDs, ",")
							global.GVA_REDIS.Set(context.Background(), waitAccPcKey, waitIDsTmp, cdTime)

							waitMsg := strings.Join([]string{waitAccPcKey, waitIDsTmp}, "_")
							err = ch.PublishWithDelay(PayCodeCDCheckDelayedExchange, PayCodeCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)

						} else {

							//给一次随机匹配其它地区的机会
							if len(resList) == 0 {
								pattern := fmt.Sprintf(global.ChanOrgPayCodePrefix, orgID, cid, v.Obj.Money)
								keys, errR := global.GVA_REDIS.Keys(context.Background(), pattern).Result()
								if errR != nil {
									global.GVA_LOG.Error("匹配预产异常", zap.Error(errR))
									_ = msg.Ack(true)
									continue
								}

								if len(keys) == 0 {
									fmt.Println("No keys found with the given prefix")
									global.GVA_LOG.Warn("当前组织无付款码可用, org", zap.Any("orgID", orgID))
									_ = msg.Ack(true)
									continue
								} else {
									var flag bool
									var keyEle string

									for _, key := range keys {
										var resTmpList []string
										resTmpList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
											Min:    "0",
											Max:    "0",
											Offset: 0,
											Count:  1,
										}).Result()

										// 检查是否为 redis.Nil 错误
										if err != nil {
											// 检查是否为 redis.Nil 错误
											if errors.Is(err, redis.Nil) {
												global.GVA_LOG.Warn("The set is empty", zap.Any("key", key))
											} else {
												global.GVA_LOG.Error("redis err", zap.Error(err))
											}
										} else if len(resTmpList) == 0 {
											global.GVA_LOG.Warn("The set is empty", zap.Any("key", key))
										} else {
											// 匹配成功
											flag = true
											keyEle = key
											break
										}

									}

									if !flag {

										global.GVA_LOG.Error("当前归属地区预产不足, redis err", zap.Error(err), zap.Any("db.region", v.Obj.PayRegion))
										// 如果解析消息失败，则直接丢弃消息
										_ = msg.Reject(false)

										//入库操作记录
										record := sysModel.SysOperationRecord{
											Ip:      v.Ctx.ClientIP,
											Method:  v.Ctx.Method,
											Path:    v.Ctx.UrlPath,
											Agent:   v.Ctx.UserAgent,
											Status:  500,
											Latency: time.Since(now),
											Resp:    fmt.Sprintf(global.ResourceNotEnough),
											UserID:  v.Ctx.UserID,
										}
										err = operationRecordService.CreateSysOperationRecord(record)
										if err != nil {
											global.GVA_LOG.Error("record 入库失败..." + err.Error())
										}

										if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0); errDB != nil {
											global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(errDB.Error))
										}
										continue
									} else {
										// 匹配成功
										resList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), keyEle, &redis.ZRangeBy{
											Min:    "0",
											Max:    "0",
											Offset: 0,
											Count:  1,
										}).Result()

										keyMem := resList[0]

										// 把预产设置为已用
										global.GVA_REDIS.ZAdd(context.Background(), keyEle, redis.Z{
											Score:  1,
											Member: keyMem,
										})

										split := strings.Split(keyMem, "_")
										ID = split[0]
										MID = split[1]
										acAccount = split[2]
										imgContent = split[3]

										global.GVA_LOG.Warn("本地区匹配不足了，随机匹配其它预产", zap.Any("db.region", v.Obj.PayRegion), zap.Any("匹配码keyEle", keyEle), zap.Any("匹配码keyMem", keyMem))

										var pcDB vbox.ChannelPayCode
										err = global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id = ?", ID).First(&pcDB).Error
										if err != nil {
											//return nil, fmt.Errorf("匹配通道账号不存在！ 请核查：%v", err.Error())

											global.GVA_LOG.Error("匹配预产异常", zap.Error(err))

											_ = msg.Ack(true)
											return
										}
										v.Obj.CreatedBy = pcDB.CreatedBy
										v.Obj.AcId = pcDB.AcId
										v.Obj.PlatId = pcDB.PlatId
										v.Ctx.UserID = int(pcDB.CreatedBy)

										err = global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id = ?", ID).Update("code_status", 1).Error

										//匹配成功之后，要把当前acAccount相同金额的所有预产置为 等候状态，不允许被取用
										var waitIDs []string

										for _, key := range keys {
											var resWaitTmpList []string
											resWaitTmpList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
												Min:    "0",
												Max:    "0",
												Offset: 0,
												Count:  -1,
											}).Result()

											for _, waitMem := range resWaitTmpList {
												if strings.Contains(waitMem, acAccount) {
													global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
														Score:  4,
														Member: waitMem,
													})

													//添加到 waitIDs 中
													id := strings.Split(waitMem, "_")[0]
													waitIDs = append(waitIDs, id)
												}

											}
										}

										// 把 pay code中除了本ID的，其它都让他进入冷却状态(包括对应通道账号)
										global.GVA_DB.Debug().Model(&vbox.ChannelPayCode{}).Where("id in ? ", waitIDs).Update("code_status", 4)
										global.GVA_DB.Debug().Model(&vbox.ChannelAccount{}).Where("ac_id in (?) ", pcDB.AcId).
											Update("status", 4).Update("sys_status", 4)

										// 把当前acAccount下所有的预产等待队列置为冷却状态
										waitAccPcKey := fmt.Sprintf(global.AccWaiting, pcDB.AcId)
										waitIDsTmp := strings.Join(waitIDs, ",")
										global.GVA_REDIS.Set(context.Background(), waitAccPcKey, waitIDsTmp, cdTime)

										waitMsg := strings.Join([]string{waitAccPcKey, waitIDsTmp}, "_")
										err = ch.PublishWithDelay(PayCodeCDCheckDelayedExchange, PayCodeCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)

									}
								}
							}

						}
					} else {
						global.GVA_LOG.Error("当前渠道非法无可用, channel", zap.Any("channel", cid))
						// 如果解析消息失败，则直接丢弃消息
						_ = msg.Reject(false)
						continue
					}

				}

				// 2.1 判断当前产品属于那种类型 1-商铺关联，2-付码关联
				eventType, _ := HandleEventType(cid)

				var eventID string
				var rsUrl string

				if eventType == 1 {
					if v.Obj.EventId == "" {
						// 获取 ID后，获取具体的跳转Url
						rsUrl, err = HandleResourceUrl2chShop(eventID)
						if err != nil {
							global.GVA_LOG.Error("商铺地址未匹配", zap.Error(err))
						}
						// 轮询获取同通道 指定金额的 商铺地址+ID
						eventID, err = HandleEventID2chShop(cid, v.Obj.Money, orgTmp)
						if err != nil {
							global.GVA_LOG.Error("商铺未匹配", zap.Error(err))
						}
					} else {
						eventID = v.Obj.EventId
						rsUrl = v.Obj.ResourceUrl
					}
					global.GVA_LOG.Info("商铺匹配", zap.Any("eventID", eventID), zap.Any("rsUrl", rsUrl))

				} else if eventType == 2 {
					eventID = MID
					rsUrl = imgContent
				}

				v.Obj.EventType = eventType
				v.Obj.EventId = eventID
				v.Obj.ResourceUrl = rsUrl

				if err := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Updates(v.Obj); err != nil {
					global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(err.Error))
				}

				od := request.PayOrderAndCtx{
					Obj: v.Obj,
					Ctx: v.Ctx,
				}

				var marshal []byte
				marshal, err = json.Marshal(od)

				odKey := fmt.Sprintf(global.PayOrderKey, v.Obj.OrderId)
				jsonString, _ := json.Marshal(v.Obj)
				global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)

				//3. 匹配账号后，更新订单信息（账号信息，订单支付链接处理）
				errQ = ch.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 30*time.Second)

				if errQ != nil {
					global.GVA_LOG.Error("订单匹配异常，消息丢弃", zap.Any("对应单号", v.Obj.OrderId), zap.Error(errQ))

					_ = msg.Reject(false)
					continue
				}

				_ = msg.Ack(true)
				global.GVA_LOG.Info("订单匹配完成，进入查询池等候付款", zap.Any("对应单号", v.Obj.OrderId))

			}
			wg.Done()
		}(i + 1)
	}
	global.GVA_LOG.Info("Vbox OrderWaitingTask 初始化搞定")
	// 等待所有消费者完成处理
	wg.Wait()
}

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

				if global.TxContains(chanID) {

					global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", *odDB.CreatedAt), zap.Any("传入的过期时间", *expTime))
					records, errQ := product.QryQQRecordsBetween(vca, *odDB.CreatedAt, *expTime)
					if errQ != nil {
						// 查单有问题，直接订单要置为超时，消息置为处理完毕
					}
					rdMap := product.Classifier(records.WaterList)
					if vm, ok := rdMap["Q币"]; !ok {
						global.GVA_LOG.Info("还没有QB的充值记录")
					} else {
						if rd, ok2 := vm[strconv.FormatInt(int64(v.Obj.Money*100), 10)]; !ok2 {
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
				} else if global.PcContains(chanID) {
					global.GVA_LOG.Info("传入的时间", zap.Any("传入的创建时间", *odDB.CreatedAt), zap.Any("传入的过期时间", *expTime))

					records, errQ := product.QryQQRecordsByID(vca, odDB.PlatId)
					if errQ != nil {
						// 查单有问题，直接订单要置为超时，消息置为处理完毕
					}
					rdMap := product.Classifier(records.WaterList)
					if vm, ok := rdMap["Q币"]; !ok {
						global.GVA_LOG.Info("还没有QB的充值记录")
					} else {
						if rd, ok2 := vm[strconv.FormatInt(int64(v.Obj.Money*100), 10)]; !ok2 {
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
										EventId:   v.Obj.EventId,
										Recharge:  -v.Obj.Money,
										Remark:    fmt.Sprintf(global.WalletEventOrderCost, v.Obj.Money, v.Obj.OrderId),
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

				// 重新丢回去 下一个30s再查一次
				marshal, err := json.Marshal(v)
				err = ch.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 30*time.Second)
			}
			wg.Done()
		}(i + 1)
	}

	// 等待所有消费者完成处理
	wg.Wait()
	//time.Sleep(time.Minute)
	global.GVA_LOG.Info("MqOrderConfirmTask 初始化搞定")

}

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
				global.GVA_LOG.Info("收到一条需要进行发起回调的订单消息", zap.Any("order ID", v.Obj.OrderId))

				//1. 筛选匹配是哪个产品

				//1.0 核查商户
				var vpa vbox.PayAccount
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

				notifyUrl := v.Obj.NotifyUrl
				client := vbHttp.NewHTTPClient()
				var headers = map[string]string{
					"Content-Type":  "application/json",
					"Authorization": "Bearer token",
				}
				var payUrl string
				payUrl, err = HandlePayUrl2PAcc(v.Obj.OrderId)

				signBody := &vboxRep.Order2PayAccountRes{
					OrderId:   v.Obj.OrderId,
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

				global.GVA_LOG.Info("请求body", zap.Any("notifyBody", notifyBody))

				var options = &vbHttp.RequestOptions{
					Headers:      headers,
					MaxRedirects: 3,
					PayloadType:  "json",
					Payload:      notifyBody,
				}

				response, errH := client.Post(notifyUrl, options)
				global.GVA_LOG.Info("回调响应消息", zap.Any("resp", response))
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
						Resp:    fmt.Sprintf(global.NotifyHandSuccess, response),
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

				if err != nil {
					global.GVA_LOG.Error("订单匹配异常，消息丢弃", zap.Any("对应单号", v.Obj.OrderId), zap.Error(err))

					_ = msg.Reject(false)
					continue
				}

				// 4.入库wallet
				var c int64
				global.GVA_DB.Model(&vbox.UserWallet{}).Where("event_id = ?", v.Obj.EventId).Count(&c)

				if c == 0 {
					wallet := vbox.UserWallet{
						Uid:       v.Obj.CreatedBy,
						CreatedBy: v.Obj.CreatedBy,
						Type:      global.WalletOrderType,
						EventId:   v.Obj.EventId,
						Recharge:  -v.Obj.Money,
						Remark:    fmt.Sprintf(global.WalletEventOrderCost, v.Obj.Money, v.Obj.OrderId),
					}

					global.GVA_DB.Model(&vbox.UserWallet{}).Save(&wallet)
				}

				_ = msg.Ack(true)
				global.GVA_LOG.Info("订单完成，回调完成", zap.Any("对应单号", v.Obj.OrderId))

			}
			wg.Done()
		}(i + 1)
	}
	global.GVA_LOG.Info("Vbox OrderCallbackTask 初始化搞定")
	// 等待所有消费者完成处理
	wg.Wait()
}

func HandleEventID2chShop(chanID string, money int, orgIDs []uint) (orgShopID string, err error) {
	// 1-商铺关联
	var vsList []vbox.ChannelShop

	var zs []redis.Z
	var key string
	for _, orgID := range orgIDs {
		key = fmt.Sprintf(global.ChanOrgShopAddrZSet, orgID, chanID, money)
		zs, err = global.GVA_REDIS.ZRangeArgsWithScores(context.Background(), redis.ZRangeArgs{
			Key:   key,
			Start: 0,
			Stop:  -1,
		}).Result()
		if err != nil {
			return "", err
		}
		if len(zs) <= 0 { // redis 没查到，查一下库
			userIDs := utils2.GetUsersByOrgId(orgID)
			err = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("cid = ? and money = ? and status = 1", chanID, money).
				Where("created_by in ?", userIDs).Find(&vsList).Error
			if err != nil {
				return "", err
			}
			if len(vsList) <= 0 {
				continue
			}

			//如果查到库里有， 设置进 redis 中
			for _, shop := range vsList {
				global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
					Score:  float64(time.Now().Unix()),
					Member: shop.ProductId + "_" + strconv.FormatUint(uint64(shop.ID), 10),
				})
			}
		} else {
			break
		}
	}

	if len(zs) <= 0 {
		return "", fmt.Errorf("该组织配置的资源不足，请核查")
	}

	z := zs[len(zs)-1] //取出最后一个，重新设置utc时间戳
	orgShopID = z.Member.(string)
	global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()), // 重新放进去，score设置最新的时间
		Member: orgShopID,
	})

	global.GVA_LOG.Info("获取引导商铺匹配信息", zap.Any("orgShopID", orgShopID))

	return orgShopID, err
}

func HandlePayUrl2PAcc(orderId string) (string, error) {
	conn := global.GVA_REDIS.Conn()
	defer conn.Close()
	key := global.PAccPay
	var url string
	//paccCreateUrl, err := global.GVA_REDIS.Ping(context.Background()).Result()
	//paccCreateUrl, err := conn.Ping(context.Background()).Result()
	//fmt.Printf(paccCreateUrl)
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		global.GVA_LOG.Warn("当前key不存在", zap.Any("key", key))

		var proxy vbox.Proxy
		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
		err = db.Where("status = ?", 1).Where("chan = ?", key).
			First(&proxy).Error
		if err != nil || proxy.Url == "" {
			return "", err
		}
		url = proxy.Url + orderId

		//global.GVA_REDIS.Set(context.Background(), key, proxy.Url, 0)
		conn.Set(context.Background(), key, proxy.Url, 0)
		global.GVA_LOG.Info("查库获取", zap.Any("商户订单地址", url))

		return url, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
	} else {
		var preUrl string
		//preUrl, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		preUrl, err = conn.Get(context.Background(), key).Result()
		url = preUrl + orderId
		global.GVA_LOG.Info("缓存池取出", zap.Any("商户订单地址", url))
	}
	return url, err
}

func HandleResourceUrl2chShop(eventID string) (addr string, err error) {
	//1. 如果是引导类的，获取引导地址 - channel shop
	split := strings.Split(eventID, "_")
	if len(split) <= 1 {
		return "", fmt.Errorf("解析商铺prod异常，param: %s", eventID)
	}
	//格式 （prodID_ID）
	ID := split[1]

	var shop vbox.ChannelShop
	db := global.GVA_DB.Model(&vbox.ChannelShop{}).Table("vbox_channel_shop")
	err = db.Where("id = ?", ID).First(&shop).Error
	if err != nil {
		return "", err
	}

	return shop.Address, nil
}

func HandleResourceUrl2payCode(eventID string) (addr string, err error) {
	//1. 付码类 - pay code
	// 格式（mid）

	var pc vbox.ChannelPayCode
	err = global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("mid = ?", eventID).First(&pc).Error
	if err != nil {
		return "", err
	}

	return pc.ImgContent, nil
}

/*
func HandleEventID2payCode(chanID string, money int, orgIDs []uint) (orgPayCodeID string, err error) {
	// 2-付码关联
	var pcList []vbox.ChannelPayCode

	var zs []redis.Z
	var key string
	for _, orgID := range orgIDs {
		key = fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgID, chanID, money)
		zs, err = global.GVA_REDIS.ZRangeArgsWithScores(context.Background(), redis.ZRangeArgs{
			Key:   key,
			Start: 0,
			Stop:  -1,
		}).Result()
		if err != nil {
			return "", err
		}
		if len(zs) <= 0 { // redis 没查到，查一下库
			userIDs := utils2.GetUsersByOrgId(orgID)
			// 当前time
			now := time.Now()

			err = global.GVA_DB.Debug().Model(&vbox.ChannelPayCode{}).Where("cid = ? and money = ? and code_status = 2", chanID, money).
				Where("created_by in ?", userIDs).Where("exp_time > ?", now).Find(&pcList).Error
			if err != nil {
				return "", err
			}
			if len(pcList) <= 0 {
				continue
			}

			//如果查到库里有， 设置进 redis 中
			for _, pc := range pcList {
				global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
					Score:  0,
					Member: fmt.Sprintf("%d", pc.ID) + "_" + pc.Mid + "_" + pc.AcAccount + "_" + pc.ImgContent,
				})
			}
		}
		break
	}

	if len(zs) <= 0 {
		return "", fmt.Errorf("该组织配置的资源不足，请核查")
	}

	return "", nil
}*/

func HandleEventType(chanID string) (int, error) {
	// 1-商铺关联，2-付码关联

	chanCode, _ := strconv.Atoi(chanID)
	if chanCode >= 1000 && chanCode <= 1099 {
		return 1, nil
	} else if chanCode >= 2000 && chanCode <= 2099 {
		return 1, nil
	} else if chanCode >= 3000 && chanCode <= 3099 {
		return 2, nil
	}
	return 0, fmt.Errorf("不存在的event类型")
}

func HandleExpTime2Product(chanID string) (time.Duration, error) {
	var key string

	if global.TxContains(chanID) {
		key = "1000"
	} else if global.J3Contains(chanID) {
		key = "2000"
	} else if global.PcContains(chanID) {
		key = "3000"
	}

	var expTimeStr string
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		global.GVA_LOG.Warn("当前key不存在", zap.Any("key", key))

		var proxy vbox.Proxy
		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
		err = db.Where("status = ?", 1).Where("chan = ?", key).
			First(&proxy).Error
		if err != nil || proxy.Url == "" {
			return 0, err
		}
		expTimeStr = proxy.Url
		seconds, _ := strconv.Atoi(expTimeStr)
		duration := time.Duration(seconds) * time.Second

		global.GVA_REDIS.Set(context.Background(), key, int64(duration.Seconds()), 0)
		global.GVA_LOG.Info("数据库取出该产品的有效时长", zap.Any("channel code", chanID), zap.Any("过期时间(s)", seconds))

		return duration, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
		return 0, err
	} else {
		expTimeStr, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		seconds, _ := strconv.Atoi(expTimeStr)

		duration := time.Duration(seconds) * time.Second

		//global.GVA_LOG.Info("缓存池取出该产品的有效时长", zap.Any("channel code", chanID), zap.Any("过期时间(s)", seconds))
		return duration, err
	}
}
