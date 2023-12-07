package vbox

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxRep "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/model"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	http2 "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"strconv"
	"strings"
	"time"
)

type PayOrderService struct {
}

// QueryOrderSimple 查询QueryOrderSimple
//
//	p := &vboxReq.QueryOrderSimple{
//			OrderId:        "123",
//		}
func (vpoService *PayOrderService) QueryOrderSimple(vpo *vboxReq.QueryOrderSimple) (rep *vboxRep.OrderSimpleRes, err error) {

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

		err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("order_id = ?", vpo.OrderId).
			Update("pay_ip", vpo.PayIp).Update("pay_region", vpo.PayRegion).Update("pay_device", vpo.PayDevice).Error
		if err != nil {
			return nil, err
		}

		// 如果event type = 2 ，搞一下码的地区匹配
		if order.EventType == 2 {
			var zs []redis.Z

			split := strings.Split(vpo.PayRegion, "|")
			province := split[2]
			ISP := split[4]

			// 查下省
			var geo model.Geo
			err = global.GVA_DB.Model(&model.Geo{}).Table("geo_provinces").
				Where("name LIKE ?", "%"+province+"%").First(&geo).Error
			if err != nil {
				return nil, err
			}

			ispPY := utils.ISP(ISP)
			provinceCode := geo.Code

			orgIDs := utils2.GetDeepOrg(order.CreatedBy)
			for _, orgID := range orgIDs {
				k := fmt.Sprintf(global.ChanOrgPayCodeZSet, orgID, order.ChannelCode, order.Money)
				zs, err = global.GVA_REDIS.ZRangeArgsWithScores(context.Background(), redis.ZRangeArgs{
					Key:   k,
					Start: 0,
					Stop:  -1,
				}).Result()

				for _, z := range zs {
					member := z.Member.(string)
					splitMem := strings.Split(member, "_")
					if len(splitMem) != 3 {
						fmt.Printf("存的值有问题，排查一下, %v", member)
						return
					} else {
						operator := splitMem[0] // 运营商
						loc := splitMem[1]      // 地区
						mid := splitMem[2]      // 码ID

						if ispPY == operator && provinceCode == loc {
							fmt.Printf("匹配到了对应地区和运营商")
							//	把付款码ID拿出来存到event ID
							order.EventId = mid
						}
					}
				}
			}

		}

		//查出来了，设置一下redis
		jsonString, err = json.Marshal(order)
		if err != nil {
			return nil, err
		}
		global.GVA_REDIS.Set(context.Background(), key, jsonString, 300*time.Second)
	} else if err != nil {
		fmt.Println("error:", err)
	} else {
		//fmt.Println("从缓存里拿result:", rdRes)
		err = json.Unmarshal(rdRes, &order)
	}

	rep = &vboxRep.OrderSimpleRes{
		OrderId:     order.OrderId,
		Account:     order.PAccount,
		Money:       order.Money,
		ResourceUrl: order.ResourceUrl,
		Status:      order.OrderStatus,
		ExpTime:     order.ExpTime,
		ChannelCode: order.ChannelCode,
	}

	return rep, err

}

