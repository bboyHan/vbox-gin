package product

import "time"

//duration time.Duration, cdTime time.Duration, preTime time.Duration

type CDTime struct {
	Duration time.Duration `json:"duration"`
	CDTime   time.Duration `json:"cdTime"`
	PreTime  time.Duration `json:"preTime"`
}
