package vbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxRep "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/system"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/product"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	http2 "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"github.com/flipped-aurora/gin-vue-admin/server/vbUtil"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type PayOrderService struct {
}

var operationRecordService system.OperationRecordService

func (vpoService *PayOrderService) QryOrderDataOverview(info vboxReq.PayOrderSearch, ids []uint) (ov []vboxRep.DataWalletOverView, err error) {
	db := global.GVA_DB.Table("vbox_pay_order").Model(&vboxRep.DataWalletOverView{})
	if info.ChannelCode != "" {
		db = db.Where("channel_code = ?", info.ChannelCode)
	}
	if info.PAccount != "" {
		db = db.Where("p_account =?", info.PAccount)
	}

	err = db.Select(
		`IFNULL(SUM(CASE WHEN order_status = 1 THEN money ELSE 0 END), 0) AS x0,
	IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() - INTERVAL 2 DAY AND order_status = 1 THEN money ELSE 0 END), 0) AS x1,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() - INTERVAL 1 DAY AND order_status = 1 THEN money ELSE 0 END), 0) AS x2,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() AND order_status = 1 THEN money ELSE 0 END), 0) AS x3,
    IFNULL(SUM(CASE WHEN cb_time >= NOW() - INTERVAL 1 HOUR AND order_status = 1 THEN money ELSE 0 END), 0) AS x4`).
		Where("created_by in ?", ids).Find(&ov).Error
	if err != nil {
		return
	}

	return ov, nil
}

func (vpoService *PayOrderService) QryOrderAccOverview(info vboxReq.PayOrderSearch, ids []uint) (ov []vboxRep.OrderAccRes, err error) {
	var acIDs []uint
	db := global.GVA_DB.Table("vbox_pay_order").Model(&vboxRep.OrderAccRes{})
	if info.ChannelCode != "" {
		db = db.Where("channel_code =?", info.ChannelCode)
	}
	if info.AcRemark != "" {
		global.GVA_DB.Debug().Model(&vbox.ChannelAccount{}).Select("ac_id").Where("ac_remark like ?", "%"+info.AcRemark+"%").Scan(&acIDs)
	}
	if info.AcId != "" && len(acIDs) == 0 {
		db = db.Where("ac_id = ?", info.AcId)
	} else if len(acIDs) > 0 {
		db = db.Where("ac_id in ?", acIDs)
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account =?", info.AcAccount)
	}
	if info.ToUid != 0 {
		err = db.Select(
			`created_by,ac_id,ac_account,channel_code,
    IFNULL(SUM(CASE WHEN order_status = 1 THEN money ELSE 0 END), 0) AS x0,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() - INTERVAL 3 DAY AND order_status = 1 THEN money ELSE 0 END), 0) AS x1,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() - INTERVAL 2 DAY AND order_status = 1 THEN money ELSE 0 END), 0) AS x2,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() - INTERVAL 1 DAY AND order_status = 1  THEN money ELSE 0 END), 0) AS x3,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() AND order_status = 1 THEN money ELSE 0 END), 0) AS x4`).
			Where("created_by = ?", info.ToUid).Group("created_by, ac_id ,channel_code").Find(&ov).Error
		if err != nil {
			return
		}
	} else {
		err = db.Select(
			`created_by,ac_id,ac_account,channel_code,
    IFNULL(SUM(CASE WHEN order_status = 1  THEN money ELSE 0 END), 0) AS x0,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() - INTERVAL 3 DAY AND order_status = 1 THEN money ELSE 0 END), 0) AS x1,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() - INTERVAL 2 DAY AND order_status = 1 THEN money ELSE 0 END), 0) AS x2,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() - INTERVAL 1 DAY AND order_status = 1  THEN money ELSE 0 END), 0) AS x3,
    IFNULL(SUM(CASE WHEN DATE(cb_time) = CURDATE() AND order_status = 1 THEN money ELSE 0 END), 0) AS x4`).
			Where("created_by in ?", ids).Group("created_by, ac_id ,channel_code").Find(&ov).Error
		if err != nil {
			return
		}
	}

	return ov, nil
}

