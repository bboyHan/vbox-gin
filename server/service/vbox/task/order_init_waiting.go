package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/vbUtil"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	OrderWaitExchange = "vbox.order.waiting_exchange"
	OrderWaitQueue    = "vbox.order.waiting_queue"
	OrderWaitKey      = "vbox.order.waiting"
)

// OrderWaitingTask 订单入库匹配
func OrderWaitingTask() {

	var operationRecordService system.OperationRecordService
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

			// 说明：执行账号匹配
			deliveries, errQ := chX.Consume(OrderWaitQueue, "", false, false, false, false, nil)
			if errQ != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", OrderWaitQueue))
			}

			for msg := range deliveries {

				v := request.PayOrderAndCtx{}
				err = json.Unmarshal(msg.Body, &v)
				if err != nil {
					global.GVA_LOG.Error("订单匹配消息处理失败，MqOrderWaitingTask...", zap.Any("err", err.Error()))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}
				orderId := v.Obj.OrderId
				pacc := v.Obj.PAccount
				global.GVA_LOG.Info("收到一条需要进行匹配的订单消息", zap.Any("order ID", orderId), zap.Any("商户", pacc))

				//1.0 核查商户
				_, errP := vbUtil.ValidPacc(pacc, &v.Obj.CreatedBy)
				if errP != nil {
					global.GVA_LOG.Error("商户校验失败", zap.Any("商户", pacc), zap.Any("err", errP))
					// 如果商户校验失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}

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
					Resp:    fmt.Sprintf(global.OrderWaitingMsg),
					UserID:  v.Ctx.UserID,
				}
				err = operationRecordService.CreateSysOperationRecord(record)

				//2. 查供应库存账号是否充足 (优先从缓存池取，取空后查库取，如果库也空了，咋报错库存不足)
				//2.0 查一下单，如果发起初始化的时候已经匹配过账号了，就不用再匹配一次了
				var ID, accID, acAccount string
				var imgContent, MID string
				//获取当前组织
				orgTmp := utils2.GetSelfOrg(v.Obj.CreatedBy)
				orgID := orgTmp[0]
				cid := v.Obj.ChannelCode
				money := v.Obj.Money

				// 设置一个CD时间（订单创建时间到过期时间之差）
				t, errD := utils.GetCDTimeByCode(cid)
				if errD != nil {
					global.GVA_LOG.Error("获取CD时间异常", zap.Any("cid", cid), zap.Any("err", errD))
					if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; errDB != nil {
						global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(errDB))
					}
					_ = msg.Reject(false)
					continue
				}
				productInfo, errP := utils.GetProductByCode(cid)
				if errP != nil || productInfo == nil {
					global.GVA_LOG.Error("获取产品信息异常", zap.Any("cid", cid), zap.Any("err", errP))
					if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; errDB != nil {
						global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(errDB))
					}
					_ = msg.Reject(false)
					continue
				}
				// 设置缓冲时间 cd时间后的释放账号时间

				var accSetKey, accSetMem string
				switch {
				case strings.Contains(productInfo.Ext, "money"):
					accSetKey = fmt.Sprintf(productInfo.Ext, orgID, cid, money)
				default:
					accSetKey = fmt.Sprintf(productInfo.Ext, orgID, cid)
				}

				// 获取集合中可用的账号
				var resList []string
				resList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), accSetKey, &redis.ZRangeBy{
					Min:    "0",
					Max:    "0",
					Offset: 0,
					Count:  -1,
				}).Result()
				if err != nil {
					global.GVA_LOG.Error("获取集合中可用的账号匹配异常, redis err", zap.Error(err))
					if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; errDB != nil {
						global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(errDB))
					}
					_ = msg.Ack(true)
					continue
				}

				if resList != nil && len(resList) > 0 {
					randomList := utils.NewRandomList(resList)
					var useAccFlag bool
					for len(randomList.List) > 0 {
						accSetMem = randomList.GetRandom()
						global.GVA_LOG.Info("取出的accSetMem值", zap.Any("accSetMem", accSetMem), zap.Any("orderId", orderId))

						// 判断一下当前的账号tm的有没有同时拉单
						split := strings.Split(accSetMem, ",")
						var pcDB vbox.ChannelPayCode
						if global.PcContains(cid) {
							ID = split[0] // pc id
							MID = split[1]
							acAccount = split[2]
							imgContent = split[3]
							err = global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id = ?", ID).First(&pcDB).Error
							if err != nil {
								global.GVA_LOG.Error("pc 匹配预产异常", zap.Error(err))
								continue
							}
							accID = pcDB.AcId
						} else {
							ID = split[0]
							accID = split[1]
							acAccount = split[2]
						}

						var accJucKey string
						switch {
						case strings.Contains(productInfo.Ext, "money"):
							accJucKey = fmt.Sprintf(global.PayAccMoneyKey, accID, money)
						default:
							accJucKey = fmt.Sprintf(global.PayAccKey, accID)
						}

						jucTTL := global.GVA_REDIS.TTL(context.Background(), accJucKey).Val()
						if jucTTL > 0 { // tm的已经有这个号在这个时间段被拉单了
							global.GVA_LOG.Error("tm的已经有这个号在这个时间段被拉单了", zap.Any("acID", accID), zap.Any("account", acAccount))
							continue
						} else {
							// 查一下限额情况
							var vca vbox.ChannelAccount
							err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("ac_id = ?", accID).First(&vca).Error
							if err != nil {
								//return nil, fmt.Errorf("匹配通道账号不存在！ 请核查：%v", err.Error())
								global.GVA_LOG.Error("数据库查账号异常", zap.Error(err))
								continue
							}

							errM := vbUtil.CheckAccLimit(vca, money)
							if errM != nil {
								global.GVA_LOG.Error(errM.Error(), zap.Any("orderId", v.Obj.OrderId), zap.Any("acID", accID), zap.Any("acAccount", vca.AcAccount), zap.Any("money", money))
								continue
							}

							if global.J3Contains(cid) {
								// 查一下当前余额
								j3Record, errQy := product.QryJ3Record(vca)
								if errQy != nil {
									// 查单有问题，直接订单要置为超时，消息置为处理完毕
									global.GVA_LOG.Error("查询充值记录异常", zap.Error(errQy))
									if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; errDB != nil {
										global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(errDB))
									}
									record = sysModel.SysOperationRecord{
										Ip:      v.Ctx.ClientIP,
										Method:  v.Ctx.Method,
										Path:    v.Ctx.UrlPath,
										Agent:   v.Ctx.UserAgent,
										Status:  500,
										Latency: time.Since(now),
										Resp:    fmt.Sprintf(global.AccQryJ3RecordsEx, accID, acAccount),
										UserID:  v.Ctx.UserID,
									}
									err = operationRecordService.CreateSysOperationRecord(record)
									continue
								}

								balance := j3Record.LeftCoins
								J3AccBalanceKey := fmt.Sprintf(global.J3AccBalanceZSet, vca.AcId)
								nowTimeUnix := v.Obj.CreatedAt.Unix()
								keyMem := fmt.Sprintf("%s,%s,%v,%d,%d,%d,%d", orderId, vca.AcAccount, money, nowTimeUnix, balance, 0, 0)

								global.GVA_REDIS.ZAdd(context.Background(), J3AccBalanceKey, redis.Z{
									Score:  float64(nowTimeUnix),
									Member: keyMem,
								})

								v.Obj.PlatId = keyMem

								global.GVA_LOG.Info("当前余额情况", zap.Any("j3Record", j3Record))
							}

							if v.Obj.EventType == 2 {
								v.Obj.EventId = MID
								v.Obj.ResourceUrl = imgContent
							}

							// 把账号设置为已用
							global.GVA_REDIS.ZAdd(context.Background(), accSetKey, redis.Z{
								Score:  1,
								Member: accSetMem,
							})

							if global.PcContains(cid) {
								global.GVA_DB.Debug().Model(&vbox.ChannelPayCode{}).Where("id = ?", pcDB.ID).Update("code_status", 1)
							}

							global.GVA_REDIS.Set(context.Background(), accJucKey, v.Obj.OrderId, t.CDTime)

							var accWaitYdKey, accInfoVal string
							switch {
							case strings.Contains(productInfo.Ext, "money"):
								accWaitYdKey = fmt.Sprintf(global.ProdAccMoneyWaiting, productInfo.ProductId, accID, money)
								accInfoVal = fmt.Sprintf("%s,%s,%s,%v", ID, accID, acAccount, money)
							default:
								accWaitYdKey = fmt.Sprintf(global.ProdAccWaiting, productInfo.ProductId, accID)
								accInfoVal = fmt.Sprintf("%s,%s,%s", ID, accID, acAccount)
							}
							// 设置一个冷却时间
							global.GVA_REDIS.Set(context.Background(), accWaitYdKey, accInfoVal, t.CDTime)
							waitMsg := strings.Join([]string{accWaitYdKey, accInfoVal}, "-")
							err = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), t.CDTime)

							global.GVA_LOG.Info("正常取用", zap.Any("acID", accID), zap.Any("account", acAccount), zap.Any("orderId", orderId))

							// 赋值匹配账号到订单中
							v.Obj.CreatedBy = vca.CreatedBy
							v.Obj.AcId = vca.AcId
							v.Obj.AcAccount = vca.AcAccount
							if cid == "1101" || cid == "6001" || cid == "2001" {

							} else if cid == "3000" {
								v.Obj.PlatId = pcDB.PlatId
							} else {
								v.Obj.PlatId = utils.GenerateID(global.WalletEventOrderPrefix)
							}
							v.Ctx.UserID = int(vca.CreatedBy)
							useAccFlag = true

							record = sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
								Type:    global.OrderType,
								Status:  200,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.OrderWaitingFinishedMsg, v.Obj.AcAccount, v.Obj.AcId),
								UserID:  v.Ctx.UserID,
							}
							err = operationRecordService.CreateSysOperationRecord(record)

							break
						}
					}
					if !useAccFlag { //轮询一遍了，还是没取到的话，就失败
						if errDB := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Update("order_status", 0).Error; errDB != nil {
							global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(errDB))
						}
						record = sysModel.SysOperationRecord{
							Ip:      v.Ctx.ClientIP,
							Method:  v.Ctx.Method,
							Path:    v.Ctx.UrlPath,
							Agent:   v.Ctx.UserAgent,
							MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
							Type:    global.OrderType,
							Status:  500,
							Latency: time.Since(now),
							Resp:    fmt.Sprintf(global.ResourceAccNotEnough, cid, money),
							UserID:  v.Ctx.UserID,
						}
						err = operationRecordService.CreateSysOperationRecord(record)
						_ = msg.Ack(true)

						continue
					}
				} else {
					global.GVA_LOG.Error("当前渠道库存不足, channel", zap.Any("code", cid), zap.Any("pacc", pacc))
					// 如果解析消息失败，则直接丢弃消息
					_ = msg.Reject(false)
					continue
				}

				if err := global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", v.Obj.ID).Updates(v.Obj); err != nil {
					global.GVA_LOG.Info("MqOrderWaitingTask...", zap.Error(err.Error))
				}

				od := request.PayOrderAndCtx{
					Obj: v.Obj,
					Ctx: v.Ctx,
				}

				var marshal []byte
				marshal, err = json.Marshal(od)

				odKey := fmt.Sprintf(global.PayOrderKey, orderId)
				jsonString, _ := json.Marshal(v.Obj)
				global.GVA_REDIS.Set(context.Background(), odKey, jsonString, 300*time.Second)

				//3. 匹配账号后，更新订单信息（账号信息，订单支付链接处理）
				errQ = chX.PublishWithDelay(OrderConfirmDelayedExchange, OrderConfirmDelayedRoutingKey, marshal, 35*time.Second)

				if errQ != nil {
					global.GVA_LOG.Error("订单匹配异常，消息丢弃", zap.Any("对应单号", orderId), zap.Error(errQ))

					_ = msg.Reject(false)
					continue
				}

				_ = msg.Ack(true)
				global.GVA_LOG.Info("订单匹配完成，进入查询池等候付款", zap.Any("对应单号", orderId))

			}
			wg.Done()
		}(i + 1)
	}
	// 等待所有消费者完成处理
	wg.Wait()
	global.GVA_LOG.Info("Vbox OrderWaitingTask 初始化搞定")
}

