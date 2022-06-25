package rider_service

import (
	"errors"
	"ess/model/address"
	"ess/model/group"
	"ess/model/rider"
	"ess/model/route"
	"ess/model/user"
	"ess/service/address_service"
	"ess/service/group_service"
	"ess/service/route_service"
	"ess/service/user_service"
	"ess/utils/db"
	"time"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func GetRiderAvailable(RiderId int) {
	usr := user_service.QueryUserById(RiderId)
	usr.UserAvailable = true
	err := db.MysqlDB.Model(&usr).Updates(usr).Error
	if err == nil {
		return
	}
}

func GetRiderNotavailable(RiderId int) {
	usr := user_service.QueryUserById(RiderId)
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

func QueryAvailableOrder() (*rider.RiderQueryNewOrdersResp, error) {
	var availableorder rider.RiderQueryNewOrdersResp
	var grou group.Group

	if err := db.MysqlDB.Where(&group.Group{GroupStatus: group.Submitted}).Or(&group.Group{GroupStatus: group.Delivering}).Order("group_updated_at ASC").First(&grou).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var err error
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
		return nil, err
	}
	availableorder.OrderRemark = grou.GroupRemark
	//availableorder.OrderDistance = 0
	_, est_end_time, err := route_service.QueryGroupTime(grou.GroupId, 0)
	if err != nil {
		return nil, err
	}
	availableorder.OrderExpectedTime = int64(-time.Since(est_end_time).Minutes())

	return &availableorder, nil
}

// func RiderFeedbackToOrder(rid int, yesorno int) {
// 	var grou group.Group
// 	db.MysqlDB.Where(&group.Group{GroupId: OrderId}).Find(&grou)
// 	if yesorno == 1 {
// 		grou.GroupSeenByRider = true
// 		grou.GroupStatus = 3
// 		usr := user_service.QueryUserById(rid)
// 		grou.GroupRiderId = usr.UserId
// 	}
// 	if yesorno == 0 {
// 		grou.GroupSeenByRider = true
// 	}
// 	err := db.MysqlDB.Model(grou).Updates(grou).Error
// 	if err != nil {
// 		return
// 	}
// }
func RefreshOrderStatus(uid int, RFTO rider.FeedbackToOrder) error {
	usr := user_service.QueryUserById(RFTO.StoreId)
	var new_status group.Status
	if usr.UserRole == user.Supplier {
		new_status = group.Delivering
	} else {
		new_status = group.Finished
	}
	RefreshRiderPosition(uid, RFTO.AddressLat, RFTO.AddressLng)
	var rout route.Route
	rout, err := route_service.QueryRouteByUserAndGroup(RFTO.StoreId, RFTO.GroupId)
	if err != nil {
		return err
	}
	rout.RouteDone = true
	t := time.Now()
	rout.RouteFinishedAt = &t
	err = db.MysqlDB.Model(rout).Updates(rout).Error
	if err != nil {
		return err
	}
	if err := group_service.UpdateGroup(&group.Group{
		GroupId:     RFTO.GroupId,
		GroupStatus: new_status,
	}); err != nil {
		return err
	}
	return nil
}
