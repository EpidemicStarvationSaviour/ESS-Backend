package rider_service

import (
	"ess/model/address"
	"ess/model/group"
	"ess/model/rider"
	"ess/model/route"
	"ess/model/user"
	"ess/service/address_service"
	"ess/service/route_service"
	"ess/service/user_service"
	"ess/utils/db"
	"time"

	"github.com/jinzhu/copier"
)

func GetRiderAvailable(RiderId int) {
	var usr user.User
	usr = user_service.QueryUserById(RiderId)
	usr.UserAvailable = true
	err := db.MysqlDB.Model(&usr).Updates(usr).Error
	if err == nil {
		return
	}
}

func GetRiderNotavailable(RiderId int) {
	var usr user.User
	usr = user_service.QueryUserById(RiderId)
	usr.UserAvailable = false
	err := db.MysqlDB.Model(&usr).Update("user_available", 0).Error
	if err == nil {
		return
	}
}

func RefreshRiderPosition(RiderId int, lat float64, lng float64) {
	var addr address.Address
	addr, _ = address_service.QueryDefaultAddressByUserId(RiderId)
	addr.AddressLat = lat
	addr.AddressLng = lng
	db.MysqlDB.Where(&address.Address{AddressId: addr.AddressId}).Updates(addr)
}

var OrderId int

func QueryAvailableOrder() (error, *rider.RiderQueryNewOrdersResp) {
	var availableorder rider.RiderQueryNewOrdersResp
	var grou group.Group

	if err := db.MysqlDB.Where(&group.Group{GroupStatus: 2, GroupSeenByRider: false}).First(&grou).Error; err != nil { //数据库的返回值
		return err, nil
	}
	grou.GroupSeenByRider = true
	err := db.MysqlDB.Model(grou).Updates(grou).Error
	if err != nil {
		return err, nil
	}
	var usr user.User
	var addr *address.Address
	OrderId = grou.GroupId
	availableorder.OrderId = grou.GroupId
	availableorder.GroupName = grou.GroupName
	usr = user_service.QueryUserById(grou.GroupCreatorId)

	availableorder.CreatorName = usr.UserName
	availableorder.CreatorPhone = usr.UserPhone
	addr, _ = address_service.QueryAddressById(usr.UserDefaultAddressId)
	_ = copier.Copy(&availableorder.CreatorAddress, &addr)
	availableorder.OrderReward, err = route_service.GetRouteItemTotalPrice(grou.GroupId)
	if err != nil {
		return err, nil
	}
	availableorder.OrderRemark = grou.GroupRemark
	//availableorder.OrderDistance = 0
	_, est_end_time, err := route_service.QueryGroupTime(grou.GroupId, 0)
	if err != nil {
		return err, nil
	}
	availableorder.OrderExpectedTime = int64(time.Since(est_end_time).Minutes())

	return nil, &availableorder
}

func RiderFeedbackToOrder(rid int, yesorno int) {
	var grou group.Group
	db.MysqlDB.Where(&group.Group{GroupId: OrderId}).Find(&grou)
	if yesorno == 1 {
		grou.GroupSeenByRider = true
		grou.GroupStatus = 3
		usr := user_service.QueryUserById(rid)
		grou.GroupRiderId = usr.UserId
	}
	if yesorno == 0 {
		grou.GroupSeenByRider = true
	}
	err := db.MysqlDB.Model(grou).Updates(grou).Error
	if err != nil {
		return
	}
}
func RefreshOrderStatus(uid int, RFTO rider.FeedbackToOrder) {
	usr := user_service.QueryUserById(uid)
	if usr.UserRole == 2 {
		RefreshRiderPosition(uid, RFTO.AddressLat, RFTO.AddressLng)
		var rout route.Route
		rout, _ = route_service.QueryRouteByUserAndGroup(RFTO.StoreId, RFTO.GroupId)
		rout.RouteDone = true
		t := time.Now()
		rout.RouteFinishedAt = &t
		db.MysqlDB.Model(rout).Updates(rout)
	}
	if usr.UserRole == 3 || usr.UserRole == 4 {
		var rout route.Route
		rout, _ = route_service.QueryRouteByUserAndGroup(RFTO.StoreId, RFTO.GroupId)
		rout.RouteDone = true
		t := time.Now()
		rout.RouteFinishedAt = &t
		db.MysqlDB.Model(rout).Updates(rout)
		var grou group.Group
		db.MysqlDB.Where(&group.Group{GroupId: OrderId}).Find(&grou)
		grou.GroupStatus = 4
		db.MysqlDB.Model(grou).Updates(grou)
	}

}
