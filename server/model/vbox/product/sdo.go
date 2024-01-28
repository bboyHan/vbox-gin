package product

import (
	"encoding/json"
	"time"
)

type SdoOrderRecord struct {
	OrderID            string `json:"orderId"`
	SellerID           int    `json:"sellerId"`
	TraceID            string `json:"traceId"`
	OrderUser          string `json:"orderUser"`
	InputOrderUserType int    `json:"inputOrderUserType"`
	InputOrderUser     string `json:"inputOrderUser"`
	PayUser            string `json:"payUser"`
	InputPayUserType   int    `json:"inputPayUserType"`
	InputPayUser       string `json:"inputPayUser"`
	ProductName        string `json:"productName"`
	ProductModel       string `json:"productModel"`
	ProductUnit        string `json:"productUnit"`
	ProductValue       string `json:"productValue"`
	ProductAmount      string `json:"productAmount"`
	OrderAmount        string `json:"orderAmount"`
	DeliverAmount      string `json:"deliverAmount"`
	ExpiredTime        string `json:"expiredTime"`
	AppID              string `json:"appId"`
	AppName            string `json:"appName"`
	AreaID             string `json:"areaId"`
	AreaName           string `json:"areaName"`
	OrderType          int    `json:"orderType"`
	State              int    `json:"state"`
	Version            string `json:"version"`
	ProductID          string `json:"productId"`

	OrderTimeStr string `json:"orderTime"`
	StateTimeStr string `json:"stateTime"`
	// ... 其他字段
	OrderTime time.Time // 用于保存解析后的时间值
	StateTime time.Time // 用于保存解析后的时间值
}

type SdoOrderResponse struct {
	ReturnCode    int     `json:"return_code"`
	ReturnMessage string  `json:"return_message"`
	Data          SdoPage `json:"data"`
}

type SdoPage struct {
	Orders         []SdoOrderRecord `json:"orders"`
	TotalCount     int              `json:"totalCount"`
	PageCount      int              `json:"pageCount"`
	Version        string           `json:"version"`
	ExpiredMinutes int              `json:"expiredMinutes"`
}

const customTimeLayout = "2006-01-02 15:04:05"

func (o *SdoOrderRecord) UnmarshalJSON(data []byte) error {
	type Alias SdoOrderRecord
	aux := &struct {
		OrderTime string `json:"orderTime"`
		StateTime string `json:"stateTime"`
		*Alias
	}{
		Alias: (*Alias)(o),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	orderTime, err := time.Parse(customTimeLayout, aux.OrderTimeStr)
	if err != nil {
		return err
	}
	o.OrderTime = orderTime

	stateTime, err := time.Parse(customTimeLayout, aux.StateTimeStr)
	if err != nil {
		return err
	}
	o.StateTime = stateTime

	return nil
}

// ToJSON 方法将 SdoOrderRecord 结构体转换为 JSON 字符串
func (o *SdoOrderRecord) ToJSON() ([]byte, error) {
	type Alias SdoOrderRecord
	return json.Marshal(&struct {
		OrderTime string `json:"orderTime"`
		StateTime string `json:"stateTime"`
		*Alias
	}{
		OrderTime: o.OrderTime.Format(customTimeLayout),
		StateTime: o.StateTime.Format(customTimeLayout),
		Alias:     (*Alias)(o),
	})
}
