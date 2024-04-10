package vbUtil

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

func HandleEventID2chShop(chanID string, money int, orgIDs []uint) (orgShopID string, err error) {
	// 1-商铺关联
	var vsList []vbox.ChannelShop
	if chanID == "6001" || chanID == "1007" {
		orgIDs = []uint{1}
	}
	var zs []redis.Z
	var key string
	for _, orgID := range orgIDs {
		key = fmt.Sprintf(global.ChanOrgShopAddrZSet, orgID, chanID, money)
		zs, err = global.GVA_REDIS.ZRangeArgsWithScores(context.Background(), redis.ZRangeArgs{
			Key:   key,
			Start: 0,
			Stop:  -1,
		}).Result()
		if err != nil {
			return "", err
		}
		if len(zs) <= 0 { // redis 没查到，查一下库
			userIDs := utils2.GetUsersByOrgId(orgID)
			err = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("cid = ? and money = ? and status = 1", chanID, money).
				Where("created_by in ?", userIDs).Find(&vsList).Error
			if err != nil {
				return "", err
			}
			if len(vsList) <= 0 {
				continue
			}

			//如果查到库里有， 设置进 redis 中
			for _, shop := range vsList {
				global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
					Score:  float64(time.Now().Unix()),
					Member: shop.ProductId + "_" + strconv.FormatUint(uint64(shop.ID), 10),
				})
			}
		}
		break
	}

	if len(zs) <= 0 {
		global.GVA_LOG.Info("该组织配置的资源不足，请核查", zap.Any("orgIDs", orgIDs), zap.Any("chanID", chanID), zap.Any("money", money))
		return "", fmt.Errorf("该组织配置的资源不足，请核查")
	}

	z := zs[len(zs)-1] //取出最后一个，重新设置utc时间戳
	orgShopID = z.Member.(string)
	global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()), // 重新放进去，score设置最新的时间
		Member: orgShopID,
	})
	global.GVA_LOG.Info("获取引导商铺匹配信息", zap.Any("orgShopID", orgShopID))

	return orgShopID, err
}

func HandleResourceUrl2chShop(eventID string) (payUrl string, err error) {
	global.GVA_LOG.Info("接收event id", zap.Any("eventID", eventID))
	//1. 如果是引导类的，获取引导地址 - channel shop
	split := strings.Split(eventID, "_")
	if len(split) <= 1 {
		return "", fmt.Errorf("解析商铺prod异常，param: %s", eventID)
	}
	//格式 （prodID_ID）
	ID := split[1]

	var shop vbox.ChannelShop
	err = global.GVA_DB.Debug().Model(&vbox.ChannelShop{}).Where("id = ?", ID).First(&shop).Error
	if err != nil {
		return "", err
	}
	global.GVA_LOG.Info("查出shop", zap.Any("shop ID", ID), zap.Any("money", shop.Money), zap.Any("status", shop.Status))

	cid := shop.Cid

	prodInfo, err := utils.GetProductByCode(cid)
	switch {
	case strings.Contains(prodInfo.ProductId, "jd"):
		payUrl, err = utils.HandleJDUrl(shop.Address)
	case strings.Contains(prodInfo.ProductId, "tb"):
		payUrl, err = utils.HandleTBUrl(shop.Address)
	case strings.Contains(prodInfo.ProductId, "dy"):
		payUrl, err = utils.HandleDYUrl(shop.Address)
	case strings.Contains(prodInfo.ProductId, "jym") || strings.Contains(prodInfo.ProductId, "zfb"):
		payUrl, err = utils.HandleAlipayUrl(shop.Address)
	case strings.Contains(prodInfo.ProductId, "pdd"):
		payUrl, err = utils.HandlePddUrl(shop.Address)
	default:
		payUrl = shop.Address
	}
	if err != nil {
		return "", err
	}
	return payUrl, nil
}
