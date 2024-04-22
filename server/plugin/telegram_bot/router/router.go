package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/telegram_bot/api"
	"github.com/gin-gonic/gin"
)

type Telegram_botRouter struct{}

func (s *Telegram_botRouter) InitTelegram_botRouter(Router *gin.RouterGroup) {
	plugRouter := Router
	plugApi := api.ApiGroupApp.Telegram_botApi
	{
		plugRouter.POST("sendMessage", plugApi.SendMessage)
		plugRouter.POST("isMember", plugApi.IsMember)
	}
}
