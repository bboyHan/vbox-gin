package vbox

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxResp "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	organization "github.com/flipped-aurora/gin-vue-admin/server/plugin/organization/model"
	"gorm.io/gorm"
	"log"
	"time"
)

type BdaChorgIndexDService struct {
}

// CreateBdaChorgIndexD 创建通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) CreateBdaChorgIndexD(bdaChorg *vbox.BdaChorgIndexD) (err error) {
	err = global.GVA_DB.Create(bdaChorg).Error
	return err
}

// DeleteBdaChorgIndexD 删除通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) DeleteBdaChorgIndexD(bdaChorg vbox.BdaChorgIndexD) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChorgIndexD{}).Where("id = ?", bdaChorg.ID).Update("deleted_by", bdaChorg.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&bdaChorg).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteBdaChorgIndexDByIds 批量删除通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) DeleteBdaChorgIndexDByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.BdaChorgIndexD{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.BdaChorgIndexD{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateBdaChorgIndexD 更新通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) UpdateBdaChorgIndexD(bdaChorg vbox.BdaChorgIndexD) (err error) {
	err = global.GVA_DB.Save(&bdaChorg).Error
	return err
}

// GetBdaChorgIndexD 根据id获取通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) GetBdaChorgIndexD(id uint) (bdaChorg vbox.BdaChorgIndexD, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&bdaChorg).Error
	return
}

