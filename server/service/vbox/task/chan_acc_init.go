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
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"sync"
	"time"
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
						Resp:    fmt.Sprintf(global.BalanceNotEnough, v.Obj.AcId, v.Obj.AcAccount),
						UserID:  v.Ctx.UserID,
					}

					err = operationRecordService.CreateSysOperationRecord(record)
					if err != nil {
						global.GVA_LOG.Error("余额不足情况下，record 入库失败..." + err.Error())
					}
					err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
						Update("sys_status", 0).Error
					// 不允许开启sys_status， 到这里结束
					_ = msg.Reject(false)
					continue
				}

				// 2. 查询账号是否有超 金额限制，或者笔数限制
				// 2.1 日限制
				if v.Obj.DailyLimit > 0 {
					dailyKey := fmt.Sprintf(global.ChanAccDailyUsed, v.Obj.AcId)
					dailyUsed, err := global.GVA_REDIS.Get(context.Background(), dailyKey).Int()
					if err == redis.Nil { // redis中无，查一下库
						var dailySum int
						// 获取今天的时间范围
						startOfDay := time.Now().UTC().Truncate(24 * time.Hour)
						endOfDay := startOfDay.Add(24 * time.Hour)

						err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("sum(money) as dailySum").
							Where("ac_id = ?", v.Obj.AcId).
							Where("order_status = ? AND created_at BETWEEN ? AND ?", 1, startOfDay, endOfDay).Scan(&dailySum).Error

						if err != nil {
							global.GVA_LOG.Error("当前账号计算日消耗查mysql错误，直接丢了..." + err.Error())
							_ = msg.Reject(false)
							err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0).Error
							continue
						}

						// 查完算一下并且给redis赋值设置一下
						global.GVA_REDIS.Set(context.Background(), dailyKey, dailySum, 0)

					} else if err != nil {
						global.GVA_LOG.Error("当前账号计算日消耗差redis错误，直接丢了..." + err.Error())
						_ = msg.Reject(false)
						continue
					} else { // redis查出来了，直接比一下
						if dailyUsed > v.Obj.DailyLimit { // 如果日消费已经超了，不允许开启了，直接结束

							//入库操作记录
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.AccDailyLimitNotEnough, v.Obj.AcId, v.Obj.AcAccount),
								UserID:  v.Ctx.UserID,
							}

							err = operationRecordService.CreateSysOperationRecord(record)
							if err != nil {
								global.GVA_LOG.Error("当前账号日消耗已经超限额情况下，record 入库失败..." + err.Error())
							}

							err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0).Error

							global.GVA_LOG.Info("当前账号日消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
							_ = msg.Reject(false)
							continue
						}
					}
				}
				// 2.2 总限制
				if v.Obj.TotalLimit > 0 {
					totalKey := fmt.Sprintf(global.ChanAccTotalUsed, v.Obj.AcId)
					totalUsed, err := global.GVA_REDIS.Get(context.Background(), totalKey).Int()
					if err == redis.Nil { // redis中无，查一下库
						var totalSum int

						err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as totalSum").
							Where("ac_id = ?", v.Obj.AcId).
							Where("order_status = ?", 1).Scan(&totalSum).Error

						if err != nil {
							global.GVA_LOG.Error("当前账号计算总消耗查mysql错误，直接丢了..." + err.Error())
							_ = msg.Reject(false)
							continue
						}

						// 查完算一下并且给redis赋值设置一下
						global.GVA_REDIS.Set(context.Background(), totalKey, totalSum, 0)

					} else if err != nil {
						global.GVA_LOG.Error("当前账号计算总消耗差redis错误，直接丢了..." + err.Error())
						_ = msg.Reject(false)

						continue
					} else { // redis查出来了，直接比一下
						if totalUsed > v.Obj.TotalLimit { // 如果总消费已经超了，不允许开启了，直接结束

							//入库操作记录
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.AccTotalLimitNotEnough, v.Obj.AcId, v.Obj.AcAccount),
								UserID:  v.Ctx.UserID,
							}

							err = operationRecordService.CreateSysOperationRecord(record)
							if err != nil {
								global.GVA_LOG.Error("当前账号总消耗已经超限额情况下，record 入库失败..." + err.Error())
							}

							err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0).Error
							global.GVA_LOG.Info("当前账号总消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
							_ = msg.Reject(false)
							continue
						}
					}
				}
				// 2.3 笔数限制
				if v.Obj.CountLimit > 0 {
					countKey := fmt.Sprintf(global.ChanAccCountUsed, v.Obj.AcId)
					countUsed, err := global.GVA_REDIS.Get(context.Background(), countKey).Int()
					if err == redis.Nil { // redis中无，查一下库
						var count int64

						err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("ac_id = ?", v.Obj.AcId).Count(&count).Error

						if err != nil {
							global.GVA_LOG.Error("当前账号笔数消耗查mysql错误，直接丢了..." + err.Error())
							_ = msg.Reject(false)
							continue
						}

						// 查完算一下并且给redis赋值设置一下
						global.GVA_REDIS.Set(context.Background(), countKey, count, 0)

					} else if err != nil {
						global.GVA_LOG.Error("当前账号笔数消耗查redis错误，直接丢了..." + err.Error())
						_ = msg.Reject(false)
						continue
					} else { // redis查出来了，直接比一下
						if countUsed > v.Obj.CountLimit { // 如果笔数消费已经超了，不允许开启了，直接结束

							//入库操作记录
							record := sysModel.SysOperationRecord{
								Ip:      v.Ctx.ClientIP,
								Method:  v.Ctx.Method,
								Path:    v.Ctx.UrlPath,
								Agent:   v.Ctx.UserAgent,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.AccCountLimitNotEnough, v.Obj.AcId, v.Obj.AcAccount),
								UserID:  v.Ctx.UserID,
							}

							err = operationRecordService.CreateSysOperationRecord(record)
							if err != nil {
								global.GVA_LOG.Error("当前账号笔数消耗已经超限额情况下，record 入库失败..." + err.Error())
							}
							err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0).Error
							global.GVA_LOG.Info("当前账号笔数消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
							_ = msg.Reject(false)
							continue
						}
					}
				}

				// 3. 筛选匹配是哪个产品，查一下对应产品的账户是否能够正常官方使用
				if global.TxContains(v.Obj.Cid) { //腾讯

					err := product.QryQQRecords(v.Obj)
					if err != nil {
						global.GVA_LOG.Error("当前账号查官方记录异常情况下，record 入库失败..." + err.Error())
						//入库操作记录
						record := sysModel.SysOperationRecord{
							Ip:      v.Ctx.ClientIP,
							Method:  v.Ctx.Method,
							Path:    v.Ctx.UrlPath,
							Agent:   v.Ctx.UserAgent,
							Status:  500,
							Latency: time.Since(now),
							Resp:    fmt.Sprintf(global.AccQryRecordsEx, v.Obj.AcId, v.Obj.AcAccount),
							UserID:  v.Ctx.UserID,
						}

						err = operationRecordService.CreateSysOperationRecord(record)
						if err != nil {
							global.GVA_LOG.Error("当前账号查官方记录异常情况下，record 入库失败..." + err.Error())
						}

						err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
							Update("sys_status", 0).Error
						global.GVA_LOG.Info("当前账号查官方记录异常了，结束...", zap.Any("ac info", v.Obj))
						_ = msg.Reject(false)
						continue

					}

				} else if global.J3Contains(v.Obj.Cid) { //剑三

				}

				//2. 校验都没啥问题，开启sys_status = 1，即可以调度订单使用
				//global.GVA_LOG.Info("消息 : ", zap.Any("msg", msg.Body))
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
