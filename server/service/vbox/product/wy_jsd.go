package product

//import (
//	"fmt"
//	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
//	"golang.org/x/text/encoding/simplifiedchinese"
//	"regexp"
//)
//
//func QryWYRecord(ck string) error {
//	client := vbHttp.NewHTTPClient()
//
//	URL := "http://ecard.163.com/account/query_balance"
//
//	h := map[string]string{
//		//"Cookie":          "_ga=GA1.1.597369612.1679625421; c=roX5vpg4-1688456833510-f8efeab11a76f-497839763; _fmdata=41Jel7rLilqZ0tq7Br1AtjJpM9juFNcPMpmpWEeZ3gRuKLm90pLDwFmcNKyG62ttGbz%2BVoySNtNPfsaAGkn0Ew%3D%3D; _ntes_nnid=1cd5d45a7fb78efb176a95093d841144,1690376217991; _ntes_nuid=1cd5d45a7fb78efb176a95093d841144; _ga_C6TGHFPQ1H=GS1.1.1692882820.6.0.1692882820.0.0.0; Qs_lvt_382223=1679625420%2C1687875665%2C1688122057%2C1688133597%2C1692882820; Qs_pv_382223=3664678399769109500%2C3456671677726777000%2C1702658761295782400%2C724781406661779100%2C3106578160469569500; _clck=gn5431|2|fef|0|1273; _ns=NS1.2.1409026210.1708001399; 9AD585D8A7CB034A=roX5vpg4-1688456833510-f8efeab11a76f-497839763; _xid=YkhkMIQ7I5j%2F0JZ3pBExAZq6pe86dvAB9%2BFZhS1dJOM%3D; 1735D64331DF397E=41Jel7rLilqZ0tq7Br1AtjJpM9juFNcPMpmpWEeZ3gQOQN0vGpm3Z99O%2BlaM85aIBR7m%2BFteMKyU5gn1lu%2Bx8Q%3D%3D; timing_user_id=time_KIOc4cTWjn; NTES_SESS=EMcGQf9_yakPmvQyqIMHBq.iA4RWT_0BYr6jvJnFbML.ybvYyrnuZzhuQsz9RFz1.1BN8C.Ro_YVzardq6sqo.HqFs0XbbZwV0yjJY1bgtgBSSipCxOWBSk38cr8TgXaiV5qSarY9sf6VVonQy36IJYasY3QMBYdoTNcD0Tum8kDqj0.On679lhyjk6tMydL52GjkrbP87nb14ZEFUTbFwYN7; S_INFO=1710942828|0|##3&80|bboyhan@yeah.net; P_INFO=bboyhan@yeah.net|1710942828|0|ecard|00&99|bej&1710940321&ecard#bej&null#10#0#0|&0|ecard&mtoken_client|bboyhan@yeah.net; sid=EYtSZhYAefWMPPfRmLVsYRzexdk58rPhdDQGJWfn",
//		"Cookie":          ck,
//		"Accept-Encoding": "deflate, br, zstd",
//		"Accept-Language": "zh-CN,zh;q=0.9",
//	}
//
//	op := &vbHttp.RequestOptions{
//		Headers: h,
//	}
//
//	resp, err := client.Get(URL, op)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	if resp != nil && resp.StatusCode != 200 {
//		fmt.Println("请求失败, code", resp.StatusCode)
//		return fmt.Errorf("ck失效或账号异常")
//	} else {
//		readerBody, _ := simplifiedchinese.GBK.NewDecoder().Bytes(resp.Body)
//
//		hb := string(readerBody)
//		fmt.Println(hb)
//
//		re := regexp.MustCompile(`<span id="urs_confirm">([^<]+)</span>`)
//		matches := re.FindAllStringSubmatch(hb, -1)
//
//		for _, match := range matches {
//			if len(match) > 1 {
//				alertMessage := match[1]
//				fmt.Println(alertMessage)
//				// 账号设置
//
//			}
//		}
//
//		reM := regexp.MustCompile(`<a\s+href="/ecard\?init_reason=2"\s+class="red"\s+style="float:right;">充值</a><span\s+class="red bold">(\d+)\s+点</span>`)
//		matchesM := reM.FindAllStringSubmatch(hb, -1)
//
//		for _, match := range matchesM {
//			if len(match) > 1 {
//				alertMessage := match[1]
//				fmt.Println(alertMessage)
//
//				// 积分余额
//			}
//		}
//	}
//}
