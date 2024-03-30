package main

//import (
//	"fmt"
//	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
//	"golang.org/x/text/encoding/simplifiedchinese"
//	"regexp"
//)
//
//func main() {
//
//	client := vbHttp.NewHTTPClient()
//
//	URL := "http://ecard.163.com/account/query_balance"
//
//	h := map[string]string{
//		"Cookie":          "__snaker__id=lWa4jeCRwRtSSuXJ;_gid=GA.0177113652.73400298616245;_ga=GA.1.24316c8142958.cdaad1ecff60f3648f89;BIDUPSID=B8BEC46820846D3CB755CDA58DD38B92;utid=umW42E0i5bELREdHp6F6kJHaiK17qy8Z;l_s_ecardTdcueVw=DD4D25C67BE43EF79EA5C33ECC6A540F64DAB43D6C016FDB20B24393584BD42177CA1590D94F89F3B7487D9B0525AB3A5320D27F5E022CD733881A817ACADD8253EC4C3AA2C7E18E960EA46B9D8D7DF51E69BD7E0B591BF2CE60819BD9813EF9A07337B7F2D524D33C6AD673EA6FBF35;NTES_WEB_FP=b19583edd163eea5c265d547689fed11;gdxidpyhxdE=dY9TOiXn%5CwwsJ25ta5%2FtQ0wwwH1%2BmEhLV7Y%5CAw%2FNJYGd3vXrrKz%2BL2eYEPnAwSOwcv3eNxSAnvjbJABiQiZN%2BHlTZCmu6tzKZNydwZkLo8exZaeHgpRtdeJGDqdLmdA%2BPXuSqLCTOPvNv4W5Gmp2ck5Tr2LXom38E7e2miL%5Cm9G5Pqu2%3A1711693664988;NTES_P_UTID=umW42E0i5bELREdHp6F6kJHaiK17qy8Z|1711692779;NTES_SESS=9hKe3OS0Va1dNcadieFmD8JQ1KsRUebBrbqbfVm7g_DxIBoZIK0.mF8.gkFD_xFvzvW1nXz_dyZHFGKJbSkbdzTwj8.Wj3B0cfqNRf8i1NkEaJT6cmoMQIy6.4oDQPg8A384MylmRAQfBC7W4I59sWoYXcWx7PzfAXvw06i3l66LBW0xUtGKlzsfgqxYdP9Ge4ZM61y7PusivZhgnSGhTEZ1T;S_INFO=1711692779|0|##3&80|bboyhan@yeah.net;P_INFO=bboyhan@yeah.net|1711692779|0|ecard|00&99|bej&1711526171&ecard#bej&null#10#0#0|&0|ecard|bboyhan@yeah.net;THE_LAST_LOGIN=bboyhan@yeah.net;sid=5WpAcDtWf6P3L2jrZVmn_Lmg-6Pxc0MsSGbskDFr;timing_user_id=time_fT2V2msGmV;",
//		"Accept-Encoding": "deflate, br, zstd",
//		"Accept-Language": "zh-CN,zh;q=0.9",
//	}
//
//	op := &vbHttp.RequestOptions{
//		Headers: h,
//	}
//	resp, err := client.Get(URL, op)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	if resp != nil && resp.StatusCode != 200 {
//		fmt.Println("请求失败")
//		fmt.Println("code", resp.StatusCode)
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
//			}
//		}
//	}
//
//}
