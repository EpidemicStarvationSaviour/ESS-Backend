package algorithm

import (
	"ess/model/item"
	"ess/model/route"
	"ess/model/user"
	"ess/service/address_service"
	"ess/service/group_service"
	"ess/service/item_service"
	"ess/service/user_service"
	"ess/utils/amap"
	"ess/utils/db"
	"ess/utils/logging"
	"ess/utils/setting"
	"fmt"
	"log"
	"time"

	pb "ess/gRPC"

	"github.com/vishalkuo/bimap"
	"gorm.io/gorm"
)

var s server
var enable bool
var timeout time.Duration

func Setup() {
	enable = setting.GRPCSetting.Enable
	if enable {
		timeout = time.Millisecond * time.Duration(setting.GRPCSetting.Timeout)
		s.Setup()

		_, err := s.Ping(&pb.PingRequest{Message: "ping"})
		if err != nil {
			log.Fatalf("failed to ping rpc server: %v", err)
		}
	}
}

func Schedule(gid, uid int) error {
	if !enable {
		return fmt.Errorf("rpc server is not enabled")
	}
	var cnt uint32
	cnt = 1
	category_map := bimap.NewBiMap()
	request_items := make(map[uint32]float64)
	grp := group_service.QueryGroupById(gid)
	grp_info, err := group_service.GetGroupDetail(grp, uid)
	if err != nil {
		return err
	}
	for _, c := range grp_info.Commodities {
		category_map.Insert(cnt, c.CategoryId)
		request_items[cnt] = c.TotalAmount
		cnt++
	}

	cnt = 1
	rider_map := make(map[uint32]int)
	riders, err := user_service.QueryAvailableRiders()
	if err != nil {
		return err
	}
	for _, r := range riders {
		rider_map[cnt] = r.UserId
		cnt++
	}

	cnt = 1
	supplier_map := make(map[uint32]int)
	var supplier_items []*pb.ItemList
	var supplier_items_tmp map[uint32]float64
	suppliers, err := user_service.QueryUsersByRole(user.Supplier)
	if err != nil {
		return err
	}
	for _, s := range suppliers {
		supplier_items_tmp = make(map[uint32]float64)
		supplier_map[cnt] = s.UserId
		items, err := item_service.QueryItemsByUserId(s.UserId)
		if err != nil {
			return err
		}
		for _, item := range items {
			val, ok := category_map.GetInverse(item.ItemCategoryId)
			if ok {
				supplier_items_tmp[val.(uint32)] = item.ItemAmount
			}
		}
		supplier_items = append(supplier_items, &pb.ItemList{
			Items: supplier_items_tmp,
		})
		cnt++
	}

	var distances []uint64
	// (0, 1) (0, 2)...(0, n)
	for _, s := range suppliers {
		dis, err := address_service.QueryDistanceCacheByAid(grp.GroupAddressId, s.UserDefaultAddressId)
		if err != nil {
			return err
		}
		distances = append(distances, dis)
	}
	// (1, 2) (1, 3)...(1, m) (2, 3) (2, 4)...(2, m) ... (n, m)
	for i, s := range suppliers {
		for j := i + 1; j < len(suppliers); j++ {
			dis, err := address_service.QueryDistanceCacheByAid(s.UserDefaultAddressId, suppliers[j].UserDefaultAddressId)
			if err != nil {
				return err
			}
			distances = append(distances, dis)
		}
		for _, r := range riders {
			dis, err := amap.DistanceByAid(s.UserDefaultAddressId, r.UserDefaultAddressId)
			if err != nil {
				return err
			}
			distances = append(distances, dis)
		}
	}

	// log.Printf("request: %+v", request_items)
	// log.Printf("riders: %+v", len(riders))
	// log.Printf("items: %+v", supplier_items)
	// log.Printf("distances: %+v", distances)
	// return fmt.Errorf("debug")

	// gRPC call
	resp, err := s.Schedule(&pb.ScheduleRequest{
		Request: &pb.ItemList{
			Items: request_items,
		},
		NumDeliverer: uint32(len(riders)),
		Itemlists:    supplier_items,
		Distance:     distances,
	})
	if err != nil {
		logging.ErrorF("could not schedule: %+v", err)
	}
	// log.Printf("%+v", resp)
	// return fmt.Errorf("debug")

	// database transaction
	tx := db.MysqlDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return err
	}

	// update group status
	grp.GroupRiderId = rider_map[resp.DelivererId]
	if err := tx.Model(&grp).Updates(grp).Error; err != nil {
		tx.Rollback()
		return err
	}

	last_supplier_mapped_id := uint32(0)
	for i, r := range resp.Route {
		rt := route.Route{
			RouteGroupId: gid,
			RouteIndex:   uint(i + 1),
			RouteUserId:  supplier_map[r.SupplierId],
		}
		if last_supplier_mapped_id != 0 {
			dist, err := amap.DistanceByAid(suppliers[last_supplier_mapped_id-1].UserDefaultAddressId, suppliers[r.SupplierId-1].UserDefaultAddressId)
			rt.RouteEstimatedTime = int64(dist)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			dist, err := amap.DistanceByAid(riders[resp.DelivererId-1].UserDefaultAddressId, suppliers[r.SupplierId-1].UserDefaultAddressId)
			rt.RouteEstimatedTime = int64(dist)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
		last_supplier_mapped_id = r.SupplierId
		if err := tx.Create(&rt).Error; err != nil {
			tx.Rollback()
			return err
		}

		for id, amount := range r.Itemlist.Items {
			if amount != 0 {
				cid, _ := category_map.Get(id)
				rt.RouteItems = append(rt.RouteItems, route.RouteItem{
					RouteId:             rt.RouteId,
					RouteItemCategoryId: cid.(int),
					RouteItemAmount:     amount,
				})
				if err := tx.Model(&item.Item{}).Where(&item.Item{ItemUserId: rt.RouteUserId, ItemCategoryId: cid.(int)}).UpdateColumn("item_amount", gorm.Expr("item_amount - ?", amount)).Error; err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		if err := tx.Model(&rt).Association("RouteItems").Replace(rt.RouteItems); err != nil {
			tx.Rollback()
			return err
		}
	}

	// the final supplier to user
	time, err := address_service.QueryDistanceCacheByAid(suppliers[last_supplier_mapped_id-1].UserDefaultAddressId, grp.GroupAddressId)
	if err != nil {
		tx.Rollback()
		return err
	}
	rt := route.Route{
		RouteGroupId:       gid,
		RouteIndex:         uint(len(resp.Route) + 1),
		RouteUserId:        grp.GroupCreatorId,
		RouteEstimatedTime: int64(time),
	}
	if err := tx.Create(&rt).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