// GetBdaChorgIndexDInfoList 分页获取通道团队统计-天更新记录
// Author [piexlmax](https://github.com/piexlmax)
func (bdaChorgService *BdaChorgIndexDService) GetBdaChorgIndexDInfoList(info vboxReq.BdaChorgIndexDSearch) (list []vbox.BdaChorgIndexD, total int64, err error) {
	fmt.Println("统计开始")
	//bdaChorgService.CronVboxBdaChOrgIndexD()
	fmt.Println("统计结束")
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.BdaChorgIndexD{})
	var bdaChorgs []vbox.BdaChorgIndexD
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.OrganizationId != 0 {
		db = db.Where("organization_id = ?", info.OrganizationId)
	}
	if info.OrganizationName != "" {
		db = db.Where("organization_name LIKE ?", "%"+info.OrganizationName+"%")
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.ChannelCode != "" {
		db = db.Where("channel_code = ?", info.ChannelCode)
	}
	if info.ProductId != "" {
		db = db.Where("product_id = ?", info.ProductId)
	}
	if info.ProductName != "" {
		db = db.Where("product_name LIKE ?", "%"+info.ProductName+"%")
	}
	if info.Dt != "" {
		db = db.Where("dt = ?", info.Dt)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Find(&bdaChorgs).Error
	return bdaChorgs, total, err
}

func (bdaChOrgIndexDService *BdaChorgIndexDService) CronVboxBdaChOrgIndexD() (err error) {

	querySql := `

			SELECT
				a.organization_id as organizationId,
				a.organization_name as organizationName,
				a.channel_code as channelCode,
				step_time as stepTime,
				c.product_id as productId,
				c.product_name as productName,
				a.cnt as orderQuantify,
				a.ok_cnt as okOrderQuantify,
				ROUND((100.0 * a.ok_cnt/a.cnt),2) as ratio,
				a.ok_income as okIncome
			FROM
			(
				SELECT
				 organization_id,
				 organization_name,
				 channel_code,
				 step_time,
				 sum(if (order_status =1,money,0)) as ok_income,
				 count(*) as cnt,
				 sum(if(order_status=1,1,0)) as ok_cnt
				FROM(
					select 
					    uid,
					    channel_code,
					    order_status,
					    money,
					    step_time,
					    organization_id,
					    organization_name
					from
					(
						SELECT
								created_by as uid,
								channel_code,
								order_status,
								money,
								DATE_FORMAT(created_at, '%Y-%m-%d') AS step_time	
						from vbox_pay_order 
						where DATE_FORMAT(created_at, '%Y-%m-%d') = ?
						and event_type = 1  
					) a1	
					join(
						 SELECT 
								organization_id,
								name as organization_name,
								sys_user_id
							from
								(
								select
									organization_id,
									sys_user_id
								from
									org_user
							) o1
							left join 
							(
								select
									id,
									name
								from
									organization
							) o2
							on
								o1.organization_id = o2.id
					                                    
					) a2
				    on a1.uid = a2.sys_user_id
					
				)aa
				GROUP BY organization_id,channel_code,organization_name
			)a 
			left join (
			SELECT DISTINCT channel_code,product_id,product_name 
			FROM vbox_channel_product where channel_code !=''
			) c
			on a.channel_code = c.channel_code  

`
	dt := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	db := global.GVA_DB.Model(&vbox.PayOrder{}).Where("DATE_FORMAT(created_at, '%Y-%m-%d') = ? ", dt)

	rows, err := db.Raw(querySql, dt).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	for rows.Next() {
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.ChOrgTableResp
		err := rows.Scan(&item.OrganizationId, &item.OrganizationName, &item.ChannelCode, &item.StepTime,
			&item.ProductId, &item.ProductName, &item.OrderQuantify, &item.OkOrderQuantify, &item.Ratio, &item.OkIncome)
		if err != nil {
			panic(err)
		}
		// 将 acId + '-' + stepTime 作为 key 存入 map，okIncome 作为 value
		entity := vbox.BdaChorgIndexD{
			OrganizationId:   item.OrganizationId,
			OrganizationName: item.OrganizationName,
			Cid:              item.ChannelCode,
			ChannelCode:      item.ChannelCode,
			ProductId:        item.ProductId,
			ProductName:      item.ProductName,
			OrderQuantify:    item.OrderQuantify,
			OkOrderQuantify:  item.OkOrderQuantify,
			Ratio:            item.Ratio,
			Income:           item.OkIncome,
			Dt:               dt,
			CreatedBy:        1,
			UpdatedBy:        1,
			DeletedBy:        1,
		}
		//fmt.Println(entity.Dt)
		log.Println("统计 ch org 结果=", entity)
		err = global.GVA_DB.Save(&entity).Error
		//fmt.Printf("---> UID: %d, AcID: %d, OkIncome: %f, StepTime: %s\n", item.Uid, item.AcID, item.OkIncome, item.StepTime)
	}

	return err
}

func (bdaChorgService *BdaChorgIndexDService) GetBdaChorgIndexRealList(info vboxReq.OrgSelectForm) (
	list []vboxResp.ChaOrgRealCardResp, total int64, err error) {
	var parentId int
	var orgs []organization.Organization
	var orgIds []int

	err = global.GVA_DB.Model(&orgs).Select("parent_id").Where("id = ?", info.OrganizationID).Pluck("parent_id", &parentId).Error
	if err != nil {
		panic(err)
	}
	if parentId != 0 {
		orgIds = append(orgIds, info.OrganizationID)
	} else {
		err = global.GVA_DB.Model(&orgs).Select("id").Where("parent_id = ?", info.OrganizationID).Pluck("id", &orgIds).Error
		if err != nil {
			panic(err)
		}
	}

	dt := time.Now().AddDate(0, 0, 0).Format("2006-01-02")
	uids, err := getYyUids(orgIds)

	fmt.Println("orgIds = ", orgIds)
	fmt.Println("uid = ", *info.SysUserID)
	fmt.Println("OrganizationID = ", info.OrganizationID)
	fmt.Println("uids = ", uids)
	if info.PAccount != "" {
		fmt.Println("1--->")
		resp, err := getPaccCardResp(dt, uids)
		return resp, int64(len(resp)), err
	}
	if info.SysUserID != nil && *info.SysUserID != 0 {
		resp, err := getUidCardResp(dt, uids)
		return resp, int64(len(resp)), err
	}
	if info.Cid != "" {
		resp, err := getCidCardResp(dt, uids)
		return resp, int64(len(resp)), err
	}
	if info.OrganizationID != 0 {
		resp, err := getOrgCardResp(dt, uids, info)
		return resp, int64(len(resp)), err
	}
	// 默认方法
	resp, err := getOrgCardResp(dt, uids, info)
	return resp, int64(len(resp)), err
}

func getYyUids(orgs []int) (sysUserIDs []int, err error) {
	// 成员维度统计
	queryUidsSql := `
		SELECT sys_user_id as sysUserID
		FROM (
			SELECT organization_id, sys_user_id
			FROM org_user
		) o1
		JOIN (
			SELECT id, name, parent_id
			FROM organization
			WHERE id IN ?
		) o2
		ON o1.organization_id = o2.id
	`

	//fmt.Println("orgs-->", orgs)
	db := global.GVA_DB.Model(&vboxResp.ChaOrgRealCardResp{})
	rows, err := db.Raw(queryUidsSql, orgs).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()

	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var sysUserID int
		err := rows.Scan(&sysUserID)
		if err != nil {
			panic(err)
		}

		sysUserIDs = append(sysUserIDs, sysUserID)
	}
	//fmt.Println("sysUserIDs = ", sysUserIDs)
	return sysUserIDs, err

	// 打印查询结果
	//fmt.Println(list)
}

