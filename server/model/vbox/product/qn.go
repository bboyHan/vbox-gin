package product

import "time"

type QNRecord struct {
	CreateTime  time.Time `json:"createTime"`
	OrderStatus string    `json:"orderStatus"`
	Money       string    `json:"money"`
	Buyer       string    `json:"buyer"`
	SkuTitle    string    `json:"skuTitle"`
}