// QueryOrder2PayAcc 查询QueryOrder2PayAcc
//
//	p := &vboxReq.QueryOrder2PayAccount{
//			Account:     "",
//			Key:         "",
//			Sign:        "123",
//		}
func (vpoService *PayOrderService) QueryOrder2PayAcc(vpo *vboxReq.QueryOrder2PayAccount) (rep *vboxRep.Order2PayAccountRes, err error) {
	// 1. 校验签名
	var vpa vbox.PayAccount
	count, err := global.GVA_REDIS.Exists(context.Background(), vpo.Account).Result()
	if count == 0 {
		global.GVA_LOG.Warn("缓存中暂无", zap.Any("当前 pacc", vpo.Account))
		jsonStr, _ := global.GVA_REDIS.Get(context.Background(), vpo.Account).Bytes()
		err = json.Unmarshal(jsonStr, &vpa)
	} else { //查库看有没有
		err = global.GVA_DB.Model(&vbox.PayAccount{}).Table("vbox_pay_account").
			Where("p_account = ?", vpo.Account).Find(&vpa).Error
		if err != nil {
			return nil, err
		} else { //有的话，更新一下redis
			jsonStr, _ := json.Marshal(vpa)
			global.GVA_REDIS.Set(context.Background(), vpa.PAccount, jsonStr, 0)
		}
	}

	vpo.Key = vpa.PKey
	signValid := utils.VerifySign(vpo)
	if !signValid {
		return nil, errors.New("请求参数或签名值不正确，请联系管理员核对")
	}

	// 2. 查单
	var order vbox.PayOrder
	err = global.GVA_DB.Model(&vbox.PayOrder{}).
		Where("order_id = ? and p_account = ?", vpo.OrderId, vpo.Account).
		First(&order).Error
	if err != nil {
		return nil, err
	}

	var payUrl string
	payUrl, err = vpoService.HandlePayUrl2PAcc(vpo.OrderId)

	rep = &vboxRep.Order2PayAccountRes{
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
func (vpoService *PayOrderService) CreateOrder2PayAcc(vpo *vboxReq.CreateOrder2PayAccount, ctx *gin.Context) (rep *vboxRep.Order2PayAccountRes, err error) {
	var accID, acAccount string
	money := vpo.Money
	rdConn := global.GVA_REDIS.Conn()
	defer rdConn.Close()

	var vpa vbox.PayAccount
	count, err := rdConn.Exists(context.Background(), global.PayAccPrefix+vpo.Account).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("当前缓存池无此商户，redis err", zap.Error(err))
		}
		global.GVA_LOG.Info("当前缓存池无此商户，查一下库。。。", zap.Any("入参商户ID", vpo.Account))

		err = global.GVA_DB.Table("vbox_pay_account").
			Where("p_account = ?", vpo.Account).First(&vpa).Error
		if err != nil {
			return nil, fmt.Errorf("无此商户，请核查！")
		}
		jsonStr, _ := json.Marshal(vpa)
		rdConn.Set(context.Background(), global.PayAccPrefix+vpo.Account, jsonStr, 10*time.Minute)
	} else {
		jsonStr, _ := rdConn.Get(context.Background(), global.PayAccPrefix+vpo.Account).Bytes()
		err = json.Unmarshal(jsonStr, &vpa)
	}

	vpo.Key = vpa.PKey

	// 1.0 校验签名
	//signValid := utils.VerifySign(vpo)
	//if !signValid {
	//	return nil, errors.New("请求参数或签名值不正确，请联系管理员核对")
	//}
	//global.GVA_LOG.Info("签名校验通过", zap.Any("商户ID", vpo.Account))

	// 1.1 ----- 校验该组织是否有此产品 -----------
	var channelCodeList []string
	// 获取组织ID
	orgIdTemp := utils2.GetSelfOrg(vpa.Uid)
	key := fmt.Sprintf(global.OrgChanSet, orgIdTemp[0])

	c, err := rdConn.Exists(context.Background(), key).Result()
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

		jsonStr, _ := json.Marshal(channelCodeList)
		rdConn.Set(context.Background(), key, jsonStr, 10*time.Minute)
	} else {
		jsonStr, _ := rdConn.Get(context.Background(), key).Bytes()
		err = json.Unmarshal(jsonStr, &channelCodeList)
	}

	global.GVA_LOG.Info("当前所拥有的产品code", zap.Any("通道编码", channelCodeList), zap.Any("vpa.Uid", vpa.Uid), zap.Any("商户", vpa.PRemark))
	global.GVA_LOG.Info("此次请求产品code", zap.Any("code", vpo.ChannelCode))
	exist := utils.Contains(channelCodeList, vpo.ChannelCode)
	if !exist {
		global.GVA_LOG.Warn("该账户不存在此产品，请核查！", zap.Any("目前支持的通道：%v", channelCodeList))
		return nil, fmt.Errorf("该账户不存在此产品，请核查！ 目前支持的通道：%v", channelCodeList)
	}
	// ----- 校验该组织是否有此产品 -----------

	// 2. 查供应库存账号是否充足 (优先从缓存池取，取空后查库取，如果库也空了，咋报错库存不足)
	orgIDs := utils2.GetDeepOrg(vpa.Uid)
	for _, orgID := range orgIDs {
		key := fmt.Sprintf(global.ChanOrgAccZSet, orgID, vpo.ChannelCode, money)

		var resList []string
		resList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
			Min:    "0",
			Max:    "0",
			Offset: 0,
			Count:  1,
		}).Result()

		if err != nil {
			fmt.Printf("当前组织无账号可用, org : %d", orgID)
			continue
		}
		if resList != nil && len(resList) > 0 {
			accTmp := resList[0]

			// 2.1 把账号设置为已用
			global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
				Score:  1,
				Member: accTmp,
			})

			// 2.2 把可用的账号给出来继续往下执行建单步骤
			split := strings.Split(accTmp, "_")
			accID = split[0]
			acAccount = split[1]
			break
		} else {
			fmt.Printf("当前组织无账号可用, org : %d", orgID)
			continue
		}
	}

	if accID == "" || acAccount == "" {
		global.GVA_LOG.Info("此次请求后台账号资源不足")
		return nil, fmt.Errorf("后台账号资源不足！ 请核查")
	}

	// 判断当前产品属于那种类型 1-商铺关联，2-付码关联
	eventType, err := vpoService.HandleEventType(vpo.ChannelCode)
	if err != nil {
		return nil, err
	}

	var eventID string
	var rsUrl string
	if eventType == 1 {
		eventID, err = vpoService.HandleEventID2chShop(vpo.ChannelCode, vpo.Money, orgIDs)
		rsUrl, err = vpoService.HandleResourceUrl2chShop(eventID)
	} else if eventType == 2 {
		eventID, err = vpoService.HandleEventID2payCode(vpo.ChannelCode, vpo.Money, orgIDs)
	}
	if err != nil {
		return nil, err
	}

	// 获取过期时间
	expTime, err := vpoService.HandleExpTime2Product(vpo.ChannelCode)
	if err != nil {
		return nil, err
	}

	global.GVA_LOG.Info("此次请求后台账号资源核查通过", zap.Any("请求金额", money))
	global.GVA_LOG.Info("匹配账号", zap.Any("acID", accID), zap.Any("acAccount", acAccount))

	count, err = rdConn.Exists(context.Background(), vpo.OrderId).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}
		global.GVA_LOG.Info("当前缓存池无此订单，可继续。。。", zap.Any("orderId", vpo.OrderId))
		//global.GVA_REDIS.Set(context.Background(), vpo.OrderId, 1, 10*time.Minute)
		rdConn.Set(context.Background(), vpo.OrderId, 1, 10*time.Minute)
		go func() {

			order := vbox.PayOrder{
				PlatformOid: utils.GenerateID("VB"),
				ChannelCode: vpo.ChannelCode,
				PAccount:    vpo.Account,
				OrderId:     vpo.OrderId,
				Money:       vpo.Money,
				NotifyUrl:   vpo.NotifyUrl,
				AcId:        accID,
				EventId:     eventID,
				EventType:   eventType,
				ExpTime:     time.Now().Add(expTime),
				ResourceUrl: rsUrl,
			}

			err = global.GVA_DB.Create(&order).Error

			conn, err := mq.MQ.ConnPool.GetConnection()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)

			ch, err := conn.Channel()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
			}

			body := http2.DoGinContextBody(ctx)

			od := vboxReq.PayOrderAndCtx{
				Obj: order,
				Ctx: vboxReq.Context{
					Body:      string(body),
					ClientIP:  ctx.ClientIP(),
					Method:    ctx.Request.Method,
					UrlPath:   ctx.Request.URL.Path,
					UserAgent: ctx.Request.UserAgent(),
				},
			}

			marshal, _ := json.Marshal(od)
			err = ch.Publish(task.OrderWaitExchange, task.OrderWaitKey, marshal)
			global.GVA_LOG.Info("发起一条资源匹配消息并入库初始化订单数据", zap.Any("od", od))
		}()
	} else {
		global.GVA_LOG.Info("订单已存在，请勿重复创建")
		return nil, errors.New("订单已存在，请勿重复创建")
	}

	var payUrl string
	payUrl, err = vpoService.HandlePayUrl2PAcc(vpo.OrderId)

	rep = &vboxRep.Order2PayAccountRes{
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
func (vpoService *PayOrderService) CreateOrderTest(vpo *vboxReq.CreateOrderTest) (rep *vboxRep.Order2PayAccountRes, err error) {

	// 1. 查供应库存账号是否充足 (优先从缓存池取，取空后查库取，如果库也空了，咋报错库存不足)

	chanID := vpo.ChannelCode
	var accID, acAccount string
	//获取当前组织
	orgIDs := utils2.GetDeepOrg(vpo.UserId)
	for _, orgID := range orgIDs {
		key := fmt.Sprintf(global.ChanOrgAccZSet, orgID, chanID, vpo.Money)

		var resList []string
		resList, err = global.GVA_REDIS.ZRangeByScore(context.Background(), key, &redis.ZRangeBy{
			Min:    "0",
			Max:    "0",
			Offset: 0,
			Count:  1,
		}).Result()

		if err != nil {
			fmt.Printf("当前组织无账号可用, org : %d", orgID)
			continue
		}
		if resList != nil && len(resList) > 0 {
			accTmp := resList[0]

			// 2.1 把账号设置为已用
			global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
				Score:  1,
				Member: accTmp,
			})

			// 2.2 把可用的账号给出来继续往下执行建单步骤
			split := strings.Split(accTmp, "_")
			accID = split[0]
			acAccount = split[1]
			break
		} else {
			fmt.Printf("当前组织无账号可用, org : %d", orgID)
			continue
		}
	}

	if accID == "" || acAccount == "" {
		return nil, fmt.Errorf("后台库存资源不足！ 请核查")
	}

	fmt.Printf("拿到了账号： %s  id: %s", accID, acAccount)
	var vca vbox.ChannelAccount
	err = global.GVA_DB.Model(&vbox.ChannelAccount{}).Where("ac_id = ?", accID).First(&vca).Error
	if err != nil {
		return nil, fmt.Errorf("匹配通道账号不存在！ 请核查：%v", err.Error())
	}

	oid := "TEST" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	vpo.NotifyUrl, _ = vpoService.HandleNotifyUrl2Test()

	// 判断当前产品属于那种类型 1-商铺关联，2-付码关联
	eventType, err := vpoService.HandleEventType(chanID)
	if err != nil {
		return nil, err
	}

	var eventID string
	var rsUrl string
	if eventType == 1 {
		eventID, err = vpoService.HandleEventID2chShop(chanID, vpo.Money, orgIDs)
		rsUrl, err = vpoService.HandleResourceUrl2chShop(eventID)
	} else if eventType == 2 {
		eventID, err = vpoService.HandleEventID2payCode(chanID, vpo.Money, orgIDs)
	}
	if err != nil {
		return nil, fmt.Errorf("后台shop资源不足！ 请核查, 金额: [%d]", vpo.Money)
	}

	// 获取过期时间
	expTime, err := vpoService.HandleExpTime2Product(chanID)
	if err != nil {
		return nil, err
	}

	order := &vbox.PayOrder{
		PlatformOid: oid,
		ChannelCode: chanID,
		PAccount:    "TEST_" + vpo.Username,
		EventType:   eventType,
		EventId:     eventID,
		AcId:        accID,
		OrderId:     oid,
		Money:       vpo.Money,
		NotifyUrl:   vpo.NotifyUrl,
		ResourceUrl: rsUrl,
		ExpTime:     time.Now().Add(expTime),
		CreatedBy:   vca.CreatedBy,
	}

	err = global.GVA_DB.Create(order).Error
	if err != nil {
		s := err.Error()
		if strings.Contains(s, "Duplicate") {
			return nil, errors.New("订单已存在，请勿重复创建")
		}
		return nil, err
	}
	//var (
	//	exchangeName = "vbox.order.direct"
	//	keyName      = "vbox.order.waiting"
	//)
	//marshal, _ := json.Marshal(order)
	//err = utils.NewChannel().Publish(exchangeName, keyName, marshal)

	var payUrl string
	payUrl, err = vpoService.HandlePayUrl2PAcc(oid)

	rep = &vboxRep.Order2PayAccountRes{
		OrderId:   oid,
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
	} else if global.J3Contains(chanID) {
		key = "2000"
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

// 付方获取支付url
func (vpoService *PayOrderService) HandlePayUrl2PAcc(orderId string) (string, error) {
	conn := global.GVA_REDIS.Conn()
	defer conn.Close()
	key := global.PAccPay
	var url string
	//paccCreateUrl, err := global.GVA_REDIS.Ping(context.Background()).Result()
	//paccCreateUrl, err := conn.Ping(context.Background()).Result()
	//fmt.Printf(paccCreateUrl)
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
			return "", err
		}
		url = proxy.Url + orderId

		//global.GVA_REDIS.Set(context.Background(), key, proxy.Url, 0)
		conn.Set(context.Background(), key, proxy.Url, 0)
		global.GVA_LOG.Info("查库获取", zap.Any("商户订单地址", url))

		return url, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
	} else {
		var preUrl string
		//preUrl, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		preUrl, err = conn.Get(context.Background(), key).Result()
		url = preUrl + orderId
		global.GVA_LOG.Info("缓存池取出", zap.Any("商户订单地址", url))
	}
	return url, err
}

