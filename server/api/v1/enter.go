package v1

import (
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/channelshop"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/example"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/student"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/system"
)

type ApiGroup struct {
	SystemApiGroup      system.ApiGroup
	ExampleApiGroup     example.ApiGroup
	StudentApiGroup     student.ApiGroup
	ChannelshopApiGroup channelshop.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
