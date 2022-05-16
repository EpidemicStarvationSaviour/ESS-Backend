package group_service

import (
	"ess/model/group"
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