//	p := &vboxReq.CallBackExtReq{
//			OrderId:        "123",
//			Ext:        "123",
//		}
func (vpoService *PayOrderService) CallbackOrderExt(vpo *vboxReq.CallBackExtReq, c *gin.Context) (rep *vboxRep.OrderSimpleRes, err error) {
	var order vbox.PayOrder
	// 校验传入卡密合法性
	if vpo.ChannelCode == "1101" {
		if _, errX := product.ParseJWCardRecord(vpo.Ext); errX != nil {
			return nil, errX
		}
	} else if vpo.ChannelCode == "6001" {
	} else {
		return nil, errors.New("该订单类型，不支持卡密提交，请联系管理员")
	}

	var jsonString []byte
	key := fmt.Sprintf(global.PayOrderKey, vpo.OrderId)
	rdRes, err := global.GVA_REDIS.Get(context.Background(), key).Bytes()
	if err == redis.Nil { // redis中还没有的情况，查一下库，并且去匹配设备信息
		err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("order_id = ?", vpo.OrderId).First(&order).Error
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		fmt.Println("error:", err)
	} else {
		//fmt.Println("从缓存里拿result:", rdRes)
		err = json.Unmarshal(rdRes, &order)
	}

	//设置提交次数限制
	limitKey := fmt.Sprintf(global.PayOrderExtLimitKey, order.OrderId)
	limitCnt, err := vbUtil.SetLimitWithTime(limitKey, 3, 8*time.Minute)
	if err != nil {
		record := sysModel.SysOperationRecord{
			Ip:      c.ClientIP(),
			Method:  c.Request.Method,
			Path:    c.Request.URL.Path,
			Agent:   c.Request.UserAgent(),
			MarkId:  fmt.Sprintf(global.OrderRecord, order.OrderId),
			Type:    global.OrderType,
			Status:  500,
			Latency: time.Since(time.Now()),
			Resp:    fmt.Sprintf(global.OrderConfirmBindLimitErrMsg),
			UserID:  int(order.CreatedBy),
		}
		err = operationRecordService.CreateSysOperationRecord(record)

		return nil, errors.New("您提交的错误次数过多，请重新下单")
	}

	if vpo.Ext == order.PlatId {
		return nil, errors.New(fmt.Sprintf("您已提交过输入的卡密，无法重复提交，当前已提交%d次", limitCnt))
	}

	//校验卡密
	if vpo.ChannelCode == "6001" {
		//1.0 核查商户
		var vpa vbox.PayAccount
		if strings.Contains(order.PAccount, "TEST") {
			global.GVA_LOG.Info("测试单，商户检测跳过", zap.Any("入参商户", order.PAccount))
			vpa = vbox.PayAccount{
				PAccount: order.PAccount,
				Uid:      order.CreatedBy,
			}
		} else {
			var count int64
			count, err = global.GVA_REDIS.Exists(context.Background(), global.PayAccPrefix+order.PAccount).Result()
			if count == 0 {
				if err != nil {
					global.GVA_LOG.Error("当前缓存池无此商户，redis err", zap.Error(err))
				}
				global.GVA_LOG.Info("当前缓存池无此商户，查一下库。。。", zap.Any("入参商户ID", order.PAccount))

				err = global.GVA_DB.Table("vbox_pay_account").
					Where("p_account = ?", order.PAccount).First(&vpa).Error
				jsonStr, _ := json.Marshal(vpa)
				global.GVA_REDIS.Set(context.Background(), global.PayAccPrefix+order.PAccount, jsonStr, 10*time.Minute)
			} else {
				jsonStr, _ := global.GVA_REDIS.Get(context.Background(), global.PayAccPrefix+order.PAccount).Bytes()
				err = json.Unmarshal(jsonStr, &vpa)
			}
			if err != nil {
				global.GVA_LOG.Error("订单匹配消息处理失败，MqOrderWaitingTask...", zap.Any("err", err.Error()))
			}
		}
		orgTmp := utils2.GetSelfOrg(vpa.Uid)

		accPoolKey := fmt.Sprintf(global.ChanOrgECPoolAccZSet, orgTmp[0], "6001")

		var resPoolList []string
		resPoolList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), accPoolKey, &redis.ZRangeBy{
			Min:    "0",
			Max:    "0",
			Offset: 0,
			Count:  -1,
		}).Result()

		if err != nil {
			global.GVA_LOG.Error("卡密类匹配查单池redis异常, redis err", zap.Error(err))
		}
		var cardID, cardAccID, cardAcAccount string

		var cardAcc vbox.ChannelCardAcc

		if resPoolList != nil && len(resPoolList) > 0 {
			//accTmp := resList[0]
			accPoolTmp := utils.RandomElement(resPoolList)
			// 2.1 把账号设置为已用
			global.GVA_REDIS.ZAdd(context.Background(), accPoolKey, redis.Z{
				Score:  1,
				Member: accPoolTmp,
			})
			split := strings.Split(accPoolTmp, ",")
			cardID = split[0]
			cardAccID = split[1]
			cardAcAccount = split[2]
			err = global.GVA_DB.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", cardID).First(&cardAcc).Error
			if err != nil {
				//return nil, fmt.Errorf("匹配通道账号不存在！ 请核查：%v", err.Error())

				global.GVA_LOG.Error("匹配查单池账号异常", zap.Error(err))
			}
			global.GVA_LOG.Info("匹配查单池账号ck", zap.Any("cardID", cardID), zap.Any("cardAccID", cardAccID), zap.Any("cardAcAccount", cardAcAccount))

			cdTime := 10 * time.Second
			accWaitYdKey := fmt.Sprintf(global.YdECPoolAccWaiting, cardAccID)
			accInfoVal := fmt.Sprintf("%s,%s,%s", cardID, cardAccID, cardAcAccount)

			// 设置一个冷却时间
			ttl := global.GVA_REDIS.TTL(context.Background(), accWaitYdKey).Val()
			if ttl > 0 {
				cdTime = ttl
				global.GVA_LOG.Info("当前添加的账号正在冷却中（有预产正在处理中）", zap.Any("accWaitYdKey", accWaitYdKey), zap.Any("ttl", cdTime))
			} else {
				cdTime += 2 * time.Second
				global.GVA_LOG.Info("当前添加的账号新一轮冷却", zap.Any("accWaitYdKey", accWaitYdKey), zap.Any("ttl", cdTime))
			}
			global.GVA_REDIS.Set(context.Background(), accWaitYdKey, accInfoVal, cdTime)

			waitMsg := strings.Join([]string{accWaitYdKey, accInfoVal}, "-")

			connX, errX := mq.MQ.ConnPool.GetConnection()
			if errX != nil {
				//log.Fatalf("Failed to get connection from pool: %v", err)
				global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(errX))
			}
			defer mq.MQ.ConnPool.ReturnConnection(connX)
			chX, _ := connX.Channel()
			err = chX.PublishWithDelay(task.CardAccCDCheckDelayedExchange, task.CardAccCDCheckDelayedRoutingKey, []byte(waitMsg), 0)

		} else {
			global.GVA_LOG.Error("查单池匹配账号CK不足, list size zero", zap.Error(err), zap.Any("channelCode", order.ChannelCode), zap.Any("money", order.Money))
			record := sysModel.SysOperationRecord{
				Ip:      c.ClientIP(),
				Method:  c.Request.Method,
				Path:    c.Request.URL.Path,
				Agent:   c.Request.UserAgent(),
				MarkId:  fmt.Sprintf(global.OrderRecord, order.OrderId),
				Type:    global.OrderType,
				Status:  500,
				Latency: time.Since(time.Now()),
				Resp:    fmt.Sprintf(global.OrderConfirmBindPoolErrMsg),
				UserID:  int(order.CreatedBy),
			}
			err = operationRecordService.CreateSysOperationRecord(record)
		}

		//a.拿账户的ck查一遍卡密，如果卡密ok则执行绑卡

		errE := product.ECardQuery(vpo.Ext, cardAcc.Token)
		if errE != nil {

			global.GVA_LOG.Error("匹配查单池账号异常,查询卡密合法性错误", zap.Error(errE))

			record := sysModel.SysOperationRecord{
				Ip:      c.ClientIP(),
				Method:  c.Request.Method,
				Path:    c.Request.URL.Path,
				Agent:   c.Request.UserAgent(),
				MarkId:  fmt.Sprintf(global.OrderRecord, order.OrderId),
				Type:    global.OrderType,
				Status:  500,
				Latency: time.Since(time.Now()),
				Resp:    fmt.Sprintf(global.OrderConfirmBindErrMsg, limitCnt, vpo.Ext, errE.Error()),
				UserID:  int(order.CreatedBy),
			}
			err = operationRecordService.CreateSysOperationRecord(record)

			return nil, errE
		}
	}

	//if order.PlatId == "" { //
	global.GVA_LOG.Info("传入卡密", zap.Any("ext", vpo.Ext))

	err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("order_id = ?", vpo.OrderId).Update("plat_id", vpo.Ext).Error
	if err != nil {
		return nil, err
	}
	order.PlatId = vpo.Ext
	//查出来了，设置一下redis
	jsonString, err = json.Marshal(order)
	if err != nil {
		return nil, err
	}
	global.GVA_REDIS.Set(context.Background(), key, jsonString, 2*time.Second)
	//} else {
	//	return nil, errors.New("该订单已提交过卡密，无法重复提交")
	//}

	var ext string
	if order.PlatId != "" {
		ext = fmt.Sprintf("cnt_%d,", limitCnt)
	}

	rep = &vboxRep.OrderSimpleRes{
		OrderId:     order.OrderId,
		Account:     order.AcAccount,
		Money:       order.Money,
		ResourceUrl: order.ResourceUrl,
		Status:      order.OrderStatus,
		ExpTime:     *order.ExpTime,
		Ext:         ext,
		Cnt:         limitCnt,
		ChannelCode: order.ChannelCode,
	}

	return rep, err
}

