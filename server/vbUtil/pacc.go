package vbUtil

import (
	"context"
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"go.uber.org/zap"
	"strings"
	"time"
)

// ValidPacc 通过redis缓存查询商户信息
func ValidPacc(account string, uid *uint) (vpa vbox.PayAccount, err error) {
	paccKey := global.PayAccPrefix + account

	if strings.Contains(account, "TEST") {
		global.GVA_LOG.Info("测试商户检测跳过", zap.Any("入参商户", account))
		vpa = vbox.PayAccount{
			PAccount:  account,
			CreatedBy: *uid,
		}
	} else {
		var count int64
		count, err = global.GVA_REDIS.Exists(context.Background(), paccKey).Result()
		if count == 0 {
			if err != nil {
				global.GVA_LOG.Error("redis err", zap.Error(err))
			}
			global.GVA_LOG.Info("当前缓存池无此商户，查一下库。。。", zap.Any("入参商户ID", account))

			err = global.GVA_DB.Table("vbox_pay_account").
				Where("p_account = ?", account).First(&vpa).Error
			jsonStr, _ := json.Marshal(vpa)
			global.GVA_REDIS.Set(context.Background(), paccKey, jsonStr, 10*time.Minute)
		} else {
			jsonStr, _ := global.GVA_REDIS.Get(context.Background(), paccKey).Bytes()
			err = json.Unmarshal(jsonStr, &vpa)
		}
	}
	return vpa, err
}

func HandlePayUrl2PAcc(orderId string) (string, error) {
	key := global.PAccPay
	var url string

	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		global.GVA_LOG.Warn("当前key不存在", zap.Any("key", key))

		var proxy vbox.Proxy
		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
		err = db.Where("status = ?", 1).Where("chan = ?", key).
			First(&proxy).Error
		if err != nil || proxy.Url == "" {
			return "", err
		}
		url = proxy.Url + orderId

		//global.GVA_REDIS.Set(context.Background(), key, proxy.Url, 0)
		global.GVA_REDIS.Set(context.Background(), key, proxy.Url, 0)
		global.GVA_LOG.Info("查库获取", zap.Any("商户订单地址", url))

		return url, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
	} else {
		var preUrl string
		//preUrl, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		preUrl, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		url = preUrl + orderId
		global.GVA_LOG.Info("缓存池取出", zap.Any("商户订单地址", url))
	}
	return url, err
}
