package utils

import (
	"bufio"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"go.uber.org/zap"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strings"
)

// IsAlphaNumericUnderscore // 匹配规则：只包含英文字母、数字和下划线
func IsAlphaNumericUnderscore(str string) bool {
	// 匹配规则：只包含英文字母、数字和下划线
	pattern := "^[a-zA-Z0-9_]+$"

	match, _ := regexp.MatchString(pattern, str)
	return match
}

// IsNumeric 判断字符串是否为数值
func IsNumeric(s string) bool {
	// 编译正则表达式
	// 这个正则表达式匹配整数、浮点数以及负数
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(s)
}

// ParseUrlContent 识别url内容，如果含有多个，则返回第一个
func ParseUrlContent(content string) (url string, err error) {
	// 1.
	// 使用正则表达式匹配 URL
	//re := regexp.MustCompile(`((https?://)[^\s]+)`)
	//urls := re.FindAllString(content, -1)
	//
	//// 输出匹配结果
	//for _, url = range urls {
	//	fmt.Println(url)
	//}

	//	2.处理查找第一个
	// 使用正则表达式匹配 URL
	//re := regexp.MustCompile(`((https?://)[^\s]+)`)
	re := regexp.MustCompile(`(\S+://[^\s]+)`)
	url = re.FindString(content)

	// 输出匹配结果
	global.GVA_LOG.Info("解析url结果", zap.Any("url", url))
	if url != "" {
		return url, nil
	} else {
		return "", fmt.Errorf("解析url失败")
	}
}

// ValidAlipayUrl 查找Alipay url合法性
func ValidAlipayUrl(requestString string) bool {
	if strings.Contains(requestString, "alipays://platformapi") {
	} else if strings.Contains(requestString, "jiaoyimao") {
	} else if strings.Contains(requestString, "https://ur.alipay.com") {
	} else {
		return false
	}
	return true
}

// ValidJDUrl 查找JD url合法性
func ValidJDUrl(requestString string) bool {
	if strings.Contains(requestString, "openapp.jdmobile://") {
	} else if strings.Contains(requestString, "3.cn") {
	} else if strings.Contains(requestString, "m.jd.com") {
	} else if strings.Contains(requestString, "item.jd.com") {
	} else {
		return false
	}
	return true
}

// ValidDYUrl 查找JD url合法性
func ValidDYUrl(requestString string) bool {
	if strings.Contains(requestString, "snssdk1128://") {
	} else if strings.Contains(requestString, "v.douyin.com") {
	} else {
		return false
	}
	return true
}

// HandleAlipayUrl 处理Alipay url
func HandleAlipayUrl(requestString string) (payUrl string, err error) {
	//global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))
	if strings.Contains(requestString, "alipays://platformapi") || strings.Contains(requestString, "jiaoyimao") {
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", payUrl))
	} else if strings.Contains(requestString, "ur.alipay.com") {
		payUrl = "alipays://platformapi/startapp?appId=20000067&url=" + url.QueryEscape(requestString)
		global.GVA_LOG.Info("处理后链接", zap.Any("payUrl", payUrl))
	} else {
		return "", fmt.Errorf("不合法的ZFB链接")
	}
	return payUrl, nil
}

// HandleJDUrl 处理jd url
func HandleJDUrl(requestString string) (payUrl string, err error) {
	//global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))

	if strings.Contains(requestString, "openapp.jdmobile://") {
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", payUrl))
	} else if strings.Contains(requestString, "3.cn") {
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", payUrl))
	} else if strings.Contains(requestString, "item.jd.com") {
		re := regexp.MustCompile(`item.jd.com/(\d+)\.html`)
		match := re.FindStringSubmatch(requestString)
		if match != nil && len(match) > 1 {
			skuId := match[1]
			schemaBody := "{\"sourceValue\":\"0_productDetail_97\",\"des\":\"productDetail\",\"skuId\":\"" + skuId + "\",\"category\":\"jump\",\"sourceType\":\"PCUBE_CHANNEL\"}"
			payUrl = "openapp.jdmobile://virtual?params=" + url.QueryEscape(schemaBody)
			global.GVA_LOG.Info("处理后链接", zap.Any("payUrl", payUrl))
		} else {
			return "", fmt.Errorf("不合法的JD链接")
		}
	} else if strings.Contains(requestString, "m.jd.com") {
		re := regexp.MustCompile(`/product/(\d+)\.html`)
		match := re.FindStringSubmatch(requestString)
		if match != nil && len(match) > 1 {
			skuId := match[1]
			schemaBody := "{\"sourceValue\":\"0_productDetail_97\",\"des\":\"productDetail\",\"skuId\":\"" + skuId + "\",\"category\":\"jump\",\"sourceType\":\"PCUBE_CHANNEL\"}"
			payUrl = "openapp.jdmobile://virtual?params=" + url.QueryEscape(schemaBody)
			global.GVA_LOG.Info("处理后链接", zap.Any("payUrl", payUrl))
		} else {
			return "", fmt.Errorf("不合法的JD链接")
		}
	} else {
		return "", fmt.Errorf("不合法的JD链接")
	}
	return payUrl, nil
}

