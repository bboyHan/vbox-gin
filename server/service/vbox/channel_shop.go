package vbox

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/redis/go-redis/v9"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"
	"time"
)

type ChannelShopService struct {
}

// CreateChannelShop 创建引导商铺记录
func (channelShopService *ChannelShopService) CreateChannelShop(channelShop *vboxReq.ChannelShop) (err error) {
	length := len(channelShop.ChannelShopList)
	if length < 1 {
		err = fmt.Errorf("创建的店铺地址至少包含1个")
		return err
	}
	pid := channelShop.ProductId
	if pid == "" {
		pid = time.Now().Format("20060102150405") + rand_string.RandomInt(3)
	}
	orgTmp := utils2.GetSelfOrg(channelShop.CreatedBy)

	cid := channelShop.Cid
	for i := range channelShop.ChannelShopList {
		// 获取当前元素的指针
		c := &channelShop.ChannelShopList[i]
		// 增加校验
		if c.Money <= 0 {
			return fmt.Errorf("传入的金额不合法")
		}
		if c.Address == "" {
			return fmt.Errorf("传入的地址不合法, %s", c.Address)
		} else {
			if cid == "1006" {
				global.GVA_LOG.Info("跳过解析http url合法性校验")
			} else {
				//	如果是口令地址，则先解析url
				addr, errA := utils.ParseUrlContent(c.Address)
				if errA != nil {
					return fmt.Errorf("传入的地址不合法, %s", c.Address)
				}
				c.Address = addr
			}
		}
		var flag bool
		switch cid {

		case "1001": //jd
			flag = utils.ValidJDUrl(c.Address)
		case "1002": //dy
			flag = utils.ValidDYUrl(c.Address)
		case "1003": //jym
			flag = utils.ValidAlipayUrl(c.Address)
		case "1004": //zfb
			flag = utils.ValidAlipayUrl(c.Address)
		case "1005": //qb tb
			flag = utils.ValidTBUrl(c.Address)
		case "1006": //wx xcx
			flag = utils.ValidXCXUrl(c.Address)
		case "1007": //qb pdd
			flag = utils.ValidPddUrl(c.Address)

		case "1101": //jw qb tb
			flag = utils.ValidTBUrl(c.Address)

		case "1201": //dnf tb
			flag = utils.ValidTBUrl(c.Address)
		case "1202": //dnf jd
			flag = utils.ValidJDUrl(c.Address)

		case "2001": //j3 tb
			flag = utils.ValidTBUrl(c.Address)
		case "4001": //sdo tb
			flag = utils.ValidTBUrl(c.Address)
		case "6001": //ec jd
			flag = utils.ValidJDUrl(c.Address)

		case "7001": //qn tb
			flag = utils.ValidTBUrl(c.Address)
		}

		if !flag {
			return fmt.Errorf("传入的地址不合法, %s", c.Address)
		}
	}

	for _, c := range channelShop.ChannelShopList {

		var shopDB vbox.ChannelShop
		if c.ID == 0 { //走创建
			chNew := vbox.ChannelShop{
				Cid:        channelShop.Cid,
				ProductId:  pid,
				ShopRemark: channelShop.ShopRemark,
				Address:    c.Address,
				Device:     c.Device,
				MarkId:     c.MarkId,
				Money:      c.Money,
				Status:     c.Status,
				CreatedBy:  channelShop.CreatedBy,
			}
			if err := global.GVA_DB.Model(&vbox.ChannelShop{}).Create(&chNew).Error; err != nil {
				return err
			}
			shopDB = chNew
			shopDB.ID = chNew.ID
			global.GVA_LOG.Info("创建店铺信息", zap.Any("shopDB", shopDB))

		} else { //走更新
			chNew := vbox.ChannelShop{
				Cid:        channelShop.Cid,
				ProductId:  pid,
				ShopRemark: channelShop.ShopRemark,
				Address:    c.Address,
				Device:     c.Device,
				MarkId:     c.MarkId,
				Money:      c.Money,
				Status:     c.Status,
				UpdatedBy:  channelShop.CreatedBy,
			}
			chNew.ID = c.ID
			if err := global.GVA_DB.Model(&vbox.ChannelShop{}).Where("id = ?", c.ID).Updates(&chNew).Error; err != nil {
				return err
			}
			shopDB = chNew
			global.GVA_LOG.Info("更新店铺信息", zap.Any("shopDB", shopDB))
		}

		key := fmt.Sprintf(global.ChanOrgShopAddrZSet, orgTmp[0], channelShop.Cid, c.Money)
		keyMem := fmt.Sprintf("%s_%v", shopDB.ProductId, shopDB.ID)
		if c.Status == 1 {

			global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
				Score:  float64(time.Now().Unix()), // 重新放进去，score设置最新的时间
				Member: keyMem,
			})
		} else {
			global.GVA_REDIS.ZRem(context.Background(), key, keyMem)
		}
	}

	if err == nil {
		conn, errC := mq.MQ.ConnPool.GetConnection()
		if errC != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
		}
		defer mq.MQ.ConnPool.ReturnConnection(conn)
		ch, errN := conn.Channel()
		if errN != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
		}

		orgTmpX := utils2.GetSelfOrg(channelShop.CreatedBy)

		moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmpX[0], channelShop.Cid)
		createBy := channelShop.CreatedBy
		waitMsg := fmt.Sprintf("%s-%d", moneyKey, createBy)
		err = ch.Publish(task.ChanAccShopUpdCheckExchange, task.ChanAccShopUpdCheckKey, []byte(waitMsg))
	}

	return err
}

