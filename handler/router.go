package handler

import (
	"ess/handler/address"
	"ess/handler/admin"
	api_info "ess/handler/api"
	"ess/handler/category"
	"ess/handler/group"
	"ess/handler/rider"
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

	categoryMod := api.Group("/commodity")
	categoryMod.GET("/list", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), category.GetCategoryList)
	categoryMod.POST("/add", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), category.CreateCate)
	categoryMod.DELETE("/delete", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), category.DeleteCate)
	categoryMod.GET("/details/:id", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), category.GetCateDetail)
	categoryMod.GET("/my", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), category.GetMyCategoryDetails)
	categoryMod.POST("/restock", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), category.ModifyCategoryNumber)

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

	riderMod := api.Group("/rider")
	riderMod.POST("/start", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), rider.RiderStartGetOrder)
	riderMod.POST("/stop", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), rider.RiderStopGetOrder)
	riderMod.POST("/pos", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), rider.RiderUploadAddressPort)
	riderMod.GET("/query", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), rider.RiderQueryNewOrder)
	riderMod.POST("/feedback", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), rider.RiderFeedbackNeworder)
	riderMod.POST("/groupfd", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), rider.OrderFeedback)

	adminMod := api.Group("/admin")
	adminMod.GET("/users", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), admin.GetAllUsers)
	adminMod.DELETE("/users", middleware.AuthenticationMiddleware(), middleware.LoginOnly(), admin.DeleteUser)

	return r
}
