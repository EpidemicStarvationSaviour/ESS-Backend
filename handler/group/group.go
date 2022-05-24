package group

import (
	"ess/define"
	"ess/model/group"
	"ess/model/order"
	"ess/service/address_service"
	"ess/service/category_service"
	"ess/service/group_service"
	"ess/service/order_service"
	"ess/service/user_service"
	"ess/utils/authUtils"
	"ess/utils/response"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func GetGroupDetail(grp *group.Group, uid int) (*group.GroupInfoData, error) {
	var resinfo group.GroupInfoData
	groupaddr, err := address_service.QueryAddressById(grp.GroupAddressId)
	if err != nil {
		return &resinfo, err
	}

	creatorinfo, err := user_service.GetUserById(grp.GroupCreatorId)
	if err != nil {
		return &resinfo, err
	}

	_ = copier.Copy(&resinfo, grp)
	_ = copier.Copy(&resinfo, &creatorinfo)
	_ = copier.Copy(&resinfo.CreatorAddr, &groupaddr)
	resinfo.UserNumber = group_service.CountGroupUserById(grp.GroupId)
	resinfo.TotalPrice = group_service.QueryGroupTotalPriceById(grp.GroupId)
	resinfo.TotalMyPrice = group_service.QueryGroupUserPriceById(grp.GroupId, uid)
	CategoryIDs := group_service.QueryGroupCategories(grp.GroupId)
	log.Printf("%+v", CategoryIDs)
	var commo group.GroupInfoCommodity
	for _, catid := range *CategoryIDs {
		copier.Copy(&commo, category_service.QueryCategoryById(catid))
		// TODO
		orders := order_service.QueryOrderByGroupCategory(grp.GroupId, catid)
		for _, ord := range *orders {
			commo.TotalAmount += ord.OrderAmount
			if ord.OrderUserId == uid {
				commo.OrderAmount = ord.OrderAmount
			}
		}

		resinfo.Commodities = append(resinfo.Commodities, commo)
	}
	return &resinfo, nil
}

// @Summary get groups I joined conditional
// @Tags	group
// @Produce json
// @Param data body group.GroupInfoReq true "Group Condition"
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

	Orders, err := order_service.QueryOrderByUser(userID)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	var result group.GroupInfoResp
	for _, order := range *Orders {
		retgroup := group_service.QueryGroupById(order.OrderGroupId)

		if retgroup.GroupStatus == group.Status(groupcondition.Type) {
			data, err := GetGroupDetail(retgroup, userinfo.UserId)
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
				c.Abort()
				return
			}
			result.Count++
			// groupinfo, userID
			// groupaddr := group_service.QueryGroupAddrById(order.OrderGroupId)

			// creatorinfo, err := user_service.GetUserById(retgroup.GroupCreatorId)
			// if err != nil {
			// 	c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			// 	c.Abort()
			// 	return
			// }
			// result.Count++
			// _ = copier.Copy(&data, retgroup)
			// _ = copier.Copy(&data, creatorinfo)
			// _ = copier.Copy(&data.CreatorAddr, groupaddr)
			// data.UserNumber = group_service.CountGroupUserById(retgroup.GroupId)
			// data.TotalPrice = group_service.QueryGroupTotalPriceById(retgroup.GroupId)
			// data.TotalMyPrice = group_service.QueryGroupUserPriceById(retgroup.GroupId, userinfo.UserId)
			// data.Commodities = *group_service.QueryGroupCategories(retgroup.GroupId)

			result.Data = append(result.Data, *data)
		}
	}
	c.Set(define.ESSRESPONSE, response.JSONData(&result))
}

// @Summary create a new group
// @Tags	group
// @Produce json
// @Param data body group.GroupCreateReq true "New Group Info"
// @Success 200 {object} group.GroupCreateResp
// @Router /group/create [post]
func LaunchNewGroup(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	var createinfo group.GroupCreateReq
	if err := c.ShouldBind(&createinfo); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	userID := policy.GetId()

	var groupinfo group.Group
	groupinfo.GroupName = createinfo.GroupName
	groupinfo.GroupDescription = createinfo.GroupDescription
	groupinfo.GroupRemark = createinfo.GroupRemark
	groupinfo.GroupAddressId = createinfo.GroupAddressId
	// groupinfo.GroupRiderId = nil

	copier.Copy(&groupinfo, &createinfo)

	groupinfo.GroupCreatorId = userID

	for _, cid := range createinfo.GroupCommodities {
		groupinfo.GroupCategories = append(groupinfo.GroupCategories, *category_service.QueryCategoryById(cid))
	}
	log.Printf("%+v", groupinfo.GroupCategories)
	if err := group_service.CreateGroup(&groupinfo); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_GROUP_CREATE_FAIL))
		return
	}
	var neworder order.Order
	neworder.OrderGroupId = groupinfo.GroupId
	neworder.OrderAmount = 0

	if createinfo.GroupUserGroupId != 0 {
		uids, err := order_service.QueryUidByGroup(createinfo.GroupUserGroupId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_GROUP_CREATE_FAIL))
			return
		}
		for _, uid := range *uids {
			for _, cid := range createinfo.GroupCommodities {
				neworder.OrderCategoryId = cid
				neworder.OrderUserId = uid
				err := order_service.CreateNewOrder(&neworder)
				if err != nil {
					c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_GROUP_CREATE_FAIL))
					return
				}
			}
		}
	}

	// for _, cid := range createinfo.GroupCommodities {
	// 	err := category_service.AddCategoryGroupRelation(groupinfo.GroupId, cid)
	// 	if err != nil {
	// 		logging.Info("INSERT Cat-Group Fail\n")
	// 		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_GROUP_CREATE_FAIL))
	// 		return
	// 	}
	// }

	var res group.GroupCreateResp
	res.GroupId = groupinfo.GroupId

	c.Set(define.ESSRESPONSE, response.JSONData(&res))

}

