package product

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"go.uber.org/zap"
	"net/url"
	"strconv"
	"strings"
	"time"
)

//var rawURL = "https://api.unipay.qq.com/v1/r/1450000186/trade_record_query?" +
//	"CmdCode=query2&SubCmdCode=default&PageNum=1&PageSize=200" +
//	"&BeginUnixTime=1659803532&EndUnixTime=1691339532&SystemType=portal&pf=2199&pfkey=pfkey" +
//	"&from_h5=1&session_token=63F728D4-74CB-4817-9F5D-3C344573837F1691339532798&webversion=MidasTradeRecord1.0&r=0.10077481030292357" +
//	"&openid=446794914&openkey=openkey&session_id=hy_gameid&session_type=st_dummy&__refer=" +
//	"&encrypt_msg=ab00dc01d7748d2ea42b2f24971b6c52ba4ecee8b4b741031ffea3e0775f5e06edb08110ebba54a8dcc93fc9a7ff0a4bee0eb4f6ad2033d3c3b2a90e5d9547d1aa96750a759652b9fe44dbcb0dce4d19&msg_len=76"

// 创建 HTTP 客户端实例
var headers = map[string]string{
	"Content-Type": "application/json",
}

var options = &vbHttp.RequestOptions{
	Headers:      headers,
	MaxRedirects: 3,
}

func QryQQRecordsByID(vca vbox.ChannelAccount, orderID string) (*product.Records, error) {
	var Url string

	c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordQBPrefix).Result()
	if c == 0 {
		var channelCode string
		if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) || global.DnfContains(vca.Cid) { // tx系
			channelCode = "qb_proxy"
		}

		err = global.GVA_DB.Model(&vbox.Proxy{}).Select("url").
			Where("status = ? and chan=?", 1, channelCode).
			First(&Url).Error

		if err != nil {
			return nil, errors.New("该信道无资源配置")
		}

		global.GVA_REDIS.Set(context.Background(), global.ProductRecordQBPrefix, Url, 10*time.Minute)

	} else {
		Url, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordQBPrefix).Result()
	}

	openID, openKey, err := Secret(vca.Token)
	if err != nil {
		return nil, err
	}
	records := RecordsByID(Url, openID, openKey, orderID, 24*time.Hour)
	if records == nil || records.Ret != 0 {
		return nil, errors.New("查询官方记录异常")
	}
	//classifier := Classifier(records.WaterList)
	return records, nil
}

// RecordsByID 获取指定时间内记录（开始时间到结束时间）
func RecordsByID(rawURL string, openID string, openKey string, orderID string, period time.Duration) *product.Records {

	// Records 获取指定时间内记录

	u, _ := url.Parse(rawURL)
	queryParams := u.Query()

	// 获取当前时间
	currentTime := time.Now()
	// 计算半小时前的时间
	//halfHourAgo := currentTime.Add(-30 * time.Minute)
	halfHourAgo := currentTime.Add(-period)

	// 当前时间秒数
	currentSeconds := currentTime.Unix()

	// 将半小时前的时间转换为秒数
	halfHourAgoSeconds := halfHourAgo.Unix()

	queryParams.Set("openid", openID)
	queryParams.Set("openkey", openKey)
	queryParams.Set("BeginUnixTime", strconv.FormatInt(halfHourAgoSeconds, 10))
	queryParams.Set("EndUnixTime", strconv.FormatInt(currentSeconds, 10))
	queryParams.Set("SerialNo", orderID)

	u.RawQuery = queryParams.Encode()
	newURL := u.String()
	client := vbHttp.NewHTTPClient()

	global.GVA_LOG.Info("RecordsByID newURL:  ->", zap.String("newURL", newURL))
	resp, err := client.Get(newURL, options)
	if err != nil {
		global.GVA_LOG.Error("err:  ->", zap.Error(err))
		return nil
	}

	var records product.Records
	err = json.Unmarshal(resp.Body, &records)
	if err != nil {
		global.GVA_LOG.Error("json.Unmarshal:  ->", zap.Error(err))
		return nil
	} else if records.Ret != 0 {
		global.GVA_LOG.Error("官方查单异常:  ->", zap.Any("resp body", string(resp.Body)))
		return nil
	}
	//fmt.Print(records)

	return &records
}

