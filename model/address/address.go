package address

import "gorm.io/gorm"

type Address struct {
	AddressId        int     `gorm:"primaryKey"`
	AddressProvince  string  `gorm:"size:31;not null"`
	AddressCity      string  `gorm:"size:31;not null"`
	AddressArea      string  `gorm:"size:31;not null"`
	AddressDetail    string  `gorm:"size:127;not null"`
	AddressLat       float64 `gorm:"not null"`
	AddressLng       float64 `gorm:"not null"`
	AddressUserId    int     `gorm:"index;not null"` // CAUTIOUS: ensure foreignKey manually
	AddressCached    bool    `gorm:"not null;default:false"`
	AddressCreatedAt int64   `gorm:"autoCreateTime"`
	AddressUpdatedAt int64   `gorm:"autoUpdateTime"`
	AddressDeleted   gorm.DeletedAt
}

type DistanceCache struct {
	DistanceId        int     `gorm:"primaryKey"`
	DistanceCost      uint64  `gorm:"not null"`
	DistanceA         Address `gorm:"foreignKey:DistanceAId"`
	DistanceAId       int     `gorm:"index;not null"`
	DistanceB         Address `gorm:"foreignKey:DistanceBId"`
	DistanceBId       int     `gorm:"index;not null;check:distance_a_id<distance_b_id"`
	DistanceCreatedAt int64   `gorm:"autoCreateTime"`
	DistanceUpdatedAt int64   `gorm:"autoUpdateTime"`
}

type AddressCreateReq struct {
	AddressProvince  string `json:"province" form:"province" binding:"required" example:"浙江省"`
	AddressCity      string `json:"city" form:"city" binding:"required" example:"杭州市"`
	AddressArea      string `json:"area" form:"area" binding:"required" example:"西湖区"`
	AddressDetail    string `json:"detail" form:"detail" binding:"required" example:"浙江大学紫金港校区"`
	IsDefaultAddress bool   `json:"is_default" form:"is_default" binding:"required" example:"true"`
}

type AddressCreateResp struct {
	AddressId int `json:"id"`
}

type AddressDeleteReq struct {
	AddressId int `json:"address_id" form:"address_id" binding:"required"`
}
