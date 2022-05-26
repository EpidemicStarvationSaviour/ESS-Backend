package route_service

import (
	"ess/model/route"
	"ess/utils/db"
)

func SupplierFinishedCount(uid int) (int64, error) {
	var count int64
	err := db.MysqlDB.Model(&route.Route{}).Where(&route.Route{RouteUserId: uid, RouteDone: true}).Count(&count).Error
	return count, err
}

func QueryRouteByUser(uid int) (*[]route.Route, error) {
	var result []route.Route
	err := db.MysqlDB.Where(&route.Route{RouteUserId: uid}).Find(&result).Error
	return &result, err
}

func DeleteRouteItemById(rid int) error {
	return db.MysqlDB.Where(&route.RouteItem{RouteId: rid}).Delete(&route.RouteItem{}).Error
}

func DeleteRouteById(rid int) error {
	return db.MysqlDB.Where(&route.Route{RouteId: rid}).Delete(&route.RouteItem{}).Error
}

func QeuryRouteByGroupId(gid int) (*[]route.Route, error) {
	var result []route.Route
	err := db.MysqlDB.Where(&route.Route{RouteGroupId: gid}).Find(&result).Error
	return &result, err
}

func DeleteRouteByGroupId(gid int) error {
	return db.MysqlDB.Where(&route.Route{RouteGroupId: gid}).Delete(&route.Route{}).Error

}
