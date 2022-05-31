package admin

import (
	"ess/define"
	"ess/model/admin"
	"ess/model/group"
	"ess/model/user"
	"ess/service/address_service"
	"ess/service/admin_service"
	"ess/service/group_service"
	"ess/service/item_service"
	"ess/service/order_service"
	"ess/service/route_service"
	"ess/service/user_service"
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

// @Summary delete certain user
// @Tags	admin
// @Produce json
// @Param data body admin.AdminDeleteUser true "User Id"
// @Success 200
// @Router /admin/users [delete]
func DeleteUser(c *gin.Context) {
	// log.Print("0\n")
	var DeleteUserId admin.AdminDeleteUser
	if err := c.ShouldBind(&DeleteUserId); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	// log.Print("1\n")
	// order
	err := order_service.DeleteOrderByUser(DeleteUserId.UserId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	// log.Print("2\n")
	// route
	deleteroutes, err := route_service.QueryRouteByUser(DeleteUserId.UserId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	for _, rt := range *deleteroutes {
		err := route_service.DeleteRouteItemById(rt.RouteId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
			return
		}
		err = route_service.DeleteRouteById(rt.RouteId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
			return
		}
	}
	// log.Print("3\n")
	// group
	createdgroup := group_service.QeuryGroupByCreatorId(DeleteUserId.UserId)
	for _, gp := range *createdgroup {
		rts, err := route_service.QeuryRouteByGroupId(gp.GroupId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
		}

		if len(*rts) > 0 {
			for _, rt := range *rts {
				err = route_service.DeleteRouteItemById(rt.RouteId)
				if err != nil {
					c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
					c.Abort()
				}
			}
		}

		err = route_service.DeleteRouteByGroupId(gp.GroupId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
		}
		err = group_service.DeleteGroupById(gp.GroupId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
		}
	}

	// log.Print("4\n")
	// address
	err = address_service.DeleteAddressByUser(DeleteUserId.UserId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
	}

	// log.Print("5\n")
	// item
	err = item_service.DeleteItemByUserId(DeleteUserId.UserId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
	}

	// log.Print("6\n")
	// user
	err = user_service.DeleteUserById(DeleteUserId.UserId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
	}
	c.Set(define.ESSRESPONSE, response.JSONData("success"))

}