func QryQQRecordsBetween(vca vbox.ChannelAccount, start time.Time, end time.Time) ([]product.Payment, error, string) {
	var Url string

	c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordQBPrefix).Result()
	if c == 0 {
		var channelCode string
		if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) || global.DnfContains(vca.Cid) { // tx系
			channelCode = "qb_proxy"
		}

		err = global.GVA_DB.Model(&vbox.Proxy{}).Select("url").
			Where("status = ? and type = ? and chan=?", 1, 1, channelCode).
			First(&Url).Error

		if err != nil {
			return nil, errors.New("该信道无资源配置"), ""
		}

		global.GVA_REDIS.Set(context.Background(), global.ProductRecordQBPrefix, Url, 10*time.Minute)

	} else {
		Url, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordQBPrefix).Result()
	}

	openID, openKey, err := Secret(vca.Token)
	if err != nil {
		return nil, err, ""
	}
	records, newURL := RecordsBetween(Url, openID, openKey, start, end)
	if records == nil || records.Ret != 0 {
		return nil, errors.New("查询官方记录异常"), ""
	}
	//classifier := Classifier(records.WaterList)
	return records.WaterList, nil, newURL
}

// 校验官方合法性用一下

func QryQQRecords(vca vbox.ChannelAccount) error {
	var Url string

	c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordQBPrefix).Result()
	if c == 0 {
		var channelCode string
		if global.TxContains(vca.Cid) || global.PcContains(vca.Cid) || global.DnfContains(vca.Cid) { // tx系
			channelCode = "qb_proxy"
		}

		err = global.GVA_DB.Debug().Model(&vbox.Proxy{}).Select("url").Where("status = ? and type = ? and chan=?", 1, 1, channelCode).
			First(&Url).Error

		if err != nil {
			return errors.New("该信道无资源配置")
		}

		global.GVA_REDIS.Set(context.Background(), global.ProductRecordQBPrefix, Url, 10*time.Minute)

	} else {
		Url, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordQBPrefix).Result()
	}

	openID, openKey, err := Secret(vca.Token)
	if err != nil {
		return err
	}
	records := Records(Url, openID, openKey, 24*30*time.Hour)
	if records == nil || records.Ret != 0 {
		return errors.New("查询官方记录异常")
	}
	//classifier := product.Classifier(records.WaterList)
	return nil
}

func Secret(token string) (string, string, error) {
	openID := vbHttp.ParseCookie(token, "openid")
	openKey := vbHttp.ParseCookie(token, "openkey")
	if openID != "" && openKey != "" {
		return openID, openKey, nil
	}
	return "", "", errors.New(fmt.Sprintf("tx secret值异常, token : %s", token))
}

// RecordsBetween 获取指定时间内记录（开始时间到结束时间）
func RecordsBetween(rawURL string, openID string, openKey string, start time.Time, end time.Time) (*product.Records, string) {

	u, _ := url.Parse(rawURL)
	queryParams := u.Query()

	// 当前时间秒数
	startSeconds := start.Unix()

	// 将半小时前的时间转换为秒数
	endSeconds := end.Unix()

	queryParams.Set("openid", openID)
	queryParams.Set("openkey", openKey)
	queryParams.Set("BeginUnixTime", strconv.FormatInt(startSeconds, 10))
	queryParams.Set("EndUnixTime", strconv.FormatInt(endSeconds, 10))

	u.RawQuery = queryParams.Encode()
	newURL := u.String()
	client := vbHttp.NewProxyHTTPClient()

	//global.GVA_LOG.Info("当前查询用的url", zap.Any("url", newURL))

	resp, err := client.Get(newURL, options)
	if err != nil {
		global.GVA_LOG.Error("err:  ->", zap.Error(err), zap.Any("resp", resp))
		s := err.Error()
		if strings.Contains(s, "connection") {
			global.GVA_LOG.Info("代理获取异常，清除缓存池，重新获取1次")
			global.GVA_REDIS.Del(context.Background(), global.SysProxyIPPrefix)
			client = vbHttp.NewProxyHTTPClient()
			resp, err = client.Get(newURL, options)
			if err != nil {
				return nil, newURL
			}
		} else {
			return nil, newURL
		}
	}

	var records product.Records
	err = json.Unmarshal(resp.Body, &records)
	if err != nil {
		global.GVA_LOG.Error("json.Unmarshal:  ->", zap.Error(err))
		return nil, newURL
	} else if records.Ret != 0 {
		global.GVA_LOG.Error("官方查单异常:  ->", zap.Any("openID", openID), zap.Any("resp body", string(resp.Body)))
		return nil, newURL
	}
	//fmt.Print(records)

	return &records, newURL
}