func (vpoService *PayOrderService) HandleResourceUrl2chShop(eventID string) (addr string, err error) {
	//1. 如果是引导类的，获取引导地址 - channel shop
	split := strings.Split(eventID, "_")
	if len(split) <= 1 {
		return "", fmt.Errorf("解析商铺prod异常，param: %s", eventID)
	}
	//格式 （prodID_ID）
	ID := split[1]

	var shop vbox.ChannelShop
	db := global.GVA_DB.Model(&vbox.ChannelShop{}).Table("vbox_channel_shop")
	err = db.Where("id = ?", ID).First(&shop).Error
	if err != nil {
		return "", err
	}

	return shop.Address, nil
}

func (vpoService *PayOrderService) HandleResourceUrl2payCode(eventID string) (addr string, err error) {
	//1. 付码类 - pay code
	// 格式（mid）

	var pc vbox.ChannelPayCode
	err = global.GVA_DB.Model(&vbox.ChannelPayCode{}).Where("mid = ?", eventID).First(&pc).Error
	if err != nil {
		return "", err
	}

	return pc.ImgContent, nil
}

// GetPayOrder 根据id获取订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *PayOrderService) GetPayOrder(id uint) (payOrder vbox.PayOrder, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&payOrder).Error
	return
}

// GetPayOrderInfoList 分页获取订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *PayOrderService) GetPayOrderInfoList(info vboxReq.PayOrderSearch, ids []uint) (list []vbox.PayOrder, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var payOrders []vbox.PayOrder
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}
	err = db.Where("created_by in ?", ids).Order("id desc").Find(&payOrders).Error
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	return payOrders, total, err
}

