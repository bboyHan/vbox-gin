package product

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"go.uber.org/zap"
	"time"
)

func QrySdoDaoYuRecordBetween(vca vbox.ChannelAccount, start time.Time, end time.Time) ([]product.SdoDaoYuOrderRecord, error) {

	client := vbHttp.NewHTTPClient()

	global.GVA_LOG.Info("sdo 查询 URL", zap.Any("ck url", vca.Token))
	httpResp, err := client.Get(vca.Token, options)
	if err != nil {
		global.GVA_LOG.Error("err:  ->", zap.Error(err))
		return nil, err
	}

	var ret product.SdoDaoYuOrderResponse
	err = json.Unmarshal(httpResp.Body, &ret)
	if err != nil {
		global.GVA_LOG.Error("json.Unmarshal:  ->", zap.Error(err), zap.Any("resp", string(httpResp.Body)))
		return nil, err
	}
	if ret.ReturnCode == 0 {
		data := ret.Data
		var newRecord []product.SdoDaoYuOrderRecord
		for _, ele := range data {

			stateTime := time.UnixMilli(ele.TimestampMs)
			//打印时间
			//global.GVA_LOG.Info("sdo 查询时间", zap.Any("stateTime", stateTime), zap.Any("start", start), zap.Any("end", end), zap.Any("PayStatus", ele.PayStatus))
			if start.Before(stateTime) && end.After(stateTime) && ele.PayStatus == 5 {
				global.GVA_LOG.Info("哟呵，满足条件", zap.Any("ele", ele))
				newRecord = append(newRecord, ele)
			}
		}
		global.GVA_LOG.Info("sdo 筛选后的结果", zap.Any("newRecord len", len(newRecord)))

		return newRecord, nil
	} else {
		global.GVA_LOG.Error("Qry Sdo DaoYu Record 官方查单异常", zap.Any("resp", string(httpResp.Body)))
		return nil, fmt.Errorf("查询系统异常: %s", ret.ReturnMessage)
	}
}

func QrySdoDaoYuRecords(vca vbox.ChannelAccount) ([]product.SdoDaoYuOrderRecord, error) {
	client := vbHttp.NewHTTPClient()

	httpResp, err := client.Get(vca.Token, options)
	if err != nil {
		global.GVA_LOG.Error("err:  ->", zap.Error(err))
		return nil, err
	}

	//fmt.Println(string(httpResp.Body))
	var ret product.SdoDaoYuOrderResponse
	err = json.Unmarshal(httpResp.Body, &ret)
	if err != nil {
		global.GVA_LOG.Error("json.Unmarshal:  ->", zap.Error(err), zap.Any("resp", string(httpResp.Body)))
		//fmt.Println(err)
		return nil, err
	}
	if ret.ReturnCode == 0 {
		data := ret.Data
		var newRecord []product.SdoDaoYuOrderRecord
		for _, ele := range data {
			if ele.PayStatus == 5 {
				newRecord = append(newRecord, ele)
			}
		}
		return newRecord, nil
	} else {
		global.GVA_LOG.Error("Qry Sdo DaoYu Record 官方查单异常", zap.Any("resp", string(httpResp.Body)))
		return nil, fmt.Errorf("查询系统异常: %s", ret.ReturnMessage)
	}
}

func QrySdoRecordBetween(vca vbox.ChannelAccount, start time.Time, end time.Time) ([]product.SdoOrderRecord, error) {
	var Url string

	c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordSdoPrefix).Result()
	if c == 0 {
		var channelCode string
		if global.SdoContains(vca.Cid) { // tx系
			channelCode = "sdo_proxy"
		}

		err = global.GVA_DB.Model(&vbox.Proxy{}).Select("url").
			Where("status = ? and type = ? and chan=?", 1, 1, channelCode).
			First(&Url).Error

		if err != nil {
			return nil, fmt.Errorf("该信道无资源配置")
		}

		global.GVA_REDIS.Set(context.Background(), global.ProductRecordSdoPrefix, Url, 10*time.Minute)

	} else {
		Url, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordSdoPrefix).Result()
	}

	client := vbHttp.NewHTTPClient()

	headers = map[string]string{
		"Content-Type": "application/json",
		"Cookie":       vca.Token,
	}
	var opt = &vbHttp.RequestOptions{
		Headers:      headers,
		MaxRedirects: 3,
	}

	httpResp, err := client.Get(Url, opt)
	if err != nil {
		global.GVA_LOG.Error("err:  ->", zap.Error(err))
		return nil, err
	}

	var ret product.SdoOrderResponse
	err = json.Unmarshal(httpResp.Body, &ret)
	if err != nil {
		global.GVA_LOG.Error("json.Unmarshal:  ->", zap.Error(err))
		return nil, err
	}
	if ret.ReturnCode == 0 {
		data := ret.Data.Orders
		var newRecord []product.SdoOrderRecord
		for _, ele := range data {
			if start.Before(ele.StateTime) && end.After(ele.StateTime) && ele.State == 5 {
				newRecord = append(newRecord, ele)
			}
		}
		return newRecord, nil
	} else {
		global.GVA_LOG.Error("Qry Sdo Record 官方查单异常", zap.Any("record", ret))
		return nil, fmt.Errorf("查询系统异常: %s", ret.ReturnMessage)
	}
}

func QrySdoRecords(vca vbox.ChannelAccount) ([]product.SdoOrderRecord, error) {
	var Url string

	c, err := global.GVA_REDIS.Exists(context.Background(), global.ProductRecordSdoPrefix).Result()
	if c == 0 {
		var channelCode string
		if global.SdoContains(vca.Cid) { // tx系
			channelCode = "sdo_proxy"
		}

		err = global.GVA_DB.Model(&vbox.Proxy{}).Select("url").
			Where("status = ? and type = ? and chan=?", 1, 1, channelCode).
			First(&Url).Error

		if err != nil {
			return nil, fmt.Errorf("该信道无资源配置")
		}

		global.GVA_REDIS.Set(context.Background(), global.ProductRecordSdoPrefix, Url, 10*time.Minute)

	} else {
		Url, _ = global.GVA_REDIS.Get(context.Background(), global.ProductRecordSdoPrefix).Result()
	}

	client := vbHttp.NewHTTPClient()

	headers = map[string]string{
		"Content-Type": "application/json",
		"Cookie":       vca.Token,
	}
	var opt = &vbHttp.RequestOptions{
		Headers:      headers,
		MaxRedirects: 3,
	}

	httpResp, err := client.Get(Url, opt)
	if err != nil {
		global.GVA_LOG.Error("err:  ->", zap.Error(err))
		return nil, err
	}

	var ret product.SdoOrderResponse
	err = json.Unmarshal(httpResp.Body, &ret)
	if err != nil {
		global.GVA_LOG.Error("json.Unmarshal:  ->", zap.Error(err))
		return nil, err
	}
	if ret.ReturnCode == 0 {
		data := ret.Data.Orders
		var newRecord []product.SdoOrderRecord
		for _, ele := range data {
			if ele.State == 5 {
				newRecord = append(newRecord, ele)
			}
		}
		return newRecord, nil
	} else {
		global.GVA_LOG.Error("Qry Sdo Record 官方查单异常", zap.Any("record", ret))
		return nil, fmt.Errorf("查询系统异常: %s", ret.ReturnMessage)
	}
}
