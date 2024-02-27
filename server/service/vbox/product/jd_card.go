package product

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
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
	// 创建 HTTP 客户端实例
	client := vbHttp.NewHTTPClient()
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

	// 判断Cookie属于Pc或app
	UrlPath := "app"
	if strings.Contains(Cookie, "thor=") {
		UrlPath = "pc"
	}
	URL := "https://mygiftcard.jd.com/giftcard/queryBindGiftCardCom/" + UrlPath

	var jsCodeMem string
	var Temp map[string]interface{}

	// 1. 先查一遍卡，不绑
	for i := 0; i < 3; i++ {
		zs, err := global.GVA_REDIS.ZRangeArgsWithScores(context.Background(), redis.ZRangeArgs{
			Key:   global.YdECJdCodeZSet,
			Start: 0,
			Stop:  -1,
		}).Result()
		if err != nil {
			global.GVA_LOG.Error("redis Ex", zap.Error(err))
			continue
		}

		if len(zs) <= 0 {
			global.GVA_LOG.Info("redis jd code 资源不足，重新拉取")

			JdCodeTemp, errJ := JdCode()
			if errJ != nil {
				global.GVA_LOG.Error("请求 jd code server 异常", zap.Error(errJ))
				continue
			}
			//JdCodeTemp["verifyCode"], JdCodeTemp["sessionId"]
			jsCodeMem = fmt.Sprintf("verifyCode=%s&sessionId=%s", JdCodeTemp["verifyCode"], JdCodeTemp["sessionId"])

			PData := fmt.Sprintf("giftCardPwd=%s&%s&doBindFlag=0&queryGiftCardType=0", Code, jsCodeMem)

			var jdOptions = &vbHttp.RequestOptions{
				Headers:      jdHeaders,
				MaxRedirects: 0,
				PayloadType:  "url",
				Payload:      PData,
			}
			resp, errR := client.Post(URL, jdOptions)
			if errR != nil {
				global.GVA_LOG.Error(fmt.Sprintf("请求 jd code server 异常, url : %s", URL), zap.Error(errR))
				continue
			}

			if resp.StatusCode == 200 {

				err = json.Unmarshal(resp.Body, &Temp)
				if err != nil {
					return err
				}

				if msg, ok := Temp["msg"].(string); ok && strings.Contains(msg, "验证码") {
					global.GVA_LOG.Error(fmt.Sprintf("请求 jd code server 异常, url : %s", URL), zap.Any("msg", msg))
					continue
				} else {
					break
				}
			} else if resp.StatusCode == 302 {
				// 卡密失效
				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> Cookie失效 时间:%s", Code, time.Now().Format("01-02 15:04:05"))
				return fmt.Errorf(ErrNote)
			}

		} else {
			z := zs[len(zs)-1] //取出最后一个，重新设置utc时间戳
			// verifyCode=xx&sessionId=xx
			jsCodeMem = z.Member.(string)

			PData := fmt.Sprintf("giftCardPwd=%s&%s&doBindFlag=0&queryGiftCardType=0", Code, jsCodeMem)
			//取出即删除，防止被重复使用
			global.GVA_REDIS.ZRem(context.Background(), global.YdECJdCodeZSet, jsCodeMem)

			var jdOptions = &vbHttp.RequestOptions{
				Headers:      jdHeaders,
				MaxRedirects: 0,
				PayloadType:  "url",
				Payload:      PData,
			}
			resp, errR := client.Post(URL, jdOptions)
			if errR != nil {
				global.GVA_LOG.Error(fmt.Sprintf("请求 jd code server 异常, url : %s", URL), zap.Error(errR))
				continue
			}

			if resp.StatusCode == 200 {

				err = json.Unmarshal(resp.Body, &Temp)
				if err != nil {
					return err
				}

				if msg, ok := Temp["msg"].(string); ok {
					if strings.Contains(msg, "被绑定") {
						global.GVA_LOG.Error(fmt.Sprintf("请求 jd card server 异常, url : %s", URL), zap.Any("msg", msg))
						return fmt.Errorf("该卡已被绑定！")
					}
					if strings.Contains(msg, "不存在") {
						global.GVA_LOG.Error(fmt.Sprintf("请求 jd card server 异常, url : %s", URL), zap.Any("msg", msg))
						return fmt.Errorf("该卡不存在，请重新输入！")
					}
				} else {
					break
				}
			} else if resp.StatusCode == 302 {
				// 卡密失效
				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> Cookie失效 时间:%s", Code, time.Now().Format("01-02 15:04:05"))
				return fmt.Errorf(ErrNote)
			}
		}

		if len(Temp) == 0 {
			return fmt.Errorf("网络异常,请稍后再试")
		}

		if data, ok := Temp["data"].(map[string]interface{}); ok {
			if amount, okA := data["amount"].(float64); !okA || amount == 0 {
				if msg, okM := Temp["msg"].(string); okM {
					ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> %s 时间:%s", Code, msg, time.Now().Format("01-02 15:04:05"))
					return fmt.Errorf(ErrNote)
				} else {
					return fmt.Errorf("网络异常,稍后再试")
				}
			}
			if cardTypeName, okE := data["cardTypeName"].(string); !okE || cardTypeName != "京东E卡" {
				ErrNote := fmt.Sprintf("卡密:%s  类型 -> %s 时间:%s", Code, cardTypeName, time.Now().Format("01-02 15:04:05"))
				return fmt.Errorf("请提交京东E卡卡密, %v", ErrNote)
			}
		}
	}
	global.GVA_LOG.Info("查验卡密通过", zap.Any("tmp info", Temp))

	return nil
}

