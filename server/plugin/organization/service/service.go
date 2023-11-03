package service

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	organization "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/model"
	organizationReq "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/model/request"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/utils"
	"gorm.io/gorm"
)

type OrganizationService struct {
}

// CreateOrganization 创建Organization记录
// Author [piexlmax](https://github.com/piexlmax)
func (orgService *OrganizationService) CreateOrganization(org organization.Organization) (err error) {
	err = global.GVA_DB.Create(&org).Error
	return err
}

// DeleteOrganization 删除Organization记录
// Author [piexlmax](https://github.com/piexlmax)
func (orgService *OrganizationService) DeleteOrganization(org organization.Organization) (err error) {
	err = global.GVA_DB.Where("parent_id = ?", org.ID).First(&organization.Organization{}).Error
	if err == nil {
		return errors.New("该组织有子组织，不能删除")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = global.GVA_DB.Delete(&org, "id = ?", org.ID).Error
	}
	return err
}

// DeleteOrganizationByIds 批量删除Organization记录
// Author [piexlmax](https://github.com/piexlmax)
func (orgService *OrganizationService) DeleteOrganizationByIds(ids request.IdsReq) (err error) {
	err = global.GVA_DB.Delete(&[]organization.Organization{}, "id in ?", ids.Ids).Error
	return err
}

// UpdateOrganization 更新Organization记录
// Author [piexlmax](https://github.com/piexlmax)
func (orgService *OrganizationService) UpdateOrganization(org organization.Organization) (err error) {
	var updatesmap = make(map[string]interface{})
	updatesmap["Name"] = org.Name
	updatesmap["ParentID"] = org.ParentID
	err = global.GVA_DB.Model(&organization.Organization{}).Where("id = ?", org.ID).Updates(updatesmap).Error
	return err
}

// GetOrganization 根据id获取Organization记录
// Author [piexlmax](https://github.com/piexlmax)
func (orgService *OrganizationService) GetOrganization(id uint) (org organization.Organization, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&org).Error
	return
}

