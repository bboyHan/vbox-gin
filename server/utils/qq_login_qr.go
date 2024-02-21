package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//func main() {
//
//	Data, err := LoopIfLogin()
//
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	ResCookie, err := GetResult(Data)
//
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	fmt.Println(ResCookie)
//
//}

// GetResult /* 获取所需参数 */
func GetResult(Data map[string]string) (map[string]string, error) {

	timestamp := strconv.Itoa(int(time.Now().Unix() * 1000))
	PostData := "response_type=token&client_id=101502376&redirect_uri=https://pay.qq.com/h5/shop.shtml&scope=all&state=&switch=&from_ptlogin=1&src=1&update_auth=1&openapi=1010&g_tk=" + Data["G_Token"] + "&auth_time=" + timestamp

	data := url.Values{}
	for _, value := range strings.Split(PostData, "&") {
		parts := strings.Split(value, "=")
		if len(parts) == 2 {
			data.Set(parts[0], parts[1])
		}
	}

	req, err := http.NewRequest("POST", "https://graph.qq.com/oauth2.0/authorize", strings.NewReader(data.Encode()))
	if err != nil {
		fmt.Println("创建请求时发生错误:", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", Data["Cookie"])

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 禁止自动跳转
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("发送请求时发生错误:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 检查是否有重定向地址
	redirectURL := resp.Header.Get("Location")

	ResData := make(map[string]string)
	ResData["OpenKey"] = ExtractParamValue(redirectURL, "access_token")
	ResData["OpenID"] = GetOpenId(ResData["OpenKey"])
	return ResData, nil

}

// GetOpenId /* 获取OpenId */
func GetOpenId(OpenKey string) string {

	url := "https://graph.qq.com/oauth2.0/me?access_token=" + OpenKey
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("请求失败:", err)
		return ""
	}

	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("获取OpenId结果时发生错误:", err)
		return ""
	}

	re := regexp.MustCompile(`"openid":"([^"]+)"`)
	// 查找匹配的子串
	match := re.FindStringSubmatch(string(body))
	if len(match) == 2 {
		openid := match[1]
		return openid
	} else {
		fmt.Println("OpenID为空：" + string(body))
	}

	return ""
}

// Credential /* 循环检查用户是否扫描成功以及是否登录成功 */
func LoopIfLogin() (map[string]string, error) {

StartLoop:
	LoginSig, err := GetLoginSig()
	if err != nil {
		return nil, err
	}

	QrData, err := GetQrImg()
	if err != nil {
		fmt.Println("请求失败:", err)
		return nil, err
	}

	QrToken := GetQrToken(QrData["QrSig"])

	var (
		isFirstLoop bool
		res         = make(map[string]string)
	)

OuterLoop:
	for {
		str, err := IfLogin(QrToken, LoginSig, QrData["QrSig"])
		if err != nil {
			return nil, err
		}

		if !strings.Contains(str, "") {
			return nil, errors.New("未知错误 Line 70，请刷新重试！")
		}

		// 间隔3秒循环一次
		if isFirstLoop {
			time.Sleep(time.Second * 3)
		}

		s := strings.Split(strings.ReplaceAll(str[strings.Index(str, "(")+1:len(str)-1], "'", ""), ",")
		// 65 二维码已失效 66 二维码未失效 67 已扫描,但还未点击确认 0  已经点击确认,并登录成功
		switch s[0] {
		case "65":
			fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "二维码失效 -> 已重新生成")
			goto StartLoop

		case "66":
			fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "登录二维码获取成功 -> 等待扫描")
			isFirstLoop = true
			continue OuterLoop

		case "67":
			fmt.Println(time.Now().Format("2006/01/02 15:04:05"), "已扫描 -> 请点击允许登录")
			isFirstLoop = true
			continue OuterLoop

		case "0":

			// 已经点击确认,并登录成功
			res["NickName"] = s[5]
			res["Location"] = s[2]
			res["QQ"] = ExtractParamValue(res["Location"], "uin")
			fmt.Println("QQ号:"+res["QQ"], "QQ昵称:"+res["NickName"])
			break OuterLoop

		default:
			return nil, errors.New("未知错误 Line 104，请刷新重试！")
		}
	}

	Res, err := Credential(res["Location"])
	if err != nil {
		return nil, err
	}
	return Res, nil
}

// Credential /* 登录成功，获取必要Cookie */
func Credential(url string) (map[string]string, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var (
		PsKey  string
		needs  = []string{"p_uin", "pt4_token", "p_skey", "pt_oauth_token"} // 需要从set-cookie取的参数
		cookie = make([]string, 0)
	)

	setCookies := resp.Header.Values("Set-Cookie")
	for _, val := range setCookies {
		c := strings.Split(strings.Split(val, ";")[0], "=")
		name := c[0]
		value := c[1]
		for _, TempKey := range needs {
			if name == TempKey && value != "" {
				if TempKey == "p_skey" {
					PsKey = value
				}
				cookie = append(cookie, fmt.Sprintf("%s=%s", name, value))
			}
		}
	}

	Data := make(map[string]string)
	Data["G_Token"] = GetToken(PsKey)
	Data["Cookie"] = strings.Join(cookie, "; ")
	return Data, nil
}

