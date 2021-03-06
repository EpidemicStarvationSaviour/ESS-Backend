package group_service

import (
	"ess/model/group"
	"ess/model/order"
	"ess/service/address_service"
	"ess/service/category_service"
	"ess/service/order_service"
	"ess/service/user_service"
	"ess/utils/db"
	"ess/utils/logging"

	"github.com/jinzhu/copier"
)

func QueryGroupByName(name string) *[]group.Group {
	groups := []group.Group{}
	resinfo := db.MysqlDB.Where("group_name LIKE ?", "%"+name+"%").Find(&groups)
	logging.InfoF("Find %d groups with name %s\n", resinfo.RowsAffected, name)
	return &groups
}

func QueryGroupByDescription(des string) *[]group.Group {
	groups := []group.Group{}
	resinfo := db.MysqlDB.Where("group_description LIKE ?", "%"+des+"%").Find(&groups)
	logging.InfoF("Find %d groups with description %s\n", resinfo.RowsAffected, des)
	return &groups
}

func QueryGroupByCreatorId(cid int) *[]group.Group {
	groups := []group.Group{}
	resinfo := db.MysqlDB.Where(&group.Group{GroupCreatorId: cid}).Find(&groups)
	logging.InfoF("Find %d groups with creatorID %d\n", resinfo.RowsAffected, cid)
	return &groups
}

func QueryGroupById(gid int) *group.Group {
	var resgroup group.Group
	resinfo := db.MysqlDB.Preload("GroupCategories").Where(&group.Group{GroupId: gid}).First(&resgroup)
	if resinfo.RowsAffected == 0 {
		logging.InfoF("No Group with gid %d !\n", gid)
	}

	return &resgroup
}

func CreateGroup(gp *group.Group) error {
	if err := db.MysqlDB.Select("group_name", "group_description", "group_remark", "group_creator_id", "group_address_id", "GroupCategories").Create(gp).Error; err != nil {
		return err
	}
	return nil
}

func UpdateGroup(gp *group.Group) error {
	if err := db.MysqlDB.Updates(gp).Error; err != nil {
		return err
	}
	if len(gp.GroupCategories) > 0 {
		if err := db.MysqlDB.Model(&gp).Association("GroupCategories").Replace(gp.GroupCategories); err != nil {
			return err
		}
	}
	return nil
}

func CountGroupUserById(gid int) int {
	var resorder []order.Order
	resinfo := db.MysqlDB.Where(&order.Order{OrderGroupId: gid}).Distinct("order_user_id").Find(&resorder)
	return int(resinfo.RowsAffected)
}

func QueryGroupTotalPriceById(gid int) float64 {
	var resorder []order.Order
	var result float64 = 0
	orderinfo := db.MysqlDB.Where(&order.Order{OrderGroupId: gid}).Find(&resorder)
	if orderinfo.RowsAffected == 0 {
		logging.Info("Group Has No Order!\n")
		return 0
	}

	for _, ord := range resorder {
		cat := category_service.QueryCategoryById(ord.OrderCategoryId)
		result += ord.OrderAmount * cat.CategoryPrice
	}
	return result
}

func QueryGroupUserPriceById(gid int, uid int) float64 {
	var resorder []order.Order
	var result float64 = 0
	orderinfo := db.MysqlDB.Where(&order.Order{OrderGroupId: gid, OrderUserId: uid}).Find(&resorder)
	if orderinfo.RowsAffected == 0 {
		logging.Info("Group Has No Order With This Uid!\n")
		return 0
	}
	for _, ord := range resorder {
		cat := category_service.QueryCategoryById(ord.OrderCategoryId)
		result += cat.CategoryPrice * ord.OrderAmount
	}
	return result
}

func QueryGroupCategories(gid int) *[]int {
	catids := []int{}

	resinfo := db.MysqlDB.Raw("SELECT category_category_id FROM group_category WHERE group_group_id = ?", gid).Scan(&catids)
	if resinfo.RowsAffected == 0 {
		logging.Info("Group Has No Category!\n")
	} else {
		logging.InfoF("Find %d categories in group %d \n", resinfo.RowsAffected, gid)
	}
	return &catids

}

func RiderFinishedCount(uid int) (int64, error) {
	var count int64
	err := db.MysqlDB.Model(&group.Group{}).Where(&group.Group{GroupRiderId: uid, GroupStatus: group.Finished}).Count(&count).Error
	return count, err
}

func PurchaserAndLeaderFinishedCount(uid int) (int64, error) {
	var ret int64
	var gids []int
	err := db.MysqlDB.Model(&order.Order{}).Select("order_group_id").Distinct("order_group_id").Where(&order.Order{OrderUserId: uid}).Find(&gids).Error
	if err != nil {
		return 0, err
	}
	for _, gid := range gids {
		var gp group.Group
		err = db.MysqlDB.Model(&group.Group{}).Where(&group.Group{GroupId: gid}).First(&gp).Error
		if err != nil {
			return 0, err
		}
		if gp.GroupStatus == group.Finished {
			ret++
		}
	}
	return ret, nil
}

func DeleteGroupById(gid int) error {
	return db.MysqlDB.Where(&group.Group{GroupId: gid}).Delete(&group.Group{}).Error
}

func QueryGroupByRider(rid int) (*[]group.Group, error) {
	result := []group.Group{}
	err := db.MysqlDB.Where(&group.Group{GroupRiderId: rid}).Find(&result).Error
	return &result, err
}

func GetGroupDetail(grp *group.Group, uid int) (*group.GroupInfoData, error) {
	var resinfo group.GroupInfoData
	resinfo.Commodities = make([]group.GroupInfoCommodity, 0)
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
	resinfo.UserNumber = CountGroupUserById(grp.GroupId)
	resinfo.TotalPrice = QueryGroupTotalPriceById(grp.GroupId)
	resinfo.TotalMyPrice = QueryGroupUserPriceById(grp.GroupId, uid)
	CategoryIDs := QueryGroupCategories(grp.GroupId)
	var commo group.GroupInfoCommodity
	for _, catid := range *CategoryIDs {
		catinfo := category_service.QueryCategoryById(catid)
		_ = copier.Copy(&commo, catinfo)
		commo.OrderAmount = 0
		commo.TotalAmount = 0
		commo.ParentId = catinfo.CategoryFatherId
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

func DeleteOrderByUserID(gid, uid int) error {
	return db.MysqlDB.Delete(&order.Order{}, "order_group_id = ? and order_user_id = ?", gid, uid).Error
}
