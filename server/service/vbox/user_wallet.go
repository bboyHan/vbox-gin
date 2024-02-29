package vbox

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxRep "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	utils2 "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"gorm.io/gorm"
)

type UserWalletService struct {
}

// CreateUserWallet 创建用户钱包记录
// Author [piexlmax](https://github.com/piexlmax)
func (userWalletService *UserWalletService) CreateUserWallet(userWallet *vbox.UserWallet) (err error) {
	err = global.GVA_DB.Create(userWallet).Error
	return err
}

// TransferUserWallet 划转积分给其它用户
func (userWalletService *UserWalletService) TransferUserWallet(userWalletTransfer *vboxReq.UserWalletTransfer) (err error) {

	// 1.校验是否给同组织成员划转
	uid := userWalletTransfer.CurrentUid
	username := userWalletTransfer.Username
	recharge := userWalletTransfer.Recharge
	toUid := userWalletTransfer.ToUid
	toUsername := userWalletTransfer.ToUsername

	orgUserIds := utils2.GetDeepUserIDs(uid)
	isExist := utils.Contains(orgUserIds, toUid)
	if !isExist {
		return errors.New("不允许给非当前团队成员划转/充值积分")
	}

	if recharge > 0 {
		// 2.
		switch userWalletTransfer.Type {
		// Type: 1 直充 (只允许超管操作)
		case global.WalletRechargeType:

			// uid 888 (超管)
			var roleID uint
			roleID, err = utils.GetRoleID(uid)
			if roleID == 888 || roleID == 1000 {
			} else {
				return errors.New("该账号无直充权限")
			}

			eventId := utils.Prefix(global.WalletEventRechargePrefix, rand_string.RandomInt(8))
			err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
				// 直充
				rowSelf := &vbox.UserWallet{
					Uid:       toUid,
					Username:  toUsername,
					Recharge:  recharge,
					Type:      global.WalletRechargeType,
					EventId:   eventId,
					Remark:    fmt.Sprintf(global.WalletEventRecharge, recharge),
					CreatedBy: toUid,
				}
				if err := tx.Model(&vbox.UserWallet{}).Create(rowSelf).Error; err != nil {
					return err
				}

				return nil
			})
		// Type: 2 划转 -> 入库（1条-本用户扣减，1条-转用户加分）
		case global.WalletTransferType:

			// 检查当前余额够不够
			var balance int
			balance, err = userWalletService.GetUserWalletSelf(uid)
			if balance <= 0 || balance < recharge {
				return errors.New("余额不足")
			}

			eventId := utils.Prefix(global.WalletEventTransferPrefix, rand_string.RandomInt(8))

			err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

				// 扣减
				rowSelf := &vbox.UserWallet{
					Uid:       uid,
					Username:  username,
					Recharge:  -recharge,
					Type:      global.WalletTransferType,
					EventId:   eventId,
					Remark:    fmt.Sprintf(global.WalletEventTransfer, recharge, toUsername),
					CreatedBy: uid,
				}
				if err := tx.Model(&vbox.UserWallet{}).Create(rowSelf).Error; err != nil {
					return err
				}

				// 划转至新账户
				rowTo := &vbox.UserWallet{
					Uid:       toUid,
					Username:  toUsername,
					Recharge:  recharge,
					Type:      global.WalletTransferType,
					EventId:   eventId,
					Remark:    fmt.Sprintf(global.WalletEventIncome, recharge, username),
					CreatedBy: uid,
				}

				if err := tx.Model(&vbox.UserWallet{}).Create(rowTo).Error; err != nil {
					return err
				}
				return nil
			})
		case global.WalletOrderType: // 2- 订单积分扣费
		default:
			return errors.New("不支持的充值类型")
		}
	} else if recharge < 0 {
		// 2.
		switch userWalletTransfer.Type {
		// Type: 1 直充 (只允许超管操作)
		case global.WalletRechargeType:

			// uid 888 (超管)
			var roleID uint
			roleID, err = utils.GetRoleID(uid)
			if roleID == 888 || roleID == 1000 {
			} else {
				return errors.New("该账号无直充权限")
			}

			eventId := utils.Prefix(global.WalletEventRechargePrefix, rand_string.RandomInt(8))
			err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
				// 直充
				rowSelf := &vbox.UserWallet{
					Uid:       toUid,
					Username:  toUsername,
					Recharge:  recharge,
					Type:      global.WalletRechargeType,
					EventId:   eventId,
					Remark:    fmt.Sprintf(global.WalletEventRecharge, recharge),
					CreatedBy: toUid,
				}
				if err := tx.Model(&vbox.UserWallet{}).Create(rowSelf).Error; err != nil {
					return err
				}

				return nil
			})
		// Type: 2 划转 -> 入库（1条-本用户扣减，1条-转用户加分）
		case global.WalletTransferType:

			revRecharge := -recharge //正

			// uid 888 (超管)
			var roleID uint
			roleID, err = utils.GetRoleID(uid)
			if roleID == 888 || roleID == 1000 || roleID == 1001 {
			} else {
				return errors.New("该账号无扣减权限")
			}

			// 检查对方余额够不够
			var balance int
			balance, err = userWalletService.GetUserWalletSelf(toUid)
			if balance <= 0 || balance < revRecharge {
				return errors.New("余额不足")
			}

			eventId := utils.Prefix(global.WalletEventTransferPrefix, rand_string.RandomInt(8))

			err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

				// 扣减
				rowSelf := &vbox.UserWallet{
					Uid:       uid,
					Username:  username,
					Recharge:  -recharge,
					Type:      global.WalletTransferType,
					EventId:   eventId,
					Remark:    fmt.Sprintf(global.WalletEventTransfer, recharge, toUsername),
					CreatedBy: uid,
				}
				if err := tx.Model(&vbox.UserWallet{}).Create(rowSelf).Error; err != nil {
					return err
				}

				// 划转至新账户
				rowTo := &vbox.UserWallet{
					Uid:       toUid,
					Username:  toUsername,
					Recharge:  recharge,
					Type:      global.WalletTransferType,
					EventId:   eventId,
					Remark:    fmt.Sprintf(global.WalletEventIncome, recharge, username),
					CreatedBy: uid,
				}

				if err := tx.Model(&vbox.UserWallet{}).Create(rowTo).Error; err != nil {
					return err
				}
				return nil
			})
		case global.WalletOrderType: // 2- 订单积分扣费
		default:
			return errors.New("不支持的充值类型")
		}
	}

	return err
}

