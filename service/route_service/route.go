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
