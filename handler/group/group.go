package group

import (
	"ess/define"
	"ess/model/group"
	"ess/model/order"
	"ess/model/user"
	"ess/service/address_service"
	"ess/service/category_service"
	"ess/service/group_service"
	"ess/service/order_service"
	"ess/service/route_service"
	"ess/service/user_service"
	"ess/utils/authUtils"
	"ess/utils/logging"
	"ess/utils/response"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

func GetUserAgentGroup(c *gin.Context, groupcondition group.GroupInfoReq, userID int) {
	// claim, _ := c.Get(define.ESSPOLICY)
	// policy, _ := claim.(authUtils.Policy)
	// var groupcondition group.GroupInfoReq
	// if err := c.ShouldBind(&groupcondition); err != nil {
	// 	c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
	// 	c.Abort()
	// 	return
	// }
	// userID := policy.GetId()

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

		if groupcondition.Type == 0 || retgroup.GroupStatus == group.Status(groupcondition.Type) {
			data, err := group_service.GetGroupDetail(retgroup, userinfo.UserId)
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
				c.Abort()
				return
			}
			result.Count++
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
				// neworder.OrderUserId = ord.OrderUserId
				err := order_service.CreateNewOrder(&neworder)
				if err != nil {
					c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_GROUP_CREATE_FAIL))
					return
				}
			}
		}
	}

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
				data, err := group_service.GetGroupDetail(&retgroup, userID)
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
				data, err := group_service.GetGroupDetail(&retgroup, userID)
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
// @Success 200
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

// @Summary Agent Get Own Groups
// @Tags	group
// @Produce json
// @Param _ query group.GroupInfoReq true "Agent Group Info"
// @Success 200 {object} group.GroupInfoResp
// @Router /group/own [get]
func AgentOwnGroup(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	var groupreq group.GroupInfoReq
	if err := c.ShouldBind(&groupreq); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	userID := policy.GetId()
	var result group.GroupInfoResp
	mygroup := group_service.QeuryGroupByCreatorId(userID)
	for _, grp := range *mygroup {
		if groupreq.Type == 0 || grp.GroupStatus == group.Status(groupreq.Type) {
			result.Count++
			data, err := group_service.GetGroupDetail(&grp, userID)
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
				c.Abort()
				return
			}
			result.Data = append(result.Data, *data)
		}
	}
	c.Set(define.ESSRESPONSE, response.JSONData(&result))

}

// @Summary Agent Edit Group
// @Tags	group
// @Produce json
// @Param id path int true "edit group id"
// @Param data body group.GroupEditReq true "Agent Group Edit"
// @Success 200
// @Router /group/edit/{id} [put]
func EditGroup(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	var newgroup group.Group
	newgroup.GroupId = c.GetInt(c.Param("id"))
	{
		var req_uri group.GroupDetailReq
		if err := c.ShouldBindUri(&req_uri); err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
			return
		}
		newgroup.GroupId = req_uri.GroupId
	}

	newgroup = *group_service.QueryGroupById(newgroup.GroupId)

	var editinfo group.GroupEditReq
	if err := c.ShouldBind(&editinfo); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	userID := policy.GetId()

	copier.Copy(&newgroup, editinfo)

	if newgroup.GroupCreatorId != userID {
		logging.ErrorF("the creator %d does not match the user %d!\n", newgroup.GroupCreatorId, userID)
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	groupords, err := order_service.QueryOrderByGroup(userID)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	for _, deleteusr := range editinfo.GroupDeteledUsers {
		for _, ord := range *groupords {
			if deleteusr == ord.OrderUserId {
				err := order_service.DeleteOrder(&ord)
				if err != nil {
					c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
					c.Abort()
					return
				}

			}
		}
	}
	for i := range editinfo.GroupCommodityIds {
		for j := range newgroup.GroupCategories {
			if editinfo.GroupCommodityIds[i] == newgroup.GroupCategories[j].CategoryId {
				editinfo.GroupCommodityIds[i] = -1
				newgroup.GroupCategories[j].CategoryId = -newgroup.GroupCategories[j].CategoryId
			}
		}
	}
	var idx = 0
	for idx < len(newgroup.GroupCategories) {
		if newgroup.GroupCategories[idx].CategoryId > 0 {
			err := order_service.DeleteOrderByGroupCategory(newgroup.GroupId, newgroup.GroupCategories[idx].CategoryId)
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
				c.Abort()
				return
			}
			newgroup.GroupCategories = append(newgroup.GroupCategories[:idx], newgroup.GroupCategories[idx+1:]...)
		} else {
			newgroup.GroupCategories[idx].CategoryId = -newgroup.GroupCategories[idx].CategoryId
			idx++
		}
	}
	for i := range editinfo.GroupCommodityIds {
		if editinfo.GroupCommodityIds[i] > 0 {
			newgroup.GroupCategories = append(newgroup.GroupCategories, *category_service.QueryCategoryById(editinfo.GroupCommodityIds[i]))
		}
	}
	err = group_service.UpdateGroup(&newgroup)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
}

