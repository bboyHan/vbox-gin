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

	var orderDBList []vbox.PayOrder
	global.GVA_DB.Model(&vbox.PayOrder{}).Table("vbox_pay_order").
		Where("order_status = ? and cb_status = ?", 1, 0).Find(&orderDBList)

	//global.GVA_LOG.Info("根据开启的商户列表，开始检测可用账号情况")

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

	return err
}
