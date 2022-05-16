package item

import (
	"ess/model/category"
	"ess/model/user"
)

type Item struct {
	ItemId         int               `gorm:"primaryKey"`
	ItemCategory   category.Category `gorm:"foreignKey:ItemCategoryId"`
	ItemCategoryId int               `gorm:"not null"`
	ItemUser       user.User         `gorm:"foreignKey:ItemUserId"`
	ItemUserId     int               `gorm:"not null"`
	ItemAmount     float64           `gorm:"not null;check:item_amount>=0"`
	ItemCreatedAt  int64             `gorm:"autoCreateTime"`
	ItemUpdatedAt  int64             `gorm:"autoUpdateTime"`
}
