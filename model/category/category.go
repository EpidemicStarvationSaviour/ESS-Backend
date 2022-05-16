package category

import "gorm.io/gorm"

type Category struct {
	CategoryId        int       `gorm:"primaryKey"`
	CategoryName      string    `gorm:"size:31;not null"`
	CategoryPrice     float64   `gorm:"not null;check:category_price>=0"`
	CategoryImageUrl  string    `gorm:"size:127"`
	CategoryLevel     uint      `gorm:"not null"` // 0: root, 1: sub
	CategoryFather    *Category `gorm:"foreignKey:CategoryFatherId"`
	CategoryFatherId  int       `gorm:""`
	CategoryCreatedAt int64     `gorm:"autoCreateTime"`
	CategoryUpdatedAt int64     `gorm:"autoUpdateTime"`
	CategoryDeleted   gorm.DeletedAt
}
