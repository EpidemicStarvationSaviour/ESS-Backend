package order_service

import (
	"ess/model/order"
	"ess/utils/db"
)

func GetOrderByUser(uid int) (*[]order.Order, error) {
	var orders []order.Order
	if err := db.MysqlDB.Find(&orders, "uid = ?", uid).Error; err != nil {
		return nil, err
	}
	return &orders, nil
}
