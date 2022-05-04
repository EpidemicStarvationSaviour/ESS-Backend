package ping

import (
	"ess/define"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
)

func Pong(c *gin.Context) {
	c.Set(define.ESSRESPONSE, response.JSONData(gin.H{
		"message": "pong!",
	}))
}
