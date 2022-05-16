package common_service

import (
	"ess/utils/db"
)

func DatabaseCount(model interface{}) (int64, error) {
	var count int64
	err := db.MysqlDB.Model(&model).Count(&count).Error
	return count, err
}