// QueryOrderSimple 查询QueryOrderSimple
//
//	p := &vboxReq.QueryOrderSimple{
//			OrderId:        "123",
//		}
func (vpoService *PayOrderService) QueryOrderSimple(vpo *vboxReq.QueryOrderSimple, ctx *gin.Context) (rep *vboxRep.OrderSimpleRes, err error) {

	// 1. 查单
	var order vbox.PayOrder
	var jsonString []byte
	key := fmt.Sprintf(global.PayOrderKey, vpo.OrderId)
	rdRes, err := global.GVA_REDIS.Get(context.Background(), key).Bytes()
	if err == redis.Nil { // redis中还没有的情况，查一下库，并且去匹配设备信息
		fmt.Println("redis key does not exist")
		err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("order_id = ?", vpo.OrderId).First(&order).Error
		if err != nil {
			return nil, err
		}

		if order.PayIp != "" {

		} else {
			// 算他第一次点开进行匹配

			err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("order_id = ?", vpo.OrderId).
				Update("pay_ip", vpo.PayIp).Update("pay_region", vpo.PayRegion).Update("pay_device", vpo.PayDevice).Error
			if err != nil {
				return nil, err
			}

			//
			conn, errX := mq.MQ.ConnPool.GetConnection()
			if errX != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)

			ch, errC := conn.Channel()
			if errC != nil || ch == nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
				//重试一次
				conn, errX = mq.MQ.ConnPool.GetConnection()
				if errX != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
				}
				ch, errC = conn.Channel()
				if errC != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
				}
				global.GVA_LOG.Warn("cn == nil 重试一次获取")
				if ch == nil {
					global.GVA_LOG.Warn("cn tnn 还是没取到")
					return nil, err
				}
			}

			body := http2.DoGinContextBody(ctx)

			err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("order_id = ?", vpo.OrderId).First(&order).Error
			if err != nil {
				return nil, err
			}

			od := vboxReq.PayOrderAndCtx{
				Obj: order,
				Ctx: vboxReq.Context{
					Body:      string(body),
					ClientIP:  ctx.ClientIP(),
					Method:    ctx.Request.Method,
					UrlPath:   ctx.Request.URL.Path,
					UserAgent: ctx.Request.UserAgent(),
					UserID:    int(order.CreatedBy),
				},
			}
			if ch == nil {
				//重试一次
				conn, errX = mq.MQ.ConnPool.GetConnection()
				if errX != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
				}
				defer mq.MQ.ConnPool.ReturnConnection(conn)
				ch, errC = conn.Channel()
				if errC != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
				}
				global.GVA_LOG.Warn("cn == nil 重试一次获取")

			}

			marshal, _ := json.Marshal(od)
			err = ch.Publish(task.OrderWaitExchange, task.OrderWaitKey, marshal)
			global.GVA_LOG.Info("发起一条资源匹配消息并入库初始化订单数据", zap.Any("od", od))
		}

		if order.ResourceUrl != "" {
			//查出来了，设置一下redis
			jsonString, err = json.Marshal(order)
			if err != nil {
				return nil, err
			}
			global.GVA_REDIS.Set(context.Background(), key, jsonString, 2*time.Second)
		}

	} else if err != nil {
		fmt.Println("error:", err)
	} else {
		//fmt.Println("从缓存里拿result:", rdRes)
		err = json.Unmarshal(rdRes, &order)
	}

	var ext string
	if order.PlatId != "" {
		ext = "_"
	}

	rep = &vboxRep.OrderSimpleRes{
		OrderId:     order.OrderId,
		Account:     order.AcAccount,
		Money:       order.Money,
		ResourceUrl: order.ResourceUrl,
		Status:      order.OrderStatus,
		ExpTime:     *order.ExpTime,
		Ext:         ext,
		ChannelCode: order.ChannelCode,
	}

	return rep, err
}

