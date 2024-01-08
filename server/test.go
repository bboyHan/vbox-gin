package main

/*
import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"sort"
	"strings"
	"time"
)

func hash33(t string) int64 {
	var e int64 = 0
	for i := 0; i < len(t); i++ {
		e += (e << 5) + int64(t[i])
	}
	return e & 2147483647
}

func ParseCookie(cookieStr string, targetKey string) string {
	pairs := strings.Split(cookieStr, ";")
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) == 2 && strings.Contains(strings.ToLower(strings.TrimSpace(kv[0])), strings.ToLower(targetKey)) {
			return strings.TrimSpace(kv[1])
		}
		if len(kv) == 1 && strings.ToLower(strings.TrimSpace(kv[0])) == strings.ToLower(targetKey) {
			return ""
		}
	}
	return ""
}


func main() {
	// 测试hash33函数
	// 模拟获取的数据
	dataJSON := `[
            {
                "ID": 23,
                "CreatedAt": "2023-12-18T19:05:26.526+08:00",
                "UpdatedAt": "2023-12-16T19:10:36.029+08:00",
                "orderId": "TZ170272473434567",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "20231216173651529",
                "eventType": 2,
                "acId": "88278125",
                "channelCode": "3000",
                "platformOid": "VB2023121619052617359794",
                "payIp": "106.39.151.103",
                "payRegion": "中国|0|北京|北京市|电信",
                "payDevice": "Windows",
                "resourceUrl": "weixin://wxpay/bizpayurl?pr=MxaBSkJzz",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 3,
                "cbStatus": 2,
                "expTime": "2023-12-16T19:10:33.378+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 22,
                "CreatedAt": "2023-12-18T19:01:59.082+08:00",
                "UpdatedAt": "2023-12-16T19:07:21.751+08:00",
                "orderId": "TZ170272452668393",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "20231216173647271",
                "eventType": 2,
                "acId": "88278125",
                "channelCode": "3000",
                "platformOid": "VB2023121619015917359785",
                "payIp": "106.39.151.103",
                "payRegion": "中国|0|北京|北京市|电信",
                "payDevice": "iOS",
                "resourceUrl": "weixin://wxpay/bizpayurl?pr=v9lVtQozz",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 3,
                "cbStatus": 2,
                "expTime": "2023-12-16T19:07:19.124+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 21,
                "CreatedAt": "2023-12-18T18:59:04.592+08:00",
                "UpdatedAt": "2023-12-16T19:04:14.978+08:00",
                "orderId": "TZ170272435267621",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "20231216101824014",
                "eventType": 2,
                "acId": "88278125",
                "channelCode": "3000",
                "platformOid": "VB2023121618590417359778",
                "payIp": "106.39.151.103",
                "payRegion": "中国|0|北京|北京市|电信",
                "payDevice": "Windows",
                "resourceUrl": "weixin://wxpay/bizpayurl?pr=MxaBSkJzz",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 3,
                "cbStatus": 2,
                "expTime": "2023-12-16T19:04:12.366+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 20,
                "CreatedAt": "2023-12-18T18:27:34.43+08:00",
                "UpdatedAt": "2023-12-16T18:32:46.893+08:00",
                "orderId": "TZ170272246250852",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "20231215230743386",
                "eventType": 2,
                "acId": "88278125",
                "channelCode": "3000",
                "platformOid": "VB2023121618273417359699",
                "payIp": "106.39.151.103",
                "payRegion": "中国|0|北京|北京市|电信",
                "payDevice": "Windows",
                "resourceUrl": "https://open.weixin.qq.com/connect/confirm?uuid=001RNyeY0eVKkl28",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 3,
                "cbStatus": 2,
                "expTime": "2023-12-16T18:32:44.427+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 19,
                "CreatedAt": "2023-12-18T18:21:32.291+08:00",
                "UpdatedAt": "2023-12-16T18:22:07.165+08:00",
                "orderId": "TZ170272210063071",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "",
                "eventType": 2,
                "acId": "",
                "channelCode": "3000",
                "platformOid": "VB2023121618213217359683",
                "payIp": "127.0.0.1",
                "payRegion": "0|0|0|内网IP|内网IP",
                "payDevice": "Windows",
                "resourceUrl": "",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 2,
                "cbStatus": 2,
                "expTime": "2023-12-16T18:26:32.265+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 18,
                "CreatedAt": "2023-12-18T17:39:32.752+08:00",
                "UpdatedAt": "2023-12-16T17:40:58.231+08:00",
                "orderId": "TZ170271958058707",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "",
                "eventType": 2,
                "acId": "",
                "channelCode": "3000",
                "platformOid": "VB2023121617393217359578",
                "payIp": "127.0.0.1",
                "payRegion": "0|0|0|内网IP|内网IP",
                "payDevice": "Windows",
                "resourceUrl": "",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 2,
                "cbStatus": 2,
                "expTime": "2023-12-16T17:44:32.723+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 17,
                "CreatedAt": "2023-12-18T10:19:14.404+08:00",
                "UpdatedAt": "2023-12-16T10:25:30.859+08:00",
                "orderId": "TZ170269316093533",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "20231215230742445",
                "eventType": 2,
                "acId": "88278125",
                "channelCode": "3000",
                "platformOid": "VB2023121610191417358470",
                "payIp": "106.39.151.103",
                "payRegion": "中国|0|北京|北京市|电信",
                "payDevice": "Android",
                "resourceUrl": "weixin://wxpay/bizpayurl?pr=MxaBSkJzz",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 3,
                "cbStatus": 2,
                "expTime": "2023-12-16T10:25:28.079+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 16,
                "CreatedAt": "2023-12-18T23:08:21.272+08:00",
                "UpdatedAt": "2023-12-15T23:14:34.845+08:00",
                "orderId": "TZ170265290376457",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "20231215230742202",
                "eventType": 2,
                "acId": "88278125",
                "channelCode": "3000",
                "platformOid": "VB2023121523082117356781",
                "payIp": "106.39.149.147",
                "payRegion": "中国|0|北京|北京市|电信",
                "payDevice": "Android",
                "resourceUrl": "weixin://wxpay/bizpayurl?pr=v9lVtQozz",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 3,
                "cbStatus": 2,
                "expTime": "2023-12-15T23:14:31.242+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 15,
                "CreatedAt": "2023-12-18T23:01:25.037+08:00",
                "UpdatedAt": "2023-12-15T23:06:55.782+08:00",
                "orderId": "TZ170265248749932",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "20231215230116329",
                "eventType": 2,
                "acId": "88278125",
                "channelCode": "3000",
                "platformOid": "VB2023121523012517356764",
                "payIp": "106.39.149.147",
                "payRegion": "中国|0|北京|北京市|电信",
                "payDevice": "Windows",
                "resourceUrl": "weixin://wxpay/bizpayurl?pr=v9lVtQozz",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 3,
                "cbStatus": 2,
                "expTime": "2023-12-15T23:06:52.383+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 14,
                "CreatedAt": "2023-12-18T22:42:02.235+08:00",
                "UpdatedAt": "2023-12-15T22:48:08.78+08:00",
                "orderId": "TZ170265132481701",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "20231215224139551",
                "eventType": 2,
                "acId": "88278125",
                "channelCode": "3000",
                "platformOid": "VB2023121522420217356715",
                "payIp": "106.39.149.147",
                "payRegion": "中国|0|北京|北京市|电信",
                "payDevice": "Windows",
                "resourceUrl": "weixin://wxpay/bizpayurl?pr=v9lVtQozz",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 3,
                "cbStatus": 2,
                "expTime": "2023-12-15T22:47:18.517+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            },
            {
                "ID": 13,
                "CreatedAt": "2023-12-18T22:14:39.402+08:00",
                "UpdatedAt": "2023-12-15T22:15:04.118+08:00",
                "orderId": "TZ170264968113604",
                "money": 10,
                "unitPrice": 0,
                "unitId": 0,
                "eventId": "",
                "eventType": 2,
                "acId": "",
                "channelCode": "3000",
                "platformOid": "VB2023121522143917356646",
                "payIp": "106.39.149.147",
                "payRegion": "中国|0|北京|北京市|电信",
                "payDevice": "Windows",
                "resourceUrl": "",
                "pAccount": "50813324",
                "notifyUrl": "http://127.0.0.1/callback",
                "orderStatus": 0,
                "cbStatus": 2,
                "expTime": "2023-12-15T22:19:39.368+08:00",
                "cbTime": "0001-01-01T00:00:00Z",
                "CreatedBy": 4,
                "UpdatedBy": 0,
                "DeletedBy": 0
            }
        ]`

	// 解析 JSON 数据
	var dataList []vbox.PayOrder
	if err := json.Unmarshal([]byte(dataJSON), &dataList); err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	// 设置东八区时区
	location, _ := time.LoadLocation("Asia/Shanghai")
	/// 指定开始和结束时间
	startTime := time.Date(2023, 12, 17, 0, 0, 0, 0, location)
	endTime := time.Date(2023, 12, 20, 0, 0, 0, 0, location)

	// 指定时间间隔
	interval := 30 * time.Minute
	//interval := 1 * time.Hour

	// 调用函数计算结果
	resultMap := calculateTotalMoney(dataList, startTime, endTime, interval)

	// 输出结果
	for key, value := range resultMap {
		fmt.Printf("Time: %v, Total Money: %v\n", key, value)
	}
}

func calculateTotalMoney(dataList []vbox.PayOrder, startTime time.Time, endTime time.Time, interval time.Duration) []response.DataOverView {
	// 初始化结果映射
	resultMap := make(map[string]int)
	var keys []string
	var sortedResult []response.DataOverView
	// 遍历数据并按时间间隔累加 money
	for _, item := range dataList {
		if item.CreatedAt.After(startTime) && item.CreatedAt.Before(endTime) {
			// 计算所属的时间间隔
			location, _ := time.LoadLocation("Asia/Shanghai")
			// 先将时间调整到当天的零时零分零秒
			intervalEnd := item.CreatedAt.Truncate(interval)

			// 再进行时区调整
			intervalEnd = intervalEnd.Add(interval).In(location)
			key := intervalEnd.Format("2006-01-02 15:04:05")
			resultMap[key] += item.Money

			// 添加 key 到有序列表
			keys = append(keys, key)
		}
	}

	// 对 keys 进行排序
	sort.Strings(keys)
	for _, key := range keys {
		sortedResult = append(sortedResult, response.DataOverView{
			Y: resultMap[key],
			X: key,
		})
	}
	return sortedResult
}
*/
