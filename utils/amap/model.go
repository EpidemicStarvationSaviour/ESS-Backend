package amap

type NavigationResp struct {
	Status   string              `json:"status"`
	Info     string              `json:"info"`
	Infocode string              `json:"infocode"`
	Count    string              `json:"count"`
	Route    NavigationRespRoute `json:"route"`
}

type NavigationRespRoute struct {
	Origin      string `json:"origin"`
	Destination string `json:"destination"`
	// TaxiCost    string               `json:"taxi_cost"`
	Paths []NavigationRespPath `json:"paths"`
}

type NavigationRespPath struct {
	Distance string `json:"distance"`
	// Restriction string               `json:"restriction"`
	Cost  NavigationRespCost   `json:"cost"`
	Steps []NavigationRespStep `json:"steps"`
}

type NavigationRespCost struct {
	Duration string `json:"duration"`
	// Tolls         string `json:"tolls"`
	// TollDistance  string `json:"toll_distance"`
	// TollRoad      string `json:"toll_road"`
	// TrafficLights string `json:"traffic_lights"`
}

type NavigationRespStep struct {
	// Instruction  string             `json:"instruction"`
	// Orientation  string             `json:"orientation"`
	// RoadName     string             `json:"road_name"`
	StepDistance string             `json:"step_distance"`
	Cost         NavigationRespCost `json:"cost"`
}

type CoordinationResp struct {
	Status   string                    `json:"status"`
	Info     string                    `json:"info"`
	Infocode string                    `json:"infocode"`
	Count    string                    `json:"count"`
	Geocodes []CoordinationRespGeocode `json:"geocodes"`
}

type CoordinationRespGeocode struct {
	FormattedAddress string `json:"formatted_address"`
	Country          string `json:"country"`
	AddressProvince  string `json:"province"`
	AddressCity      string `json:"city"`
	AddressArea      string `json:"district"`
	Adcode           string `json:"adcode"`
	Location         string `json:"location"`
}
