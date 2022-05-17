package amap

import "ess/service/address_service"

func DistanceByAid(src_aid, dst_aid int) (float64, error) {
	src, err := address_service.QueryAddressById(src_aid)
	if err != nil {
		return 0, err
	}
	dst, err := address_service.QueryAddressById(dst_aid)
	if err != nil {
		return 0, err
	}
	return DistanceByPos(src.AddressLng, src.AddressLat, dst.AddressLng, dst.AddressLat)
}

func DistanceByPos(src_lng, src_lat, dst_lng, dst_lat float64) (float64, error) {
	// TODO: (TO/GA)
	return 12.5, nil
}
