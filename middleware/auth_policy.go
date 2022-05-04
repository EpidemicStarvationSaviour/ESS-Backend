package middleware

import (
	"ess/define"
	"ess/utils/authUtils"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
)

// check the policy and return ERROR_NOT_ADMIN if forbidden
// CAUTION: use it after jwt middleware
func SysAdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, _ := c.Get(define.ESSPOLICY)
		if policy, ok := claim.(authUtils.Policy); !ok || !policy.SysAdminOnly() {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_NOT_ADMIN))
			c.Abort()
			return
		}
	}
}

// check the policy and return ERROR_NOT_ADMIN if forbidden
// CAUTION: use it after jwt middleware
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, _ := c.Get(define.ESSPOLICY)
		if policy, ok := claim.(authUtils.Policy); !ok || !policy.AdminOnly() {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_NOT_ADMIN))
			c.Abort()
			return
		}
	}
}

// check the policy and return ERROR_NOT_LOGIN if forbidden
// CAUTION: use it after jwt middleware
func LoginOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		claim, _ := c.Get(define.ESSPOLICY)
		if policy, ok := claim.(authUtils.Policy); !ok || !policy.LoginOnly() {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_NOT_LOGIN))
			c.Abort()
			return
		}
	}
}
