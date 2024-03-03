package product

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"golang.org/x/text/encoding/simplifiedchinese"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

func QNValidCookie(ck string) (bool, error) {
	content := "https://trade.taobao.com/trade/itemlist/asyncSold.htm?event_submit_do_query=1&_input_charset=utf8&prePageNo=1&sifg=0&action=itemlist%2FSoldQueryAction&tabCode=success&buyerNick=&dateBegin=1709222400000&dateEnd=1709305200000&orderStatus=SUCCESS&rateStatus=ALL&pageSize=15&rxOldFlag=0&rxSendFlag=0&useCheckcode=false&tradeTag=0&rxHasSendFlag=0&auctionType=0&close=0&sellerNick=&notifySendGoodsType=ALL&sellerMemoFlag=0&useOrderInfo=false&logisticsService=ALL&isQnNew=true&pageNum=1&o2oDeliveryType=ALL&rxAuditFlag=0&queryOrder=desc&rxElectronicAuditFlag=0&queryMore=false&rxWaitSendflag=0&sellerMemo=0&rxElectronicAllFlag=0&rxSuccessflag=0&refund=ALL&errorCheckcode=false&mailNo=&yushouStatus=ALL&orderType=ALL&deliveryTimeType=ALL&queryTag=0&itemTitle=%E5%86%92%E9%99%A9%E5%B2%9B10000%E7%82%B9%E5%88%B8%E5%86%92%E9%99%A9%E5%B2%9B100%E5%85%83%E7%82%B9%E5%8D%A1%E7%9B%9B%E8%B6%A3%E6%B8%B8%E6%88%8F%E4%B8%80%E5%8D%A1%E9%80%9A%E8%87%AA%E5%8A%A8%E5%85%85%E5%80%BC&buyerEncodeId=&queryBizType=ALL&isHideNick=true"
	client := http.NewHTTPClient()
	var qnHeaders = map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"Accept-Encoding": "deflate, br",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Referer":         "https://myseller.taobao.com/",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Cookie":          ck,
	}

	var qnOptions = &http.RequestOptions{
		Headers:      qnHeaders,
		MaxRedirects: 0,
		PayloadType:  "url",
	}

	resp, _ := client.Post(content, qnOptions)
	if resp.StatusCode == 302 {
		fmt.Println("302")
		return false, fmt.Errorf("ck已失效")
	}
	return true, nil
}

func FindAccNick(ck string) (nickname string, err error) {
	re := regexp.MustCompile(`lid=([^;]+)`)
	match := re.FindStringSubmatch(ck)
	if len(match) > 1 {
		lidValue := match[1]
		// 使用 QueryUnescape 函数进行解码
		decodedLidValue, errX := url.QueryUnescape(lidValue)
		if errX != nil {
			fmt.Println("解码失败:", errX)
			return "", errX
		}

		fmt.Println("lid的值为:", lidValue)
		fmt.Println("lid的解码值为:", decodedLidValue)
		return decodedLidValue, nil
	} else {
		fmt.Println("未找到lid的值")
		return "", fmt.Errorf("未找到lid的值")
	}
}

