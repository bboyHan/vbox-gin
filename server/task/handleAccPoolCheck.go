package task

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"go.uber.org/zap"
	"time"
)

func HandleAccPoolCheck() (err error) {

	conn, errC := mq.MQ.ConnPool.GetConnection()
	if errC != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)
	ch, errN := conn.Channel()
	if errN != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
	}

	orgIDList := utils2.GetAllOrgID()
	for _, orgID := range orgIDList {

		if orgID == 1 {
			continue
		}
		var channelCodeList []string

		userIDs := utils2.GetUsersByOrgId(orgID)
		if len(userIDs) == 0 {
			continue
		}

		key := fmt.Sprintf(global.OrgChanSet, orgID)
		c, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
		if c == 0 {
			var productIds []uint
			if err != nil {
				global.GVA_LOG.Error("当前缓存池无此用户对应的orgIds，redis err", zap.Error(err))
			}
			if err = global.GVA_DB.Model(&vbox.OrgProduct{}).Distinct("channel_product_id").Select("channel_product_id").Where("organization_id in (?)", orgID).Find(&productIds).Error; err != nil {
				global.GVA_LOG.Error("OrgProduct查该组织下数据channel code异常", zap.Error(err))
				continue
			}
			if err = global.GVA_DB.Model(&vbox.ChannelProduct{}).Select("channel_code").Where("id in (?)", productIds).Find(&channelCodeList).Error; err != nil {
				global.GVA_LOG.Error("ChannelProduct查channelCodeList 库数据异常", zap.Error(err))
				continue
			}

			for _, cid := range channelCodeList {
				global.GVA_REDIS.SAdd(context.Background(), key, cid)
			}
			//jsonStr, _ := json.Marshal(channelCodeList)
			//rdConn.Set(context.Background(), key, jsonStr, 10*time.Minute)
			global.GVA_REDIS.Expire(context.Background(), key, 1*time.Minute)
		} else {
			cidList, _ := global.GVA_REDIS.SMembers(context.Background(), key).Result()
			channelCodeList = cidList
		}

		for _, cid := range channelCodeList {
			var moneyList []string
			moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgID, cid)
			cm, err := global.GVA_REDIS.Exists(context.Background(), moneyKey).Result()

			if cm == 0 {
				if err != nil {
					global.GVA_LOG.Error("redis err", zap.Error(err))
				}
				if err = global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
					Where("cid = ? and status = ? and created_by in ?", cid, 1, userIDs).Scan(&moneyList).Error; err != nil {
					global.GVA_LOG.Error("查该组织下数据money异常", zap.Error(err))
					continue
				}

				if moneyList == nil || len(moneyList) == 0 {
					continue
				} else {
					for _, m := range moneyList {
						global.GVA_REDIS.SAdd(context.Background(), moneyKey, m)
					}
					//jsonStr, _ := json.Marshal(moneyList)
					//rdConn.Set(context.Background(), moneyKey, jsonStr, 10*time.Minute)
					global.GVA_REDIS.Expire(context.Background(), moneyKey, 1*time.Minute)
				}
			} else {
				moneyMem, _ := global.GVA_REDIS.SMembers(context.Background(), moneyKey).Result()
				moneyList = moneyMem
			}

			waitMsg := fmt.Sprintf("%s-%d", moneyKey, 1)
			err = ch.Publish(task.ChanAccShopUpdCheckExchange, task.ChanAccShopUpdCheckKey, []byte(waitMsg))
			global.GVA_LOG.Info("发起一次池子刷新任务", zap.Any("waitMsg", waitMsg))
		}

	}

	return err
}
