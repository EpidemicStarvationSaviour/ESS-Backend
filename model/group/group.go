package group

import (
	"ess/model/address"
	"ess/model/category"
	"ess/model/user"
)

type Group struct {
	GroupId          int                 `gorm:"primaryKey"`
	GroupName        string              `gorm:"size:63;not null"`
	GroupDescription string              `gorm:"size:255;not null"`
	GroupRemark      string              `gorm:"size:255;not null"`
	GroupCreator     user.User           `gorm:"foreignKey:GroupCreatorId"`
	GroupCreatorId   int                 `gorm:"not null"`
	GroupAddress     address.Address     `gorm:"foreignKey:GroupAddressId"`
	GroupAddressId   int                 `gorm:"not null"`
	GroupCategories  []category.Category `gorm:"many2many:group_category;"`
	GroupRider       user.User           `gorm:"foreignKey:GroupRiderId"`
	GroupRiderId     int                 `gorm:""`
	GroupStatus      Status              `gorm:"not null;default:1"`
	GroupSeenByRider bool                `gorm:"not null;default:false"`
	GroupCreatedAt   int64               `gorm:"autoCreateTime"`
	GroupUpdatedAt   int64               `gorm:"autoUpdateTime"`
}

type GroupInfoReq struct {
	Type     int `json:"type"`
	PageNum  int `json:"page_num"`
	PageSize int `json:"page_size"`
}

type GroupInfoAddress struct {
	AddressId       int     `json:"id"`
	AddressLat      float64 `json:"lat"`
	AddressLng      float64 `json:"lng"`
	AddressProvince string  `json:"province"`
	AddressCity     string  `json:"city"`
	AddressArea     string  `json:"area"`
	AddressDetail   string  `json:"detail"`
}

type GroupInfoCommodity struct {
	CategoryId       int     `json:"type_id"`
	OrderId          int     `json:"id"`
	CategoryName     string  `json:"name"`
	CategoryImageUrl string  `json:"avatar"`
	CategoryPrice    float64 `json:"price"`
	OrderAmount      float64 `json:"number"`
	TotalAmount      float64 `json:"total_number"`
}

type GroupInfoData struct {
	GroupId          int                  `json:"id"`
	GroupName        string               `json:"name"`
	GroupStatus      Status               `json:"type"`
	UserId           int                  `json:"creator_id"`
	UserName         string               `json:"creator_name"`
	UserPhone        string               `json:"creator_phone"`
	CreatorAddr      GroupInfoAddress     `json:"creator_address"`
	UserNumber       int                  `json:"user_number"`
	TotalPrice       float64              `json:"total_price"`
	TotalMyPrice     float64              `json:"total_my_price"`
	Commodities      []GroupInfoCommodity `json:"commodity_detail"`
	GroupDescription string               `json:"description"`
	GroupRemark      string               `json:"remark"`
	GroupCreatedAt   int64                `json:"created_time"`
}

type GroupInfoResp struct {
	Count int             `json:"count"`
	Data  []GroupInfoData `json:"data"`
}
