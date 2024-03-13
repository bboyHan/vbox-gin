package vbUtil

import (
	"context"
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"time"
)

// SetLimitWithTime 设置访问次数
func SetLimitWithTime(key string, limit int, expiration time.Duration) (cnt int, err error) {
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if err != nil {
		return cnt, err
	}
	if count == 0 {
		pipe := global.GVA_REDIS.TxPipeline()
		pipe.Incr(context.Background(), key)
		pipe.Expire(context.Background(), key, expiration)
		_, err = pipe.Exec(context.Background())
		return cnt + 1, err
	} else {
		// 次数
		if times, err := global.GVA_REDIS.Get(context.Background(), key).Int(); err != nil {
			return times, err
		} else {
			if times >= limit {
				if t, err := global.GVA_REDIS.PTTL(context.Background(), key).Result(); err != nil {
					return times, errors.New("请求太过频繁，请稍后再试")
				} else {
					return times, errors.New("请求太过频繁, 请 " + t.String() + " 秒后尝试")
				}
			} else {
				return times + 1, global.GVA_REDIS.Incr(context.Background(), key).Err()
			}
		}
	}
}
