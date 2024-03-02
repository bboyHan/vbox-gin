package vbUtil

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

// GetProdType 获取产品类型
func GetProdType(chanID string) (Type int, err error) {
	key := fmt.Sprintf(global.ProdTypeKey, chanID)
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		var proxy vbox.Proxy
		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
		err = db.Where("status = ?", 1).Where("chan = ?", chanID).First(&proxy).Error
		if err != nil {
			return 0, fmt.Errorf("不存在的event类型")
		}
		Type = *proxy.Type

		global.GVA_REDIS.Set(context.Background(), key, Type, 0)
		global.GVA_LOG.Info("数据库取出该产品的类型", zap.Any("code", chanID), zap.Any("类型", Type))

		return Type, err
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
		return 0, fmt.Errorf("不存在的event类型")
	} else {
		typeStr, errR := global.GVA_REDIS.Get(context.Background(), key).Result()
		if errR != nil {
			return 0, fmt.Errorf("不存在的event类型")
		}
		Type, _ = strconv.Atoi(typeStr)

		return Type, err
	}
}

// GetProductByCode 获取产品信息
func GetProductByCode(chanID string) (product *vbox.ChannelProduct, err error) {
	key := fmt.Sprintf(global.ProductKey, chanID)
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		db := global.GVA_DB.Model(&vbox.ChannelProduct{}).Table("vbox_channel_product")
		err = db.Where("channel_code = ?", chanID).First(&product).Error
		if err != nil {
			return nil, err
		}

		jsonStr, _ := json.Marshal(product)
		global.GVA_REDIS.Set(context.Background(), key, jsonStr, 0)

		global.GVA_LOG.Info("数据库取出该产品的信息", zap.Any("code", chanID), zap.Any("product", product))

		return product, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
		return nil, err
	} else {
		jsonStr, _ := global.GVA_REDIS.Get(context.Background(), key).Bytes()
		err = json.Unmarshal(jsonStr, &product)

		global.GVA_LOG.Info("缓存池取出该产品的信息", zap.Any(" code", chanID), zap.Any("product", product))
		return product, err
	}
}

// GetCDTimeByCode 获取产品code对应的过期时间，CD时间，前置时间
func GetCDTimeByCode(chanID string) (t *product.CDTime, err error) {
	key := fmt.Sprintf(global.ProdCodeKey, chanID)
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		var proxy vbox.Proxy
		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
		err = db.Where("status = ?", 1).Where("chan = ?", chanID).First(&proxy).Error
		if err != nil || proxy.Url == "" {
			return nil, fmt.Errorf("proxy查询产品异常, %v", chanID)
		}
		timeSplit := strings.Split(proxy.Url, ",")

		var dTimeSecond, cdTimeSecond, preTimeSecond int
		if len(timeSplit) == 3 {
			dTimeSecond, _ = strconv.Atoi(timeSplit[0])
			cdTimeSecond, _ = strconv.Atoi(timeSplit[1])
			preTimeSecond, _ = strconv.Atoi(timeSplit[2])
			t = &product.CDTime{
				Duration: time.Duration(dTimeSecond) * time.Second,
				CDTime:   time.Duration(cdTimeSecond) * time.Second,
				PreTime:  time.Duration(preTimeSecond) * time.Second,
			}
		} else {
			durationStr := proxy.Url
			dTimeSecond, _ = strconv.Atoi(durationStr)
			t = &product.CDTime{
				Duration: time.Duration(dTimeSecond) * time.Second,
				CDTime:   0,
				PreTime:  0,
			}
		}
		jsonStr, _ := json.Marshal(t)

		global.GVA_REDIS.Set(context.Background(), key, jsonStr, 0)
		global.GVA_LOG.Info("数据库取出该产品的有效时长", zap.Any("code", chanID),
			zap.Any("过期时间(s)", dTimeSecond), zap.Any("CD时间(s)", cdTimeSecond), zap.Any("前置时间(s)", preTimeSecond))

		return t, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
		return nil, err
	} else {
		jsonByte, errJ := global.GVA_REDIS.Get(context.Background(), key).Bytes()
		errJ = json.Unmarshal(jsonByte, &t)
		if errJ != nil {
			return nil, errJ
		}

		global.GVA_LOG.Info("缓存池取出该产品的有效时长", zap.Any(" code", chanID), zap.Any("过期时间(s)", t.Duration),
			zap.Any("CD时间(s)", t.CDTime), zap.Any("前置时间(s)", t.PreTime))
		return t, err
	}
}
