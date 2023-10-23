package utils

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

func GetOwnerUserIdsList(id uint) (list []system.SysUser, total int64, err error) {

	db := global.GVA_DB.Model(&system.SysUser{})
	var userList []system.SysUser
	if id == 1 {
		err = db.Count(&total).Error
		if err != nil {
			return
		}
		err = db.Preload("Authorities").Preload("Authority").Find(&userList).Error

		return userList, total, err
	}
	if err != nil {
		return nil, 0, err
	}

	querySub := `
		WITH RECURSIVE cte (id) AS (SELECT id FROM sys_users WHERE id = ? UNION ALL SELECT sys_users.id FROM sys_users 
			JOIN cte ON sys_users.parent_id = cte.id)
		SELECT id,uuid,username,password,nick_name,side_mode,header_img,base_color,active_color,authority_id,phone,email,enable,parent_id
		FROM sys_users WHERE id IN (SELECT id FROM cte) and deleted_at is null;
    `
	//var userListSub []system.SysUser
	err = db.Raw(querySub, id).Preload("Authorities").Preload("Authority").Find(&userList).Error

	total = int64(len(userList))
	return userList, total, err
}
