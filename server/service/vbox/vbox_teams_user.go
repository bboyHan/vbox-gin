package vbox

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	organization "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/model"
	"gorm.io/gorm"
)

type VboxTeamsUserService struct {
}

// CreateVboxTeamsUser 创建VboxTeamsUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (tUsersService *VboxTeamsUserService) CreateVboxTeamsUser(tUsers vbox.TeamsUserReq) (err error) {
	var Users []vbox.VboxTeamsUser
	var CUsers []vbox.VboxTeamsUser
	err = global.GVA_DB.Find(&Users, "team_id = ?", tUsers.TeamID).Error
	if err != nil {
		return err
	}
	var UserIdMap = make(map[uint]bool)
	for i := range Users {
		UserIdMap[Users[i].Uid] = true
	}

	for i := range tUsers.SysUserIDS {
		if !UserIdMap[tUsers.SysUserIDS[i]] {
			CUsers = append(CUsers, vbox.VboxTeamsUser{Uid: tUsers.SysUserIDS[i], Team_id: tUsers.TeamID})
		}
	}
	err = global.GVA_DB.Create(&CUsers).Error
	//err = global.GVA_DB.Create(tUsers).Error
	return err
}

// DeleteVboxTeamsUser 删除VboxTeamsUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (tUsersService *VboxTeamsUserService) DeleteVboxTeamsUser(ids []uint, orgID uint) (err error) {
	return global.GVA_DB.Where("uid in (?) and team_id = ?", ids, orgID).Delete(&[]vbox.VboxTeamsUser{}).Error
}

// DeleteVboxTeamsUserByIds 批量删除VboxTeamsUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (tUsersService *VboxTeamsUserService) DeleteVboxTeamsUserByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.VboxTeamsUser{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.VboxTeamsUser{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateVboxTeamsUser 更新VboxTeamsUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (tUsersService *VboxTeamsUserService) UpdateVboxTeamsUser(tUsers vbox.VboxTeamsUser) (err error) {
	err = global.GVA_DB.Save(&tUsers).Error
	return err
}

// GetVboxTeamsUser 根据id获取VboxTeamsUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (tUsersService *VboxTeamsUserService) GetVboxTeamsUser(id uint) (tUsers vbox.VboxTeamsUser, err error) {
	err = global.GVA_DB.Where("team_id = ?", id).First(&tUsers).Error
	return tUsers, err
}

// GetVboxTeamsUserInfoList 分页获取VboxTeamsUser记录
// Author [piexlmax](https://github.com/piexlmax)
func (tUsersService *VboxTeamsUserService) GetVboxTeamsUserInfoList(info vboxReq.VboxTeamsUserSearch) (list []vbox.VboxTeamsUserRep, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var orgs []vbox.VboxTeamsUserRep
	// 如果有条件搜索 下方会自动创建搜索语句
	db := global.GVA_DB.Model(&vbox.VboxTeamsUser{}).Where("team_id = ?", info.Team_id)
	// 创建db
	db = db.Select("vbox_teams_user.*, t1.username,t2.authority_name").
		Joins("left join sys_users t1 on t1.id = vbox_teams_user.uid").
		Joins("left join sys_authorities t2 on t1.authority_id = t2.authority_id")

	//if info.Team_name != "" {
	//	db = db.Where("SysUser.nick_name LIKE ?", "%"+info.UserName+"%")
	//}

	err = db.Limit(limit).Offset(offset).Find(&orgs).Error
	//err = db.Count(&total).Error
	//if err != nil {
	//	return
	//}
	total = int64(len(orgs))
	return orgs, total, err
	//limit := info.PageSize
	//offset := info.PageSize * (info.Page - 1)
	//// 创建db
	//db := global.GVA_DB.Model(&vbox.VboxTeamsUser{})
	//var tUserss []vbox.VboxTeamsUser
	//// 如果有条件搜索 下方会自动创建搜索语句
	//if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
	//	db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	//}
	//if info.Team_id != 99999 {
	//	db = db.Where("team_id = ?", info.Team_id)
	//}
	//if info.Team_name != "" {
	//	db = db.Where("team_name LIKE ?", "%"+info.Team_name+"%")
	//}
	//if info.Uid != 99999 {
	//	db = db.Where("uid = ?", info.Uid)
	//}
	//if info.Leader_id != "" {
	//	db = db.Where("leader_id = ?", info.Leader_id)
	//}
	//err = db.Count(&total).Error
	//if err != nil {
	//	return
	//}
	//
	//err = db.Limit(limit).Offset(offset).Find(&tUserss).Error
	//return tUserss, total, err
}

