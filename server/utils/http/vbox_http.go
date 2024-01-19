package http

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/gin-gonic/gin"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RequestOptions struct {
	Headers      map[string]string // 请求头
	Proxy        string            // 代理地址
	MaxRedirects int               // 重定向次数
	QueryParams  map[string]string // URL参数
	PayloadType  string            // 数据传参方式，可选值为 "form"、"json"、"url"
	Payload      interface{}       // 数据传参内容
}

type Response struct {
	StatusCode int               // 响应状态码
	Headers    map[string]string // 响应头
	Body       []byte            // 响应体
}

// FastHttpClient 结构体
type FastHttpClient struct {
	client *fasthttp.Client
}

const (
	Default = 0
	Redis   = 1
	DB      = 2
)

func IsValidCookie(cookieString string) bool {
	// 解析Cookie字符串
	request := &http.Request{Header: http.Header{"Cookie": {cookieString}}}
	cookies := request.Cookies()

	// 验证解析后的 Cookie 是否合法
	for _, cookie := range cookies {
		if cookie.Name == "" || cookie.Value == "" {
			return false
		}
	}

	return true
}

func ParseCookie(cookieStr string, targetKey string) string {
	fg := IsValidCookie(cookieStr)
	if !fg {
		global.GVA_LOG.Warn("cookie不合法", zap.String("cookie", cookieStr))
		return ""
	}
	pairs := strings.Split(cookieStr, ";")
	var flag bool
	var valueX string
	var valueY string
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		lowerKey := strings.ToLower(strings.TrimSpace(kv[0]))
		lowerTargetKey := strings.ToLower(targetKey)
		if len(kv) == 2 && strings.Contains(lowerKey, lowerTargetKey) {
			if lowerKey == lowerTargetKey {
				flag = true
				valueX = strings.TrimSpace(kv[1])
				break
			}
			valueY = strings.TrimSpace(kv[1])
		}
		if len(kv) == 1 && lowerKey == lowerTargetKey {
			return ""
		}
	}

	if flag {
		return valueX
	} else {
		return valueY
	}

}

func ProxyAddress2DB() string {
	var proxyDB vbox.Proxy
	err := global.GVA_DB.Where("status = ? and chan = ?", 1, "proxy").First(&proxyDB).Error
	if err != nil {
		log.Fatal("Proxy URL from DB parsing error:", err)
	}
	log.Printf("xxxxxxx: %v", proxyDB)

	c := NewHTTPClient()
	// 创建 HTTP 客户端实例
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer token",
	}
	options := &RequestOptions{
		Headers:      headers,
		MaxRedirects: 3,
	}
	res, err := c.Get(proxyDB.Url, options)
	s := string(res.Body)
	log.Printf("pppppppp: %v", strings.TrimSpace(s))

	return s
}

func NewProxyHTTPClient() *FastHttpClient {

	//1. cache
	//2. db
	// 设置数据库连接参数
	//dsn := "vbox_admin:Vbox123qwe@tcp(rm-cn-pe33bubix0001wko.rwlb.rds.aliyuncs.com:3306)/vbox_gin?charset=utf8mb4&parseTime=True&loc=Local"
	//// 连接数据库
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	//var proxyDB map[string]interface{}
	//err = db.Where("status = ?", 1).First(&proxyDB).Error

	s := ProxyAddress2DB()

	return NewHTTPClient(s)
}

// NewHTTPClient 创建一个新的 httpClient 实例
func NewHTTPClient(proxyAddr ...string) *FastHttpClient {
	var dialFunc fasthttp.DialFunc

	switch len(proxyAddr) {
	case 0: // 不使用代理
		dialFunc = nil
	case 1: // 使用指定的代理地址，格式为 "IP:Port"
		proxyURL, err := url.Parse("http://" + proxyAddr[0])
		if err != nil {
			log.Fatal("Proxy URL parsing error:", err)
		}

		// 创建一个定制的Dial函数，用于设置代理
		dialFunc = func(addr string) (net.Conn, error) {
			proxyConn, err := net.Dial("tcp", proxyURL.Host)
			if err != nil {
				return nil, fmt.Errorf("proxy connection error: %v", err)
			}

			// 连接成功后发送Connect请求，告知代理要连接的目标地址
			proxyReaderWriter := bufio.NewReadWriter(bufio.NewReader(proxyConn), bufio.NewWriter(proxyConn))
			fmt.Fprintf(proxyReaderWriter, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, addr)
			if err := proxyReaderWriter.Flush(); err != nil {
				return nil, fmt.Errorf("proxy write error: %v", err)
			}

			// 读取响应，确认代理是否连接成功
			resp, err := http.ReadResponse(proxyReaderWriter.Reader, &http.Request{Method: "CONNECT"})
			if err != nil {
				return nil, fmt.Errorf("proxy read response error: %v", err)
			}
			resp.Body.Close()

			// 返回与目标服务器的连接
			return proxyConn, nil
		}
	case 2: // 传入IP和Port，用于支持 http.NewHTTPClient("1.1.1.1","80") 的方式
		proxyIP := proxyAddr[0]
		proxyPort := proxyAddr[1]

		proxyURL, err := url.Parse("http://" + proxyIP + ":" + proxyPort)
		if err != nil {
			log.Fatal("Proxy URL parsing error:", err)
		}

		// 创建一个定制的Dial函数，用于设置代理
		dialFunc = func(addr string) (net.Conn, error) {
			proxyConn, err := net.Dial("tcp", proxyURL.Host)
			if err != nil {
				return nil, fmt.Errorf("proxy connection error: %v", err)
			}

			// 连接成功后发送Connect请求，告知代理要连接的目标地址
			proxyReaderWriter := bufio.NewReadWriter(bufio.NewReader(proxyConn), bufio.NewWriter(proxyConn))
			fmt.Fprintf(proxyReaderWriter, "CONNECT %s HTTP/1.1\r\nHost: %s\r\n\r\n", addr, addr)
			if err := proxyReaderWriter.Flush(); err != nil {
				return nil, fmt.Errorf("proxy write error: %v", err)
			}

			// 读取响应，确认代理是否连接成功
			resp, err := http.ReadResponse(proxyReaderWriter.Reader, &http.Request{Method: "CONNECT"})
			if err != nil {
				return nil, fmt.Errorf("proxy read response error: %v", err)
			}
			resp.Body.Close()

			// 返回与目标服务器的连接
			return proxyConn, nil
		}
	default:
		log.Fatal("Invalid proxy address")
	}
	return &FastHttpClient{
		client: &fasthttp.Client{
			Dial:                dialFunc,
			MaxConnsPerHost:     100,
			ReadBufferSize:      4096,
			WriteBufferSize:     4096,
			MaxIdleConnDuration: 10 * time.Second,
		},
	}
}

