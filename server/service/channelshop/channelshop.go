package channelshop

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/channelshop"
	channelshopReq "github.com/flipped-aurora/gin-vue-admin/server/model/channelshop/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"gorm.io/gorm"
	"strconv"
)

type ChannelShopService struct {
}

// CreateChannelShop 创建ChannelShop记录
// Author [piexlmax](https://github.com/piexlmax)
func (chShopService *ChannelShopService) CreateChannelShop(chShop *channelshop.ChannelShop) (err error) {
	err = global.GVA_DB.Create(chShop).Error
	return err
}

// DeleteChannelShop 删除ChannelShop记录
// Author [piexlmax](https://github.com/piexlmax)
func (chShopService *ChannelShopService) DeleteChannelShop(chShop channelshop.ChannelShop) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&channelshop.ChannelShop{}).Where("id = ?", chShop.ID).Update("deleted_by", chShop.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&chShop).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteChannelShopByIds 批量删除ChannelShop记录
// Author [piexlmax](https://github.com/piexlmax)
func (chShopService *ChannelShopService) DeleteChannelShopByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&channelshop.ChannelShop{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&channelshop.ChannelShop{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannelShop 更新ChannelShop记录
// Author [piexlmax](https://github.com/piexlmax)
func (chShopService *ChannelShopService) UpdateChannelShop(chShop channelshop.ChannelShop) (err error) {
	err = global.GVA_DB.Save(&chShop).Error
	return err
}

// GetChannelShop 根据id获取ChannelShop记录
// Author [piexlmax](https://github.com/piexlmax)
func (chShopService *ChannelShopService) GetChannelShop(id uint) (chShop channelshop.ChannelShop, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&chShop).Error
	return
}

// GetChannelShopInfoListByChanelRemark 获取同一通道下同一店铺的ChannelShop列表
// Author yoga
func (chShopService *ChannelShopService) GetChannelShopInfoListByChanelRemark(info channelshopReq.ChannelShopSearch, userId uint) (list []channelshop.ChannelShop, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	channel := info.Channel
	shopRemark := info.Shop_remark
	// 创建db
	db := global.GVA_DB.Model(&channelshop.ChannelShop{})
	var chShops []channelshop.ChannelShop
	// 如果有条件搜索 下方会自动创建搜索语句
	if userId != 1 {
		db = db.Where("uid = ?", userId)
	} else {
		db = db.Where("1 = 1", userId)
	}
	db = db.Where("channel = ? and shop_remark = ?", channel, shopRemark)
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&chShops).Error
	return chShops, total, err
}

// GetTreeChannelShopInfoList 分页获取ChannelShop记录
// Author yoga
func (chShopService *ChannelShopService) GetTreeChannelShopInfoList(info channelshopReq.ChannelShopSearch, userId uint) (list []channelshop.ChannelShop, total int64, err error) {
	//limit := info.PageSize
	//offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&channelshop.ChannelShop{})
	var chShops []channelshop.ChannelShop
	// 查出所有该用户下的数据
	if userId != 1 {
		db = db.Where("uid = ?", userId)
	} else {
		db = db.Where("1 = 1", userId)
	}

	err = db.Find(&chShops).Error
	if err != nil {
		return
	}

	// 返回查询该用户下的每个channel下的最大的一条数据
	queryP := `SELECT id,created_at,uid,cid,channel,shop_remark,address,money,status
		FROM vbox_channel_shop
		WHERE (uid = ?) AND (id in (
			SELECT MAX(id)
			FROM vbox_channel_shop
			WHERE uid = ?
			GROUP BY channel
		));`
	var chShopsP []channelshop.ChannelShop
	db = global.GVA_DB.Raw(queryP, userId, userId).Find(&chShopsP)
	fmt.Println(len(chShopsP))

	// 返回查询该用户下的相同channel下的数据
	queryC := `
        SELECT id,created_at,uid,cid,channel,shop_remark,address,money,status
		FROM vbox_channel_shop
		WHERE channel = ? and uid = ?;
    `

	queryD := `
		SELECT id,created_at,uid,cid,channel,shop_remark,address,money,status
		FROM vbox_channel_shop
		WHERE (uid = ?)  AND (id in (
			SELECT MAX(id)
			FROM vbox_channel_shop
			WHERE uid = ?
			GROUP BY channel,shop_remark
		)) AND channel = ?
		;
	`

	queryE := `
        SELECT id,created_at,uid,cid,channel,shop_remark,address,money,status
		FROM vbox_channel_shop
		WHERE channel = ? and uid = ? and shop_remark = ?;
    `

	for i, shop := range chShopsP {
		channel := shop.Channel
		var chShopsC []channelshop.ChannelShop
		err = global.GVA_DB.Raw(queryC, channel, userId).Find(&chShopsC).Error
		if err != nil {
			return
		}
		var chShopsD []channelshop.ChannelShop
		err = global.GVA_DB.Raw(queryD, userId, userId, channel).Find(&chShopsD).Error
		for i, shopMark := range chShopsD {
			remark := shopMark.Shop_remark

			var chShopsE []channelshop.ChannelShop
			err = global.GVA_DB.Raw(queryE, channel, userId, remark).Find(&chShopsE).Error
			if err != nil {
				return
			}
			ukMoney, openCnt, totalCnt := GetUniqueMoneyString(chShopsE)
			output := fmt.Sprintf("已启动【%d】个,共【%d】个", openCnt, totalCnt)
			chShopsD[i].DisMoney = ukMoney
			chShopsD[i].OpenAndClose = output
		}
		if err != nil {
			return
		}

		result, resultOpenCnt, resultTotalCnt := GetUniqueMoneyString(chShopsC)
		resultOutput := fmt.Sprintf("已启动【%d】个,共【%d】个", resultOpenCnt, resultTotalCnt)
		//fmt.Println(result)
		// 将 chShopsC 赋值给 chShopsP[i] 的 Children 字段
		chShopsP[i].Children = chShopsD
		chShopsP[i].DisMoney = result
		chShopsP[i].OpenAndClose = resultOutput
		fmt.Println("dismoney=", chShopsP[i].DisMoney)
	}

	total = int64(len(chShopsP))
	fmt.Println("total=", total)
	//err = db.Limit(limit).Offset(offset).Find(&chShopsP).Error
	return chShopsP, total, err
}

// GetChannelShopInfoList 分页获取ChannelShop记录
// Author [piexlmax](https://github.com/piexlmax)
func (chShopService *ChannelShopService) GetChannelShopInfoList(info channelshopReq.ChannelShopSearch, userId uint) (list []channelshop.ChannelShop, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&channelshop.ChannelShop{})
	var chShops []channelshop.ChannelShop
	// 如果有条件搜索 下方会自动创建搜索语句
	if userId != 1 {
		db = db.Where("uid = ?", userId)
	} else {
		db = db.Where("1 = 1", userId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&chShops).Error
	return chShops, total, err
}
func GetUniqueMoneyString(chShopsC []channelshop.ChannelShop) (string, int, int) {
	uniqueMoney := make(map[int]struct{})
	var result string
	openCnt := 0
	totalCnt := 0

	for _, shop := range chShopsC {
		money := *shop.Money

		if _, ok := uniqueMoney[money]; !ok {
			uniqueMoney[money] = struct{}{}
			if result == "" {
				result = strconv.Itoa(money)
			} else {
				result += "-" + strconv.Itoa(money)
			}
		}

		if *shop.Status == 1 {
			openCnt++
		}
		totalCnt++
	}

	return result, openCnt, totalCnt
}