// CallbackOrder2PayAcc 补单 2 PayAcc
func (vpoService *PayOrderService) CallbackOrder2PayAcc(orderID string, ctx *gin.Context) (err error) {
	var order vbox.PayOrder

	//1.0 查单
	odKey := fmt.Sprintf(global.PayOrderKey, orderID)
	rdRes, err := global.GVA_REDIS.Get(context.Background(), odKey).Bytes()
	if err == redis.Nil { // redis中还没有的情况，查一下库
		//	查一下库
		err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("order_id =?", orderID).First(&order).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.GVA_LOG.Warn("订单不存在", zap.Any("当前订单", orderID))
			return errors.New("订单不存在")
		} else if err != nil {
			global.GVA_LOG.Error("订单不存在", zap.Any("当前订单", orderID), zap.Error(err))
			return err
		}

	} else if err != nil {
		global.GVA_LOG.Error("redis err", zap.Error(err))
	} else {
		//fmt.Println("从缓存里拿result:", rdRes)
		err = json.Unmarshal(rdRes, &order)
	}

	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		//log.Fatalf("Failed to get connection from pool: %v", err)
		global.GVA_LOG.Error("Failed to get connection from pool", zap.Error(err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	// ------------- 订单查询（30s延迟，第一次查单） 消息处理 --------------------
	ch, _ := conn.Channel()

	//	看一下订单状态，如果是已支付（order_status=1）,则简单发起http再回调一次（callback mq任务）
	if order.OrderStatus == 1 {
		global.GVA_LOG.Info("【候补单】订单处于已支付状态", zap.Any("order ID", orderID))

	} else if order.OrderStatus == 2 {
		global.GVA_LOG.Info("【候补单】订单处于待支付状态", zap.Any("order ID", orderID))

		//	更新状态为等待补单
		order.HandStatus = 3

	} else if order.OrderStatus == 3 {
		global.GVA_LOG.Info("【候补单】订单处于超时失效状态", zap.Any("order ID", orderID))

		//	更新状态为等待补单
		order.HandStatus = 3

	}

	body := http2.DoGinContextBody(ctx)

	v := vboxReq.PayOrderAndCtx{
		Obj: order,
		Ctx: vboxReq.Context{
			Body:      string(body),
			ClientIP:  ctx.ClientIP(),
			Method:    ctx.Request.Method,
			UrlPath:   ctx.Request.URL.Path,
			UserAgent: ctx.Request.UserAgent(),
			UserID:    int(order.CreatedBy),
		},
	}

	// 并且发起一个回调通知的消息
	marshal, _ := json.Marshal(v)
	err = ch.Publish(task.OrderCallbackExchange, task.OrderCallbackKey, marshal)
	global.GVA_LOG.Info("【候补单】发起一条回调消息等待处理", zap.Any("pa", v.Obj.PAccount), zap.Any("order ID", v.Obj.OrderId))

	return err
}

// QueryOrder2PayAcc 查询QueryOrder2PayAcc
//
//	p := &vboxReq.QueryOrder2PayAccount{
//			Account:     "",
//			Key:         "",
//			Sign:        "123",
//		}
func (vpoService *PayOrderService) QueryOrder2PayAcc(vpo *vboxReq.QueryOrder2PayAccount) (rep *vboxRep.OrderSimple2PayAccountRes, err error) {
	// 1. 校验签名
	var vpa vbox.PayAccount
	count, err := global.GVA_REDIS.Exists(context.Background(), global.PayAccPrefix+vpo.Account).Result()
	if count != 0 {
		global.GVA_LOG.Warn("缓存中暂无", zap.Any("当前 pacc", vpo.Account))
		jsonStr, _ := global.GVA_REDIS.Get(context.Background(), global.PayAccPrefix+vpo.Account).Bytes()
		err = json.Unmarshal(jsonStr, &vpa)
	} else { //查库看有没有
		err = global.GVA_DB.Model(&vbox.PayAccount{}).Table("vbox_pay_account").
			Where("p_account = ?", vpo.Account).Find(&vpa).Error
		if err != nil {
			return nil, err
		} else { //有的话，更新一下redis
			jsonStr, _ := json.Marshal(vpa)
			global.GVA_REDIS.Set(context.Background(), global.PayAccPrefix+vpa.PAccount, jsonStr, 10*time.Minute)
		}
	}

	vpo.Key = vpa.PKey
	signValid := utils.VerifySign(vpo)
	if !signValid {
		return nil, errors.New("请求参数或签名值不正确，请联系管理员核对")
	}

	// 2. 查单
	key := fmt.Sprintf(global.PayOrderKey, vpo.OrderId)
	rdRes, err := global.GVA_REDIS.Get(context.Background(), key).Bytes()
	var order vbox.PayOrder
	var jsonString []byte
	if err == redis.Nil { // redis中还没有的情况，查一下库，并且去匹配设备信息
		err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("order_id = ?", vpo.OrderId).First(&order).Error
		if err != nil {
			return nil, err
		}
		//查出来了，设置一下redis
		jsonString, err = json.Marshal(order)
		if err != nil {
			return nil, err
		}
		global.GVA_REDIS.Set(context.Background(), key, jsonString, 2*time.Second)
	} else if err != nil {
		global.GVA_LOG.Error("error:", zap.Error(err))
	} else {
		//fmt.Println("从缓存里拿result:", rdRes)
		err = json.Unmarshal(rdRes, &order)
	}

	var payUrl string
	payUrl, err = vbUtil.HandlePayUrl2PAcc(vpo.OrderId)

	rep = &vboxRep.OrderSimple2PayAccountRes{
		OrderId:   vpo.OrderId,
		Money:     order.Money,
		PayUrl:    payUrl,
		Status:    order.OrderStatus,
		NotifyUrl: order.NotifyUrl,
	}

	return rep, err
}

// CreateOrder2PayAcc 创建CreateOrder2PayAcc
//
//	p := &vbox.CreateOrder2PayAccount{
//			Account:     "",
//			Key:         "",
//			Money:       10,
//			Sign:        "123",
//			ChannelCode: "600",
//			NotifyUrl:   "http://1.1.1.1",
//			OrderId:     "P1234",
//		}
func (vpoService *PayOrderService) CreateOrder2PayAcc(vpo *vboxReq.CreateOrder2PayAccount, cg *gin.Context) (rep *vboxRep.OrderSimple2PayAccountRes, err error) {
	money := vpo.Money
	cid := vpo.ChannelCode

	vpa, err := vbUtil.ValidPacc(vpo.Account, nil)
	if err != nil {
		return nil, err
	}

	// 0.1 核验商户是否开启
	if vpa.Status != 1 {
		return nil, fmt.Errorf("商户未启用，请核查！")
	}
	vpo.Key = vpa.PKey
	uidTmp := vpa.Uid

	// 1.0 校验签名
	signValid := utils.VerifySign(vpo)
	if !signValid {
		return nil, errors.New("请求参数或签名值不正确，请联系管理员核对")
	}
	global.GVA_LOG.Info("签名校验通过", zap.Any("商户ID", vpo.Account))

	orgTmp := utils2.GetSelfOrg(vpa.Uid)
	if len(orgTmp) < 1 {
		global.GVA_LOG.Warn("该账户组织归属异常，请核查！")
		return nil, fmt.Errorf("该账户组织归属异常，请核查")
	}
	orgID := orgTmp[0]

	// 1.1 ----- 校验该组织是否有此产品 -----------

	var channelCodeList []string
	// 获取组织ID
	cidKey := fmt.Sprintf(global.OrgChanSet, orgID)

	c, err := global.GVA_REDIS.Exists(context.Background(), cidKey).Result()
	if c == 0 {
		var productIds []uint
		if err != nil {
			global.GVA_LOG.Error("当前缓存池无此用户对应的orgIds，redis err", zap.Error(err))
		}
		global.GVA_LOG.Info("当前缓存池此用户对应的orgIds，查一下库。。。", zap.Any("商户", vpa.PRemark))
		orgIds := utils2.GetDeepOrg(vpa.Uid)
		db := global.GVA_DB.Model(&vbox.OrgProduct{})
		if err = db.Debug().Distinct("channel_product_id").Select("channel_product_id").Where("organization_id in ?", orgIds).Find(&productIds).Error; err != nil {
			return nil, err
		}
		if err = db.Debug().Model(&vbox.ChannelProduct{}).Select("channel_code").Where("id in ?", productIds).Find(&channelCodeList).Error; err != nil {
			return nil, err
		}

		for _, cidX := range channelCodeList {
			global.GVA_REDIS.SAdd(context.Background(), cidKey, cidX)
		}
		global.GVA_REDIS.Expire(context.Background(), cidKey, 1*time.Minute)
	} else {
		members, _ := global.GVA_REDIS.SMembers(context.Background(), cidKey).Result()
		channelCodeList = members
	}

	global.GVA_LOG.Info("当前所拥有的产品code", zap.Any("此次请求产品", vpo.ChannelCode),
		zap.Any("通道编码", channelCodeList), zap.Any("vpa.Uid", vpa.Uid), zap.Any("商户", vpa.PRemark))

	exist := utils.Contains(channelCodeList, vpo.ChannelCode)
	if !exist {
		global.GVA_LOG.Warn("该账户不存在此产品，请核查！", zap.Any("目前支持的通道", channelCodeList))
		return nil, fmt.Errorf("该账户不存在此产品，请核查！ 目前支持的通道：%v", channelCodeList)
	}

	prodInfo, err := utils.GetProductByCode(cid)
	if err != nil {
		return nil, fmt.Errorf("产品信息异常，请联系管理员")
	}

	var accSetKey string

	switch {
	case strings.Contains(prodInfo.Ext, "money"):
		accSetKey = fmt.Sprintf(prodInfo.Ext, orgID, cid, money)
	default:
		accSetKey = fmt.Sprintf(prodInfo.Ext, orgID, cid)
	}
	global.GVA_LOG.Info("ext", zap.Any("accSetKey", accSetKey))

	canUseCount := global.GVA_REDIS.ZCount(context.Background(), accSetKey, "0", "0").Val()
	if canUseCount > 0 {
		global.GVA_LOG.Info("当前后台剩余可以匹配资源", zap.Any("canUseCount", canUseCount))
	} else {
		fmt.Printf("当前组织无账号可用, org : %d", orgID)
		return nil, fmt.Errorf("库存不足，请联系对接人")
	}

	if cid == "5001" {
		shopSetKey := fmt.Sprintf(global.ChanOrgQNShopZSet, orgID, cid, money)
		canShopUseCount := global.GVA_REDIS.ZCount(context.Background(), shopSetKey, "0", "0").Val()
		if canShopUseCount > 0 {
			global.GVA_LOG.Info("当前后台剩余可以匹配资源", zap.Any("canShopUseCount", canShopUseCount))
		} else {
			fmt.Printf("当前组织无qn shop可用, org : %d", orgID)
			return nil, fmt.Errorf("库存不足，请联系对接人")
		}
	}

	var eventID, rsUrl string
	if prodInfo.Type == 1 { // 引导类，查
		eventID, err = vbUtil.HandleEventID2chShop(vpo.ChannelCode, vpo.Money, orgTmp)
		if err != nil {
			global.GVA_LOG.Error("HandleEventID2chShop该组织配置的资源不足，请核查", zap.Error(err))
			return nil, err
		}
		rsUrl, err = vbUtil.HandleResourceUrl2chShop(eventID)
		if err != nil {
			global.GVA_LOG.Error("HandleResourceUrl2chShop该组织配置的资源不足，请核查", zap.Error(err))
			return nil, err
		}
	}

	// 获取过期时间
	t, err := utils.GetCDTimeByCode(vpo.ChannelCode)
	if err != nil {
		return nil, err
	}

	global.GVA_LOG.Info("此次请求后台账号资源核查通过，订单进入等待匹配", zap.Any("orderID", vpo.OrderId), zap.Any("请求金额", money))

	jucKey := fmt.Sprintf(global.PayOrderJUCKey, vpo.OrderId)
	var count int64
	count, err = global.GVA_REDIS.Exists(context.Background(), jucKey).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}
		global.GVA_LOG.Info("当前缓存池无此订单，可继续。。。", zap.Any("orderId", vpo.OrderId))
		//global.GVA_REDIS.Set(context.Background(), vpo.OrderId, 1, 10*time.Minute)
		global.GVA_REDIS.Set(context.Background(), jucKey, 1, 5*time.Minute)
		go func() {

			expTime := time.Now().Add(t.Duration)
			order := vbox.PayOrder{
				ChannelCode: vpo.ChannelCode,
				PAccount:    vpo.Account,
				OrderId:     vpo.OrderId,
				Money:       vpo.Money,
				NotifyUrl:   vpo.NotifyUrl,
				EventId:     eventID,
				EventType:   prodInfo.Type,
				ExpTime:     &expTime,
				ResourceUrl: rsUrl,
				CreatedBy:   uidTmp,
			}

			err = global.GVA_DB.Create(&order).Error

			record := sysModel.SysOperationRecord{
				Ip:      cg.ClientIP(),
				Method:  cg.Request.Method,
				Path:    cg.Request.URL.Path,
				Agent:   cg.Request.UserAgent(),
				MarkId:  fmt.Sprintf(global.OrderRecord, vpo.OrderId),
				Type:    global.OrderType,
				Body:    "",
				Status:  200,
				Latency: time.Since(time.Now()),
				Resp:    fmt.Sprintf(global.OrderStartMsg),
				UserID:  int(uidTmp),
			}
			err = operationRecordService.CreateSysOperationRecord(record)

			// 设置检查长时间未匹配的订单
			conn, errC := mq.MQ.ConnPool.GetConnection()
			if errC != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", errC))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)
			ch, errC := conn.Channel()
			if errC != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", errC))
				return
			}

			waitOrderMsg := fmt.Sprintf("%s-%d", order.OrderId, order.ID)
			err = ch.PublishWithDelay(task.OrderStatusCheckDelayedExchange, task.OrderStatusCheckDelayedRoutingKey, []byte(waitOrderMsg), 20*time.Minute)

		}()
	} else {
		global.GVA_LOG.Info("订单已存在，请勿重复创建")
		return nil, errors.New("订单已存在，请勿重复创建")
	}

	var payUrl string
	payUrl, err = vbUtil.HandlePayUrl2PAcc(vpo.OrderId)

	rep = &vboxRep.OrderSimple2PayAccountRes{
		OrderId:   vpo.OrderId,
		Money:     vpo.Money,
		PayUrl:    payUrl,
		Status:    2,
		NotifyUrl: vpo.NotifyUrl,
	}
	return rep, err
}

