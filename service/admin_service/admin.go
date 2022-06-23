package admin_service

import (
	"ess/model/user"
	"ess/utils/db"
)

func QueryUserByRoll(roll int) (*[]user.User, error) {
	users := []user.User{}
	if err := db.MysqlDB.Where(&user.User{UserRole: user.Role(roll)}).Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}

func QueryAllUser() (*[]user.User, error) {
	users := []user.User{}
	if err := db.MysqlDB.Find(&users).Error; err != nil {
		return nil, err
	}
	return &users, nil
}
