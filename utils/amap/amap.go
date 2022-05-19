package amap

import (
	"encoding/json"
	"ess/model/address"
	"ess/service/address_service"
	"ess/utils/logging"
	"ess/utils/setting"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/jinzhu/copier"
)

var client *http.Client
var inf uint64

func Setup() {
	if setting.AmapSetting.Enable {
		client = &http.Client{}
		inf = math.MaxUint32
	}
}

// seconds
func DistanceByAid(src_aid, dst_aid int) (uint64, error) {
	src, err := address_service.QueryAddressById(src_aid)
	if err != nil {
		return inf, err
	}
	dst, err := address_service.QueryAddressById(dst_aid)
	if err != nil {
		return inf, err
	}
	return DistanceByCoordination(src.AddressLng, src.AddressLat, dst.AddressLng, dst.AddressLat)
}

// seconds
func DistanceByCoordination(src_lng, src_lat, dst_lng, dst_lat float64) (uint64, error) {
	if !setting.AmapSetting.Enable {
		return uint64(rand.Int63n(3600)), nil
	}
	req, err := http.NewRequest("POST", "https://restapi.amap.com/v5/direction/driving", nil)
	if err != nil {
		return inf, err
	}
	query := req.URL.Query()
	query.Add("key", setting.AmapSetting.WebAPIKey)
	query.Add("origin", fmt.Sprintf("%.6f,%.6f", src_lng, src_lat))
	query.Add("destination", fmt.Sprintf("%.6f,%.6f", dst_lng, dst_lat))
	query.Add("show_fields", "cost")
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return inf, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return inf, err
	}

	var result NavigationResp
	if err := json.Unmarshal(body, &result); err != nil {
		return inf, err
	}

	if result.Status != "1" || result.Infocode != "10000" {
		logging.ErrorF("[amap] request error: %+v", result)
		return inf, fmt.Errorf("高德 Web 服务路线规划 2.0 API 调用失败")
	}

	ret, err := strconv.ParseUint(result.Route.Paths[0].Cost.Duration, 10, 64)
	if err != nil {
		return inf, err
	}

	return ret, nil
}

func GetCoordination(addr *address.Address) error {
	if !setting.AmapSetting.Enable {
		addr.AddressLng = 30.263842 + (rand.Float64()-0.5)/0.5*0.2  // +- 0.2
		addr.AddressLat = 120.123077 + (rand.Float64()-0.5)/0.5*0.2 // +- 0.2
	}
	req, err := http.NewRequest("GET", "https://restapi.amap.com/v3/geocode/geo", nil)
	if err != nil {
		return err
	}
	query := req.URL.Query()
	query.Add("key", setting.AmapSetting.WebAPIKey)
	query.Add("city", addr.AddressCity)
	query.Add("address", fmt.Sprintf("%s%s%s%s", addr.AddressProvince, addr.AddressCity, addr.AddressArea, addr.AddressDetail))
	req.URL.RawQuery = query.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result CoordinationResp
	if err := json.Unmarshal(body, &result); err != nil {
		return err
	}

	if result.Status != "1" || result.Infocode != "10000" {
		logging.ErrorF("[amap] request error: %+v", result)
		return fmt.Errorf("高德 Web 地理编码 API 调用失败")
	}

	_ = copier.Copy(addr, result.Geocodes[0])
	_, err = fmt.Sscanf(result.Geocodes[0].Location, "%f,%f", &addr.AddressLng, &addr.AddressLat)
	if err != nil {
		return err
	}

	return nil
}