func (tUsersService *VboxTeamsUserService) FindTeamUserAll(teamID string) ([]uint, error) {
	var Users []vbox.VboxTeamsUser
	var ids []uint
	err := global.GVA_DB.Find(&Users, "team_id = ?", teamID).Error
	if err != nil {
		return ids, err
	}
	for i := range Users {
		ids = append(ids, Users[i].Uid)
	}
	return ids, err
}

//
//func (tUsersService *VboxTeamsUserService) GetOrgAuthority() (authorityData []organization.DataAuthority, err error) {
//	err = global.GVA_DB.Preload("Authority").Find(&authorityData).Error
//	return authorityData, err
//}
//
//func (tUsersService *VboxTeamsUserService) SyncAuthority() (err error) {
//	var authData []system.SysAuthority
//	var auth []organization.DataAuthority
//	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
//		var idMap = make(map[uint]*bool)
//		err := tx.Find(&authData).Error
//		if err != nil {
//			return err
//		}
//		for _, datum := range authData {
//			idMap[datum.AuthorityId] = utils.GetBoolPointer(true)
//		}
//		err = tx.Find(&auth).Error
//		if err != nil {
//			return err
//		}
//		for _, datum := range auth {
//			if idMap[datum.AuthorityID] != nil {
//				idMap[datum.AuthorityID] = utils.GetBoolPointer(false)
//			} else {
//				idMap[datum.AuthorityID] = nil
//			}
//
//		}
//		var ayncData []organization.DataAuthority
//		var deleteAuth []organization.DataAuthority
//
//		for k, _ := range idMap {
//			if idMap[k] != nil && *idMap[k] {
//				ayncData = append(ayncData, organization.DataAuthority{
//					AuthorityID:   k,
//					AuthorityType: 0,
//				})
//			}
//			if idMap[k] == nil {
//				deleteAuth = append(deleteAuth, organization.DataAuthority{
//					AuthorityID:   k,
//					AuthorityType: 0,
//				})
//			}
//		}
//		if len(ayncData) > 0 || len(deleteAuth) > 0 {
//			if len(ayncData) > 0 {
//				err := tx.Create(&ayncData).Error
//
//				if err != nil {
//					return err
//				}
//			}
//
//			if len(deleteAuth) > 0 {
//				var deleteAuthIds []uint
//				for i := range deleteAuth {
//					deleteAuthIds = append(deleteAuthIds, deleteAuth[i].AuthorityID)
//				}
//				err = tx.Delete(&deleteAuth, "authority_id in (?)", deleteAuthIds).Error
//				if err != nil {
//					return err
//				}
//			}
//			return nil
//		} else {
//			return errors.New("当前无需同步")
//		}
//	})
//}

func (tUsersService *VboxTeamsUserService) TransferTeamUser(ids []uint, orgID, toOrgID uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var CUsers []vbox.VboxTeamsUser
		err := global.GVA_DB.Where("uid in (?) and team_id in (?)", ids, []uint{orgID, toOrgID}).Delete(&[]organization.OrgUser{}).Error
		if err != nil {
			return err
		}
		for i := range ids {
			CUsers = append(CUsers, vbox.VboxTeamsUser{Uid: ids[i], Team_id: toOrgID})
		}
		err = global.GVA_DB.Create(&CUsers).Error
		if err != nil {
			return err
		}
		return nil
	})
}