func getPaccCardResp(dt string, orgs []int) (list []vboxResp.ChaOrgRealCardResp, err error) {
	// 成员维度统计
	queryPaccSql := `
		SELECT
			a.p_account as pAccount,
			'' as stepTime,
			coalesce(c.p_remark, a.p_account) as userName,
			a.cnt as orderQuantify,
			a.ok_cnt as okOrderQuantify,
			a.ok_income as okIncome
		FROM
		(
			SELECT
			 p_account,
			 sum(if (order_status=1 and cb_status=1,money,0)) as ok_income,
			 count(*) as cnt,
			 sum(if(order_status=1 and cb_status=1,1,0)) as ok_cnt
			FROM(
				SELECT
							created_by as uid,
							p_account,
							channel_code,
							order_status,
							cb_status,
							money	
					from vbox_pay_order 
					where (DATE_FORMAT(cb_time, '%Y-%m-%d') = ? or DATE_FORMAT(created_at, '%Y-%m-%d') = ?)
					and event_type = 1  and created_by in ?
				
			)aa
			GROUP BY p_account
		)a 
		left join (
		SELECT DISTINCT p_account,p_remark
		FROM vbox_pay_account   
		) c
		on a.p_account = c.p_account  
	`

	//fmt.Println("dt-->", dt, "uid-->", uid, "querySql-->", querySql)
	db := global.GVA_DB.Model(&vboxResp.ChaOrgRealCardResp{})
	rows, err := db.Raw(queryPaccSql, dt, dt, orgs).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.PaccRealStatisicsResp
		err := rows.Scan(&item.PAccount, &item.StepTime, &item.UserName, &item.OrderQuantify, &item.OkOrderQuantify, &item.OkIncome)
		if err != nil {
			panic(err)
		}
		defaultItem := vboxResp.ChaOrgRealCardResp{
			Title:           item.UserName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           0,
			Income:          item.OkIncome,
			Dt:              dt,
		}
		list = append(list, defaultItem)

	}

	return list, err

	// 打印查询结果
	//fmt.Println(list)
}

func getUidCardResp(dt string, orgs []int) (list []vboxResp.ChaOrgRealCardResp, err error) {
	// 成员维度统计
	queryUidSql := `

			SELECT
				a.uid as uid,
				'' as stepTime,
				c.nickname as userName,
				a.cnt as orderQuantify,
				a.ok_cnt as okOrderQuantify,
				a.ok_income as okIncome
			FROM
			(
				SELECT
				 uid,
				 sum(if (order_status=1 and cb_status=1,money,0)) as ok_income,
				 count(*) as cnt,
				 sum(if(order_status=1 and cb_status=1,1,0)) as ok_cnt
				FROM(
					SELECT
								created_by as uid,
								p_account,
								channel_code,
								order_status,
								cb_status,
								money
						from vbox_pay_order 
						where (DATE_FORMAT(cb_time, '%Y-%m-%d') = ? or DATE_FORMAT(created_at, '%Y-%m-%d') = ?)
						and event_type = 1   and created_by in ?
					
				)aa
				GROUP BY uid
			)a 
			left join (
			SELECT DISTINCT id,nickname
			FROM sys_users 
			) c
			on a.uid = c.id  
`

	//fmt.Println("dt-->", dt, "uid-->", uid, "querySql-->", querySql)
	db := global.GVA_DB.Model(&vboxResp.ChaOrgRealCardResp{})
	rows, err := db.Raw(queryUidSql, dt, dt, orgs).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.UidRealStatisicsResp
		err := rows.Scan(&item.Uid, &item.StepTime, &item.UserName, &item.OrderQuantify, &item.OkOrderQuantify, &item.OkIncome)
		if err != nil {
			panic(err)
		}
		defaultItem := vboxResp.ChaOrgRealCardResp{
			Title:           item.UserName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           0,
			Income:          item.OkIncome,
			Dt:              dt,
		}
		list = append(list, defaultItem)

	}

	// 打印查询结果
	//fmt.Println(list)
	return list, err

	// 打印查询结果
	//fmt.Println(list)
}

