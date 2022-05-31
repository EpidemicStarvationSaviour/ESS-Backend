package handler

import (
	"ess/handler/address"
	"ess/handler/admin"
	api_info "ess/handler/api"
	"ess/handler/group"
	"ess/handler/token"
	"ess/handler/user"
	"ess/middleware"
	"ess/utils/setting"
	"ess/utils/swagger"

	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	swagger.Setup()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	r.Use(middleware.CorsMiddleware())

	gin.SetMode(setting.ServerSetting.RunMode)

	api := r.Group("/api",
		middleware.RecoverMiddleware(),
		middleware.ResponseMiddleware(),
		middleware.RewriteToken())

	api.GET("/ping", api_info.Ping)
	api.GET("/version", api_info.Version)

	userMod := api.Group("/user")
	userMod.GET("/me", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), user.GetInfo)
	userMod.POST("/register", user.CreateUser)
	userMod.PUT("/me", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), user.ModifyInfo)
	userMod.POST("/address", middleware.AuthenticationMiddleware(), middleware.PurchaserOnly(), address.CreateAddr)
	userMod.DELETE("/address", middleware.AuthenticationMiddleware(), middleware.PurchaserOnly(), address.DeleteAddr)
	userMod.GET("/workinfo", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), user.GetDashboard)

	tokenMod := api.Group("/token")
	tokenMod.POST("/login", token.Login)
	tokenMod.GET("/logout", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), token.Logout)
	tokenMod.POST("/refresh", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), token.Refresh)

	groupMod := api.Group("/group")
	groupMod.GET("/list", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), group.GroupInfo)
	groupMod.POST("/create", middleware.AuthenticationMiddleware(), middleware.LeaderOnly(), group.LaunchNewGroup)
	groupMod.GET("/search", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), group.SearchGroup)
	groupMod.POST("/join", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), group.JoinGroup)
	groupMod.GET("/own", middleware.AuthenticationMiddleware(), middleware.LeaderOnly(), group.AgentOwnGroup)
	groupMod.PUT("/edit/:id", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), group.EditGroup)
	groupMod.GET("/details/:id", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), group.GetDetailInfo)

	adminMod := api.Group("/admin")
	adminMod.GET("/users", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), admin.GetAllUsers)
	adminMod.DELETE("/users", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), admin.DeleteUser)

	return r
}
