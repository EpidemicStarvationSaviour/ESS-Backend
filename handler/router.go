package handler

import (
	"ess/handler/ping"
	"ess/middlware"
	"ess/utils/setting"

	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.Use(middlware.CorsMiddleware())

	gin.SetMode(setting.ServerSetting.RunMode)

	api := r.Group("/api",
		middlware.RecoverMiddleware(),
		middlware.ResponseMiddleware(),
		middlware.RewriteToken())

	api.GET("/ping", ping.Pong)

	return r
}
