package main

//
//import (
//	"bufio"
//	"fmt"
//	"net/http"
//	"strings"
//)
//
//func main() {
//	requestString := "POST /security_extend_server/helper/balance/queryBalance?gameCode=jx3&account=18210889498&accountType=&zoneCode=z22&SN=98710648156&remark=&sign=C0879271B5E8F47E3775B119FFD40024 HTTP/1.1\n\n"
//
//	request, err := parseRequest(requestString)
//	if err != nil {
//		fmt.Println("无效的请求:", err)
//		return
//	}
//
//	// 输出解析后的请求信息
//	fmt.Printf("Method: %s\n", request.Method)
//	fmt.Printf("URL: %s\n", request.URL)
//	fmt.Printf("Protocol: %s\n", request.Proto)
//	fmt.Println("Headers:")
//	for name, values := range request.Header {
//		for _, value := range values {
//			fmt.Printf("%s: %s\n", name, value)
//		}
//	}
//}
//
//func parseRequest(requestString string) (*http.Request, error) {
//	reader := bufio.NewReader(strings.NewReader(requestString))
//
//	// 使用 http.ReadRequest 函数解析请求
//	request, err := http.ReadRequest(reader)
//	if err != nil {
//		return nil, err
//	}
//
//	return request, nil
//}
