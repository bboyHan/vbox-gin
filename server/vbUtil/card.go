package vbUtil

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strings"
)

func GetOrgPoolCK(orgId uint, chanID string) (cardAcc *vbox.ChannelCardAcc, err error) {
	accPoolKey := fmt.Sprintf(global.ChanOrgECPoolAccZSet, orgId, chanID)

	var resPoolList []string
	resPoolList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), accPoolKey, &redis.ZRangeBy{
		Min:    "0",
		Max:    "0",
		Offset: 0,
		Count:  -1,
	}).Result()

	if err != nil {
		global.GVA_LOG.Error("引导类匹配异常, redis err", zap.Error(err))

	}
	var cardID, cardAccID, cardAcAccount string

	if resPoolList != nil && len(resPoolList) > 0 {
		//accTmp := resList[0]
		accPoolTmp := utils.RandomElement(resPoolList)
		// 2.1 把账号设置为已用
		global.GVA_REDIS.ZAdd(context.Background(), accPoolKey, redis.Z{
			Score:  1,
			Member: accPoolTmp,
		})
		split := strings.Split(accPoolTmp, ",")
		cardID = split[0]
		cardAccID = split[1]
		cardAcAccount = split[2]
		err = global.GVA_DB.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", cardID).First(&cardAcc).Error
		if err != nil {
			global.GVA_LOG.Error("匹配查单池账号异常", zap.Error(err))
			return nil, fmt.Errorf("匹配查单池账号异常！ err：%v", err.Error())
		}
		global.GVA_LOG.Info("匹配查单池账号ck", zap.Any("cardID", cardID), zap.Any("cardAccID", cardAccID), zap.Any("cardAcAccount", cardAcAccount))
		
	} else {
		global.GVA_LOG.Error("查单池匹配账号CK不足, list size zero", zap.Error(err))
		return nil, fmt.Errorf("查单池匹配账号CK不足, list size zero")
	}

	return cardAcc, err
}
