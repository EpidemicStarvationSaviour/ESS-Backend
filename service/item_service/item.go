package item_service

import (
	"ess/model/item"
	"ess/utils/db"
)

func DeleteItemByUserId(uid int) error {
	return db.MysqlDB.Where(&item.Item{ItemUserId: uid}).Delete(&item.Item{}).Error
}

func QueryItemsByUserId(uid int) ([]item.Item, error) {
	var items []item.Item
	err := db.MysqlDB.Where(&item.Item{ItemUserId: uid}).Find(&items).Error
	return items, err
}

func QueryItemsByCategoryId(cid int) ([]item.Item, error) {
	var items []item.Item
	err := db.MysqlDB.Where(&item.Item{ItemUserId: cid}).Find(&items).Error
	return items, err
}
