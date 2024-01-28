package main

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
//	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
//	product2 "github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
//	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
//	"strings"
//	"time"
//)
//
//func main() {
//	stateTime := time.UnixMilli(1706443120000)
//	fmt.Println(stateTime)
//	start := time.UnixMilli(1706444545000)
//	end := time.Now()
//	record, _ := product2.QrySdoDaoYuRecordBetween(vbox.ChannelAccount{
//		Token: "https://yaoshi.sdo.com/apipool?system_deviceId=868410048918892-40caee54f545f373&sequence=7&isHttp=0&netFlag=WIFI&method=txz_bs_mixed.dqOrder.list&ticket=cFFHWkDhsIxes6NhvcX09FFaDUm3p%2F4Pk146eTWCS7IdU34mtVVI8rgoCVXnJmTJ9kfSiHkT0BP8SXK1sdeLempgKsItWc2F3FPn3BMsa6stXomxFjDNyaOieJADp3NapOnjl9Qnh7n9zi%2BavTlWAxE45Y9R38iCZz6x98tLMu0%3D&txzDeviceId=868410048918892&sndaId=3577890350&maxCount=10&version=a.9.4.8&timestampMs=-1",
//	}, start, end)
//	fmt.Println(record)
//	vm := product2.Classifier(record)
//	m := vm["彩虹岛"]
//	if rd, ok2 := m["1"]; !ok2 {
//		fmt.Println("不存在")
//
//	} else {
//		//fmt.Println(rd)
//		var flag bool
//		for _, tg := range rd {
//			if strings.Contains(tg, "18210889498") {
//				flag = true
//				break
//			}
//		}
//		if flag {
//			fmt.Println("存在")
//		}
//	}
//
//}
//
//func QrySdoRecord(vca vbox.ChannelAccount) ([]product.SdoOrderRecord, error) {
//	var Url = "https://pay.sdo.com/api/orderlist?page=1&range=7"
//
//	client := vbHttp.NewHTTPClient()
//
//	var headers = map[string]string{
//		"Content-Type": "application/json",
//		"Cookie":       "SNDA_ADRefererSystem_MachineTicket=e123ff5f-38db-4c11-a6f5-4eaa0cc76471; sdo_dw_track=RmT7rL1euukxlRroUkngLw==; nsessionid=dfe97c5604af7dfdb8ea203fba310d15; CAS_LOGIN_STATE=1; SECURE_CAS_LOGIN_STATE=1\n",
//	}
//	var opt = &vbHttp.RequestOptions{
//		Headers:      headers,
//		MaxRedirects: 3,
//	}
//
//	httpResp, err := client.Get(Url, opt)
//	if err != nil {
//		return nil, err
//	}
//
//	var ret product.SdoOrderResponse
//	//fmt.Println(string(httpResp.Body))
//	//
//	err = json.Unmarshal(httpResp.Body, &ret)
//	if err != nil {
//		fmt.Println("err: ", err)
//		return nil, err
//	}
//	if ret.ReturnCode == 0 {
//		return ret.Data.Orders, nil
//	} else {
//		return nil, fmt.Errorf("查询系统异常: %s", ret.ReturnMessage)
//	}
//}
//
//func QryDaoYuRecord(vca vbox.ChannelAccount) ([]product.SdoOrderRecord, error) {
//	var Url = "https://pay.sdo.com/api/orderlist?page=1&range=7"
//
//	client := vbHttp.NewHTTPClient()
//
//	var headers = map[string]string{
//		"Content-Type": "application/json",
//		"Cookie":       "SNDA_ADRefererSystem_MachineTicket=e123ff5f-38db-4c11-a6f5-4eaa0cc76471; sdo_dw_track=RmT7rL1euukxlRroUkngLw==; nsessionid=dfe97c5604af7dfdb8ea203fba310d15; CAS_LOGIN_STATE=1; SECURE_CAS_LOGIN_STATE=1\n",
//	}
//	var opt = &vbHttp.RequestOptions{
//		Headers:      headers,
//		MaxRedirects: 3,
//	}
//
//	httpResp, err := client.Get(Url, opt)
//	if err != nil {
//		return nil, err
//	}
//
//	var ret product.SdoOrderResponse
//	//fmt.Println(string(httpResp.Body))
//	//
//	err = json.Unmarshal(httpResp.Body, &ret)
//	if err != nil {
//		fmt.Println("err: ", err)
//		return nil, err
//	}
//	if ret.ReturnCode == 0 {
//		return ret.Data.Orders, nil
//	} else {
//		return nil, fmt.Errorf("查询系统异常: %s", ret.ReturnMessage)
//	}
//}
//
////
////// Classifier 计算不同类型 - 不同金额 - 记录集合
////func Classifier(payments []product.SdoOrderRecord) map[string]map[string][]string {
////	// 使用map存储不同充值类型下的支付金额和充值账号ID集合（去重）
////	paymentsByTypeAndAmount := make(map[string]map[string][]string)
////	for _, payment := range payments {
////		orderAmount := payment.OrderAmount
////		parts := strings.Split(orderAmount, ".")
////
////		var intStrAmount string
////		if len(parts) > 0 {
////			intStrAmount = parts[0]
////		} else {
////			intStrAmount = orderAmount
////		}
////
////		appName := payment.AppName
////		account := payment.InputOrderUser
////
////		// 检查是否存在对应的充值类型的map
////		if _, ok := paymentsByTypeAndAmount[appName]; !ok {
////			paymentsByTypeAndAmount[appName] = make(map[string][]string)
////		}
////
////		// 添加充值账号ID到对应的支付金额中（去重）
////		ids := paymentsByTypeAndAmount[appName][intStrAmount]
////		exists := false
////		for _, id := range ids {
////			if id == account {
////				exists = true
////				break
////			}
////		}
////		if !exists {
////			paymentsByTypeAndAmount[appName][intStrAmount] = append(ids, account)
////		}
////	}
////
////	// 输出结果
////	//for showName, amounts := range paymentsByTypeAndAmount {
////	//	fmt.Printf("充值类型：%s\n", showName)
////	//	for amount, ids := range amounts {
////	//		fmt.Printf("支付金额：%s，充值账号ID集合：%v\n", amount, ids)
////	//	}
////	//	fmt.Println()
////	//}
////	return paymentsByTypeAndAmount
////}
