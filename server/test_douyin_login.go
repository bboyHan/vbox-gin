package main

//import (
//	"fmt"
//	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
//)
//
//func main() {
//	url := "https://sso.douyin.com/get_qrcode/?need_logo=true"
//
//	client := vbHttp.NewHTTPClient()
//	headers := map[string]string{
//		"Content-Type": "application/json",
//	}
//	options := &vbHttp.RequestOptions{
//		MaxRedirects: 0,
//		Headers:      headers,
//	}
//
//	resp, _ := client.Get(url, options)
//	fmt.Println(resp.Headers)
//	fmt.Println(string(resp.Body))
//}
