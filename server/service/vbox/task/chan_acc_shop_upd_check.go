package task

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"log"
	"strconv"
	"strings"
	"sync"
)

const (
	ChanAccShopUpdCheckExchange = "vbox.order.acc_shop_upd_check_exchange"
	ChanAccShopUpdCheckQueue    = "vbox.order.acc_shop_upd_check_queue"
	ChanAccShopUpdCheckKey      = "vbox.order.acc_shop_upd_check"
)

// ChanAccShopUpdCheckTask 店铺地址更新后，查看通道账号资源池
func ChanAccShopUpdCheckTask() {

	// 示例：发送消息
	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		log.Fatalf("Failed to get connection from pool: %v", err)
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 创建 账号开启检查 初始化 消息处理 --------------------
	ch, _ := conn.Channel()
	if err := ch.ExchangeDeclare(ChanAccShopUpdCheckExchange, "direct"); err != nil {
		global.GVA_LOG.Error("create exchange err:", zap.Any("err", err))
	}
	if err := ch.QueueDeclare(ChanAccShopUpdCheckQueue); err != nil {
		global.GVA_LOG.Error("create queue err:", zap.Any("err", err))
	}
	if err := ch.QueueBind(ChanAccShopUpdCheckQueue, ChanAccShopUpdCheckKey, ChanAccShopUpdCheckExchange); err != nil {
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
			deliveries, err := chX.Consume(ChanAccShopUpdCheckQueue, "", false, false, false, false, nil)
			if err != nil {
				global.GVA_LOG.Error("err", zap.Error(err), zap.Any("queue", ChanAccShopUpdCheckQueue))
			}

			for msg := range deliveries {
				//v := &map[string]interface{}{}
				//err := json.Unmarshal(msg.Body, v)
				//global.GVA_LOG.Info(fmt.Sprintf("%v", msg.Body))

				v := string(msg.Body)

				global.GVA_LOG.Info("【引导商铺】收到一条需要处理账号对应的资源情况", zap.Any("info", v))

				split := strings.Split(v, "-")

				//uid := split[1]
				// moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], channelShop.Cid)
				moneyKey := split[0]
				k1 := strings.Split(moneyKey, ":")
				cid := strings.Split(k1[2], "_")[1]
				orgIDStr := strings.Split(k1[1], "_")[1]

				orgID, _ := strconv.Atoi(orgIDStr)
				userIDs := utils2.GetUsersByOrgIds([]uint{uint(orgID)})

				var moneyList []string
				if err = global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
					Where("cid = ? and status = ? and created_by in ?", cid, 1, userIDs).Scan(&moneyList).Error; err != nil {
					global.GVA_LOG.Error("查该组织下数据money异常", zap.Error(err))
					_ = msg.Ack(true)
					continue
				}
				/*cm, err := global.GVA_REDIS.Exists(context.Background(), moneyKey).Result()

				if cm == 0 {
					if err != nil {
						global.GVA_LOG.Error("redis err", zap.Error(err))
					}
					if err = global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
						Where("cid = ? and status = ? and created_by in ?", cid, 1, userIDs).Scan(&moneyList).Error; err != nil {
						global.GVA_LOG.Error("查该组织下数据money异常", zap.Error(err))
						_ = msg.Ack(true)
						continue
					}

					for _, m := range moneyList {
						global.GVA_REDIS.SAdd(context.Background(), moneyKey, m)
					}
					global.GVA_REDIS.Expire(context.Background(), moneyKey, 1*time.Minute)
				} else {
					moneyList = global.GVA_REDIS.SMembers(context.Background(), moneyKey).Val()
				}*/

				// 查一下所有开启状态的通道账号
				var acDBList []vbox.ChannelAccount
				if err = global.GVA_DB.Model(&vbox.ChannelAccount{}).
					Where("cid = ? and status = ? and sys_status = ? and created_by in ?", cid, 1, 1, userIDs).Find(&acDBList).Error; err != nil {
					global.GVA_LOG.Error("查该组织下acDB list 异常", zap.Error(err))
					_ = msg.Ack(true)
					continue
				}

				var moneyOffList []string
				if err = global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
					Where("cid = ? and status = ? and created_by in ?", cid, 0, userIDs).Scan(&moneyOffList).Error; err != nil {
					global.GVA_LOG.Error("查该组织下数据money异常", zap.Error(err))
					_ = msg.Ack(true)
					continue
				}
				global.GVA_LOG.Info("查到该组织下【商铺】money", zap.Any("moneyKey", moneyKey), zap.Any("【开启中】money", moneyList), zap.Any("【关闭中】money", moneyOffList))

				if global.TxContains(cid) { //QB yd

					// 处理重置开启的money
					for _, money := range moneyList {
						accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgID, cid, money)
						//先删掉所有的，重新加一次
						global.GVA_REDIS.Del(context.Background(), accKey)

						for _, acDB := range acDBList {
							accDBTmp := acDB
							moneyTmp := money
							go func(accDBTmp vbox.ChannelAccount, moneyTmp string) {
								ID := accDBTmp.ID
								acId := accDBTmp.AcId
								acAccount := accDBTmp.AcAccount
								waitAccYdKey := fmt.Sprintf(global.YdQBAccWaiting, acId, moneyTmp)
								waitAccMem := fmt.Sprintf("%v,%s,%s,%v", ID, acId, acAccount, moneyTmp)
								waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
								ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
								if ttl > 0 { //该账号正在冷却中
									global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
									cdTime := ttl
									_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
									global.GVA_LOG.Info("商铺开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
								} else {
									global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
									global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
									global.GVA_LOG.Info("商铺开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
								}
							}(accDBTmp, moneyTmp)
						}
					}

					// 处理重置关闭的money
					for _, money := range moneyOffList {
						//for _, acDB := range acDBList {
						//	accDBTmp := acDB
						moneyTmp := money
						go func(moneyTmp string) {
							accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgID, cid, moneyTmp)
							//先删掉所有的，重新加一次
							global.GVA_REDIS.Del(context.Background(), accKey)

							//ID := accDBTmp.ID
							//acId := accDBTmp.AcId
							//acAccount := accDBTmp.AcAccount
							//waitAccYdKey := fmt.Sprintf(global.YdQBAccWaiting, acId, moneyTmp)
							//waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
							////waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
							//ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
							//if ttl > 0 { //该账号正在冷却中，直接处理删掉
							//	global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							//	global.GVA_LOG.Info("商铺关闭过程校验..账号在冷却中..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							//} else {
							//	global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							//	global.GVA_LOG.Info("商铺关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							//}
						}(moneyTmp)
						//}
					}
				} else if global.DnfContains(cid) { // dnf
					// 处理重置开启的money
					for _, money := range moneyList {
						accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgID, cid, money)
						//先删掉所有的，重新加一次
						global.GVA_REDIS.Del(context.Background(), accKey)

						for _, acDB := range acDBList {
							accDBTmp := acDB
							moneyTmp := money
							go func(accDBTmp vbox.ChannelAccount, moneyTmp string) {
								ID := accDBTmp.ID
								acId := accDBTmp.AcId
								acAccount := accDBTmp.AcAccount
								waitAccYdKey := fmt.Sprintf(global.YdDnfAccWaiting, acId, moneyTmp)
								waitAccMem := fmt.Sprintf("%v,%s,%s,%v", ID, acId, acAccount, moneyTmp)
								waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
								ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
								if ttl > 0 { //该账号正在冷却中
									global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
									cdTime := ttl
									_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
									global.GVA_LOG.Info("商铺开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
								} else {
									global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
									global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
									global.GVA_LOG.Info("商铺开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
								}
							}(accDBTmp, moneyTmp)
						}
					}

					// 处理重置关闭的money
					for _, money := range moneyOffList {
						//for _, acDB := range acDBList {
						//	accDBTmp := acDB
						moneyTmp := money
						go func(moneyTmp string) {
							accKey := fmt.Sprintf(global.ChanOrgDnfAccZSet, orgID, cid, moneyTmp)
							//先删掉所有的，重新加一次
							global.GVA_REDIS.Del(context.Background(), accKey)

							//ID := accDBTmp.ID
							//acId := accDBTmp.AcId
							//acAccount := accDBTmp.AcAccount
							//waitAccYdKey := fmt.Sprintf(global.YdDnfAccWaiting, acId, moneyTmp)
							//waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, moneyTmp)
							////waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
							//ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
							//if ttl > 0 { //该账号正在冷却中，直接处理删掉
							//	global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							//	global.GVA_LOG.Info("商铺关闭过程校验..账号在冷却中..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							//} else {
							//	global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							//	global.GVA_LOG.Info("商铺关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							//}
						}(moneyTmp)
						//}
					}
				} else if global.SdoContains(cid) { // sdo

					/*pattern := fmt.Sprintf(global.ChanOrgSdoAccZSetPrefix, orgID, cid)
					keys := global.GVA_REDIS.Keys(context.Background(), pattern).Val()

					global.GVA_LOG.Info("", zap.Any("keys", keys))

					if len(moneyList) >= len(keys) { // 把多出来的新的子集过滤出来，把可用的账号都添加进去
						var subMoneyList []string
						for _, key := range keys {
							m := strings.Split(key, ":")[3]
							myEle := strings.Split(m, "_")[1]
							subMoneyList = append(subMoneyList, myEle)
						}
						global.GVA_LOG.Info("", zap.Any("subMoneyList", subMoneyList))
						// 先过滤出
						notList := utils.FilterNotContains(moneyList, subMoneyList)

						for _, money := range notList {
							for _, acDB := range acDBList {
								ID := acDB.ID
								acId := acDB.AcId
								acAccount := acDB.AcAccount
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
									accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgID, cid, money)
									global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
									global.GVA_LOG.Info("开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
								}
							}
						}

					} else { // 把多出来的新的子集过滤出来，把账号都清除掉
						var keyMoneyList []string
						for _, key := range keys {
							m := strings.Split(key, ":")[3]
							myEle := strings.Split(m, "_")[1]
							keyMoneyList = append(keyMoneyList, myEle)
						}
						global.GVA_LOG.Info("", zap.Any("subMoneyList", keyMoneyList))

						notList := utils.FilterNotContains(keyMoneyList, moneyList)

						for _, money := range notList {
							for _, acDB := range acDBList {
								ID := acDB.ID
								acId := acDB.AcId
								acAccount := acDB.AcAccount
								waitAccYdKey := fmt.Sprintf(global.YdSdoAccWaiting, acId, money)
								waitAccMem := fmt.Sprintf("%v_%s_%s_%v", ID, acId, acAccount, money)
								//waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
								ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
								if ttl > 0 { //该账号正在冷却中，直接处理删掉
									accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgID, cid, money)
									global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
									global.GVA_LOG.Info("关闭过程校验..账号在冷却中..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
								} else {
									accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgID, cid, money)
									global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
									global.GVA_LOG.Info("关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
								}
							}
						}
					}*/

					for _, money := range moneyList {
						accKey := fmt.Sprintf(global.ChanOrgQBAccZSet, orgID, cid, money)
						//先删掉所有的，重新加一次
						global.GVA_REDIS.Del(context.Background(), accKey)

						for _, acDB := range acDBList {
							accDBTmp := acDB
							moneyTmp := money
							go func(accDBTmp vbox.ChannelAccount, moneyTmp string) {
								ID := accDBTmp.ID
								acId := accDBTmp.AcId
								acAccount := accDBTmp.AcAccount
								waitAccYdKey := fmt.Sprintf(global.YdSdoAccWaiting, acId, moneyTmp)

								waitAccMem := fmt.Sprintf("%v,%s,%s,%v", ID, acId, acAccount, moneyTmp)
								waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
								ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
								if ttl > 0 { //该账号正在冷却中
									global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
									cdTime := ttl
									_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
									global.GVA_LOG.Info("商铺开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
								} else {
									global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
									global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
									global.GVA_LOG.Info("商铺开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
								}
							}(accDBTmp, moneyTmp)
						}
					}

					// 处理重置关闭的money
					for _, money := range moneyOffList {
						//for _, acDB := range acDBList {
						//	accDBTmp := acDB
						moneyTmp := money
						go func(moneyTmp string) {
							accKey := fmt.Sprintf(global.ChanOrgSdoAccZSet, orgID, cid, moneyTmp)
							//先删掉所有的，重新加一次
							global.GVA_REDIS.Del(context.Background(), accKey)

							//ID := accDBTmp.ID
							//acId := accDBTmp.AcId
							//acAccount := accDBTmp.AcAccount
							//waitAccYdKey := fmt.Sprintf(global.YdSdoAccWaiting, acId, moneyTmp)
							//waitAccMem := fmt.Sprintf("%v,%s,%s,%v", ID, acId, acAccount, moneyTmp)
							////waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
							//ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
							//if ttl > 0 { //该账号正在冷却中，直接处理删掉
							//	global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							//	global.GVA_LOG.Info("商铺关闭过程校验..账号在冷却中..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							//} else {
							//	global.GVA_REDIS.ZRem(context.Background(), accKey, waitAccMem)
							//	global.GVA_LOG.Info("商铺关闭过程校验..处理掉waitAccMem", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							//}
						}(moneyTmp)
						//}
					}
				} else if global.J3Contains(cid) { // J3

					accKey := fmt.Sprintf(global.ChanOrgJ3AccZSet, orgID, cid)

					//先删掉所有的，重新加一次
					global.GVA_REDIS.Del(context.Background(), accKey)

					for _, acDB := range acDBList {
						acDBTmp := acDB
						go func(acDBTmp vbox.ChannelAccount) {
							ID := acDBTmp.ID
							acId := acDBTmp.AcId
							acAccount := acDBTmp.AcAccount
							waitAccYdKey := fmt.Sprintf(global.YdJ3AccWaiting, acId)
							waitAccMem := fmt.Sprintf("%v,%s,%s", ID, acId, acAccount)
							waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
							ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
							if ttl > 0 { //该账号正在冷却中
								global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
								cdTime := ttl
								_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
								global.GVA_LOG.Info("商铺开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
							} else {
								global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
								global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
								global.GVA_LOG.Info("商铺开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							}
						}(acDBTmp)

					}

				} else if global.WYContains(cid) { // WY

					accKey := fmt.Sprintf(global.ChanOrgWYAccZSet, orgID, cid)

					//先删掉所有的，重新加一次
					global.GVA_REDIS.Del(context.Background(), accKey)

					for _, acDB := range acDBList {
						acDBTmp := acDB
						go func(acDBTmp vbox.ChannelAccount) {
							ID := acDBTmp.ID
							acId := acDBTmp.AcId
							acAccount := acDBTmp.AcAccount
							waitAccYdKey := fmt.Sprintf(global.YdWYAccWaiting, acId)
							waitAccMem := fmt.Sprintf("%v,%s,%s", ID, acId, acAccount)
							waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
							ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
							if ttl > 0 { //该账号正在冷却中
								global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
								cdTime := ttl
								_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
								global.GVA_LOG.Info("商铺开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
							} else {
								global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
								global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
								global.GVA_LOG.Info("商铺开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							}
						}(acDBTmp)

					}

				} else if global.DyContains(cid) { // 抖音

					accKey := fmt.Sprintf(global.ChanOrgDyAccZSet, orgID, cid)

					//先删掉所有的，重新加一次
					global.GVA_REDIS.Del(context.Background(), accKey)

					for _, acDB := range acDBList {
						acDBTmp := acDB
						go func(acDBTmp vbox.ChannelAccount) {
							ID := acDBTmp.ID
							acId := acDBTmp.AcId
							acAccount := acDBTmp.AcAccount
							waitAccYdKey := fmt.Sprintf(global.YdDyAccWaiting, acId)
							waitAccMem := fmt.Sprintf("%v,%s,%s", ID, acId, acAccount)
							waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
							ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
							if ttl > 0 { //该账号正在冷却中
								global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
								cdTime := ttl
								_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
								global.GVA_LOG.Info("商铺开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
							} else {
								global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
								global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
								global.GVA_LOG.Info("商铺开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
							}
						}(acDBTmp)

					}

				} else if global.QNContains(cid) { // qn

					//accKey := fmt.Sprintf(global.ChanOrgQNAccZSet, orgID, cid)
					//
					////先删掉所有的，重新加一次
					//global.GVA_REDIS.Del(context.Background(), accKey)
					//
					//for _, acDB := range acDBList {
					//	acDBTmp := acDB
					//	go func(acDBTmp vbox.ChannelAccount) {
					//		ID := acDBTmp.ID
					//		acId := acDBTmp.AcId
					//		acAccount := acDBTmp.AcAccount
					//		waitAccYdKey := fmt.Sprintf(global.YdQNAccWaiting, acId)
					//		waitAccMem := fmt.Sprintf("%v,%s,%s", ID, acId, acAccount)
					//		waitMsg := strings.Join([]string{waitAccYdKey, waitAccMem}, "-")
					//		ttl := global.GVA_REDIS.TTL(context.Background(), waitAccYdKey).Val()
					//		if ttl > 0 { //该账号正在冷却中
					//			global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 2)
					//			cdTime := ttl
					//			_ = chX.PublishWithDelay(AccCDCheckDelayedExchange, AccCDCheckDelayedRoutingKey, []byte(waitMsg), cdTime)
					//			global.GVA_LOG.Info("商铺开启过程校验..账号在冷却中,发起cd校验任务", zap.Any("waitMsg", waitMsg), zap.Any("cdTime", cdTime))
					//		} else {
					//			global.GVA_DB.Unscoped().Model(&vbox.ChannelAccount{}).Where("id = ?", ID).Update("cd_status", 1)
					//			global.GVA_REDIS.ZAdd(context.Background(), accKey, redis.Z{Score: 0, Member: waitAccMem})
					//			global.GVA_LOG.Info("商铺开启过程校验..置为可用", zap.Any("accKey", accKey), zap.Any("waitAccMem", waitAccMem))
					//		}
					//	}(acDBTmp)
					//
					//}

				} else if global.PcContains(cid) {
					global.GVA_LOG.Info("非引导类，无需处理", zap.Any("cid", cid))
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
	global.GVA_LOG.Info("Vbox Acc shop upd check 初始化搞定")
}
