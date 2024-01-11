package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"time"
)

type ChannelPayCodeSearch struct {
	vbox.ChannelPayCode
	StartCreatedAt *time.Time `json:"startCreatedAt" form:"startCreatedAt"`
	EndCreatedAt   *time.Time `json:"endCreatedAt" form:"endCreatedAt"`
	request.PageInfo
}

type BatchPayCode struct {
	Cid       string            `json:"cid" form:"cid" `             //通道id
	AcAccount string            `json:"acAccount" form:"acAccount" ` //通道账户名
	AcId      string            `json:"acId" form:"acId" `           //账户ID
	AcRemark  string            `json:"acRemark" form:"acRemark" `   //账户备注
	ExpTime   *time.Time        `json:"expTime" form:"expTime" `     //过期时间
	Type      uint              `json:"type" form:"type"`            //上传方式
	CreatedBy uint              `json:"createdBy" form:"createdBy"`  //创建者
	List      []BatchPayCodeEle `json:"list" form:"list" `           //上传方式
}

type BatchPayCodeEle struct {
	Money      int    `json:"money" form:"money"`           //金额
	Operator   string `json:"operator" form:"operator"`     //运营商
	Location   string `json:"location" form:"location"`     //省市
	ImgBaseStr string `json:"imgBaseStr" form:"imgBaseStr"` //图片base64编码
	ImgContent string `json:"imgContent" form:"imgContent"` //图片base64编码
	Mid        string `json:"mid" form:"mid"`               //标识id
	PlatId     string `json:"platId" form:"platId"`         //平台id
}
