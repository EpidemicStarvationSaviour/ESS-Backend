package route_service

import (
	"ess/model/category"
	"ess/model/route"
	"ess/service/group_service"
	"ess/service/user_service"
	"ess/utils/amap"
	"ess/utils/db"
	"fmt"
	"time"
)

func SupplierFinishedCount(uid int) (int64, error) {
	var count int64
	err := db.MysqlDB.Model(&route.Route{}).Where(&route.Route{RouteUserId: uid, RouteDone: true}).Count(&count).Error
	return count, err
}

func QueryRouteByUser(uid int) (*[]route.Route, error) {
	result := []route.Route{}
	err := db.MysqlDB.Where(&route.Route{RouteUserId: uid}).Order("route_index").Find(&result).Error
	return &result, err
}

func QueryRouteByUserAndGroup(uid, gid int) (route.Route, error) {
	var result route.Route
	err := db.MysqlDB.Where(&route.Route{RouteUserId: uid, RouteGroupId: gid}).First(&result).Error
	return result, err
}

func DeleteRouteItemById(rid int) error {
	return db.MysqlDB.Where(&route.RouteItem{RouteId: rid}).Delete(&route.RouteItem{}).Error
}

func DeleteRouteById(rid int) error {
	return db.MysqlDB.Model(&route.Route{}).Where(&route.Route{RouteId: rid}).Delete(&route.RouteItem{}).Error
}

func QueryRouteByGroupId(gid int) (*[]route.Route, error) {
	result := []route.Route{}
	err := db.MysqlDB.Model(&route.Route{}).Where(&route.Route{RouteGroupId: gid}).Order("route_index").Find(&result).Error
	return &result, err
}

func DeleteRouteByGroupId(gid int) error {
	return db.MysqlDB.Where(&route.Route{RouteGroupId: gid}).Delete(&route.Route{}).Error

}
func GetRouteItemTotalPrice(gid int) (float64, error) {
	var result []route.Route
	var pri float64
	pri = 0
	err := db.MysqlDB.Where(&route.Route{RouteGroupId: gid}).Order("route_id").Find(&result).Error
	if err != nil {
		return 0, err
	}
	for _, res := range result {
		var resu []route.RouteItem
		err := db.MysqlDB.Where(&route.RouteItem{RouteId: res.RouteId}).Find(&resu).Error
		if err != nil {
			return 0, err
		}
		for _, resul := range resu {
			var cat category.Category
			err := db.MysqlDB.Where(&category.Category{CategoryId: resul.RouteItemCategoryId}).First(&cat).Error
			if err != nil {
				return 0, err
			}
			pri = pri + cat.CategoryPrice*resul.RouteItemAmount
		}
	}
	return pri, nil
}
func QueryRouteItem(rid int) (*[]route.RouteItem, error) {
	result := []route.RouteItem{}
	err := db.MysqlDB.Where(&route.RouteItem{RouteId: rid}).Find(&result).Error
	return &result, err
}

func QuerySetupTime(gid, aid int) (int64, error) {
	src := user_service.QueryUserById(group_service.QueryGroupById(gid).GroupRiderId)
	if src.UserId < 0 {
		return 0, fmt.Errorf("user not found")
	}

	d, e := amap.DistanceByAid(src.UserDefaultAddressId, aid)
	return int64(d), e
}

// uid: dst user id
// (duration, time, err)
func QueryGroupTime(gid int, uid int) (int64, time.Time, error) {
	var result int64 = 0
	routes, err := QueryRouteByGroupId(gid)
	if err != nil {
		return 0, time.Now(), err
	}
	if len(*routes) == 0 {
		return 0, time.Now(), nil
	}
	var start_at, end_at time.Time
	if !(*routes)[0].RouteDone {
		start_at = time.Now()
		end_at = start_at

		dst := user_service.QueryUserById((*routes)[0].RouteUserId)
		if dst.UserId < 0 {
			return 0, time.Now(), fmt.Errorf("user not found")
		}
		result, err = QuerySetupTime(gid, dst.UserDefaultAddressId)
		if err != nil {
			return 0, time.Now(), err
		}
	} else {
		start_at = *(*routes)[0].RouteFinishedAt
		end_at = start_at
	}

	for _, rt := range *routes {
		if rt.RouteDone {
			end_at = *rt.RouteFinishedAt
		} else {
			result += rt.RouteEstimatedTime
		}
		if rt.RouteUserId == uid {
			break
		}
	}

	result += end_at.Unix() - start_at.Unix()
	end_at = start_at.Add(time.Duration(result * 1e9))

	return result, end_at, nil
}
