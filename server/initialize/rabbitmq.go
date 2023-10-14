package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
)

func MQ() {
	//mqCfg := global.GVA_CONFIG.RabbitMQ
	//conn, err := utils.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/",
	//	mqCfg.Username,
	//	mqCfg.Password,
	//	mqCfg.Addr,
	//	mqCfg.Port))
	//if err != nil {
	//	global.GVA_LOG.Error("new mq conn err:", zap.Error(err))
	//} else {
	//	global.GVA_MQ_CONN = conn
	//	channel, err := conn.Channel()
	//	global.GVA_MQ_CHANNEL = channel
	//	if err != nil {
	//		global.GVA_LOG.Error("new mq conn err:", zap.Error(err))
	//	}
	//	global.GVA_LOG.Info("mq connect success")
	//}
	utils.Init()
}
