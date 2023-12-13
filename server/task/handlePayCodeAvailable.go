package task

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"go.uber.org/zap"
)

func HandlePayCodeAvailable() (err error) {
	//自定义 通道账号可用检测
	var idList []uint
	// 拿出现在所有付方可用的账户
	err = global.GVA_DB.Model(&vbox.ChannelPayCode{}).Table("vbox_channel_pay_code").
		Where("code_status = ?", 2).Find(&idList).Error
	if err != nil {
		global.GVA_LOG.Error("查付方库数据异常", zap.Error(err))
		return
	}

	/*for _, uid := range idList {
		var channelCodeList []string

		// 获取组织ID
		orgIdTemp := utils2.GetSelfOrg(uid)

		orgIds := utils2.GetDeepOrg(uid)
		userIDs := utils2.GetUsersByOrgIds(orgIdTemp)

		// 获取当前组织所拥有的产品列表
		key := fmt.Sprintf(global.OrgChanSet, orgIdTemp[0])
		c, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
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

			for _, cid := range channelCodeList {
				rdConn.SAdd(context.Background(), key, cid)
			}
			//jsonStr, _ := json.Marshal(channelCodeList)
			//rdConn.Set(context.Background(), key, jsonStr, 10*time.Minute)
		} else {
			cidList, _ := rdConn.SMembers(context.Background(), key).Result()
			channelCodeList = cidList
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
				if err = global.GVA_DB.Model(&vbox.ChannelShop{}).Distinct("money").Select("money").
					Where("cid = ? and status = ? and created_by in ?", cid, 1, userIDs).Scan(&moneyList).Error; err != nil {
					global.GVA_LOG.Error("查该组织下数据money异常", zap.Error(err))
					continue
				}

				if moneyList == nil || len(moneyList) == 0 {
					continue
				} else {
					for _, m := range moneyList {
						rdConn.SAdd(context.Background(), moneyKey, m)
					}
					//jsonStr, _ := json.Marshal(moneyList)
					//rdConn.Set(context.Background(), moneyKey, jsonStr, 10*time.Minute)
				}
			} else {
				moneyMem, _ := rdConn.SMembers(context.Background(), moneyKey).Result()
				moneyList = moneyMem
			}
		}

	}*/

	return err
}
