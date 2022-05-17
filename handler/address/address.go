package address

import (
	"ess/define"
	"ess/model/address"
	"ess/model/user"
	"ess/service/address_service"
	"ess/service/user_service"
	"ess/utils/amap"
	"ess/utils/authUtils"
	"ess/utils/logging"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// @Summary add address
// @Tags    user
// @Produce json
// @Param data body address.AddressCreateReq true "new address information"
// @Success 200 {object} address.AddressCreateResp
// @Router  /user/address [post]
func CreateAddr(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)

	var req address.AddressCreateReq
	if err := c.ShouldBind(&req); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	var addr address.Address
	_ = copier.Copy(&addr, &req)
	addr.AddressUserId = policy.GetId()

	if err := amap.GetCoordination(&addr); err != nil {
		logging.ErrorF("failed to get coordination (addr: %+v): %+v\n", addr, err)
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	if err := address_service.CreateAddress(&addr); err != nil {
		logging.ErrorF("failed to add addresses (addr: %+v): %+v\n", addr, err)
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
		c.Abort()
		return
	}

	if req.IsDefaultAddress {
		if err := user_service.UpdateUser(&user.User{UserId: addr.AddressUserId, UserDefaultAddressId: addr.AddressId}); err != nil {
			logging.ErrorF("failed to update user's default address (aid: %+v): %+v\n", addr.AddressId, err)
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
			c.Abort()
			return
		}
	}

	c.Set(define.ESSRESPONSE, response.JSONData(address.AddressCreateResp{AddressId: addr.AddressId}))
}

// @Summary delete address
// @Tags    user
// @Produce json
// @Param data body address.AddressDeleteReq true "address information"
// @Success 200 {string} string "'success'"
// @Router  /user/address [delete]
func DeleteAddr(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	uid := policy.GetId()

	var req address.AddressDeleteReq
	if err := c.ShouldBind(&req); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	valid, err := address_service.CheckAddressByUserId(req.AddressId, uid)
	if err != nil {
		logging.ErrorF("failed to check address owner(aid: %+v, uid: %+v): %+v\n", req.AddressId, uid, err)
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
		c.Abort()
		return
	}
	if !valid {
		c.Set(define.ESSRESPONSE, response.JSONErrorWithMsg("地址不存在"))
		c.Abort()
		return
	}

	err = address_service.ModifyDefaultAddressIfNeeded(req.AddressId)
	if err != nil {
		logging.ErrorF("failed to modify user's default address (aid: %+v): %+v\n", req.AddressId, err)
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
		c.Abort()
		return
	}

	if err := address_service.DeleteAddress(&address.Address{AddressId: req.AddressId}); err != nil {
		logging.ErrorF("failed to delete address (aid: %+v): %+v\n", req.AddressId, err)
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_DATABASE_QUERY))
		c.Abort()
		return
	}

	c.Set(define.ESSRESPONSE, response.JSONData("success"))
}
