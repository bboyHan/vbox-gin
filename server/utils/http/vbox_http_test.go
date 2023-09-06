package http

import (
	"testing"
)

func TestHttp(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		url     string
		wantErr bool
	}{
		{
			name:    "GET TX test",
			args:    args{"446794914"},
			url:     "https://api.unipay.qq.com/v1/r/1450000186/trade_record_query?CmdCode=query2&SubCmdCode=default&PageNum=1&BeginUnixTime=1659803532&EndUnixTime=1691339532&PageSize=10&SystemType=portal&pf=2199&pfkey=pfkey&from_h5=1&session_token=63F728D4-74CB-4817-9F5D-3C344573837F1691339532798&webversion=MidasTradeRecord1.0&r=0.10077481030292357&openid=446794914&openkey=openkey&session_id=hy_gameid&session_type=st_dummy&__refer=&encrypt_msg=ab00dc01d7748d2ea42b2f24971b6c52ba4ecee8b4b741031ffea3e0775f5e06edb08110ebba54a8dcc93fc9a7ff0a4bee0eb4f6ad2033d3c3b2a90e5d9547d1aa96750a759652b9fe44dbcb0dce4d19&msg_len=76",
			wantErr: false,
		},
		{
			name:    "GET 22 test",
			args:    args{"446794914"},
			url:     "https://api.unipay.qq.com/v1/r/1450000186/trade_record_query?CmdCode=query2&SubCmdCode=default&PageNum=1&BeginUnixTime=1659803532&EndUnixTime=1691339532&PageSize=10&SystemType=portal&pf=2199&pfkey=pfkey&from_h5=1&session_token=63F728D4-74CB-4817-9F5D-3C344573837F1691339532798&webversion=MidasTradeRecord1.0&r=0.10077481030292357&openid=446794914&openkey=openkey&session_id=hy_gameid&session_type=st_dummy&__refer=&encrypt_msg=ab00dc01d7748d2ea42b2f24971b6c52ba4ecee8b4b741031ffea3e0775f5e06edb08110ebba54a8dcc93fc9a7ff0a4bee0eb4f6ad2033d3c3b2a90e5d9547d1aa96750a759652b9fe44dbcb0dce4d19&msg_len=76",
			wantErr: false,
		},
	}
	client := NewHTTPClient()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// 创建 HTTP 客户端实例
			headers := map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer token",
			}
			options := &RequestOptions{
				Headers:      headers,
				MaxRedirects: 3,
			}
			got, err := client.Get(tt.url, options)
			if err != nil {
				t.Errorf("GET() error = %v", err)
				return
			}

			// 检查响应体不是空的
			if len(got.Body) == 0 {
				t.Errorf("empty response body")
			}

			t.Logf("GET() got = %v", string(got.Body))
		})
	}
}

func TestNewProxyHTTPClientHttp(t *testing.T) {
	type args struct {
		d string
	}
	tests := []struct {
		name    string
		args    args
		url     string
		wantErr bool
	}{
		{
			name:    "GET TX test",
			args:    args{"446794914"},
			url:     "https://api.unipay.qq.com/v1/r/1450000186/trade_record_query?CmdCode=query2&SubCmdCode=default&PageNum=1&BeginUnixTime=1659803532&EndUnixTime=1691339532&PageSize=10&SystemType=portal&pf=2199&pfkey=pfkey&from_h5=1&session_token=63F728D4-74CB-4817-9F5D-3C344573837F1691339532798&webversion=MidasTradeRecord1.0&r=0.10077481030292357&openid=446794914&openkey=openkey&session_id=hy_gameid&session_type=st_dummy&__refer=&encrypt_msg=ab00dc01d7748d2ea42b2f24971b6c52ba4ecee8b4b741031ffea3e0775f5e06edb08110ebba54a8dcc93fc9a7ff0a4bee0eb4f6ad2033d3c3b2a90e5d9547d1aa96750a759652b9fe44dbcb0dce4d19&msg_len=76",
			wantErr: false,
		},
	}

	client := NewProxyHTTPClient()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			// 创建 HTTP 客户端实例
			headers := map[string]string{
				"Content-Type":  "application/json",
				"Authorization": "Bearer token",
			}
			options := &RequestOptions{
				Headers:      headers,
				MaxRedirects: 3,
			}
			got, err := client.Get(tt.url, options)
			if err != nil {
				t.Errorf("GET() error = %v", err)
				return
			}

			// 检查响应体不是空的
			if len(got.Body) == 0 {
				t.Errorf("empty response body")
			}

			t.Logf("GET() got = %v", string(got.Body))
		})
	}
}
