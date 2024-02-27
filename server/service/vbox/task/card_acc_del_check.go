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
	"log"
	"sync"
	"time"
)

// 账号开启查询
const (
	ChanCardAccDelCheckExchange = "vbox.order.pool_card_acc_del_check_exchange"
	ChanCardAccDelCheckQueue    = "vbox.order.pool_card_acc_del_check_queue"
	ChanCardAccDelCheckKey      = "vbox.order.pool_card_acc_del_check"
)

// ChanCardAccDelCheckTask 通道账号删除对通道剩余资源核查
func ChanCardAccDelCheckTask() {
	var operationRecordService system.OperationRecordService

	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 创建 账号开启检查 初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(ChanCardAccDelCheckExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(ChanCardAccDelCheckQueue); err != nil {
		global.GVA_LOG.Error("create queue err:", zap.Any("err", err))
	}
	if err := ch.QueueBind(ChanCardAccDelCheckQueue, ChanCardAccDelCheckKey, ChanCardAccDelCheckExchange); err != nil {
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
			deliveries, err := chX.Consume(ChanCardAccDelCheckQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", ChanCardAccDelCheckQueue))
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

				//var moneyList []string
				//userIDs := utils2.GetUsersByOrgIds(orgTmp)
				//
				//if err = global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
				//	Where("cid = ? and created_by in ?", cid, userIDs).Scan(&moneyList).Error; err != nil {
				//	global.GVA_LOG.Error("查该组织下数据money异常", zap.Error(err))
				//	_ = msg.Ack(true)
				//	continue
				//}

				//2.
				if global.ECContains(cid) { //QB引导，

					accKey := fmt.Sprintf(global.ChanOrgECPoolAccZSet, orgTmp[0], cid)
					waitAccMem := fmt.Sprintf("%v,%s,%s", ID, acId, acAccount)
					global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
					global.GVA_LOG.Info("账号删除过程..处理删除剩余资源", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
				}

				//3. 校验资源都删的差不多了，然后删号
				if errD := global.GVA_DB.Where("id = ?", ID).Delete(&vbox.ChannelCardAcc{}).Error; errD != nil {
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
	global.GVA_LOG.Info("Vbox card Acc pool DEL init 初始化搞定")
}
