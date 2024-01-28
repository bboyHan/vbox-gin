package main

//
//import (
//	"fmt"
//	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
//	"net/url"
//	"regexp"
//	"strings"
//)
//
//func main() {
//
//	client := vbHttp.NewHTTPClient()
//	content := "https://v.douyin.com/iLm9L8P2/"
//	var options = &vbHttp.RequestOptions{
//		MaxRedirects: 0,
//	}
//	re := regexp.MustCompile(`((https?://)[^\s]+)`)
//	urlX := re.FindString(content)
//	resp, err := client.Get(urlX, options)
//	if err != nil {
//		//global.GVA_LOG.Error("err:  ->", zap.Error(err))
//	}
//	//htmlBody := string(resp.Body)
//	respHeaders := resp.Headers
//	loc := respHeaders["Location"]
//	//fmt.Printf("%s", loc)
//
//	parsedURL, _ := url.Parse(loc)
//
//	query := parsedURL.Query()
//	detail_schema := query.Get("detail_schema")
//	//fmt.Printf("%s", detail_schema)
//
//	//decodedURL, err := url.QueryUnescape(detail_schema)
//
//	replace := strings.ReplaceAll(detail_schema, "sslocal://", "snssdk1128://")
//	fmt.Printf("\n%s", replace)
//
//	//if strings.Contains(htmlBody, "v.douyin.com") {
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
//
//}