// ProxyDialer 设置代理
func ProxyDialer(proxyURL string) fasthttp.DialFunc {
	return func(addr string) (netConn net.Conn, err error) {
		if len(proxyURL) > 0 {
			proxyConn, err := fasthttp.Dial(proxyURL)
			if err != nil {
				return nil, err
			}
			return proxyConn, nil
		}

		return fasthttp.Dial(addr)
	}
}

// Get 发送 GET 请求
func (c *FastHttpClient) Get(url string, options *RequestOptions) (*Response, error) {
	return c.sendRequest("GET", url, options)
}

// Post 发送 POST 请求
func (c *FastHttpClient) Post(url string, options *RequestOptions) (*Response, error) {
	return c.sendRequest("POST", url, options)
}

// Put 发送 PUT 请求
func (c *FastHttpClient) Put(url string, options *RequestOptions) (*Response, error) {
	return c.sendRequest("PUT", url, options)
}

// Delete 发送 DELETE 请求
func (c *FastHttpClient) Delete(url string, options *RequestOptions) (*Response, error) {
	return c.sendRequest("DELETE", url, options)
}

// sendRequest 发送 HTTP 请求
func (c *FastHttpClient) sendRequest(method, url string, options *RequestOptions) (*Response, error) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod(method)

	// 设置请求头
	if options != nil && len(options.Headers) > 0 {
		for key, value := range options.Headers {
			req.Header.Set(key, value)
		}
	}

	// 设置代理
	if options != nil && options.Proxy != "" {
		global.GVA_LOG.Info("使用代理: ->", zap.Any("addr", options.Proxy))
		req.SetHost(options.Proxy)
	}

	// 设置URL参数
	if options != nil && len(options.QueryParams) > 0 {
		queryArgs := req.URI().QueryArgs()
		for key, value := range options.QueryParams {
			queryArgs.Add(key, value)
		}
	}

	// 设置数据传参
	if options != nil && options.PayloadType != "" && options.Payload != nil {
		switch options.PayloadType {
		case "form":
			formArgs := req.PostArgs()
			if payloadMap, ok := options.Payload.(map[string]string); ok {
				for key, value := range payloadMap {
					formArgs.Add(key, value)
				}
			}
		case "json":
			if payloadBytes, err := json.Marshal(options.Payload); err == nil {
				req.Header.SetContentType("application/json")
				req.SetBody(payloadBytes)
			} else {
				return nil, fmt.Errorf("failed to marshal JSON payload: %v", err)
			}
		case "url":
			if payloadStr, ok := options.Payload.(string); ok {
				req.Header.SetContentType("application/x-www-form-urlencoded")
				req.SetBodyString(payloadStr)
			}
		default:
			return nil, fmt.Errorf("unsupported payload type: %s", options.PayloadType)
		}
	}

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	var err error
	// 设置重定向次数
	if options != nil && options.MaxRedirects > 0 {
		err = c.client.DoRedirects(req, resp, options.MaxRedirects)
	} else {
		err = c.client.Do(req, resp)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %v", err)
	}

	body := resp.Body()

	response := &Response{
		StatusCode: resp.StatusCode(),
		Headers:    make(map[string]string),
		Body:       body,
	}

	// 获取响应头
	resp.Header.VisitAll(func(key, value []byte) {
		response.Headers[string(key)] = string(value)
	})

	return response, nil
}

func DoGinContextBody(c *gin.Context) []byte {
	var body []byte
	if c.Request.Method != http.MethodGet {
		var err error
		body, err = io.ReadAll(c.Request.Body)
		if err != nil {
			global.GVA_LOG.Error("read body from request error:", zap.Error(err))
		} else {
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
	} else {
		query := c.Request.URL.RawQuery
		query, _ = url.QueryUnescape(query)
		split := strings.Split(query, "&")
		m := make(map[string]string)
		for _, v := range split {
			kv := strings.Split(v, "=")
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			}
		}
		body, _ = json.Marshal(&m)
	}

	return body
}