// GetOrganizationInfoList 分页获取Organization记录
// Author [piexlmax](https://github.com/piexlmax)
func (orgService *OrganizationService) GetOrganizationInfoList(info organizationReq.OrganizationSearch) (list interface{}, total int64, err error) {
	// 创建db
	db := global.GVA_DB.Model(&organization.Organization{})
	var orgs []organization.Organization
	// 如果有条件搜索 下方会自动创建搜索语句
	db = db.Where("parent_id = ?", info.ParentID)
	if info.Name != "" {
		db = db.Where("name LIKE ?", "%"+info.Name+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Find(&orgs).Error
	return orgs, total, err
}

func (orgService *OrganizationService) CreateOrgUser(orgUser organization.OrgUserReq) error {
	var Users []organization.OrgUser
	var CUsers []organization.OrgUser
	err := global.GVA_DB.Find(&Users, "organization_id = ?", orgUser.OrganizationID).Error
	if err != nil {
		return err
	}
	var UserIdMap = make(map[uint]bool)
	for i := range Users {
		UserIdMap[Users[i].SysUserID] = true
	}

	for i := range orgUser.SysUserIDS {
		if !UserIdMap[orgUser.SysUserIDS[i]] {
			CUsers = append(CUsers, organization.OrgUser{SysUserID: orgUser.SysUserIDS[i], OrganizationID: orgUser.OrganizationID})
		}
	}
	err = global.GVA_DB.Create(&CUsers).Error
	return err
}

func (orgService *OrganizationService) FindOrgUserAll(orgID string) ([]uint, error) {
	var Users []organization.OrgUser
	var ids []uint
	err := global.GVA_DB.Find(&Users, "organization_id = ?", orgID).Error
	if err != nil {
		return ids, err
	}
	for i := range Users {
		ids = append(ids, Users[i].SysUserID)
	}
	return ids, err
}

// GetOrganizationInfoList 分页获取当前组织下用户记录
// Author [piexlmax](https://github.com/piexlmax)
func (orgService *OrganizationService) GetOrgUserList(info organizationReq.OrgUserSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&organization.OrgUser{}).Joins("SysUser").Preload("SysUser.Authority")
	var orgs []organization.OrgUser
	// 如果有条件搜索 下方会自动创建搜索语句
	db = db.Where("organization_id = ?", info.OrganizationID)
	if info.UserName != "" {
		db = db.Where("SysUser.nick_name LIKE ?", "%"+info.UserName+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&orgs).Error
	return orgs, total, err
}

// GetOrgUserListSelf 分页获取当前用户所属组织下用户记录
func (orgService *OrganizationService) GetOrgUserListSelf(info organizationReq.OrgUserSearch) (list interface{}, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&organization.OrgUser{}).Joins("SysUser").Preload("SysUser.Authority")
	var orgs []organization.OrgUser
	// 如果有条件搜索 下方会自动创建搜索语句
	db = db.Where("organization_id in ?", info.OrgIds)
	if info.UserName != "" {
		db = db.Where("SysUser.nick_name LIKE ?", "%"+info.UserName+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&orgs).Error
	return orgs, total, err
}

func (orgService *OrganizationService) SetOrgUserAdmin(id uint, flag bool) (err error) {
	return global.GVA_DB.Model(&organization.OrgUser{}).Where("sys_user_id = ?", id).Update("is_admin", flag).Error
}

func (orgService *OrganizationService) SetOrgAuthority(authID uint, authorityType int) (err error) {
	return global.GVA_DB.Model(&organization.DataAuthority{}).Where("authority_id = ?", authID).Update("authority_type", authorityType).Error
}

func (orgService *OrganizationService) GetOrgAuthority() (authorityData []organization.DataAuthority, err error) {
	err = global.GVA_DB.Preload("Authority").Find(&authorityData).Error
	return authorityData, err
}

func (orgService *OrganizationService) SyncAuthority() (err error) {
	var authData []system.SysAuthority
	var auth []organization.DataAuthority
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var idMap = make(map[uint]*bool)
		err := tx.Find(&authData).Error
		if err != nil {
			return err
		}
		for _, datum := range authData {
			idMap[datum.AuthorityId] = utils.GetBoolPointer(true)
		}
		err = tx.Find(&auth).Error
		if err != nil {
			return err
		}
		for _, datum := range auth {
			if idMap[datum.AuthorityID] != nil {
				idMap[datum.AuthorityID] = utils.GetBoolPointer(false)
			} else {
				idMap[datum.AuthorityID] = nil
			}

		}
		var ayncData []organization.DataAuthority
		var deleteAuth []organization.DataAuthority

		for k, _ := range idMap {
			if idMap[k] != nil && *idMap[k] {
				ayncData = append(ayncData, organization.DataAuthority{
					AuthorityID:   k,
					AuthorityType: 0,
				})
			}
			if idMap[k] == nil {
				deleteAuth = append(deleteAuth, organization.DataAuthority{
					AuthorityID:   k,
					AuthorityType: 0,
				})
			}
		}
		if len(ayncData) > 0 || len(deleteAuth) > 0 {
			if len(ayncData) > 0 {
				err := tx.Create(&ayncData).Error

				if err != nil {
					return err
				}
			}

			if len(deleteAuth) > 0 {
				var deleteAuthIds []uint
				for i := range deleteAuth {
					deleteAuthIds = append(deleteAuthIds, deleteAuth[i].AuthorityID)
				}
				err = tx.Delete(&deleteAuth, "authority_id in (?)", deleteAuthIds).Error
				if err != nil {
					return err
				}
			}
			return nil
		} else {
			return errors.New("当前无需同步")
		}
	})
}

func (orgService *OrganizationService) DeleteOrgUser(ids []uint, orgID uint) (err error) {
	return global.GVA_DB.Where("sys_user_id in (?) and organization_id = ?", ids, orgID).Delete(&[]organization.OrgUser{}).Error
}

func (orgService *OrganizationService) TransferOrgUser(ids []uint, orgID, toOrgID uint) (err error) {
	return global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		var CUsers []organization.OrgUser
		err := global.GVA_DB.Where("sys_user_id in (?) and organization_id in (?)", ids, []uint{orgID, toOrgID}).Delete(&[]organization.OrgUser{}).Error
		if err != nil {
			return err
		}
		for i := range ids {
			CUsers = append(CUsers, organization.OrgUser{SysUserID: ids[i], OrganizationID: toOrgID})
		}
		err = global.GVA_DB.Create(&CUsers).Error
		if err != nil {
			return err
		}
		return nil
	})
}