// DeleteUserWallet 删除用户钱包记录
// Author [piexlmax](https://github.com/piexlmax)
func (userWalletService *UserWalletService) DeleteUserWallet(userWallet vbox.UserWallet) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.UserWallet{}).Where("id = ?", userWallet.ID).Update("deleted_by", userWallet.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&userWallet).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteUserWalletByIds 批量删除用户钱包记录
// Author [piexlmax](https://github.com/piexlmax)
func (userWalletService *UserWalletService) DeleteUserWalletByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.UserWallet{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.UserWallet{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateUserWallet 更新用户钱包记录
// Author [piexlmax](https://github.com/piexlmax)
func (userWalletService *UserWalletService) UpdateUserWallet(userWallet vbox.UserWallet) (err error) {
	err = global.GVA_DB.Save(&userWallet).Error
	return err
}

// GetUserWallet 根据id获取用户钱包记录
// Author [piexlmax](https://github.com/piexlmax)
func (userWalletService *UserWalletService) GetUserWallet(id uint) (userWallet vbox.UserWallet, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&userWallet).Error
	return
}

// GetUserWalletSelf 获取当前用户钱包余额
func (userWalletService *UserWalletService) GetUserWalletSelf(id uint) (balance int, err error) {
	err = global.GVA_DB.Model(&vbox.UserWallet{}).Select("IFNULL(sum(recharge),0) as balance").Where("uid = ?", id).First(&balance).Error
	return
}

// GetUserWalletCostOV 获取指定用户3日内消费情况
func (userWalletService *UserWalletService) GetUserWalletCostOV(info vboxReq.UserWalletSearch, id uint) (ov []vboxRep.DataWalletOverView, err error) {
	//err = global.GVA_DB.Model(&vbox.UserWallet{}).Select("IFNULL(sum(recharge),0) as balance").Where("uid = ?", id).First(&balance).Error
	db := global.GVA_DB.Model(&vbox.UserWallet{})
	if info.ToUid == 0 {
		//未传值，则默认查自己
		info.ToUid = id
	}
	orgUserIds := utils2.GetDeepUserIDs(id)
	isExist := utils.Contains(orgUserIds, info.ToUid)
	if !isExist {
		return nil, errors.New("不允许给非当前团队的积分数据")
	}

	err = db.Select(
		`uid AS x0,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() - INTERVAL 3 DAY AND recharge < 0 AND type = 3 THEN recharge ELSE 0 END), 0) AS x1,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() - INTERVAL 2 DAY AND recharge < 0 AND type = 3 THEN recharge ELSE 0 END), 0) AS x2,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() - INTERVAL 1 DAY AND recharge < 0 AND type = 3 THEN recharge ELSE 0 END), 0) AS x3,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() AND recharge < 0 AND type = 3 THEN recharge ELSE 0 END), 0) AS x4,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() - INTERVAL 1 DAY AND recharge > 0 AND type in (1,2) THEN recharge ELSE 0 END), 0) AS x5,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() - INTERVAL 1 DAY AND recharge < 0 AND type in (1,2) THEN recharge ELSE 0 END), 0) AS x6,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() AND recharge > 0 AND type in (1,2) THEN recharge ELSE 0 END), 0) AS x7,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() AND recharge < 0 AND type in (1,2) THEN recharge ELSE 0 END), 0) AS x8,
    IFNULL(SUM(recharge), 0) AS x9`).
		Where("uid = ?", info.ToUid).Find(&ov).Error
	if err != nil {
		return nil, err
	}

	return ov, err
}

// GetUserWalletOverview 获取当前用户钱包余额
func (userWalletService *UserWalletService) GetUserWalletOverview(info vboxReq.UserWalletSearch, ids []uint) (ov []vboxRep.DataWalletOverView, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.UserWallet{})

	err = db.Select(
		`uid AS x0,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() - INTERVAL 2 DAY AND recharge > 0 THEN recharge ELSE 0 END), 0) AS x1,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() - INTERVAL 2 DAY AND recharge < 0 THEN recharge ELSE 0 END), 0) AS x2,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() - INTERVAL 1 DAY AND recharge > 0 THEN recharge ELSE 0 END), 0) AS x3,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() - INTERVAL 1 DAY AND recharge < 0 THEN recharge ELSE 0 END), 0) AS x4,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() AND recharge > 0 THEN recharge ELSE 0 END), 0) AS x5,
    IFNULL(SUM(CASE WHEN DATE(created_at) = CURDATE() AND recharge < 0 THEN recharge ELSE 0 END), 0) AS x6,
    IFNULL(SUM(CASE WHEN recharge > 0 THEN recharge ELSE 0 END), 0) AS x7,
    IFNULL(SUM(CASE WHEN recharge < 0 THEN recharge ELSE 0 END), 0) AS x8,
    IFNULL(SUM(recharge), 0) AS x9`).
		Where("uid in ?", ids).Group("x0").Find(&ov).Error
	if err != nil {
		return
	}

	return ov, err

}

// GetUserWalletInfoList 分页获取用户钱包记录
// Author [piexlmax](https://github.com/piexlmax)
func (userWalletService *UserWalletService) GetUserWalletInfoList(info vboxReq.UserWalletSearch, ids []uint) (list []vbox.UserWallet, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.UserWallet{})
	var userWallets []vbox.UserWallet
	// 如果有条件搜索 下方会自动创建搜索语句
	db.Where("uid in ?", ids)
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Username != "" {
		db = db.Where("username LIKE ?", "%"+info.Username+"%")
	}
	if info.Type != 0 {
		db = db.Where("type = ?", info.Type)
	}
	if info.Remark != "" {
		db = db.Where("remark LIKE ?", "%"+info.Remark+"%")
	}

	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Debug().Order("id desc").Find(&userWallets).Error

	return userWallets, total, err
}
