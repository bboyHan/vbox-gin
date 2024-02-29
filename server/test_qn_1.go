package main

//import (
//	"fmt"
//	"golang.org/x/text/encoding/simplifiedchinese"
//	"io/ioutil"
//	"net/http"
//	"strings"
//)
//
//func main() {
//
//	url := "https://trade.taobao.com/trade/itemlist/asyncSold.htm?event_submit_do_query=1&_input_charset=utf8&prePageNo=1&sifg=0&action=itemlist%2FSoldQueryAction&queryMore=false&close=0&pageNum=1&isQnNew=true&isHideNick=true&pageSize=30"
//	method := "POST"
//
//	client := &http.Client{}
//	req, err := http.NewRequest(method, url, nil)
//
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	req.Header.Add("Cookie", "Cookie: DI_T_=CvCyYhs4fx1SHMLxwCHxDHh5AoWsb; unb=291897500; lgc=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; cancelledSubSites=empty; dnk=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; tracknick=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; _l_g_=Ug%3D%3D; sg=%E6%B3%AE00; cookie1=BYe81wL1aEcXISEX05eGKFLHATjgHSMwpdZLh%2Bonxys%3D; lid=%E5%AE%9D%E5%AE%9D%E6%B3%AE%E6%B3%AE; cna=FGlaHnsb10QCAXU+tSXEiMHd; thw=cn; cookie2=10623e205f58c7e968229f67304d1d02; tbcp=e=UoM%2BHFG%2BH40YFva9%2BW9MM%2Bo%3D&f=UUjZeloosIiw2%2BCvtr5iVE1G0QM%3D; xlly_s=1; cookie17=UUGjOpdJllU9; _nk_=%5Cu5B9D%5Cu5B9D%5Cu6CEE%5Cu6CEE; uc1=existShop=true&cookie21=U%2BGCWk%2F7owY2UcgNjKoRpw%3D%3D&cookie15=VT5L2FSpMGV7TQ%3D%3D&cookie14=UoYenMybgOipZw%3D%3D&pas=0&cookie16=W5iHLLyFPlMGbLDwA%2BdvAGZqLg%3D%3D; sn=; uc3=nk2=0uNrG6CNYqo%3D&vt3=F8dD3er%2F0loI8jNL%2Fno%3D&lg2=VFC%2FuZ9ayeYq2g%3D%3D&id2=UUGjOpdJllU9; csg=1fe052a3; t=a4b6b35c23815dd0ed51c362682a4ba9; skt=2e755f87f775e356; existShop=MTcwOTEyMDgzNg%3D%3D; uc4=nk4=0%400FJ7kRcJ2hk1GuZTgLlt5gCbyA%3D%3D&id4=0%40U2OU9SmOE7zVKGEpEatYkPdWeiI%3D; publishItemObj=; _cc_=Vq8l%2BKCLiw%3D%3D; _tb_token_=e3333e9164eee; sgcookie=P100y%2FFjvfaJ%2FFNhpfPE1vitYYowDBTiaTDBy77VTLYDlikqsdMRVGsQ0u%2FZwzdZOIp4ZwfMa702n1HY3k1iZkJfrWNuLsbvSbPBvs19V%2B%2FTj89Bn3ihwJsDkpklVDcXAk%2BT; lc=V3ic9Tykb4JHIbVnVQ%3D%3D; mtop_partitioned_detect=1; _m_h5_tk=967854cac29d1dba60bf9474c0bd6e25_1709128037680; _m_h5_tk_enc=f23e6c71538f5a41a7597826731973e6; tfstk=eXvHRev66B5CKfyXXeBIBPtmvVhORy65mUeRyTQr_N71v_KPyTYkrUbda26pEFYB56dp4LIyEN-V2DgSAdxy5U2Jw9CJEFxX4wydpUIuAdY0eJ3IO3yleTuxkxHvAMB5UqUt7umpAoMFa-kxHHtUvMovbxnBkpd1x7D-3a49EDXGoo2JpJ9T7K_2Yw-pvdSQkZ-F58yP-Mfh9H7gU8J142Vag9c0FGoJQ7N5TGsGkDH9aYdTsUYrjcVvV6S1cZnij7iNTGsMwcmgGGCFfMVA.; isg=BLm5V_y9-HcaNqTvMIujGDVF0CWTxq14kI6uPdvuN-BfYtn0Ixb7SFY04GaUWkWw; DI_T_=CvCyYhs4fx1SHMLxwCHxDHh5AoWsb; 3PcFlag=1709118191551; _samesite_flag_=true; _tb_token_=361d8db3b935e; cookie2=10623e205f58c7e968229f67304d1d02; t=c2b992be677a4265cbe3861511edfdcf; uc1=cookie21=Vq8l%2BKCLiv0NZbsVmuOI%2Bw%3D%3D&cookie16=WqG3DMC9UpAPBHGz5QBErFxlCA%3D%3D&cookie14=UoYenM23TMAiFw%3D%3D&existShop=true&pas=0")
//	req.Header.Add("Referer", " https://qn.taobao.com/")
//
//	res, err := client.Do(req)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	defer res.Body.Close()
//
//	reader := simplifiedchinese.GB18030.NewDecoder().Reader(res.Body)
//	body, err := ioutil.ReadAll(reader)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	a := string(body)
//
//	fmt.Println(a)
//	if strings.Contains(a, "发货") {
//		fmt.Println("xxxxx")
//	}
//}
