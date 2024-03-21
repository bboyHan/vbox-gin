package product

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"go.uber.org/zap"
)

func QryDyRecord(Cookie string) (*product.DyWalletInfoRecord, error) {
	var dyHeaders = map[string]string{
		"Content-Type": "application/json",
		"Cookie":       Cookie,
	}
	var dyOptions = &vbHttp.RequestOptions{
		Headers:      dyHeaders,
		MaxRedirects: 0,
		PayloadType:  "url",
	}
	URL := "https://www.douyin.com/webcast/wallet/info/?account_type=0&aid=1128&fp=verify_ltn6aev0_ac0e6781_4dc4_f5cd_0898_3949e4d4e049&msToken=M3s0ZH3LcA8tFkmaDsbwgb4gWvbaUoT67QIuGDCJ_6C0UHqlMH4iV6Cl7Guoe7KbuxSGJ6CdJJXVX0x4DhoDco4Ild0EqOdc9gf4ORd6HfJr9iBO_f5rAEraskMgL4A=&X-Bogus=DFSzswVOLTxANGyItLh8eORXoR8b&_signature=_02B4Z6wo0000143SVaAAAIDCY.eRXRtit.ON0lEAAIaJnrtEq5NyECCkynLnKA4oGKd49c9XAQTOKVexwGPqsRsHfqF7ARB8RQM44ii6sAWhAc5n4kkU0NgJhdhi5pCGdVzqO2emvJa2Hg8M7c"

	client := vbHttp.NewHTTPClient()

	resp, err := client.Get(URL, dyOptions)
	if err != nil {
		global.GVA_LOG.Error("err:  ->", zap.Error(err))
		return nil, err
	}

	var record product.DyWalletInfoResponse
	err = json.Unmarshal(resp.Body, &record)
	if err != nil {
		global.GVA_LOG.Error("json.Unmarshal:  ->", zap.Error(err))
		return nil, err
	}
	if record.StatusCode == 0 {
		resp := &product.DyWalletInfoRecord{
			Money:   record.Data.Money,
			UserID:  record.Data.UserID,
			Diamond: record.Data.Diamond,
		}
		return resp, nil
	} else {
		global.GVA_LOG.Error("QryDyRecord 官方查单异常", zap.Any("record", record))
		return nil, fmt.Errorf("查询系统异常: %s", record.Data.Message)
	}
}
