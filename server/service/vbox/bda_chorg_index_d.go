package vbox

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxResp "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
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
