package group_service

import (
	"ess/model/group"
	"ess/model/order"
	"ess/utils/db"
	"ess/utils/logging"
)

func QeuryGroupByName(name string) []group.Group {
	var groups []group.Group
	resinfo := db.MysqlDB.Where("name = ?", name).Find(&groups)
	logging.InfoF("Find %d groups\n", resinfo.RowsAffected)
	return groups
}

func QueryGroupById(gid int) group.Group {
	var resgroup group.Group
	resinfo := db.MysqlDB.Where("gid = ?", gid).First(&resgroup)
	if resinfo.RowsAffected == 0 {
		logging.InfoF("No Group with gid %d !\n", gid)
	}

	return resgroup
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

func RiderFinishedCount(uid int) (int64, error) {
	var count int64
	err := db.MysqlDB.Where(&group.Group{GroupRiderId: uid, GroupStatus: group.Finished}).Count(&count).Error
	return count, err
}

func PurchaserAndLeaderFinishedCount(uid int) (int64, error) { // FIX(TO/GA)
	var ret int64
	var gids []int
	err := db.MysqlDB.Select("order_group_id").Distinct("order_group_id").Where(&order.Order{OrderUserId: uid}).Find(&gids).Error
	if err != nil {
		return 0, err
	}
	for _, gid := range gids {
		var gp group.Group
		err = db.MysqlDB.Where(&group.Group{GroupId: gid}).First(&gp).Error
		if err != nil {
			return 0, err
		}
		if gp.GroupStatus == group.Finished {
			ret++
		}
	}
	return ret, nil
}
