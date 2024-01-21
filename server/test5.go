package main

import (
	"fmt"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"regexp"
)

func main() {

	client := vbHttp.NewHTTPClient()
	content := "https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG&product=AAJSUPHNZG010CZ"
	var options = &vbHttp.RequestOptions{
		MaxRedirects: 0,
	}
	re := regexp.MustCompile(`((https?://)[^\s]+)`)
	urlX := re.FindString(content)
	resp, err := client.Get(urlX, options)
	if err != nil {
		//global.GVA_LOG.Error("err:  ->", zap.Error(err))
	}
	htmlBody := string(resp.Body)
	//respHeaders := resp.Headers
	//loc := respHeaders["Location"]
	fmt.Printf("%s", htmlBody)

	//if strings.Contains(htmlBody, "v.douyin.com") {
	//	// 先请求一次,获取 html body
	//
	re = regexp.MustCompile(`<input[^>]*id="__VIEWSTATE"[^>]*value="([^"]+)"`)
	match := re.FindStringSubmatch(htmlBody)
	if match != nil && len(match) > 1 {
		tmpUrl := match[1]
		fmt.Printf("%s", tmpUrl)
	}
	//
	//		parsedURL, errX := url.Parse(tmpUrl)
	//		if errX != nil {
	//			global.GVA_LOG.Warn("无效的 URL:", zap.Error(errX))
	//		}
	//
	//		query := parsedURL.Query()
	//		itemId := query.Get("id")
	//		fmt.Printf("\nitemId: %s", itemId)
	//
	//		schema := "https://main.m.taobao.com/order/index.html?buildOrderVersion=3.0&skuId=undefined&exParams=%7B%22id%22%3A%22674305212211%22%7D&quantity=1&itemId=" + itemId
	//
	//		payUrl := "tbopen://m.taobao.com/tbopen/index.html?h5Url=" + url.QueryEscape(schema)
	//		fmt.Printf("\npayUrl: %s", payUrl)
	//
	//	}
	//}

}
