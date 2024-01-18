package utils

import (
	"fmt"
	"reflect"
	"regexp"
)

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
