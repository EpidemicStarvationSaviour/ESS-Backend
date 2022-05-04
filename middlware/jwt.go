package middlware

import (
	"ess/define"
	"ess/utils/authUtils"
	"ess/utils/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		// Note that rewrite_token has been called
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_NOT_LOGIN))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_AUTH_NO_VALID_HEADER))
			c.Abort()
			return
		}

		policy, err := authUtils.ParseToken(parts[1])
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_TOKEN_NOT_VAILD))
			c.Abort()
			return
		} else if !policy.CheckExpired() {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_TOKEN_EXPIRED))
			c.Abort()
			return
		}

		c.Set(define.ESSPOLICY, policy)
		c.Next() // handlers can use c.Get(define.ESSPOLICY) to get the policy
	}

}
