package task

import (
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"go.uber.org/zap"
)

func HandleOrderCallCheck() (err error) {
	//
	//// 获取今天的时间范围
	//defaultOfDay := time.Now().UTC().Truncate(24 * time.Hour)
	//startOfDay := defaultOfDay.Add(24 * time.Hour)
	//endOfDay := defaultOfDay.Add(-24 * time.Hour)
	//// 获取本地时区
	//loc, _ := time.LoadLocation("Asia/Shanghai") // 请替换为你实际使用的时区
	//startOfDay = startOfDay.In(loc)
	//endOfDay = endOfDay.In(loc)

	var orderDBList []vbox.PayOrder
	global.GVA_DB.Model(&vbox.PayOrder{}).Table("vbox_pay_order").
		Where("order_status = ? and cb_status in (0,2)", 1).Find(&orderDBList)

	//if len(orderDBList) == 0 {
	//	return nil
	//} else {
	//	//查出来有订单已支付，未回调
	//	global.GVA_LOG.Info("查出来有订单已支付，未回调，当前情况", zap.Any("len", len(orderDBList)))
	//	global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Table("vbox_pay_order").
	//		Where("order_status = ? and cb_status = ?", 1, 2).Find(&orderDBList)
	//
	//	global.GVA_LOG.Info("查出来有订单已支付，未回调，再debug查一遍看看", zap.Any("len", len(orderDBList)))
	//}

	for _, orderDB := range orderDBList {
		orderDBTmp := orderDB
		go func() {
			//查一下订单是否超出账户限制
			v := request.PayOrderAndCtx{
				Obj: orderDBTmp,
				Ctx: vboxReq.Context{
					Body:      "系统补单",
					ClientIP:  "127.0.0.1",
					Method:    "POST",
					UrlPath:   "/payOrder/callback2Pa",
					UserAgent: "",
					UserID:    int(orderDBTmp.CreatedBy),
				},
			}

			conn, errC := mq.MQ.ConnPool.GetConnection()
			if errC != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)
			ch, errN := conn.Channel()
			if errN != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
			}

			// 并且发起一个回调通知的消息
			marshal, _ := json.Marshal(v)
			err = ch.Publish(task.OrderCallbackExchange, task.OrderCallbackKey, marshal)
			global.GVA_LOG.Info("【系统补单】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))

		}()
	}

	var orderDBFixList []vbox.PayOrder
	global.GVA_DB.Model(&vbox.PayOrder{}).Table("vbox_pay_order").
		Where("order_status != ? and cb_status = ?", 1, 1).Find(&orderDBFixList)
	for _, orderDB := range orderDBFixList {
		orderDBTmp := orderDB
		go func() {
			global.GVA_DB.Debug().Model(&vbox.PayOrder{}).Where("id = ?", orderDBTmp.ID).
				Update("order_status", 1).Update("hand_status", 1)
			global.GVA_LOG.Info("【系统修复】更新已回调确显示未支付的订单", zap.Any("order ID", orderDBTmp.OrderId))
		}()
	}
	return err
}
