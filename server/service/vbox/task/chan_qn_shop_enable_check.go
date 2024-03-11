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
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"
)

// 账号开启查询
const (
	ChanQNShopEnableCheckExchange = "vbox.order.qn_shop_enable_check_exchange"
	ChanQNShopEnableCheckQueue    = "vbox.order.qn_shop_enable_check_queue"
	ChanQNShopEnableCheckKey      = "vbox.order.qn_shop_enable_check"
)

// ChanQNShopEnableCheckTask 通道账号开启状态核查
func ChanQNShopEnableCheckTask() {

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
				uid := v.Obj.CreatedBy
				money := v.Obj.Money
				mid := v.Obj.ProductId
				markID := v.Obj.MarkId
				URL := v.Obj.Address

				var vcaList []vbox.ChannelAccount
				global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("created_by =?", uid).Scan(&vcaList)

				if len(vcaList) == 0 {
					global.GVA_LOG.Error("len(vcaList) == 0， 无法开启qn shop")
					_ = msg.Ack(true)
					continue
				}

				if v.Obj.Status == 1 {
					global.GVA_LOG.Info("收到一条需要处理的QN商品【开启】", zap.Any("ID", v.Obj.ID), zap.Any("cid", v.Obj.Cid), zap.Any("markID", v.Obj.MarkId), zap.Any("money", v.Obj.Money))

					waitQNShopYdKey := fmt.Sprintf(global.YdQNShopWaiting, mid, ID)
					//TODO 重要
					waitQNShopMem := fmt.Sprintf("%v,%v,%v,%v,%v,%v", ID, uid, money, mid, markID, URL)
					waitMsg := strings.Join([]string{waitQNShopYdKey, waitQNShopMem}, "-")
					ttl := global.GVA_REDIS.TTL(context.Background(), waitQNShopYdKey).Val()
					if ttl > 0 { //该账号正在冷却中
						global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
						cdTime := ttl
						_ = chX.PublishWithDelay(QNShopCDCheckDelayedExchange, QNShopCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
						global.GVA_LOG.Info("开启过程校验..该QN shop在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
					} else {
						global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
						QNShopKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgTmp[0], "5001", money)
						global.GVA_REDIS.ZAdd(context.Background(), QNShopKey, redis.Z{Score: 0, Member: waitQNShopMem})
						global.GVA_LOG.Info("开启过程校验..QN shop置为可用", zap.Any("QNShopKey", QNShopKey), zap.Any("waitQNShopMem", waitQNShopMem))
					}

				} else {
					global.GVA_LOG.Info("收到一条需要处理的QN商品【关闭】", zap.Any("ID", v.Obj.ID), zap.Any("cid", v.Obj.Cid), zap.Any("markID", v.Obj.MarkId), zap.Any("money", v.Obj.Money))

					waitQNShopYdKey := fmt.Sprintf(global.YdQNShopWaiting, mid, ID)
					waitQNShopMem := fmt.Sprintf("%v,%v,%v,%v,%v,%v", ID, uid, money, mid, markID, URL)

					//waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
					ttl := global.GVA_REDIS.TTL(context.Background(), waitQNShopYdKey).Val()
					if ttl > 0 { //该账号正在冷却中，直接处理删掉
						QNShopKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgTmp[0], "5001", money)
						global.GVA_REDIS.ZRem(context.Background(), QNShopKey, waitQNShopMem)
						global.GVA_LOG.Info("关闭过程校验..QN shop在冷却中..处理掉waitAccMem", zap.Any("QNShopKey", QNShopKey), zap.Any("waitAccMem", waitQNShopMem), zap.Any("ttl", ttl))
					} else {
						QNShopKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgTmp[0], "5001", money)
						global.GVA_REDIS.ZRem(context.Background(), QNShopKey, waitQNShopMem)
						global.GVA_LOG.Info("关闭过程校验..QN shop 处理掉waitAccMem", zap.Any("QNShopKey", QNShopKey), zap.Any("waitQNShopMem", waitQNShopMem))
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