// Records 获取指定时间内记录
func Records(rawURL string, openID string, openKey string, period time.Duration) *product.Records {

	u, _ := url.Parse(rawURL)
	queryParams := u.Query()

	// 获取当前时间
	currentTime := time.Now()
	// 计算半小时前的时间
	//halfHourAgo := currentTime.Add(-30 * time.Minute)
	halfHourAgo := currentTime.Add(-period)

	// 当前时间秒数
	currentSeconds := currentTime.Unix()

	// 将半小时前的时间转换为秒数
	halfHourAgoSeconds := halfHourAgo.Unix()

	queryParams.Set("openid", openID)
	queryParams.Set("openkey", openKey)
	queryParams.Set("BeginUnixTime", strconv.FormatInt(halfHourAgoSeconds, 10))
	queryParams.Set("EndUnixTime", strconv.FormatInt(currentSeconds, 10))

	u.RawQuery = queryParams.Encode()
	newURL := u.String()
	client := vbHttp.NewHTTPClient()

	resp, err := client.Get(newURL, options)
	//global.GVA_LOG.Info("Records newURL:  ->", zap.String("newURL", newURL), zap.Any("resp", string(resp.Body)))
	if err != nil {
		if resp == nil {
			global.GVA_LOG.Error("err:  ->", zap.Error(err), zap.Any("resp", resp))
			return nil
		}
		global.GVA_LOG.Error("err:  ->", zap.Error(err), zap.Any("resp", string(resp.Body)))
		return nil
	}

	var records product.Records
	err = json.Unmarshal(resp.Body, &records)
	if err != nil {
		global.GVA_LOG.Error("json.Unmarshal:  ->", zap.Error(err))
		return nil
	} else if records.Ret != 0 {
		global.GVA_LOG.Error("官方查单异常:  ->", zap.Any("openID", openID), zap.Any("resp body", string(resp.Body)))
		return nil
	}
	//fmt.Print(records)

	return &records
}

// ClassifierTx 计算不同类型 - 不同金额 - 记录集合
func ClassifierTx(payments []product.Payment) map[string]map[string][]string {
	// 使用map存储不同充值类型下的支付金额和充值账号ID集合（去重）
	paymentsByTypeAndAmount := make(map[string]map[string][]string)
	for _, payment := range payments {
		amount := payment.PayAmt
		showName := payment.ShowName
		provideID := payment.ProvideID

		// 检查是否存在对应的充值类型的map
		if _, ok := paymentsByTypeAndAmount[showName]; !ok {
			paymentsByTypeAndAmount[showName] = make(map[string][]string)
		}

		// 添加充值账号ID到对应的支付金额中（去重）
		ids := paymentsByTypeAndAmount[showName][amount]
		exists := false
		for _, id := range ids {
			if id == provideID {
				exists = true
				break
			}
		}
		if !exists {
			paymentsByTypeAndAmount[showName][amount] = append(ids, provideID)
		}
	}

	// 输出结果
	//for showName, amounts := range paymentsByTypeAndAmount {
	//	fmt.Printf("充值类型：%s\n", showName)
	//	for amount, ids := range amounts {
	//		fmt.Printf("支付金额：%s，充值账号ID集合：%v\n", amount, ids)
	//	}
	//	fmt.Println()
	//}
	return paymentsByTypeAndAmount
}
