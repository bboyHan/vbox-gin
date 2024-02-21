package product

import (
	"encoding/json"
	"fmt"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func JDValidCookie(Cookie string) (bool, error) {
	client := vbHttp.NewHTTPClient()

	// 判断Cookie属于Pc或app
	UrlPath := "app"
	if strings.Contains(Cookie, "thor=") {
		UrlPath = "pc"
	}

	sessionId := "PWAXywABAAAATFuSvfkAMHlf-9L27tjijMqe2aOEW9LgxuKqnSHIyxyl-nuTSy8nHCj4V9Q3LHMYM8g11Xmk4QAAAAA"
	verifyCode := "VXv-5AABAAABjcw5xAYAgB8aLQGhPZIByfRKwrRcy-DjLvdXo7iCih06G_rsGbfpdokG5qZHlqZW05GTXRbUaBSbT9UWOz6HSaJ7q5q7azr4r9DU jYwaEVMnRQmEOsDiFS51jwQtc8_38xPxQ8xDkRy0zGHoSGCvWd2cVqOyeCod-wbqSchj_mzOcd-OBUe9"
	Pdata := fmt.Sprintf("giftCardPwd=%s&verifyCode=%s&sessionId=%s&doBindFlag=0&queryGiftCardType=0", "E71D-5E47-33ED-50BA", verifyCode, sessionId)

	// 创建 HTTP 客户端实例
	//"Content-type": {"application/x-www-form-urlencoded"},
	//			"User-agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48"},
	//			"Referer":      {"https://ma.taobao.com/consume/code.htm?spm=0.0.0.0.tCxxtm"},
	//			"Cookie":       {Cookie},
	var jdHeaders = map[string]string{
		"Content-type": "application/x-www-form-urlencoded",
		"User-agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48",
		"Referer":      "https://ma.taobao.com/consume/code.htm?spm=",
		"Cookie":       Cookie,
	}

	var jdOptions = &vbHttp.RequestOptions{
		Headers:      jdHeaders,
		MaxRedirects: 0,
		PayloadType:  "url",
		Payload:      Pdata,
	}

	resp, err := client.Post("https://mygiftcard.jd.com/giftcard/queryBindGiftCardCom/"+UrlPath, jdOptions)
	if err != nil {
		fmt.Println(err)
	}
	body := string(resp.Body)
	fmt.Printf("headers: %v \n", resp.Headers)
	fmt.Printf("statusCode: %v \n", resp.StatusCode)
	fmt.Printf("body: %v \n", body)
	if resp.StatusCode == 302 {
		return false, fmt.Errorf("ck过期, %v, resp: %v", resp.StatusCode, body)
	}
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("code error, %v, resp: %v", resp.StatusCode, body)
	}
	if resp.StatusCode == 200 {
		if strings.Contains(body, "NotLogin") {
			return false, fmt.Errorf("ck过期, %v, resp: %v", resp.StatusCode, body)
		}
	}
	return true, nil
}

func ECardQuery(Code string, Cookie string) error {
	httpRequest := http.DefaultClient

	// 查询卡密
	var Temp map[string]interface{}

	for i := 0; i <= 1; i++ {
		JdCode, err := JdCode()
		if err != nil {
			return err
		}

		header := map[string][]string{
			"Content-type": {"application/x-www-form-urlencoded"},
			"User-agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48"},
			"Referer":      {"https://ma.taobao.com/consume/code.htm?spm=0.0.0.0.tCxxtm"},
			"Cookie":       {Cookie},
		}

		// 判断Cookie属于Pc或app
		UrlPath := "app"
		if strings.Contains(Cookie, "thor=") {
			UrlPath = "pc"
		}

		Pdata := fmt.Sprintf("giftCardPwd=%s&verifyCode=%s&sessionId=%s&doBindFlag=0&queryGiftCardType=0", Code, JdCode["verifyCode"], JdCode["sessionId"])

		request, err := http.NewRequest("POST", "https://mygiftcard.jd.com/giftcard/queryBindGiftCardCom/"+UrlPath, strings.NewReader(Pdata))
		if err != nil {
			return err
		}
		request.Header = header

		response, err := httpRequest.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return err
			}
			err = json.Unmarshal(body, &Temp)
			if err != nil {
				return err
			}

			fmt.Printf("查询成功, %v", Temp)
		} else if response.StatusCode == 302 {
			// 删除查单Cookie
			// Do something
		}
	}

	if len(Temp) == 0 {
		return fmt.Errorf("网络异常,请稍后再试")
	}

	if data, ok := Temp["data"].(map[string]interface{}); ok {
		if amount, ok := data["amount"].(float64); !ok || amount == 0 {
			if msg, ok := Temp["msg"].(string); ok {
				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> %s 时间:%s", Code, msg, time.Now().Format("01-02 15:04:05"))
				return fmt.Errorf(ErrNote)
			} else {
				return fmt.Errorf("网络异常,稍后再试")
			}
		}
		if cardTypeName, ok := data["cardTypeName"].(string); !ok || cardTypeName != "京东E卡" {
			ErrNote := fmt.Sprintf("卡密:%s  类型 -> %s 时间:%s", Code, cardTypeName, time.Now().Format("01-02 15:04:05"))
			return fmt.Errorf("请提交京东E卡卡密, %v", ErrNote)
		}
	}

	return nil
}

