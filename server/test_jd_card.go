package main

//import (
//	"encoding/json"
//	"fmt"
//	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
//	"io/ioutil"
//	"net/http"
//	"strings"
//	"time"
//)
//
//func ECardQuery(Code string, Cookie string) error {
//	httpRequest := http.DefaultClient
//
//	// 查询卡密
//	var Temp map[string]interface{}
//
//	for i := 0; i <= 1; i++ {
//		JdCode, err := JdCode()
//		if err != nil {
//			return err
//		}
//
//		header := map[string][]string{
//			"Content-type": {"application/x-www-form-urlencoded"},
//			"User-agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48"},
//			"Referer":      {"https://ma.taobao.com/consume/code.htm?spm=0.0.0.0.tCxxtm"},
//			"Cookie":       {Cookie},
//		}
//
//		// 判断Cookie属于Pc或app
//		UrlPath := "app"
//		if strings.Contains(Cookie, "thor=") {
//			UrlPath = "pc"
//		}
//
//		Pdata := fmt.Sprintf("giftCardPwd=%s&verifyCode=%s&sessionId=%s&doBindFlag=0&queryGiftCardType=0", Code, JdCode["verifyCode"], JdCode["sessionId"])
//
//		request, err := http.NewRequest("POST", "https://mygiftcard.jd.com/giftcard/queryBindGiftCardCom/"+UrlPath, strings.NewReader(Pdata))
//		if err != nil {
//			return err
//		}
//		request.Header = header
//
//		response, err := httpRequest.Do(request)
//		if err != nil {
//			return err
//		}
//		defer response.Body.Close()
//
//		if response.StatusCode == 200 {
//			body, err := ioutil.ReadAll(response.Body)
//			if err != nil {
//				return err
//			}
//			err = json.Unmarshal(body, &Temp)
//			if err != nil {
//				return err
//			}
//
//			fmt.Printf("查询成功, %v", Temp)
//		} else if response.StatusCode == 302 {
//			// 删除查单Cookie
//			// Do something
//		}
//	}
//
//	if len(Temp) == 0 {
//		return fmt.Errorf("网络异常,请稍后再试")
//	}
//
//	if data, ok := Temp["data"].(map[string]interface{}); ok {
//		if amount, ok := data["amount"].(float64); !ok || amount == 0 {
//			if msg, ok := Temp["msg"].(string); ok {
//				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> %s 时间:%s", Code, msg, time.Now().Format("01-02 15:04:05"))
//				return fmt.Errorf(ErrNote)
//			} else {
//				return fmt.Errorf("网络异常,稍后再试")
//			}
//		}
//		if cardTypeName, ok := data["cardTypeName"].(string); !ok || cardTypeName != "京东E卡" {
//			ErrNote := fmt.Sprintf("卡密:%s  类型 -> %s 时间:%s", Code, cardTypeName, time.Now().Format("01-02 15:04:05"))
//			return fmt.Errorf("请提交京东E卡卡密, %v", ErrNote)
//		}
//	}
//
//	return nil
//}
//
//func ECardBind(Code string, Cookie string) (map[string]interface{}, error) {
//	httpRequest := http.DefaultClient
//
//	// 判断Cookie属于Pc或app
//	UrlPath := "app"
//	if strings.Contains(Cookie, "thor=") {
//		UrlPath = "pc"
//	}
//
//	var Temp map[string]interface{}
//
//	for i := 0; i <= 1; i++ {
//		JdCode, err := JdCode()
//		if err != nil {
//			return nil, err
//		}
//
//		header := map[string][]string{
//			"Content-type": {"application/x-www-form-urlencoded"},
//			"User-agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48"},
//			"Referer":      {"https://ma.taobao.com/consume/code.htm?spm=0.0.0.0.tCxxtm"},
//			"Cookie":       {Cookie},
//		}
//
//		Pdata := fmt.Sprintf("giftCardPwd=%s&verifyCode=%s&sessionId=%s&doBindFlag=0&queryGiftCardType=0", Code, JdCode["verifyCode"], JdCode["sessionId"])
//
//		request, err := http.NewRequest("POST", "https://mygiftcard.jd.com/giftcard/queryBindGiftCardCom/"+UrlPath, strings.NewReader(Pdata))
//		if err != nil {
//			return nil, err
//		}
//		request.Header = header
//
//		response, err := httpRequest.Do(request)
//		if err != nil {
//			return nil, err
//		}
//		defer response.Body.Close()
//
//		if response.StatusCode == 200 {
//			body, err := ioutil.ReadAll(response.Body)
//			if err != nil {
//				return nil, err
//			}
//			err = json.Unmarshal(body, &Temp)
//			if err != nil {
//				return nil, err
//			}
//
//			if msg, ok := Temp["msg"].(string); ok && strings.Contains(msg, "验证码") {
//				continue
//			} else {
//				break
//			}
//		} else if response.StatusCode == 302 {
//			// 卡密失效
//			ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> Cookie失效 时间:%s", Code, time.Now().Format("01-02 15:04:05"))
//			return nil, fmt.Errorf(ErrNote)
//		}
//	}
//
//	if len(Temp) == 0 {
//		return nil, fmt.Errorf("网络异常,请稍后再试")
//	}
//
//	if data, ok := Temp["data"].(map[string]interface{}); ok {
//		if amount, ok := data["amount"].(float64); !ok || amount == 0 {
//			if msg, ok := Temp["msg"].(string); ok {
//				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> %s 时间:%s", Code, msg, time.Now().Format("01-02 15:04:05"))
//				return nil, fmt.Errorf(ErrNote)
//			} else {
//				return nil, fmt.Errorf("网络异常,稍后再试")
//			}
//		}
//		pwdKey := data["pwdKey"].(string)
//		Pdata := fmt.Sprintf("pwdKey=%s&giftCardPwd=%s&doBindFlag=1&queryGiftCardType=0", pwdKey, Code)
//		request, err := http.NewRequest("POST", "https://mygiftcard.jd.com/giftcard/queryBindGiftCardCom/"+UrlPath, strings.NewReader(Pdata))
//		if err != nil {
//			return nil, err
//		}
//		request.Header = map[string][]string{
//			"Cookie":       {Cookie},
//			"Content-type": {"application/x-www-form-urlencoded"},
//			"User-agent":   {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36 Edg/112.0.1722.48"},
//			"Referer":      {"https://ma.taobao.com/consume/code.htm?spm=0.0.0.0.tCxxtm"}}
//		response, err := httpRequest.Do(request)
//		if err != nil {
//			return nil, err
//		}
//		defer response.Body.Close()
//
//		if response.StatusCode != 200 {
//			if response.StatusCode == 302 {
//				ErrNote := fmt.Sprintf("卡密:%s  错误信息 -> Cookie失效 时间:%s", Code, time.Now().Format("01-02 15:04:05"))
//				return nil, fmt.Errorf(ErrNote)
//			}
//			return nil, fmt.Errorf("网络异常,请稍后再试")
//		}
//		body, err := ioutil.ReadAll(response.Body)
//		if err != nil {
//			return nil, err
//		}
//		err = json.Unmarshal(body, &Temp)
//		if err != nil {
//			return nil, err
//		}
//		if _, ok := Temp["bindNextUrl"]; ok {
//			return nil, fmt.Errorf("卡密兑换异常")
//		}
//		if data, ok := Temp["data"].(map[string]interface{}); ok {
//			if _, ok := data["amount"].(float64); !ok {
//				return nil, fmt.Errorf("卡密兑换网络异常")
//			}
//			if _, ok := data["giftCardId"].(string); !ok {
//				return nil, fmt.Errorf("卡密兑换网络异常1")
//			}
//			return map[string]interface{}{
//				"pwdKey":     pwdKey,
//				"CardId":     data["giftCardId"].(string),
//				"OrderMoney": data["amount"].(float64),
//			}, nil
//		}
//	}
//
//	return nil, nil
//}
//
//func JdCode() (map[string]interface{}, error) {
//	// 可用验证码服务器
//	url := "http://1.12.50.148:8887/jd/slide"
//	response, err := http.Get(url)
//	if err != nil {
//		return nil, err
//	}
//	defer response.Body.Close()
//
//	body, err := ioutil.ReadAll(response.Body)
//	if err != nil {
//		return nil, err
//	}
//
//	var Temp map[string]interface{}
//	err = json.Unmarshal(body, &Temp)
//	if err != nil {
//		return nil, err
//	}
//
//	if _, ok := Temp["verifyCode"]; !ok {
//		return nil, fmt.Errorf("验证码获取失败")
//	}
//	if _, ok := Temp["sessionId"]; !ok {
//		return nil, fmt.Errorf("验证码获取失败")
//	}
//
//	return Temp, nil
//}
//
//func main() {
//	//ck := "__jda=122270672.17079125173481453510891.1707912517.1707912517.1707912517.1; __jdb=122270672.1.17079125173481453510891|1.1707912517; __jdv=122270672%7Cdirect%7C-%7Cnone%7C-%7C1707912517365; __jdc=122270672; mba_muid=17079125173481453510891; shshshfpb=BApXeTE2MpOhA8z0HYKS5t4hGSpkr4RWcB8dhFRlX9xJ1MipNk4O2; shshshfpa=43059317-92e3-5703-3a8f-23ee9027a446-1685201236; shshshfpx=43059317-92e3-5703-3a8f-23ee9027a446-1685201236; guid=24835ddfad01ef4556cf05cf4a2aeffba32acd789a9decc0db124187064b5f55; lang=chs; lsid=31834414192075z4oyvazev6f5qcylebnzhmr8tua3ih81707912518604; lstoken=gh7gw3mb; mba_sid=17079125179117037738453322324.2; __jd_ref_cls=MLoginRegister_LoginSucExpo; jcap_dvzw_fp=RGrbYDpQKDXyePRyMMfxgTaW7BDu8Og01EtZQzzub2tt7gq31yY4kcXixia9-x281byBwGxKl-5f8amC4E9BDg==; TrackerID=pNcyN-4arWSEgEgmIRVq-kIu86iWHup_orm1xLYRzZouOoRtjzua1HmkhxJcQC8ZOQ4sEvrwn9XB-2dPpg_CT2ZykKTLjT0CIF2pbTw9Oqc; pt_key=AAJlzK1eADBrXSqqGJvIg56cByMESDckyrIFnu-EzqeRUoKxahMv_TOn4jmpSUi_m7R9AdCs1sc; pt_pin=jd_NSVBmYiyiCtV; pt_token=ie42572b; pwdt_id=jd_NSVBmYiyiCtV; s_key=AAJlzK1eADCcR0ZEg6ygVzbQm7pQK7PZ2sgrhrLaxsIcPjm6yvognSqMJ4ueHzj2lFbgx3MNnZc; s_pin=jd_NSVBmYiyiCtV; sfstoken=tk01m9ca61c10a8sMisyeDNpZ29q3/KNCVBcOPZH3In/ifULxWphk/R1X5zvIr16bOh3N8bQnNKgrP9dzo56cOokz7ie; whwswswws=;"
//	//ck := "__jdu=1679706653460551248286; shshshfpa=0808d58d-22ea-558a-b032-881272e4f06c-1679706655; shshshfpx=0808d58d-22ea-558a-b032-881272e4f06c-1679706655; pinId=jGMQX5feeOE; shshshfp=40f1858c1a0126c2b9acaf93773a607f; mba_muid=1679706653460551248286; webp=1; visitkey=6283653559125451873; __wga=1705806356760.1705806356760.1705724619720.1705724619720.1.2; __jdv=95931165|direct|-|none|-|1707893209697; TrackID=1IZZVBqzzs7Ltm-iVMl5E5l0mYJZD1Hemzs0Lw34yFjsbX2lKMWWjdbhCwJJ03afar6HAnHGjru3g1mMsxTxDl_RcEanySCx3r5lEQ1ikCV4; thor=FB6800089A1E5822E48B6F2FA2E91C2EF4B2BDE3F5682EA8ADFAA8DF1ABA21E85157CDD7DF7E3FD190169EB884E78054CD7CC5D35621C9FAFAEF41E811080843F15BD0BB6F6BA629CA76CDE0817E08B38D81502D3393715AF0D07BAA29182D3B597D7214F28A064D3ADD494A362E0B387BEC128076DE43658DA9ECDDBBC962AB38E9F9E3D243FED0A68E791FE7E03DB7; flash=2_pvNB3rq3odAg6OMBJt0_GPo_ohYv3jVlpuAhZ19OAyxbdafEt6rq70WRhP2981fcsl2V0Dy3mPxwlElmaypRsBj05q0dLwntzesxTuVWFDL*; pin=han0373; unick=han0373; _tp=egJi1suo0lkSpN%2BtzDQJ%2Bg%3D%3D; _pst=han0373; areaId=1; ipLoc-djd=1-2901-0-0; source=PC; platform=pc; mba_sid=17084393308571457793655091157.1; __jd_ref_cls=LoginDisposition_Go; x-rp-evtoken=N-nAb5Oj6OS1u8hkvixIgIH7Kgael4VpIq2AGqEIUsmxtY0FPhYhkZzGFK3G2V_9qV_Uu6RSQvCxBWbrigJotm_c0zub0EIe_oD5BUDr2STES3AjoufH8INICURB2uHT6Xc_vf5wjxMaA5RGxhneNMq6oGvlnRIBYFYFKqygWWuZqNPzCxnoNLcN0kpfb9lYeeunRQM8FXuybgrWDMuuJIyz8hN7AayGuF7asUgdKM8%3D; lpkLoginType=3; joyya=1708439347.1708439387.31.1xys0p3; 3AB9D23F7A4B3CSS=jdd03M7T7ECP5B4S5IGEKZ4YOUPMXYLODHELJ2K3ZZFJ3ZSZEB4S3Y2FIAXFFIRFVM7TBFOMGY3KCUR4BKJX3Z7UVIRIBP4AAAAMNY3WNPJIAAAAAD22D57GMFFPSDUX; user-key=6d77891c-2a37-4c90-90a0-c71fefdf499f; cn=29; shshshsID=b3956aaed01fa0ad349edecad66bc16d_5_1708439441923; shshshfpb=BApXe-5DlxehA_nud0K27PgL74WXJIkG2BzsxYXp79xJ1MrOxsYO2; qid_seq=1; qid_uid=c1ffec2a-b884-4b5c-a859-a41bc64e4f80; qid_fs=1708439521777; qid_ls=1708439521777; qid_ts=1708439521785; qid_vis=1; qid_sid=c1ffec2a-b884-4b5c-a859-a41bc64e4f80-1; _distM=290154380108; __jda=261478171.1679706653460551248286.1679706653.1708157249.1708277724.56; __jdb=261478171.13.1679706653460551248286|56.1708277724; __jdc=261478171; 3AB9D23F7A4B3C9B=M7T7ECP5B4S5IGEKZ4YOUPMXYLODHELJ2K3ZZFJ3ZSZEB4S3Y2FIAXFFIRFVM7TBFOMGY3KCUR4BKJX3Z7UVIRIBP4"
//	ck := "a=1;thor=2"
//	//ECardBind("E71D-5E47-33ED-50BA", ck)
//	//ECardQuery("515F-2CB9-34CD-DFE1", ck)
//	//code, _ := JdCode()
//	//fmt.Println(code)
//
//	cookie, err := product.JDValidCookie(ck)
//	if err != nil {
//		fmt.Println(err)
//	} else {
//		fmt.Println(cookie)
//	}
//}
