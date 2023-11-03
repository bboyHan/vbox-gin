package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"time"
)

type UserWalletSearch struct {
	vbox.UserWallet
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

type UserWalletTransfer struct {
	CurrentUid uint   `json:"currentUid" form:"currentUid" ` //当前用户ID
	Username   string `json:"username" form:"username"`      //当前用户名
	ToUid      uint   `json:"toUid" form:"toUid" `           //（划转至）用户ID
	ToUsername string `json:"toUsername" form:"toUsername"`  //（划转至）用户名
	Recharge   int    `json:"recharge" form:"recharge"`      //积分
	Type       int    `json:"type" form:"type"`              //事件类型， 1 - 充值/划转、2 - 订单积分消费
}
