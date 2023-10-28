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
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
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
	//var order vboxRep.VboxPayOrderRes
	var order vbox.PayOrder
	err = global.GVA_DB.Model(&vbox.PayOrder{}).
		Where("order_id = ?", vpo.OrderId).
		First(&order).Error
	if err != nil {
		return nil, err
	}

	fmt.Println("vpo=", vpo)
	err = global.GVA_DB.Model(&vbox.PayOrder{}).Where("order_id = ?", vpo.OrderId).
		Update("pay_ip", vpo.PayIp).Update("pay_region", vpo.PayRegion).Update("pay_device", vpo.PayDevice).Error

	if err != nil {
		return nil, err
	}

	var resUrl string
	resUrl, err = HandelResourceUrl(order)

	rep = &vboxRep.OrderSimpleRes{
		OrderId:     order.OrderId,
		Account:     order.PAccount,
		Money:       order.Money,
		ResourceUrl: resUrl,
		Status:      order.OrderStatus,
		ChannelCode: order.ChannelCode,
	}

	return rep, err

}

func HandelResourceUrl(order vbox.PayOrder) (string, error) {
	return order.ResourceUrl, nil
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
	payUrl, err = HandelPayUrl2Pacc(vpo.OrderId)

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
func (vpoService *PayOrderService) CreateOrder2PayAcc(vpo *vboxReq.CreateOrder2PayAccount) (rep *vboxRep.Order2PayAccountRes, err error) {
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

	// 1. 校验签名
	/*signValid := utils.VerifySign(vpo)
	if !signValid {
		return nil, errors.New("请求参数或签名值不正确，请联系管理员核对")
	}*/

	// 2. 查供应库存账号是否充足 (优先从缓存池取，取空后查库取，如果库也空了，咋报错库存不足)
	/*var total int64 = 0
	userList, tot, err := GetOwnerUserIdsList(vpa.Uid)
	var idList []int
	for _, user := range userList {
		idList = append(idList, int(user.ID))
	}
	if err != nil || tot == 0 {
		return
	}
	db := global.GVA_DB.Model(&vbox.ChannelAccount{}).Table("vbox_channel_account").
		Where("uid in (?)", idList).Count(&total)

	limit, offset := utils.RandSize2DB(int(total), 20)
	var vcas []vbox.ChannelAccount
	err = db.Debug().Where("status = ? and sys_status = ?", 1, 1).Where("cid = ?", vpo.ChannelCode).
		Where("uid in (?)", idList).Limit(limit).Offset(offset).
		Find(&vcas).Error
	if err != nil || len(vcas) == 0 {
		if len(vcas) == 0 {
			err = errors.New("库存不足！ 请联系对接人。")
		}
		return nil, err
	}

	vca := vcas[0]*/

	// ----- 校验该组织是否有此产品 -----------
	var channelCodeList []string
	c, err := rdConn.Exists(context.Background(), global.UserOrgChannelCodePrefix+strconv.FormatUint(uint64(vpa.Uid), 10)).Result()
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
		rdConn.Set(context.Background(), global.UserOrgChannelCodePrefix+strconv.FormatUint(uint64(vpa.Uid), 10), jsonStr, 10*time.Minute)
	} else {
		jsonStr, _ := rdConn.Get(context.Background(), global.UserOrgChannelCodePrefix+strconv.FormatUint(uint64(vpa.Uid), 10)).Bytes()
		err = json.Unmarshal(jsonStr, &channelCodeList)
	}

	global.GVA_LOG.Info("当前所拥有的产品code", zap.Any("通道编码", channelCodeList), zap.Any("vpa.Uid", vpa.Uid), zap.Any("商户", vpa.PRemark))
	exist := utils.Contains(channelCodeList, vpo.ChannelCode)
	if !exist {
		global.GVA_LOG.Warn("该账户不存在此产品，请核查！", zap.Any("目前支持的通道：%v", channelCodeList))
		return nil, fmt.Errorf("该账户不存在此产品，请核查！ 目前支持的通道：%v", channelCodeList)
	}
	/*var checkTotal int64
	if err = global.GVA_DB.Debug().Model(&vbox.ChannelProduct{}).Where("id in ?", productIds).Where("channel_code = ?", vpo.ChannelCode).Count(&checkTotal).Error; err != nil {
		return nil, err
	}
	if checkTotal < 1 {
		return nil, fmt.Errorf("该账户不存在此产品，请核查！")
	}*/
	// ----- 校验该组织是否有此产品 -----------

	count, err = rdConn.Exists(context.Background(), vpo.OrderId).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}
		global.GVA_LOG.Info("当前缓存池无此订单，可继续。。。", zap.Any("orderId", vpo.OrderId))
		//global.GVA_REDIS.Set(context.Background(), vpo.OrderId, 1, 10*time.Minute)
		rdConn.Set(context.Background(), vpo.OrderId, 1, 10*time.Minute)
		go func() {
			order := &vbox.PayOrder{
				PlatformOid:    utils.GenerateID("VB"),
				ChannelCode:    vpo.ChannelCode,
				PAccount:       vpo.Account,
				OrderId:        vpo.OrderId,
				Money:          vpo.Money,
				NotifyUrl:      vpo.NotifyUrl,
				OrderStatus:    2,
				CallbackStatus: 2,
			}

			err = global.GVA_DB.Create(order).Error
			go func() {
				marshal, _ := json.Marshal(order)
				conn, err := mq.MQ.ConnPool.GetConnection()
				if err != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
				}
				defer mq.MQ.ConnPool.ReturnConnection(conn)

				ch, err := conn.Channel()
				if err != nil {
					global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
				}

				err = ch.Publish(task.OrderWaitExchange, task.OrderWaitKey, marshal)
			}()
		}()
	} else {
		return nil, errors.New("订单已存在，请勿重复创建")
	}

	var payUrl string
	payUrl, err = HandelPayUrl2Pacc(vpo.OrderId)

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

	var total int64 = 0
	// 1. 查供应库存账号是否充足 (优先从缓存池取，取空后查库取，如果库也空了，咋报错库存不足)

	idList := utils2.GetCurrentUserIDs(vpo.UserId)
	db := global.GVA_DB.Model(&vbox.ChannelAccount{}).Table("vbox_channel_account").
		Where("uid in ?", idList).Count(&total)

	limit, offset := utils.RandSize2DB(int(total), 20)
	var vcas []vbox.ChannelAccount
	err = db.Debug().Where("status = ? and sys_status = ?", 1, 1).Where("cid = ?", vpo.ChannelCode).
		Where("uid in (?)", idList).Limit(limit).Offset(offset).
		Find(&vcas).Error
	if err != nil || len(vcas) == 0 {
		if len(vcas) == 0 {
			err = errors.New("无库存账号！ 请联系对接人。")
		}
		return nil, err
	}

	vca := vcas[0]

	oid := "TEST" + strconv.FormatInt(time.Now().UnixMilli(), 10)

	vpo.NotifyUrl, _ = HandelNotifyUrl2Test(oid)

	order := &vbox.PayOrder{
		PlatformOid:    oid,
		ChannelCode:    vca.Cid,
		Uid:            vca.Uid,
		PAccount:       "TEST_" + vpo.Username,
		OrderId:        oid,
		Money:          vpo.Money,
		NotifyUrl:      vpo.NotifyUrl,
		OrderStatus:    2,
		CallbackStatus: 2,
		ResourceUrl:    HandleResourceUrl(vca, vca.Cid, vpo.Money),
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
	payUrl, err = HandelPayUrl2Pacc(oid)

	rep = &vboxRep.Order2PayAccountRes{
		OrderId:   oid,
		Money:     vpo.Money,
		PayUrl:    payUrl,
		Status:    2,
		NotifyUrl: vpo.NotifyUrl,
	}
	return rep, err
}

