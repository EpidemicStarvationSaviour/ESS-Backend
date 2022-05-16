package address_service

import (
	"ess/model/address"
	"ess/utils/db"
)

func CreateAddress(addr *address.Address) error {
	return db.MysqlDB.Create(addr).Error
}

func UpdateAddress(addr *address.Address) error {
	return db.MysqlDB.Model(addr).Updates(addr).Error
}