func GetSupplierGroup(c *gin.Context, groupcondition group.GroupInfoReq, userID int) {
	// userinfo, err := user_service.GetUserById(userID)
	// if err != nil {
	// 	c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
	// 	c.Abort()
	// 	return
	// }
	routes, err := route_service.QueryRouteByUser(userID)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	var result group.GroupInfoSupplierResp
	for _, rt := range *routes {
		retgroup := group_service.QueryGroupById(rt.RouteGroupId)

		if groupcondition.Type == 0 || ((retgroup.GroupStatus == 1 || retgroup.GroupStatus == 2) && groupcondition.Type == 1) || (retgroup.GroupStatus == 3 && groupcondition.Type == 2) || (retgroup.GroupStatus == 4 && groupcondition.Type == 3) {
			result.Count++
			var data group.GroupInfoSupplierData
			copier.Copy(&data, &retgroup)
			creator := user_service.QueryUserById(retgroup.GroupCreatorId)
			data.GroupCreatorPhone = creator.UserPhone
			data.GroupCreatorName = creator.UserName
			items, err := route_service.QueryRouteItem(rt.RouteId)
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
				c.Abort()
				return
			}
			for _, it := range *items {
				cat := category_service.QueryCategoryById(it.RouteItemCategoryId)
				data.GroupTotalPrice += cat.CategoryPrice * it.RouteItemAmount
				var commo group.GroupInfoSupplierCommodity
				copier.Copy(&commo, cat)
				commo.ParentId = cat.CategoryFatherId
				commo.RouteAmount = it.RouteItemAmount
				data.GroupCommodity = append(data.GroupCommodity, commo)
			}
			if retgroup.GroupRiderId != 0 {
				grouprider := user_service.QueryUserById(retgroup.GroupRiderId)
				data.GroupRiderName = grouprider.UserName
				data.GroupRiderPhone = grouprider.UserPhone
				riderpos, err := address_service.QueryAddressById(grouprider.UserDefaultAddressId)
				if err != nil {
					c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
					c.Abort()
					return
				}
				copier.Copy(&data.GroupRiderPos, riderpos)
				data.GroupRiderPos.RouteEstimatedTime = rt.RouteEstimatedTime
			}
			result.GroupData = append(result.GroupData, data)
		}
	}
	c.Set(define.ESSRESPONSE, response.JSONData(&result))
}

func GetRiderGroup(c *gin.Context, groupcondition group.GroupInfoReq, userID int) {
	mygroup, err := group_service.QueryGroupByRider(userID)
	var result group.GroupInfoRiderResp
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	for _, gp := range *mygroup {
		if groupcondition.Type == 0 || groupcondition.Type+1 == int(gp.GroupStatus) {
			var data group.GroupInfoRiderData
			copier.Copy(&data, &gp)
			creator := user_service.QueryUserById(gp.GroupCreatorId)
			data.GroupCreatorName = creator.UserName
			data.GroupCreatorPhone = creator.UserPhone
			gpaddr, err := address_service.QueryAddressById(gp.GroupAddressId)
			if err != nil {
				c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
				c.Abort()
				return
			}
			copier.Copy(&data.GroupCreatorAddress, gpaddr)
			price := group_service.QueryGroupTotalPriceById(gp.GroupId)
			log.Printf("price= %f\n", price)
			data.GroupReward = 0.1 * price
			result.Count++
			result.GroupData = append(result.GroupData, data)
		}
	}
	c.Set(define.ESSRESPONSE, response.JSONData(&result))
}

// @Summary get groups related to me
// @Tags	group
// @Produce json
// @Param _ query group.GroupInfoReq true "Group Condition"
// @Success 200 {object} group.GroupInfoResp
// @Router /group/list [get]
func GroupInfo(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)
	var groupcondition group.GroupInfoReq
	if err := c.ShouldBind(&groupcondition); err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	userID := policy.GetId()
	role := policy.ConvertToUser().UserRole
	log.Printf("userID: %d\n", userID)
	log.Printf("role: %d\n", role)
	switch role {
	case user.Purchaser:
		{
			GetUserAgentGroup(c, groupcondition, userID)
		}
	case user.Leader:
		{
			GetUserAgentGroup(c, groupcondition, userID)
		}
	case user.Supplier:
		{
			GetSupplierGroup(c, groupcondition, userID)
		}
	case user.Rider:
		{
			GetRiderGroup(c, groupcondition, userID)
		}
	}
}