// ValidTBUrl 校验tb url合法性
func ValidTBUrl(requestString string) bool {
	if strings.Contains(requestString, "m.tb.cn") {
	} else if strings.Contains(requestString, "tbopen://") {
	} else if strings.Contains(requestString, "m.taobao.com") {
	} else if strings.Contains(requestString, "item.taobao.com") {
	} else {
		return false
	}
	return true
}

// ValidXCXUrl 校验微信小程序合法性
func ValidXCXUrl(requestString string) bool {
	if strings.Contains(requestString, "#小程序://") {
	} else {
		return false
	}
	return true
}

// ValidPddUrl 校验PDD合法性
func ValidPddUrl(requestString string) bool {
	if strings.Contains(requestString, "pingduoduo://") {
	} else if strings.Contains(requestString, "pinduoduo:") {
	} else if strings.Contains(requestString, "mobile.yangkeduo.com") {
	} else {
		return false
	}
	return true
}

// HandleDYUrl 处理dy url
func HandleDYUrl(requestString string) (payUrl string, err error) {
	//global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))

	if strings.Contains(requestString, "snssdk1128://") {
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", payUrl))
	} else if strings.Contains(requestString, "v.douyin.com") {
		client := vbHttp.NewHTTPClient()
		var options = &vbHttp.RequestOptions{
			MaxRedirects: 0,
		}
		resp, errQ := client.Get(requestString, options)
		if errQ != nil {
			global.GVA_LOG.Error("err:  ->", zap.Error(errQ))
			return "", fmt.Errorf("不合法的链接")
		}
		respHeaders := resp.Headers
		loc := respHeaders["Location"]
		parsedURL, errX := url.Parse(loc)
		if errX != nil {
			global.GVA_LOG.Error("err:  ->", zap.Error(errX))
			return "", fmt.Errorf("不合法的链接")
		}
		query := parsedURL.Query()
		detailSchema := query.Get("detail_schema")
		payUrl = strings.ReplaceAll(detailSchema, "sslocal://", "snssdk1128://")
		global.GVA_LOG.Info("处理后链接", zap.Any("payUrl", payUrl))
	} else {
		return "", fmt.Errorf("不合法的DY链接")
	}

	return payUrl, nil
}

// HandleXCXUrl 处理wx xcx
func HandleXCXUrl(requestString string) (payUrl string, err error) {
	//global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))
	if strings.Contains(requestString, "#小程序://") {
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", requestString))
	} else {
		return "", fmt.Errorf("不合法的XCX链接")
	}
	return payUrl, nil
}

// HandlePddUrl 处理pdd url
func HandlePddUrl(requestString string) (payUrl string, err error) {
	//global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))
	if strings.Contains(requestString, "pingduoduo://") {
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", requestString))
	} else if strings.Contains(requestString, "mobile.yangkeduo.com") {
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", requestString))
	} else if strings.Contains(requestString, "pinduoduo:") {
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", requestString))
	} else {
		return "", fmt.Errorf("不合法的PDD链接")
	}
	return payUrl, nil
}

