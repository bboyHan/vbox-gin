package tx

type Payment struct {
	CurrentState    string `json:"CurrentState" form:"CurrentState"`       // 当前状态
	LoginID         string `json:"LoginID" form:"LoginID"`                 // 登录ID
	OfferID         string `json:"OfferID" form:"OfferID"`                 // 报价ID
	OrigPayAmt      string `json:"OrigPayAmt" form:"OrigPayAmt"`           // 原始支付金额
	PayAmt          string `json:"PayAmt" form:"PayAmt"`                   // 支付金额
	PayChannel      string `json:"PayChannel" form:"PayChannel"`           // 支付渠道
	PayChannelSubId string `json:"PayChannelSubId" form:"PayChannelSubId"` // 支付渠道子ID
	PayItem         string `json:"PayItem" form:"PayItem"`                 // 支付项目
	PayTime         string `json:"PayTime" form:"PayTime"`                 // 支付时间
	PayUnit         string `json:"PayUnit" form:"PayUnit"`                 // 支付单位
	ProductName     string `json:"ProductName" form:"ProductName"`         // 产品名称
	ProvideID       string `json:"ProvideID" form:"ProvideID"`             // 提供商ID
	ProvideOpenID   string `json:"ProvideOpenID" form:"ProvideOpenID"`     // 提供商OpenID
	ProvideUin      string `json:"Provide Uin" form:"Provide Uin"`         // 提供商Uin
	QQUin           string `json:"QQUin" form:"QQUin"`                     // QQ Uin
	SerialNo        string `json:"SerialNo" form:"SerialNo"`               // 流水号
	ServiceCode     string `json:"ServiceCode" form:"ServiceCode"`         // 服务代码
	ServiceName     string `json:"ServiceName" form:"ServiceName"`         // 服务名称
	ServiceType     string `json:"ServiceType" form:"ServiceType"`         // 服务类型
	SessionType     string `json:"SessionType" form:"SessionType"`         // 会话类型
	ShowName        string `json:"ShowName" form:"ShowName"`               // 显示名称
	SystemType      string `json:"SystemType" form:"SystemType"`           // 系统类型
	UsrOpenID       string `json:"UsrOpenID" form:"UsrOpenID"`             // 用户OpenID
}

type Records struct {
	Payments []Payment `json:"WaterList" form:"WaterList"`
	Ret      int       `json:"ret" form:"ret"`
	Msg      string    `json:"msg" form:"msg"`
	PageNum  int       `json:"PageNum" form:"PageNum"`
	PageSize int       `json:"PageSize" form:"PageSize"`
}
