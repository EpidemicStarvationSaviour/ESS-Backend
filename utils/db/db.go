package db

import (
	"ess/model/user"
	"ess/utils/logging"
	"ess/utils/setting"
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	MysqlDB *gorm.DB
)

func Setup() {
	var (
		dbType, dbName, dbUser, password, host, tablePrefix string
		err                                                 error
	)
	dbType = setting.DatabaseSetting.Type
	dbName = setting.DatabaseSetting.DbName
	dbUser = setting.DatabaseSetting.User
	password = setting.DatabaseSetting.Password
	host = setting.DatabaseSetting.Host
	tablePrefix = setting.DatabaseSetting.TablePrefix

	var MySqlDNS = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, password, host, dbName)

	MysqlDB, err = gorm.Open(mysql.Open(MySqlDNS), &gorm.Config{})
	if err != nil {
		logging.Fatal("failed to init db")
	}

	if len(tablePrefix) > 0 {
		fmt.Printf("[warning] tablePrefix '%s' will be nothing to do in current version", tablePrefix)
	}
	if dbType != "mysql" {
		fmt.Printf("[warning] '%s' will be not be use in current version", dbType)
		os.Exit(-1)
	}

	// MysqlDB.Debug()

	// auto migrate  it can't handle the dependency relations, so you need handle it by yourself
	if err = MysqlDB.AutoMigrate(&user.User{}); err != nil {
		logging.Fatal("failed to auto migrate user.User: ", err)
	}

	logging.InfoF("[server] database %s@tcp(%s)/%s connected", dbUser, host, dbName)
}
