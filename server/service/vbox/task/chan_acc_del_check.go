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
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"strings"
	"sync"
	"time"
)

// 账号开启查询
const (
	ChanAccDelCheckExchange = "vbox.order.acc_del_check_exchange"
	ChanAccDelCheckQueue    = "vbox.order.acc_del_check_queue"
	ChanAccDelCheckKey      = "vbox.order.acc_del_check"
)

// ChanAccDelCheckTask 通道账号删除对通道剩余资源核查
func ChanAccDelCheckTask() {
	var operationRecordService system.OperationRecordService

	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 创建 账号开启检查 初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(ChanAccDelCheckExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(ChanAccDelCheckQueue); err != nil {
		global.GVA_LOG.Error("create queue err:", zap.Any("err", err))
	}
	if err := ch.QueueBind(ChanAccDelCheckQueue, ChanAccDelCheckKey, ChanAccDelCheckExchange); err != nil {
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
			deliveries, err := ch.Consume(ChanAccDelCheckQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", ChanAccDelCheckQueue))
			}

			for msg := range deliveries {

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

				global.GVA_LOG.Info("收到一条需要处理的账号【删除】", zap.Any("v", v))

				//2.
				if global.TxContains(cid) { //QB引导，

					moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
					moneyList := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
					pattern := fmt.Sprintf(global.ChanOrgQBAccZSet, orgTmp[0], cid, "*")
					keys := global.GVA_REDIS.Keys(context.Background(), pattern).Val()

					if len(moneyList) >= len(keys) {
						for _, money := range moneyList {
							moneyTmp := money
							go func() {
								accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgTmp[0], cid, moneyTmp)
								waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
								global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
								global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							}()
						}
					} else {
						for _, key := range keys {
							keyTmp := key
							go func() {
								money := strings.Split(keyTmp, ":")[4]
								waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, money)
								global.GVA_REDIS.ZRem(context.Background(), keyTmp, waitAccMem)
								global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", keyTmp), zap.Any("waitAccMem", waitAccMem))
							}()
						}
					}

				} else if global.J3Contains(cid) { //QB引导，

					moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
					moneyList := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
					accKey := fmt.Sprintf(global.ChanOrgJ3AccZSet, orgTmp[0], cid)

					for _, money := range moneyList {
						moneyTmp := money
						go func() {
							waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
						}()
					}

				} else if global.PcContains(cid) { //QB直付，查一下有没有还没剩余的预产，处理掉
					var pcDBList []vbox.ChannelPayCode
					global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("ac_id = ?", acId).Find(&pcDBList)
					if len(pcDBList) == 0 {
						global.GVA_LOG.Info("账号删除过程...暂无需要处理的预产", zap.Any("当前账号", acId))
					} else {

						var pcIDs []uint
						for _, pcDB := range pcDBList {

							pcKey := fmt.Sprintf(global.ChanOrgPayCodeLocZSet, orgTmp[0], pcDB.Cid, pcDB.Money, pcDB.Operator, pcDB.Location)
							pcMem := fmt.Sprintf("%d", pcDB.ID) + "_" + pcDB.Mid + "_" + pcDB.AcAccount + "_" + pcDB.ImgContent

							global.GVA_LOG.Info("账号删除过程...处理删除预产", zap.Any("pcKey", pcKey), zap.Any("pcMem", pcMem))

							global.GVA_REDIS.ZRem(context.Background(), pcKey, pcMem)
							pcIDs = append(pcIDs, pcDB.ID)
						}
						if len(pcIDs) > 0 {
							//删库
							err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

								if errD := tx.Model(&vbox.ChannelPayCode{}).Where("id in ?", pcIDs).Update("deleted_by", v.Ctx.UserID).Error; errD != nil {
									return errD
								}
								if errD := tx.Where("id in ?", pcIDs).Delete(&vbox.ChannelPayCode{}).Error; errD != nil {
									return errD
								}
								return nil
							})
							if err != nil {
								global.GVA_LOG.Error("账号删除过程...处理删除预产,删除预产失败", zap.Error(err))
							}
						}
					}
				}

				//3. 校验资源都删的差不多了，然后删号
				if errD := global.GVA_DB.Where("id = ?", ID).Delete(&vbox.ChannelAccount{}).Error; errD != nil {
					global.GVA_LOG.Error("账号删除过程...处理删除预产,删除预产失败", zap.Error(errD))
				}

				if err != nil {
					_ = msg.Reject(false)
					continue
				}

				record := sysModel.SysOperationRecord{
					Ip:      v.Ctx.ClientIP,
					Method:  v.Ctx.Method,
					Path:    v.Ctx.UrlPath,
					Agent:   v.Ctx.UserAgent,
					Status:  200,
					Latency: time.Since(time.Now()),
					Resp:    fmt.Sprintf(global.AccDelSuccess, ID, acAccount),
					UserID:  v.Ctx.UserID,
				}
				err = operationRecordService.CreateSysOperationRecord(record)

				global.GVA_LOG.Info("删号，处理资源成功", zap.Any("acc info", v.Obj))
				_ = msg.Ack(true)
			}
			wg.Done()
		}(i + 1)
	}
	global.GVA_LOG.Info("Vbox Acc DEL init 初始化搞定")
	// 等待所有消费者完成处理
	wg.Wait()
}
