package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"
	"time"
)

// 账号开启查询
const (
	ChanAccEnableCheckExchange = "vbox.order.acc_enable_check_exchange"
	ChanAccEnableCheckQueue    = "vbox.order.acc_enable_check_queue"
	ChanAccEnableCheckKey      = "vbox.order.acc_enable_check"
)

// ChanAccEnableCheckTask 通道账号开启状态核查
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
	consumerCount := 10
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
			now := time.Now()

			for msg := range deliveries {
				//v := &map[string]interface{}{}
				//err := json.Unmarshal(msg.Body, v)
				//global.GVA_LOG.Info(fmt.Sprintf("%v", msg.Body))

				v := vboxReq.ChanAccAndCtx{}
				err := json.Unmarshal(msg.Body, &v)
				if err != nil {
					global.GVA_LOG.Error("错了，直接丢了..." + err.Error())
					_ = msg.Reject(false)
					continue
				}

				orgTmp := utils2.GetSelfOrg(v.Obj.CreatedBy)
				ID := v.Obj.ID
				cid := v.Obj.Cid
				acId := v.Obj.AcId
				acAccount := v.Obj.AcAccount
				if v.Obj.Status == 1 {
					global.GVA_LOG.Info("收到一条需要处理的账号【开启】", zap.Any("v", v))

					// 1. 查询该用户的余额是否充足
					var balance int
					err = global.GVA_DB.Model(&vbox.UserWallet{}).Select("IFNULL(sum(recharge), 0) as balance").
						Where("uid = ?", v.Obj.CreatedBy).Scan(&balance).Error
					if err != nil {
						global.GVA_LOG.Info("查询该用户的余额错了，直接丢了..." + err.Error())
						_ = msg.Reject(false)
						continue
					}

					if balance <= 0 { //余额不足，则 log 一条
						//入库操作记录
						record := sysModel.SysOperationRecord{
							Ip:      v.Ctx.ClientIP,
							Method:  v.Ctx.Method,
							Path:    v.Ctx.UrlPath,
							Agent:   v.Ctx.UserAgent,
							Status:  500,
							Latency: time.Since(now),
							Resp:    fmt.Sprintf(global.BalanceNotEnough, acId, v.Obj.AcAccount),
							UserID:  v.Ctx.UserID,
						}

						err = operationRecordService.CreateSysOperationRecord(record)
						if err != nil {
							global.GVA_LOG.Error("余额不足情况下，record 入库失败..." + err.Error())
						}
						err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
							Update("sys_status", 0).Error
						// 不允许开启sys_status， 到这里结束
						_ = msg.Reject(false)
						continue
					}

					// 2. 查询账号是否有超 金额限制，或者笔数限制
					// 2.1 日限制
					if v.Obj.DailyLimit > 0 {
						dailyKey := fmt.Sprintf(global.ChanAccDailyUsed, acId)
						dailyUsed, err := global.GVA_REDIS.Get(context.Background(), dailyKey).Int()
						if err == redis.Nil { // redis中无，查一下库
							var dailySum int
							// 获取今天的时间范围
							startOfDay := time.Now().UTC().Truncate(24 * time.Hour)
							endOfDay := startOfDay.Add(24 * time.Hour)

							err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("sum(money) as dailySum").
								Where("ac_id = ?", acId).
								Where("order_status = ? AND created_at BETWEEN ? AND ?", 1, startOfDay, endOfDay).Scan(&dailySum).Error

							if err != nil {
								global.GVA_LOG.Error("当前账号计算日消耗查mysql错误，直接丢了..." + err.Error())
								_ = msg.Reject(false)
								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
									Update("sys_status", 0).Error
								continue
							}

							// 查完算一下并且给redis赋值设置一下
							global.GVA_REDIS.Set(context.Background(), dailyKey, dailySum, 1*time.Minute)

						} else if err != nil {
							global.GVA_LOG.Error("当前账号计算日消耗差redis错误，直接丢了..." + err.Error())
							_ = msg.Reject(false)
							continue
						} else {                               // redis查出来了，直接比一下
							if dailyUsed >= v.Obj.DailyLimit { // 如果日消费已经超了，不允许开启了，直接结束

								//入库操作记录
								record := sysModel.SysOperationRecord{
									Ip:      v.Ctx.ClientIP,
									Method:  v.Ctx.Method,
									Path:    v.Ctx.UrlPath,
									Agent:   v.Ctx.UserAgent,
									Status:  500,
									Latency: time.Since(now),
									Resp:    fmt.Sprintf(global.AccDailyLimitNotEnough, acId, v.Obj.AcAccount),
									UserID:  v.Ctx.UserID,
								}

								err = operationRecordService.CreateSysOperationRecord(record)
								if err != nil {
									global.GVA_LOG.Error("当前账号日消耗已经超限额情况下，record 入库失败..." + err.Error())
								}

								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
									Update("sys_status", 0).Error

								global.GVA_LOG.Info("当前账号日消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
								_ = msg.Reject(false)
								continue
							}
						}
					}
					// 2.2 总限制
					if v.Obj.TotalLimit > 0 {
						totalKey := fmt.Sprintf(global.ChanAccTotalUsed, acId)
						totalUsed, err := global.GVA_REDIS.Get(context.Background(), totalKey).Int()
						if err == redis.Nil { // redis中无，查一下库
							var totalSum int

							err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as totalSum").
								Where("ac_id = ?", acId).
								Where("order_status = ?", 1).Scan(&totalSum).Error

							if err != nil {
								global.GVA_LOG.Error("当前账号计算总消耗查mysql错误，直接丢了..." + err.Error())
								_ = msg.Reject(false)
								continue
							}

							// 查完算一下并且给redis赋值设置一下
							global.GVA_REDIS.Set(context.Background(), totalKey, totalSum, 1*time.Minute)

						} else if err != nil {
							global.GVA_LOG.Error("当前账号计算总消耗差redis错误，直接丢了..." + err.Error())
							_ = msg.Reject(false)

							continue
						} else {                               // redis查出来了，直接比一下
							if totalUsed >= v.Obj.TotalLimit { // 如果总消费已经超了，不允许开启了，直接结束

								//入库操作记录
								record := sysModel.SysOperationRecord{
									Ip:      v.Ctx.ClientIP,
									Method:  v.Ctx.Method,
									Path:    v.Ctx.UrlPath,
									Agent:   v.Ctx.UserAgent,
									Status:  500,
									Latency: time.Since(now),
									Resp:    fmt.Sprintf(global.AccTotalLimitNotEnough, acId, v.Obj.AcAccount),
									UserID:  v.Ctx.UserID,
								}

								err = operationRecordService.CreateSysOperationRecord(record)
								if err != nil {
									global.GVA_LOG.Error("当前账号总消耗已经超限额情况下，record 入库失败..." + err.Error())
								}

								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
									Update("sys_status", 0).Error
								global.GVA_LOG.Info("当前账号总消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
								_ = msg.Reject(false)
								continue
							}
						}
					}
					// 2.3 笔数限制
					if v.Obj.CountLimit > 0 {

						var count int64

						err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("ac_id = ? and order_status = ?", acId, 1).Count(&count).Error

						if err != nil {
							global.GVA_LOG.Error("当前账号笔数消耗查mysql错误，直接丢了..." + err.Error())
							_ = msg.Reject(false)
							continue
						}

						if int(count) >= v.Obj.CountLimit { // 如果笔数消费已经超了，不允许开启了，直接结束

							//入库操作记录
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.AccCountLimitNotEnough, acId, v.Obj.AcAccount),
								UserID:  v.Ctx.UserID,
							}

							err = operationRecordService.CreateSysOperationRecord(record)
							if err != nil {
								global.GVA_LOG.Error("当前账号笔数消耗已经超限额情况下，record 入库失败..." + err.Error())
							}
							err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0).Error
							global.GVA_LOG.Warn("当前账号笔数消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
							_ = msg.Reject(false)
							continue
						}
					}

					// 3. 筛选匹配是哪个产品，查一下对应产品的账户是否能够正常官方使用
					if global.TxContains(cid) { //腾讯

						errQ := product.QryQQRecords(v.Obj)
						if errQ != nil {
							global.GVA_LOG.Error("当前账号查官方记录异常情况下，record 入库失败...", zap.Error(errQ))

							//入库操作记录
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.AccQryRecordsEx, acId, v.Obj.AcAccount),
								UserID:  v.Ctx.UserID,
							}

							err = operationRecordService.CreateSysOperationRecord(record)
							if err != nil {
								global.GVA_LOG.Error("当前账号查官方记录异常情况下，record 入库失败..." + err.Error())
							}

							err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0).Error
							global.GVA_LOG.Warn("当前账号查官方记录异常了，结束...", zap.Any("ac info", v.Obj))
							_ = msg.Reject(false)
							continue

						}

					} else if global.J3Contains(cid) { //剑三
						_, errQ := product.QryJ3Record(v.Obj)
						if errQ != nil {
							global.GVA_LOG.Error("当前账号查官方记录异常情况下，record 入库失败..." + errQ.Error())
							//入库操作记录
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.AccQryJ3RecordsEx, acId, v.Obj.AcAccount),
								UserID:  v.Ctx.UserID,
							}

							errR := operationRecordService.CreateSysOperationRecord(record)
							if errR != nil {
								global.GVA_LOG.Error("当前账号查官方记录异常情况下，record 入库失败..." + errR.Error())
							}

							global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0)
							global.GVA_LOG.Warn("当前账号查官方记录异常了，结束...", zap.Any("ac info", v.Obj))
							_ = msg.Reject(false)
							continue
						}
					} else if global.PcContains(cid) { //QB直付
						errQ := product.QryQQRecords(v.Obj)
						if errQ != nil {
							global.GVA_LOG.Error("当前账号查官方记录异常情况下，record 入库失败..." + errQ.Error())
							//入库操作记录
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.AccQryRecordsEx, acId, v.Obj.AcAccount),
								UserID:  v.Ctx.UserID,
							}

							errR := operationRecordService.CreateSysOperationRecord(record)
							if errR != nil {
								global.GVA_LOG.Error("当前账号查官方记录异常情况下，record 入库失败..." + errR.Error())
							}

							global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0)
							global.GVA_LOG.Warn("当前账号查官方记录异常了，结束...", zap.Any("ac info", v.Obj))
							_ = msg.Reject(false)
							continue

						}
					}

					//2.
					if global.TxContains(cid) { //QB引导，

						moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
						moneyList := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
						pattern := fmt.Sprintf(global.ChanOrgQBAccZSetPrefix, orgTmp[0], cid)
						keys := global.GVA_REDIS.Keys(context.Background(), pattern).Val()

						if len(moneyList) >= len(keys) {
							for _, money := range moneyList {
								moneyTmp := money
								go func() {
									waitAccYdKey := fmt.Sprintf(global.YdQBAccWaiting, acId, moneyTmp)
									waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
									waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
									ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
									if ttl > 0 { //该账号正在冷却中
										global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
										cdTime := ttl
										_ = ch.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
										global.GVA_LOG.Info("开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
									} else {
										global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
										accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgTmp[0], cid, moneyTmp)
										global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
										global.GVA_LOG.Info("开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									}
								}()
							}
						} else {
							for _, key := range keys {
								keyTmp := key
								go func() {
									money := strings.Split(keyTmp, ":")[3]
									waitAccYdKey := fmt.Sprintf(global.YdQBAccWaiting, acId, money)
									waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, money)
									waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
									ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
									if ttl > 0 { //该账号正在冷却中
										global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
										cdTime := ttl
										_ = ch.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
										global.GVA_LOG.Info("开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
									} else {
										global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
										accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgTmp[0], cid, money)
										global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
										global.GVA_LOG.Info("开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									}
								}()
							}
						}

					} else if global.SdoContains(cid) { //sdo

						moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
						moneyList := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
						pattern := fmt.Sprintf(global.ChanOrgSdoAccZSetPrefix, orgTmp[0], cid)
						keys := global.GVA_REDIS.Keys(context.Background(), pattern).Val()

						if len(moneyList) >= len(keys) {
							for _, money := range moneyList {
								moneyTmp := money
								go func() {
									waitAccYdKey := fmt.Sprintf(global.YdSdoAccWaiting, acId, moneyTmp)
									waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
									waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
									ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
									if ttl > 0 { //该账号正在冷却中
										global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
										cdTime := ttl
										_ = ch.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
										global.GVA_LOG.Info("开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
									} else {
										global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
										accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgTmp[0], cid, moneyTmp)
										global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
										global.GVA_LOG.Info("开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									}
								}()
							}
						} else {
							for _, key := range keys {
								keyTmp := key
								go func() {
									money := strings.Split(keyTmp, ":")[3]
									waitAccYdKey := fmt.Sprintf(global.YdSdoAccWaiting, acId, money)
									waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, money)
									waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
									ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
									if ttl > 0 { //该账号正在冷却中
										global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
										cdTime := ttl
										_ = ch.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
										global.GVA_LOG.Info("开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
									} else {
										global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
										accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgTmp[0], cid, money)
										global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
										global.GVA_LOG.Info("开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									}
								}()
							}
						}

					} else if global.J3Contains(cid) { //J3

						moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
						moneyList := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
						if len(moneyList) < 1 {
							global.GVA_LOG.Error("商铺没有匹配资源...")
							//入库操作记录
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.AccQryJ3ShopEx, acId, v.Obj.AcAccount),
								UserID:  v.Ctx.UserID,
							}

							errR := operationRecordService.CreateSysOperationRecord(record)
							if errR != nil {
								global.GVA_LOG.Error("商铺没有匹配资源，record 入库失败..." + errR.Error())
							}

							global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0)
							global.GVA_LOG.Warn("商铺没有匹配资源，结束...", zap.Any("ac info", v.Obj))
							_ = msg.Reject(false)
							continue
						}

						accKey := fmt.Sprintf(global.ChanOrgJ3AccZSet, orgTmp[0], cid)
						waitAccYdKey := fmt.Sprintf(global.YdJ3AccWaiting, acId)

						waitAccMem := fmt.Sprintf("%v_%s_%s", ID, acId, acAccount)
						waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
						ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
						if ttl > 0 { //该账号正在冷却中
							global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
							cdTime := ttl
							_ = ch.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
							global.GVA_LOG.Info("开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
						} else {
							global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)

							global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
							global.GVA_LOG.Info("开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
						}

					} else if global.PcContains(cid) { //QB直付，查一下有没有还没被禁用的预产，把还没过期的恢复
						var pcDBList []vbox.ChannelPayCode
						global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("ac_id = ? and code_status = 5", acId).Find(&pcDBList)
						if len(pcDBList) == 0 {
							global.GVA_LOG.Info("开启过程校验..暂无需要处理的预产", zap.Any("当前账号", acId))
						} else {

							for _, pcDB := range pcDBList {

								pcKey := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0], pcDB.Cid, pcDB.Money, pcDB.Operator, pcDB.Location)
								pcMem := fmt.Sprintf("%d", pcDB.ID) + "_" + pcDB.Mid + "_" + pcDB.AcAccount + "_" + pcDB.ImgContent

								global.GVA_LOG.Info("开启过程校验..处理预产开启匹配", zap.Any("pcKey", pcKey), zap.Any("pcMem", pcMem))

								if pcDB.ExpTime.After(time.Now()) { // 设置的过期时间比当前时间晚，如果账号冷却中，置预产为4，发起cd任务
									waitAccPcKey := fmt.Sprintf(global.PcAccWaiting, pcDB.AcId)
									// 查一下redis key的过期时间
									ttl := global.GVA_REDIS.TTL(context.Background(), waitAccPcKey).Val()
									if ttl > 0 {
										global.GVA_REDIS.ZAdd(context.Background(), pcKey, redis.Z{Score: 4, Member: pcMem})
										global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id =?", pcDB.ID).Update("code_status", 4)

										waitIDsTmp := fmt.Sprintf("%d", pcDB.ID)
										waitMsg := strings.Join([]string{waitAccPcKey, waitIDsTmp}, "-")
										cdTime := ttl
										_ = ch.PublishWithDelay(PayCodeCDCheckDelayedExchange, PayCodeCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
										global.GVA_LOG.Info("开启过程校验..处理预产开启匹配,账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
									} else {
										global.GVA_REDIS.ZAdd(context.Background(), pcKey, redis.Z{Score: 0, Member: pcMem})
										global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id =?", pcDB.ID).Update("code_status", 2)
										global.GVA_LOG.Info("开启过程校验..处理预产开启匹配,账号正常状态,更新库状态为2")
									}

								} else { // 过期了，直接删除redis，并且状态置为失效 3
									global.GVA_REDIS.ZRem(context.Background(), pcKey, pcMem)
									global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id =?", pcDB.ID).Update("code_status", 3)
									global.GVA_LOG.Info("开启过程校验..处理预产过期了,更新库状态为3")
								}

							}
						}
					}
					//3. 校验都没啥问题，开启sys_status = 1，即可以调度订单使用
					//global.GVA_LOG.Info("消息 : ", zap.Any("msg", msg.Body))
					err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).Update("sys_status", 1).Error

				} else {
					global.GVA_LOG.Info("收到一条需要处理的账号【关闭】", zap.Any("v", v))

					if global.TxContains(cid) {
						moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
						moneyList := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
						pattern := fmt.Sprintf(global.ChanOrgQBAccZSetPrefix, orgTmp[0], cid)
						keys := global.GVA_REDIS.Keys(context.Background(), pattern).Val()

						if len(moneyList) >= len(keys) {
							for _, money := range moneyList {
								moneyTmp := money
								go func() {
									waitAccYdKey := fmt.Sprintf(global.YdQBAccWaiting, acId, moneyTmp)
									waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
									//waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
									ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
									if ttl > 0 { //该账号正在冷却中，直接处理删掉
										accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgTmp[0], cid, moneyTmp)
										global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
										global.GVA_LOG.Info("关闭过程校验..账号在冷却中..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									} else {
										accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgTmp[0], cid, moneyTmp)
										global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
										global.GVA_LOG.Info("关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									}
								}()

							}
						} else {

							for _, key := range keys {
								keyTmp := key
								go func() {
									money := strings.Split(keyTmp, ":")[3]
									waitAccYdKey := fmt.Sprintf(global.YdQBAccWaiting, acId, money)
									waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, money)
									waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
									ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
									if ttl > 0 { //该账号正在冷却中，直接处理删掉
										//global.GVA_REDIS.Del(context.Background(), waitAccYdKey)
										global.GVA_LOG.Info("关闭过程校验..账号在冷却中,把账号ttl key干掉", zap.Any("waitMsg", waitMsg))
									} else {
										accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgTmp[0], cid, money)
										global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
										global.GVA_LOG.Info("关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									}
								}()
							}
						}

					} else if global.SdoContains(cid) {

						moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
						moneyList := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
						pattern := fmt.Sprintf(global.ChanOrgSdoAccZSetPrefix, orgTmp[0], cid)
						keys := global.GVA_REDIS.Keys(context.Background(), pattern).Val()

						if len(moneyList) >= len(keys) {
							for _, money := range moneyList {
								moneyTmp := money
								go func() {
									waitAccYdKey := fmt.Sprintf(global.YdSdoAccWaiting, acId, moneyTmp)
									waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
									//waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
									ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
									if ttl > 0 { //该账号正在冷却中，直接处理删掉
										accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgTmp[0], cid, moneyTmp)
										global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
										global.GVA_LOG.Info("关闭过程校验..账号在冷却中..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									} else {
										accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgTmp[0], cid, moneyTmp)
										global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
										global.GVA_LOG.Info("关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									}
								}()

							}
						} else {

							for _, key := range keys {
								keyTmp := key
								go func() {
									money := strings.Split(keyTmp, ":")[3]
									waitAccYdKey := fmt.Sprintf(global.YdSdoAccWaiting, acId, money)
									waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, money)
									waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
									ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
									if ttl > 0 { //该账号正在冷却中，直接处理删掉
										//global.GVA_REDIS.Del(context.Background(), waitAccYdKey)
										global.GVA_LOG.Info("关闭过程校验..账号在冷却中,把账号ttl key干掉", zap.Any("waitMsg", waitMsg))
									} else {
										accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgTmp[0], cid, money)
										global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
										global.GVA_LOG.Info("关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
									}
								}()
							}
						}

					} else if global.J3Contains(cid) {

						waitAccYdKey := fmt.Sprintf(global.YdJ3AccWaiting, acId)
						waitAccMem := fmt.Sprintf("%v_%s_%s", ID, acId, acAccount)
						//waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
						ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
						if ttl > 0 { //该账号正在冷却中，直接处理删掉
							accKey := fmt.Sprintf(global.ChanOrgJ3AccZSet, orgTmp[0], cid)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("关闭过程校验..账号在冷却中..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
						} else {
							accKey := fmt.Sprintf(global.ChanOrgJ3AccZSet, orgTmp[0], cid)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
						}

					} else if global.PcContains(cid) { //QB直付，查一下有没有还没关闭的预产，处理掉
						var pcDBList []vbox.ChannelPayCode
						global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("ac_id = ? and code_status in (2,4)", acId).Find(&pcDBList)
						if len(pcDBList) == 0 {
							global.GVA_LOG.Info("关闭过程校验..暂无需要处理的预产", zap.Any("当前账号", acId))
						} else {

							for _, pcDB := range pcDBList {

								pcKey := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0], pcDB.Cid, pcDB.Money, pcDB.Operator, pcDB.Location)
								pcMem := fmt.Sprintf("%d", pcDB.ID) + "_" + pcDB.Mid + "_" + pcDB.AcAccount + "_" + pcDB.ImgContent

								global.GVA_LOG.Info("关闭过程校验..处理预产关闭匹配", zap.Any("pcKey", pcKey), zap.Any("pcMem", pcMem))

								if pcDB.ExpTime.After(time.Now()) { // 设置的过期时间比当前时间晚，先禁用
									global.GVA_REDIS.ZRem(context.Background(), pcKey, pcMem)
									global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id =?", pcDB.ID).Update("code_status", 5)
									global.GVA_LOG.Info("关闭过程校验..处理预产禁用,更新库状态为5")
								} else { // 过期了，直接删除redis，并且状态置为失效 3
									global.GVA_REDIS.ZRem(context.Background(), pcKey, pcMem)
									global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("id =?", pcDB.ID).Update("code_status", 3)
									global.GVA_LOG.Info("关闭过程校验..处理预产过期了,更新库状态为3")
								}

							}
						}
					}
				}

				if err != nil {
					_ = msg.Reject(false)
					continue
				}
				_ = msg.Ack(true)
			}
			wg.Done()
		}(i + 1)
	}
	global.GVA_LOG.Info("Vbox Acc init 初始化搞定")
	// 等待所有消费者完成处理
	wg.Wait()
}
