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
	AddressCreatedAt int64   `gorm:"autoCreateTime"`
	AddressUpdatedAt int64   `gorm:"autoUpdateTime"`
	AddressDeleted   gorm.DeletedAt
}
