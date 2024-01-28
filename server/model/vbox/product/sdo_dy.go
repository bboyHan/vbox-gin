package product

type SdoDaoYuOrderRecord struct {
	OrderType       int     `json:"orderType" form:"orderType" url:"orderType"`
	CanPay          int     `json:"canPay" form:"canPay" url:"canPay"`
	OrderId         string  `json:"orderId" form:"orderId" url:"orderId"`
	AppIcon         string  `json:"appIcon" form:"appIcon" url:"appIcon"`
	OrderAmount     float64 `json:"orderAmount" form:"orderAmount" url:"orderAmount"`
	AreaName        string  `json:"areaName" form:"areaName" url:"areaName"`
	AppId           string  `json:"appId" form:"appId" url:"appId"`
	ConsumeType     int     `json:"consumeType" form:"consumeType" url:"consumeType"`
	PayUserType     string  `json:"payUserType" form:"payUserType" url:"payUserType"`
	TraceId         string  `json:"traceId" form:"traceId" url:"traceId"`
	ConsumeAmount   float64 `json:"consumeAmount" form:"consumeAmount" url:"consumeAmount"`
	DqCount         int     `json:"dqCount" form:"dqCount" url:"dqCount"`
	AppName         string  `json:"appName" form:"appName" url:"appName"`
	SndaId          string  `json:"sndaId" form:"sndaId" url:"sndaId"`
	DisplayAccount  string  `json:"displayAccount" form:"displayAccount" url:"displayAccount"`
	PayProductName  string  `json:"payProductName" form:"payProductName" url:"payProductName"`
	AreaId          string  `json:"areaId" form:"areaId" url:"areaId"`
	PayProductModel string  `json:"payProductModel" form:"payProductModel" url:"payProductModel"`
	TransferType    int     `json:"transferType" form:"transferType" url:"transferType"`
	TargetAccount   string  `json:"targetAccount" form:"targetAccount" url:"targetAccount"`
	PayStatus       int     `json:"payStatus" form:"payStatus" url:"payStatus"`
	DqTypeName      string  `json:"dqTypeName" form:"dqTypeName" url:"dqTypeName"`
	TimestampMs     int64   `json:"timestampMs" form:"timestampMs" url:"timestampMs"`

	//OrderTimeStr string    `json:"orderTime"`
	//OrderTime    time.Time // 用于保存解析后的时间值

}

type SdoDaoYuOrderResponse struct {
	ReturnCode    int                   `json:"return_code"`
	ReturnMessage string                `json:"return_message"`
	ErrHint       string                `json:"ErrHint"`
	Data          []SdoDaoYuOrderRecord `json:"data"`
}

//func (o *SdoDaoYuOrderRecord) UnmarshalJSON(data []byte) error {
//	type Alias SdoDaoYuOrderRecord
//	aux := &struct {
//		OrderTime string `json:"orderTime"`
//		*Alias
//	}{
//		Alias: (*Alias)(o),
//	}
//	if err := json.Unmarshal(data, &aux); err != nil {
//		return err
//	}
//	orderTime, err := time.Parse(customTimeLayout, aux.OrderTimeStr)
//	if err != nil {
//		return err
//	}
//	o.OrderTime = orderTime
//
//	return nil
//}
//
//// ToDaoYuJSON 方法将 SdoDaoYuOrderRecord 结构体转换为 JSON 字符串
//func (o *SdoDaoYuOrderRecord) ToDaoYuJSON() ([]byte, error) {
//	type Alias SdoDaoYuOrderRecord
//	return json.Marshal(&struct {
//		OrderTime string `json:"orderTime"`
//		*Alias
//	}{
//		OrderTime: o.OrderTime.Format(customTimeLayout),
//		Alias:     (*Alias)(o),
//	})
//}
//
//// ToJSON 方法将 SdoOrderRecord 结构体转换为 JSON 字符串
//func (o *SdoDaoYuOrderRecord) ToJSON() ([]byte, error) {
//	type Alias SdoDaoYuOrderRecord
//	return json.Marshal(&struct {
//		OrderTime string `json:"orderTime"`
//		*Alias
//	}{
//		OrderTime: o.OrderTime.Format(customTimeLayout),
//		Alias:     (*Alias)(o),
//	})
//}
