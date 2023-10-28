package initialize

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/task"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"log"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/global"
)

func Timer() {
	go func() {
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

		_, err := global.GVA_Timer.AddTaskByFunc("listAccountAvailable", "@every 5s", func() {
			//自定义 通道账号可用检测
			rdConn := global.GVA_REDIS.Conn()
			defer rdConn.Close()
			var idList []uint
			// 拿出现在所有付方可用的账户
			err := global.GVA_DB.Debug().Model(&vbox.PayAccount{}).Table("vbox_pay_account").
				Select("uid").Where("status = ?", 1).Find(&idList).Error
			if err != nil {
				global.GVA_LOG.Error("查付方库数据异常", zap.Error(err))
				return
			}

			global.GVA_LOG.Info("我开始检测有没有可用的账号了")

			for _, uid := range idList {
				var channelCodeList []string

				orgIds := utils2.GetDeepOrg(uid)
				c, err := rdConn.Exists(context.Background(), global.UserOrgChannelCodePrefix+strconv.FormatUint(uint64(uid), 10)).Result()
				if c == 0 {
					var productIds []uint
					if err != nil {
						global.GVA_LOG.Error("当前缓存池无此用户对应的orgIds，redis err", zap.Error(err))
					}
					db := global.GVA_DB.Model(&vbox.OrgProduct{})
					if err = db.Debug().Distinct("channel_product_id").Select("channel_product_id").Where("organization_id in ?", orgIds).Find(&productIds).Error; err != nil {
						global.GVA_LOG.Error("OrgProduct查该组织下数据channel code异常", zap.Error(err))
						continue
					}
					if err = db.Debug().Model(&vbox.ChannelProduct{}).Select("channel_code").Where("id in ?", productIds).Find(&channelCodeList).Error; err != nil {
						global.GVA_LOG.Error("ChannelProduct查channelCodeList 库数据异常", zap.Error(err))
						continue
					}

					jsonStr, _ := json.Marshal(channelCodeList)
					rdConn.Set(context.Background(), global.UserOrgChannelCodePrefix+strconv.FormatUint(uint64(uid), 10), jsonStr, 10*time.Minute)
				} else {
					jsonStr, _ := rdConn.Get(context.Background(), global.UserOrgChannelCodePrefix+strconv.FormatUint(uint64(uid), 10)).Bytes()
					err = json.Unmarshal(jsonStr, &channelCodeList)
				}

				for _, channelCode := range channelCodeList {
					var total int64
					deepUIDs := utils2.GetDeepUserIDs(uid)
					db := global.GVA_DB.Debug().Model(&vbox.ChannelAccount{}).Table("vbox_channel_account").
						Where("created_by in ?", deepUIDs).Count(&total)

					log.Printf("查出来总号有 %d 个", total)
					limit, offset := utils.RandSize2DB(int(total), 20)
					var vcas []vbox.ChannelAccount
					err = db.Debug().Where("status = ? and sys_status = ?", 1, 1).Where("cid = ?", channelCode).
						Where("created_by in (?)", deepUIDs).Limit(limit).Offset(offset).
						Find(&vcas).Error
					if err != nil || len(vcas) == 0 {
						if len(vcas) == 0 {
							global.GVA_LOG.Error("ChannelAccount查数据，没号不管", zap.Error(err))
							continue
						}
					}

					global.GVA_LOG.Info("vca available", zap.String("channel code", channelCode), zap.Any("可用数", len(vcas)), zap.Any("list", vcas))
				}
			}

		})
		if err != nil {
			fmt.Println("add timer error:", err)
		}
	}()
}
