package main

//import (
//	"fmt"
//	"net/url"
//	"regexp"
//	"strings"
//)
//
//func main() {
//	htmlContent := `<form method="GET" name="form1" action="https://mapi.alipay.com/gateway.do">
//<input type=hidden name="_input_charset" value="GBK"><input type=hidden name="body" value="为qqu****8888巨人游戏充值10元"><input type=hidden name="notify_url" value="http://paybg.mztgame.com/alipay/fillin.php"><input type=hidden name="out_trade_no" value="317118974348548650"><input type=hidden name="partner" value="2088511387849079"><input type=hidden name="payment_type" value="1"><input type=hidden name="return_url" value="//pay.ztgame.com/p.php"><input type=hidden name="seller_email" value="giant101@ztgame.com"><input type=hidden name="service" value="create_direct_pay_by_user"><input type=hidden name="subject" value="为qqu****8888巨人游戏充值10元"><input type=hidden name="total_fee" value="10.00"><input type=hidden name="paymethod" value="directPay"><input type=hidden name="sign" value="oKTIovH1QqodbBAdCegmhlSiJj/Dcg7YXUq4Mi2LqlSJFmgWaHcGrCcxAGvemvfD2yXp7BlAui8pRUhdWJUNQL1U5Qm6gvLEN/zgRDOfij0i1zfZGC5GP+HAc+nrL3RzkAJlYBMo8W0TCzGD9pavdkCoJBxY8XOeX0aXE8M3wyo="><input type=hidden name="sign_type" value="RSA"><input type=hidden name="exter_invoke_ip" value="106.39.149.122"><input type=hidden name="anti_phishing_key" value="KP7vhI6RefplnOiLAA=="><input type=hidden name="extend_param" value="security_wbinfo^{'regDate':'2024-03-15','regName':'637561247','depAccount':'637561247'}"><input type=hidden name="it_b_pay" value="5m">
//<script>`
//
//	// 定义正则表达式
//	re := regexp.MustCompile(`<input[^>]*?name="(.*?)"[^>]*?value="(.*?)"[^>]*?>`)
//
//	// 使用正则表达式提取input标签中的name和value
//	matches := re.FindAllStringSubmatch(htmlContent, -1)
//
//	// 构建URL查询参数
//	formData := make(map[string]string)
//	for _, match := range matches {
//		name := match[1]
//		value := match[2]
//		formData[name] = value
//	}
//
//	// 构建URL查询参数
//	var queryParams []string
//	for key, value := range formData {
//		queryParams = append(queryParams, fmt.Sprintf("%s=%s", key, url.QueryEscape(value)))
//	}
//
//	// 将查询参数拼接成URL
//	urlString := "https://mapi.alipay.com/gateway.do?" + strings.Join(queryParams, "&")
//	//urlString2 := "https://mapi.alipay.com/gateway.do?" + url.QueryEscape(strings.Join(queryParams, "&"))
//	fmt.Println(urlString)
//	//fmt.Println(urlString2)
//}