// CreateOrderTest 创建CreateOrderTest
//
//	p := &vbox.CreateOrderTest{
//			Money:       10,
//			ChannelCode: "6001",
//			AuthCaptcha: "P1234",
//		}
func (vpoService *PayOrderService) CreateOrderTest(vpo *vboxReq.CreateOrderTest, cg *gin.Context) (rep *vboxRep.Order2PayAccountRes, err error) {

	// 1. 查供应库存账号是否充足 (优先从缓存池取，取空后查库取，如果库也空了，咋报错库存不足)
	money := vpo.Money
	cid := vpo.ChannelCode
	//获取当前组织
	orgTmp := utils2.GetSelfOrg(vpo.UserId)
	orgID := orgTmp[0]

	// 1.1 ----- 校验该组织是否有此产品 -----------

	var channelCodeList []string
	// 获取组织ID
	cidKey := fmt.Sprintf(global.OrgChanSet, orgID)

	c, err := global.GVA_REDIS.Exists(context.Background(), cidKey).Result()
	if c == 0 {
		var productIds []uint
		if err != nil {
			global.GVA_LOG.Error("当前缓存池无此用户对应的orgIds，redis err", zap.Error(err))
		}
		global.GVA_LOG.Info("当前缓存池此用户对应的orgIds，查一下库。。。", zap.Any("账户", vpo.Username))
		orgIds := utils2.GetDeepOrg(vpo.UserId)
		db := global.GVA_DB.Model(&vbox.OrgProduct{})
		if err = db.Debug().Distinct("channel_product_id").Select("channel_product_id").Where("organization_id in ?", orgIds).Find(&productIds).Error; err != nil {
			return nil, err
		}
		if err = db.Debug().Model(&vbox.ChannelProduct{}).Select("channel_code").Where("id in ?", productIds).Find(&channelCodeList).Error; err != nil {
			return nil, err
		}

		for _, cid := range channelCodeList {
			global.GVA_REDIS.SAdd(context.Background(), cidKey, cid)
		}
		global.GVA_REDIS.Expire(context.Background(), cidKey, 1*time.Minute)
	} else {
		members, _ := global.GVA_REDIS.SMembers(context.Background(), cidKey).Result()
		channelCodeList = members
	}

	global.GVA_LOG.Info("当前所拥有的产品code", zap.Any("此次请求产品", vpo.ChannelCode),
		zap.Any("通道编码", channelCodeList), zap.Any("vpa.Uid", vpo.UserId), zap.Any("账户", vpo.Username))

	exist := utils.Contains(channelCodeList, vpo.ChannelCode)
	if !exist {
		global.GVA_LOG.Warn("该账户不存在此产品，请核查！", zap.Any("目前支持的通道", channelCodeList))
		return nil, fmt.Errorf("该账户不存在此产品，请核查！ 目前支持的通道：%v", channelCodeList)
	}

	prodInfo, err := utils.GetProductByCode(cid)
	if err != nil {
		return nil, fmt.Errorf("产品信息异常，请联系管理员")
	}

	var accSetKey string

	switch {
	case strings.Contains(prodInfo.Ext, "money"):
		accSetKey = fmt.Sprintf(prodInfo.Ext, orgID, cid, money)
	default:
		accSetKey = fmt.Sprintf(prodInfo.Ext, orgID, cid)
	}

	canUseCount := global.GVA_REDIS.ZCount(context.Background(), accSetKey, "0", "0").Val()
	if canUseCount > 0 {
		global.GVA_LOG.Info("当前后台剩余可以匹配资源", zap.Any("canUseCount", canUseCount))
	} else {
		fmt.Printf("当前组织无账号可用, org : %d", orgID)
		return nil, fmt.Errorf("库存不足，请联系对接人")
	}

	var eventID, rsUrl string
	if prodInfo.Type == 1 { // 引导类，查
		eventID, err = vbUtil.HandleEventID2chShop(vpo.ChannelCode, vpo.Money, orgTmp)
		if err != nil {
			global.GVA_LOG.Error("HandleEventID2chShop该组织配置的资源不足，请核查", zap.Error(err))
			return nil, err
		}
		rsUrl, err = vbUtil.HandleResourceUrl2chShop(eventID)
		if err != nil {
			global.GVA_LOG.Error("HandleResourceUrl2chShop该组织配置的资源不足，请核查", zap.Error(err))
			return nil, err
		}
	}

	// 获取过期时间
	t, err := utils.GetCDTimeByCode(vpo.ChannelCode)
	if err != nil {
		return nil, err
	}
	orderId := utils.GenerateID(global.WalletEventOrderPrefix)

	global.GVA_LOG.Info("此次请求后台账号资源核查通过，订单进入等待匹配", zap.Any("orderID", orderId), zap.Any("请求金额", money))

	jucKey := fmt.Sprintf(global.PayOrderJUCKey, orderId)
	var count int64
	count, err = global.GVA_REDIS.Exists(context.Background(), jucKey).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}
		global.GVA_LOG.Info("当前缓存池无此订单，可继续。。。", zap.Any("orderId", orderId))
		//global.GVA_REDIS.Set(context.Background(), vpo.OrderId, 1, 10*time.Minute)
		global.GVA_REDIS.Set(context.Background(), jucKey, 1, 5*time.Minute)
		go func() {
			vpo.NotifyUrl, _ = vpoService.HandleNotifyUrl2Test()

			expTime := time.Now().Add(t.Duration)
			order := vbox.PayOrder{
				ChannelCode: vpo.ChannelCode,
				PAccount:    "TEST_" + vpo.Username,
				OrderId:     orderId,
				Money:       vpo.Money,
				NotifyUrl:   vpo.NotifyUrl,
				EventId:     eventID,
				EventType:   prodInfo.Type,
				ExpTime:     &expTime,
				ResourceUrl: rsUrl,
				CreatedBy:   vpo.UserId,
			}

			err = global.GVA_DB.Create(&order).Error

			record := sysModel.SysOperationRecord{
				Ip:      cg.ClientIP(),
				Method:  cg.Request.Method,
				Path:    cg.Request.URL.Path,
				Agent:   cg.Request.UserAgent(),
				MarkId:  fmt.Sprintf(global.OrderRecord, orderId),
				Type:    global.OrderType,
				Body:    "",
				Status:  200,
				Latency: time.Since(time.Now()),
				Resp:    fmt.Sprintf(global.OrderStartMsg),
				UserID:  int(vpo.UserId),
			}
			err = operationRecordService.CreateSysOperationRecord(record)

			// 设置检查长时间未匹配的订单
			conn, errC := mq.MQ.ConnPool.GetConnection()
			if errC != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", errC))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)
			if conn == nil {
				global.GVA_LOG.Error("conn is nil")
				return
			}
			ch, errC := conn.Channel()
			if errC != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", errC))
				return
			}

			waitOrderMsg := fmt.Sprintf("%s-%d", order.OrderId, order.ID)
			err = ch.PublishWithDelay(task.OrderStatusCheckDelayedExchange, task.OrderStatusCheckDelayedRoutingKey, []byte(waitOrderMsg), 20*time.Minute)

		}()
	} else {
		global.GVA_LOG.Info("订单已存在，请勿重复创建")
		return nil, errors.New("订单已存在，请勿重复创建")
	}

	var payUrl string
	payUrl, err = vbUtil.HandlePayUrl2PAcc(orderId)

	rep = &vboxRep.Order2PayAccountRes{
		OrderId:   orderId,
		Money:     vpo.Money,
		PayUrl:    payUrl,
		Status:    2,
		NotifyUrl: vpo.NotifyUrl,
	}
	return rep, err
}

