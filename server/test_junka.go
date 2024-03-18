package main

//import (
//	"encoding/base64"
//	"fmt"
//	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
//	"net/url"
//	"regexp"
//)
//
//func main() {
//
//	client := vbHttp.NewHTTPClient()
//	content := `https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG&product=AAJSUPHNZG010CZ`
//	payload := `__EVENTTARGET=&__EVENTARGUMENT=&__LASTFOCUS=&__VIEWSTATE=rNLyvB4r%2FTh%2B5j0%2BViV6yBNaoTfdGp4C9iZtcQzMAKk%2FJsGH1hOjnMUYeb3vjv6Qnh28T2VzLlg%2FjfQwH1wXhIVEotSoJoZjmxHcaQd%2FgK0WfPmD9dG4GKQitRoS8hkbtIEgHHfivAke8JNgO7CnDf7uL8SpSCzBgxuJ96zfAsr7GjJx4kgIUMsS%2Bns%2Bq%2FdIYh8n2qEz82SQFLm%2FEZqPqKZFH9WQ%2F%2Fl%2BaLX0eqT%2B9ECYFvgbqY3U1GV%2FemdmVQtFuM0PPI8M6p4TUg8Qt127hMaPLS6tk8TQecWDCgzGq6P2jVU%2F%2B1UoRPMU4LnnEw%2BuOzvxhUoGmYwipN%2BZtkir%2FPUDpc5jcZ0AlilIB2YyylFf%2FGOJ4ncZzKGdDbfpDg7gz3Ln7VF7pUorKWmo9uMQ1UClg35HkeVexJH7wlnoAtBRwlQjAholyGZqDwNFSHIBNaj%2F1rDm7U65WhAFUZUiFpOW%2FXqCNClgq6N3YHzZuwOu6hLek0AwFZCaErRZxJw2UDROSQNC1IhJbhQ%2Ff%2Boz0W2iL3lxdVqi25AGkV4yJbPKusTqwdcyxOipMafOKWjPLj6%2Fj%2FNSdTePhPUU4gR1%2B5FijhgupcTv6xEdSCDwUyypKOBg4kOqO6PrnwEHY%2FvyE%2B09RMJX%2BKi3APHgraxbwO0Qa9TOLo%2FVuJZpvNEIZrg8NFPXQDEIC9dHZs6HQJJDDne32uwncshXDHKg0rvcHoheIR14sUJshRyhwkKWqCtoRkzAAfFeBL%2BI6hSKkP8B3j4P5vGYMVhLeQCq6qARxdtsRAQrLGc2LDc9BzSmbWiSuLHtKdnOl8T7jrKKPKsYnaDGpfvr4fjiVfWXb%2B1GuGLKtC2k8ASOK5B1gdNr00oZ2PpXWH1cS8fxroFvMYgWC8mXmOIsguOS3xPiqb4D%2BWMe8xkdLKtq70cn6M4sZqj99f5wd9Ut%2BGK%2ByEIiqJnZXsTMKKIEKAAAvuKf5b1%2BynCIFpDNX9rLF4nVD6IwagSc8Jhwr5VPter5eytjd6k%2Fqz6tEe3ZKvC97qblZwlpqnRAuIS5Tt4bSGALXNjKhoGOLtjniZ6uXxRgr8asoxpJlGmbaAW4h%2FFBXTfu0uBQjDkG4C%2FI71%2FsD6IsmOX5rE08VrYlEcvQfTHvlioDQc7LHymCXW5nMDiL4DyMBaPLyjA3OXX234IWUOcOoaPNbOC%2FVNkv1wg5TbG4tqIvWKnB1SyaTjd0ly7K0PhWrMtCYEE8Xk3zzm%2Bdho7MUxNC9Cdq12XwmmorzpJ92h05gFV4nnzFrvlQhpsLbGNsmNM9Y2p4Ao3fKqzSg0oHrmK8bAUjk1ZW7TGKDJ1YrL%2FWDPALLsKcSRWPeZeVZyJNQp7HDJ3ogge5tiyZOtQlhvL%2BmOcZ90s1PtI1y2dVmMzlHbkUEFnY8zYUWTian3CHkiTyopGrx1iQDwtuGMWzMsEe9KNEqc4mqWzcugd7VWrcZeMHHCpJJchMHaxLHDrKVwUm8ESV2O391KhKKs%2FCseyPMQ%2FeqSqepymz62EyqiEJitby%2Fc3sRhTApbah%2B4ez11wffxZEf9St9gGY5pkdGkeb8Cb1qmL7SAF5ejKFYHVh4XlWkKTKXTC58SOxig0cPnvQwOqRnWQ%2BCaYiXqpQTKaiX3cPE9hz1ciKTFgq%2FxmEjvrQAagT67A%2BQbQycsi%2FB%2Fddly5IhxWR867S%2BFvEN%2B%2FnaCQK%2BJUl192tLrIcykXunDjCKXmd0RogQtsMTUpz5Iwzv0oDL0x%2Bz2NnwmiIRJm4yd3cA7F5moJSidzGMrHf0z5%2BnSdnkg%3D%3D&__VIEWSTATEGENERATOR=DCA79EE2&__VIEWSTATEENCRYPTED=&ctl00%24rblSaleType=3&ctl00%24cboProductList=AAJSUPHNZG010CZ&ctl00%24cboProductNum=1&ctl00%24hidCardPwdSign=&ctl00%24junUCard%24UCardType=OneUCardType&ctl00%24junUCard%24FirstUCardNo=2312081559857109&ctl00%24junUCard%24FirstUCardPassword=2312081559857109&ctl00%24junUCard%24SecondUCardNo=&ctl00%24junUCard%24SecondUCardPassword=&ctl00%24junUCard%24ThirdUCardNo=&ctl00%24junUCard%24ThirdUCardPassword=&ctl00%24MiddleTemplate%24hidFromType=&ctl00%24MiddleTemplate%24hidCategory=AAJSUPHNZG&ctl00%24MiddleTemplate%24hidRegionName=&ctl00%24MiddleTemplate%24hidRegionValue=&ctl00%24MiddleTemplate%24hidServerName=&ctl00%24MiddleTemplate%24hidServerValue=&ctl00%24MiddleTemplate%24junCharge%24txtUserAccount=446794914&ctl00%24MiddleTemplate%24junCharge%24txtUserAccountOk=446794914&ctl00%24MiddleTemplate%24txtRandomCode=l1jq&ctl00%24MiddleTemplate%24btnChargeOK=%E7%AB%8B%E5%8D%B3%E5%85%85%E5%80%BC&ctl00%24txtRealityUserName=&ctl00%24txtRealityIDCard=&ctl00%24txtEmail=`
//	var options = &vbHttp.RequestOptions{
//		MaxRedirects: 0,
//		Payload:      payload,
//		PayloadType:  "url",
//	}
//	re := regexp.MustCompile(`((https?://)[^\s]+)`)
//	urlX := re.FindString(content)
//	resp, err := client.Post(urlX, options)
//	if err != nil {
//		//global.GVA_LOG.Error("err:  ->", zap.Error(err))
//	}
//	htmlBody := string(resp.Body)
//	respHeaders := resp.Headers
//
//	loc := respHeaders["Set-Cookie"]
//	//fmt.Printf("%s", htmlBody)
//	fmt.Println(loc)
//	//fmt.Printf("%s", respHeaders)
//
//	imgURL := "https://www.junka.com/Modules/RandomImage.aspx?ran=5972736"
//	var imgHeaders = map[string]string{
//		"Cookie":                    loc,
//		"sec-ch-ua":                 "Not_A Brand;v=8, Chromium;v=120, Google Chrome;v=120",
//		"sec-ch-ua-mobile":          "?0",
//		"sec-ch-ua-platform":        "Windows",
//		"Upgrade-Insecure-Requests": "1",
//		"Origin":                    "https://www.junka.com",
//		"Content-Type":              "application/x-www-form-urlencoded",
//		"User-Agent":                "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
//		"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
//		"Sec-Fetch-Site":            "same-origin",
//		"Sec-Fetch-Mode":            "navigate",
//		"Sec-Fetch-User":            "?1",
//		"Sec-Fetch-Dest":            "document",
//		"Referer":                   "https://www.junka.com/Official/JNetMapSup/CommonCharge.aspx?category=AAJSUPHNZG&product=AAJSUPHNZG010CZ",
//		"Accept-Encoding":           "gzip, deflate, br, zstd",
//		"Accept-Language":           "zh-CN,zh;q=0.9",
//	}
//	var imgOptions = &vbHttp.RequestOptions{
//		Headers:      imgHeaders,
//		MaxRedirects: 0,
//		Payload:      payload,
//		PayloadType:  "url",
//	}
//	respImg, err := client.Get(imgURL, imgOptions)
//	if err != nil {
//		//global.GVA_LOG.Error("err:  ->", zap.Error(err))
//	}
//
//	// 将字节图片转换为 Base64 字符串
//	base64String := base64.StdEncoding.EncodeToString(respImg.Body)
//
//	//fmt.Println(base64String)
//
//	//localhost:9877/api/cls/100000/file
//	imgParseUrl := "http://127.0.0.1:9877/api/cls/100000/b64"
//
//	var imgPayload = map[string]interface{}{
//		"t_img": base64String,
//	}
//	var imgParseOptions = &vbHttp.RequestOptions{
//		MaxRedirects: 0,
//		Payload:      imgPayload,
//		PayloadType:  "json",
//	}
//	respParse, err := client.Post(imgParseUrl, imgParseOptions)
//	if err != nil {
//		//global.GVA_LOG.Error("err:  ->", zap.Error(err))
//	}
//	if respParse == nil {
//		fmt.Println("请求识别码错误")
//		return
//	}
//	code := string(respParse.Body)
//	fmt.Println(code)
//
//	//if strings.Contains(htmlBody, "v.douyin.com") {
//	//	// 先请求一次,获取 html body
//	//
//	var __VIEWSTATE string
//	re = regexp.MustCompile(`<input[^>]*id="__VIEWSTATE"[^>]*value="([^"]+)"`)
//	match := re.FindStringSubmatch(htmlBody)
//	if match != nil && len(match) > 1 {
//		__VIEWSTATE = match[1]
//		//fmt.Printf("%s", __VIEWSTATE)
//	}
//
//	//<script language="javascript">alert(
//	re = regexp.MustCompile(`<script language="javascript">alert\('(.*?)'\);?</script>`)
//	matches2 := re.FindAllStringSubmatch(htmlBody, -1)
//
//	for _, match2 := range matches2 {
//		if len(match2) > 1 {
//			alertMessage := match2[1]
//			fmt.Println(alertMessage)
//		}
//	}
//
//	account := "446794914"
//	card := "2402071320021191"
//	pwd := "1257101300227187"
//	p := `__EVENTTARGET=&__EVENTARGUMENT=&__LASTFOCUS=&__VIEWSTATE=` + url.QueryEscape(__VIEWSTATE) + `&__VIEWSTATEGENERATOR=DCA79EE2&__VIEWSTATEENCRYPTED=&ctl00%24rblSaleType=3&ctl00%24cboProductList=AAJSUPHNZG010CZ&ctl00%24cboProductNum=1&ctl00%24hidCardPwdSign=&ctl00%24junUCard%24UCardType=OneUCardType
//&ctl00%24junUCard%24FirstUCardNo=` + card + `&ctl00%24junUCard%24FirstUCardPassword=` + pwd + `&ctl00%24junUCard%24SecondUCardNo=&ctl00%24junUCard%24SecondUCardPassword=&ctl00%24junUCard%24ThirdUCardNo=&ctl00%24junUCard%24ThirdUCardPassword=&ctl00%24MiddleTemplate%24hidFromType=&ctl00%24MiddleTemplate%24hidCategory=AAJSUPHNZG&ctl00%24MiddleTemplate%24hidRegionName=&ctl00%24MiddleTemplate%24hidRegionValue=&ctl00%24MiddleTemplate%24hidServerName=&ctl00%24MiddleTemplate%24hidServerValue=
//&ctl00%24MiddleTemplate%24junCharge%24txtUserAccount=` + account + `
//&ctl00%24MiddleTemplate%24junCharge%24txtUserAccountOk=` + account + `
//&ctl00%24MiddleTemplate%24txtRandomCode=` + code + `&ctl00%24MiddleTemplate%24btnChargeOK=%E7%AB%8B%E5%8D%B3%E5%85%85%E5%80%BC&ctl00%24txtRealityUserName=&ctl00%24txtRealityIDCard=&ctl00%24txtEmail=`
//
//	var o = &vbHttp.RequestOptions{
//		Headers:      imgHeaders,
//		MaxRedirects: 0,
//		Payload:      p,
//		PayloadType:  "url",
//	}
//
//	respEnd, err := client.Post(urlX, o)
//	if err != nil {
//		//global.GVA_LOG.Error("err:  ->", zap.Error(err))
//	}
//
//	hb := string(respEnd.Body)
//	//fmt.Println("body", hb)
//
//	re = regexp.MustCompile(`<script language="javascript">alert\('(.*?)'\);?</script>`)
//	matches3 := re.FindAllStringSubmatch(hb, -1)
//
//	for _, match3 := range matches3 {
//		if len(match3) > 1 {
//			alertMessage := match3[1]
//			fmt.Println(alertMessage)
//		}
//	}
//}
