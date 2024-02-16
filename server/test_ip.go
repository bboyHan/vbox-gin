package main

//import (
//	"fmt"
//	"net"
//	"strconv"
//)
//
//func isIP(input string) bool {
//	// 尝试直接解析为IP地址
//	if net.ParseIP(input) != nil {
//		return true
//	}
//
//	// 尝试解析为IP地址加端口号
//	host, portStr, err := net.SplitHostPort(input)
//	if err != nil {
//		return false // 不能正确分割为host和port部分
//	}
//	if net.ParseIP(host) == nil {
//		return false // host部分不是合法的IP地址
//	}
//	port, err := strconv.Atoi(portStr)
//	if err != nil || port < 1 || port > 65535 {
//		return false // 端口号不是有效的范围
//	}
//
//	return true
//}
//
//func main() {
//	tests := []string{
//		"http://192.168.1.1",
//		"192.168.1.1:8080",
//		"256.256.256.256",   // 不合法的IP
//		"192.168.1.1:99999", // 不合法的端口
//	}
//
//	for _, test := range tests {
//		fmt.Printf("%s is valid: %t\n", test, isIP(test))
//	}
//}
