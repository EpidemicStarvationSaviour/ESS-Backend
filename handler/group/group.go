package group

import (
	"ess/define"
	"ess/model/group"
	"ess/service/group_service"
	"ess/service/order_service"
	"ess/service/user_service"
	"ess/utils/authUtils"
	"ess/utils/response"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

// @Summary get groups I joined
// @Tags	group
// @Produce json
// @Success 200 {object} group.GroupInfoResp
// @Router /group/list [get]
func GetMyGroup(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	var groupcondition group.GroupInfoReq
	if err := c.ShouldBind(&groupcondition); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	userID := policy.GetId()

	userinfo, err := user_service.GetUserById(userID)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	Orders, err := order_service.GetOrderByUser(userID)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	var result group.GroupInfoResp
	for _, order := range *Orders {
		retgroup := group_service.QueryGroupById(order.OrderGroupId)

		if retgroup.GroupStatus == group.Status(groupcondition.Type) {
			var data group.GroupInfoData

			groupaddr := group_service.QueryGroupAddrById(order.OrderGroupId)

			creatorinfo, err := user_service.GetUserById(retgroup.GroupCreatorId)
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
				c.Abort()
				return
			}
			result.Count++
			_ = copier.Copy(data, *retgroup)
			_ = copier.Copy(data, *creatorinfo)
			_ = copier.Copy(data.CreatorAddr, *groupaddr)
			data.UserNumber = group_service.CountGroupUserById(retgroup.GroupId)
			data.TotalPrice = group_service.QueryGroupTotalPriceById(retgroup.GroupId)
			data.TotalMyPrice = group_service.QueryGroupUserPriceById(retgroup.GroupId, userinfo.UserId)
			data.Commodities = *group_service.QueryGroupCategories(retgroup.GroupId)
		}
	}
}
