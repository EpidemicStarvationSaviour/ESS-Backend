package ping

import (
	"ess/define"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
)

// @Summary ping example
// @Description test connection
// @Success 200 {string} pong!
// @Router /ping [get]
func Ping(c *gin.Context) {
	c.Set(define.ESSRESPONSE, response.JSONData(gin.H{
		"message": "pong!",
	}))
}
