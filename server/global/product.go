package global

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// tx qb

func TxContains(target string) bool {
	set := map[string]bool{
		"1000": true,
		"1001": true, //qb jd
		"1002": true, //qb dy
		"1003": true, //qb jym
		"1004": true, //qb zfb
		"1005": true, //qb tb
		"1006": true, //qb wx xcx
		"1007": true, //qb mt
		"1008": true,
		"1009": true,

		"1101": true, // jun ka
	}
	_, found := set[target]
	return found
}

// tx dnf

func DnfContains(target string) bool {
	set := map[string]bool{
		"1200": true,
		"1201": true, //dnf jd
		"1202": true, //dnf tb
		"1203": true, //dnf jym
		"1204": true, //dnf zfb
		"1205": true, //dnf tb
		"1206": true, //dnf wx xcx
		"1207": true, //dnf mt
		"1208": true,
		"1209": true,
	}
	_, found := set[target]
	return found
}

// 剑三

func J3Contains(target string) bool {
	set := map[string]bool{
		"2000": true,
		"2001": true,
		"2002": true,
		"2003": true,
		"2004": true,
		"2005": true,
		"2006": true,
		"2007": true,
		"2008": true,
		"2009": true,
	}
	_, found := set[target]
	return found
}

// pay code

func PcContains(target string) bool {
	set := map[string]bool{
		"3000": true,
		"3001": true,
		"3002": true,
		"3003": true,
		"3004": true,
		"3005": true,
		"3006": true,
		"3007": true,
		"3008": true,
		"3009": true,
	}
	_, found := set[target]
	return found
}

func SdoContains(target string) bool {
	set := map[string]bool{
		"4000": true,
		"4001": true,
		"4002": true,
		"4003": true,
		"4004": true,
		"4005": true,
		"4006": true,
		"4007": true,
		"4008": true,
		"4009": true,
	}
	_, found := set[target]
	return found
}

func QNContains(target string) bool {
	set := map[string]bool{
		"5000": true,
		"5001": true,
		"5002": true,
		"5003": true,
		"5004": true,
		"5005": true,
		"5006": true,
		"5007": true,
		"5008": true,
		"5009": true,
	}
	_, found := set[target]
	return found
}

func ECContains(target string) bool {
	set := map[string]bool{
		"6000": true,
		"6001": true,
		"6002": true,
		"6003": true,
		"6004": true,
		"6005": true,
		"6006": true,
		"6007": true,
		"6008": true,
		"6009": true,
	}
	_, found := set[target]
	return found
}

func ISPContains(target string) bool {
	//写一个函数判断是否包含 yidong|liantong|dianxin,包含返回true，反之false
	if target == "yidong" || target == "liantong" || target == "dianxin" {
		return true
	}
	return false
}

func ISPTranslate(target string) string {
	//写一个函数判断是否包含 yidong|liantong|dianxin,包含返回true，反之false
	if strings.Contains(target, "移动") {
		return "yidong"
	} else if strings.Contains(target, "联通") {
		return "liantong"
	} else if strings.Contains(target, "电信") {
		return "dianxin"
	}
	return RandISP()
}

func ProvinceContains(target string) bool {
	if strings.Contains(target, "北京") || strings.Contains(target, "天津") || strings.Contains(target, "上海") || strings.Contains(target, "河北省") || strings.Contains(target, "山西省") || strings.Contains(target, "内蒙古自治区") || strings.Contains(target, "辽宁省") || strings.Contains(target, "吉林省") || strings.Contains(target, "黑龙江省") ||
		strings.Contains(target, "上海市") || strings.Contains(target, "江苏省") || strings.Contains(target, "浙江省") || strings.Contains(target, "安徽省") || strings.Contains(target, "福建省") || strings.Contains(target, "江西省") || strings.Contains(target, "山东省") || strings.Contains(target, "河南省") || strings.Contains(target, "湖北省") ||
		strings.Contains(target, "湖南省") || strings.Contains(target, "广东省") || strings.Contains(target, "广西壮族自治区") || strings.Contains(target, "海南省") || strings.Contains(target, "重庆市") || strings.Contains(target, "四川省") || strings.Contains(target, "贵州省") || strings.Contains(target, "云南省") || strings.Contains(target, "西藏自治区") ||
		strings.Contains(target, "陕西省") || strings.Contains(target, "甘肃省") || strings.Contains(target, "青海省") || strings.Contains(target, "宁夏回族自治区") || strings.Contains(target, "新疆维吾尔自治区") {
		return true
	}
	return false
}

func RandISP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 创建本地的随机生成器
	isps := []string{"yidong", "liantong", "dianxin"}    // 可选的ISP列表
	return isps[r.Intn(len(isps))]                       // 随机选择一个ISP并返回
}

func RandIP() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 创建本地的随机生成器
	ip := make([]byte, 4)                                // 创建字节切片存储IP地址
	for i := range ip {
		ip[i] = byte(r.Intn(255)) // 生成0~255之间的随机数作为IP地址的每个分段
	}
	return fmt.Sprintf("%d.%d.%d.%d", ip[0], ip[1], ip[2], ip[3]) // 格式化IP地址并返回
}

func RandProvince() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano())) // 创建本地的随机生成器
	provinces := []string{
		"北京市", "天津市", "河北省", "山西省", "内蒙古自治区", "辽宁省", "吉林省", "黑龙江省",
		"上海市", "江苏省", "浙江省", "安徽省", "福建省", "江西省", "山东省", "河南省", "湖北省",
		"湖南省", "广东省", "广西壮族自治区", "海南省", "重庆市", "四川省", "贵州省", "云南省",
		"西藏自治区", "陕西省", "甘肃省", "青海省", "宁夏回族自治区", "新疆维吾尔自治区",
	} // 可选的省份列表
	return provinces[r.Intn(len(provinces))] // 随机选择一个省份并返回
}
