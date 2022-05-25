package admin

import (
	"ess/define"
	"ess/model/admin"
	"ess/model/group"
	"ess/model/user"
	"ess/service/address_service"
	"ess/service/admin_service"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// @Summary get all users conditionally
// @Tags	admin
// @Produce json
// @Param _ query group.GroupInfoReq true "User Condition"
// @Success 200 {object} admin.AllUserResp
// @Router /admin/users [get]
func GetAllUsers(c *gin.Context) {
	// claim, _ := c.Get(define.ESSPOLICY)
	// policy, _ := claim.(authUtils.Policy)
	var UserCondition group.GroupInfoReq
	if err := c.ShouldBind(&UserCondition); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	var users *[]user.User
	var err error

	switch UserCondition.Type {
	case 0:
		users, err = admin_service.QueryUserByRoll(1)
	case 1:
		users, err = admin_service.QueryUserByRoll(2)
	case 2:
		users, err = admin_service.QueryUserByRoll(3)
	case 3:
		users, err = admin_service.QueryUserByRoll(4)
	case 4:
		users, err = admin_service.QueryAllUser()
	}

	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	var result admin.AllUserResp

	for _, usr := range *users {
		result.UserCount++
		var data admin.UserData
		copier.Copy(&data, &usr)
		data.UserRole++
		Adds, err := address_service.QueryAddressesByUserId(usr.UserId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
			return
		}
		for _, addr := range Adds {
			var usraddr admin.UserAddress
			copier.Copy(&usraddr, &addr)
			usraddr.AddressIsDefault = (addr.AddressId == usr.UserDefaultAddressId)
			data.UserAddresses = append(data.UserAddresses, usraddr)
		}
		result.UserList = append(result.UserList, data)
	}
	c.Set(define.ESSRESPONSE, response.JSONData(&result))
}
