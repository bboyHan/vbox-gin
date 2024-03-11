package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"go.uber.org/zap"
	"log"
	"sync"
)

// 账号开启查询
const (
	ChanQNShopDelCheckExchange = "vbox.order.qn_shop_del_check_exchange"
	ChanQNShopDelCheckQueue    = "vbox.order.qn_shop_del_check_queue"
	ChanQNShopDelCheckKey      = "vbox.order.qn_shop_del_check"
)

// ChanQNShopDelCheckTask 通道账号删除对通道剩余资源核查
func ChanQNShopDelCheckTask() {

	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 创建 账号开启检查 初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(ChanQNShopDelCheckExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(ChanQNShopDelCheckQueue); err != nil {
		global.GVA_LOG.Error("create queue err:", zap.Any("err", err))
	}
	if err := ch.QueueBind(ChanQNShopDelCheckQueue, ChanQNShopDelCheckKey, ChanQNShopDelCheckExchange); err != nil {
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
			if connX == nil {
				global.GVA_LOG.Error("connX is nil")
				return
			}
			chX, _ := connX.Channel()

			// 说明：执行账号匹配
			deliveries, err := chX.Consume(ChanQNShopDelCheckQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", ChanQNShopDelCheckQueue))
			}

			for msg := range deliveries {

				v := vboxReq.ChanQNShopAndCtx{}
				err := json.Unmarshal(msg.Body, &v)
				if err != nil {
					global.GVA_LOG.Error("错了，直接丢了..." + err.Error())
					_ = msg.Reject(false)
					continue
				}

				orgTmp := utils2.GetSelfOrg(v.Obj.CreatedBy)
				ID := v.Obj.ID
				uid := v.Obj.CreatedBy
				money := v.Obj.Money
				mid := v.Obj.ProductId
				markID := v.Obj.MarkId
				URL := v.Obj.Address

				global.GVA_LOG.Info("收到一条需要处理的qn shop【删除】", zap.Any("v", v))

				QNShopKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgTmp[0], "5001", money)
				waitQNShopMem := fmt.Sprintf("%v,%v,%v,%v,%v,%v", ID, uid, money, mid, markID, URL)
				global.GVA_REDIS.ZRem(context.Background(), QNShopKey, waitQNShopMem)
				global.GVA_LOG.Info("qn shop删除过程..处理删除剩余资源", zap.Any("QNShopKey", QNShopKey), zap.Any("waitQNShopMem", waitQNShopMem))

				//3. 校验资源都删的差不多了，然后删号
				if errD := global.GVA_DB.Where("id = ?", ID).Delete(&vbox.ChannelShop{}).Error; errD != nil {
					global.GVA_LOG.Error("qn shop删除过程...处理删除预产,删除预产失败", zap.Error(errD))
				}

				if err != nil {
					_ = msg.Reject(false)
					continue
				}

				global.GVA_LOG.Info("删qn shop，处理资源成功", zap.Any("acc info", v.Obj))
				_ = msg.Ack(true)
			}
			wg.Done()
		}(i + 1)
	}
	// 等待所有消费者完成处理
	wg.Wait()
	global.GVA_LOG.Info("Vbox Acc DEL init 初始化搞定")
}