func AgentGetDetail(c *gin.Context, uid int, gid int) {
	var result group.GroupAgentDetail
	gp := group_service.QueryGroupById(gid)
	copier.Copy(&result, gp)
	result.GroupTotalPrice = group_service.QueryGroupTotalPriceById(gid)
	result.GroupUserNumber = group_service.CountGroupUserById(gid)
	creator := user_service.QueryUserById(gp.GroupCreatorId)
	result.GroupCreatorName = creator.UserName
	result.GroupCreatorPhone = creator.UserPhone
	gpaddr, err := address_service.QueryAddressById(gp.GroupAddressId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	copier.Copy(&result.GroupCreatorAddress, gpaddr)
	rider := user_service.QueryUserById(gp.GroupRiderId)
	result.GroupRiderName = rider.UserName
	result.GroupRiderPhone = rider.UserPhone
	rideraddr, err := address_service.QueryAddressById(rider.UserDefaultAddressId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	copier.Copy(&result.GroupRiderPos, rideraddr)
	result.GroupRiderPos.RouteEstimatedTime, err = route_service.QueryGroupTime(gp.GroupId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	CategoryIDs := group_service.QueryGroupCategories(gp.GroupId)

	for _, cid := range *CategoryIDs {
		var commo group.GroupAgentCommodity
		catinfo := category_service.QueryCategoryById(cid)
		copier.Copy(&commo, catinfo)
		commo.Id = catinfo.CategoryFatherId
		grouporders := order_service.QueryOrderByGroupCategory(gp.GroupId, cid)
		commo.TotalAmount = 0
		for _, ord := range *grouporders {
			var catuser group.GroupAgentCommodityUser
			commo.TotalAmount += ord.OrderAmount
			userinfo := user_service.QueryUserById(ord.OrderUserId)
			copier.Copy(&catuser, &userinfo)
			catuser.UserAmount = ord.OrderAmount
			commo.CategoryUser = append(commo.CategoryUser, catuser)
		}
		result.GroupCommodities = append(result.GroupCommodities, commo)
	}
	c.Set(define.ESSRESPONSE, response.JSONData(&result))
}

func RiderGetDetail(c *gin.Context, uid int, gid int) {
	var result group.GroupRiderDetail
	gp := group_service.QueryGroupById(gid)
	copier.Copy(&result, gp)

	creator := user_service.QueryUserById(gp.GroupCreatorId)
	result.GroupCreatorName = creator.UserName
	result.GroupCreatorPhone = creator.UserPhone
	gpaddr, err := address_service.QueryAddressById(gp.GroupAddressId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	copier.Copy(&result.GroupCreatorAddress, gpaddr)
	myuserinfo := user_service.QueryUserById(uid)
	myaddr, err := address_service.QueryAddressById(myuserinfo.UserDefaultAddressId)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}
	copier.Copy(&result.GroupRiderPos, myaddr)
	log.Printf("gid= %d\n", gid)
	grouproutes, err := route_service.QueryRouteByGroupId(gid)
	if err != nil {
		c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
		c.Abort()
		return
	}

	myroutes := *grouproutes
	myroutes = myroutes[:len(myroutes)-1]
	for _, rt := range myroutes {
		var rtdetail group.GroupRiderRoute
		storeinfo := user_service.QueryUserById(rt.RouteUserId)
		rtdetail.RouteUserId = rt.RouteUserId
		rtdetail.RouteUserPhone = storeinfo.UserPhone
		storeaddr, err := address_service.QueryAddressById(storeinfo.UserDefaultAddressId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
			return
		}
		copier.Copy(&rtdetail.RouteUserPos, storeaddr)
		rtdetail.RouteVisited = rt.RouteDone
		if rt.RouteDone {
			rtdetail.RouteVisitedTime = rt.RouteFinishedAt.Unix()
		} else {
			rtdetail.RouteVisitedTime = rt.RouteEstimatedTime
		}
		items, err := route_service.QueryRouteItem(rt.RouteId)
		if err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
			return
		}
		log.Printf("rid: %d\n", rt.RouteId)
		for _, it := range *items {
			var commo group.GroupInfoSupplierCommodity
			commo.CategoryId = it.RouteItemCategoryId
			cat := category_service.QueryCategoryById(it.RouteItemCategoryId)
			commo.CategoryName = cat.CategoryName
			commo.ParentId = cat.CategoryFatherId
			commo.RouteAmount = it.RouteItemAmount
			rtdetail.RouteCommodities = append(rtdetail.RouteCommodities, commo)
		}
		result.GroupRouteDetail = append(result.GroupRouteDetail, rtdetail)
	}

	c.Set(define.ESSRESPONSE, response.JSONData(&result))
}

// @Summary Agent/Rider get detail
// @Tags	group
// @Produce json
// @Param id path int true "edit group id"
// @Success 200
// @Router /group/details/{id} [get]
func GetDetailInfo(c *gin.Context) {
	claim, _ := c.Get(define.ESSPOLICY)
	policy, _ := claim.(authUtils.Policy)

	var GroupId int
	{
		var req_uri group.GroupDetailReq
		if err := c.ShouldBindUri(&req_uri); err != nil {
			c.Set(define.ESSRESPONSE, response.JSONError(response.ERROR_PARAM_FAIL))
			c.Abort()
			return
		}
		GroupId = req_uri.GroupId
	}
	UserId := policy.GetId()

	switch policy.ConvertToUser().UserRole {
	case user.Leader:
		{
			AgentGetDetail(c, UserId, GroupId)
		}
	case user.Rider:
		{
			RiderGetDetail(c, UserId, GroupId)
		}
	}

}
