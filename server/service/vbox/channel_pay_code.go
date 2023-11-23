package vbox

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/vbox"
	vboxReq "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/request"
	vboxResp "github.com/flipped-aurora/gin-vue-admin/server/model/vbox/response"
	"github.com/songzhibin97/gkit/tools/rand_string"
	"gorm.io/gorm"
	"math"
	"time"
)

type ChannelPayCodeService struct {
}

// CreateChannelPayCode 创建通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) CreateChannelPayCode(vboxChannelPayCode *vbox.ChannelPayCode) (err error) {
	mid := time.Now().Format("20060102150405") + rand_string.RandomInt(3)
	vboxChannelPayCode.Mid = mid
	err = global.GVA_DB.Create(vboxChannelPayCode).Error
	return err
}

// DeleteChannelPayCode 删除通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) DeleteChannelPayCode(vboxChannelPayCode vbox.ChannelPayCode) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.ChannelPayCode{}).Where("id = ?", vboxChannelPayCode.ID).Update("deleted_by", vboxChannelPayCode.DeletedBy).Error; err != nil {
			return err
		}
		if err = tx.Delete(&vboxChannelPayCode).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// DeleteChannelPayCodeByIds 批量删除通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) DeleteChannelPayCodeByIds(ids request.IdsReq, deleted_by uint) (err error) {
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&vbox.ChannelPayCode{}).Where("id in ?", ids.Ids).Update("deleted_by", deleted_by).Error; err != nil {
			return err
		}
		if err := tx.Where("id in ?", ids.Ids).Delete(&vbox.ChannelPayCode{}).Error; err != nil {
			return err
		}
		return nil
	})
	return err
}

// UpdateChannelPayCode 更新通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) UpdateChannelPayCode(vboxChannelPayCode vbox.ChannelPayCode) (err error) {
	err = global.GVA_DB.Save(&vboxChannelPayCode).Error
	return err
}

// GetChannelPayCode 根据id获取通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) GetChannelPayCode(id uint) (vboxChannelPayCode vbox.ChannelPayCode, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&vboxChannelPayCode).Error
	return
}

// GetChannelPayCodeInfoList 分页获取通道账户付款二维码记录
// Author [piexlmax](https://github.com/piexlmax)
func (channelPayCodeService *ChannelPayCodeService) GetChannelPayCodeInfoList(info vboxReq.ChannelPayCodeSearch, ids []uint) (list []vbox.ChannelPayCode, total int64, err error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelPayCode{})
	var vboxChannelPayCodes []vbox.ChannelPayCode
	// 如果有条件搜索 下方会自动创建搜索语句
	if info.StartCreatedAt != nil && info.EndCreatedAt != nil {
		db = db.Where("created_at BETWEEN ? AND ?", info.StartCreatedAt, info.EndCreatedAt)
	}
	if info.Cid != "" {
		db = db.Where("cid = ?", info.Cid)
	}
	if info.AcAccount != "" {
		db = db.Where("ac_account = ?", info.AcAccount)
	}
	if info.Location != "" {
		db = db.Where("location = ?", info.Location)
	}
	if info.Operator != "" {
		db = db.Where("operator = ?", info.Location)
	}
	if info.Mid != "" {
		db = db.Where("mid = ?", info.Mid)
	}
	if info.CodeStatus != 0 {
		db = db.Where("code_status = ?", info.CodeStatus)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}

	if limit != 0 {
		db = db.Limit(limit).Offset(offset)
	}

	err = db.Where("created_by in ?", ids).Order("id desc").Find(&vboxChannelPayCodes).Error
	return vboxChannelPayCodes, total, err
}

