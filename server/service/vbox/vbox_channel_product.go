package vbox

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ChannelProductService struct {
}

var ErrRoleExistence = errors.New("存在相同通道编码id")

// CreateVboxChannelProduct 创建VboxChannelProduct记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) CreateVboxChannelProduct(vcp *vbox.ChannelProduct) (err error) {
	var vcpDB vbox.ChannelProduct
	if !errors.Is(global.GVA_DB.Where("channel_code = ?", vcp.ChannelCode).First(&vcpDB).Error, gorm.ErrRecordNotFound) {
		return ErrRoleExistence
	}
	err = global.GVA_DB.Create(vcp).Error
	return err
}

// DeleteVboxChannelProduct 删除VboxChannelProduct记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) DeleteVboxChannelProduct(vcp vbox.ChannelProduct) (err error) {
	if !errors.Is(global.GVA_DB.Where("parent_id = ?", vcp.ChannelCode).First(&vbox.ChannelProduct{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("此产品存在子产品不允许删除")
	}
	global.GVA_LOG.Info("DeleteVboxChannelProduct", zap.Any("channel_code", vcp.ChannelCode))
	err = global.GVA_DB.Delete(&vcp).Error
	return err
}

// UpdateVboxChannelProduct 更新VboxChannelProduct记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) UpdateVboxChannelProduct(vcp vbox.ChannelProduct) (err error) {
	err = global.GVA_DB.Where("channel_code = ?", vcp.ChannelCode).First(&vbox.ChannelProduct{}).Updates(&vcp).Error
	return err
}

// GetVboxChannelProduct 根据id获取VboxChannelProduct记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) GetVboxChannelProduct(id uint) (vcp vbox.ChannelProduct, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vcp).Error
	return
}

// GetVboxChannelProductInfoList 分页获取VboxChannelProduct记录
// Author [piexlmax](https://github.com/piexlmax)
func (vcpService *ChannelProductService) GetVboxChannelProductInfoList(info vboxReq.VboxChannelProductSearch) (list []vbox.ChannelProduct, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelProduct{})
	if err = db.Where("parent_id = ?", "0").Count(&total).Error; total == 0 || err != nil {
		return
	}
	var vcps []vbox.ChannelProduct
	err = db.Limit(limit).Offset(offset).Where("parent_id = ?", "0").Find(&vcps).Error
	for k := range vcps {
		err = vcpService.findChildrenChannelProduct(&vcps[k])
	}
	return vcps, total, err
}

//@author: [piexlmax](https://github.com/piexlmax)
//@function: findChildrenAuthority
//@description: 查询子角色
//@param: authority *model.SysAuthority
//@return: err error

func (vcpService *ChannelProductService) findChildrenChannelProduct(vcp *vbox.ChannelProduct) (err error) {
	err = global.GVA_DB.Where("parent_id = ?", vcp.ChannelCode).Find(&vcp.Children).Error
	if len(vcp.Children) > 0 {
		for k := range vcp.Children {
			err = vcpService.findChildrenChannelProduct(&vcp.Children[k])
		}
	}
	return err
}
