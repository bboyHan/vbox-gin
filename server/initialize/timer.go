package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/service"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/robfig/cron/v3"
)

func Timer() {
	if global.GVA_CONFIG.Timer.Start {
		fmt.Println(global.GVA_CONFIG.Timer.Detail)
		for i := range global.GVA_CONFIG.Timer.Detail {
			go func(detail config.Detail) {
				var option []cron.Option
				if global.GVA_CONFIG.Timer.WithSeconds {
					option = append(option, cron.WithSeconds())
				}
				_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", global.GVA_CONFIG.Timer.Spec, func() {
					err := utils.ClearTable(global.GVA_DB, detail.TableName, detail.CompareField, detail.Interval)
					if err != nil {
						fmt.Println("timer error:", err)
					}
				}, option...)
				if err != nil {
					fmt.Println("add timer error:", err)
				}
			}(global.GVA_CONFIG.Timer.Detail[i])
		}
	}
}

var vuwService = service.ServiceGroupApp.VboxServiceGroup.VboxUserWalletService

func MockFunc() {
	//time.Sleep(time.Second)
	//fmt.Println("1s...")
	//log.Printf("1s...")
	//func() {
	// 在匿名函数中调用带参数的函数
	vuwService.GetVboxUserWalletAvailablePoints(3)
	//}
}

func CronTableAnalysisTimer() {
	if global.GVA_CONFIG.Timer.Start {
		//fmt.Println(global.GVA_CONFIG.Timer.Detail)
		var option []cron.Option
		if global.GVA_CONFIG.Timer.WithSeconds {
			option = append(option, cron.WithSeconds())
		}
		_, err := global.GVA_Timer.AddTaskByFunc("payOrderAnalysisTask", "@every 10s", func() {
			MockFunc()
		}, option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}
	}
}
