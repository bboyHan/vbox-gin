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
					and operator in ? and code_status in ?
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
			    location,
			   codeNums
			from (
			SELECT code AS location,
				count( mid ) AS codeNums 
			FROM
				(
				SELECT
					t.mid,
					a.name AS code 
				FROM
					( SELECT 
					      location,
					      mid 
					  FROM vbox_channel_pay_code 
					  WHERE location != '' 
						   AND created_by IN ? 
						   AND SUBSTRING( location, 1,? ) = ?
					  	and operator in ? and code_status in ?
					  ) t
					JOIN geo_cities a ON a.code = SUBSTRING( t.location, 1,? ) 
					AND LENGTH( location ) >= ? 
				) b 
			GROUP BY code
			UNION ALL
			SELECT
				'无具体省市区' AS location,
				count( mid ) AS codeNums 
			FROM
				( SELECT location, mid FROM vbox_channel_pay_code
			   WHERE location != '' AND created_by IN ? AND location = ?
			   and operator in ? and code_status in ?
			   ) c 
			) d 
		where codeNums > 0
		ORDER BY
				codeNums DESC
		;
		`
	querySubS := `
			SELECT
				location,
				codeNums 
			FROM
				(
				SELECT code AS location,
					count(mid) AS codeNums 
				FROM
					(
					SELECT
						t.mid,
						a.name AS code 
					FROM
						(
						SELECT
							location,
							mid 
						FROM
							vbox_channel_pay_code 
						WHERE location != '' 
							AND created_by IN ? 
							AND SUBSTRING( location, 1,? ) = ? 
							AND LENGTH( location ) = ? 
						and operator in ? and code_status in ?
						) t
						JOIN geo_areas a ON a.code = SUBSTRING( t.location, 1,? ) 
					) b 
				GROUP BY code
				UNION ALL
				SELECT
					'无具体省市区' AS location,
					count( mid ) AS codeNums 
				FROM
					( SELECT location, mid FROM vbox_channel_pay_code
				   WHERE location != '' AND created_by IN ? AND location = ? 
				   and operator in ? and code_status in ?
				   ) c 
				) d 
			where codeNums > 0
			ORDER BY
				codeNums DESC
		;
		`

	querySubCity := `
		SELECT
			a.name AS location ,
			codeNums
		FROM
			(
			SELECT
				location,
				count(mid) as codeNums 
			FROM
				vbox_channel_pay_code 
			WHERE location != '' 
				AND created_by IN ? 
				AND location = ?
			 and operator in ? and code_status in ?
			group by location
			) t
			JOIN geo_areas a
			ON a.code =  t.location 
;
		`

	fmt.Println("CodeStatus=", info.CodeStatus, "operators", info.Operator)
	codeStatus := []uint{1, 2, 3}
	if info.CodeStatus != 0 {
		codeStatus = []uint{info.CodeStatus}
	}
	operators := []string{"yidong", "liantong", "dianxin"}
	if info.Operator != "" {
		operators = []string{info.Operator}
	}

	// 创建db
	db := global.GVA_DB.Model(&vbox.ChannelPayCode{})

	var totalGroup int64 = 0

	//fmt.Println("info.Location:", info.Location)

	var codeStatisResultList []vboxResp.ChannelPayCodeStatisticsResult
	// 全国各省
	if info.Location == "" {
		//fmt.Println("0 >>>>>")
		rows, err := db.Raw(query, ids, operators, codeStatus, 2).Rows()
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
		//fmt.Println("2 >>>>>")
		rows, err := db.Raw(querySubF, ids, 2, info.Location, operators, codeStatus, 4, 4, ids, info.Location, operators, codeStatus).Rows()
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
		//fmt.Println("4 >>>>>")
		rows, err := db.Raw(querySubS, ids, 4, info.Location, 6, operators, codeStatus, 6, ids, info.Location, operators, codeStatus).Rows()
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
	if len(info.Location) == 6 {
		//fmt.Println("6 >>>>>")
		rows, err := db.Raw(querySubCity, ids, info.Location, operators, codeStatus).Rows()
		if err != nil {
			// 处理错误
		}
		defer rows.Close()
		for rows.Next() {
			var result vboxResp.ChannelPayCodeStatisticsResult
			// scan中有严格的字段顺序 要和sql中的一致
			err := rows.Scan(&result.Location, &result.CodeNums)
			if err != nil {
				// 处理错误
			}
			codeStatisResultList = append(codeStatisResultList, result)
			totalGroup += int64(result.CodeNums)
		}

	}
	for i, statis := range codeStatisResultList {
		//fmt.Println("A num: ", i, "code: ", statis.Location, "total", totalGroup, "codeNums", statis.CodeNums)
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