//查的第一页,60单以内
//
//func QryQNRecords(vca vbox.ChannelAccount) ([]product.QNRecord, error) {
//
//	content := "https://trade.taobao.com/trade/itemlist/asyncSold.htm?event_submit_do_query=1&_input_charset=utf8&prePageNo=1&sifg=0&action=itemlist%2FSoldQueryAction&queryMore=false&close=0&pageNum=1&isQnNew=true&isHideNick=true&pageSize=60"
//	client := http.NewHTTPClient()
//
//	// 创建 HTTP 客户端实例
//
//	var qnHeaders = map[string]string{
//		"Content-Type":    "application/x-www-form-urlencoded",
//		"Accept-Encoding": "deflate, br",
//		"Accept-Language": "zh-CN,zh;q=0.9",
//		"Referer":         "https://qn.taobao.com/",
//		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
//		//"Cookie":          "Cookie: DI_T_=CvCyYhs4fx1SHMLxwCHxDHh5AoWsb; unb=291897500; lgc=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; cancelledSubSites=empty; dnk=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; tracknick=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; _l_g_=Ug%3D%3D; sg=%E6%B3%AE00; cookie1=BYe81wL1aEcXISEX05eGKFLHATjgHSMwpdZLh%2Bonxys%3D; lid=%E5%AE%9D%E5%AE%9D%E6%B3%AE%E6%B3%AE; cna=FGlaHnsb10QCAXU+tSXEiMHd; thw=cn; cookie2=10623e205f58c7e968229f67304d1d02; tbcp=e=UoM%2BHFG%2BH40YFva9%2BW9MM%2Bo%3D&f=UUjZeloosIiw2%2BCvtr5iVE1G0QM%3D; xlly_s=1; cookie17=UUGjOpdJllU9; _nk_=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; uc1=existShop=true&cookie21=U%2BGCWk%2F7owY2UcgNjKoRpw%3D%3D&cookie15=VT5L2FSpMGV7TQ%3D%3D&cookie14=UoYenMybgOipZw%3D%3D&pas=0&cookie16=W5iHLLyFPlMGbLDwA%2BdvAGZqLg%3D%3D; sn=; uc3=nk2=0uNrG6CNYqo%3D&vt3=F8dD3er%2F0loI8jNL%2Fno%3D&lg2=VFC%2FuZ9ayeYq2g%3D%3D&id2=UUGjOpdJllU9; csg=1fe052a3; t=a4b6b35c23815dd0ed51c362682a4ba9; skt=2e755f87f775e356; existShop=MTcwOTEyMDgzNg%3D%3D; uc4=nk4=0%400FJ7kRcJ2hk1GuZTgLlt5gCbyA%3D%3D&id4=0%40U2OU9SmOE7zVKGEpEatYkPdWeiI%3D; publishItemObj=; _cc_=Vq8l%2BKCLiw%3D%3D; _tb_token_=e3333e9164eee; sgcookie=P100y%2FFjvfaJ%2FFNhpfPE1vitYYowDBTiaTDBy77VTLYDlikqsdMRVGsQ0u%2FZwzdZOIp4ZwfMa702n1HY3k1iZkJfrWNuLsbvSbPBvs19V%2B%2FTj89Bn3ihwJsDkpklVDcXAk%2BT; lc=V3ic9Tykb4JHIbVnVQ%3D%3D; mtop_partitioned_detect=1; _m_h5_tk=967854cac29d1dba60bf9474c0bd6e25_1709128037680; _m_h5_tk_enc=f23e6c71538f5a41a7597826731973e6; tfstk=eXvHRev66B5CKfyXXeBIBPtmvVhORy65mUeRyTQr_N71v_KPyTYkrUbda26pEFYB56dp4LIyEN-V2DgSAdxy5U2Jw9CJEFxX4wydpUIuAdY0eJ3IO3yleTuxkxHvAMB5UqUt7umpAoMFa-kxHHtUvMovbxnBkpd1x7D-3a49EDXGoo2JpJ9T7K_2Yw-pvdSQkZ-F58yP-Mfh9H7gU8J142Vag9c0FGoJQ7N5TGsGkDH9aYdTsUYrjcVvV6S1cZnij7iNTGsMwcmgGGCFfMVA.; isg=BLm5V_y9-HcaNqTvMIujGDVF0CWTxq14kI6uPdvuN-BfYtn0Ixb7SFY04GaUWkWw",
//		"Cookie": vca.Token,
//	}
//
//	var qnOptions = &http.RequestOptions{
//		Headers:      qnHeaders,
//		MaxRedirects: 0,
//		PayloadType:  "url",
//		//Payload:      query,
//		//QueryParams:  params,
//	}
//
//	resp, _ := client.Post(content, qnOptions)
//
//	if resp.StatusCode == 302 {
//		return nil, fmt.Errorf("ck过期，acc info, %v", vca.AcAccount)
//	}
//
//	readerBody, err := simplifiedchinese.GB18030.NewDecoder().Bytes(resp.Body)
//
//	s := string(readerBody)
//	fmt.Println(s)
//	fmt.Println(resp.StatusCode)
//	fmt.Println(resp.Headers)
//
//	var m map[string]interface{}
//	err = json.Unmarshal(readerBody, &m)
//	var qnRecords []product.QNRecord
//
//	if i, ok := m["mainOrders"]; ok {
//		var qnRecord product.QNRecord
//		for _, ele := range i.([]interface{}) {
//			var flag bool
//			if e, ok := ele.(map[string]interface{}); ok {
//				if statusInfo, ok := e["statusInfo"]; ok {
//					if text, ok := statusInfo.(map[string]interface{})["text"]; ok {
//						//fmt.Printf(text.(string))
//						qnRecord.OrderStatus = text.(string)
//						//if (strings.Contains(s0, "卖家已发货")) || (strings.Contains(s0, "卖家已发货")) {
//						//	qnRecord.OrderStatus = 1
//						//}
//					} else {
//						fmt.Println("Error:", err)
//						flag = true
//					}
//				} else {
//					fmt.Println("Error:", err)
//					flag = true
//				}
//
//				if payInfo, ok := e["payInfo"]; ok {
//					if actualFee, ok := payInfo.(map[string]interface{})["actualFee"]; ok {
//						qnRecord.Money = actualFee.(string)
//						fmt.Printf(", %v", actualFee.(string))
//					} else {
//						fmt.Println("Error:", err)
//						flag = true
//					}
//				} else {
//					fmt.Println("Error:", err)
//					flag = true
//				}
//
//				if orderInfo, ok := e["orderInfo"]; ok {
//					if createTime, ok := orderInfo.(map[string]interface{})["createTime"]; ok {
//						fmt.Printf(", %v", createTime.(string))
//						t, err := time.Parse("2006-01-02 15:04:05", createTime.(string))
//						if err != nil {
//							fmt.Println("解析时间失败:", err)
//							flag = true
//							continue
//						}
//						qnRecord.CreateTime = t
//					} else {
//						flag = true
//						fmt.Println("Error:", err)
//					}
//				} else {
//					flag = true
//					fmt.Println("Error:", err)
//				}
//
//				if buyer, ok := e["buyer"]; ok {
//					if decodeNick, ok := buyer.(map[string]interface{})["decodeNick"]; ok {
//						qnRecord.Buyer = decodeNick.(string)
//						fmt.Printf(", %v", decodeNick.(string))
//					} else {
//						fmt.Println("Error:", err)
//						flag = true
//					}
//				} else {
//					fmt.Println("Error:", err)
//					flag = true
//				}
//
//				if subOrders, ok := e["subOrders"]; ok {
//					for _, subOrder := range subOrders.([]interface{}) {
//						if itemInfo, ok := subOrder.(map[string]interface{})["itemInfo"]; ok {
//							if title, ok := itemInfo.(map[string]interface{})["title"]; ok {
//								qnRecord.SkuTitle = title.(string)
//								fmt.Printf(", %v \n", title.(string))
//							} else {
//								fmt.Println("Error:", err)
//								flag = true
//							}
//						} else {
//							fmt.Println("Error:", err)
//							flag = true
//						}
//					}
//
//				} else {
//					flag = true
//					fmt.Println("Error:", err)
//				}
//				if !flag {
//					qnRecords = append(qnRecords, qnRecord)
//				}
//			} else {
//				flag = true
//				fmt.Println("Error:", err)
//			}
//		}
//
//	} else {
//		fmt.Println("Error:", err)
//	}
//
//	return qnRecords, err
//}

