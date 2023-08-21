package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/channelshop"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/student"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
)

type RouterGroup struct {
	System      system.RouterGroup
	Example     example.RouterGroup
	Student     student.RouterGroup
	Channelshop channelshop.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
