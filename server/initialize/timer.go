package initialize

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/robfig/cron/v3"
)

func Timer() {
	go func() {
		defer func() { //保持程序运行态，而不是中断程序
			if err := recover(); err != nil {
				fmt.Println("Timer() 发生了 panic ex：", err)
				// 可以在这里进行一些处理操作
			}
		}()

		var option []cron.Option
		option = append(option, cron.WithSeconds())
		// 清理DB定时任务
		_, err := global.GVA_Timer.AddTaskByFunc("ClearDB", "@daily", func() {
			err := task.ClearTable(global.GVA_DB) // 定时任务方法定在task文件包中
			if err != nil {
				fmt.Println("timer error:", err)
			}
		}, option...)
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		// 其他定时任务定在这里 参考上方使用方法

		//_, err := global.GVA_Timer.AddTaskByFunc("定时任务标识", "corn表达式", func() {
		//	具体执行内容...
		//  ......
		//}, option...)
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}

		//_, err = global.GVA_Timer.AddTaskByFunc("handleAccPoolCheck", "@every 60s", func() {
		//	err = task.HandleAccPoolCheck()
		//	if err != nil {
		//		fmt.Println("timer error:", err)
		//	}
		//})
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}

		//_, err = global.GVA_Timer.AddTaskByFunc("handleEcJDCodeAdd", "@every 3s", func() {
		//	err = task.HandleEcJDCodeAdd()
		//	//if err != nil {
		//	//	fmt.Println("timer error:", err)
		//	//}
		//	err = task.HandleEcJDCodeDel()
		//	//if err != nil {
		//	//	fmt.Println("timer error:", err)
		//	//}
		//})
		//if err != nil {
		//	fmt.Println("add timer error:", err)
		//}

		_, err = global.GVA_Timer.AddTaskByFunc("handleOrderCallCheck", "@every 30s", func() {
			err = task.HandleOrderCallCheck()
			if err != nil {
				fmt.Println("timer error:", err)
			}
		})
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFunc("handleAccLimitCheck", "@every 5s", func() {
			err = task.HandleAccLimitCheck()
			if err != nil {
				fmt.Println("timer error:", err)
			}
		})
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFunc("handleShopMoneyAvailable", "@every 5s", func() {
			err = task.HandleShopMoneyAvailable()
			if err != nil {
				fmt.Println("timer error:", err)
			}
		})
		if err != nil {
			fmt.Println("add timer error:", err)
		}

		_, err = global.GVA_Timer.AddTaskByFunc("handleChannelStatisTask", "55 17 * * *", func() {
			err = task.HandleChannelStatisTask()
			if err != nil {
				fmt.Println("timer error:", err)
			}
		})
		if err != nil {
			fmt.Println("add timer error:", err)
		}
	}()

}