func QryQNRecords(rawURL string, vca vbox.ChannelAccount, start time.Time, end time.Time, titleID string) ([]product.QNRecord, error) {

	//content := "https://trade.taobao.com/trade/itemlist/asyncSold.htm?event_submit_do_query=1&_input_charset=utf8&prePageNo=1&sifg=0&action=itemlist%2FSoldQueryAction&tabCode=success&buyerNick=" +
	//	"&dateBegin=1709222400000&dateEnd=1709305200000" +
	//	"&orderStatus=SUCCESS" +
	//	"&rateStatus=ALL&pageSize=15&rxOldFlag=0&rxSendFlag=0&useCheckcode=false&tradeTag=0&rxHasSendFlag=0&auctionType=0&close=0&sellerNick=&notifySendGoodsType=ALL&sellerMemoFlag=0&useOrderInfo=false&logisticsService=ALL&isQnNew=true&pageNum=1&o2oDeliveryType=ALL&rxAuditFlag=0&queryOrder=desc&rxElectronicAuditFlag=0&queryMore=false&rxWaitSendflag=0&sellerMemo=0&rxElectronicAllFlag=0&rxSuccessflag=0&refund=ALL&errorCheckcode=false&mailNo=&yushouStatus=ALL&orderType=ALL&deliveryTimeType=ALL&queryTag=0" +
	//	"&itemTitle=" + titleID +
	//	"&buyerEncodeId=&queryBizType=ALL&isHideNick=true"

	u, _ := url.Parse(rawURL)
	queryParams := u.Query()

	// 当前时间毫秒数
	startSeconds := start.UnixMilli()
	endSeconds := end.UnixMilli()

	// orderStatus - SEND / SUCCESS
	queryParams.Set("orderStatus", "SEND")
	queryParams.Set("itemTitle", titleID)
	queryParams.Set("dateBegin", strconv.FormatInt(startSeconds, 10))
	queryParams.Set("dateEnd", strconv.FormatInt(endSeconds, 10))

	u.RawQuery = queryParams.Encode()
	newURL := u.String()

	client := http.NewHTTPClient()

	// 创建 HTTP 客户端实例

	var qnHeaders = map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"Accept-Encoding": "deflate, br",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Referer":         "https://qn.taobao.com/",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		//"Cookie":          "Cookie: DI_T_=CvCyYhs4fx1SHMLxwCHxDHh5AoWsb; unb=291897500; lgc=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; cancelledSubSites=empty; dnk=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; tracknick=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; _l_g_=Ug%3D%3D; sg=%E6%B3%AE00; cookie1=BYe81wL1aEcXISEX05eGKFLHATjgHSMwpdZLh%2Bonxys%3D; lid=%E5%AE%9D%E5%AE%9D%E6%B3%AE%E6%B3%AE; cna=FGlaHnsb10QCAXU+tSXEiMHd; thw=cn; cookie2=10623e205f58c7e968229f67304d1d02; tbcp=e=UoM%2BHFG%2BH40YFva9%2BW9MM%2Bo%3D&f=UUjZeloosIiw2%2BCvtr5iVE1G0QM%3D; xlly_s=1; cookie17=UUGjOpdJllU9; _nk_=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; uc1=existShop=true&cookie21=U%2BGCWk%2F7owY2UcgNjKoRpw%3D%3D&cookie15=VT5L2FSpMGV7TQ%3D%3D&cookie14=UoYenMybgOipZw%3D%3D&pas=0&cookie16=W5iHLLyFPlMGbLDwA%2BdvAGZqLg%3D%3D; sn=; uc3=nk2=0uNrG6CNYqo%3D&vt3=F8dD3er%2F0loI8jNL%2Fno%3D&lg2=VFC%2FuZ9ayeYq2g%3D%3D&id2=UUGjOpdJllU9; csg=1fe052a3; t=a4b6b35c23815dd0ed51c362682a4ba9; skt=2e755f87f775e356; existShop=MTcwOTEyMDgzNg%3D%3D; uc4=nk4=0%400FJ7kRcJ2hk1GuZTgLlt5gCbyA%3D%3D&id4=0%40U2OU9SmOE7zVKGEpEatYkPdWeiI%3D; publishItemObj=; _cc_=Vq8l%2BKCLiw%3D%3D; _tb_token_=e3333e9164eee; sgcookie=P100y%2FFjvfaJ%2FFNhpfPE1vitYYowDBTiaTDBy77VTLYDlikqsdMRVGsQ0u%2FZwzdZOIp4ZwfMa702n1HY3k1iZkJfrWNuLsbvSbPBvs19V%2B%2FTj89Bn3ihwJsDkpklVDcXAk%2BT; lc=V3ic9Tykb4JHIbVnVQ%3D%3D; mtop_partitioned_detect=1; _m_h5_tk=967854cac29d1dba60bf9474c0bd6e25_1709128037680; _m_h5_tk_enc=f23e6c71538f5a41a7597826731973e6; tfstk=eXvHRev66B5CKfyXXeBIBPtmvVhORy65mUeRyTQr_N71v_KPyTYkrUbda26pEFYB56dp4LIyEN-V2DgSAdxy5U2Jw9CJEFxX4wydpUIuAdY0eJ3IO3yleTuxkxHvAMB5UqUt7umpAoMFa-kxHHtUvMovbxnBkpd1x7D-3a49EDXGoo2JpJ9T7K_2Yw-pvdSQkZ-F58yP-Mfh9H7gU8J142Vag9c0FGoJQ7N5TGsGkDH9aYdTsUYrjcVvV6S1cZnij7iNTGsMwcmgGGCFfMVA.; isg=BLm5V_y9-HcaNqTvMIujGDVF0CWTxq14kI6uPdvuN-BfYtn0Ixb7SFY04GaUWkWw",
		"Cookie": vca.Token,
	}

	var qnOptions = &http.RequestOptions{
		Headers:      qnHeaders,
		MaxRedirects: 0,
		PayloadType:  "url",
		//Payload:      query,
		//QueryParams:  params,
	}

	resp, _ := client.Post(newURL, qnOptions)

	if resp.StatusCode == 302 {
		return nil, fmt.Errorf("ck过期，acc info, %v", vca.AcAccount)
	}

	readerBody, err := simplifiedchinese.GB18030.NewDecoder().Bytes(resp.Body)

	var m map[string]interface{}
	err = json.Unmarshal(readerBody, &m)
	var qnRecords []product.QNRecord

	if i, ok := m["mainOrders"]; ok {
		var qnRecord product.QNRecord
		for _, ele := range i.([]interface{}) {
			var flag bool
			if e, ok := ele.(map[string]interface{}); ok {
				if statusInfo, ok := e["statusInfo"]; ok {
					if text, ok := statusInfo.(map[string]interface{})["text"]; ok {
						//fmt.Printf(text.(string))
						qnRecord.OrderStatus = text.(string)
						//if (strings.Contains(s0, "卖家已发货")) || (strings.Contains(s0, "卖家已发货")) {
						//	qnRecord.OrderStatus = 1
						//}
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if payInfo, ok := e["payInfo"]; ok {
					if actualFee, ok := payInfo.(map[string]interface{})["actualFee"]; ok {
						qnRecord.Money = actualFee.(string)
						fmt.Printf(", %v", actualFee.(string))
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if orderInfo, ok := e["orderInfo"]; ok {
					if createTime, ok := orderInfo.(map[string]interface{})["createTime"]; ok {
						fmt.Printf(", %v", createTime.(string))
						t, err := time.Parse("2006-01-02 15:04:05", createTime.(string))
						if err != nil {
							fmt.Println("解析时间失败:", err)
							flag = true
							continue
						}
						qnRecord.CreateTime = t
					} else {
						flag = true
						fmt.Println("Error:", err)
					}
				} else {
					flag = true
					fmt.Println("Error:", err)
				}

				if buyer, ok := e["buyer"]; ok {
					if decodeNick, ok := buyer.(map[string]interface{})["decodeNick"]; ok {
						qnRecord.Buyer = decodeNick.(string)
						fmt.Printf(", %v", decodeNick.(string))
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if subOrders, ok := e["subOrders"]; ok {
					for _, subOrder := range subOrders.([]interface{}) {
						if itemInfo, ok := subOrder.(map[string]interface{})["itemInfo"]; ok {
							if title, ok := itemInfo.(map[string]interface{})["title"]; ok {
								qnRecord.SkuTitle = title.(string)
								fmt.Printf(", %v \n", title.(string))
							} else {
								fmt.Println("Error:", err)
								flag = true
							}
						} else {
							fmt.Println("Error:", err)
							flag = true
						}
					}

				} else {
					flag = true
					fmt.Println("Error:", err)
				}
				if !flag {
					qnRecords = append(qnRecords, qnRecord)
				}
			} else {
				flag = true
				fmt.Println("Error:", err)
			}
		}

	} else {
		fmt.Println("Error:", err)
	}

	// 2. 查 交易成功
	queryParams.Set("orderStatus", "SUCCESS")
	u.RawQuery = queryParams.Encode()
	newURL = u.String()

	resp, _ = client.Post(newURL, qnOptions)

	if resp.StatusCode == 302 {
		return nil, fmt.Errorf("ck过期，acc info, %v", vca.AcAccount)
	}

	readerBody, err = simplifiedchinese.GB18030.NewDecoder().Bytes(resp.Body)
	err = json.Unmarshal(readerBody, &m)

	if i, ok := m["mainOrders"]; ok {
		var qnRecord product.QNRecord
		for _, ele := range i.([]interface{}) {
			var flag bool
			if e, ok := ele.(map[string]interface{}); ok {
				if statusInfo, ok := e["statusInfo"]; ok {
					if text, ok := statusInfo.(map[string]interface{})["text"]; ok {
						//fmt.Printf(text.(string))
						qnRecord.OrderStatus = text.(string)
						//if (strings.Contains(s0, "卖家已发货")) || (strings.Contains(s0, "卖家已发货")) {
						//	qnRecord.OrderStatus = 1
						//}
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if payInfo, ok := e["payInfo"]; ok {
					if actualFee, ok := payInfo.(map[string]interface{})["actualFee"]; ok {
						qnRecord.Money = actualFee.(string)
						fmt.Printf(", %v", actualFee.(string))
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if orderInfo, ok := e["orderInfo"]; ok {
					if createTime, ok := orderInfo.(map[string]interface{})["createTime"]; ok {
						fmt.Printf(", %v", createTime.(string))
						t, err := time.Parse("2006-01-02 15:04:05", createTime.(string))
						if err != nil {
							fmt.Println("解析时间失败:", err)
							flag = true
							continue
						}
						qnRecord.CreateTime = t
					} else {
						flag = true
						fmt.Println("Error:", err)
					}
				} else {
					flag = true
					fmt.Println("Error:", err)
				}

				if buyer, ok := e["buyer"]; ok {
					if decodeNick, ok := buyer.(map[string]interface{})["decodeNick"]; ok {
						qnRecord.Buyer = decodeNick.(string)
						fmt.Printf(", %v", decodeNick.(string))
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if subOrders, ok := e["subOrders"]; ok {
					for _, subOrder := range subOrders.([]interface{}) {
						if itemInfo, ok := subOrder.(map[string]interface{})["itemInfo"]; ok {
							if title, ok := itemInfo.(map[string]interface{})["title"]; ok {
								qnRecord.SkuTitle = title.(string)
								fmt.Printf(", %v \n", title.(string))
							} else {
								fmt.Println("Error:", err)
								flag = true
							}
						} else {
							fmt.Println("Error:", err)
							flag = true
						}
					}

				} else {
					flag = true
					fmt.Println("Error:", err)
				}
				if !flag {
					qnRecords = append(qnRecords, qnRecord)
				}
			} else {
				flag = true
				fmt.Println("Error:", err)
			}
		}

	} else {
		fmt.Println("Error:", err)
	}

	return qnRecords, err
}

func QryQNRecordsByTitleID(rawURL string, vca vbox.ChannelAccount, start time.Time, end time.Time, titleID string) ([]product.QNRecord, error) {

	//content := "https://trade.taobao.com/trade/itemlist/asyncSold.htm?event_submit_do_query=1&_input_charset=utf8&prePageNo=1&sifg=0&action=itemlist%2FSoldQueryAction&tabCode=success&buyerNick=" +
	//	"&dateBegin=1709222400000&dateEnd=1709305200000" +
	//	"&orderStatus=SUCCESS" +
	//	"&rateStatus=ALL&pageSize=15&rxOldFlag=0&rxSendFlag=0&useCheckcode=false&tradeTag=0&rxHasSendFlag=0&auctionType=0&close=0&sellerNick=&notifySendGoodsType=ALL&sellerMemoFlag=0&useOrderInfo=false&logisticsService=ALL&isQnNew=true&pageNum=1&o2oDeliveryType=ALL&rxAuditFlag=0&queryOrder=desc&rxElectronicAuditFlag=0&queryMore=false&rxWaitSendflag=0&sellerMemo=0&rxElectronicAllFlag=0&rxSuccessflag=0&refund=ALL&errorCheckcode=false&mailNo=&yushouStatus=ALL&orderType=ALL&deliveryTimeType=ALL&queryTag=0" +
	//	"&itemTitle=" + titleID +
	//	"&buyerEncodeId=&queryBizType=ALL&isHideNick=true"

	u, _ := url.Parse(rawURL)
	queryParams := u.Query()

	// 当前时间毫秒数
	startSeconds := start.UnixMilli()
	endSeconds := end.UnixMilli()

	// orderStatus - SEND / SUCCESS
	queryParams.Set("orderStatus", "SEND")
	queryParams.Set("itemTitle", titleID)
	queryParams.Set("dateBegin", strconv.FormatInt(startSeconds, 10))
	queryParams.Set("dateEnd", strconv.FormatInt(endSeconds, 10))

	u.RawQuery = queryParams.Encode()
	newURL := u.String()

	client := http.NewHTTPClient()

	// 创建 HTTP 客户端实例

	var qnHeaders = map[string]string{
		"Content-Type":    "application/x-www-form-urlencoded",
		"Accept-Encoding": "deflate, br",
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Referer":         "https://qn.taobao.com/",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		//"Cookie":          "Cookie: DI_T_=CvCyYhs4fx1SHMLxwCHxDHh5AoWsb; unb=291897500; lgc=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; cancelledSubSites=empty; dnk=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; tracknick=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; _l_g_=Ug%3D%3D; sg=%E6%B3%AE00; cookie1=BYe81wL1aEcXISEX05eGKFLHATjgHSMwpdZLh%2Bonxys%3D; lid=%E5%AE%9D%E5%AE%9D%E6%B3%AE%E6%B3%AE; cna=FGlaHnsb10QCAXU+tSXEiMHd; thw=cn; cookie2=10623e205f58c7e968229f67304d1d02; tbcp=e=UoM%2BHFG%2BH40YFva9%2BW9MM%2Bo%3D&f=UUjZeloosIiw2%2BCvtr5iVE1G0QM%3D; xlly_s=1; cookie17=UUGjOpdJllU9; _nk_=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; uc1=existShop=true&cookie21=U%2BGCWk%2F7owY2UcgNjKoRpw%3D%3D&cookie15=VT5L2FSpMGV7TQ%3D%3D&cookie14=UoYenMybgOipZw%3D%3D&pas=0&cookie16=W5iHLLyFPlMGbLDwA%2BdvAGZqLg%3D%3D; sn=; uc3=nk2=0uNrG6CNYqo%3D&vt3=F8dD3er%2F0loI8jNL%2Fno%3D&lg2=VFC%2FuZ9ayeYq2g%3D%3D&id2=UUGjOpdJllU9; csg=1fe052a3; t=a4b6b35c23815dd0ed51c362682a4ba9; skt=2e755f87f775e356; existShop=MTcwOTEyMDgzNg%3D%3D; uc4=nk4=0%400FJ7kRcJ2hk1GuZTgLlt5gCbyA%3D%3D&id4=0%40U2OU9SmOE7zVKGEpEatYkPdWeiI%3D; publishItemObj=; _cc_=Vq8l%2BKCLiw%3D%3D; _tb_token_=e3333e9164eee; sgcookie=P100y%2FFjvfaJ%2FFNhpfPE1vitYYowDBTiaTDBy77VTLYDlikqsdMRVGsQ0u%2FZwzdZOIp4ZwfMa702n1HY3k1iZkJfrWNuLsbvSbPBvs19V%2B%2FTj89Bn3ihwJsDkpklVDcXAk%2BT; lc=V3ic9Tykb4JHIbVnVQ%3D%3D; mtop_partitioned_detect=1; _m_h5_tk=967854cac29d1dba60bf9474c0bd6e25_1709128037680; _m_h5_tk_enc=f23e6c71538f5a41a7597826731973e6; tfstk=eXvHRev66B5CKfyXXeBIBPtmvVhORy65mUeRyTQr_N71v_KPyTYkrUbda26pEFYB56dp4LIyEN-V2DgSAdxy5U2Jw9CJEFxX4wydpUIuAdY0eJ3IO3yleTuxkxHvAMB5UqUt7umpAoMFa-kxHHtUvMovbxnBkpd1x7D-3a49EDXGoo2JpJ9T7K_2Yw-pvdSQkZ-F58yP-Mfh9H7gU8J142Vag9c0FGoJQ7N5TGsGkDH9aYdTsUYrjcVvV6S1cZnij7iNTGsMwcmgGGCFfMVA.; isg=BLm5V_y9-HcaNqTvMIujGDVF0CWTxq14kI6uPdvuN-BfYtn0Ixb7SFY04GaUWkWw",
		"Cookie": vca.Token,
	}

	var qnOptions = &http.RequestOptions{
		Headers:      qnHeaders,
		MaxRedirects: 0,
		PayloadType:  "url",
		//Payload:      query,
		//QueryParams:  params,
	}

	resp, _ := client.Post(newURL, qnOptions)

	if resp.StatusCode == 302 {
		return nil, fmt.Errorf("ck过期，acc info, %v", vca.AcAccount)
	}

	readerBody, err := simplifiedchinese.GB18030.NewDecoder().Bytes(resp.Body)

	var m map[string]interface{}
	err = json.Unmarshal(readerBody, &m)
	var qnRecords []product.QNRecord

	if i, ok := m["mainOrders"]; ok {
		var qnRecord product.QNRecord
		for _, ele := range i.([]interface{}) {
			var flag bool
			if e, ok := ele.(map[string]interface{}); ok {
				if statusInfo, ok := e["statusInfo"]; ok {
					if text, ok := statusInfo.(map[string]interface{})["text"]; ok {
						//fmt.Printf(text.(string))
						qnRecord.OrderStatus = text.(string)
						//if (strings.Contains(s0, "卖家已发货")) || (strings.Contains(s0, "卖家已发货")) {
						//	qnRecord.OrderStatus = 1
						//}
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if payInfo, ok := e["payInfo"]; ok {
					if actualFee, ok := payInfo.(map[string]interface{})["actualFee"]; ok {
						qnRecord.Money = actualFee.(string)
						fmt.Printf(", %v", actualFee.(string))
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if orderInfo, ok := e["orderInfo"]; ok {
					if createTime, ok := orderInfo.(map[string]interface{})["createTime"]; ok {
						fmt.Printf(", %v", createTime.(string))
						t, err := time.Parse("2006-01-02 15:04:05", createTime.(string))
						if err != nil {
							fmt.Println("解析时间失败:", err)
							flag = true
							continue
						}
						qnRecord.CreateTime = t
					} else {
						flag = true
						fmt.Println("Error:", err)
					}
				} else {
					flag = true
					fmt.Println("Error:", err)
				}

				if buyer, ok := e["buyer"]; ok {
					if decodeNick, ok := buyer.(map[string]interface{})["decodeNick"]; ok {
						qnRecord.Buyer = decodeNick.(string)
						fmt.Printf(", %v", decodeNick.(string))
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if subOrders, ok := e["subOrders"]; ok {
					for _, subOrder := range subOrders.([]interface{}) {
						if itemInfo, ok := subOrder.(map[string]interface{})["itemInfo"]; ok {
							if title, ok := itemInfo.(map[string]interface{})["title"]; ok {
								qnRecord.SkuTitle = title.(string)
								fmt.Printf(", %v \n", title.(string))
							} else {
								fmt.Println("Error:", err)
								flag = true
							}
						} else {
							fmt.Println("Error:", err)
							flag = true
						}
					}

				} else {
					flag = true
					fmt.Println("Error:", err)
				}
				if !flag {
					qnRecords = append(qnRecords, qnRecord)
				}
			} else {
				flag = true
				fmt.Println("Error:", err)
			}
		}

	} else {
		fmt.Println("Error:", err)
	}

	// 2. 查 交易成功
	queryParams.Set("orderStatus", "SUCCESS")
	u.RawQuery = queryParams.Encode()
	newURL = u.String()

	resp, _ = client.Post(newURL, qnOptions)

	if resp.StatusCode == 302 {
		return nil, fmt.Errorf("ck过期，acc info, %v", vca.AcAccount)
	}

	readerBody, err = simplifiedchinese.GB18030.NewDecoder().Bytes(resp.Body)

	if i, ok := m["mainOrders"]; ok {
		var qnRecord product.QNRecord
		for _, ele := range i.([]interface{}) {
			var flag bool
			if e, ok := ele.(map[string]interface{}); ok {
				if statusInfo, ok := e["statusInfo"]; ok {
					if text, ok := statusInfo.(map[string]interface{})["text"]; ok {
						//fmt.Printf(text.(string))
						qnRecord.OrderStatus = text.(string)
						//if (strings.Contains(s0, "卖家已发货")) || (strings.Contains(s0, "卖家已发货")) {
						//	qnRecord.OrderStatus = 1
						//}
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if payInfo, ok := e["payInfo"]; ok {
					if actualFee, ok := payInfo.(map[string]interface{})["actualFee"]; ok {
						qnRecord.Money = actualFee.(string)
						fmt.Printf(", %v", actualFee.(string))
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if orderInfo, ok := e["orderInfo"]; ok {
					if createTime, ok := orderInfo.(map[string]interface{})["createTime"]; ok {
						fmt.Printf(", %v", createTime.(string))
						t, err := time.Parse("2006-01-02 15:04:05", createTime.(string))
						if err != nil {
							fmt.Println("解析时间失败:", err)
							flag = true
							continue
						}
						qnRecord.CreateTime = t
					} else {
						flag = true
						fmt.Println("Error:", err)
					}
				} else {
					flag = true
					fmt.Println("Error:", err)
				}

				if buyer, ok := e["buyer"]; ok {
					if decodeNick, ok := buyer.(map[string]interface{})["decodeNick"]; ok {
						qnRecord.Buyer = decodeNick.(string)
						fmt.Printf(", %v", decodeNick.(string))
					} else {
						fmt.Println("Error:", err)
						flag = true
					}
				} else {
					fmt.Println("Error:", err)
					flag = true
				}

				if subOrders, ok := e["subOrders"]; ok {
					for _, subOrder := range subOrders.([]interface{}) {
						if itemInfo, ok := subOrder.(map[string]interface{})["itemInfo"]; ok {
							if title, ok := itemInfo.(map[string]interface{})["title"]; ok {
								qnRecord.SkuTitle = title.(string)
								fmt.Printf(", %v \n", title.(string))
							} else {
								fmt.Println("Error:", err)
								flag = true
							}
						} else {
							fmt.Println("Error:", err)
							flag = true
						}
					}

				} else {
					flag = true
					fmt.Println("Error:", err)
				}
				if !flag {
					qnRecords = append(qnRecords, qnRecord)
				}
			} else {
				flag = true
				fmt.Println("Error:", err)
			}
		}

	} else {
		fmt.Println("Error:", err)
	}

	return qnRecords, err
}