/*
*
todo 当只选了省或者市包含到下游个数的展示
*/
func (channelPayCodeService *ChannelPayCodeService) GetChannelPayCodeNumsByLocation(info vboxReq.ChannelPayCodeSearch, ids []uint) (list []vboxResp.ChannelPayCodeStatistics, total int64, err error) {
	query := `
		    SELECT
			 code as location,count(mid) as codeNums
			FROM(
			
				SELECT  
					t.mid,
					a.name as code
				from (
					SELECT location, mid
					from vbox_channel_pay_code
					where  location !='' and created_by in ?
				) t 
				join geo_provinces a 
				on a.code = SUBSTRING(t.location,1,?)
			
			)b
			GROUP BY code
		    ORDER BY codeNums desc
		;
		`

	querySubF := `
		    SELECT
			 code as location,count(mid) as codeNums
			FROM(
				SELECT  
					t.mid,
					a.name as code
				from (
					SELECT location, mid
					from vbox_channel_pay_code
					where  location !='' and created_by in ?
					and SUBSTRING(location,1,?) = ?
				) t 
				join geo_cities a 
				on a.code = SUBSTRING(t.location,1,?)
				and LENGTH(location) >= ?
			
			)b
			GROUP BY code
		    ORDER BY codeNums desc
		;
		`
	querySubS := `
		    SELECT
			 code as location,count(mid) as codeNums
			FROM(
				SELECT  
					t.mid,
					a.name as code
				from (
					SELECT location, mid
					from vbox_channel_pay_code
					where  location !='' and created_by in ?
					and SUBSTRING(location,1,?) = ?
					and LENGTH(location) = ?
				) t 
				join geo_areas a 
				on a.code = SUBSTRING(t.location,1,?)
			
			)b
			GROUP BY code
		    ORDER BY codeNums desc
		;
		`
	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelPayCode{})

	var totalGroup int64 = 0

	fmt.Println("info.Location:", info.Location)

	var codeStatisResultList []vboxResp.ChannelPayCodeStatisticsResult
	// 全国各省
	if info.Location == "" {
		fmt.Println("0>>>>>")
		rows, err := db.Raw(query, ids, 2).Rows()
		if err != nil {
			// 处理错误
		}

		defer rows.Close()
		for rows.Next() {
			var result vboxResp.ChannelPayCodeStatisticsResult
			err := rows.Scan(&result.Location, &result.CodeNums)
			if err != nil {
				// 处理错误
			}
			codeStatisResultList = append(codeStatisResultList, result)
			totalGroup += int64(result.CodeNums)
		}

	}

	if len(info.Location) == 2 {
		fmt.Println("4>>>>>")
		rows, err := db.Raw(querySubF, ids, 2, info.Location, 4, 4).Rows()
		if err != nil {
			// 处理错误
		}
		defer rows.Close()
		for rows.Next() {
			var result vboxResp.ChannelPayCodeStatisticsResult
			err := rows.Scan(&result.Location, &result.CodeNums)
			if err != nil {
				// 处理错误
			}
			codeStatisResultList = append(codeStatisResultList, result)
			totalGroup += int64(result.CodeNums)
		}
	}

	if len(info.Location) == 4 {
		fmt.Println("6>>>>>")
		rows, err := db.Raw(querySubS, ids, 4, info.Location, 6, 6).Rows()
		if err != nil {
			// 处理错误
		}
		defer rows.Close()
		for rows.Next() {
			var result vboxResp.ChannelPayCodeStatisticsResult
			err := rows.Scan(&result.Location, &result.CodeNums)
			if err != nil {
				// 处理错误
			}
			codeStatisResultList = append(codeStatisResultList, result)
			totalGroup += int64(result.CodeNums)
		}

	}
	for i, statis := range codeStatisResultList {
		fmt.Println("A num: ", i, "code: ", statis.Location, "total", totalGroup, "codeNums", statis.CodeNums)
		ratio := math.Round(float64(statis.CodeNums)/float64(totalGroup)*10000) / 10000
		entity := vboxResp.ChannelPayCodeStatistics{
			Order:    uint(i + 1),
			Location: statis.Location,
			CodeNums: statis.CodeNums,
			Ratio:    ratio,
		}
		list = append(list, entity)
	}
	total = int64(len(list))
	fmt.Println("total:", total)
	if total == 0 {
		entity := vboxResp.ChannelPayCodeStatistics{
			Order:    1,
			Location: "无",
			CodeNums: 0,
			Ratio:    1,
		}
		list = append(list, entity)
	}
	return list, total, err

}