// @Summary search exist groups
// @Tags	group
// @Produce json
// @Param _ query group.GroupSearchReq true "Search Group Info"
// @Success 200 {object} group.GroupInfoResp
// @Router /group/search [get]
func SearchGroup(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	var searchinfo group.GroupSearchReq
	if err := c.ShouldBind(&searchinfo); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	userID := policy.GetId()
	var result group.GroupInfoResp

	if searchinfo.SearchType == 0 {
		groups := group_service.QeuryGroupByName(searchinfo.SearchValue)

		for _, retgroup := range *groups {
			if searchinfo.GroupType == 0 || retgroup.GroupStatus == group.Status(searchinfo.GroupType) {
				data, err := GetGroupDetail(&retgroup, userID)
				if err != nil {
					c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
					c.Abort()
					return
				}
				result.Count++
				// groupinfo, userID
				// groupaddr := group_service.QueryGroupAddrById(order.OrderGroupId)

				// creatorinfo, err := user_service.GetUserById(retgroup.GroupCreatorId)
				// if err != nil {
				// 	c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
				// 	c.Abort()
				// 	return
				// }
				// result.Count++
				// _ = copier.Copy(&data, retgroup)
				// _ = copier.Copy(&data, creatorinfo)
				// _ = copier.Copy(&data.CreatorAddr, groupaddr)
				// data.UserNumber = group_service.CountGroupUserById(retgroup.GroupId)
				// data.TotalPrice = group_service.QueryGroupTotalPriceById(retgroup.GroupId)
				// data.TotalMyPrice = group_service.QueryGroupUserPriceById(retgroup.GroupId, userinfo.UserId)
				// data.Commodities = *group_service.QueryGroupCategories(retgroup.GroupId)

				result.Data = append(result.Data, *data)
			}
		}

		log.Printf("%+v", result)

		c.Set(define.ESSRESPONSE, response.JSONData(&result))
		return
	}

	if searchinfo.SearchType == 1 {
		CreatorId := user_service.Name2Id(searchinfo.SearchValue)
		groups := group_service.QeuryGroupByCreatorId(CreatorId)

		for _, retgroup := range *groups {
			if retgroup.GroupStatus == group.Status(searchinfo.GroupType) {
				data, err := GetGroupDetail(&retgroup, userID)
				if err != nil {
					c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
					c.Abort()
					return
				}
				result.Count++
				// groupinfo, userID
				// groupaddr := group_service.QueryGroupAddrById(order.OrderGroupId)

				// creatorinfo, err := user_service.GetUserById(retgroup.GroupCreatorId)
				// if err != nil {
				// 	c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
				// 	c.Abort()
				// 	return
				// }
				// result.Count++
				// _ = copier.Copy(&data, retgroup)
				// _ = copier.Copy(&data, creatorinfo)
				// _ = copier.Copy(&data.CreatorAddr, groupaddr)
				// data.UserNumber = group_service.CountGroupUserById(retgroup.GroupId)
				// data.TotalPrice = group_service.QueryGroupTotalPriceById(retgroup.GroupId)
				// data.TotalMyPrice = group_service.QueryGroupUserPriceById(retgroup.GroupId, userinfo.UserId)
				// data.Commodities = *group_service.QueryGroupCategories(retgroup.GroupId)

				result.Data = append(result.Data, *data)
			}
		}

		c.Set(define.ESSRESPONSE, response.JSONData(&result))
		return
	}
}

// @Summary Join a group (create new order)
// @Tags	group
// @Produce json
// @Param data body group.GroupJoinReq true "Join Group Info"
// @Success 200 {object} group.GroupInfoResp
// @Router /group/join [post]
func JoinGroup(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	var joininfo group.GroupJoinReq
	if err := c.ShouldBind(&joininfo); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	userID := policy.GetId()
	// joingroup = group_service.QueryGroupById(joininfo.GroupId)
	groupuserIDs, err := order_service.QueryUidByGroup(joininfo.GroupId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	for _, usrid := range *groupuserIDs {
		if usrid == userID {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
			return
		}
	}
	for _, joindata := range joininfo.OrderData {
		var neworder order.Order
		copier.Copy(&neworder, &joindata)
		neworder.OrderGroupId = joininfo.GroupId
		neworder.OrderUserId = userID
		err := order_service.CreateNewOrder(&neworder)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
			return
		}

	}
}
