package vbUtil

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"go.uber.org/zap"
)

// CheckAccLimit 查询账号限额情况
func CheckAccLimit(vca vbox.ChannelAccount, money int) (err error) {

	global.GVA_LOG.Info("当前账号限制情况", zap.Any("dailyLimit", vca.DailyLimit), zap.Any("totalLimit", vca.TotalLimit),
		zap.Any("inLimit", vca.InCntLimit), zap.Any("cntLimit", vca.CountLimit), zap.Any("当前请求金额", money))

	if vca.TotalLimit > 0 {
		var totalSum int

		err = global.GVA_DB.Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as totalSum").
			Where("ac_id = ?", vca.AcId).
			Where("channel_code = ?", vca.Cid).
			Where("order_status = ?", 1).Scan(&totalSum).Error

		if totalSum+money > vca.TotalLimit {
			//global.GVA_LOG.Error("总限额不足", zap.Any("orderId", v.Obj.OrderId), zap.Any("acID", accID), zap.Any("acAccount", vca.AcAccount), zap.Any("totalSum", totalSum), zap.Any("money", money), zap.Any("TotalLimit", vca.TotalLimit))
			return fmt.Errorf("总限额不足, 限额参数:%v, 当前消费:%v, 请求金额:%v", vca.TotalLimit, totalSum, money)
		}
	}
	if vca.DailyLimit > 0 {
		var dailySum int

		err = global.GVA_DB.Model(&vbox.PayOrder{}).Select("IFNULL(sum(money), 0) as dailySum").
			Where("ac_id = ?", vca.AcId).
			Where("channel_code = ?", vca.Cid).
			Where("order_status = ? AND created_at BETWEEN CURDATE() AND CURDATE() + INTERVAL 1 DAY - INTERVAL 1 SECOND", 1).Scan(&dailySum).Error

		if dailySum+money > vca.DailyLimit {
			//global.GVA_LOG.Error("每日限额不足", zap.Any("acID", accID), zap.Any("acAccount", vca.AcAccount), zap.Any("dailySum", dailySum), zap.Any("money", money), zap.Any("dailyLimit", vca.DailyLimit))
			return fmt.Errorf("每日限额不足, 限额参数:%v, 当前消费:%v, 请求金额:%v", vca.DailyLimit, dailySum, money)
		}
	}
	if vca.InCntLimit > 0 {
		var inCnt int
		err = global.GVA_DB.Model(&vbox.PayOrder{}).Select("count(1) as inCnt").
			Where("ac_id = ?", vca.AcId).
			Where("channel_code = ?", vca.Cid).
			Where("order_status = ?", 1).Scan(&inCnt).Error

		if inCnt+1 > vca.InCntLimit {
			//global.GVA_LOG.Error("入单数限额不足", zap.Any("orderId", v.Obj.OrderId), zap.Any("acID", accID), zap.Any("acAccount", vca.AcAccount), zap.Any("inCnt", inCnt), zap.Any("money", money), zap.Any("inCntLimit", vca.InCntLimit))
			return fmt.Errorf("入单数限额不足, 限额参数:%v, 当前笔数:%v", vca.InCntLimit, inCnt)
		}
	}
	if vca.DlyCntLimit > 0 {
		var dlyCnt int
		err = global.GVA_DB.Model(&vbox.PayOrder{}).Select("count(1) as dlyCnt").
			Where("ac_id = ?", vca.AcId).
			Where("channel_code = ?", vca.Cid).
			Where("order_status = ? AND created_at BETWEEN CURDATE() AND CURDATE() + INTERVAL 1 DAY - INTERVAL 1 SECOND", 1).Scan(&dlyCnt).Error
		if dlyCnt+1 > vca.DlyCntLimit {
			//global.GVA_LOG.Error("每日入单数限额不足", zap.Any("acID", accID), zap.Any("acAccount", vca.AcAccount), zap.Any("dlyCnt", dlyCnt), zap.Any("money", money), zap.Any("dlyCntLimit", vca.DlyCntLimit))
			return fmt.Errorf("每日入单数限额不足, 限额参数:%v, 当前笔数:%v", vca.DlyCntLimit, dlyCnt)
		}
	}
	if vca.CountLimit > 0 {
		var cnt int
		err = global.GVA_DB.Model(&vbox.PayOrder{}).Select("count(1) as cnt").
			Where("ac_id = ?", vca.AcId).
			Where("channel_code = ?", vca.Cid).
			Scan(&cnt).Error

		if cnt+1 > vca.CountLimit {
			//global.GVA_LOG.Error("总拉单限额不足", zap.Any())
			return fmt.Errorf("总拉单限额不足, 限额参数:%v, 当前笔数:%v", vca.CountLimit, cnt)
		}
	}
	return nil
}
