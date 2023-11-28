package utils

import (
	"fmt"
	"math"
	"time"
)

func FormatTimeToStr(times time.Time, pattern string) string {
	var d string
	if pattern != "" {
		d = times.Format(pattern)
	} else {
		d = times.Format("2006-01-02 15:04:05")
	}
	return d
}

func Get5MinNearlyOneHour() []string {
	// 获取当前时间
	now := time.Now()

	// 计算当前时间的分钟数
	currentMinute := now.Minute()

	// 计算最新的时间点
	latestMinute := math.Floor(float64(currentMinute)/5) * 5

	// 创建一个空数组来存储时间点
	var timePoints []string

	// 循环生成时间点
	for i := 0; i < 12; i++ { // 12 表示一个小时内的时间点数量
		// 计算每个时间点的小时和分钟数
		hour := now.Hour()
		minute := int(latestMinute) - (i * 5)

		// 处理小时和分钟的边界情况
		if minute < 0 {
			hour -= 1
			minute += 60
		}

		// 格式化小时和分钟，确保是两位数
		hourStr := fmt.Sprintf("%02d", hour)
		minuteStr := fmt.Sprintf("%02d", minute)

		// 拼接时间点字符串并添加到数组中
		timePoint := hourStr + ":" + minuteStr
		timePoints = append([]string{timePoint}, timePoints...) // 使用 append() 方法将时间点添加到数组的开头
	}

	// 输出时间点数组
	return timePoints
}

func main() {
	t := time.Now()
	fmt.Println(FormatTimeToStr(t, "yyyy-MM-dd hh:mm:ss")) // 示例调用
	fmt.Println(Get5MinNearlyOneHour())                    // 示例调用
}
