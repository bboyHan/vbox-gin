package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

const expirationTime = 240 * time.Second

//func HandleEcAccCKPool() (err error) {
//
//	var accList []vbox.ChannelCardAcc
//	err = global.GVA_DB.Model(&vbox.ChannelCardAcc{}).
//		Where("cid = 6001 and status = 1 AND sys_status = 1").Scan(&accList).Error
//
//	if len(accList) > 0 {
//		for _, acc := range accList {
//			isCkExp, errV := product.JDValidCookie(acc.Token)
//			if errV != nil {
//				global.GVA_LOG.Error("校验ck失败", zap.Error(errV), zap.Any("ID", acc.ID), zap.Any("account", acc.AcAccount))
//				global.GVA_DB.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", acc.ID).Update("sys_status", 0)
//
//				continue
//			}
//
//			if !isCkExp {
//				global.GVA_LOG.Error("校验ck expire", zap.Any("ID", acc.ID), zap.Any("account", acc.AcAccount))
//				global.GVA_DB.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", acc.ID).Update("sys_status", 0)
//				continue
//			}
//
//			//global.GVA_REDIS.SAdd()
//		}
//	}
//
//}

func HandleEcJDCodeAdd() (err error) {

	var count int64
	var preCnt int64

	err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("channel_code = 6001 and created_at BETWEEN NOW() - INTERVAL 1 MINUTE AND NOW()").Count(&count).Error

	preCnt = 30 // 默认值
	if count > 15 {
		// 当count大于60时，每多20个，preCnt也多加20
		extra := (count - 60) / 20 * 20
		preCnt += extra
	}

	nowCnt := global.GVA_REDIS.ZCount(context.Background(), global.YdECJdCodeZSet, "-inf", "inf").Val()

	//global.GVA_LOG.Info("当前jd code 池子", zap.Any("nowCnt", nowCnt), zap.Any("当前并发量级", preCnt), zap.Any("当前并发", count))
	if nowCnt > preCnt {
		//global.GVA_LOG.Warn(fmt.Sprintf("当前池子数量: %d, 并发量级: %d，无需继续获取验证码", nowCnt, preCnt))
		return fmt.Errorf("当前jd code池子够用，无需继续获取验证码")
	}
	/*可用验证码服务器 http://1.12.50.148:8887/jd/slide http://43.136.111.242:8887/jd/slide http://159.75.241.132:8887/jd/slide http://43.138.239.132:8887/jd/slide */
	UrlList := []string{"http://1.12.50.148:8887/jd/slide", "http://43.136.111.242:8887/jd/slide", "http://159.75.241.132:8887/jd/slide", "http://43.138.239.132:8887/jd/slide"}
	index := rand.Intn(4)
	//随机取UrlList中的一个
	client := vbHttp.NewHTTPClient()
	url := UrlList[index]
	//global.GVA_LOG.Info("取用jd code服务器地址", zap.Any("url", url))
	resp, errG := client.Get(url, nil)
	if errG != nil {
		global.GVA_LOG.Error("验证码获取失败", zap.Error(errG))

		return fmt.Errorf("验证码获取失败")
	}

	// 获取返回值
	var result map[string]interface{}
	err = json.Unmarshal(resp.Body, &result)
	if err != nil {
		global.GVA_LOG.Error("验证码获取失败", zap.Error(err))
		return fmt.Errorf("验证码获取失败")
	}

	if _, ok := result["verifyCode"]; !ok {
		global.GVA_LOG.Error("没有verifyCode")
		return fmt.Errorf("验证码获取失败")
	}
	if _, ok := result["sessionId"]; !ok {
		global.GVA_LOG.Error("没有sessionId")
		return fmt.Errorf("验证码获取失败")
	}

	mem := fmt.Sprintf("verifyCode=%s&sessionId=%s", result["verifyCode"], result["sessionId"])
	global.GVA_REDIS.ZAdd(context.Background(), global.YdECJdCodeZSet, redis.Z{
		Score:  float64(time.Now().Add(expirationTime).Unix()),
		Member: mem,
	})
	return nil
}

func HandleEcJDCodeDel() (err error) {
	now := time.Now().Unix()
	_, err = global.GVA_REDIS.ZRemRangeByScore(context.Background(), global.YdECJdCodeZSet, "-inf", fmt.Sprintf("%d", now)).Result()
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return err
	}
	//global.GVA_LOG.Info("删除过期的验证码", zap.Any("result", result))
	return nil
}
