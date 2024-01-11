package utils

import (
	"fmt"
	"reflect"
	"regexp"
)

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
