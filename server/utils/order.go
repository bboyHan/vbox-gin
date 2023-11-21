package utils

/*
import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

func HandleNotifyUrl2Test() (string, error) {
	var proxy vbox.Proxy
	db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
	err := db.Where("status = ?", 1).Where("chan = ?", "test_notify").
		First(&proxy).Error
	if err != nil || proxy.Url == "" {
		return "", err
	}
	var url = proxy.Url
	return url, nil
}

func HandleEventType(chanID string) (int, error) {
	// 1-商铺关联，2-付码关联

	chanCode, _ := strconv.Atoi(chanID)
	if chanCode >= 1000 && chanCode <= 1099 {
		return 1, nil
	} else if chanCode >= 2000 && chanCode <= 2099 {
		return 1, nil
	} else if chanCode >= 3000 && chanCode <= 3099 {
		return 2, nil
	}
	return 0, fmt.Errorf("不存在的event类型")
}

func HandleEventID(chanID string, money int, orgIDs []uint) (string, error) {
	// 1-商铺关联，2-付码关联

	chanCode, err := strconv.Atoi(chanID)
	var vsList []vbox.ChannelShop

	var orgShopID string
	var zs []redis.Z
	var key string
	for _, orgID := range orgIDs {
		key = fmt.Sprintf(global.ChanOrgShopZSet, orgID, chanID, money)
		zs, err = global.GVA_REDIS.ZRangeArgsWithScores(context.Background(), redis.ZRangeArgs{
			Key:   key,
			Start: 0,
			Stop:  -1,
		}).Result()
		if err != nil {
			return "", err
		}
		if len(zs) <= 0 {
			continue
		}
		break
	}

	if len(zs) <= 0 {
		return "", fmt.Errorf("该组织配置的资源不足，请核查")
	}

	z := zs[len(zs)-1] //取出最后一个，重新设置utc时间戳
	orgShopID = z.Member.(string)
	global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()), // 重新放进去，score设置最新的时间
		Member: orgShopID,
	})

	if orgShopID != "" {
		return orgShopID, err
	}

	// 假设 redis中没查到，重新查一下库，如果库里也没有，那就真的没有了，直接报错

	for _, orgID := range orgIDs {
		userIDs := utils2.GetUsersByOrgId(orgID)
		err := global.GVA_DB.Model(&vbox.ChannelShop{}).Where("cid = ? and money = ? and status = 1", chanID, money).
			Where("create_by in ?", userIDs).Find(&vsList).Error
		if err != nil {
			return "", err
		}
		if len(vsList) <= 0 {
			continue
		}

		//TODO 设置进 redis 中
		for _, shop := range vsList {
			k := fmt.Sprintf(global.ChanOrgShopZSet, orgID, chanID, money)

			global.GVA_REDIS.ZAdd(context.Background(), k, redis.Z{
				Score:  float64(time.Now().Unix()),
				Member: shop.ProductId + strconv.FormatUint(uint64(shop.ID), 10),
			})
		}
	}

	if chanCode >= 1000 && chanCode <= 1099 {

		return "", nil
	} else if chanCode >= 2000 && chanCode <= 2099 {

		return "", nil
	} else if chanCode >= 3000 && chanCode <= 3099 {

		return "", nil
	}
	return "", err
}
*/
