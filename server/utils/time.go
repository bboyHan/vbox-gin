package utils

import "time"

func GenerateTimePoints(startTime, endTime time.Time, interval time.Duration) []string {
	var timePoints []string

	// 初始化当前时间为 startTime
	currentTime := startTime

	// 循环生成时间点，直到当前时间超过 endTime
	for currentTime.Before(endTime) {
		timePoints = append(timePoints, currentTime.Format("2006-01-02 15:04:05"))

		// 增加 interval 到当前时间
		currentTime = currentTime.Add(interval)
	}

	return timePoints
}
