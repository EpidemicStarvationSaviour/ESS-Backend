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
	err := db.MysqlDB.Where(&route.Route{RouteUserId: uid}).Order("route_index").Find(&result).Error
	return &result, err
}

func DeleteRouteItemById(rid int) error {
	return db.MysqlDB.Where(&route.RouteItem{RouteId: rid}).Delete(&route.RouteItem{}).Error
}

func DeleteRouteById(rid int) error {
	return db.MysqlDB.Model(&route.Route{}).Where(&route.Route{RouteId: rid}).Delete(&route.RouteItem{}).Error
}

func QueryRouteByGroupId(gid int) (*[]route.Route, error) {
	var result []route.Route
	err := db.MysqlDB.Model(&route.Route{}).Where(&route.Route{RouteGroupId: gid}).Order("route_index").Find(&result).Error
	return &result, err
}

func DeleteRouteByGroupId(gid int) error {
	return db.MysqlDB.Where(&route.Route{RouteGroupId: gid}).Delete(&route.Route{}).Error

}

func QueryRouteItem(rid int) (*[]route.RouteItem, error) {
	var result []route.RouteItem
	err := db.MysqlDB.Where(&route.RouteItem{RouteId: rid}).Find(&result).Error
	return &result, err
}

func QueryGroupTime(gid int) (int64, error) {
	var result int64 = 0
	routes, err := QueryRouteByGroupId(gid)
	if err != nil {
		return 0, err
	}
	for _, rt := range *routes {
		result += rt.RouteEstimatedTime
	}
	return result, nil
}