// DeleteChannelShop 删除引导商铺记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelShopService *ChannelShopService) DeleteChannelShop(channelShop vbox.ChannelShop) (err error) {
	var shopDB vbox.ChannelShop
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.ID).First(&shopDB)

		orgTmp := utils2.GetSelfOrg(shopDB.CreatedBy)

		key := fmt.Sprintf(global.ChanOrgShopAddrZSet, orgTmp[0], shopDB.Cid, shopDB.Money)
		keyMem := fmt.Sprintf("%s_%v", shopDB.ProductId, shopDB.ID)
		global.GVA_REDIS.ZRem(context.Background(), key, keyMem)

		if err := tx.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.ID).Update("deleted_by", channelShop.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&channelShop).Error; err != nil {
			return err
		}
		return nil
	})

	if err == nil {
		conn, errC := mq.MQ.ConnPool.GetConnection()
		if errC != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
		}
		defer mq.MQ.ConnPool.ReturnConnection(conn)
		ch, errN := conn.Channel()
		if errN != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
		}

		orgTmp := utils2.GetSelfOrg(shopDB.CreatedBy)

		moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], shopDB.Cid)
		createBy := shopDB.CreatedBy
		waitMsg := fmt.Sprintf("%s-%d", moneyKey, createBy)
		err = ch.Publish(task.ChanAccShopUpdCheckExchange, task.ChanAccShopUpdCheckKey, []byte(waitMsg))
	}
	return err
}

