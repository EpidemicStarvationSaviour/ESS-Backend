package user

import (
	"ess/define"
	"ess/model/user"
	"ess/service/address_service"
	"ess/service/user_service"
	"ess/utils/authUtils"
	"ess/utils/crypto"
	"ess/utils/logging"
	"ess/utils/response"
	"ess/utils/setting"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// @Summary get user info
// @Tags    user
// @Produce json
// @Success 200 {object} user.UserInfoResp
// @Router  /user/me [get]
func GetInfo(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)

	if policy.SysAdminOnly() {
		sysAdminResp := user.UserInfoResp{
			UserId:    setting.AdminSetting.UserId,
			UserName:  setting.AdminSetting.Name,
			UserRole:  user.SysAdmin,
			UserPhone: setting.AdminSetting.Phone,
		}
		c.Set(define.ESSRESPONSE, response.JSONData(sysAdminResp))
		return
	}

	userID := policy.GetId()
	userRec := user_service.QueryUserById(userID)
	if userRec.UserId < 0 {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_USER_NOT_FOUND))
		c.Abort()
		return
	}

	addr, err := address_service.QueryAddressesByUserId(userID)
	if err != nil {
		logging.ErrorF("failed to retrieve addresses (uid: %v): %+v\n", userID, err)
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
		c.Abort()
		return
	}

	resp := user.UserInfoResp{}
	copier.Copy(&resp, &userRec)
	for _, v := range addr {
		var address user.UserInfoRespAddress
		copier.Copy(&address, &v)
		address.IsDefaultAddress = (v.AddressId == userRec.UserDefaultAddressId)
		resp.UserAddress = append(resp.UserAddress, address)
	}

	c.Set(define.ESSRESPONSE, response.JSONData(resp))
}

// @Summary modify user info
// @Tags    user
// @Produce json
// @Param data body user.UserModifyReq true "user's new information"
// @Success 200 {string} string "'success'"
// @Router  /user/me [put]
func ModifyInfo(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)

	if policy.SysAdminOnly() {
		c.Set(define.ESSRESPONSE, response.JSONErrorWithMsg("系统管理员不支持修改资料"))
		c.Abort()
		return
	}

	userID := policy.GetId()

	userRec := user_service.QueryUserById(userID)
	oldUser := userRec // clean cache

	if userRec.UserId < 0 {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_USER_NOT_FOUND))
		c.Abort()
		return
	}

	req := user.UserModifyReq{}
	if err := c.ShouldBind(&req); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	if req.UserRole == user.SysAdmin || req.UserRole == user.Admin {
		c.Set(define.ESSRESPONSE, response.JSONErrorWithMsg("不允许提权"))
		c.Abort()
		return
	}

	if req.UserDefaultAddressId != 0 {
		valid, err := address_service.CheckAddressByUserId(req.UserDefaultAddressId, userID)
		if err != nil {
			logging.ErrorF("failed to check address owner(aid: %+v, uid: %+v): %+v\n", req.UserDefaultAddressId, userID, err)
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
			c.Abort()
			return
		}
		if !valid {
			c.Set(define.ESSRESPONSE, response.JSONErrorWithMsg("默认地址不存在"))
			c.Abort()
			return
		}
	}

	copier.Copy(&userRec, &req)

	err := user_service.UpdateUser(&userRec)
	if err != nil {
		logging.ErrorF("failed to retrieve addresses (uid: %v): %+v\n", userID, err)
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_UPDATE_FAIL))
		c.Abort()
		return
	}

	user_service.CleanUserCache(oldUser)

	c.Set(define.ESSRESPONSE, response.JSONData("success"))
}

// @Summary register
// @Tags    user
// @Produce json
// @Param data body user.UserCreateReq true "register information"
// @Success 200 {object} user.UserCreateResp
// @Router  /user/register [post]
func CreateUser(c *gin.Context) {
	var req user.UserCreateReq
	if err := c.ShouldBind(&req); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	addr, valid := user_service.ValidUser(req)
	if !valid {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_NOT_VALID_USER_PARAM))
		c.Abort()
		return
	}

	req.UserSecret = crypto.Password2Secret(req.UserSecret)

	var usr user.User
	copier.Copy(&usr, &req)

	if err := user_service.CreateUserWithAddress(&usr, addr); err != nil {
		logging.ErrorF("failed to create user(%+v): %+v\n", usr, err)
		c.Set(define.ESSRESPONSE, response.JSONErrorWithMsg(err.Error()))
		c.Abort()
		return
	}

	logging.InfoF("create a new user: %+v with address:%+v\n", usr, *addr)

	jwt, err := authUtils.GetUserToken(usr)
	if err != nil {
		logging.ErrorF("generate token error for user:%+v\n", usr)
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_TOKEN_GENERATE_FAIL))
		c.Abort()
	}
	c.SetCookie(define.ESSTOKEN, "Bearer "+jwt, int(setting.ServerSetting.JwtExpireTime.Seconds()), "/", "", false, true)

	resp := user.UserCreateResp{UserId: usr.UserId}
	c.Set(define.ESSRESPONSE, response.JSONData(resp))
}
