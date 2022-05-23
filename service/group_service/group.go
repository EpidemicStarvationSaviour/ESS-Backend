package group_service

import (
	"ess/model/group"
	"ess/model/order"
	"ess/service/category_service"
	"ess/utils/db"
	"ess/utils/logging"
)

func QeuryGroupByName(name string) *[]group.Group {
	var groups []group.Group
	resinfo := db.MysqlDB.Where("group_name Like ?", "%"+name+"%").Find(&groups)
	logging.InfoF("Find %d groups with name %s\n", resinfo.RowsAffected, name)
	return &groups
}

func QeuryGroupByCreatorId(cid int) *[]group.Group {
	var groups []group.Group
	resinfo := db.MysqlDB.Where(&group.Group{GroupCreatorId: cid}).Find(&groups)
	logging.InfoF("Find %d groups with creatorID %d\n", resinfo.RowsAffected, cid)
	return &groups
}

func QueryGroupById(gid int) *group.Group {
	var resgroup group.Group
	resinfo := db.MysqlDB.Where(&group.Group{GroupId: gid}).First(&resgroup)
	if resinfo.RowsAffected == 0 {
		logging.InfoF("No Group with gid %d !\n", gid)
	}

	return &resgroup
}

func CreateGroup(gp *group.Group) error {
	if err := db.MysqlDB.Select("group_name", "group_description", "group_remark", "group_creator_id", "group_address_id").Create(gp).Error; err != nil {
		return err
	}
	return nil
}

func UpdateGroup(gp *group.Group) error {
	if err := db.MysqlDB.Select("*").Updates(gp).Error; err != nil {
		return err
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
	var catids []int

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