func HandleTBUrl(requestString string) (payUrl string, err error) {
	//global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))

	if strings.Contains(requestString, "m.tb.cn") {
		// 先请求一次,获取 html body
		client := vbHttp.NewHTTPClient()
		var options = &vbHttp.RequestOptions{
			MaxRedirects: 3,
		}
		resp, errQ := client.Get(requestString, options)
		if errQ != nil {
			global.GVA_LOG.Error("err:  ->", zap.Error(errQ))
			return "", fmt.Errorf("不合法的链接")
		}
		htmlBody := string(resp.Body)
		//fmt.Printf("%s", htmlBody)

		// 先请求一次,获取 html body
		re := regexp.MustCompile(`var url = '([^']*)'`)
		match := re.FindStringSubmatch(htmlBody)
		if match != nil && len(match) > 1 {
			tmpUrl := match[1]
			global.GVA_LOG.Info("", zap.Any("tmpUrl", (tmpUrl)))

			parsedURL, errX := url.Parse(tmpUrl)
			if errX != nil {
				global.GVA_LOG.Warn("无效的 URL:", zap.Error(errX))
				return "", fmt.Errorf("不合法的链接")
			}

			query := parsedURL.Query()
			itemId := query.Get("id")
			global.GVA_LOG.Info("", zap.Any("itemID", itemId))
			if itemId == "" {
				// m.taobao.com/i(\d+)\.htm
				re = regexp.MustCompile(`m.taobao.com/i(\d+)\.htm`)
				match = re.FindStringSubmatch(htmlBody)
				if match != nil && len(match) > 1 {
					itemId = match[1]
					global.GVA_LOG.Info("二次修复查找", zap.Any("itemID", itemId))
				}
			}

			if itemId == "" {
				return "", fmt.Errorf("不合法的链接")
			}

			schema := "https://h5.m.taobao.com/awp/core/detail.htm?id=" + itemId

			payUrl = "tbopen://m.taobao.com/tbopen/index.html?h5Url=" + url.QueryEscape(schema)
			global.GVA_LOG.Info("处理后链接", zap.Any("payUrl", payUrl))
		}

	} else if strings.Contains(requestString, "item.taobao.com") {
		parsedURL, errX := url.Parse(requestString)
		if errX != nil {
			global.GVA_LOG.Warn("无效的 URL:", zap.Error(errX))
			return "", fmt.Errorf("不合法的链接")
		}

		query := parsedURL.Query()
		itemId := query.Get("id")
		global.GVA_LOG.Info("", zap.Any("itemID", itemId))
		if itemId == "" {
			return "", fmt.Errorf("不合法的链接")
		}

		schema := "https://h5.m.taobao.com/awp/core/detail.htm?id=" + itemId

		payUrl = "tbopen://m.taobao.com/tbopen/index.html?h5Url=" + url.QueryEscape(schema)
		global.GVA_LOG.Info("处理后链接", zap.Any("payUrl", payUrl))

	} else if strings.Contains(requestString, "tbopen://") {
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", payUrl))
	} else if strings.Contains(requestString, "main.m.taobao.com") {
		//payUrl = "tbopen://m.taobao.com/tbopen/index.html?h5Url=" + url.QueryEscape(requestString)
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", payUrl))
		//global.GVA_LOG.Info("处理后链接", zap.Any("payUrl", payUrl))
	} else if strings.Contains(requestString, "market.m.taobao.com") {
		//payUrl = "tbopen://m.taobao.com/tbopen/index.html?h5Url=" + url.QueryEscape(requestString)
		payUrl = requestString
		global.GVA_LOG.Info("无需处理", zap.Any("payUrl", payUrl))
		//global.GVA_LOG.Info("处理后链接", zap.Any("payUrl", payUrl))
	} else {
		return "", fmt.Errorf("不合法的TB链接")
	}
	return payUrl, nil
}

// ParseRequest 解析http 报文内容
func ParseRequest(requestString string) (*http.Request, error) {
	reader := bufio.NewReader(strings.NewReader(requestString))

	// 使用 http.ReadRequest 函数解析请求
	request, err := http.ReadRequest(reader)
	if err != nil {
		return nil, err
	}

	return request, nil
}

// FilterNotContains a 全集，b 子集，筛选出 a 中不在 b 中的元素
func FilterNotContains(a, b []string) []string {
	// 创建一个 map，用于存储 b 数组中的元素
	bMap := make(map[string]bool)
	for _, v := range b {
		bMap[v] = true
	}

	// 遍历 a 数组，将不在 b 中的元素添加到结果数组中
	var result []string
	for _, v := range a {
		if !bMap[v] {
			result = append(result, v)
		}
	}

	return result
}

func Contains(slice interface{}, target interface{}) bool {
	sliceValue := reflect.ValueOf(slice)
	itemValue := reflect.ValueOf(target)

	if sliceValue.Kind() != reflect.Slice {
		panic("slice parameter is not a slice")
	}

	for i := 0; i < sliceValue.Len(); i++ {
		if reflect.DeepEqual(sliceValue.Index(i).Interface(), itemValue.Interface()) {
			return true
		}
	}

	return false
}

func FindJsonValueByKey(jsonStr, key string) (string, error) {
	// 构建正则表达式，查找指定 key 的值
	regexPattern := fmt.Sprintf(`"%s":"(.*?)"`, key)
	keyRegex := regexp.MustCompile(regexPattern)

	// 在字符串中查找匹配项
	matches := keyRegex.FindStringSubmatch(jsonStr)

	// 提取匹配项的值
	var value string
	if len(matches) > 1 {
		value = matches[1]
	} else {
		return "", fmt.Errorf("未找到 key 为 %s 的值", key)
	}

	return value, nil
}

func FindUrlValueByKey(urlStr, key string) (string, error) {
	// 解析URL
	u, err := url.Parse(urlStr)
	if err != nil {
		fmt.Println("解析URL时出错:", err)
		return "", fmt.Errorf("不合法的支付链接")
	}
	// 获取查询参数
	queryParams := u.Query()

	// 获取contextId参数的值
	contextId := queryParams.Get(key)
	global.GVA_LOG.Info("contextId的值为:", zap.Any(key, contextId))
	return contextId, nil
}
