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
	var req user.AuthReq
	if err := c.ShouldBind(&req); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	if (req.Type == "name" && req.Account == setting.AdminSetting.Name) || (req.Type == "phone" && req.Account == setting.AdminSetting.Phone) {
		if req.Secret == setting.AdminSetting.Password {
			adminToken, err := authUtils.GetSysAdminToken()
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_TOKEN_GENERATE_FAIL))
				c.Abort()
				return
			}
			c.SetCookie(define.ESSTOKEN, "Bearer "+adminToken, int(setting.ServerSetting.JwtExpireTime.Seconds()), "/", "", false, true)
			c.Set(define.ESSRESPONSE, response.JSONData(user_service.NewLoginResp(user.User{
				UserName: setting.AdminSetting.Name,
				UserRole: user.SysAdmin,
			}, adminToken, req.Type)))
			return
		} else {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_NOT_ADMIN))
			c.Abort()
			return
		}
	}

	secret := crypto.Password2Secret(req.Secret)
	var queryUser user.User
	if req.Type == "phone" {
		queryUser = user_service.QueryUserByPhone(req.Account)
	} else {
		queryUser = user_service.QueryUserByName(req.Account)
	}

	if queryUser.UserId <= 0 { // not exist
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_USERID))
		return
	}
	// check secret
	if secret == queryUser.UserSecret {
		jwt, err := authUtils.GetUserToken(queryUser)
		if err != nil {
			logging.ErrorF("generate token error for user:%+v\n", queryUser)
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_TOKEN_GENERATE_FAIL))
			c.Abort()
		} else {
			c.SetCookie(define.ESSTOKEN, "Bearer "+jwt, int(setting.ServerSetting.JwtExpireTime.Seconds()), "/", "", false, true)
			c.Set(define.ESSRESPONSE, response.JSONData(user_service.NewLoginResp(queryUser, jwt, req.Type)))
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
	c.Set(define.ESSRESPONSE, response.JSONData(user_service.NewLoginResp(policy.ConvertToUser(), token, "phone")))
}
