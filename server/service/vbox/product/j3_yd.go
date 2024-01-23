package product

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/product"
	vbHttp "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"go.uber.org/zap"
)

//var rawURL = "https://security.seasungame.com/security_extend_server/helper/balance/queryBalance?
//gameCode=jx3&account=18210889498&accountType=&zoneCode=z05&SN=98710648156&remark=&sign=36A360706FD189A2BF867D70F656C7BE"

// 校验官方合法性用一下

func QryJ3Record(vca vbox.ChannelAccount) (*product.J3BalanceData, error) {
	client := vbHttp.NewHTTPClient()

	resp, err := client.Post(vca.Token, options)
	if err != nil {
		global.GVA_LOG.Error("err:  ->", zap.Error(err))
		return nil, err
	}

	var record product.J3BalanceRecord
	err = json.Unmarshal(resp.Body, &record)
	if err != nil {
		global.GVA_LOG.Error("json.Unmarshal:  ->", zap.Error(err))
		return nil, err
	}
	if record.Code == 0 {
		return &record.Data, nil
	} else {
		global.GVA_LOG.Error("QryJ3Record 官方查单异常", zap.Any("record", record))
		return nil, fmt.Errorf("查询系统异常: %s", record.Message)
	}
}
