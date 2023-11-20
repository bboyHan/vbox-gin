package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/geo"
	"github.com/flipped-aurora/gin-vue-admin/server/plugin/organization"
	swaggerFiles "github.com/swaggo/files"
	"net/http"

	"github.com/flipped-aurora/gin-vue-admin/server/docs"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/flipped-aurora/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Routers() *gin.Engine {

	if global.GVA_CONFIG.System.Env == "public" {
		gin.SetMode(gin.ReleaseMode)
	}
	Router := gin.New()

	if global.GVA_CONFIG.System.Env != "public" {
		Router.Use(gin.Logger(), gin.Recovery())
	}

	InstallPlugin(Router)
	systemRouter := router.RouterGroupApp.System
	exampleRouter := router.RouterGroupApp.Example

	Router.StaticFS(global.GVA_CONFIG.Local.StorePath, http.Dir(global.GVA_CONFIG.Local.StorePath))

	Router.Use(middleware.Cors())

	docs.SwaggerInfo.BasePath = global.GVA_CONFIG.System.RouterPrefix
	Router.GET(global.GVA_CONFIG.System.RouterPrefix+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	global.GVA_LOG.Info("register swagger handler")

	PublicGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	{

		PublicGroup.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	{
		systemRouter.InitBaseRouter(PublicGroup)
		systemRouter.InitInitRouter(PublicGroup)
	}
	PrivateGroup := Router.Group(global.GVA_CONFIG.System.RouterPrefix)
	PrivateGroup.Use(middleware.JWTAuth()).Use(middleware.CasbinHandler())
	{
		systemRouter.InitApiRouter(PrivateGroup, PublicGroup)
		systemRouter.InitJwtRouter(PrivateGroup)
		systemRouter.InitUserRouter(PrivateGroup)
		systemRouter.InitMenuRouter(PrivateGroup)
		systemRouter.InitSystemRouter(PrivateGroup)
		systemRouter.InitCasbinRouter(PrivateGroup)
		systemRouter.InitAutoCodeRouter(PrivateGroup)
		systemRouter.InitAuthorityRouter(PrivateGroup)
		systemRouter.InitSysDictionaryRouter(PrivateGroup)
		systemRouter.InitAutoCodeHistoryRouter(PrivateGroup)
		systemRouter.InitSysOperationRecordRouter(PrivateGroup)
		systemRouter.InitSysDictionaryDetailRouter(PrivateGroup)
		systemRouter.InitAuthorityBtnRouterRouter(PrivateGroup)
		systemRouter.InitChatGptRouter(PrivateGroup)

		exampleRouter.InitCustomerRouter(PrivateGroup)
		exampleRouter.InitFileUploadAndDownloadRouter(PrivateGroup)

	}

	PluginInit(PublicGroup, organization.CreateOrganizationPlug())
	PluginInit(PublicGroup, geo.CreateGeoPlug())
	{
		vboxRouter := router.RouterGroupApp.Vbox
		vboxRouter.InitPubAccessRouter(PublicGroup)

		vboxRouter.InitChannelAccountRouter(PrivateGroup)
		vboxRouter.InitChannelProductRouter(PrivateGroup)
		vboxRouter.InitPayOrderRouter(PrivateGroup)
		vboxRouter.InitPayAccountRouter(PrivateGroup)
		vboxRouter.InitVboxProxyRouter(PrivateGroup)
		vboxRouter.InitChannelShopRouter(PrivateGroup)
		vboxRouter.InitOrgProductRouter(PrivateGroup)
		vboxRouter.InitUserWalletRouter(PrivateGroup)
		vboxRouter.InitChannelPayCodeRouter(PrivateGroup)
		vboxRouter.InitBdaChIndexDRouter(PrivateGroup)
		vboxRouter.InitBdaChaccIndexDRouter(PrivateGroup)
		vboxRouter.InitBdaChShopIndexDRouter(PrivateGroup)

	}
	global.GVA_LOG.Info("router register success")
	return Router
}
