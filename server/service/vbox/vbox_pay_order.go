package vbox

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/channelshop"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxRep "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"strings"
	"time"
)

type VboxPayOrderService struct {
}

// QueryOrderSimple 查询QueryOrderSimple
//
//	p := &vboxReq.QueryOrderSimple{
//			OrderId:        "123",
//		}
func (vpoService *VboxPayOrderService) QueryOrderSimple(vpo *vboxReq.QueryOrderSimple) (rep *vboxRep.OrderSimpleRes, err error) {

	// 1. 查单
	var order vbox.VboxPayOrder
	err = global.GVA_DB.Model(&vbox.VboxPayOrder{}).
		Where("order_id = ?", vpo.OrderId).
		First(&order).Error
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

func HandelResourceUrl(order vbox.VboxPayOrder) (string, error) {
	return order.ResourceUrl, nil
}

// QueryOrder2PayAcc 查询QueryOrder2PayAcc
//
//	p := &vboxReq.QueryOrder2PayAccount{
//			Account:     "",
//			Key:         "",
//			Sign:        "123",
//		}
func (vpoService *VboxPayOrderService) QueryOrder2PayAcc(vpo *vboxReq.QueryOrder2PayAccount) (rep *vboxRep.Order2PayAccountRes, err error) {

	// 1. 校验签名
	signValid := utils.VerifySign(vpo)
	if !signValid {
		return nil, errors.New("请求参数或签名值不正确，请联系管理员核对")
	}

	// 2. 查单
	var order vbox.VboxPayOrder
	err = global.GVA_DB.Model(&vbox.VboxPayOrder{}).
		Where("order_id = ? and p_account = ?", vpo.OrderId, vpo.Account).
		First(order).Error
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
func (vpoService *VboxPayOrderService) CreateOrder2PayAcc(vpo *vboxReq.CreateOrder2PayAccount) (rep *vboxRep.Order2PayAccountRes, err error) {

	// 1. 校验签名
	signValid := utils.VerifySign(vpo)
	if !signValid {
		return nil, errors.New("请求参数或签名值不正确，请联系管理员核对")
	}

	var total int64 = 0
	// 2. 查供应库存账号是否充足 (优先从缓存池取，取空后查库取，如果库也空了，咋报错库存不足)
	db := global.GVA_DB.Model(&vbox.ChannelAccount{}).Table("vbox_channel_account").Count(&total)

	limit, offset := utils.RandSize2DB(int(total), 20)
	var vcas []vbox.ChannelAccount
	err = db.Debug().Where("status = ? and sys_status = ?", 1, 1).Where("cid = ?", vpo.ChannelCode).
		Limit(limit).Offset(offset).
		Find(&vcas).Error
	if err != nil || len(vcas) == 0 {
		return nil, err
	}

	vca := vcas[0]

	order := &vbox.VboxPayOrder{
		PlatformOid:    time.Now().String(),
		ChannelCode:    vca.Cid,
		Uid:            vca.Uid,
		PAccount:       vpo.Account,
		OrderId:        vpo.OrderId,
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

func HandelPayUrl2Pacc(orderId string) (string, error) {
	var proxy vbox.Proxy
	db := global.GVA_DB.Model(&vbox.Proxy{}).Table("vbox_proxy")
	err := db.Where("status = ?", 1).Where("chan = ?", "pacc_create").
		First(&proxy).Error
	if err != nil || proxy.Url == "" {
		return "", err
	}
	var url = proxy.Url + orderId
	return url, nil
}

func HandleResourceUrl(vca vbox.ChannelAccount, cid string, money int) string {
	return ""
}

// GetVboxPayOrder 根据id获取VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) GetVboxPayOrder(id uint) (vpo vbox.VboxPayOrder, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vpo).Error
	return
}

// GetVboxPayOrderInfoList 分页获取VboxPayOrder记录
// Author [piexlmax](https://github.com/piexlmax)
func (vpoService *VboxPayOrderService) GetVboxPayOrderInfoList(info vboxReq.VboxPayOrderSearch, ids []int) (list []vboxRep.VboxPayOrderRes, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vboxRep.VboxPayOrderRes{}).Table("vbox_pay_order")
	var vpos []vboxRep.VboxPayOrderRes

	db = db.Where("vbox_pay_order.uid in (?)", ids)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("vbox_pay_order.created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.OrderId != "" {
		db = db.Where("vbox_pay_order.order_id LIKE ?", "%"+info.OrderId+"%")
	}
	if info.PAccount != "" {
		db = db.Where("vbox_pay_order.p_account = ?", info.PAccount)
	}
	if info.AcId != "" {
		db = db.Where("vbox_pay_order.ac_id = ?", info.AcId)
	}
	if info.ChannelCode != "" {
		db = db.Where("vbox_pay_order.channel_code = ?", info.ChannelCode)
	}
	if info.PlatformOid != "" {
		db = db.Where("vbox_pay_order.platform_oid = ?", info.PlatformOid)
	}
	if info.PayRegion != "" {
		db = db.Where("vbox_pay_order.pay_region = ?", info.PayRegion)
	}
	if info.ResourceUrl != "" {
		db = db.Where("vbox_pay_order.resource_url = ?", info.ResourceUrl)
	}
	if info.NotifyUrl != "" {
		db = db.Where("vbox_pay_order.notify_url = ?", info.NotifyUrl)
	}
	if info.OrderStatus != 0 {
		db = db.Where("vbox_pay_order.order_status = ?", info.OrderStatus)
	}
	if info.CallbackStatus != 0 {
		db = db.Where("vbox_pay_order.callback_status = ?", info.CallbackStatus)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Select("vbox_pay_order.*, vbox_channel_account.ac_remark, vbox_channel_account.ac_account, vbox_channel_account.uid as user_id").
		Joins("LEFT JOIN vbox_channel_account ON vbox_channel_account.ac_id = vbox_pay_order.ac_id").
		Scan(&vpos).Error

	return vpos, total, err
}

// GetUsersVboxPayOrderInfoList 分页获取VboxPayOrder记录
// Author youga
func (vpoService *VboxPayOrderService) GetUsersVboxPayOrderInfoList(ids []int, num int) (list []vbox.VboxPayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxPayOrder{})
	var vpos []vbox.VboxPayOrder
	err = db.Where("uid in (?) and DATE(create_time) = (CURDATE() - INTERVAL ? DAY)", ids, num).Find(&vpos).Error
	total = int64(len(vpos))
	return vpos, total, err
}

// Author youga
func (vpoService *VboxPayOrderService) GetUsersVboxPayOrderInfoHList(ids []int, num int) (list []vbox.VboxPayOrder, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxPayOrder{})
	var vpos []vbox.VboxPayOrder
	if num == 1 {
		err = db.Where("uid in (?) and created_at >= (CURDATE() - INTERVAL ? HOUR)", ids, num).Find(&vpos).Error
	} else {
		err = db.Where("uid in (?) and created_at >= (CURDATE() - INTERVAL ? HOUR) and created_at < (CURDATE() - INTERVAL 1 HOUR)", ids, num).Find(&vpos).Error
	}
	total = int64(len(vpos))
	return vpos, total, err
}

func (vpoService *VboxPayOrderService) GetVboxUserPayOrderAnalysis(id uint, idList []int) (list []vboxRep.VboxUserOrderPayAnalysis, total int64, err error) {
	query := `
        SELECT count(1)
		FROM vbox_channel_shop
		WHERE  uid = ? and status = ?;
    `

	tInfoList, tOrderTotal, err := vpoService.GetUsersVboxPayOrderInfoList(idList, 0)
	if err != nil {
		return
	}
	yInfoList, yOrderTotal, err := vpoService.GetUsersVboxPayOrderInfoList(idList, 1)
	if err != nil {
		return
	}

	tGroupedCounts := make(map[int]int16)
	tOkGroupedCounts := make(map[int]int16)
	tOkGroupedCosts := make(map[int]int)

	yGroupedCounts := make(map[int]int16)
	yOkGroupedCounts := make(map[int]int16)
	yOkGroupedCosts := make(map[int]int)

	for _, order := range tInfoList {
		uid := order.Uid
		tGroupedCounts[int(uid)]++
		if order.OrderStatus == 1 {
			tOkGroupedCounts[int(uid)]++
			tOkGroupedCosts[int(uid)] += order.Money
		}
	}

	for _, order := range yInfoList {
		uid := order.Uid
		yGroupedCounts[int(uid)]++
		if order.OrderStatus == 1 {
			yOkGroupedCounts[int(uid)]++
			yOkGroupedCosts[int(uid)] += order.Money
		}
	}

	for _, uid := range idList {
		tOrderQuantify := tOrderTotal
		tOkOrderQuantify := 0
		tOkRate := 0
		tInCome := 0
		yOrderQuantify := yOrderTotal
		yOkOrderQuantify := 0
		yOkRate := 0
		yInCome := 0
		// 判断 tGroupedCounts 中是否包含指定的 uid 键
		_, tContainsUID := tGroupedCounts[uid]
		_, tOkContainsUID := tOkGroupedCounts[uid]
		_, yContainsUID := yGroupedCounts[uid]
		_, yOkContainsUID := yOkGroupedCounts[uid]

		if tContainsUID {
			tOrderQuantify = int64(tGroupedCounts[uid])
			if tOkContainsUID {
				tOkOrderQuantify = int(tOkGroupedCounts[uid])
			}

			if tOrderQuantify > 0 {
				result := float64(tOkOrderQuantify) / float64(tOrderQuantify)
				tOkRate = int(result * 100)
				tInCome = tOkGroupedCosts[uid]
			}

		}

		if yContainsUID {

			yOrderQuantify = int64(yGroupedCounts[uid])
			if yOkContainsUID {
				yOkOrderQuantify = int(yOkGroupedCounts[uid])
			}

			if yOrderQuantify > 0 {
				result := float64(yOkOrderQuantify) / float64(yOrderQuantify)
				yOkRate = int(result * 100)
				yInCome = yOkGroupedCosts[uid]
			}

		}

		var userInfo system.SysUser
		err = global.GVA_DB.Where("`id` = ?", uid).First(&userInfo).Error
		if err != nil {
			return
		}

		var openTotal int64
		var closeTotal int64

		// 创建db
		err = global.GVA_DB.Model(&channelshop.ChannelShop{}).Raw(query, uid, 1).Count(&openTotal).Error
		if err != nil {
			return
		}
		err = global.GVA_DB.Model(&channelshop.ChannelShop{}).Raw(query, uid, 0).Count(&closeTotal).Error
		if err != nil {
			return
		}
		rechargeData, errB := GetVboxUserWalletAvailablePoints(uint(uid))
		if errB != nil {
			return
		}

		entity := vboxRep.VboxUserOrderPayAnalysis{
			Uid:              uint(uid),
			Username:         userInfo.Username,
			Balance:          rechargeData,
			ChIdCnt:          int(openTotal + closeTotal),
			OpenChId:         int(openTotal),
			YOrderQuantify:   int(yOrderQuantify),
			YOkOrderQuantify: yOkOrderQuantify,
			YOkRate:          yOkRate,
			YInCome:          yInCome,
			TOrderQuantify:   int(tOrderQuantify),
			TOkOrderQuantify: tOkOrderQuantify,
			TOkRate:          tOkRate,
			TInCome:          tInCome,
		}
		list = append(list, entity)

	}
	return list, int64(len(idList)), err

}

func (vpoService *VboxPayOrderService) GetVboxUserPayOrderAnalysisH(id uint, idList []int) (list []vboxRep.VboxUserOrderPayAnalysisH, total int64, err error) {
	query := `
        SELECT count(1)
		FROM vbox_channel_shop
		WHERE  uid = ? and status = ?;
    `
	queryChB := `
        SELECT count(1)
		FROM vbox_channel_shop
		WHERE  uid = ? and status = ? and created_at >= (CURDATE() - INTERVAL ? HOUR);
    `

	tInfoList, tOrderTotal, err := vpoService.GetUsersVboxPayOrderInfoHList(idList, 1)
	if err != nil {
		return
	}
	yInfoList, yOrderTotal, err := vpoService.GetUsersVboxPayOrderInfoHList(idList, 2)
	if err != nil {
		return
	}

	tGroupedCounts := make(map[int]int16)
	tOkGroupedCounts := make(map[int]int16)
	tOkGroupedCosts := make(map[int]int)

	yGroupedCounts := make(map[int]int16)
	yOkGroupedCounts := make(map[int]int16)
	yOkGroupedCosts := make(map[int]int)

	for _, order := range tInfoList {
		uid := order.Uid
		tGroupedCounts[int(uid)]++
		if order.OrderStatus == 1 {
			tOkGroupedCounts[int(uid)]++
			tOkGroupedCosts[int(uid)] += order.Money
		}
	}

	for _, order := range yInfoList {
		uid := order.Uid
		yGroupedCounts[int(uid)]++
		if order.OrderStatus == 1 {
			yOkGroupedCounts[int(uid)]++
			yOkGroupedCosts[int(uid)] += order.Money
		}
	}

	for _, uid := range idList {
		tOrderQuantify := tOrderTotal
		tOkOrderQuantify := 0
		tOkRate := 0
		tInCome := 0
		yOrderQuantify := yOrderTotal
		yOkOrderQuantify := 0
		yOkRate := 0
		yInCome := 0
		// 判断 tGroupedCounts 中是否包含指定的 uid 键
		_, tContainsUID := tGroupedCounts[uid]
		_, tOkContainsUID := tOkGroupedCounts[uid]
		_, yContainsUID := yGroupedCounts[uid]
		_, yOkContainsUID := yOkGroupedCounts[uid]

		if tContainsUID {
			tOrderQuantify = int64(tGroupedCounts[uid])
			if tOkContainsUID {
				tOkOrderQuantify = int(tOkGroupedCounts[uid])
			}

			if tOrderQuantify > 0 {
				result := float64(tOkOrderQuantify) / float64(tOrderQuantify)
				tOkRate = int(result * 100)
				tInCome = tOkGroupedCosts[uid]
			}

		}

		if yContainsUID {

			yOrderQuantify = int64(yGroupedCounts[uid])
			if yOkContainsUID {
				yOkOrderQuantify = int(yOkGroupedCounts[uid])
			}

			if yOrderQuantify > 0 {
				result := float64(yOkOrderQuantify) / float64(yOrderQuantify)
				yOkRate = int(result * 100)
				yInCome = yOkGroupedCosts[uid]
			}

		}

		var userInfo system.SysUser
		err = global.GVA_DB.Where("`id` = ?", uid).First(&userInfo).Error
		if err != nil {
			return
		}

		var openTotal int64
		var newOpenTotal int64
		var closeTotal int64

		// 创建db
		err = global.GVA_DB.Model(&channelshop.ChannelShop{}).Raw(queryChB, uid, 1, 1).Count(&newOpenTotal).Error
		if err != nil {
			return
		}
		err = global.GVA_DB.Model(&channelshop.ChannelShop{}).Raw(query, uid, 1).Count(&openTotal).Error
		if err != nil {
			return
		}
		err = global.GVA_DB.Model(&channelshop.ChannelShop{}).Raw(query, uid, 0).Count(&closeTotal).Error
		if err != nil {
			return
		}
		rechargeData, errB := GetVboxUserWalletAvailablePoints(uint(uid))
		if errB != nil {
			return
		}

		entity := vboxRep.VboxUserOrderPayAnalysisH{
			Uid:              uid,
			Username:         userInfo.Username,
			Balance:          rechargeData,
			ChIdCnt:          int(openTotal + closeTotal),
			OpenChId:         int(openTotal),
			NewOpenChId:      int(newOpenTotal),
			YOrderQuantify:   int(yOrderQuantify),
			YOkOrderQuantify: yOkOrderQuantify,
			YOkRate:          yOkRate,
			YInCome:          yInCome,
			TOrderQuantify:   int(tOrderQuantify),
			TOkOrderQuantify: tOkOrderQuantify,
			TOkRate:          tOkRate,
			TInCome:          tInCome,
		}
		list = append(list, entity)

	}
	return list, int64(len(idList)), err

}

func (vpoService *VboxPayOrderService) GetVboxUserPayOrderAnalysisIncomeCharts(id uint, idList []int) (resData vboxRep.LineChartData, err error) {
	query := `
        SELECT coalesce(sum(cost),0) as totalIncome
		FROM vbox_pay_order
		WHERE  uid = ? and order_status = ? and  DATE(created_at) = ?;
    `
	queryB := `
		select DATE_FORMAT(dt, '%Y-%m-%d') as dt from (
		    SELECT distinct DATE(created_at)  as dt
			FROM vbox_pay_order
			WHERE  uid in (?)
		)t
        order by dt;
    `

	rowDt, err := global.GVA_DB.Raw(queryB, idList).Rows()
	var dts []string
	if err != nil {
		return
	}
	for rowDt.Next() {
		var dt string
		scanErr := rowDt.Scan(&dt)
		if scanErr != nil {
			return
		}
		dts = append(dts, dt)
	}
	resData.XData = dts

	for _, uid := range idList {
		var userInfo system.SysUser
		err = global.GVA_DB.Where("`id` = ?", uid).First(&userInfo).Error
		if err != nil {
			return
		}
		resData.LegendData = append(resData.LegendData, userInfo.Username)

		var income []int
		for _, dt := range dts {
			result := 0
			err = global.GVA_DB.Model(&vbox.VboxPayOrder{}).Raw(query, uid, 1, dt).Scan(&result).Error
			if err != nil {
				return
			}
			income = append(income, result)
		}
		entity := vboxRep.LineChartDataYSeries{
			Name:   userInfo.Username,
			Type:   "line",
			Smooth: true,
			Data:   income,
		}
		resData.Lists = append(resData.Lists, entity)
	}

	return resData, err
}

func (vpoService *VboxPayOrderService) GetSelectPayOrderAnalysisQuantifyCharts(uid int) (resData vboxRep.LineChartData, err error) {
	query := `
        SELECT coalesce(count(1),0) as totalIncome
		FROM vbox_pay_order
		WHERE  uid = ? and order_status = ? and  DATE(created_at) = ? and c_channel_id= ?;
    `
	queryB := `
		select DATE_FORMAT(dt, '%Y-%m-%d') as dt from (
		    SELECT distinct DATE(create_time)  as dt
			FROM vbox_pay_order
			WHERE  uid = ?  and DATE(created_at) >= (CURDATE() - INTERVAL 30 DAY)
		)t
        order by dt;
    `
	queryC := `
		SELECT distinct channel
		FROM vbox_channel_shop
		WHERE  uid = ?
		;
    `
	rowCh, err := global.GVA_DB.Raw(queryC, uid).Rows()
	var channels []string
	if err != nil {
		return
	}
	for rowCh.Next() {
		var ch string
		scanErr := rowCh.Scan(&ch)
		if scanErr != nil {
			return
		}
		channels = append(channels, ch)
	}
	resData.LegendData = channels

	rowDt, err := global.GVA_DB.Raw(queryB, uid).Rows()
	var dts []string
	if err != nil {
		return
	}
	for rowDt.Next() {
		var dt string
		scanErr := rowDt.Scan(&dt)
		if scanErr != nil {
			return
		}
		dts = append(dts, dt)
	}
	resData.XData = dts

	for _, channelId := range channels {

		var income []int
		for _, dt := range dts {
			result := 0
			err = global.GVA_DB.Model(&vbox.VboxPayOrder{}).Raw(query, uid, 1, dt, channelId).Scan(&result).Error
			if err != nil {
				return
			}
			income = append(income, result)
		}
		entity := vboxRep.LineChartDataYSeries{
			Name:   channelId,
			Type:   "line",
			Smooth: true,
			Data:   income,
		}
		resData.Lists = append(resData.Lists, entity)
	}

	return resData, err
}

func (vpoService *VboxPayOrderService) GetSelectPayOrderAnalysisIncomeBarCharts(uid int) (resData vboxRep.CustomBarChartData, err error) {
	query := `
        SELECT coalesce(sum(cost),0) as totalIncome
		FROM vbox_pay_order
		WHERE  uid = ? and order_status = ? and  DATE(created_at) = ?;
    `
	queryB := `
		select DATE_FORMAT(dt, '%Y-%m-%d') as dt from (
		    SELECT distinct DATE(created_at)  as dt
			FROM vbox_pay_order
			WHERE  uid = ?  and DATE(created_at) >= (CURDATE() - INTERVAL 30 DAY)
		)t
        order by dt;
    `

	rowDt, err := global.GVA_DB.Raw(queryB, uid).Rows()
	var dts []string
	if err != nil {
		return
	}
	for rowDt.Next() {
		var dt string
		scanErr := rowDt.Scan(&dt)
		if scanErr != nil {
			return
		}
		dts = append(dts, dt)
	}
	resData.XData = dts

	var income []int
	for _, dt := range dts {
		result := 0
		err = global.GVA_DB.Model(&vbox.VboxPayOrder{}).Raw(query, uid, 1, dt).Scan(&result).Error
		if err != nil {
			return
		}
		income = append(income, result)
	}
	resData.Lists = income

	return resData, err
}

func (vpoService *VboxPayOrderService) GetSelectPayOrderAnalysisChannelIncomeCharts(uid int) (resData vboxRep.LineChartData, err error) {
	query := `
        SELECT coalesce(sum(cost),0) as totalIncome
		FROM vbox_pay_order
		WHERE  uid = ? and order_status = ? and  DATE(created_at) = ? and c_channel_id= ?;
    `
	queryB := `
		select DATE_FORMAT(dt, '%Y-%m-%d') as dt from (
		    SELECT distinct DATE(created_at)  as dt
			FROM vbox_pay_order
			WHERE  uid = ?  and DATE(created_at) >= (CURDATE() - INTERVAL 30 DAY)
		)t
        order by dt;
    `
	queryC := `
		SELECT distinct channel
		FROM vbox_channel_shop
		WHERE  uid = ?
		;
    `
	rowCh, err := global.GVA_DB.Raw(queryC, uid).Rows()
	var channels []string
	if err != nil {
		return
	}
	for rowCh.Next() {
		var ch string
		scanErr := rowCh.Scan(&ch)
		if scanErr != nil {
			return
		}
		channels = append(channels, ch)
	}
	resData.LegendData = channels

	rowDt, err := global.GVA_DB.Raw(queryB, uid).Rows()
	var dts []string
	if err != nil {
		return
	}
	for rowDt.Next() {
		var dt string
		scanErr := rowDt.Scan(&dt)
		if scanErr != nil {
			return
		}
		dts = append(dts, dt)
	}
	resData.XData = dts

	for _, channelId := range channels {

		var income []int
		for _, dt := range dts {
			result := 0
			err = global.GVA_DB.Model(&vbox.VboxPayOrder{}).Raw(query, uid, 1, dt, channelId).Scan(&result).Error
			if err != nil {
				return
			}
			income = append(income, result)
		}
		entity := vboxRep.LineChartDataYSeries{
			Name:   channelId,
			Type:   "line",
			Smooth: true,
			Data:   income,
		}
		resData.Lists = append(resData.Lists, entity)
	}

	return resData, err
}

func (vpoService *VboxPayOrderService) GetHomePagePayOrderAnalysis(uid uint, idList []int) (resData vboxRep.VboxUserOrderPayAnalysis, err error) {
	list, total, err := vpoService.GetVboxUserPayOrderAnalysis(uid, idList)
	if total > 0 {
		resData = list[0]
	} else {
		global.GVA_LOG.Error("为获取到数据")
		return
	}
	return resData, err
}
func (vpoService *VboxPayOrderService) GetHomePagePayOrderAnalysisH(uid uint, idList []int) (resData vboxRep.VboxUserOrderPayAnalysisH, err error) {
	list, total, err := vpoService.GetVboxUserPayOrderAnalysisH(uid, idList)
	if total > 0 {
		resData = list[0]
	} else {
		global.GVA_LOG.Error("为获取到数据")
		return
	}
	return resData, err
}

// 判断是否是今天
func isToday(createTime *time.Time, now time.Time) bool {
	year, month, day := now.Date()
	createYear, createMonth, createDay := createTime.Date()

	return year == createYear && month == createMonth && day == createDay
}

// 判断是否是昨天
func isYesterday(createTime *time.Time, now time.Time) bool {
	yesterday := now.AddDate(0, 0, -1)
	year, month, day := yesterday.Date()
	createYear, createMonth, createDay := createTime.Date()

	return year == createYear && month == createMonth && day == createDay
}

func GetVboxUserWalletAvailablePoints(uid uint) (rechargeData int, err error) {

	queryA := `
        SELECT coalesce(sum(recharge),0) as recharge
		FROM vbox_user_wallet
		WHERE  uid = ? ;
    `

	queryB := `
        SELECT coalesce(sum(cost),0) as recharge
		FROM vbox_pay_order
		WHERE  uid = ? and order_status = ?;
    `

	resultA := 0
	err = global.GVA_DB.Model(&vbox.VboxUserWallet{}).Raw(queryA, uid).Scan(&resultA).Error
	if err != nil {
		return
	}

	resultB := 0
	err = global.GVA_DB.Model(&vbox.VboxPayOrder{}).Raw(queryB, uid, 1).Scan(&resultB).Error
	if err != nil {
		return
	}
	rechargeData = resultA - resultB
	return rechargeData, err

}
