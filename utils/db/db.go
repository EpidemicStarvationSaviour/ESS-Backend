package db

import (
	"ess/model/address"
	"ess/model/category"
	"ess/model/group"
	"ess/model/item"
	"ess/model/order"
	"ess/model/route"
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
	if err = MysqlDB.AutoMigrate(&address.Address{}); err != nil {
		logging.Fatal("failed to auto migrate address.Address: ", err)
	}
	if err = MysqlDB.AutoMigrate(&address.DistanceCache{}); err != nil {
		logging.Fatal("failed to auto migrate address.DistanceCache: ", err)
	}
	if err = MysqlDB.AutoMigrate(&user.User{}); err != nil {
		logging.Fatal("failed to auto migrate user.User: ", err)
	}
	if err = MysqlDB.AutoMigrate(&category.Category{}); err != nil {
		logging.Fatal("failed to auto migrate category.Category: ", err)
	}
	if err = MysqlDB.AutoMigrate(&item.Item{}); err != nil {
		logging.Fatal("failed to auto migrate item.Item: ", err)
	}
	if err = MysqlDB.AutoMigrate(&group.Group{}); err != nil {
		logging.Fatal("failed to auto migrate group.Group: ", err)
	}
	if err = MysqlDB.AutoMigrate(&order.Order{}); err != nil {
		logging.Fatal("failed to auto migrate order.Order: ", err)
	}
	if err = MysqlDB.AutoMigrate(&route.Route{}); err != nil {
		logging.Fatal("failed to auto migrate route.Route: ", err)
	}
	if err = MysqlDB.AutoMigrate(&route.RouteItem{}); err != nil {
		logging.Fatal("failed to auto migrate route.RouteItem: ", err)
	}

	logging.InfoF("[server] database %s@tcp(%s)/%s connected", dbUser, host, dbName)
}
