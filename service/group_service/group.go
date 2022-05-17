package group_service

import (
	"ess/model/address"
	"ess/model/category"
	"ess/model/group"
	"ess/model/order"
	"ess/utils/db"
	"ess/utils/logging"

	"github.com/jinzhu/copier"
)

func QeuryGroupByName(name string) []group.Group {
	var groups []group.Group
	resinfo := db.MysqlDB.Where("name = ?", name).Find(&groups)
	logging.InfoF("Find %d groups\n", resinfo.RowsAffected)
	return groups
}

func QueryGroupById(gid int) *group.Group {
	var resgroup group.Group
	resinfo := db.MysqlDB.Where("gid = ?", gid).First(&resgroup)
	if resinfo.RowsAffected == 0 {
		logging.InfoF("No Group with gid %d !\n", gid)
	}

	return &resgroup
}

func QueryGroupAddrById(gid int) *address.Address {
	var resaddr address.Address
	var resgroup group.Group
	resinfo := db.MysqlDB.Model(&resgroup).Where("gid = ?", gid).Association("GroupAddress").Find(&resaddr)
	logging.Info(resinfo.Error())
	return &resaddr
}

func CreateGroup(gp *group.Group) error {
	if err := db.MysqlDB.Create(gp).Error; err != nil {
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
	var resgroup []group.Group
	resinfo := db.MysqlDB.Where("gid = ?", gid).Distinct("uid").Find(&resgroup)
	return int(resinfo.RowsAffected)
}

func QueryGroupTotalPriceById(gid int) float64 {
	var resorder []order.Order
	var rescat []category.Category
	var result float64 = 0
	orderinfo := db.MysqlDB.Where("gid = ?", gid).Find(&resorder)
	if orderinfo.RowsAffected == 0 {
		logging.Info("Group Has No Order!\n")
		return 0
	}
	catinfo := db.MysqlDB.Model(&order.Order{}).Where("gid = ?", gid).Association("OrderCategory").Find(&rescat)
	logging.Info(catinfo.Error())
	if len(rescat) != len(resorder) {
		logging.Fatal("category number != order number\n")
		return 0
	}
	for i := range rescat {
		result += rescat[i].CategoryPrice * resorder[i].OrderAmount
	}
	return result
}

func QueryGroupUserPriceById(gid int, uid int) float64 {
	var resorder []order.Order
	var rescat []category.Category
	var result float64 = 0
	orderinfo := db.MysqlDB.Where("gid = ? AND uid = ?", gid, uid).Find(&resorder)
	if orderinfo.RowsAffected == 0 {
		logging.Info("Group Has No Order With This Uid!\n")
		return 0
	}
	catinfo := db.MysqlDB.Model(&order.Order{}).Where("gid = ? AND uid = ?", gid, uid).Association("OrderCategory").Find(&rescat)
	logging.Info(catinfo.Error())
	if len(rescat) != len(resorder) {
		logging.Fatal("category number != order number\n")
		return 0
	}
	for i := range rescat {
		result += rescat[i].CategoryPrice * resorder[i].OrderAmount
	}
	return result
}

func QueryGroupCategories(gid int) *[]group.GroupInfoCommodity {
	var groupcat []category.Category
	var rescat []group.GroupInfoCommodity
	var tmp group.GroupInfoCommodity

	resinfo := db.MysqlDB.Model(&group.Group{}).Where("gid = ?", gid).Find(&groupcat)
	if resinfo.RowsAffected == 0 {
		return &rescat
	}

	for _, catinfo := range groupcat {
		copier.Copy(tmp, catinfo)
		rescat = append(rescat, tmp)
	}
	return &rescat
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
