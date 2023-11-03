package service

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/model"
	"gorm.io/gorm"
	"strconv"
)

type GeoService struct{}

func (e *GeoService) CreateGeo(geo model.Geo) (g model.Geo, err error) {
	// 写你的业务逻辑
	level := geo.Level
	geo.Level = level + 1
	err = global.GVA_DB.Scopes(searchChildren(model.Geo{Level: level})).Create(&geo).Error
	return geo, err
}

func (e *GeoService) GetGeo(geo model.Geo) (g model.Geo, err error) {
	// 写你的业务逻辑
	err = global.GVA_DB.Scopes(searchChildren(model.Geo{Level: geo.Level - 1})).Where("id = ?", geo.ID).First(&geo).Error
	return geo, err
}

func (e *GeoService) GetGeos(level string, code string) (g []model.Geo, err error) {
	// 写你的业务逻辑
	var Geos []model.Geo
	l, err := strconv.Atoi(level)
	if err != nil {
		l = -1
	}
	err = global.GVA_DB.Scopes(searchChildren(model.Geo{Level: l})).Where("parentCode = ?", code).Find(&Geos).Error
	return Geos, err
}

func (e *GeoService) EditGeo(geo model.Geo) (g model.Geo, err error) {
	// 写你的业务逻辑
	err = global.GVA_DB.Scopes(searchChildren(model.Geo{Level: geo.Level - 1})).Where("code = ?", geo.Code).Updates(&geo).Error
	return geo, err
}

func (e *GeoService) DeleteGeo(geo model.Geo) (err error) {
	// 写你的业务逻辑
	var count int64
	if geo.Level != 2 {
		global.GVA_DB.Scopes(searchChildren(model.Geo{Level: geo.Level})).Where("parentCode = ?", geo.Code).Count(&count)
		if count > 0 {
			return errors.New("存在子区域不允许删除")
		}
	}
	err = global.GVA_DB.Scopes(searchChildren(model.Geo{Level: geo.Level - 1})).Where("id = ?", geo.ID).Delete(&model.Geo{}).Error
	return err
}

func searchChildren(get model.Geo) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {
		switch get.Level {
		case 0:
			return tx.Table("geo_cities")
		case 1:
			return tx.Table("geo_areas")
		default:
			return tx.Table("geo_provinces")
		}
	}
}
