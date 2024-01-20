package main

//
//import (
//	"fmt"
//	"github.com/flipped-aurora/gin-vue-admin/server/global"
//	"github.com/flipped-aurora/gin-vue-admin/server/utils"
//	"go.uber.org/zap"
//)
//
//func main() {
//
//	//client := vbHttp.NewHTTPClient()
//	token := "https://m.tb.cn/h.5qxWGxH4tGbl2Yk?tk=6WmhWi4OGI2"
//	//var options = &vbHttp.RequestOptions{
//	//	MaxRedirects: 3,
//	//}
//	//resp, err := client.Get(token, options)
//	//if err != nil {
//	//	global.GVA_LOG.Error("err:  ->", zap.Error(err))
//	//}
//	//htmlBody := string(resp.Body)
//	////fmt.Printf("%s", htmlBody)
//	//
//	//if strings.Contains(htmlBody, "m.tb.cn") {
//	//	// 先请求一次,获取 html body
//	//
//	//	re := regexp.MustCompile(`var url = '([^']*)'`)
//	//	match := re.FindStringSubmatch(htmlBody)
//	//	if match != nil && len(match) > 1 {
//	//		tmpUrl := match[1]
//	//		fmt.Printf("%s", tmpUrl)
//	//
//	//		parsedURL, errX := url.Parse(tmpUrl)
//	//		if errX != nil {
//	//			global.GVA_LOG.Warn("无效的 URL:", zap.Error(errX))
//	//		}
//	//
//	//		query := parsedURL.Query()
//	//		itemId := query.Get("id")
//	//		fmt.Printf("\nitemId: %s", itemId)
//	//
//	//		schema := "https://main.m.taobao.com/order/index.html?buildOrderVersion=3.0&skuId=undefined&exParams=%7B%22id%22%3A%22674305212211%22%7D&quantity=1&itemId=" + itemId
//	//
//	//		payUrl := "tbopen://m.taobao.com/tbopen/index.html?h5Url=" + url.QueryEscape(schema)
//	//		fmt.Printf("\npayUrl: %s", payUrl)
//	//
//	//	}
//	//}
//	payUrl, err := utils.HandleTBUrl(token)
//	if err != nil {
//		global.GVA_LOG.Error("err:  ->", zap.Error(err))
//	}
//	fmt.Printf("\npayUrl: %s", payUrl)
//}
