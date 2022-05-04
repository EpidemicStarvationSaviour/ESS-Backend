package token

import (
	"ess/define"
	"ess/model/user"
	"ess/service/user_service"
	"ess/utils/authUtils"
	"ess/utils/crypto"
	"ess/utils/logging"
	"ess/utils/response"
	"ess/utils/setting"

	"github.com/gin-gonic/gin"
)

// @Summary login
// @Tags    token
// @Produce json
// @Param data body user.AuthReq true "login information"
// @Success 200 {object} user.AuthResp
// @Router  /token/login [post]
func Login(c *gin.Context) {
	var userAuth user.AuthReq
	if err := c.ShouldBind(&userAuth); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	if (userAuth.Type == "email" && len(userAuth.Email) == 0) || userAuth.Type == "account" && len(userAuth.Account) == 0 {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	loginStatus := false
	defer c.Set(define.ESSLOGINSTATUS, loginStatus)

	if userAuth.Email == setting.AdminSetting.Email {
		if userAuth.Secret == setting.AdminSetting.Password {
			loginStatus = true //nolint
			adminToken, err := authUtils.GetSysAdminToken()
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_TOKEN_GENERATE_FAIL))
				c.Abort()
				return
			}
			c.SetCookie(define.ESSTOKEN, "Bearer "+adminToken, int(setting.ServerSetting.JwtExpireTime.Seconds()), "/", "", false, true)
			c.Set(define.ESSRESPONSE, response.JSONData(user_service.NewLoginResp(user.User{
				UserEmail: setting.AdminSetting.Email,
				UserName:  setting.AdminSetting.Name,
				UserType:  user.SysAdmin,
			}, adminToken, userAuth.Type)))
			return
		} else {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_NOT_ADMIN))
			c.Abort()
			return
		}
	}

	secret := crypto.Password2Secret(userAuth.Secret)

	var queryUser user.User
	if userAuth.Type == "email" {
		queryUser = user_service.QueryUserByEmail(userAuth.Email)
	} else {
		queryUser = user_service.QueryUserByName(userAuth.Account)

	}

	if queryUser.UserId == 0 { // not exist
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_USERID))
		return
	}
	// check secret
	if secret == queryUser.UserSecret {
		loginStatus = true //nolint
		jwt, err := authUtils.GetUserToken(queryUser)
		if err != nil {
			logging.ErrorF("generate token error for user:%+v\n", queryUser)
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_TOKEN_GENERATE_FAIL))
			c.Abort()
		} else {
			c.SetCookie(define.ESSTOKEN, "Bearer "+jwt, int(setting.ServerSetting.JwtExpireTime.Seconds()), "/", "", false, true)
			c.Set(define.ESSRESPONSE, response.JSONData(user_service.NewLoginResp(queryUser, jwt, userAuth.Type)))
		}
	} else {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PASSWORD))
		c.Abort()
	}
}

// @Summary logout
// @Tags    token
// @Produce json
// @Success 200 {string} string "'logout'"
// @Router  /token/logout [get]
func Logout(c *gin.Context) {
	c.SetCookie(define.ESSTOKEN, "", -1, "/", "", false, true)
	c.Set(define.ESSRESPONSE, response.JSONData("logout"))
}

// @Summary refresh token
// @Tags    token
// @Produce json
// @Success 200 {object} user.AuthResp
// @Router  /token/refresh [post]
func Refresh(c *gin.Context) {
	tmp, _ := c.Get(define.ESSPOLICY)
	policy := tmp.(authUtils.Policy)
	var token string
	if policy.SysAdminOnly() {
		token, _ = authUtils.GetSysAdminToken()
	} else {
		token, _ = authUtils.GetUserToken(policy.ConvertToUser())
	}
	c.SetCookie(define.ESSTOKEN, "Bearer "+token, int(setting.ServerSetting.JwtExpireTime.Seconds()), "/", "", false, true)
	c.Set(define.ESSRESPONSE, response.JSONData(user_service.NewLoginResp(policy.ConvertToUser(), token, "email")))
}
