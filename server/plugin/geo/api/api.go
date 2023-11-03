package api

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/model"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GeoApi struct{}

// @Tags Geo
// @Summary 请手动填写接口功能
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /geo/getGeos [get]
func (p *GeoApi) GetGeos(c *gin.Context) {
	level := c.Query("level")
	code := c.Query("code")
	if geos, err := service.ServiceGroupApp.GetGeos(level, code); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("失败", c)
	} else {
		response.OkWithData(geos, c)
	}
}

// @Tags Geo
// @Summary 请手动填写接口功能
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /geo/createGeo [post]
func (p *GeoApi) CreateGeo(c *gin.Context) {
	var geo model.Geo
	c.ShouldBindJSON(&geo)
	if g, err := service.ServiceGroupApp.CreateGeo(geo); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("失败", c)
	} else {
		response.OkWithData(g, c)
	}
}

// @Tags Geo
// @Summary 请手动填写接口功能
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /geo/getGeo [get]
func (p *GeoApi) GetGeo(c *gin.Context) {
	var geo model.Geo
	c.ShouldBindQuery(&geo)
	if g, err := service.ServiceGroupApp.GetGeo(geo); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("失败", c)
	} else {
		response.OkWithData(g, c)
	}
}

// @Tags Geo
// @Summary 请手动填写接口功能
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /geo/editGeo [put]
func (p *GeoApi) EditGeo(c *gin.Context) {
	var geo model.Geo
	c.ShouldBindJSON(&geo)
	if g, err := service.ServiceGroupApp.EditGeo(geo); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage("失败", c)
	} else {
		response.OkWithData(g, c)
	}
}

// @Tags Geo
// @Summary 请手动填写接口功能
// @Produce  application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"发送成功"}"
// @Router /geo/deleteGeo [delete]
func (p *GeoApi) DeleteGeo(c *gin.Context) {
	var geo model.Geo
	c.ShouldBindQuery(&geo)
	if err := service.ServiceGroupApp.DeleteGeo(geo); err != nil {
		global.GVA_LOG.Error("失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
