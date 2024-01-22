package utils

import (
	"context"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/model"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"time"
)

// string 去重方法
func UniqStr(array []string) []string {
	var uintMap = make(map[string]bool)
	var res []string
	for _, u := range array {
		if !uintMap[u] {
			uintMap[u] = true
			res = append(res, u)
		}
	}
	return res
}

// uint 去重方法
func Uniq(array []uint) []uint {
	var uintMap = make(map[uint]bool)
	var res []uint
	for _, u := range array {
		if !uintMap[u] {
			uintMap[u] = true
			res = append(res, u)
		}
	}
	return res
}

const (
	Node    = 0 // 无资源权限
	Self    = 1 // 仅自己
	Current = 2 // 当前部门
	Deep    = 3 // 当前部门及以下
	All     = 4 // 所有
)

// 获取当前部门所有产品
func GetChannelCodeByOrgID(orgID uint) []string {
	var Products []vbox.OrgProduct
	var ids []string
	err := global.GVA_DB.Model(&vbox.OrgProduct{}).Joins("ChannelProduct").Preload("ChannelProduct").
		Find(&Products, "organization_id = ?", orgID).Error
	if err != nil {
		return []string{}
	}
	for i := range Products {
		ids = append(ids, Products[i].ChannelProduct.ChannelCode)
	}
	return UniqStr(ids)
}

// GetSelfOrg 获取当前部门ID(传user ID)
func GetSelfOrg(id uint) []uint {
	orgKey := fmt.Sprintf(global.SysUserOrgPrefix, id)
	orgVal := global.GVA_REDIS.SMembers(context.Background(), orgKey).Val()
	if len(orgVal) == 0 { //查库
		var orgUser []model.OrgUser
		err := global.GVA_DB.Find(&orgUser, "sys_user_id = ?", id).Error
		if err != nil {
			return []uint{}
		}
		var orgId []uint
		for _, m := range orgUser {
			orgId = append(orgId, m.OrganizationID)
		}

		uniq := Uniq(orgId)

		global.GVA_REDIS.SAdd(context.Background(), orgKey, uniq)
		global.GVA_REDIS.Expire(context.Background(), orgKey, 5*time.Minute)

		return uniq
	} else {
		orgIds, err := utils.ConvertStringSliceToUintSlice(orgVal)
		if err != nil {
			return []uint{}
		}
		return orgIds
	}
}

// 获取所有部门
func GetAllOrg() []model.Organization {
	var orgUser []model.Organization
	err := global.GVA_DB.Find(&orgUser).Error
	if err != nil {
		return []model.Organization{}
	}
	return orgUser
}

// 获取所有部门ID
func GetAllOrgID() []uint {
	orgUser := GetAllOrg()
	if len(orgUser) == 0 {
		return []uint{}
	}
	var orgids []uint
	for _, organization := range orgUser {
		orgids = append(orgids, organization.ID)
	}
	return Uniq(orgids)
}

// 获取当前部门及以下部门id (传参userId)
func GetDeepOrg(userId uint) []uint {
	orgId := GetSelfOrg(userId)
	if len(orgId) == 0 {
		return []uint{}
	}
	orgs := GetAllOrg()
	if len(orgs) == 0 {
		return []uint{}
	}
	orgids := findChildren(orgId, orgs)
	return Uniq(append(orgids, orgId...))
}

// 获取当前部门及以下部门的递归方法
func findChildren(ids []uint, orgs []model.Organization) []uint {
	var idsMap = make(map[uint]bool)
	var resIDs []uint
	for _, id := range ids {
		idsMap[id] = true
	}
	for _, org := range orgs {
		if idsMap[org.ParentID] {
			resIDs = append(resIDs, org.ID)
		}
	}
	if len(resIDs) == 0 {
		resIDs = append(resIDs, ids...)
		return resIDs
	}
	dids := findChildren(resIDs, orgs)
	resIDs = append(resIDs, ids...)
	resIDs = append(resIDs, dids...)
	return resIDs
}

// 获取当前部门的用户id
func GetCurrentUserIDs(id uint) []uint {
	orgId := GetSelfOrg(id)
	if len(orgId) == 0 {
		return []uint{}
	}
	return GetUsersByOrgIds(orgId)
}

// 获取当前部门及以下的用户id
func GetDeepUserIDs(id uint) []uint {
	orgids := GetDeepOrg(id)
	if len(orgids) == 0 {
		return []uint{}
	}
	return GetUsersByOrgIds(orgids)
}

// GetUsersByOrgId 根据部门获取部门下用户ID(单个)
func GetUsersByOrgId(orgId uint) []uint {
	var orgUser []model.OrgUser
	err := global.GVA_DB.Find(&orgUser, "organization_id = ?", orgId).Error
	if err != nil {
		return []uint{}
	}
	var userIDS []uint
	for _, m := range orgUser {
		userIDS = append(userIDS, m.SysUserID)
	}
	return Uniq(userIDS)
}

// 根据部门获取部门下用户ID
func GetUsersByOrgIds(orgIds []uint) []uint {
	var orgUser []model.OrgUser
	err := global.GVA_DB.Find(&orgUser, "organization_id in (?)", orgIds).Error
	if err != nil {
		return []uint{}
	}
	var userIDS []uint
	for _, m := range orgUser {
		userIDS = append(userIDS, m.SysUserID)
	}
	return Uniq(userIDS)
}

// 获取所有用户ID
func GetAllUserIDs() []uint {
	var users []system.SysUser
	err := global.GVA_DB.Find(&users).Error
	if err != nil {
		return []uint{}
	}
	var usersID []uint
	for _, sysUser := range users {
		usersID = append(usersID, sysUser.ID)
	}
	return Uniq(usersID)
}

// GetUserIDS 自动获取当前用户拥有的权限的用户ID
func GetUserIDS(c *gin.Context) []uint {
	user := utils.GetUserInfo(c)
	var data model.DataAuthority
	err := global.GVA_DB.Debug().First(&data, "authority_id = ?", user.AuthorityId).Error
	if err != nil {
		return []uint{}
	}
	switch data.AuthorityType {
	case Node:
		return []uint{}
	case Self:
		return []uint{user.BaseClaims.ID}
	case Current:
		return GetCurrentUserIDs(user.BaseClaims.ID)
	case Deep:
		return GetDeepUserIDs(user.BaseClaims.ID)
	case All:
		return GetAllUserIDs()
	}
	return []uint{}
}

// 自动获取当前用户拥有的权限的部门ID
func GetOrgIDS(c *gin.Context) []uint {
	user := utils.GetUserInfo(c)
	var data model.DataAuthority
	err := global.GVA_DB.First(&data, "authority_id = ?", user.AuthorityId).Error
	if err != nil {
		global.GVA_LOG.Error(fmt.Sprintf("err: %s", err))
		return []uint{}
	}
	switch data.AuthorityType {
	case Node:
		return []uint{}
	case Self:
		return GetSelfOrg(user.BaseClaims.ID)
	case Current:
		return GetSelfOrg(user.BaseClaims.ID)
	case Deep:
		return GetDeepOrg(user.BaseClaims.ID)
	case All:
		return GetAllOrgID()
	}
	return []uint{}
}
