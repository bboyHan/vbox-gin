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
	re := regexp.MustCompile(`((https?://)[^\s]+)`)
	url = re.FindString(content)

	// 输出匹配结果
	global.GVA_LOG.Info("解析url结果", zap.Any("url", url))
	if url != "" {
		return url, nil
	}
	return url, err

}

// ValidAlipayUrl 查找Alipay url合法性
func ValidAlipayUrl(requestString string) bool {
	if strings.Contains(requestString, "alipays://platformapi") {
	} else if strings.Contains(requestString, "https://ur.alipay.com") {
	} else {
		return false
	}
	return true
}

// ValidJDUrl 查找JD url合法性
func ValidJDUrl(requestString string) bool {
	if strings.Contains(requestString, "openapp.jdmobile://") {
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
	global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))
	if strings.Contains(requestString, "alipays://platformapi") {
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
	global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))

	if strings.Contains(requestString, "openapp.jdmobile://") {
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

// HandleDYUrl 处理dy url
func HandleDYUrl(requestString string) (payUrl string, err error) {
	global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))

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

// HandleTBUrl 处理tb url
func HandleTBUrl(requestString string) (payUrl string, err error) {
	global.GVA_LOG.Info("处理前链接", zap.Any("payUrl", requestString))

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
