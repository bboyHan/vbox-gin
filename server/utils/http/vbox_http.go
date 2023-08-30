package http

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/core"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"net"
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

// httpClient 结构体
type HttpClient struct {
	client *fasthttp.Client
}

// NewHTTPClient 创建一个新的 httpClient 实例
func NewHTTPClient() *HttpClient {
	return &HttpClient{
		client: &fasthttp.Client{
			//Addr:                ProxyDialer("113.142.58.204:51022"),
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
func (c *HttpClient) Get(url string, options *RequestOptions) (*Response, error) {
	return c.sendRequest("GET", url, options)
}

// Post 发送 POST 请求
func (c *HttpClient) Post(url string, options *RequestOptions) (*Response, error) {
	return c.sendRequest("POST", url, options)
}

// Put 发送 PUT 请求
func (c *HttpClient) Put(url string, options *RequestOptions) (*Response, error) {
	return c.sendRequest("PUT", url, options)
}

// Delete 发送 DELETE 请求
func (c *HttpClient) Delete(url string, options *RequestOptions) (*Response, error) {
	return c.sendRequest("DELETE", url, options)
}

// sendRequest 发送 HTTP 请求
func (c *HttpClient) sendRequest(method, url string, options *RequestOptions) (*Response, error) {
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
		global.GVA_LOG = core.Zap()
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
