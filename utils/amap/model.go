package amap

type NavigationResp struct {
	Status   string              `json:"status"`
	Info     string              `json:"info"`
	Infocode string              `json:"infocode"`
	Count    string              `json:"count"`
	Route    NavigationRespRoute `json:"route"`
}

type NavigationRespRoute struct {
	Origin      string               `json:"origin"`
	Destination string               `json:"destination"`
	TaxiCost    string               `json:"taxi_cost"`
	Paths       []NavigationRespPath `json:"paths"`
}

type NavigationRespPath struct {
	Distance    string               `json:"distance"`
	Restriction string               `json:"restriction"`
	Cost        NavigationRespCost   `json:"cost"`
	Steps       []NavigationRespStep `json:"steps"`
}

type NavigationRespCost struct {
	Duration      string `json:"duration"`
	Tolls         string `json:"tolls"`
	TollDistance  string `json:"toll_distance"`
	TollRoad      string `json:"toll_road"`
	TrafficLights string `json:"traffic_lights"`
}

type NavigationRespStep struct {
	Instruction  string             `json:"instruction"`
	Orientation  string             `json:"orientation"`
	RoadName     string             `json:"road_name"`
	StepDistance string             `json:"step_distance"`
	Cost         NavigationRespCost `json:"cost"`
}
