package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"go.uber.org/zap"
	"time"
)

func HandleShopMoneyAvailable() (err error) {
	//自定义 通道账号可用检测
	rdConn := global.GVA_REDIS.Conn()
	defer rdConn.Close()
	var idList []uint
	// 拿出现在所有付方可用的账户
	err = global.GVA_DB.Model(&vbox.PayAccount{}).Table("vbox_pay_account").
		Select("uid").Where("status = ?", 1).Find(&idList).Error
	if err != nil {
		global.GVA_LOG.Error("查付方库数据异常", zap.Error(err))
		return
	}

	for _, uid := range idList {
		var channelCodeList []string

		// 获取组织ID
		orgIdTemp := utils2.GetSelfOrg(uid)

		orgIds := utils2.GetDeepOrg(uid)
		userIDs := utils2.GetUsersByOrgIds(orgIdTemp)

		// 获取当前组织所拥有的产品列表
		key := fmt.Sprintf(global.OrgChanSet, orgIdTemp[0])
		c, err := rdConn.Exists(context.Background(), key).Result()
		if c == 0 {
			var productIds []uint
			if err != nil {
				global.GVA_LOG.Error("当前缓存池无此用户对应的orgIds，redis err", zap.Error(err))
			}
			if err = global.GVA_DB.Model(&vbox.OrgProduct{}).Distinct("channel_product_id").Select("channel_product_id").Where("organization_id in ?", orgIds).Find(&productIds).Error; err != nil {
				global.GVA_LOG.Error("OrgProduct查该组织下数据channel code异常", zap.Error(err))
				continue
			}
			if err = global.GVA_DB.Model(&vbox.ChannelProduct{}).Select("channel_code").Where("id in ?", productIds).Find(&channelCodeList).Error; err != nil {
				global.GVA_LOG.Error("ChannelProduct查channelCodeList 库数据异常", zap.Error(err))
				continue
			}

			jsonStr, _ := json.Marshal(channelCodeList)
			rdConn.Set(context.Background(), key, jsonStr, 10*time.Minute)
		} else {
			jsonStr, _ := rdConn.Get(context.Background(), key).Bytes()
			err = json.Unmarshal(jsonStr, &channelCodeList)
		}

		// 获取每个产品编码下拥有的money list
		for _, cid := range channelCodeList {
			var moneyList []string
			moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgIdTemp[0], cid)
			cm, err := rdConn.Exists(context.Background(), moneyKey).Result()

			if cm == 0 {
				if err != nil {
					global.GVA_LOG.Error("redis err", zap.Error(err))
				}
				if err = global.GVA_DB.Debug().Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
					Where("cid = ? and status = ? and created_by in ?", cid, 1, userIDs).Scan(&moneyList).Error; err != nil {
					global.GVA_LOG.Error("查该组织下数据money异常", zap.Error(err))
					continue
				}

				jsonStr, _ := json.Marshal(moneyList)
				rdConn.Set(context.Background(), moneyKey, jsonStr, 10*time.Minute)
			} else {
				jsonStr, _ := rdConn.Get(context.Background(), moneyKey).Bytes()
				err = json.Unmarshal(jsonStr, &moneyList)
			}
		}

	}

	return err
}
