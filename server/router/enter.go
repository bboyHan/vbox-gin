package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/student"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/vbox"
)

type RouterGroup struct {
	System  system.RouterGroup
	Example example.RouterGroup
	Student student.RouterGroup
	Vbox    vbox.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