func (vpoService *PayOrderService) GetPayOrderListByDt(info vboxReq.OrdersDtData, ids []uint) (list []vbox.PayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var payOrders []vbox.PayOrder
	dt := info.Dt
	if info.ChannelCode != "" {
		db = db.Where("channel_code = ?", info.ChannelCode)
	}
	err = db.Where("created_by in ? and DATE_FORMAT(created_at, '%Y-%m-%d') = ?", ids, dt).
		Order("id desc").
		Find(&payOrders).Error
	if err != nil {
		return
	}
	err = db.Count(&total).Error
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

func (vpoService *PayOrderService) HandleEventType(chanID string) (int, error) {
	// 1-商铺关联，2-付码关联

	chanCode, _ := strconv.Atoi(chanID)
	if chanCode >= 1000 && chanCode <= 1099 {
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
		return "", fmt.Errorf("该组织配置的资源不足，请核查")
	}

	z := zs[len(zs)-1] //取出最后一个，重新设置utc时间戳
	orgShopID = z.Member.(string)
	global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
		Score:  float64(time.Now().Unix()), // 重新放进去，score设置最新的时间
		Member: orgShopID,
	})

	return orgShopID, err
}

// HandleEventID2payCode 获取付码关联ID （productId_ID）
func (vpoService *PayOrderService) HandleEventID2payCode(chanID string, money int, orgIDs []uint) (orgPayCodeID string, err error) {
	// 2-付码关联
	var pcList []vbox.ChannelPayCode

	var zs []redis.Z
	var key string
	for _, orgID := range orgIDs {
		key = fmt.Sprintf(global.ChanOrgPayCodeZSet, orgID, chanID, money)
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
			// 当前time
			now := time.Now()

			err = global.GVA_DB.Debug().Model(&vbox.ChannelPayCode{}).Where("cid = ? and money = ? and code_status = 2", chanID, money).
				Where("created_by in ?", userIDs).Where("time_limit < ?", now).Find(&pcList).Error
			if err != nil {
				return "", err
			}
			if len(pcList) <= 0 {
				continue
			}

			//如果查到库里有， 设置进 redis 中
			for _, pc := range pcList {
				global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
					Score:  0,
					Member: pc.Operator + "_" + pc.Location + "_" + pc.Mid,
				})
			}
		}
		break
	}

	if len(zs) <= 0 {
		return "", fmt.Errorf("该组织配置的资源不足，请核查")
	}

	return "", nil
}
