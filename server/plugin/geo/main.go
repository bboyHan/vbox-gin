package geo

import (
	sysModel "github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/router"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/plugin-tool/utils"
	"github.com/gin-gonic/gin"
)

type GeoPlugin struct {
}

func CreateGeoPlug() *GeoPlugin {
	utils.RegisterMenus(
		sysModel.SysBaseMenu{
			Name:      "geoGroup",
			Path:      "geoGroup",
			Hidden:    false,
			Component: "view/routerHolder.vue",
			Sort:      1000,
			Meta:      sysModel.Meta{Title: "城市管理", Icon: "school"},
		},
		sysModel.SysBaseMenu{
			Name:      "geo",
			Path:      "geo",
			Hidden:    false,
			Component: "plugin/geo/view/index.vue",
			Sort:      0,
			Meta:      sysModel.Meta{Title: "城市管理", Icon: "school"},
		},
	)

	utils.RegisterApis(
		sysModel.SysApi{
			Path:        "/geo/getGeos",
			Description: "获取城市列表",
			ApiGroup:    "城市管理",
			Method:      "GET",
		},
		sysModel.SysApi{
			Path:        "/geo/getGeo",
			Description: "获取单一城市",
			ApiGroup:    "城市管理",
			Method:      "GET",
		},
		sysModel.SysApi{
			Path:        "/geo/editGeo",
			Description: "修改城市",
			ApiGroup:    "城市管理",
			Method:      "PUT",
		},
		sysModel.SysApi{
			Path:        "/geo/createGeo",
			Description: "创建城市",
			ApiGroup:    "城市管理",
			Method:      "POST",
		},
		sysModel.SysApi{
			Path:        "/geo/deleteGeo",
			Description: "删除城市",
			ApiGroup:    "城市管理",
			Method:      "DELETE",
		})
	return &GeoPlugin{}
}

func (*GeoPlugin) Register(group *gin.RouterGroup) {
	router.RouterGroupApp.InitGeoRouter(group)
}

func (*GeoPlugin) RouterPath() string {
	return "geo"
}