func getCidCardResp(dt string, orgs []int) (list []vboxResp.ChaOrgRealCardResp, err error) {
	// 通道产品维度统计
	queryCidSql := `
		SELECT
			a.channel_code as channelCode,
			'11' as stepTime,
			c.product_id as productId,
			c.product_name as productName,
			a.cnt as orderQuantify,
			a.ok_cnt as okOrderQuantify,
			a.ok_income as okIncome
		FROM
		(
			SELECT
			 channel_code,
			 sum(if (order_status=1 and cb_status=1,money,0)) as ok_income,
			 count(*) as cnt,
			 sum(if(order_status=1 and cb_status=1,1,0)) as ok_cnt
			FROM(
				SELECT
							created_by as uid,
							p_account,
							channel_code,
							order_status,
							cb_status,
							money
					from vbox_pay_order 
					where (DATE_FORMAT(cb_time, '%Y-%m-%d') = ? or DATE_FORMAT(created_at, '%Y-%m-%d') = ?)
					and event_type = 1   and created_by in ?
				
			)aa
			GROUP BY channel_code
		)a 
		left join (
		SELECT DISTINCT channel_code,product_id,product_name 
		FROM vbox_channel_product where channel_code !=''
		) c
		on a.channel_code = c.channel_code  
	`

	//fmt.Println("dt-->", dt, "uid-->", orgs, "querySql-->", queryCidSql)
	db := global.GVA_DB.Model(&vboxResp.ChaOrgRealCardResp{})
	rows, err := db.Raw(queryCidSql, dt, dt, orgs).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.CidRealStatisicsResp
		err := rows.Scan(&item.ChannelCode, &item.StepTime, &item.ProductId,
			&item.ProductName, &item.OrderQuantify, &item.OkOrderQuantify, &item.OkIncome)
		//fmt.Println("OrderQuantify", item.OrderQuantify)
		if err != nil {
			panic(err)
		}
		defaultItem := vboxResp.ChaOrgRealCardResp{
			Title:           item.ProductName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           0,
			Income:          item.OkIncome,
			Dt:              dt,
		}
		list = append(list, defaultItem)
		// 打印查询结果
		fmt.Println(defaultItem)
	}

	return list, err

	// 打印查询结果
	//fmt.Println(list)
}

func getOrgCardResp(dt string, orgs []int, info vboxReq.OrgSelectForm) (list []vboxResp.ChaOrgRealCardResp, err error) {
	// 团队维度统计
	queryOrgSql := `
			
			SELECT
				 organization_id as organizationId,
				 organization_name as organizationName,
				 '' as stepTime,
				 sum(if (order_status=1 and cb_status=1,money,0)) as okIncome,
				 count(*) as orderQuantify,
				 sum(if(order_status=1 and cb_status=1,1,0)) as okOrderQuantify
				FROM(
					select 
						uid,
						channel_code,
						order_status,
						cb_status,
						money,
						organization_id,
						organization_name
					from
					(
						SELECT
								created_by as uid,
								p_account,
								channel_code,
								order_status,
								cb_status,
								money
						from vbox_pay_order 
						where (DATE_FORMAT(cb_time, '%Y-%m-%d') = ? or DATE_FORMAT(created_at, '%Y-%m-%d') = ?)
						and event_type = 1  
					) a1	
					join(
						 SELECT 
								organization_id,
								name as organization_name,
								sys_user_id,
								parent_id
							from
								(
								select
									organization_id,
									sys_user_id
								from
									org_user
								where sys_user_id in ?
							) o1
							join 
							(
								select
									id,
									name,
									parent_id
								from
									organization
							) o2
							on
								o1.organization_id = o2.id
														
					) a2
					on a1.uid = a2.sys_user_id
					
				)aa
				GROUP BY organization_id,organization_name
			`

	//fmt.Println("dt-->", dt, "uid-->", uid, "querySql-->", querySql)
	db := global.GVA_DB.Model(&vboxResp.ChaOrgRealCardResp{})
	rows, err := db.Raw(queryOrgSql, dt, dt, orgs).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.OrgRealStatisicsResp
		err := rows.Scan(&item.OrganizationId, &item.OrganizationName, &item.StepTime, &item.OkIncome, &item.OrderQuantify, &item.OkOrderQuantify)
		if err != nil {
			panic(err)
		}
		defaultItem := vboxResp.ChaOrgRealCardResp{
			Title:           item.OrganizationName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           0,
			Income:          item.OkIncome,
			Dt:              dt,
		}
		list = append(list, defaultItem)

	}

	return list, err

	// 打印查询结果
	//fmt.Println(list)
}

