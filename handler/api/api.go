package api

import (
	"ess/define"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
)

// @Summary ping
// @Description test connection
// @Tags    api
// @Produce json
// @Success 200 {string} string "'pong'"
// @Router  /ping [get]
func Ping(c *gin.Context) {
	c.Set(define.ESSRESPONSE, response.JSONData("pong"))
}

// @Summary get API version
// @Tags    api
// @Produce json
// @Success 200 {string} string "version"
// @Router  /version [get]
func Version(c *gin.Context) {
	c.Set(define.ESSRESPONSE, response.JSONData(gin.H{
		"version": define.ESSAPIVERSION,
	}))
}
