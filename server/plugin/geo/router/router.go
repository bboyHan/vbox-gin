package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/api"
	"github.com/gin-gonic/gin"
)

type GeoRouter struct {
}

func (s *GeoRouter) InitGeoRouter(Router *gin.RouterGroup) {
	plugRouter := Router
	plugApi := api.ApiGroupApp.GeoApi
	{
		plugRouter.GET("getGeos", plugApi.GetGeos)
		plugRouter.GET("getGeo", plugApi.GetGeo)
		plugRouter.PUT("editGeo", plugApi.EditGeo)
		plugRouter.POST("createGeo", plugApi.CreateGeo)
		plugRouter.DELETE("deleteGeo", plugApi.DeleteGeo)
	}
}