func (bdaChorgService *BdaChorgIndexDService) GetBdaChorgIndexRealListBySelect(info vboxReq.OrgSelectForm) (
	list []vboxResp.ChaOrgRealCardResp, total int64, err error) {
	//var yys []organization.Organization
	//err = global.GVA_DB.Model(&yys).Where("id = ?", info.OrganizationID).Error
	//if err != nil {
	//	panic(err)
	//}
	//var yy = yys[0]
	var parentId int
	var orgs []organization.Organization
	var orgIds []int
	err = global.GVA_DB.Model(&orgs).Select("parent_id").Where("id = ?", info.OrganizationID).Pluck("parent_id", &parentId).Error
	if err != nil {
		panic(err)
	}
	//fmt.Println("OrganizationID pre = ", info.OrganizationID)
	//fmt.Println("parentId = ", parentId)
	if parentId == 0 {
		err = global.GVA_DB.Model(&orgs).Select("id").Where("parent_id = ?", info.OrganizationID).Pluck("id", &orgIds).Error
		if err != nil {
			panic(err)
		}
	} else {
		err = global.GVA_DB.Model(&orgs).Select("id").Where("id = ?", info.OrganizationID).Pluck("id", &orgIds).Error
		if err != nil {
			panic(err)
		}
	}

	dt := time.Now().AddDate(0, 0, 0).Format("2006-01-02")
	uids, err := getYyUids(orgIds)

	//fmt.Println("orgIds = ", orgIds)
	//fmt.Println("uid = ", *info.SysUserID)
	//fmt.Println("OrganizationID = ", info.OrganizationID)
	//fmt.Println("uids = ", uids)
	if info.PAccount != "" {
		//fmt.Println("1", info.PAccount)
		resp, err := getPaccSelectCardResp(dt, uids, info)
		return resp, int64(len(resp)), err
	}
	if info.SysUserID != nil && *info.SysUserID != 0 {
		//fmt.Println("2", info.PAccount)
		resp, err := getUidSelectCardResp(dt, uids, info)
		return resp, int64(len(resp)), err
	}
	if info.Cid != "" {
		//fmt.Println("3", info.PAccount)
		resp, err := getCidSelectCardResp(dt, uids, info)
		return resp, int64(len(resp)), err
	}
	if info.OrganizationID != 0 {
		//fmt.Println("4", info.PAccount)
		resp, err := getOrgSelectCardResp(dt, uids, info)
		return resp, int64(len(resp)), err
	}
	// 默认方法
	resp, err := getOrgSelectCardResp(dt, uids, info)
	return resp, int64(len(resp)), err
}

