package amap

import (
	"ess/service/address_service"
	"ess/utils/amap_base"
)

// seconds
func DistanceByAid(src_aid, dst_aid int) (uint64, error) {
	src, err := address_service.QueryAddressById(src_aid)
	if err != nil {
		return amap_base.Inf, err
	}
	dst, err := address_service.QueryAddressById(dst_aid)
	if err != nil {
		return amap_base.Inf, err
	}
	return amap_base.DistanceByCoordination(src.AddressLng, src.AddressLat, dst.AddressLng, dst.AddressLat)
}
