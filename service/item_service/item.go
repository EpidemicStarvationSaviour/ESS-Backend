package item_service

import (
	"ess/model/item"
	"ess/utils/db"
)

func DeleteItemByUserId(uid int) error {
	return db.MysqlDB.Where(&item.Item{ItemUserId: uid}).Delete(&item.Item{}).Error
}