func getPaccSelectCardResp(dt string, orgs []int, info vboxReq.OrgSelectForm) (list []vboxResp.ChaOrgRealCardResp, err error) {
	// 成员维度统计
	queryPaccSql := `
		SELECT
			a.p_account as pAccount,
			'' as stepTime,
			c.p_remark as userName,
			a.cnt as orderQuantify,
			a.ok_cnt as okOrderQuantify,
			a.ok_income as okIncome
		FROM
		(
			SELECT
			 p_account,
			 sum(if (order_status=1 and cb_status=1,money,0)) as ok_income,
			 count(*) as cnt,
			 sum(if(order_status=1 and cb_status=1,1,0)) as ok_cnt
			FROM(
				SELECT
							created_by as uid,
							p_account,
							channel_code,
							order_status,
							cb_status,
							money	
					from vbox_pay_order 
					where (DATE_FORMAT(cb_time, '%Y-%m-%d') = ? or DATE_FORMAT(created_at, '%Y-%m-%d') = ?)
					and event_type = 1  and created_by in ?
					and p_account = ?
			)aa
			GROUP BY p_account
		)a 
		left join (
		SELECT DISTINCT p_account,p_remark
		FROM vbox_pay_account   
		) c
		on a.p_account = c.p_account  
	`

	//fmt.Println("dt-->", dt, "uid-->", uid, "querySql-->", querySql)
	db := global.GVA_DB.Model(&vboxResp.ChaOrgRealCardResp{})
	rows, err := db.Raw(queryPaccSql, dt, dt, orgs, info.PAccount).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.PaccRealStatisicsResp
		err := rows.Scan(&item.PAccount, &item.StepTime, &item.UserName, &item.OrderQuantify, &item.OkOrderQuantify, &item.OkIncome)
		if err != nil {
			panic(err)
		}
		defaultItem := vboxResp.ChaOrgRealCardResp{
			Title:           item.UserName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           0,
			Income:          item.OkIncome,
			Dt:              dt,
		}
		list = append(list, defaultItem)

	}

	return list, err

	// 打印查询结果
	//fmt.Println(list)
}

func getUidSelectCardResp(dt string, orgs []int, info vboxReq.OrgSelectForm) (list []vboxResp.ChaOrgRealCardResp, err error) {
	// 成员维度统计
	queryUidSql := `

			SELECT
				a.uid as uid,
				'' as stepTime,
				c.nickname as userName,
				a.cnt as orderQuantify,
				a.ok_cnt as okOrderQuantify,
				a.ok_income as okIncome
			FROM
			(
				SELECT
				 uid,
				 sum(if (order_status=1 and cb_status=1,money,0)) as ok_income,
				 count(*) as cnt,
				 sum(if(order_status=1 and cb_status=1,1,0)) as ok_cnt
				FROM(
					SELECT
								created_by as uid,
								p_account,
								channel_code,
								order_status,
								cb_status,
								money
						from vbox_pay_order 
						where (DATE_FORMAT(cb_time, '%Y-%m-%d') = ? or DATE_FORMAT(created_at, '%Y-%m-%d') = ?)
						and event_type = 1   and created_by in ?
						and created_by = ?
					
				)aa
				GROUP BY uid
			)a 
			left join (
			SELECT DISTINCT id,nickname
			FROM sys_users 
			) c
			on a.uid = c.id  
`

	//fmt.Println("dt-->", dt, "uid-->", uid, "querySql-->", querySql)
	db := global.GVA_DB.Model(&vboxResp.ChaOrgRealCardResp{})
	rows, err := db.Raw(queryUidSql, dt, dt, orgs, info.SysUserID).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.UidRealStatisicsResp
		err := rows.Scan(&item.Uid, &item.StepTime, &item.UserName, &item.OrderQuantify, &item.OkOrderQuantify, &item.OkIncome)
		if err != nil {
			panic(err)
		}
		defaultItem := vboxResp.ChaOrgRealCardResp{
			Title:           item.UserName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           0,
			Income:          item.OkIncome,
			Dt:              dt,
		}
		list = append(list, defaultItem)

	}

	// 打印查询结果
	//fmt.Println(list)
	return list, err

	// 打印查询结果
	//fmt.Println(list)
}

