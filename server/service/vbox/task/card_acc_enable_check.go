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
	ChanCardAccEnableCheckExchange = "vbox.order.pool_card_acc_enable_check_exchange"
	ChanCardAccEnableCheckQueue    = "vbox.order.pool_card_acc_enable_check_queue"
	ChanCardAccEnableCheckKey      = "vbox.order.pool_card_acc_enable_check"
)

// ChanCardAccEnableCheckTask 查单池账号开启状态核查
func ChanCardAccEnableCheckTask() {
	var operationRecordService system.OperationRecordService

	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 创建 账号开启检查 初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(ChanCardAccEnableCheckExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(ChanCardAccEnableCheckQueue); err != nil {
		global.GVA_LOG.Error("create queue err:", zap.Any("err", err))
	}
	if err := ch.QueueBind(ChanCardAccEnableCheckQueue, ChanCardAccEnableCheckKey, ChanCardAccEnableCheckExchange); err != nil {
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
			chX, _ := connX.Channel()

			// 说明：执行账号匹配
			deliveries, err := chX.Consume(ChanCardAccEnableCheckQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", ChanCardAccEnableCheckQueue))
			}
			now := time.Now()

			for msg := range deliveries {
				//v := &map[string]interface{}{}
				//err := json.Unmarshal(msg.Body, v)
				//global.GVA_LOG.Info(fmt.Sprintf("%v", msg.Body))

				v := vboxReq.ChanCardAccAndCtx{}
				err := json.Unmarshal(msg.Body, &v)
				if err != nil {
					global.GVA_LOG.Error("错了，直接丢了..." + err.Error())
					_ = msg.Reject(false)
					continue
				}

				orgTmp := utils2.GetSelfOrg(v.Obj.CreatedBy)
				if len(orgTmp) == 0 {
					global.GVA_LOG.Error("当前用户没有组织信息，无法开启账号", zap.Any("ID", v.Obj.ID), zap.Any("cid", v.Obj.Cid), zap.Any("acId", v.Obj.AcId), zap.Any("acAccount", v.Obj.AcAccount))
					_ = msg.Reject(false)
					continue
				}
				ID := v.Obj.ID
				cid := v.Obj.Cid
				acId := v.Obj.AcId
				acAccount := v.Obj.AcAccount
				if v.Obj.Status == 1 {
					global.GVA_LOG.Info("收到一条需要处理的查单池账号【开启】", zap.Any("ID", v.Obj.ID), zap.Any("cid", v.Obj.Cid), zap.Any("acId", v.Obj.AcId), zap.Any("acAccount", v.Obj.AcAccount))

					// 3. 筛选匹配是哪个产品，查一下对应产品的账户是否能够正常官方使用
					if global.ECContains(cid) { //ec
						_, errQ := product.JDValidCookie(v.Obj.Token)
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
								Resp:    fmt.Sprintf(global.CardAccQryRecordsEx, acId, v.Obj.AcAccount),
								UserID:  v.Ctx.UserID,
							}

							errR := operationRecordService.CreateSysOperationRecord(record)
							if errR != nil {
								global.GVA_LOG.Error("当前账号查官方记录异常情况下，record 入库失败..." + errR.Error())
							}

							global.GVA_DB.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", v.Obj.ID).
								Update("sys_status", 0)
							global.GVA_LOG.Warn("当前账号查官方记录异常了，结束...", zap.Any("ac info", v.Obj))
							_ = msg.Reject(false)
							continue
						}
					}

					//2.
					if global.ECContains(cid) { //e card

						accKey := fmt.Sprintf(global.ChanOrgECPoolAccZSet, orgTmp[0], cid)
						waitAccYdKey := fmt.Sprintf(global.YdECPoolAccWaiting, acId)

						waitAccMem := fmt.Sprintf("%v,%s,%s", ID, acId, acAccount)
						waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
						ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
						if ttl > 0 { //该账号正在冷却中
							global.GVA_DB.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", ID).Update("cd_status", 2)
							cdTime := ttl
							_ = chX.PublishWithDelay(CardAccCDCheckDelayedExchange, CardAccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
							global.GVA_LOG.Info("开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
						} else {
							global.GVA_DB.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", ID).Update("cd_status", 1)

							global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
							global.GVA_LOG.Info("开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
						}

					}

					//3. 校验都没啥问题，开启sys_status = 1，即可以调度订单使用
					err = global.GVA_DB.Debug().Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", v.Obj.ID).Update("sys_status", 1).Error
					global.GVA_LOG.Info("校验没啥问题，更新账号sys_status = 1", zap.Any("id", v.Obj.ID))

				} else {
					global.GVA_LOG.Info("收到一条需要处理的账号【关闭】", zap.Any("ID", v.Obj.ID), zap.Any("cid", v.Obj.Cid), zap.Any("acId", v.Obj.AcId), zap.Any("acAccount", v.Obj.AcAccount))

					if global.ECContains(cid) {
						waitAccYdKey := fmt.Sprintf(global.YdECAccWaiting, acId)
						waitAccMem := fmt.Sprintf("%v,%s,%s", ID, acId, acAccount)
						//waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
						ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
						if ttl > 0 { //该账号正在冷却中，直接处理删掉
							accKey := fmt.Sprintf(global.ChanOrgECAccZSet, orgTmp[0], cid)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("关闭过程校验..账号在冷却中..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem), zap.Any("ttl", ttl))
						} else {
							accKey := fmt.Sprintf(global.ChanOrgECAccZSet, orgTmp[0], cid)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
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
	global.GVA_LOG.Info("Vbox Acc card enable init 初始化搞定")
}
