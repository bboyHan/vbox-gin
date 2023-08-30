package product

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/tx"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"net/url"
	"strconv"
	"time"
)

var rawURL = "https://api.unipay.qq.com/v1/r/1450000186/trade_record_query?" +
	"CmdCode=query2&SubCmdCode=default&PageNum=1&PageSize=200" +
	"&BeginUnixTime=1659803532&EndUnixTime=1691339532&SystemType=portal&pf=2199&pfkey=pfkey" +
	"&from_h5=1&session_token=63F728D4-74CB-4817-9F5D-3C344573837F1691339532798&webversion=MidasTradeRecord1.0&r=0.10077481030292357" +
	"&openid=446794914&openkey=openkey&session_id=hy_gameid&session_type=st_dummy&__refer=" +
	"&encrypt_msg=ab00dc01d7748d2ea42b2f24971b6c52ba4ecee8b4b741031ffea3e0775f5e06edb08110ebba54a8dcc93fc9a7ff0a4bee0eb4f6ad2033d3c3b2a90e5d9547d1aa96750a759652b9fe44dbcb0dce4d19&msg_len=76"

// 创建 HTTP 客户端实例
var headers = map[string]string{
	"Content-Type":  "application/json",
	"Authorization": "Bearer token",
}

var options = &http.RequestOptions{
	Headers:      headers,
	MaxRedirects: 3,
}

// Records 获取半小时内记录
func Records(qq string) *tx.Records {

	u, _ := url.Parse(rawURL)
	queryParams := u.Query()

	// 获取当前时间
	currentTime := time.Now()
	// 计算半小时前的时间
	halfHourAgo := currentTime.Add(-30 * time.Minute)

	// 当前时间秒数
	currentSeconds := currentTime.Unix()

	// 将半小时前的时间转换为秒数
	halfHourAgoSeconds := halfHourAgo.Unix()

	queryParams.Set("openid", qq)
	queryParams.Set("BeginUnixTime", strconv.FormatInt(halfHourAgoSeconds, 10))
	queryParams.Set("EndUnixTime", strconv.FormatInt(currentSeconds, 10))

	u.RawQuery = queryParams.Encode()
	newURL := u.String()
	client := http.NewHTTPClient()
	//options.Proxy = "43.248.99.24:61000"

	resp, _ := client.Get(newURL, options)
	fmt.Print(string(resp.Body))

	var records tx.Records
	json.Unmarshal(resp.Body, &records)
	//fmt.Print(records)

	return &records
}
