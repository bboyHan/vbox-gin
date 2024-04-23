package vbox

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxResp "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/flipped-aurora/gin-vue-admin/server/mq"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/service/vbox/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	http2 "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"github.com/gin-gonic/gin"
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
func (channelShopService *ChannelShopService) CreateChannelShop(channelShop *vboxReq.ChannelShop, cc *gin.Context) (err error) {
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
		case "1102": //jw qb tb
			flag = utils.ValidTBUrl(c.Address)
		case "1103": //jw qb tb
			flag = utils.ValidTBUrl(c.Address)
		case "1104": //jw qb tb
			flag = utils.ValidTBUrl(c.Address)
		case "1105": //jw qb tb
			flag = utils.ValidTBUrl(c.Address)

		case "1201": //dnf tb
			flag = utils.ValidTBUrl(c.Address)
		case "1202": //dnf jd
			flag = utils.ValidJDUrl(c.Address)
		case "2001": //j3 tb
			flag = utils.ValidTBUrl(c.Address)
		case "4001": //sdo tb
			flag = utils.ValidTBUrl(c.Address)
		case "5001": //qn tb
			flag = utils.ValidTBUrl(c.Address)
		case "6001": //ec jd
			flag = utils.ValidJDUrl(c.Address)

		case "7001": //qn tb
			flag = utils.ValidTBUrl(c.Address)

		case "8001": //抖币 jd
			flag = utils.ValidJDUrl(c.Address)
		case "9001": //网易 tb
			flag = utils.ValidTBUrl(c.Address)
		case "9002": //网易 pdd
			flag = utils.ValidPddUrl(c.Address)
		case "9003": //网易 jd
			flag = utils.ValidJDUrl(c.Address)
		}

		if !flag {
			return fmt.Errorf("传入的地址不合法, %s", c.Address)
		}
	}

	conn, err := mq.MQ.ConnPool.GetConnection()
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)

	ch, err := conn.Channel()
	if err != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
	}
	body := http2.DoGinContextBody(cc)

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

		if channelShop.Cid != "5001" {
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
		} else {
			if c.Status == 1 {
				var sDB vbox.ChannelShop
				err = global.GVA_DB.Model(vbox.ChannelShop{}).Where("id =?", shopDB.ID).First(&sDB).Error

				oc := vboxReq.ChanQNShopAndCtx{
					Obj: sDB,
					Ctx: vboxReq.Context{
						Body:      string(body),
						ClientIP:  cc.ClientIP(),
						Method:    cc.Request.Method,
						UrlPath:   cc.Request.URL.Path,
						UserAgent: cc.Request.UserAgent(),
						UserID:    int(sDB.CreatedBy),
					},
				}

				marshal, _ := json.Marshal(oc)

				err = ch.Publish(task.ChanQNShopEnableCheckExchange, task.ChanQNShopEnableCheckKey, marshal)

			}

		}

	}

	if err == nil {
		//conn, errC := mq.MQ.ConnPool.GetConnection()
		//if errC != nil {
		//	global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
		//}
		//defer mq.MQ.ConnPool.ReturnConnection(conn)
		//ch, errN := conn.Channel()
		//if errN != nil {
		//	global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
		//}

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
func (channelShopService *ChannelShopService) DeleteChannelShop(channelShop vbox.ChannelShop, c *gin.Context) (err error) {
	var shopDB vbox.ChannelShop
	conn, errC := mq.MQ.ConnPool.GetConnection()
	if errC != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)
	ch, errN := conn.Channel()
	if errN != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.ID).First(&shopDB)

		orgTmp := utils2.GetSelfOrg(shopDB.CreatedBy)

		if shopDB.Cid != "5001" {
			key := fmt.Sprintf(global.ChanOrgShopAddrZSet, orgTmp[0], shopDB.Cid, shopDB.Money)
			keyMem := fmt.Sprintf("%s_%v", shopDB.ProductId, shopDB.ID)
			global.GVA_REDIS.ZRem(context.Background(), key, keyMem)

			if err := tx.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.ID).Update("deleted_by", channelShop.DeletedBy).Error; err != nil {
				return err
			}
			if err = tx.Delete(&channelShop).Error; err != nil {
				return err
			}

		} else {

			body := http2.DoGinContextBody(c)

			oc := vboxReq.ChanQNShopAndCtx{
				Obj: shopDB,
				Ctx: vboxReq.Context{
					Body:      string(body),
					ClientIP:  c.ClientIP(),
					Method:    c.Request.Method,
					UrlPath:   c.Request.URL.Path,
					UserAgent: c.Request.UserAgent(),
					UserID:    int(channelShop.DeletedBy),
				},
			}
			marshal, _ := json.Marshal(oc)

			err = ch.Publish(task.ChanQNShopDelCheckExchange, task.ChanQNShopDelCheckKey, marshal)
		}

		return nil
	})

	if err == nil {

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
func (channelShopService *ChannelShopService) DeleteChannelShopByIds(ids request.IdsReq, c *gin.Context, deletedBy uint) (err error) {
	var cidList []string
	conn, errC := mq.MQ.ConnPool.GetConnection()
	if errC != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
	}
	defer mq.MQ.ConnPool.ReturnConnection(conn)
	ch, errN := conn.Channel()
	if errN != nil {
		global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
	}

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

		var shopDBList []vbox.ChannelShop
		tx.Model(&vbox.ChannelShop{}).Where("id in ?", ids).Find(&shopDBList)

		for _, shopDB := range shopDBList {
			orgTmp := utils2.GetSelfOrg(shopDB.CreatedBy)

			if shopDB.Cid != "5001" {
				key := fmt.Sprintf(global.ChanOrgShopAddrZSet, orgTmp[0], shopDB.Cid, shopDB.Money)
				keyMem := fmt.Sprintf("%s_%v", shopDB.ProductId, shopDB.ID)
				global.GVA_REDIS.ZRem(context.Background(), key, keyMem)
				cidList = append(cidList, shopDB.Cid)
			} else {
				body := http2.DoGinContextBody(c)

				oc := vboxReq.ChanQNShopAndCtx{
					Obj: shopDB,
					Ctx: vboxReq.Context{
						Body:      string(body),
						ClientIP:  c.ClientIP(),
						Method:    c.Request.Method,
						UrlPath:   c.Request.URL.Path,
						UserAgent: c.Request.UserAgent(),
						UserID:    int(deletedBy),
					},
				}
				marshal, _ := json.Marshal(oc)

				err = ch.Publish(task.ChanQNShopDelCheckExchange, task.ChanQNShopDelCheckKey, marshal)
			}

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
func (channelShopService *ChannelShopService) UpdateChannelShop(channelShop vboxReq.ChannelShopReq, c *gin.Context) (err error) {
	// 1-更新店名 2-开关单条 3-开关整个店
	if channelShop.Type == 1 {
		err = global.GVA_DB.Model(&vbox.ChannelShop{}).Where("product_id = ?", channelShop.ProductId).
			Update("shop_remark", channelShop.ShopRemark).
			Update("updated_by", channelShop.UpdatedBy).Error
	} else if channelShop.Type == 2 {

		conn, err := mq.MQ.ConnPool.GetConnection()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
		}
		defer mq.MQ.ConnPool.ReturnConnection(conn)

		ch, err := conn.Channel()
		if err != nil {
			global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
		}

		err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

			var shopDB vbox.ChannelShop
			tx.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.ID).First(&shopDB)
			orgTmp := utils2.GetSelfOrg(shopDB.CreatedBy)

			if shopDB.Cid != "5001" {
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

			} else {

				body := http2.DoGinContextBody(c)
				shopDB.Status = channelShop.Status
				oc := vboxReq.ChanQNShopAndCtx{
					Obj: shopDB,
					Ctx: vboxReq.Context{
						Body:      string(body),
						ClientIP:  c.ClientIP(),
						Method:    c.Request.Method,
						UrlPath:   c.Request.URL.Path,
						UserAgent: c.Request.UserAgent(),
						UserID:    int(shopDB.CreatedBy),
					},
				}
				marshal, _ := json.Marshal(oc)

				err = ch.Publish(task.ChanQNShopEnableCheckExchange, task.ChanQNShopEnableCheckKey, marshal)
			}

			err = tx.Model(&vbox.ChannelShop{}).Where("id = ?", channelShop.ID).
				Update("status", channelShop.Status).
				Update("updated_by", channelShop.UpdatedBy).Error

			return nil
		})

		if err == nil {

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

			conn, err := mq.MQ.ConnPool.GetConnection()
			if err != nil {
				global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
			}
			defer mq.MQ.ConnPool.ReturnConnection(conn)

			for _, shopDB := range shopDBList {
				orgTmp := utils2.GetSelfOrg(shopDB.CreatedBy)

				if shopDB.Cid != "5001" {
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
				} else {
					ch, err := conn.Channel()
					if err != nil {
						global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
						continue
					}
					body := http2.DoGinContextBody(c)
					shopDB.Status = channelShop.Status

					oc := vboxReq.ChanQNShopAndCtx{
						Obj: shopDB,
						Ctx: vboxReq.Context{
							Body:      string(body),
							ClientIP:  c.ClientIP(),
							Method:    c.Request.Method,
							UrlPath:   c.Request.URL.Path,
							UserAgent: c.Request.UserAgent(),
							UserID:    int(shopDB.CreatedBy),
						},
					}
					marshal, err := json.Marshal(oc)

					err = ch.Publish(task.ChanQNShopEnableCheckExchange, task.ChanQNShopEnableCheckKey, marshal)
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

	num := 0
	for productID, v := range shopMap {
		// 将获取的值进行反序列化到 ShopIncomeResp 类型
		var item vboxResp.ShopIncomeResp
		key := "statis:chShop:" + productID
		if num == 0 {

			result, err := global.GVA_REDIS.Get(context.Background(), key).Result()
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println("result=", result == "")
			if result == "" {
				getShopOkstatisResp(ids)
			}

		}

		resultExists, errExists := global.GVA_REDIS.Get(context.Background(), key).Result()
		if errExists != nil {
			fmt.Println(errExists)
		}
		//fmt.Println("resultExists=", resultExists)
		err = json.Unmarshal([]byte(resultExists), &item)
		if err != nil {
			fmt.Println(err)
		}

		if len(v) == 0 {
			continue
		}
		csNew := vboxReq.ChannelShop{
			ProductId:       productID,
			ShopRemark:      v[0].ShopRemark,
			Cid:             v[0].Cid,
			ChannelShopList: []vboxReq.ChannelShopSub{},
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           item.Ratio,
			OkIncome:        item.OkIncome,
		}

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
		num++
	}

	sort.Sort(vboxReq.ChannelShopList(res))

	return res, err
}

func getShopOkstatisResp(ids []uint) (err error) {
	querySql := `
		SELECT
			c.product_id as shopId,
			c.shop_remark as shopName,
			COALESCE (a.cnt,0) as orderQuantify,
			COALESCE (a.ok_cnt,0) as okOrderQuantify,
			0 as ratio,
			COALESCE (a.ok_income,0) as okIncome
		FROM
		(
			SELECT
			 event_id,
			 sum(if (order_status=1 and cb_status=1,money,0)) as ok_income,
			 count(*) as cnt,
			 sum(if(order_status=1 and cb_status=1,1,0)) as ok_cnt
			FROM(
					SELECT
							substring_index(event_id,'_',1) as event_id,
							order_status,
							cb_status,
							money,
							created_at,
							cb_time	
					from vbox_pay_order 
					where DATE_FORMAT(cb_time, '%Y-%m-%d') = ? or DATE_FORMAT(created_at, '%Y-%m-%d') = ?
					and event_type = 1 and created_by in  ?
			)a1
			GROUP BY event_id
		)a 
		right join (
		SELECT DISTINCT product_id,shop_remark
		FROM vbox_channel_shop 
		) c
		on a.event_id = c.product_id  
	`
	//shopMap := make(map[string]vboxResp.ShopIncomeResp)
	dt := time.Now().AddDate(0, 0, 0).Format("2006-01-02")

	//fmt.Println("dt-->", dt, "uid-->", uid, "querySql-->", querySql)
	db := global.GVA_DB.Model(&vboxResp.ShopIncomeResp{})
	rows, err := db.Raw(querySql, dt, dt, ids).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.ShopIncomeResp
		err := rows.Scan(&item.ShopId, &item.ShopName, &item.OrderQuantify, &item.OkOrderQuantify, &item.Ratio, &item.OkIncome)
		if err != nil {
			panic(err)
		}
		defaultItem := vboxResp.ShopIncomeResp{
			ShopId:          item.ShopId,
			ShopName:        item.ShopName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           item.Ratio,
			OkIncome:        item.OkIncome,
		}

		//shopMap[item.ShopId] = defaultItem
		key := "statis:chShop:" + item.ShopId
		jsonStr, err := json.Marshal(defaultItem)
		//fmt.Println("redis set,", jsonStr)
		err = global.GVA_REDIS.Set(context.Background(), key, jsonStr, 5*time.Minute).Err()
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return err
}
