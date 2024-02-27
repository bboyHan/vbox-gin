package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
)

type ChanAccAndCtx struct {
	Obj vbox.ChannelAccount `json:"obj" form:"obj" url:"obj"`
	Ctx Context             `json:"ctx" form:"ctx" url:"ctx"`
}

type ChanCardAccAndCtx struct {
	Obj vbox.ChannelCardAcc `json:"obj" form:"obj" url:"obj"`
	Ctx Context             `json:"ctx" form:"ctx" url:"ctx"`
}

type PayOrderAndCtx struct {
	Obj vbox.PayOrder `json:"obj" form:"obj" url:"obj"`
	Ctx Context       `json:"ctx" form:"ctx" url:"ctx"`
}

type Context struct {
	ClientIP  string `json:"clientIP" form:"clientIP" url:"clientIP"`
	Method    string `json:"method" form:"method" url:"method"`
	UrlPath   string `json:"urlPath" form:"urlPath" url:"urlPath"`
	UserAgent string `json:"userAgent" form:"userAgent" url:"userAgent"`
	Body      string `json:"body" form:"body" url:"body"`
	RawQuery  string `json:"rawQuery" form:"rawQuery" url:"rawQuery"`
	UserID    int    `json:"userID" form:"userID" url:"userID"`
}
