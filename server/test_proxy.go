package main

//
//import (
//	"fmt"
//	"github.com/flipped-aurora/gin-vue-admin/server/utils"
//	"github.com/flipped-aurora/gin-vue-admin/server/utils/http"
//)
//
//func main() {
//	//content := "9.43 s@E.uS xfb:/ 01/12 【抖音商城】https://v.douyin.com/iLm9L8P2/ 【24小时自动秒充值】腾讯QQ币充值10QB直充可充王者荣耀DNF点券 长按复制此条消息，打开抖音搜索，查看商品详情！"
//	content := "http://api.shenlongip.com/ip?key=iah1c7fo&pattern=txt&count=1&mr=1&protocol=1&sign=94f84cf83d512b135be2a82f9028d353"
//	////content := "https://v.douyin.com/iLm9L8P2/"
//	//// 使用正则表达式匹配 URL
//	client := http.NewHTTPClient()
//	resp, _ := client.Get(content, nil)
//	s := string(resp.Body)
//	// 输出匹配结果
//	fmt.Println(s)
//	//
//	//client = http.NewHTTPClient(s)
//	trim := utils.Trim(s)
//	fmt.Println(trim)
//
//	clientX := http.NewHTTPClient(trim)
//
//	respX, _ := clientX.Get("https://api.unipay.qq.com/v1/r/1450000186/trade_record_query?CmdCode=query2&SubCmdCode=default&PageNum=1&BeginUnixTime=1675232520&EndUnixTime=1706768520&PageSize=100&SystemType=portal&pf=2199&pfkey=pfkey&session_token=0314CB67-5997-430C-B738-D0BC0B3F8E2E1706768520129&webversion=MidasTradeRecord1.0&openid=637B9AB026859C80DBA62C916619893D&openkey=7362D42AA93C5C3E3B71DBDCAA79C23F&session_id=openid&session_type=kp_accesstoken", nil)
//
//	if respX == nil {
//		fmt.Println("代理不可用了")
//	} else {
//		sX := string(respX.Body)
//		// 输出匹配结果
//		fmt.Println(sX)
//	}
//}