/*func HandleEventID2chShop(chanID string, money int, orgIDs []uint) (orgShopID string, err error) {
	// 1-商铺关联
	var vsList []vbox.ChannelShop

	if chanID == "6001" {
		orgIDs = []uint{1}
	}
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
		global.GVA_REDIS.Set(context.Background(), key, proxy.Url, 0)
		global.GVA_LOG.Info("查库获取", zap.Any("商户订单地址", url))

		return url, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
	} else {
		var preUrl string
		//preUrl, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		preUrl, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		url = preUrl + orderId
		global.GVA_LOG.Info("缓存池取出", zap.Any("商户订单地址", url))
	}
	return url, err
}*/

/*func HandleResourceUrl2chShop(eventID string) (addr string, err error) {
	global.GVA_LOG.Info("接收event id", zap.Any("eventID", eventID))
	//1. 如果是引导类的，获取引导地址 - channel shop
	split := strings.Split(eventID, "_")
	if len(split) <= 1 {
		return "", fmt.Errorf("解析商铺prod异常，param: %s", eventID)
	}
	//格式 （prodID_ID）
	ID := split[1]

	var shop vbox.ChannelShop
	err = global.GVA_DB.Debug().Model(&vbox.ChannelShop{}).Where("id = ?", ID).First(&shop).Error
	if err != nil {
		return "", err
	}
	global.GVA_LOG.Info("查出shop", zap.Any("shop", shop))

	cid := shop.Cid

	var payUrl string
	switch cid {
	case "1001": //jd
		payUrl, err = utils.HandleJDUrl(shop.Address)
		if err != nil {
			return "", err
		}
	case "1002": //dy
		payUrl, err = utils.HandleDYUrl(shop.Address)
		if err != nil {
			return "", err
		}
	case "1003": //jym
		payUrl, err = utils.HandleAlipayUrl(shop.Address)
		if err != nil {
			return "", err
		}
	case "1004": //zfb
		payUrl, err = utils.HandleAlipayUrl(shop.Address)
		if err != nil {
			return "", err
		}
	case "1005": //qb tb
		payUrl, err = utils.HandleTBUrl(shop.Address)
		if err != nil {
			return "", err
		}
	case "1006": //qb xcx
		payUrl, err = utils.HandleXCXUrl(shop.Address)
		if err != nil {
			return "", err
		}
	case "1007": //qb xcx
		payUrl, err = utils.HandlePddUrl(shop.Address)
		if err != nil {
			return "", err
		}

	case "1101": //jun ka qb tb
		payUrl, err = utils.HandleTBUrl(shop.Address)
		if err != nil {
			return "", err
		}

	case "1201": //dnf tb
		payUrl, err = utils.HandleTBUrl(shop.Address)
		if err != nil {
			return "", err
		}
	case "1202": //dnf jd
		payUrl, err = utils.HandleJDUrl(shop.Address)
		if err != nil {
			return "", err
		}

	case "2001": //j3 tb
		payUrl, err = utils.HandleTBUrl(shop.Address)
		if err != nil {
			return "", err
		}
	case "4001": //sdo tb
		payUrl, err = utils.HandleTBUrl(shop.Address)
		if err != nil {
			return "", err
		}
	case "6001": //ec jd
		payUrl, err = utils.HandleJDUrl(shop.Address)
		if err != nil {
			return "", err
		}
	default:
		payUrl = shop.Address
	}

	global.GVA_LOG.Info("处理后", zap.Any("payUrl", payUrl))

	return payUrl, nil

}*/

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

