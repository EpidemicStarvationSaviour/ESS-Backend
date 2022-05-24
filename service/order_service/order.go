package order_service

import (
	"ess/model/order"
	"ess/utils/db"
	"ess/utils/logging"
)

func QueryOrderByUser(uid int) (*[]order.Order, error) {
	var orders []order.Order
	if err := db.MysqlDB.Where(&order.Order{OrderUserId: uid}).Find(&orders).Error; err != nil {
		return nil, err
	}
	return &orders, nil
}

func CreateNewOrder(ord *order.Order) error {
	if err := db.MysqlDB.Create(ord).Error; err != nil {
		return err
	}
	return nil
}

func QueryUidByGroup(gid int) (*[]int, error) {
	var resid []int
	if err := db.MysqlDB.Model(&order.Order{}).Select([]string{"order_user_id"}).Distinct([]string{"order_user_id"}).Where(&order.Order{OrderGroupId: gid}).Find(resid).Error; err != nil {
		return &resid, err
	}
	return &resid, nil
}

func QueryOrderByGroupCategory(gid int, cid int) *[]order.Order {
	var resorder []order.Order
	orderinfo := db.MysqlDB.Where(&order.Order{OrderGroupId: gid, OrderCategoryId: cid}).Find(&resorder)
	if orderinfo.RowsAffected == 0 {
		logging.InfoF("Group Has No Order With gid %d and cid %d!\n", gid, cid)

	}
	return &resorder
}

func QueryOrderByGroup(gid int) (*[]order.Order, error) {
	var orders []order.Order
	if err := db.MysqlDB.Where(&order.Order{OrderGroupId: gid}).Find(&orders).Error; err != nil {
		return nil, err
	}
	return &orders, nil
}

func DeleteOrder(ord *order.Order) error {
	return db.MysqlDB.Delete(&ord).Error
}

func DeleteOrderByGroupCategory(gid int, cid int) error {
	return db.MysqlDB.Where(&order.Order{OrderGroupId: gid, OrderCategoryId: cid}).Delete(&order.Order{}).Error
}