// DeleteChannelShopByIds 批量删除引导商铺记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelShopService *ChannelShopService) DeleteChannelShopByIds(ids request.IdsReq, deletedBy uint) (err error) {
	var cidList []string
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

		var shopDBList []vbox.ChannelShop
		tx.Model(&vbox.ChannelShop{}).Where("id in ?", ids).Find(&shopDBList)

		for _, shopDB := range shopDBList {
			orgTmp := utils2.GetSelfOrg(shopDB.CreatedBy)

			key := fmt.Sprintf(global.ChanOrgShopAddrZSet, orgTmp[0], shopDB.Cid, shopDB.Money)
			keyMem := fmt.Sprintf("%s_%v", shopDB.ProductId, shopDB.ID)
			global.GVA_REDIS.ZRem(context.Background(), key, keyMem)
			cidList = append(cidList, shopDB.Cid)
		}

		if err := tx.Model(&vbox.ChannelShop{}).Where("id in ?", ids.Ids).Update("deleted_by", deletedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.ChannelShop{}).Error; err != nil {
			return err
		}
		return nil
	})

	uniqCIDs := utils2.UniqStr(cidList)

	if err == nil {
		conn, errC := mq.MQ.ConnPool.GetConnection()
		if errC != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
		}
		defer mq.MQ.ConnPool.ReturnConnection(conn)
		ch, errN := conn.Channel()
		if errN != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
		}

		orgTmp := utils2.GetSelfOrg(deletedBy)

		for _, cid := range uniqCIDs {
			moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], cid)
			createBy := deletedBy
			waitMsg := fmt.Sprintf("%s-%d", moneyKey, createBy)
			err = ch.Publish(task.ChanAccShopUpdCheckExchange, task.ChanAccShopUpdCheckKey, []byte(waitMsg))
		}
	}

	return err
}

// UpdateChannelShop 更新引导商铺记录（ type: 1-更新店名 2-开关单条 3-开关整个店 ）
func (channelShopService *ChannelShopService) UpdateChannelShop(channelShop vboxReq.ChannelShopReq) (err error) {
	// 1-更新店名 2-开关单条 3-开关整个店
	if channelShop.Type == 1 {
		err = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("product_id = ?", channelShop.ProductId).
			Update("shop_remark", channelShop.ShopRemark).
			Update("updated_by", channelShop.UpdatedBy).Error
	} else if channelShop.Type == 2 {
		err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

			var shopDB vbox.ChannelShop
			tx.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.ID).First(&shopDB)
			orgTmp := utils2.GetSelfOrg(shopDB.CreatedBy)
			key := fmt.Sprintf(global.ChanOrgShopAddrZSet, orgTmp[0], shopDB.Cid, shopDB.Money)
			keyMem := fmt.Sprintf("%s_%v", shopDB.ProductId, shopDB.ID)
			channelShop.Cid = shopDB.Cid
			if channelShop.Status == 1 {

				global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
					Score:  float64(time.Now().Unix()), // 重新放进去，score设置最新的时间
					Member: keyMem,
				})
			} else {
				global.GVA_REDIS.ZRem(context.Background(), key, keyMem)
			}

			err = tx.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.ID).
				Update("status", channelShop.Status).
				Update("updated_by", channelShop.UpdatedBy).Error

			return nil
		})

		if err == nil {
			conn, errC := mq.MQ.ConnPool.GetConnection()
			if errC != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)
			ch, errN := conn.Channel()
			if errN != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
			}

			orgTmp := utils2.GetSelfOrg(channelShop.UpdatedBy)

			moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], channelShop.Cid)
			createBy := channelShop.UpdatedBy
			waitMsg := fmt.Sprintf("%s-%d", moneyKey, createBy)
			err = ch.Publish(task.ChanAccShopUpdCheckExchange, task.ChanAccShopUpdCheckKey, []byte(waitMsg))
		}
	} else if channelShop.Type == 3 {
		err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

			var shopDBList []vbox.ChannelShop
			tx.Model(&vbox.ChannelShop{}).Where("product_id = ?", channelShop.ProductId).Find(&shopDBList)

			for _, shopDB := range shopDBList {
				orgTmp := utils2.GetSelfOrg(shopDB.CreatedBy)

				key := fmt.Sprintf(global.ChanOrgShopAddrZSet, orgTmp[0], shopDB.Cid, shopDB.Money)
				keyMem := fmt.Sprintf("%s_%v", shopDB.ProductId, shopDB.ID)

				if channelShop.Status == 1 {

					global.GVA_REDIS.ZAdd(context.Background(), key, redis.Z{
						Score:  float64(time.Now().Unix()), // 重新放进去，score设置最新的时间
						Member: keyMem,
					})
				} else {
					global.GVA_REDIS.ZRem(context.Background(), key, keyMem)
				}

				channelShop.Cid = shopDB.Cid
			}

			err = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("product_id = ?", channelShop.ProductId).
				Update("status", channelShop.Status).
				Update("updated_by", channelShop.UpdatedBy).Error

			return nil

		})

		if err == nil {
			conn, errC := mq.MQ.ConnPool.GetConnection()
			if errC != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)
			ch, errN := conn.Channel()
			if errN != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
			}

			orgTmp := utils2.GetSelfOrg(channelShop.UpdatedBy)

			moneyKey := fmt.Sprintf(global.OrgShopMoneySet, orgTmp[0], channelShop.Cid)
			createBy := channelShop.UpdatedBy
			waitMsg := fmt.Sprintf("%s-%d", moneyKey, createBy)
			err = ch.Publish(task.ChanAccShopUpdCheckExchange, task.ChanAccShopUpdCheckKey, []byte(waitMsg))
		}
	} else {
		err = fmt.Errorf("不支持的操作，type检查")
	}

	return err
}