func ECardBind(Code string, Cookie string) (map[string]interface{}, error) {

	// 创建 HTTP 客户端实例
	client := vbHttp.NewHTTPClient()
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

	// 判断Cookie属于Pc或app
	UrlPath := "app"
	if strings.Contains(Cookie, "thor=") {
		UrlPath = "pc"
	}
	URL := "https://mygiftcard.jd.com/giftcard/queryBindGiftCardCom/" + UrlPath

	var jsCodeMem string
	var Temp map[string]interface{}

	// 1. 先查一遍卡，不绑
	for i := 0; i < 3; i++ {
		zs, err := global.GVA_REDIS.ZRangeArgsWithScores(context.Background(), redis.ZRangeArgs{
			Key:   global.YdECJdCodeZSet,
			Start: 0,
			Stop:  -1,
		}).Result()
		if err != nil {
			global.GVA_LOG.Error("redis Ex", zap.Error(err))
			continue
		}

		if len(zs) <= 0 {
			global.GVA_LOG.Info("redis jd code 资源不足，重新拉取")

			JdCodeTemp, errJ := JdCode()
			if errJ != nil {
				global.GVA_LOG.Error("请求 jd code server 异常", zap.Error(errJ))
				continue
			}
			//JdCodeTemp["verifyCode"], JdCodeTemp["sessionId"]
			jsCodeMem = fmt.Sprintf("verifyCode=%s&sessionId=%s", JdCodeTemp["verifyCode"], JdCodeTemp["sessionId"])

			PData := fmt.Sprintf("giftCardPwd=%s&%s&doBindFlag=0&queryGiftCardType=0", Code, jsCodeMem)

			var jdOptions = &vbHttp.RequestOptions{
				Headers:      jdHeaders,
				MaxRedirects: 0,
				PayloadType:  "url",
				Payload:      PData,
			}
			resp, errR := client.Post(URL, jdOptions)
			if errR != nil {
				global.GVA_LOG.Error(fmt.Sprintf("请求 jd code server 异常, url : %s", URL), zap.Error(errR))
				continue
			}

			if resp.StatusCode == 200 {

				err = json.Unmarshal(resp.Body, &Temp)
				if err != nil {
					return nil, err
				}

				if msg, ok := Temp["msg"].(string); ok && strings.Contains(msg, "验证码") {
					global.GVA_LOG.Error(fmt.Sprintf("请求 jd code server 异常, url : %s", URL), zap.Any("msg", msg))
					continue
				} else {
					break
				}
			} else if resp.StatusCode == 302 {
				// 卡密失效
				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> Cookie失效 时间:%s", Code, time.Now().Format("01-02 15:04:05"))
				return nil, fmt.Errorf(ErrNote)
			}

		} else {
			z := zs[len(zs)-1] //取出最后一个，重新设置utc时间戳
			// verifyCode=xx&sessionId=xx
			jsCodeMem = z.Member.(string)

			PData := fmt.Sprintf("giftCardPwd=%s&%s&doBindFlag=0&queryGiftCardType=0", Code, jsCodeMem)
			//取出即删除，防止被重复使用
			global.GVA_REDIS.ZRem(context.Background(), global.YdECJdCodeZSet, jsCodeMem)

			var jdOptions = &vbHttp.RequestOptions{
				Headers:      jdHeaders,
				MaxRedirects: 0,
				PayloadType:  "url",
				Payload:      PData,
			}
			resp, errR := client.Post(URL, jdOptions)
			if errR != nil {
				global.GVA_LOG.Error(fmt.Sprintf("请求 jd code server 异常, url : %s", URL), zap.Error(errR))
				continue
			}

			if resp.StatusCode == 200 {

				err = json.Unmarshal(resp.Body, &Temp)
				if err != nil {
					return nil, err
				}

				if msg, ok := Temp["msg"].(string); ok && strings.Contains(msg, "验证码") {
					global.GVA_LOG.Error(fmt.Sprintf("请求 jd code server 异常, url : %s", URL), zap.Any("msg", msg))
					continue
				} else {
					break
				}
			} else if resp.StatusCode == 302 {
				// 卡密失效
				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> Cookie失效 时间:%s", Code, time.Now().Format("01-02 15:04:05"))
				return nil, fmt.Errorf(ErrNote)
			}
		}

		if len(Temp) == 0 {
			return nil, fmt.Errorf("网络异常,请稍后再试")
		}

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
		PData := fmt.Sprintf("pwdKey=%s&giftCardPwd=%s&doBindFlag=1&queryGiftCardType=0", pwdKey, Code)

		var jdBindHeaders = map[string]string{
			"Content-type": "application/x-www-form-urlencoded",
			"User-agent":   "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48",
			"Referer":      "https://ma.taobao.com/consume/code.htm?spm=0.0.0.0.tCxxtm",
			"Cookie":       Cookie,
		}

		var jdOptions = &vbHttp.RequestOptions{
			Headers:      jdBindHeaders,
			MaxRedirects: 0,
			PayloadType:  "url",
			Payload:      PData,
		}

		resp, errB := client.Post(URL, jdOptions)
		if errB != nil {
			global.GVA_LOG.Error(fmt.Sprintf("请求 jd code server 异常, url : %s", URL), zap.Error(errB))
			return nil, errB
		}
		if resp.StatusCode != 200 {
			if resp.StatusCode == 302 {
				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> Cookie失效 时间:%s", Code, time.Now().Format("01-02 15:04:05"))
				return nil, fmt.Errorf(ErrNote)
			}
			return nil, fmt.Errorf("网络异常,请稍后再试")
		}

		if errB != nil {
			return nil, errB
		}
		var bindTemp map[string]interface{}
		err := json.Unmarshal(resp.Body, &bindTemp)
		if err != nil {
			return nil, err
		}
		if dataB, okB := bindTemp["data"].(map[string]interface{}); okB {
			if _, okN := bindTemp["bindNextUrl"]; okN {
				return nil, fmt.Errorf("卡密兑换异常 not bindNextUrl")
			}
			if _, okA := dataB["amount"].(float64); !okA {
				return nil, fmt.Errorf("卡密兑换网络异常 not amount")
			}
			if _, okG := dataB["giftCardId"].(string); !okG {
				return nil, fmt.Errorf("卡密兑换网络异常 not giftCardId")
			}
			return map[string]interface{}{
				"pwdKey":     pwdKey,
				"CardId":     dataB["giftCardId"].(string),
				"OrderMoney": dataB["amount"].(float64),
			}, nil
		}
	}

	return nil, nil
}

func JdCode() (Temp map[string]interface{}, err error) {
	/*可用验证码服务器 http://1.12.50.148:8887/jd/slide http://43.136.111.242:8887/jd/slide http://159.75.241.132:8887/jd/slide http://43.138.239.132:8887/jd/slide */
	client := vbHttp.NewHTTPClient()

	url := "http://1.12.50.148:8887/jd/slide"
	resp, err := client.Get(url, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(resp.Body, &Temp)
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
