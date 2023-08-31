package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
)

type ProxyService struct {
}

// CreateVboxProxy 创建VboxProxy记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxProxyService *ProxyService) CreateVboxProxy(vboxProxy *vbox.Proxy) (err error) {
	err = global.GVA_DB.Create(vboxProxy).Error
	return err
}

// DeleteVboxProxy 删除VboxProxy记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxProxyService *ProxyService) DeleteVboxProxy(vboxProxy vbox.Proxy) (err error) {
	err = global.GVA_DB.Delete(&vboxProxy).Error
	return err
}

// DeleteVboxProxyByIds 批量删除VboxProxy记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxProxyService *ProxyService) DeleteVboxProxyByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]vbox.Proxy{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateVboxProxy 更新VboxProxy记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxProxyService *ProxyService) UpdateVboxProxy(vboxProxy vbox.Proxy) (err error) {
	err = global.GVA_DB.Save(&vboxProxy).Error
	return err
}

// GetVboxProxy 根据id获取VboxProxy记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxProxyService *ProxyService) GetVboxProxy(id uint) (vboxProxy vbox.Proxy, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vboxProxy).Error
	return
}

// GetVboxProxyInfoList 分页获取VboxProxy记录
// Author [piexlmax](https://github.com/piexlmax)
func (vboxProxyService *ProxyService) GetVboxProxyInfoList(info vboxReq.VboxProxySearch) (list []vbox.Proxy, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.Proxy{})
	var vboxProxys []vbox.Proxy
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	err = db.Limit(limit).Offset(offset).Find(&vboxProxys).Error
	return vboxProxys, total, err
}
