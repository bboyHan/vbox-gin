package vbox

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	http2 "github.com/flipped-aurora/gin-vue-admin/server/utils/http"
	"github.com/gin-gonic/gin"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChannelCardAccService struct {
}

// CreateChannelCardAcc 创建通道账号记录
func (cardAccService *ChannelCardAccService) CreateChannelCardAcc(vca *vbox.ChannelCardAcc, c *gin.Context) (err error) {
	vca.AcId = rand_string.RandomInt(8)
	token := vca.Token
	//增加校验
	if global.ECContains(vca.Cid) {
		isCK := http2.IsValidCookie(token)
		if !isCK {
			return errors.New("ck信息不合法")
		}
	} else {
		return errors.New("该信道暂不支持创建账号")
	}

	if vca.AcAccount == "" {
		pin := http2.ParseCookie(token, "pin")
		if pin != "" {
			vca.AcAccount = pin
		}
	}

	if vca.AcAccount == "" {
		return errors.New("ck中不完整，不包含账号信息")
	} else {
		var count int64
		if err = global.GVA_DB.Debug().Model(&vbox.ChannelCardAcc{}).Where("ac_account = ? and cid = ?", vca.AcAccount, vca.Cid).Count(&count).Error; err != nil {
			global.GVA_LOG.Warn("系统查询异常", zap.Error(err))
			return errors.New("系统查询异常，请重试或联系管理员")
		}
		if count > 0 {
			global.GVA_LOG.Warn("账号在当前通道已存在，不允许重复添加，请核实", zap.Any("ac_account", vca.AcAccount))
			return errors.New("账号在当前通道已存在，不允许重复添加，请核查")
		}
		vca.AcAccount = utils.Trim(vca.AcAccount)
	}

	err = global.GVA_DB.Create(vca).Error
	//vca传入的所有值 转化成 vcaDB vbox.ChannelCardAcc存放

	return err
}

// DeleteChannelCardAcc 删除通道账号记录
func (cardAccService *ChannelCardAccService) DeleteChannelCardAcc(vca vbox.ChannelCardAcc, c *gin.Context) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		//var vcaDB vbox.ChannelCardAcc
		//if err := tx.Model(&vbox.ChannelCardAcc{}).Where("id = ?", vca.ID).First(&vcaDB).Error; err != nil {
		//	return err
		//}

		if err := tx.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", vca.ID).Update("sys_status", 2).Error; err != nil {
			return err
		}
		if errD := global.GVA_DB.Where("id = ?", vca.ID).Delete(&vbox.ChannelCardAcc{}).Error; errD != nil {
			global.GVA_LOG.Error("账号删除过程...处理删除预产,删除预产失败", zap.Error(errD))
		}

		return nil
	})
	return err
}

// DeleteChannelCardAccByIds 批量删除通道账号记录
func (cardAccService *ChannelCardAccService) DeleteChannelCardAccByIds(ids request.IdsReq, c *gin.Context, deletedBy uint) (err error) {

	if len(ids.Ids) < 1 {
		return fmt.Errorf("传入的id为空")
	} else {
		for _, ID := range ids.Ids {
			err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
				//var vcaDB vbox.ChannelCardAcc
				//if err := tx.Model(&vbox.ChannelCardAcc{}).Where("id = ?", ID).First(&vcaDB).Error; err != nil {
				//	return err
				//}
				//
				//conn, err := mq.MQ.ConnPool.GetConnection()
				//if err != nil {
				//	global.GVA_LOG.Warn(fmt.Sprintf("Failed to get connection from pool: %v", err))
				//}
				//defer mq.MQ.ConnPool.ReturnConnection(conn)
				//
				//ch, err := conn.Channel()
				//if err != nil {
				//	global.GVA_LOG.Warn(fmt.Sprintf("new mq channel err: %v", err))
				//}
				//
				//body := http2.DoGinContextBody(c)
				//
				//oc := vboxReq.ChanAccAndCtx{
				//	Obj: vcaDB,
				//	Ctx: vboxReq.Context{
				//		Body:      string(body),
				//		ClientIP:  c.ClientIP(),
				//		Method:    c.Request.Method,
				//		UrlPath:   c.Request.URL.Path,
				//		UserAgent: c.Request.UserAgent(),
				//		UserID:    int(deletedBy),
				//	},
				//}
				//marshal, err := json.Marshal(oc)
				//
				//err = ch.Publish(task.ChanAccDelCheckExchange, task.ChanAccDelCheckKey, marshal)

				if err := tx.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", ID).Update("sys_status", 2).Error; err != nil {
					return err
				}
				if errD := global.GVA_DB.Where("id = ?", ID).Delete(&vbox.ChannelCardAcc{}).Error; errD != nil {
					global.GVA_LOG.Error("账号删除过程...处理删除预产,删除预产失败", zap.Error(errD))
				}

				return nil
			})
		}
	}

	return err
}