func ECardBind(Code string, Cookie string) (map[string]interface{}, error) {
	httpRequest := http.DefaultClient

	// 判断Cookie属于Pc或app
	UrlPath := "app"
	if strings.Contains(Cookie, "thor=") {
		UrlPath = "pc"
	}

	var Temp map[string]interface{}

	for i := 0; i <= 1; i++ {
		JdCode, err := JdCode()
		if err != nil {
			return nil, err
		}

		header := map[string][]string{
			"Content-type": {"application/x-www-form-urlencoded"},
			"User-agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48"},
			"Referer":      {"https://ma.taobao.com/consume/code.htm?spm=0.0.0.0.tCxxtm"},
			"Cookie":       {Cookie},
		}

		Pdata := fmt.Sprintf("giftCardPwd=%s&verifyCode=%s&sessionId=%s&doBindFlag=0&queryGiftCardType=0", Code, JdCode["verifyCode"], JdCode["sessionId"])

		request, err := http.NewRequest("POST", "https://mygiftcard.jd.com/giftcard/queryBindGiftCardCom/"+UrlPath, strings.NewReader(Pdata))
		if err != nil {
			return nil, err
		}
		request.Header = header

		response, err := httpRequest.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode == 200 {
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(body, &Temp)
			if err != nil {
				return nil, err
			}

			if msg, ok := Temp["msg"].(string); ok && strings.Contains(msg, "验证码") {
				continue
			} else {
				break
			}
		} else if response.StatusCode == 302 {
			// 卡密失效
			ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> Cookie失效 时间:%s", Code, time.Now().Format("01-02 15:04:05"))
			return nil, fmt.Errorf(ErrNote)
		}
	}

	if len(Temp) == 0 {
		return nil, fmt.Errorf("网络异常,请稍后再试")
	}

	if data, ok := Temp["data"].(map[string]interface{}); ok {
		if amount, ok := data["amount"].(float64); !ok || amount == 0 {
			if msg, ok := Temp["msg"].(string); ok {
				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> %s 时间:%s", Code, msg, time.Now().Format("01-02 15:04:05"))
				return nil, fmt.Errorf(ErrNote)
			} else {
				return nil, fmt.Errorf("网络异常,稍后再试")
			}
		}
		pwdKey := data["pwdKey"].(string)
		Pdata := fmt.Sprintf("pwdKey=%s&giftCardPwd=%s&doBindFlag=1&queryGiftCardType=0", pwdKey, Code)
		request, err := http.NewRequest("POST", "https://mygiftcard.jd.com/giftcard/queryBindGiftCardCom/"+UrlPath, strings.NewReader(Pdata))
		if err != nil {
			return nil, err
		}
		request.Header = map[string][]string{
			"Cookie":       {Cookie},
			"Content-type": {"application/x-www-form-urlencoded"},
			"User-agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48"},
			"Referer":      {"https://ma.taobao.com/consume/code.htm?spm=0.0.0.0.tCxxtm"}}
		response, err := httpRequest.Do(request)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()

		if response.StatusCode != 200 {
			if response.StatusCode == 302 {
				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> Cookie失效 时间:%s", Code, time.Now().Format("01-02 15:04:05"))
				return nil, fmt.Errorf(ErrNote)
			}
			return nil, fmt.Errorf("网络异常,请稍后再试")
		}
		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &Temp)
		if err != nil {
			return nil, err
		}
		if _, ok := Temp["bindNextUrl"]; ok {
			return nil, fmt.Errorf("卡密兑换异常")
		}
		if data, ok := Temp["data"].(map[string]interface{}); ok {
			if _, ok := data["amount"].(float64); !ok {
				return nil, fmt.Errorf("卡密兑换网络异常")
			}
			if _, ok := data["giftCardId"].(string); !ok {
				return nil, fmt.Errorf("卡密兑换网络异常1")
			}
			return map[string]interface{}{
				"pwdKey":     pwdKey,
				"CardId":     data["giftCardId"].(string),
				"OrderMoney": data["amount"].(float64),
			}, nil
		}
	}

	return nil, nil
}

func JdCode() (map[string]interface{}, error) {
	/*可用验证码服务器 http://1.12.50.148:8887/jd/slide http://43.136.111.242:8887/jd/slide http://159.75.241.132:8887/jd/slide http://43.138.239.132:8887/jd/slide */

	url := "http://1.12.50.148:8887/jd/slide"
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var Temp map[string]interface{}
	err = json.Unmarshal(body, &Temp)
	if err != nil {
		return nil, err
	}

	if _, ok := Temp["verifyCode"]; !ok {
		return nil, fmt.Errorf("验证码获取失败")
	}
	if _, ok := Temp["sessionId"]; !ok {
		return nil, fmt.Errorf("验证码获取失败")
	}

	return Temp, nil
}