// GetChannelShop 根据id获取引导商铺记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelShopService *ChannelShopService) GetChannelShop(id uint) (channelShop vbox.ChannelShop, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&channelShop).Error
	return
}

// GetChannelShopByProductId 根据ProductId获取引导商铺记录
func (channelShopService *ChannelShopService) GetChannelShopByProductId(productId string) (res vboxReq.ChannelShop, err error) {

	var csDB []vbox.ChannelShop
	err = global.GVA_DB.Where("product_id = ?", productId).Find(&csDB).Error

	if len(csDB) == 0 { //无记录
		return vboxReq.ChannelShop{}, nil
	}

	res = vboxReq.ChannelShop{
		ProductId:       productId,
		ShopRemark:      csDB[0].ShopRemark,
		Cid:             csDB[0].Cid,
		ChannelShopList: []vboxReq.ChannelShopSub{}}

	for _, record := range csDB {
		res.ChannelShopList = append(res.ChannelShopList, vboxReq.ChannelShopSub{
			ID:      record.ID,
			Address: record.Address,
			Money:   record.Money,
			Device:  record.Device,
			MarkId:  record.MarkId,
			Status:  record.Status,
		})
	}

	return res, nil
}

// GetChannelShopInfoList 分页获取引导商铺记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelShopService *ChannelShopService) GetChannelShopInfoList(info vboxReq.ChannelShopSearch, ids []uint) (res []vboxReq.ChannelShop, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelShop{})
	var channelShops []vbox.ChannelShop

	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.ShopRemark != "" {
		db = db.Where("shop_remark LIKE ?", "%"+info.ShopRemark+"%")
	}

	if info.ProductId != "" {
		db = db.Where("product_id = ?", info.ProductId)
	}

	err = db.Where("created_by in (?)", ids).Find(&channelShops).Error

	shopMap := make(map[string][]vbox.ChannelShop)

	for _, shop := range channelShops {
		if _, ok := shopMap[shop.ProductId]; !ok { // 如果map中没有productID这个key，则添加一个空切片
			shopMap[shop.ProductId] = []vbox.ChannelShop{}
		}
		shopMap[shop.ProductId] = append(shopMap[shop.ProductId], shop)
	}

	for productID, v := range shopMap {
		if len(v) == 0 {
			continue
		}
		csNew := vboxReq.ChannelShop{
			ProductId:       productID,
			ShopRemark:      v[0].ShopRemark,
			Cid:             v[0].Cid,
			ChannelShopList: []vboxReq.ChannelShopSub{}}

		for _, record := range v {
			csNew.ChannelShopList = append(csNew.ChannelShopList, vboxReq.ChannelShopSub{
				ID:      record.ID,
				Address: record.Address,
				Money:   record.Money,
				Device:  record.Device,
				MarkId:  record.MarkId,
				Status:  record.Status,
			})
		}

		res = append(res, csNew)
	}

	sort.Sort(vboxReq.ChannelShopList(res))

	return res, err
}
