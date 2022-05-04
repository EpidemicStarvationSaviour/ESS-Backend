package middleware

import (
	"ess/define"
	"strings"

	"github.com/gin-gonic/gin"
)

// copy token in Cookie or Header into Authorization
func RewriteToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if !(len(authHeader) > 0 && strings.HasPrefix(authHeader, "Bearer ")) {
			// find in Cookie
			cookie, err := c.Cookie(define.ESSTOKEN)
			if err == nil {
				// dont't forget to "add Bearer "
				if strings.Contains(cookie, "Bearer ") {
					c.Request.Header.Set("Authorization", cookie)
				} else {
					c.Request.Header.Set("Authorization", "Bearer "+cookie)
				}
			} else {
				// find in Header
				token := c.Request.Header.Get(define.ESSTOKEN)
				if len(token) > 0 {
					// dont't forget to "add Bearer "
					if strings.Contains(token, "Bearer ") {
						c.Request.Header.Set("Authorization", token)
					} else {
						c.Request.Header.Set("Authorization", "Bearer "+token)
					}
				}
			}

		}
	}
}