func (vpoService *PayOrderService) HandleExpTime2Product(chanID string) (time.Duration, error) {
	var key string

	if global.TxContains(chanID) {
		key = "1000"
		if chanID == "1101" {
			key = "1100"
		}
	} else if global.DnfContains(chanID) {
		key = "1200"
	} else if global.J3Contains(chanID) {
		key = "2000"
	} else if global.PcContains(chanID) {
		key = "3000"
	} else if global.SdoContains(chanID) {
		key = "4000"
	} else if global.ECContains(chanID) {
		key = "6000"
	}

	var expTimeStr string
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		global.GVA_LOG.Warn("当前key不存在", zap.Any("key", key))

		var proxy vbox.Proxy
		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
		err = db.Where("status = ?", 1).Where("chan = ?", key).
			First(&proxy).Error
		if err != nil || proxy.Url == "" {
			return 0, err
		}
		expTimeStr = proxy.Url
		seconds, _ := strconv.Atoi(expTimeStr)
		duration := time.Duration(seconds) * time.Second

		global.GVA_REDIS.Set(context.Background(), key, int64(duration.Seconds()), 0)
		return duration, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
		return 0, err
	} else {
		expTimeStr, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		seconds, _ := strconv.Atoi(expTimeStr)

		duration := time.Duration(seconds) * time.Second

		global.GVA_LOG.Info("缓存池取出：", zap.Any("HandleExpTime2Product", chanID))
		return duration, err
	}
}

//// 付方获取支付url
//func (vpoService *PayOrderService) HandlePayUrl2PAcc(orderId string) (string, error) {
//	conn := global.GVA_REDIS.Conn()
//	defer conn.Close()
//	key := global.PAccPay
//	var url string
//	//paccCreateUrl, err := global.GVA_REDIS.Ping(context.Background()).Result()
//	//paccCreateUrl, err := conn.Ping(context.Background()).Result()
//	//fmt.Printf(paccCreateUrl)
//	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
//	if count == 0 {
//		if err != nil {
//			global.GVA_LOG.Error("redis ex：", zap.Error(err))
//		}
//
//		global.GVA_LOG.Warn("当前key不存在", zap.Any("key", key))
//
//		var proxy vbox.Proxy
//		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
//		err = db.Where("status = ?", 1).Where("chan = ?", key).
//			First(&proxy).Error
//		if err != nil || proxy.Url == "" {
//			return "", err
//		}
//		url = proxy.Url + orderId
//
//		//global.GVA_REDIS.Set(context.Background(), key, proxy.Url, 0)
//		conn.Set(context.Background(), key, proxy.Url, 0)
//		global.GVA_LOG.Info("查库获取", zap.Any("商户订单地址", url))
//
//		return url, nil
//	} else if err != nil {
//		global.GVA_LOG.Error("redis ex：", zap.Error(err))
//	} else {
//		var preUrl string
//		//preUrl, err = global.GVA_REDIS.Get(context.Background(), key).Result()
//		preUrl, err = conn.Get(context.Background(), key).Result()
//		url = preUrl + orderId
//		global.GVA_LOG.Info("缓存池取出", zap.Any("商户订单地址", url))
//	}
//	return url, err
//}

//func (vpoService *PayOrderService) HandleResourceUrl2chShop(eventID string) (addr string, err error) {
//	global.GVA_LOG.Info("接收event id", zap.Any("eventID", eventID))
//	//1. 如果是引导类的，获取引导地址 - channel shop
//	split := strings.Split(eventID, "_")
//	if len(split) <= 1 {
//		return "", fmt.Errorf("解析商铺prod异常，param: %s", eventID)
//	}
//	//格式 （prodID_ID）
//	ID := split[1]
//
//	var shop vbox.ChannelShop
//	err = global.GVA_DB.Debug().Model(&vbox.ChannelShop{}).Where("id = ?", ID).First(&shop).Error
//	if err != nil {
//		return "", err
//	}
//	global.GVA_LOG.Info("查出shop", zap.Any("shop", shop))
//
//	cid := shop.Cid
//
//	var payUrl string
//	switch cid {
//	case "1001": //jd
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleJDUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	case "1002": //jd
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleDYUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	case "1003": //jym
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleAlipayUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	case "1004": //zfb
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleAlipayUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	case "1005": //qb tb
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleTBUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	case "1006": //wx xcx
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleXCXUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	case "1007": //qb pdd
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandlePddUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//
//	case "1101": //jw qb tb
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleTBUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//
//	case "1201": //dnf tb
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleTBUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	case "1202": //dnf jd
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleJDUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//
//	case "2001": //j3 tb
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleTBUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	case "4001": //sdo tb
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleTBUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	case "6001": //ec jd
//		global.GVA_LOG.Info("到这一步匹配", zap.Any("cid", cid), zap.Any("payUrl", shop.Address))
//		payUrl, err = utils.HandleJDUrl(shop.Address)
//		if err != nil {
//			return "", err
//		}
//	default:
//		payUrl = shop.Address
//	}
//
//	return payUrl, nil
//}

// GetPayOrder 根据id获取订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *PayOrderService) GetPayOrder(id uint) (payOrder vbox.PayOrder, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&payOrder).Error
	rdOrderId := fmt.Sprintf(global.OrderRecord, payOrder.OrderId)
	var records []sysModel.SysOperationRecord
	global.GVA_DB.Model(&sysModel.SysOperationRecord{}).Distinct("resp,created_at").Where("mark_id = ?", rdOrderId).Scan(&records)

	ext := map[string]interface{}{
		"records": records,
	}

	if payOrder.EventId != "" && payOrder.EventType == 1 {
		split := strings.Split(payOrder.EventId, "_")
		if len(split) == 2 {
			shopId := split[1]
			var shop vbox.ChannelShop
			_ = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("id = ?", shopId).Find(&shop).Error

			var vo []vboxRep.DataWalletOverView
			global.GVA_DB.Model(&vbox.PayOrder{}).Select(`
				IFNULL(COUNT(CASE WHEN created_at >= NOW() - INTERVAL 1 HOUR THEN 1 ELSE NULL END), 0) AS x1,
				IFNULL(COUNT(CASE WHEN created_at >= NOW() - INTERVAL 1 HOUR AND order_status = 1 THEN 1 ELSE NULL END), 0) AS x2,
				IFNULL(COUNT(CASE WHEN DATE(created_at) = CURDATE() THEN 1 ELSE NULL END), 0) AS x3,
				IFNULL(COUNT(CASE WHEN DATE(created_at) = CURDATE() AND order_status = 1 THEN 1 ELSE NULL END), 0) AS x4
			`).
				Where("event_id like ?", split[0]+"%").Scan(&vo)
			if len(vo) > 0 {
				ext["shop"] = shop
				ext["dv"] = vo[0]
			}
		}
	}
	payOrder.Ext = ext
	return
}

