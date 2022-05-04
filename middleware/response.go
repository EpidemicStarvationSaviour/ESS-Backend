package middleware

import (
	"ess/define"
	"ess/utils/logging"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
)

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		value, exist := c.Get(define.ESSRESPONSE)
		if !exist {
			logging.Warn("response not set")
			return
		}
		resp, ok := value.(response.Response)
		if !ok {
			logging.Warn("response type invalid!")
			return
		}
		resp.Write(c)
	}
}