func getCidSelectCardResp(dt string, orgs []int, info vboxReq.OrgSelectForm) (list []vboxResp.ChaOrgRealCardResp, err error) {
	// 通道产品维度统计
	queryCidSql := `
		SELECT
			a.channel_code as channelCode,
			'11' as stepTime,
			c.product_id as productId,
			c.product_name as productName,
			a.cnt as orderQuantify,
			a.ok_cnt as okOrderQuantify,
			a.ok_income as okIncome
		FROM
		(
			SELECT
			 channel_code,
			 sum(if (order_status=1 and cb_status=1,money,0)) as ok_income,
			 count(*) as cnt,
			 sum(if(order_status=1 and cb_status=1,1,0)) as ok_cnt
			FROM(
				SELECT
							created_by as uid,
							p_account,
							channel_code,
							order_status,
							cb_status,
							money
					from vbox_pay_order 
					where (DATE_FORMAT(cb_time, '%Y-%m-%d') = ? or DATE_FORMAT(created_at, '%Y-%m-%d') = ?)
					and event_type = 1   and created_by in ?
					and channel_code = ?
				
			)aa
			GROUP BY channel_code
		)a 
		left join (
		SELECT DISTINCT channel_code,product_id,product_name 
		FROM vbox_channel_product where channel_code !=''
		) c
		on a.channel_code = c.channel_code  
	`

	//fmt.Println("dt-->", dt, "uid-->", orgs, "querySql-->", queryCidSql)
	db := global.GVA_DB.Model(&vboxResp.ChaOrgRealCardResp{})
	rows, err := db.Raw(queryCidSql, dt, dt, orgs, info.Cid).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.CidRealStatisicsResp
		err := rows.Scan(&item.ChannelCode, &item.StepTime, &item.ProductId,
			&item.ProductName, &item.OrderQuantify, &item.OkOrderQuantify, &item.OkIncome)
		//fmt.Println("OrderQuantify", item.OrderQuantify)
		if err != nil {
			panic(err)
		}
		defaultItem := vboxResp.ChaOrgRealCardResp{
			Title:           item.ProductName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           0,
			Income:          item.OkIncome,
			Dt:              dt,
		}
		list = append(list, defaultItem)
		// 打印查询结果
		fmt.Println(defaultItem)
	}

	return list, err

	// 打印查询结果
	//fmt.Println(list)
}

func getOrgSelectCardResp(dt string, orgs []int, info vboxReq.OrgSelectForm) (list []vboxResp.ChaOrgRealCardResp, err error) {
	// 团队维度统计
	queryOrgSql := `
			
			SELECT
				 organization_id as organizationId,
				 organization_name as organizationName,
				 '' as stepTime,
				 sum(if (order_status=1 and cb_status=1,money,0)) as okIncome,
				 count(*) as orderQuantify,
				 sum(if(order_status=1 and cb_status=1,1,0)) as okOrderQuantify
				FROM(
					select 
						uid,
						channel_code,
						order_status,
						cb_status,
						money,
						organization_id,
						organization_name
					from
					(
						SELECT
								created_by as uid,
								p_account,
								channel_code,
								order_status,
								cb_status,
								money
						from vbox_pay_order 
						where (DATE_FORMAT(cb_time, '%Y-%m-%d') = ? or DATE_FORMAT(created_at, '%Y-%m-%d') = ?)
						and event_type = 1  
						
					) a1	
					join(
						 SELECT 
								organization_id,
								name as organization_name,
								sys_user_id,
								parent_id
							from
								(
								select
									organization_id,
									sys_user_id
								from
									org_user
								where sys_user_id in ?
							) o1
							join 
							(
								select
									id,
									name,
									parent_id
								from
									organization
								
							) o2
							on
								o1.organization_id = o2.id
														
					) a2
					on a1.uid = a2.sys_user_id
					
				)aa
				GROUP BY organization_id,organization_name
			`

	//fmt.Println("dt-->", dt, "uid-->", uid, "querySql-->", querySql)
	db := global.GVA_DB.Model(&vboxResp.ChaOrgRealCardResp{})
	rows, err := db.Raw(queryOrgSql, dt, dt, orgs).Rows()
	if err != nil {
		panic(err)
	}
	//fmt.Println(rows.Next())
	defer rows.Close()
	// 如果有下一行数据，继续循环
	for rows.Next() {
		//fmt.Println("--->,,")
		// 遍历查询结果并将值映射到结构体中
		var item vboxResp.OrgRealStatisicsResp
		err := rows.Scan(&item.OrganizationId, &item.OrganizationName, &item.StepTime, &item.OkIncome, &item.OrderQuantify, &item.OkOrderQuantify)
		if err != nil {
			panic(err)
		}
		defaultItem := vboxResp.ChaOrgRealCardResp{
			Title:           item.OrganizationName,
			OrderQuantify:   item.OrderQuantify,
			OkOrderQuantify: item.OkOrderQuantify,
			Ratio:           0,
			Income:          item.OkIncome,
			Dt:              dt,
		}
		list = append(list, defaultItem)

	}

	return list, err

	// 打印查询结果
	//fmt.Println(list)
}
