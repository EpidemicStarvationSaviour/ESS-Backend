package address_service

import (
	"ess/model/address"
	"ess/model/user"
	"ess/utils/amap_base"
	"ess/utils/db"
)

func CreateAddress(addr *address.Address) error {
	if !addr.AddressCached {
		return db.MysqlDB.Create(addr).Error
	}

	tx := db.MysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Create(addr).Error; err != nil {
		tx.Rollback()
		return err
	}

	var addrs []address.Address
	if err := tx.Where(&address.Address{AddressCached: true}).Find(&addrs).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, other := range addrs {
		if other.AddressId != addr.AddressId {
			dis, err := amap_base.DistanceByCoordination(addr.AddressLng, addr.AddressLat, other.AddressLng, other.AddressLat)
			if err != nil {
				tx.Rollback()
				return err
			}
			err = tx.Create(&address.DistanceCache{DistanceCost: dis, DistanceAId: other.AddressId, DistanceBId: addr.AddressId}).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	return tx.Commit().Error
}

func QueryAddressById(aid int) (*address.Address, error) {
	var addr address.Address
	err := db.MysqlDB.Model(&addr).Where(&address.Address{AddressId: aid}).First(&addr).Error
	return &addr, err
}

func UpdateAddress(addr *address.Address) error {
	return db.MysqlDB.Model(addr).Updates(addr).Error
}

func DeleteAddress(addr *address.Address) error {
	return db.MysqlDB.Delete(addr).Error
}

func QueryAddressesByUserId(uid int) ([]address.Address, error) {
	var addresses []address.Address
	err := db.MysqlDB.Where(&address.Address{AddressUserId: uid}).Find(&addresses).Error
	return addresses, err
}

func QueryDefaultAddressByUserId(uid int) (address.Address, error) {
	var addresses address.Address
	var usr user.User
	db.MysqlDB.Where(&user.User{UserId: uid}).First(&usr)
	err := db.MysqlDB.Where(&address.Address{AddressId: usr.UserDefaultAddressId}).First(&addresses).Error
	return addresses, err
}

func CheckAddressByUserId(aid int, uid int) (bool, error) {
	var count int64
	err := db.MysqlDB.Model(&address.Address{}).Where(&address.Address{AddressId: aid, AddressUserId: uid}).Count(&count).Error
	if err != nil {
		return true, err
	}
	return (count == 1), nil
}

func ModifyDefaultAddressIfNeeded(aid int) error {
	tx := db.MysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var addr address.Address
	if err := tx.Where(&address.Address{AddressId: aid}).First(&addr).Error; err != nil {
		tx.Rollback()
		return err
	}

	var usr user.User
	if err := tx.Where(&user.User{UserId: addr.AddressUserId}).First(&usr).Error; err != nil {
		tx.Rollback()
		return err
	}

	if usr.UserDefaultAddressId == aid {
		var new_addr address.Address
		if err := tx.Where(&address.Address{AddressUserId: usr.UserId}).Not(&address.Address{AddressId: aid}).First(&new_addr).Error; err != nil {
			tx.Rollback()
			return err
		}
		usr.UserDefaultAddressId = new_addr.AddressId
		if err := tx.Model(&usr).Updates(&usr).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func DeleteAddressByUser(uid int) error {
	return db.MysqlDB.Where(&address.Address{AddressUserId: uid}).Delete(&address.Address{}).Error
}

func QueryDistanceCacheByAid(a_aid, b_aid int) (uint64, error) {
	if a_aid > b_aid {
		a_aid, b_aid = b_aid, a_aid
	}
	var cache address.DistanceCache
	err := db.MysqlDB.Where(&address.DistanceCache{DistanceAId: a_aid, DistanceBId: b_aid}).First(&cache).Error
	if err != nil {
		return 0, err
	}
	return cache.DistanceCost, nil
}
