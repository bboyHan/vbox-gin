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
			connX, errX := mq.MQ.ConnPool.GetConnection()
			if errX != nil {
				//log.Fatalf("Failed to get connection from pool: %v", err)
				global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(errX))
			}
			defer mq.MQ.ConnPool.ReturnConnection(connX)
			chX, _ := connX.Channel()

			// 说明：执行账号匹配
			deliveries, err := chX.Consume(ChanAccDelCheckQueue, "", false, false, false, false, nil)
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

				/*msgID := fmt.Sprintf(global.MsgFilterMem, msg.MessageId, acId)
				// 检查消息是否已经被处理过
				exists, errR := global.GVA_REDIS.SIsMember(context.Background(), global.MsgFilterKey, msgID).Result()
				if errR != nil {
					global.GVA_LOG.Error("redis ex", zap.Error(errR))
				}

				if exists {
					// 消息已经被处理过，直接返回
					global.GVA_LOG.Info("消息已经被处理过", zap.Any("msgID", msgID))
					// 消息已经处理过，不再处理
					_ = msg.Ack(false)
					continue
				}
				// 将消息ID添加到已处理集合
				errR = global.GVA_REDIS.SAdd(context.Background(), global.MsgFilterKey, msgID).Err()
				if errR != nil {
					global.GVA_LOG.Error("redis ex", zap.Error(errR))
				}
				global.GVA_LOG.Info("消息首次被处理", zap.Any("msgID", msgID))*/

				var moneyList []string
				userIDs := utils2.GetUsersByOrgIds(orgTmp)

				if err = global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
					Where("cid = ? and created_by in ?", cid, userIDs).Scan(&moneyList).Error; err != nil {
					global.GVA_LOG.Error("查该组织下数据money异常", zap.Error(err))
					_ = msg.Ack(true)
					continue
				}

				//2.
				if global.TxContains(cid) { //QB引导，

					for _, money := range moneyList {
						moneyTmp := money
						go func(moneyTmp string) {
							accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgTmp[0], cid, moneyTmp)
							waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
						}(moneyTmp)
					}
					/*moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
					moneyList := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
					pattern := fmt.Sprintf(global.ChanOrgQBAccZSetPrefix, orgTmp[0], cid)
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
								money := strings.Split(keyTmp, ":")[3]
								waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, money)
								global.GVA_REDIS.ZRem(context.Background(), keyTmp, waitAccMem)
								global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", keyTmp), zap.Any("waitAccMem", waitAccMem))
							}()
						}
					}*/

				} else if global.DnfContains(cid) { //dnf引导，
					for _, money := range moneyList {
						moneyTmp := money
						go func(moneyTmp string) {
							accKey := fmt.Sprintf(global.ChanOrgDnfAccZSet, orgTmp[0], cid, moneyTmp)
							waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
						}(moneyTmp)
					}
				} else if global.SdoContains(cid) { //sdo引导，

					/*moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
					moneyList := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
					pattern := fmt.Sprintf(global.ChanOrgSdoAccZSetPrefix, orgTmp[0], cid)
					keys := global.GVA_REDIS.Keys(context.Background(), pattern).Val()

					if len(moneyList) >= len(keys) {
						for _, money := range moneyList {
							moneyTmp := money
							go func() {
								accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgTmp[0], cid, moneyTmp)
								waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
								global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
								global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							}()
						}
					} else {
						for _, key := range keys {
							keyTmp := key
							go func() {
								money := strings.Split(keyTmp, ":")[3]
								waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, money)
								global.GVA_REDIS.ZRem(context.Background(), keyTmp, waitAccMem)
								global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", keyTmp), zap.Any("waitAccMem", waitAccMem))
							}()
						}
					}*/
					for _, money := range moneyList {
						moneyTmp := money
						go func(moneyTmp string) {
							accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgTmp[0], cid, moneyTmp)
							waitAccMem := fmt.Sprintf("%v,%s,%s,%v", ID, acId, acAccount, moneyTmp)
							global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
						}(moneyTmp)
					}

				} else if global.J3Contains(cid) { //QB引导，

					accKey := fmt.Sprintf(global.ChanOrgJ3AccZSet, orgTmp[0], cid)
					waitAccMem := fmt.Sprintf("%v,%s,%s", ID, acId, acAccount)
					global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
					global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))

				} else if global.ECContains(cid) { //QB引导，

					accKey := fmt.Sprintf(global.ChanOrgECAccZSet, orgTmp[0], cid)
					waitAccMem := fmt.Sprintf("%v,%s,%s", ID, acId, acAccount)
					global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
					global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))

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
	// 等待所有消费者完成处理
	wg.Wait()
	global.GVA_LOG.Info("Vbox Acc DEL init 初始化搞定")
}