// GetPayOrderInfoList 分页获取订单记录
func (vpoService *PayOrderService) GetPayOrderInfoList(info vboxReq.PayOrderSearch, ids []uint) (list []vbox.PayOrder, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var payOrders []vbox.PayOrder
	// 如果有条件搜索 下方会自动创建搜索语句
	if !info.StartCreatedAt.IsZero() && !info.EndCreatedAt.IsZero() {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.OrderId != "" {
		db = db.Where("order_id like ?", info.OrderId+"%")
	}
	if info.PAccount != "" {
		db = db.Where("p_account =?", info.PAccount)
	}
	if info.HandStatus != 0 {
		db = db.Where("hand_status =?", info.HandStatus)
	}
	if info.ChannelCode != "" {
		db = db.Where("channel_code =?", info.ChannelCode)
	}
	if info.OrderStatus != 0 {
		if info.OrderStatus == -1 {
			info.OrderStatus = 0
		}
		db = db.Where("order_status =?", info.OrderStatus)
	}
	if info.CbStatus != 0 {
		db = db.Where("cb_status =?", info.CbStatus)
	}
	if info.AcId != "" {
		db = db.Where("ac_id =?", info.AcId)
	}
	if info.PlatId != "" {
		db = db.Where("plat_id =?", info.PlatId)
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account = ?", info.AcAccount)
	}
	db.Where("created_by in ?", ids)

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Order("id desc").Find(&payOrders).Error

	return payOrders, total, err
}

// GetPayOrderRate 计算成率数据（成功/总数）
func (vpoService *PayOrderService) GetPayOrderRate(info vboxReq.PayOrderSearch, ids []uint) (ov []vboxRep.DataRateOverView, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	if info.StartTime > 0 && info.EndTime > 0 && info.StartTime < info.EndTime {
		location, _ := time.LoadLocation("Asia/Shanghai")
		start := time.Unix(info.StartTime, 0).In(location)
		end := time.Unix(info.EndTime, 0).In(location)
		db = db.Where("created_at BETWEEN ? AND ?", start, end)

		if info.ChannelCode != "" {
			db = db.Where("channel_code =?", info.ChannelCode)
		}

		if info.PAccount != "" {
			db = db.Where("p_account =?", info.PAccount)
		}

		if info.Keyword == "cas" {
			err = db.Select(
				"IFNULL(COUNT(CASE WHEN order_status = 1 THEN 1 ELSE NULL END), 0) AS x1,"+
					"IFNULL(COUNT(*), 0) AS x2,"+
					"IFNULL(SUM(CASE WHEN order_status = 1 THEN money ELSE 0 END), 0) AS x3,"+
					"IFNULL(SUM(money), 0) AS x4").Where("created_by in ?", ids).Find(&ov).Error
			if err != nil {
				return
			}
		} else if info.Keyword == "sum" {
			err = db.Select(
				"IFNULL(SUM(CASE WHEN order_status = 1 THEN money ELSE 0 END), 0) AS x3,"+
					"IFNULL(SUM(money), 0) AS x4").Where("created_by in ?", ids).Find(&ov).Error
			if err != nil {
				return
			}
		} else if info.Keyword == "cnt" {
			err = db.Select(
				"IFNULL(COUNT(CASE WHEN order_status = 1 THEN 1 ELSE NULL END), 0) AS x1,"+
					"IFNULL(COUNT(*), 0) AS x2").Where("created_by in ?", ids).Find(&ov).Error
			if err != nil {
				return
			}
		} else {
			return nil, fmt.Errorf("未知的统计类型")
		}

		return ov, err

	} else {
		return nil, fmt.Errorf("未传入合规时间参数")
	}

}

// GetPayOrderOverview 计算趋势图概览数据
func (vpoService *PayOrderService) GetPayOrderOverview(info vboxReq.PayOrderSearch, ids []uint) (list []vboxRep.DataOverView, err error) {
	var overViews []vboxRep.DataOverView
	jsonData, err := json.Marshal(&info)
	idsStr := utils.UintSliceToString(ids)
	md5VKey := utils.MD5V(jsonData, []byte(idsStr)...)
	key := fmt.Sprintf(global.PayOrderVOKey, md5VKey)
	rdRes, err := global.GVA_REDIS.Get(context.Background(), key).Bytes()
	if err == redis.Nil { // redis中还没有的情况，查一下库，并且去匹配设备信息
		// 创建db
		db := global.GVA_DB.Model(&vbox.PayOrder{})
		var payOrders []vbox.PayOrder
		// 如果有条件搜索 下方会自动创建搜索语句
		//空值时，默认设置为5m
		if info.Interval == "" {
			info.Interval = "5m"
		}
		if info.StartTime > 0 && info.EndTime > 0 && info.StartTime < info.EndTime {
			location, _ := time.LoadLocation("Asia/Shanghai")
			start := time.Unix(info.StartTime, 0).In(location)
			end := time.Unix(info.EndTime, 0).In(location)
			db = db.Where("cb_time BETWEEN ? AND ?", start, end)

			if info.ChannelCode != "" {
				db = db.Where("channel_code =?", info.ChannelCode)
			}
			if info.PAccount != "" {
				db = db.Where("p_account =?", info.PAccount)
			}
			if info.OrderStatus != 0 {
				db = db.Where("order_status =?", info.OrderStatus)
			}

			err = db.Debug().Where("created_by in ?", ids).Find(&payOrders).Error
			if err != nil {
				return
			}

			// 解析参数为 time.Duration 类型
			interval, errParse := time.ParseDuration(info.Interval)
			if errParse != nil {
				return nil, errParse
			}

			if info.Keyword == "sum" {
				overViews = calculateTotalMoney(payOrders, start, end, interval)
			} else if info.Keyword == "cnt" {
				overViews = calculateTotalCount(payOrders, start, end, interval)
			} else {
				return nil, fmt.Errorf("未知的统计类型")
			}

			//set redis
			jsonString, errJ := json.Marshal(overViews)
			if errJ != nil {
				return nil, fmt.Errorf("系统错误r")
			}
			global.GVA_REDIS.Set(context.Background(), key, jsonString, 60*time.Second)

			return overViews, err

		} else {
			return nil, fmt.Errorf("未传入合规时间参数")
		}
	} else if err != nil {
		fmt.Println("error:", err)
	} else {
		//fmt.Println("从缓存里拿result:", rdRes)
		err = json.Unmarshal(rdRes, &overViews)
	}

	return overViews, err

}

func calculateTotalMoney(dataList []vbox.PayOrder, startTime time.Time, endTime time.Time, interval time.Duration) []vboxRep.DataOverView {
	// 初始化结果映射
	resultMap := make(map[string]int)
	// 记录已经出现的 key，用于去重
	seenKeys := make(map[string]struct{})
	var sortedResult []vboxRep.DataOverView
	// 遍历数据并按时间间隔累加 money
	for _, item := range dataList {
		if item.CreatedAt.After(startTime) && item.CreatedAt.Before(endTime) {
			// 计算所属的时间间隔
			location, _ := time.LoadLocation("Asia/Shanghai")
			// 先将时间调整到当天的零时零分零秒
			intervalEnd := item.CreatedAt.Truncate(interval)

			// 再进行时区调整
			intervalEnd = intervalEnd.Add(interval).In(location)

			key := intervalEnd.Format("2006-01-02 15:04:05")
			resultMap[key] += item.Money

			// 检查是否已经出现过该 key，如果已经出现则跳过
			if _, exists := seenKeys[key]; !exists {
				seenKeys[key] = struct{}{}
			}

		}
	}

	//// 将 keys 提取到一个切片中，并排序
	//var keys []string
	//for key := range seenKeys {
	//	keys = append(keys, key)
	//}
	//// 对 keys 进行排序
	//sort.Strings(keys)
	points := utils.GenerateTimePoints(startTime, endTime, interval)

	for _, key := range points {
		sortedResult = append(sortedResult, vboxRep.DataOverView{
			Y: resultMap[key],
			X: key,
		})
	}
	return sortedResult
}

func calculateTotalCount(dataList []vbox.PayOrder, startTime time.Time, endTime time.Time, interval time.Duration) []vboxRep.DataOverView {
	// 初始化结果映射
	resultMap := make(map[string]int)
	// 记录已经出现的 key，用于去重
	seenKeys := make(map[string]struct{})
	var sortedResult []vboxRep.DataOverView
	// 遍历数据并按时间间隔累加 money
	for _, item := range dataList {
		if item.CreatedAt.After(startTime) && item.CreatedAt.Before(endTime) {
			// 计算所属的时间间隔
			location, _ := time.LoadLocation("Asia/Shanghai")
			// 先将时间调整到当天的零时零分零秒
			intervalEnd := item.CreatedAt.Truncate(interval)

			// 再进行时区调整
			intervalEnd = intervalEnd.Add(interval).In(location)
			key := intervalEnd.Format("2006-01-02 15:04:05")
			resultMap[key] += 1

			// 检查是否已经出现过该 key，如果已经出现则跳过
			if _, exists := seenKeys[key]; !exists {
				seenKeys[key] = struct{}{}
			}

		}
	}

	//// 将 keys 提取到一个切片中，并排序
	//var keys []string
	//for key := range seenKeys {
	//	keys = append(keys, key)
	//}
	//// 对 keys 进行排序
	//sort.Strings(keys)
	points := utils.GenerateTimePoints(startTime, endTime, interval)

	for _, key := range points {
		sortedResult = append(sortedResult, vboxRep.DataOverView{
			Y: resultMap[key],
			X: key,
		})
	}
	return sortedResult
}

func (vpoService *PayOrderService) GetPayOrderListByDt(info vboxReq.OrdersDtData, ids []uint) (list []vbox.PayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var payOrders []vbox.PayOrder
	dt := info.Dt
	if info.ChannelCode != "" {
		db = db.Where("channel_code = ?", info.ChannelCode)
	}
	db.Where("created_by in ? and DATE_FORMAT(created_at, '%Y-%m-%d') = ?", ids, dt)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Order("id desc").Find(&payOrders).Error
	if err != nil {
		return
	}

	return payOrders, total, err
}

func (vpoService *PayOrderService) GetPayOrderListLatestHour(info vboxReq.OrdersDtData, ids []uint) (list []vboxRep.OrderStatisRes, total int64, err error) {
	queryA := `
			SELECT 
				concat(
						if(HOUR(created_at)>=10,HOUR(created_at),concat('0',HOUR(created_at))) ,
						':',
						if(FLOOR(MINUTE(created_at) / 5) * 5>=10,FLOOR(MINUTE(created_at) / 5) * 5,CONCAT('0',FLOOR(MINUTE(created_at) / 5) * 5))
					) as state_time,
			    'all' as channel_code,
				SUM(CASE WHEN order_status = 1 THEN money ELSE 0 END) AS money,
				SUM(CASE WHEN order_status = 1 THEN 1 ELSE 0 END) AS cnt_nums
			FROM 
				vbox_pay_order
			WHERE 
				created_at >= DATE_SUB(NOW(), INTERVAL ? HOUR) and created_by in ?
			GROUP BY 
				state_time
			ORDER BY 
				state_time
;
`
	queryB := `
			SELECT 
				concat(
						if(HOUR(created_at)>=10,HOUR(created_at),concat('0',HOUR(created_at))) ,
						':',
						if(FLOOR(MINUTE(created_at) / 5) * 5>=10,FLOOR(MINUTE(created_at) / 5) * 5,CONCAT('0',FLOOR(MINUTE(created_at) / 5) * 5))
					) as state_time,
			    'all' as channel_code,
				SUM(CASE WHEN order_status = 1 THEN money ELSE 0 END) AS money,
				SUM(CASE WHEN order_status = 1 THEN 1 ELSE 0 END) AS cnt_nums
			FROM 
				vbox_pay_order
			WHERE 
				created_at >= DATE_SUB(NOW(), INTERVAL ? HOUR) and created_by in ? and channel_code = ?
			GROUP BY 
				state_time
			ORDER BY 
				state_time
;
`

	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	//var results []vboxRep.OrderStatisRes
	if info.ChannelCode == "" {
		rows, err := db.Raw(queryA, info.Interval, ids).Rows()
		if err != nil {
			// 处理错误
			fmt.Println(err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var result vboxRep.OrderStatisRes
			err := rows.Scan(&result.StateTime, &result.ChannelCode, &result.Money, &result.CntNums)
			if err != nil {
				// 处理错误
				fmt.Println(err.Error())
			}
			list = append(list, result)
		}
	} else {
		rows, err := db.Raw(queryB, info.Interval, ids, info.ChannelCode).Rows()
		if err != nil {
			// 处理错误
			fmt.Println(err.Error())
		}
		defer rows.Close()
		for rows.Next() {
			var result vboxRep.OrderStatisRes
			err := rows.Scan(&result.StateTime, &result.ChannelCode, &result.Money, &result.CntNums)
			if err != nil {
				// 处理错误
				fmt.Println(err.Error())
			}
			list = append(list, result)
		}
	}
	total = int64(len(list))
	return list, total, err
}

func (vpoService *PayOrderService) HandleNotifyUrl2Test() (string, error) {
	var proxy vbox.Proxy
	db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
	err := db.Where("status = ?", 1).Where("chan = ?", "test_notify").
		First(&proxy).Error
	if err != nil || proxy.Url == "" {
		return "", err
	}
	var url = proxy.Url
	return url, nil
}

/*func (vpoService *PayOrderService) HandleEventType(chanID string) (int, error) {
	// 1-商铺关联，2-付码关联

	chanCode, _ := strconv.Atoi(chanID)
	if chanCode >= 1000 && chanCode <= 1099 {
		return 1, nil
	} else if chanCode >= 1100 && chanCode <= 1199 {
		return 1, nil
	} else if chanCode >= 1200 && chanCode <= 1299 {
		return 1, nil
	} else if chanCode >= 4000 && chanCode <= 4099 {
		return 1, nil
	} else if chanCode >= 6000 && chanCode <= 6099 {
		return 1, nil
	} else if chanCode >= 2000 && chanCode <= 2099 {
		return 1, nil
	} else if chanCode >= 3000 && chanCode <= 3099 {
		return 2, nil
	}
	return 0, fmt.Errorf("不存在的event类型")
}

// HandleEventID2chShop 获取商铺关联ID （productId_ID）
func (vpoService *PayOrderService) HandleEventID2chShop(chanID string, money int, orgIDs []uint) (orgShopID string, err error) {
	// 1-商铺关联
	var vsList []vbox.ChannelShop
	if chanID == "6001" {
		orgIDs = []uint{1}
	}
	var zs []redis.Z
	var key string
	for _, orgID := range orgIDs {
		key = fmt.Sprintf(global.ChanOrgShopAddrZSet, orgID, chanID, money)
		zs, err = global.GVA_REDIS.ZRangeArgsWithScores(context.Background(), redis.ZRangeArgs{
			Key:   key,
			Start: 0,
			Stop:  -1,
		}).Result()
		if err != nil {
			return "", err
		}
		if len(zs) <= 0 { // redis 没查到，查一下库
			userIDs := utils2.GetUsersByOrgId(orgID)
			err = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("cid = ? and money = ? and status = 1", chanID, money).
				Where("created_by in ?", userIDs).Find(&vsList).Error
			if err != nil {
				return "", err
			}
			if len(vsList) <= 0 {
				continue
			}

			//如果查到库里有， 设置进 redis 中
			for _, shop := range vsList {
				global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
					Score:  float64(time.Now().Unix()),
					Member: shop.ProductId + "_" + strconv.FormatUint(uint64(shop.ID), 10),
				})
			}
		}
		break
	}

	if len(zs) <= 0 {
		global.GVA_LOG.Info("该组织配置的资源不足，请核查", zap.Any("orgIDs", orgIDs), zap.Any("chanID", chanID), zap.Any("money", money))
		return "", fmt.Errorf("该组织配置的资源不足，请核查")
	}

	z := zs[len(zs)-1] //取出最后一个，重新设置utc时间戳
	orgShopID = z.Member.(string)
	global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()), // 重新放进去，score设置最新的时间
		Member: orgShopID,
	})
	global.GVA_LOG.Info("获取引导商铺匹配信息", zap.Any("orgShopID", orgShopID))

	return orgShopID, err
}*/
