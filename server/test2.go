package main

/*
import (
	"fmt"
	"regexp"
)

func main() {
	// 这是你的JSON字符串
	jsonData := `{"ret":0,"err_code":"","msg":"","flex_attack_info":{"fk_info":"","verify_url":"","fk_amt":""},"info":{"spbank_partnerid":"1000023401","billno":"","channel_info":"","count":"1","mburl":"","sms_info":"","remain":0,"cfturl":"C3B629A905F772AF56EDB2FB53C20E20AA90DED0E9D1379BF2E28CE93194092A2854DE311FFA335FF4E8CE6E0AA51659_36","mcardserialno":"","verifycode":"","verifysession":"","qkbalance":"","mcard_balance":"","present_count":"0","qqacct_difference":"","success_url":"","mb_only_sms":"","need_change_key":"0","portal_serial_no":"QQACCT_SAVE-20240111-AIH0DIDUTYi2","is_allow_change":"","user_msg":"998AF0A5DAF0D1989C03E32F30A0EB904CD181033CB8D6B7F78720FAA11EC66E2F4610A11AF989751AC15F6A1AAB92D41B80C776AB39244797B635C99AA638B7_60","comm_config":{},"drm_resource":""},"wx_info":{"wx_token":"wx11101817483454400bb5b34bea9f9e0000","wx_partner":"","wx_appid":"","wx_time":"","wx_noncenum":"","wx_sign":"weixin://wxpay/bizpayurl?pr=kR2C6ETzz","wx_package":"","wx_signtype":"","wx_productid":"","wx_sign_url":"","wx_serialno":"product_SZ_20240111178150553"},"hf_info":{"tips":"","vtype":"","province":"","billno":"","channel":"","price":"","hfpay_pay_flow":"","hf_interfacetype":"","hfpay_url":"","outOrderId":"","solutionType":"","contentId":"","channelId":"","itemIndex":"","hfpay_channelInfo":"","accessnum":"","accessmsg":"","hf_extend":""},"qqwallet_info":{"qqwallet_appId":"","qqwallet_nonce":"","qqwallet_tokenId":"","qqwallet_pubAcc":"","qqwallet_bargainorId":"","qqwallet_sigType":"","qqwallet_sig":"","qqwallet_timestamp":"","qqwallet_seq":""},"sp_info":{"drm_goldcoupons_acct":"","channel_orderid":"","out_trade_no":""}, "rc_info": {"rc_type":"", "rc_msg":"", "rc_amt":"", "rc_policyid":""},"wechatapp_info":{"price_level":"","orderid":"","desc":"","sign":"","ext_info":""},"wechat_quickpass_info":{"sp_id":"","sys_provider":"","order_info":""}}`

	// 使用正则表达式提取 wx_sign 和 portal_serial_no 的值
	wxSignRegex := regexp.MustCompile(`"wx_sign":"(.*?)"`)
	portalSerialNoRegex := regexp.MustCompile(`"portal_serial_no":"(.*?)"`)

	// 在字符串中查找匹配项
	wxSignMatches := wxSignRegex.FindStringSubmatch(jsonData)
	portalSerialNoMatches := portalSerialNoRegex.FindStringSubmatch(jsonData)

	// 提取匹配项的值
	var wxSign, portalSerialNo string

	if len(wxSignMatches) > 1 {
		wxSign = wxSignMatches[1]
	}

	if len(portalSerialNoMatches) > 1 {
		portalSerialNo = portalSerialNoMatches[1]
	}

	// 打印结果
	fmt.Println("wx_sign:", wxSign)
	fmt.Println("portal_serial_no:", portalSerialNo)
}
*/
