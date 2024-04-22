package telegram_bot

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/telegram_bot/router"
	"github.com/gin-gonic/gin"
)

type Telegram_botPlugin struct{}

func CreateTelegram_botPlug() *Telegram_botPlugin {
	return &Telegram_botPlugin{}
}

func (*Telegram_botPlugin) Register(group *gin.RouterGroup) {
	router.RouterGroupApp.InitTelegram_botRouter(group)
}

func (*Telegram_botPlugin) RouterPath() string {
	return "telegram_bot"
}
