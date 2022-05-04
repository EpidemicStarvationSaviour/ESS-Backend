package middleware

import (
	"ess/define"
	"ess/service/ipblocker"
	"ess/utils/logging"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
)

func IPBlock() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if result := ipblocker.IsLoginable(ip); !result {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_IP_BLOCK))
			c.Abort()
			return
		}
		c.Next()
		status, exist := c.Get(define.ESSLOGINSTATUS)
		if statusBool, ok := status.(bool); !exist || !ok || !statusBool {
			logging.InfoF("ip %s try to login and fail\n", ip)
			ipblocker.Fail(ip)
		} else {
			ipblocker.Success(ip)
		}

	}
}
