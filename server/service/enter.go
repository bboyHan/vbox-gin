package service

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/channelshop"
	"github.com/flipped-aurora/gin-vue-admin/server/service/example"
	"github.com/flipped-aurora/gin-vue-admin/server/service/student"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
)

type ServiceGroup struct {
	SystemServiceGroup      system.ServiceGroup
	ExampleServiceGroup     example.ServiceGroup
	StudentServiceGroup     student.ServiceGroup
	ChannelshopServiceGroup channelshop.ServiceGroup
}

var ServiceGroupApp = new(ServiceGroup)
