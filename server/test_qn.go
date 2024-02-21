package main

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
//	"github.com/flipped-aurora/gin-vue-admin/server/utils/http"
//)
//
//func main() {
//
//	content := "https://trade.taobao.com/trade/itemlist/asyncSold.htm?event_submit_do_query=1&_input_charset=utf8&prePageNo=1&sifg=0&tabCode=success&rateStatus=&orderStatus=SUCCESS&payDateBegin=0&buyerNick=&dateEnd=0&rxOldFlag=0&rxSendFlag=0&useCheckcode=false&dateBegin=0&tradeTag=0&rxHasSendFlag=0&auctionType=0&close=0&sellerNick=&notifySendGoodsType=ALL&sellerMemoFlag=0&useOrderInfo=false&logisticsService=&isQnNew=true&pageNum=1&o2oDeliveryType=ALL&rxAuditFlag=0&queryOrder=desc&rxElectronicAuditFlag=0&queryMore=false&payDateEnd=0&rxWaitSendflag=0&sellerMemo=0&rxElectronicAllFlag=0&rxSuccessflag=0&refund=&errorCheckcode=false&isHideNick=true&action=itemlist%2FSoldQueryAction&pageSize=15"
//	var params = map[string]string{
//		"event_submit_do_query": "1",
//		"_input_charset":        "utf8",
//		"prePageNo":             "1",
//		"sifg":                  "0",
//		"tabCode":               "success",
//		"rateStatus":            "",
//		"orderStatus":           "SUCCESS",
//		"payDateBegin":          "0",
//		"buyerNick":             "",
//		"dateEnd":               "0",
//		"rxOldFlag":             "0",
//		"rxSendFlag":            "0",
//		"useCheckcode":          "false",
//		"dateBegin":             "0",
//		"tradeTag":              "0",
//		"rxHasSendFlag":         "0",
//		"auctionType":           "0",
//		"close":                 "0",
//		"sellerNick":            "",
//		"notifySendGoodsType":   "ALL",
//		"sellerMemoFlag":        "0",
//		"useOrderInfo":          "false",
//		"logisticsService":      "",
//		"isQnNew":               "true",
//		"pageNum":               "1",
//		"o2oDeliveryType":       "ALL",
//		"rxAuditFlag":           "0",
//		"queryOrder":            "desc",
//		"rxElectronicAuditFlag": "0",
//		"queryMore":             "false",
//		"payDateEnd":            "0",
//		"rxWaitSendflag":        "0",
//		"sellerMemo":            "0",
//		"rxElectronicAllFlag":   "0",
//		"rxSuccessflag":         "0",
//		"refund":                "",
//		"errorCheckcode":        "false",
//		"isHideNick":            "true",
//		"action":                "itemlist/SoldQueryAction",
//		"pageSize":              "15",
//	}
//	client := http.NewHTTPClient()
//
//	// 创建 HTTP 客户端实例
//
//	var headers = map[string]string{
//		"Content-Type":    "application/json; charset=UTF-8",
//		"Accept-Encoding": "",
//		"Referer":         "https://myseller.taobao.com/",
//		"Cookie":          "DI_T_=CvCyYhs4fx1SHMLxwCHxDHh5AoWsb; unb=291897500; uc1=existShop=true&cookie14=UoYenb4YFmamvg%3D%3D&cookie21=U%2BGCWk%2F7owY2VX8Kt53S0g%3D%3D&cookie16=VFC%2FuZ9az08KUQ56dCrZDlbNdA%3D%3D&pas=0&cookie15=WqG3DMC9VAQiUQ%3D%3D; sn=; uc3=nk2=0uNrG6CNYqo%3D&id2=UUGjOpdJllU9&lg2=URm48syIIVrSKA%3D%3D&vt3=F8dD3eu483YvaNgCmEU%3D; csg=8f8ad041; lgc=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; cancelledSubSites=empty; t=e6bb864b96e9abc006a2fb4c2904f8fa; cookie17=UUGjOpdJllU9; dnk=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; skt=741c052db5d73df4; cookie2=137bc6162691ad18cd7062834d727742; existShop=MTcwODM0MDUxMA%3D%3D; uc4=id4=0%40U2OU9SmOE7zVKGEpEarf3j4m2IY%3D&nk4=0%400FJ7kRcJ2hk1GuZTgWndZh4KIg%3D%3D; publishItemObj=; tracknick=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; _cc_=UtASsssmfA%3D%3D; _l_g_=Ug%3D%3D; sg=%E6%B3%AE00; _nk_=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; cookie1=BYe81wL1aEcXISEX05eGKFLHATjgHSMwpdZLh%2Bonxys%3D; _tb_token_=fb6eb53e76fe1; sgcookie=P100Eh8gaj7W%2FfEtCISh1N7mD%2FjKYPBkkoLgpdCI2aJjLYV8SvsP%2BLXHYcWlzqOD2sD7DQtV1nlRzSEPDv0tiVpcJ%2FvXwHTV3OYEz7RTcGMUINgMExxONdFbTF4Nmh3PmfEf; tbcp=e=UoM%2BHFG%2BH40YFva9%2BW9MM%2Bo%3D&f=UUjZeloosIiw2%2BCvtr5iVE1G0QM%3D; lc=V3ibmdqSa63d8CKAAA%3D%3D; lid=%E5%AE%9D%E5%AE%9D%E6%B3%AE%E6%B3%AE; mtop_partitioned_detect=1; _m_h5_tk=a0eb682d4ca19818d9225628ec43661b_1708348432522; _m_h5_tk_enc=5f1c44697cfd5a2125fc2b47fc696a96; cna=ICNaHp3jZFgCAWonluC97snQ; xlly_s=1; v=0; tfstk=eI1HJ-f1B9JC7zhfWMdQ6bbHr3yTRDOWkghJ2QKzQh-6JLQy2QjlqgxpUkOdZGjCfpBdz_LPZGtswUeQVbfGPGzIwzOpU_xNu3KevBUkrB1OpbKLAQbPDQqYDSFARwOBaoEYYMIl6Q_6Thz4Mw_I72mzkOPxZnLcJiwxdW1-esUHHn0mbXDC-cmwvw-iiN1M8v-aGhhra1YhSFbFFblPsevC4xkZ3PKKNFzR_YMWLFTMD_mwJ7cULnW3SPDbdp86JoUgSxsScbll0P4i36ve5eWV.; isg=BDQ0bfJW7bR8W3l4S6VeIsm7HfKmDVj3Zpw1Fc6Rk7_hOdeD9h7bh00_vXHhwZBP",
//	}
//
//	//prePageNo=1&sifg=0&action=itemlist%2FSoldQueryAction&tabCode=success&rateStatus=&orderStatus=SUCCESS&payDateBegin=0&buyerNick=&pageSize=15&dateEnd=0&rxOldFlag=0&rxSendFlag=0&useCheckcode=false&dateBegin=0&tradeTag=0&rxHasSendFlag=0&auctionType=0&close=0&sellerNick=&notifySendGoodsType=ALL&sellerMemoFlag=0&useOrderInfo=false&logisticsService=&isQnNew=true&pageNum=1&o2oDeliveryType=ALL&rxAuditFlag=0&queryOrder=desc&rxElectronicAuditFlag=0&queryMore=false&payDateEnd=0&rxWaitSendflag=0&sellerMemo=0&rxElectronicAllFlag=0&rxSuccessflag=0&refund=&errorCheckcode=false&isHideNick=true
//	data := `{
//		"prePageNo":             1,
//		"sifg":                  0,
//		"action":                "itemlist/SoldQueryAction",
//		"tabCode":               "success",
//		"rateStatus":            "",
//		"orderStatus":           "SUCCESS",
//		"payDateBegin":          0,
//		"buyerNick":             "",
//		"pageSize":              15,
//		"dateEnd":               0,
//		"rxOldFlag":             0,
//		"rxSendFlag":            0,
//		"useCheckcode":          false,
//		"dateBegin":             0,
//		"tradeTag":              0,
//		"rxHasSendFlag":         0,
//		"auctionType":           0,
//		"close":                 0,
//		"sellerNick":            "",
//		"notifySendGoodsType":   "ALL",
//		"sellerMemoFlag":        0,
//		"useOrderInfo":          false,
//		"logisticsService":      "",
//		"isQnNew":               true,
//		"pageNum":               1,
//		"o2oDeliveryType":       "ALL",
//		"rxAuditFlag":           0,
//		"queryOrder":            "desc",
//		"rxElectronicAuditFlag": 0,
//		"queryMore":             false,
//		"payDateEnd":            0,
//		"rxWaitSendflag":        0,
//		"sellerMemo":            0,
//		"rxElectronicAllFlag":   0,
//		"rxSuccessflag":         0,
//		"refund":                "",
//		"errorCheckcode":        false,
//		"isHideNick":            true
//	}`
//
//	var query product.QueryData
//	err := json.Unmarshal([]byte(data), &query)
//	if err != nil {
//		fmt.Println("Error:", err)
//		return
//	}
//
//	var options = &http.RequestOptions{
//		Headers:      headers,
//		MaxRedirects: 3,
//		PayloadType:  "url",
//		Payload:      query,
//		QueryParams:  params,
//	}
//
//	resp, _ := client.Post(content, options)
//	s := string(resp.Body)
//
//	fmt.Println(s)
//}
