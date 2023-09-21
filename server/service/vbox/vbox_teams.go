package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	"gorm.io/gorm"
)

type VboxTeamsService struct {
}

// CreateVboxTeams 创建VboxTeams记录
// Author [piexlmax](https://github.com/piexlmax)
func (teamsService *VboxTeamsService) CreateVboxTeams(teams *vbox.VboxTeams) (err error) {
	err = global.GVA_DB.Create(teams).Error
	return err
}

// DeleteVboxTeams 删除VboxTeams记录
// Author [piexlmax](https://github.com/piexlmax)
func (teamsService *VboxTeamsService) DeleteVboxTeams(teams vbox.VboxTeams) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxTeams{}).Where("id = ?", teams.ID).Update("deleted_by", teams.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&teams).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteVboxTeamsByIds 批量删除VboxTeams记录
// Author [piexlmax](https://github.com/piexlmax)
func (teamsService *VboxTeamsService) DeleteVboxTeamsByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxTeams{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxTeams{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateVboxTeams 更新VboxTeams记录
// Author [piexlmax](https://github.com/piexlmax)
func (teamsService *VboxTeamsService) UpdateVboxTeams(teams vbox.VboxTeams) (err error) {
	err = global.GVA_DB.Save(&teams).Error
	return err
}

// GetVboxTeams 根据id获取VboxTeams记录
// Author [piexlmax](https://github.com/piexlmax)
func (teamsService *VboxTeamsService) GetVboxTeams(id uint) (teams vbox.VboxTeams, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&teams).Error
	return
}

// GetVboxTeamsInfoList 分页获取VboxTeams记录
// Author [piexlmax](https://github.com/piexlmax)
func (teamsService *VboxTeamsService) GetVboxTeamsInfoList(info vboxReq.VboxTeamsSearch) (list []vbox.VboxTeams, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&vbox.VboxTeams{})
	var orgs []vbox.VboxTeams
	// 如果有条件搜索 下方会自动创建搜索语句
	db = db.Where("parent_id = ?", info.Parent_id)
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&orgs).Error
	return orgs, total, err
	//limit := info.PageSize
	//offset := info.PageSize * (info.Page - 1)
	//// 创建db
	//db := global.GVA_DB.Model(&vbox.VboxTeams{})
	//var teamss []vbox.VboxTeams
	//// 如果有条件搜索 下方会自动创建搜索语句
	//if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
	//	db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	//}
	//if info.Name != "" {
	//	db = db.Where("name LIKE ?", "%"+info.Name+"%")
	//}
	//if info.Parent_id != 99999 {
	//	db = db.Where("parent_id = ?", info.Parent_id)
	//}
	//err = db.Count(&total).Error
	//if err != nil {
	//	return
	//}
	//
	//err = db.Limit(limit).Offset(offset).Find(&teamss).Error
	//return teamss, total, err
}
