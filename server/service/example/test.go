package example

import (
	"fmt"
	"time"
)

func main() {
	// 获取当前日期
	now := time.Now()

	// 设置开始时间为 00:00:00
	startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 遍历一天中的每一个五分钟时间点
	for i := 0; i < 24*60/5; i++ {
		// 计算当前时间点
		currentTime := startTime.Add(time.Minute * 5 * time.Duration(i))

		// 格式化为 'HH:mm'
		formattedTime := currentTime.Format("15:04")

		// 输出
		fmt.Println(formattedTime)
	}
}
