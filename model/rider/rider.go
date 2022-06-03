package rider

type RiderStart struct {
	AddressLat float64 `json:"lat"`
	AddressLng float64 `json:"lng"`
}

type RiderStop struct {
}

type RiderUploadAddress struct {
	AddressLat float64 `json:"lat"`
	AddressLng float64 `json:"lng"`
}

type RiderQueryNewOrdersResp struct {
	OrderId        int            `json:"id"`
	GroupName      string         `json:"name"`
	CreatorName    string         `json:"creator_name"`
	CreatorPhone   string         `json:"creator_phone"`
	CreatorAddress CreatorAddress `json:"creator_address"`
	OrderReward    float64        `json:"reward"`
	OrderRemark    string         `json:"remark"`
	//OrderDistance     float64        `json:"distance"`
	OrderExpectedTime int64 `json:"expected_time"`
}
type CreatorAddress struct {
	AddressId       int     `gorm:"primaryKey"`
	AddressLat      float64 `gorm:"not null"`
	AddressLng      float64 `gorm:"not null"`
	AddressProvince string  `gorm:"size:31;not null"`
	AddressCity     string  `gorm:"size:31;not null"`
	AddressArea     string  `gorm:"size:31;not null"`
	AddressDetail   string  `gorm:"size:127;not null"`
}
type OrderAddress struct {
	AddressId       int     `json:"id"`
	AddressLat      float64 `json:"lat"`
	AddressLng      float64 `json:"lng"`
	AddressProvince string  `json:"province"`
	AddressCity     string  `json:"city"`
	AddressArea     string  `json:"area"`
	AddressDetail   string  `json:"detail"`
}

type RiderFeedbackToNewOrder struct {
	YesOrNo int `json:"type"`
}

type FeedbackToOrder struct {
	AddressLat float64 `json:"lat"`
	AddressLng float64 `json:"lng"`
	StoreId    int     `json:"store_id"`
	GroupId    int     `json:"group_id"`
}
