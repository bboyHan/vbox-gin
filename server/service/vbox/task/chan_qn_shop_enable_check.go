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
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"
	"time"
)

// 账号开启查询
const (
	ChanQNShopEnableCheckExchange = "vbox.order.qn_shop_enable_check_exchange"
	ChanQNShopEnableCheckQueue    = "vbox.order.qn_shop_enable_check_queue"
	ChanQNShopEnableCheckKey      = "vbox.order.qn_shop_enable_check"
)

// ChanQNShopEnableCheckTask 通道账号开启状态核查
func ChanQNShopEnableCheckTask() {
	var operationRecordService system.OperationRecordService

	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 创建 账号开启检查 初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(ChanQNShopEnableCheckExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(ChanQNShopEnableCheckQueue); err != nil {
		global.GVA_LOG.Error("create queue err:", zap.Any("err", err))
	}
	if err := ch.QueueBind(ChanQNShopEnableCheckQueue, ChanQNShopEnableCheckKey, ChanQNShopEnableCheckExchange); err != nil {
		global.GVA_LOG.Error("bind queue err:", zap.Any("err", err))
	}

	// 设置初始消费者数量
	consumerCount := 5
	// 使用 WaitGroup 来等待所有消费者完成处理
	var wg sync.WaitGroup
	wg.Add(consumerCount)

	// 启动多个消费者
	for i := 0; i < consumerCount; i++ {
		go func(consumerID int) {
			connX, errX := mq.MQ.ConnPool.GetConnection()
			if errX != nil {
				//log.Fatalf("Failed to get connection from pool: %v", err)
				global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(errX))
			}
			defer mq.MQ.ConnPool.ReturnConnection(connX)
			if connX == nil {
				global.GVA_LOG.Error("connX == nil", zap.Any("err", errX))
				return
			}
			chX, _ := connX.Channel()

			// 说明：执行账号匹配
			deliveries, err := chX.Consume(ChanQNShopEnableCheckQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", ChanQNShopEnableCheckQueue))
			}
			now := time.Now()

			for msg := range deliveries {

				v := vboxReq.ChanQNShopAndCtx{}
				err := json.Unmarshal(msg.Body, &v)
				if err != nil {
					global.GVA_LOG.Error("错了，直接丢了..." + err.Error())
					_ = msg.Reject(false)
					continue
				}

				orgTmp := utils2.GetSelfOrg(v.Obj.CreatedBy)
				if len(orgTmp) == 0 {
					global.GVA_LOG.Error("当前用户没有组织信息，无法开启账号", zap.Any("ID", v.Obj.ID), zap.Any("cid", v.Obj.Cid), zap.Any("prodId", v.Obj.ProductId), zap.Any("id", v.Obj.ID))
					_ = msg.Reject(false)
					continue
				}

				ID := v.Obj.ID
				cid := v.Obj.Cid
				mid := v.Obj.ProductId
				money := v.Obj.Money
				markID := v.Obj.MarkId
				uid := v.Obj.CreatedBy

				var vcaList []vbox.ChannelAccount
				global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("created_by =?", uid).Scan(&vcaList)

				if len(vcaList) == 0 {
					global.GVA_LOG.Error("len(vcaList) == 0")
					_ = msg.Ack(true)
					continue
				}

				accDB := vcaList[0]
				acId := accDB.AcId
				acAccount := accDB.AcAccount

				if global.QNContains(cid) { //QN引导，
					if v.Obj.Status == 1 {
						global.GVA_LOG.Info("收到一条需要处理的QN商品【开启】", zap.Any("ID", v.Obj.ID), zap.Any("cid", v.Obj.Cid), zap.Any("markID", v.Obj.MarkId), zap.Any("money", v.Obj.Money))

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
								MarkId:  fmt.Sprintf(global.AccRecord, acId),
								Type:    global.AccType,
								Status:  500,
								Latency: time.Since(now),
								Resp:    fmt.Sprintf(global.BalanceNotEnough, acId, acAccount),
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
						if accDB.DailyLimit > 0 {

							var dailySum int

							err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as dailySum").
								Where("ac_id = ?", acId).
								Where("channel_code = ?", cid).
								Where("order_status = ? AND created_at BETWEEN CURDATE() AND CURDATE() + INTERVAL 1 DAY - INTERVAL 1 SECOND", 1).Scan(&dailySum).Error

							if err != nil {
								global.GVA_LOG.Error("当前账号计算日消耗查mysql错误，直接丢了..." + err.Error())
								_ = msg.Reject(false)
								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
									Update("sys_status", 0).Error
								continue
							}

							if dailySum >= accDB.DailyLimit { // 如果日消费已经超了，不允许开启了，直接结束

								//入库操作记录
								record := sysModel.SysOperationRecord{
									Ip:      v.Ctx.ClientIP,
									Method:  v.Ctx.Method,
									Path:    v.Ctx.UrlPath,
									Agent:   v.Ctx.UserAgent,
									MarkId:  fmt.Sprintf(global.AccRecord, acId),
									Type:    global.AccType,
									Status:  500,
									Latency: time.Since(now),
									Resp:    fmt.Sprintf(global.AccDailyLimitNotEnough, acId, accDB.AcAccount, dailySum, accDB.DailyLimit),
									UserID:  v.Ctx.UserID,
								}

								err = operationRecordService.CreateSysOperationRecord(record)
								if err != nil {
									global.GVA_LOG.Error("当前账号日消耗已经超限额情况下，record 入库失败..." + err.Error())
								}

								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
									Update("sys_status", 0).Error

								global.GVA_LOG.Error("【DailyLimit】当前账号日消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
								_ = msg.Reject(false)
								continue
							}
						}
						// 2.2 总限制
						if accDB.TotalLimit > 0 {
							var totalSum int

							err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as totalSum").
								Where("ac_id = ?", acId).
								Where("channel_code = ?", cid).
								Where("order_status = ?", 1).Scan(&totalSum).Error

							if err != nil {
								global.GVA_LOG.Error("当前账号计算总消耗查mysql错误，直接丢了..." + err.Error())
								_ = msg.Reject(false)
								continue
							}

							if totalSum >= accDB.TotalLimit { // 如果总消费已经超了，不允许开启了，直接结束

								//入库操作记录
								record := sysModel.SysOperationRecord{
									Ip:      v.Ctx.ClientIP,
									Method:  v.Ctx.Method,
									Path:    v.Ctx.UrlPath,
									Agent:   v.Ctx.UserAgent,
									MarkId:  fmt.Sprintf(global.AccRecord, acId),
									Type:    global.AccType,
									Status:  500,
									Latency: time.Since(now),
									Resp:    fmt.Sprintf(global.AccTotalLimitNotEnough, acId, accDB.AcAccount, totalSum, accDB.TotalLimit),
									UserID:  v.Ctx.UserID,
								}

								err = operationRecordService.CreateSysOperationRecord(record)
								if err != nil {
									global.GVA_LOG.Error("当前账号总消耗已经超限额情况下，record 入库失败..." + err.Error())
								}

								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
									Update("sys_status", 0).Error
								global.GVA_LOG.Error("【TotalLimit】当前账号总消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
								_ = msg.Reject(false)
								continue
							}
						}
						// 2.3 笔数限制
						if accDB.InCntLimit > 0 {

							var count int64

							err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).
								Where("channel_code = ?", cid).
								Where("ac_id = ? and order_status = ?", acId, 1).Count(&count).Error

							if err != nil {
								global.GVA_LOG.Error("当前账号笔数消耗查mysql错误，直接丢了..." + err.Error())
								_ = msg.Reject(false)
								continue
							}

							if int(count) >= accDB.InCntLimit { // 如果笔数消费已经超了，不允许开启了，直接结束

								//入库操作记录
								record := sysModel.SysOperationRecord{
									Ip:      v.Ctx.ClientIP,
									Method:  v.Ctx.Method,
									Path:    v.Ctx.UrlPath,
									Agent:   v.Ctx.UserAgent,
									MarkId:  fmt.Sprintf(global.AccRecord, acId),
									Type:    global.AccType,
									Status:  500,
									Latency: time.Since(now),
									Resp:    fmt.Sprintf(global.AccInCntLimitNotEnough, acId, accDB.AcAccount, count, accDB.InCntLimit),
									UserID:  v.Ctx.UserID,
								}

								err = operationRecordService.CreateSysOperationRecord(record)
								if err != nil {
									global.GVA_LOG.Error("当前账号笔数消耗已经超限额情况下，record 入库失败..." + err.Error())
								}
								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
									Update("sys_status", 0).Error
								global.GVA_LOG.Error("【InCntLimit】当前账号笔数消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
								_ = msg.Reject(false)
								continue
							}
						}

						// 2.4 拉单限制
						if accDB.CountLimit > 0 {

							var count int64

							err = global.GVA_DB.Debug().Model(&vbox.PayOrder{}).
								Where("channel_code = ?", cid).
								Where("ac_id = ?", acId).Count(&count).Error

							if err != nil {
								global.GVA_LOG.Error("当前账号笔数消耗查mysql错误，直接丢了..." + err.Error())
								_ = msg.Reject(false)
								continue
							}

							if int(count) >= accDB.CountLimit { // 如果笔数消费已经超了，不允许开启了，直接结束

								//入库操作记录
								record := sysModel.SysOperationRecord{
									Ip:      v.Ctx.ClientIP,
									Method:  v.Ctx.Method,
									Path:    v.Ctx.UrlPath,
									Agent:   v.Ctx.UserAgent,
									MarkId:  fmt.Sprintf(global.AccRecord, acId),
									Type:    global.AccType,
									Status:  500,
									Latency: time.Since(now),
									Resp:    fmt.Sprintf(global.AccCountLimitNotEnough, acId, accDB.AcAccount, count, accDB.CountLimit),
									UserID:  v.Ctx.UserID,
								}

								err = operationRecordService.CreateSysOperationRecord(record)
								if err != nil {
									global.GVA_LOG.Error("当前账号笔数消耗已经超限额情况下，record 入库失败..." + err.Error())
								}
								err = global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", v.Obj.ID).
									Update("sys_status", 0).Error
								global.GVA_LOG.Error("【CountLimit】当前账号笔数消耗已经超限额了，结束...", zap.Any("ac info", v.Obj))
								_ = msg.Reject(false)
								continue
							}
						}

						waitAccYdKey := fmt.Sprintf(global.YdQNShopWaiting, mid, ID)
						waitAccMem := fmt.Sprintf("%v,%s,%s,%v,%v", ID, mid, markID, money, uid)
						waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
						ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
						if ttl > 0 { //该账号正在冷却中
							global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
							cdTime := ttl
							_ = chX.PublishWithDelay(QNShopCDCheckDelayedExchange, QNShopCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
							global.GVA_LOG.Info("开启过程校验..该QN shop在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
						} else {
							global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
							accKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgTmp[0], cid, money)
							global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
							global.GVA_LOG.Info("开启过程校验..QN shop置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
						}

					} else {
						global.GVA_LOG.Info("收到一条需要处理的QN商品【关闭】", zap.Any("ID", v.Obj.ID), zap.Any("cid", v.Obj.Cid), zap.Any("markID", v.Obj.MarkId), zap.Any("money", v.Obj.Money))

						waitAccYdKey := fmt.Sprintf(global.YdQNShopWaiting, mid, ID)
						waitAccMem := fmt.Sprintf("%v,%s,%s,%v,%v", ID, mid, markID, money, uid)

						//waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
						ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
						if ttl > 0 { //该账号正在冷却中，直接处理删掉
							accKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgTmp[0], cid, money)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("关闭过程校验..QN shop在冷却中..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem), zap.Any("ttl", ttl))
						} else {
							accKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgTmp[0], cid, money)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("关闭过程校验..QN shop 处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
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
	// 等待所有消费者完成处理
	wg.Wait()
	global.GVA_LOG.Info("Vbox Acc enable init 初始化搞定")
}