func HandelPayUrl2Pacc(orderId string) (string, error) {
	conn := global.GVA_REDIS.Conn()
	defer conn.Close()
	key := "pacc_create"
	var url string
	//paccCreateUrl, err := global.GVA_REDIS.Ping(context.Background()).Result()
	paccCreateUrl, err := conn.Ping(context.Background()).Result()
	fmt.Printf(paccCreateUrl)
	count, err := global.GVA_REDIS.Exists(context.Background(), key).Result()
	if count == 0 {
		if err != nil {
			global.GVA_LOG.Error("redis ex：", zap.Error(err))
		}

		global.GVA_LOG.Warn("当前key不存在", zap.Any("key", key))

		var proxy vbox.Proxy
		db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
		err = db.Where("status = ?", 1).Where("chan = ?", "pacc_create").
			First(&proxy).Error
		if err != nil || proxy.Url == "" {
			return "", err
		}
		url = proxy.Url + orderId

		//global.GVA_REDIS.Set(context.Background(), key, proxy.Url, 0)
		conn.Set(context.Background(), key, proxy.Url, 0)
		return url, nil
	} else if err != nil {
		global.GVA_LOG.Error("redis ex：", zap.Error(err))
	} else {
		var preUrl string
		//preUrl, err = global.GVA_REDIS.Get(context.Background(), key).Result()
		preUrl, err = conn.Get(context.Background(), key).Result()
		url = preUrl + orderId
		global.GVA_LOG.Info("缓存池取出：", zap.Any("pacc create url", url))
	}
	return url, err
}

func HandelNotifyUrl2Test(orderId string) (string, error) {
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

func HandleResourceUrl(vca vbox.ChannelAccount, cid string, money int) string {
	//1. 如果是引导类的，获取引导地址 - channel shop
	var shop vbox.ChannelShop
	db := global.GVA_DB.Model(&vbox.ChannelShop{}).Table("vbox_channel_shop")
	err := db.Where("status = ?", 1).
		Where("money = ?", money).
		Where("cid = ?", cid).
		First(&shop).Error
	if err == nil {
		return shop.Address
	}

	return ""
}

// GetPayOrder 根据id获取订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *PayOrderService) GetPayOrder(id uint) (payOrder vbox.PayOrder, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&payOrder).Error
	return
}

// GetPayOrderInfoList 分页获取订单记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *PayOrderService) GetPayOrderInfoList(info vboxReq.PayOrderSearch) (list []vbox.PayOrder, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.PayOrder{})
	var payOrders []vbox.PayOrder
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&payOrders).Error
	return payOrders, total, err
}
