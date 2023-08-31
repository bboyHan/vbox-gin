package product

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/tx"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"testing"
	"time"
)

func TestHttp(t *testing.T) {
	global.GVA_LOG = core.Zap()
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		url     string
		wantErr bool
	}{
		{
			name:    "GET 11 test",
			args:    args{"731268079"},
			wantErr: false,
		},
		{
			name:    "GET 22 test",
			args:    args{"1337088888"},
			wantErr: false,
		},
	}

	var payments []tx.Payment
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			records := Records(tt.args.d)

			global.GVA_LOG.Info("ret:  ->", zap.Any("qq", tt.args.d), zap.Any("water list", records.WaterList))
			payments = append(payments, records.WaterList...)
		})
	}

	// 建立 Redis 连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // 假设 Redis 服务器地址为 localhost:6379
		Password: "",               // 如果使用密码连接 Redis，请填写对应的密码
		DB:       0,                // 使用默认数据库
	})

	// 存储充值记录
	for _, payment := range payments {
		payTime, _ := strconv.ParseInt(payment.PayTime, 10, 64)
		currentTime := time.Now().Unix()
		period := time.Minute * 23 // 默认过期时间

		// 如果充值记录超过半小时，不进行存储
		var expirationTime int64
		if expirationTime = currentTime - payTime - int64(period.Seconds()); expirationTime < 0 {
			continue
		}

		// 更新过期时间为剩余时间
		duration := time.Duration(expirationTime) * time.Second
		global.GVA_LOG.Info("time:  ->", zap.Any("qq", payment.ProvideID), zap.Any("payTime", payTime), zap.Any("currentTime", currentTime), zap.Any("duration", duration))

		key := fmt.Sprintf("%s_%s_%s", payment.ShowName, payment.OrigPayAmt, payment.QQUin)
		err := rdb.Set(context.Background(), key, payment.OrigPayAmt, duration).Err()
		if err != nil {
			fmt.Printf("存储充值记录失败：%v\n", err)
		}
	}

	keys, err := rdb.Keys(context.Background(), "Q币_1000_*").Result()
	if err != nil {
		global.GVA_LOG.Info("ret:  ->", zap.Any("查询充值记录失败：%v\n", err))
	} else {
		global.GVA_LOG.Info("ret:  ->", zap.Any("keys", keys))
	}

	m := cal(payments)
	global.GVA_LOG.Info("ret:  ->", zap.Any("row", m))
	global.GVA_LOG.Info("ret:  ->", zap.Any("QB", m["Q币"]))
}
