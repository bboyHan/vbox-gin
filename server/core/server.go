package core

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/initialize"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"go.uber.org/zap"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	if global.GVA_CONFIG.System.UseMultipoint || global.GVA_CONFIG.System.UseRedis {
		// 初始化redis服务
		initialize.Redis()
	}
	if global.GVA_CONFIG.System.UseMongo {
		err := initialize.Mongo.Initialization()
		if err != nil {
			zap.L().Error(fmt.Sprintf("%+v", err))
		}
	}
	if global.GVA_CONFIG.System.UseMQ {
		err := mq.MQ.Init()
		if err != nil {
			zap.L().Error(fmt.Sprintf("%+v", err))
		}
	}
	//MQ任务入口
	if global.GVA_CONFIG.System.UseMQTask {
		if global.GVA_CONFIG.MQTask.OrderWait {
			go task.OrderWaitingTask()
		}
		if global.GVA_CONFIG.MQTask.OrderConfirm {
			go task.OrderConfirmTask()
		}
		if global.GVA_CONFIG.MQTask.OrderCallback {
			go task.OrderCallbackTask()
		}
		if global.GVA_CONFIG.MQTask.AccEnableCheck {
			go task.ChanAccEnableCheckTask()
			go task.ChanAccDelCheckTask()
		}
		if global.GVA_CONFIG.MQTask.PayCodeExpCheck {
			go task.PayCodeExpCheck()
		}
		if global.GVA_CONFIG.MQTask.PayCodeCdCheck {
			go task.PayCodeCDCheckTask()
		}
		if global.GVA_CONFIG.MQTask.AccCDCheck {
			go task.AccCDCheckTask()
		}
	}

	//定时任务入口
	if global.GVA_CONFIG.System.TimerTask {
		initialize.Timer()
	}
	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := initialize.Routers()
	Router.Static("/form-generator", "./resource/page")

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
`, address)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}