// SwitchEnableChannelCardAcc 开关通道账号
// Author [bboyhan](https://github.com/bboyhan)
func (cardAccService *ChannelCardAccService) SwitchEnableChannelCardAcc(vca vboxReq.ChannelCardAccUpd, c *gin.Context) (err error) {
	var vcaDB vbox.ChannelCardAcc
	err = global.GVA_DB.Where("id = ?", vca.ID).First(&vcaDB).Error
	if err != nil {
		return fmt.Errorf("不存在的账号，请核查")
	}

	err = global.GVA_DB.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id = ?", vca.ID).Update("status", vca.Status).Update("updated_by", vca.UpdatedBy).Error
	return err
}

// SwitchEnableChannelCardAccByIds 批量开关通道账号记录
// Author [bboyhan](https://github.com/bboyhan)
func (cardAccService *ChannelCardAccService) SwitchEnableChannelCardAccByIds(upd vboxReq.ChannelCardAccUpd, updatedBy uint, c *gin.Context) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {

		if err := tx.Unscoped().Model(&vbox.ChannelCardAcc{}).Where("id in ?", upd.Ids).Update("status", upd.Status).Update("updated_by", updatedBy).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", upd.Ids).Updates(&vbox.ChannelCardAcc{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannelCardAcc 更新通道账号记录
func (cardAccService *ChannelCardAccService) UpdateChannelCardAcc(vca vbox.ChannelCardAcc) (err error) {
	token := vca.Token
	var vcaDB vbox.ChannelCardAcc
	err = global.GVA_DB.Where("id = ?", vca.ID).First(&vcaDB).Error
	if err != nil {
		return fmt.Errorf("不存在的账号，请核查")
	}

	//增加校验
	if global.ECContains(vca.Cid) {
		b := http2.IsValidCookie(token)
		if !b {
			return errors.New("传入的ck不合法，请核查")
		}
	} else {
		return errors.New("该信道暂不支持创建查卡池账号")
	}
	err = global.GVA_DB.Save(&vca).Error
	return err
}

// GetChannelCardAcc 根据id获取通道账号记录
func (cardAccService *ChannelCardAccService) GetChannelCardAcc(id uint) (vca vbox.ChannelCardAcc, err error) {
	err = global.GVA_DB.Unscoped().Where("id = ?", id).First(&vca).Error
	var sysUser sysModel.SysUser
	err = global.GVA_DB.Unscoped().Where("id = ?", vca.CreatedBy).First(&sysUser).Error
	vca.Username = sysUser.Username
	return
}

// GetChannelCardAccByAcId 根据AcId获取通道账号记录
func (cardAccService *ChannelCardAccService) GetChannelCardAccByAcId(acId string) (vca vbox.ChannelCardAcc, err error) {
	err = global.GVA_DB.Unscoped().Where("ac_id = ?", acId).First(&vca).Error
	var sysUser sysModel.SysUser
	err = global.GVA_DB.Unscoped().Where("id = ?", vca.CreatedBy).First(&sysUser).Error
	vca.Username = sysUser.Username
	return
}

// GetChannelCardAccInfoList 分页获取通道账号记录
func (cardAccService *ChannelCardAccService) GetChannelCardAccInfoList(info vboxReq.ChannelCardAccSearch, ids []uint) (list []vbox.ChannelCardAcc, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelCardAcc{})
	var vcas []vbox.ChannelCardAcc
	db.Where("created_by in (?)", ids)
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Username != "" {
		var sysUserIDs []uint
		err = global.GVA_DB.Unscoped().Table("sys_users").Select("id").Where("username LIKE ?", "%"+info.Username+"%").Scan(&sysUserIDs).Error
		if len(sysUserIDs) > 0 {
			db = db.Where("created_by in (?)", sysUserIDs)
		}
	}
	if info.AcRemark != "" {
		db = db.Where("ac_remark LIKE ?", "%"+info.AcRemark+"%")
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account LIKE ?", "%"+info.AcAccount+"%")
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.Status != nil {
		db = db.Where("status = ?", info.Status)
	}
	if info.SysStatus != nil {
		db = db.Where("sys_status = ?", info.SysStatus)
	}
	if info.AcId != "" {
		db = db.Where("ac_id = ?", info.AcId)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Order("id desc").Find(&vcas).Error
	return vcas, total, err
}
