package http

import (
	"fmt"
	"net/url"
	"strconv"
	"testing"
	"time"
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
			name:    "GET 22 test",
			args:    args{"446794914"},
			url:     "https://api.unipay.qq.com/v1/r/1450000490/trade_record_query?pf=mds_storeopen_qb-__mds_default_v1_0_0.qb-html5&pfkey=pfkey&from_h5=1&session_token=90324BE6-E4D9-4DA9-AE57-44D3DE03C8F8&webversion=stdV2.16.0.1.other.other&r=0.4950957161982934&openid=69C952958A68C2BA80F5983C91C8662D&openkey=B14B6F35A990F78A17921EA32EDC8E11&session_id=openid&session_type=kp_accesstoken&qq_appid=&SerialNo=QQACCT_SAVE-20230618-DVM0DIDbwOOj&CmdCode=query2&SubCmdCode=default&SystemType=portal&BeginUnixTime=1690810697&EndUnixTime=1697069897&Order=desc&PageNum=1&PageSize=100&anti_auto_script_token_id=E86CDBACCB84586D35C73C0B5FD0869D0CC23B6014F0D5ED09E42B823F4338E8C5D606F78B48CC9B2444B720F45277DE26CFE054DAD6BC06EDF407CE52FCF1E3&__refer=https%3A%2F%2Fpay.qq.com%2Fh5%2Findex.shtml%3Fr%3D0.7360455474285279&encrypt_msg=ddcb93f583700dcf845ebb3a54dca62b5d623ce6bcafa3af58d41604917c8bb3a7cdc70bad9c406c009852ee9abc07b389da2b9f0041dc51d5655cc7679bc05739d4d4b4af72150ebdc63a1a4051c81931137759a5276911279136a0a141c6bde6982a6b383e1cb998661455244b20b775b270e1f8d9a6b0083b7895d1a4d267&msg_len=126",
			wantErr: false,
		},
	}
	client := NewHTTPClient()
	// 获取当前时间
	currentTime := time.Now()
	// 计算半小时前的时间
	halfHourAgo := currentTime.Add(-30 * time.Minute)
	// 将半小时前的时间转换为秒数
	halfHourAgoSeconds := halfHourAgo.Unix()

	// 当前时间秒数
	currentSeconds := currentTime.Unix()

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
			//rawUrl := fmt.Sprintf(tt.url, "69C952958A68C2BA80F5983C91C8662D", "B14B6F35A990F78A17921EA32EDC8E11", strconv.FormatInt(halfHourAgoSeconds, 10), strconv.FormatInt(currentSeconds, 10))
			u, _ := url.Parse(tt.url)
			queryParams := u.Query()
			queryParams.Set("openid", "69C952958A68C2BA80F5983C91C8662D")
			queryParams.Set("openkey", "B14B6F35A990F78A17921EA32EDC8E11")
			queryParams.Set("BeginUnixTime", strconv.FormatInt(halfHourAgoSeconds, 10))
			queryParams.Set("EndUnixTime", strconv.FormatInt(currentSeconds, 10))
			u.RawQuery = queryParams.Encode()
			fmt.Println(u.String())
			got, err := client.Get(u.String(), options)
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
			url:     "https://api.unipay.qq.com/v1/r/1450000490/trade_record_query?pf=mds_storeopen_qb-__mds_default_v1_0_0.qb-html5&pfkey=pfkey&from_h5=1&session_token=90324BE6-E4D9-4DA9-AE57-44D3DE03C8F8&webversion=stdV2.16.0.1.other.other&r=0.4950957161982934&openid=%s&openkey=%s&session_id=openid&session_type=kp_accesstoken&qq_appid=&SerialNo=QQACCT_SAVE-20230618-DVM0DIDbwOOj&CmdCode=query2&SubCmdCode=default&SystemType=portal&BeginUnixTime=%s&EndUnixTime=%s&Order=desc&PageNum=1&PageSize=300&anti_auto_script_token_id=E86CDBACCB84586D35C73C0B5FD0869D0CC23B6014F0D5ED09E42B823F4338E8C5D606F78B48CC9B2444B720F45277DE26CFE054DAD6BC06EDF407CE52FCF1E3&__refer=https%3A%2F%2Fpay.qq.com%2Fh5%2Findex.shtml%3Fr%3D0.7360455474285279&encrypt_msg=ddcb93f583700dcf845ebb3a54dca62b5d623ce6bcafa3af58d41604917c8bb3a7cdc70bad9c406c009852ee9abc07b389da2b9f0041dc51d5655cc7679bc05739d4d4b4af72150ebdc63a1a4051c81931137759a5276911279136a0a141c6bde6982a6b383e1cb998661455244b20b775b270e1f8d9a6b0083b7895d1a4d267&msg_len=126",
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