//func HandleEventType(chanID string) (int, error) {
//	// 1-商铺关联，2-付码关联
//
//	chanCode, _ := strconv.Atoi(chanID)
//	if chanCode >= 1000 && chanCode <= 1099 {
//		return 1, nil
//	} else if chanCode >= 1100 && chanCode <= 1199 {
//		return 1, nil
//	} else if chanCode >= 1200 && chanCode <= 1299 {
//		return 1, nil
//	} else if chanCode >= 2000 && chanCode <= 2099 {
//		return 1, nil
//	} else if chanCode >= 3000 && chanCode <= 3099 {
//		return 2, nil
//	} else if chanCode >= 4000 && chanCode <= 4099 {
//		return 1, nil
//	} else if chanCode >= 5000 && chanCode <= 5099 {
//		return 1, nil
//	} else if chanCode >= 6000 && chanCode <= 6099 {
//		return 1, nil
//	}
//	return 0, fmt.Errorf("不存在的event类型")
//}

//func HandleExpTime2Product(chanID string) (time.Duration, error) {
//	key := fmt.Sprintf(global.ProdCodeKey, chanID)
//	var expTimeStr string
//	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
//	if count == 0 {
//		if err != nil {
//			global.GVA_LOG.Error("redis ex：", zap.Error(err))
//		}
//
//		var proxy vbox.Proxy
//		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
//		err = db.Where("status = ?", 1).Where("chan = ?", chanID).First(&proxy).Error
//		if err != nil || proxy.Url == "" {
//			return 0, err
//		}
//		expTimeStr = proxy.Url
//		seconds, _ := strconv.Atoi(expTimeStr)
//		duration := time.Duration(seconds) * time.Second
//
//		global.GVA_REDIS.Set(context.Background(), key, int64(duration.Seconds()), 0)
//		global.GVA_LOG.Info("数据库取出该产品的有效时长", zap.Any("channel code", chanID), zap.Any("过期时间(s)", seconds))
//
//		return duration, nil
//	} else if err != nil {
//		global.GVA_LOG.Error("redis ex：", zap.Error(err))
//		return 0, err
//	} else {
//		expTimeStr, err = global.GVA_REDIS.Get(context.Background(), key).Result()
//		seconds, _ := strconv.Atoi(expTimeStr)
//
//		duration := time.Duration(seconds) * time.Second
//
//		//global.GVA_LOG.Info("缓存池取出该产品的有效时长", zap.Any("channel code", chanID), zap.Any("过期时间(s)", seconds))
//		return duration, err
//	}
//}
