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

func QueryAddressesByUserId(uid int) ([]address.Address, error) {
	var addresses []address.Address
	err := db.MysqlDB.Where(&address.Address{AddressUserId: uid}).Find(&addresses).Error
	return addresses, err
}