// IfLogin /* 检查用户是否扫描成功以及是否登录成功 */
func IfLogin(QrToken string, LoginSig string, QrSig string) (string, error) {

	timestamp := fmt.Sprintf("0-0-%d", time.Now().Unix()*1000)
	url := fmt.Sprintf("https://ssl.ptlogin2.qq.com/ptqrlogin?u1=%s&ptqrtoken=%v&ptredirect=0&h=1&t=1&g=1&from_ui=1&ptlang=2052&action=%v&js_ver=23111510&js_type=1&login_sig=%v&pt_uistyle=40&aid=716027609&daid=383&pt_3rd_aid=101502376&&o1vId=&pt_js_version=v1.48.1", url.QueryEscape("https://graph.qq.com/oauth2.0/login_jump"), QrToken, timestamp, LoginSig)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("检查登录状态创建请求时发生错误:", err)
		return "", err
	}

	req.Header.Set("cookie", fmt.Sprintf("qrsig=%s;", QrSig))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("检查登录状态发生错误:", err)
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取登录状态检查结果发生错误:", err)
		return "", err
	}

	return string(body), nil

}

// GetQrImg /* 获取登录二维码 */
func GetQrImg() (map[string]string, error) {

	t := strconv.FormatFloat(rand.Float64(), 'g', -1, 64)
	url := "https://ssl.ptlogin2.qq.com/ptqrshow?appid=716027609&e=2&l=M&s=3&d=72&v=4&t=" + t + "&daid=383&pt_3rd_aid=101502376&u1=https%3A%2F%2Fgraph.qq.com%2Foauth2.0%2Flogin_jump"
	resp, err := http.Get(url)

	if err != nil {
		fmt.Println("获取登录二维码请求失败:", err)
		return nil, err
	}

	defer resp.Body.Close()

	// 读取二维码
	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("读取二维码失败", err)
		return nil, err
	}

	// 将登录二维码内容转换为base64编码
	imageBase64 := "data:image/png;base64," + base64.StdEncoding.EncodeToString(imageData)

	ResMap := make(map[string]string)
	ResMap["QrSig"] = strings.Replace(strings.Split(resp.Header.Get("Set-Cookie"), ";")[0], "qrsig=", "", 1)
	ResMap["QrImg"] = imageBase64
	/* base64转换成图片 */

	//fmt.Println("登录二维码 -> Base64编码:  " + ResMap["QrImg"])
	return ResMap, nil
}

// GetLoginSig /* 获取LoginSig参数 */
func GetLoginSig() (string, error) {

	url := "https://xui.ptlogin2.qq.com/cgi-bin/xlogin?appid=716027609&daid=383&style=33&login_text=%E7%99%BB%E5%BD%95&hide_title_bar=1&hide_border=1&target=self&s_url=https%3A%2F%2Fgraph.qq.com%2Foauth2.0%2Flogin_jump&pt_3rd_aid=101502376&pt_feedback_link=https%3A%2F%2Fsupport.qq.com%2Fproducts%2F77942%3FcustomInfo%3D.appid101502376&theme=2&verify_theme="
	resp, err := http.Get(url)
	if err != nil {
		return "", errors.New(err.Error())
	}
	resp.Body.Close()

	setCookies := resp.Header.Values("Set-Cookie")
	if len(setCookies) < 1 {
		return "", errors.New("获取LoginSig参数错误，请稍后重试")
	}

	var loginSig string
	for _, val := range setCookies {
		if strings.Contains(val, "pt_login_sig=") {
			s := strings.Split(val, ";")
			for _, v := range s {
				if strings.Contains(v, "pt_login_sig=") {
					loginSig = strings.Replace(v, "pt_login_sig=", "", 1)
				}
			}
		}
	}

	if loginSig == "" {
		return "", errors.New("获取LoginSig参数错误，请稍后重试")
	}

	return loginSig, nil
}

// GetToken /* 计算GToken */
func GetToken(Pskey string) string {

	hash := 5381
	for i := 0; i < len(Pskey); i++ {
		char := int(Pskey[i])
		hash = ((hash << 5) + hash) + char
	}
	return strconv.Itoa(hash & 0x7fffffff)

}

// GetQrToken /* 计算QrToken */
func GetQrToken(QrSig string) string {

	e := 0
	for i := 0; i < len(QrSig); i++ {
		e += (e << 5) + int(QrSig[i])
	}
	return strconv.Itoa(2147483647 & e)

}

// ExtractParamValue /* 取出中间字符串 */
func ExtractParamValue(url string, param string) string {

	startIndex := strings.Index(url, param+"=")
	if startIndex == -1 {
		return ""
	}

	startIndex += len(param) + 1
	endIndex := strings.Index(url[startIndex:], "&")
	if endIndex == -1 {
		return url[startIndex:]
	}

	return url[startIndex : startIndex+endIndex]

}
