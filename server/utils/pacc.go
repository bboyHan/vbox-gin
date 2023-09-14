package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"math/rand"
	"reflect"
	"sort"
	"strings"
)

// VerifySign 校验签名合法性的通用方法
func VerifySign(data interface{}) bool {
	// 获取结构体字段及值
	v := reflect.ValueOf(data).Elem() // 获取结构体的反射对象
	t := v.Type()
	fields := make([]string, 0)
	values := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := fmt.Sprintf("%v", v.Field(i))
		//if field.Name != "Sign" && value != "" {
		fields = append(fields, field.Tag.Get("json"))
		values[field.Tag.Get("json")] = value
		//}
	}

	// 按照字段名称的ASCII码升序进行排序
	sort.Strings(fields)

	// 拼接排序后的键值对
	var signStr string
	for _, field := range fields {
		if field != "key" && field != "sign" {
			signStr += field + "=" + fmt.Sprintf("%v", values[field]) + "&"
		}
	}
	signStr += "key=" + fmt.Sprintf("%v", values["key"])

	fmt.Printf("计算前: %v", signStr)
	// 计算MD5签名
	sign := SignMD5(signStr)

	fmt.Printf("计算sign: %v", sign)
	fmt.Printf("传入sign: %v", values["sign"])

	// 将计算得到的签名与请求参数中的签名进行比较
	return sign == fmt.Sprintf("%v", values["sign"])
}

// SignMD5 计算MD5签名
func SignMD5(data string) string {
	hash := md5.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}

func RandSize2DB(total, limit int) (int, int) {
	offset := 0

	if limit == 0 {
		limit = 20 // 如果未传入limit，默认设置为20
	}

	if total > limit {
		// 如果总数大于需要获取的数量，计算limit和offset的值
		offset = rand.Intn(total - limit + 1) // 随机生成起始位置
	} else {
		limit = total // 如果总数不足20个，则按实际数量取出
	}

	return limit, offset
}

func GetDeviceSimpleInfo(userAgent string) string {
	var os string
	// 操作系统检查
	if strings.Contains(userAgent, "Windows") {
		os = "Windows"
	} else if strings.Contains(userAgent, "Android") {
		os = "Android"
	} else if strings.Contains(userAgent, "iPhone") || strings.Contains(userAgent, "iPad") || strings.Contains(userAgent, "iPod") {
		os = "iOS"
	} else if strings.Contains(userAgent, "BlackBerry") {
		os = "BlackBerry"
	} else if strings.Contains(userAgent, "hpwOS") {
		os = "WebOS"
	} else if strings.Contains(userAgent, "SymbianOS") {
		os = "Symbian"
	} else {
		os = "Other"
	}
	return os
}

func ToLowerCamelCase(s string) string {
	// 将字符串按照下划线分割为单词，并且将每个单词的首字母小写
	c := cases.Title(language.Und)

	words := strings.Split(s, "_")
	for i := range words {
		if i > 0 {
			words[i] = c.String(words[i])
		}
	}

	// 特殊处理第一个单词，将其首字母转换为小写
	if len(words) > 0 {
		words[0] = strings.ToLower(words[0])
	}
	return strings.Join(words, "")
}
