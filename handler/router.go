package handler

import (
	"ess/handler/ping"
	"ess/middlware"
	"ess/utils/setting"

	"github.com/gin-gonic/gin"

	docs "ess/docs"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	docs.SwaggerInfo.BasePath = "/api"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(middlware.CorsMiddleware())

	gin.SetMode(setting.ServerSetting.RunMode)

	api := r.Group("/api",
		middlware.RecoverMiddleware(),
		middlware.ResponseMiddleware(),
		middlware.RewriteToken())

	api.GET("/ping", ping.Ping)

	return r
}
