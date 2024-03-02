package task

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"strings"
	"time"
)

func HandleShopMoneyAvailable() (err error) {

	var ml []map[string]interface{}
	global.GVA_DB.Debug().Model(&vbox.ChannelShop{}).Select(`
			cid,
			created_by as uid,
			GROUP_CONCAT(DISTINCT money ORDER BY money) AS money
		`).
		Where("status = 1").Group("cid,created_by").Scan(&ml)

	for _, m := range ml {
		cid := m["cid"]
		uid := m["uid"]
		moneyCollect := m["money"]
		mList := strings.Split(moneyCollect.(string), ",")

		orgTmp := utils2.GetSelfOrg(uint(uid.(int64)))
		var orgId uint
		if len(orgTmp) > 0 {
			orgId = orgTmp[0]
		} else {
			continue
		}
		moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgId, cid)
		//global.GVA_REDIS.Del(context.Background(), moneyKey)

		for _, money := range mList {
			global.GVA_REDIS.SAdd(context.Background(), moneyKey, money)
		}
		global.GVA_REDIS.Expire(context.Background(), moneyKey, 1*time.Minute)
	}

	return err
}
