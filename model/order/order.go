package order

import (
	"ess/model/category"
	"ess/model/group"
	"ess/model/user"
)

type Order struct {
	OrderId         int               `gorm:"primaryKey"`
	OrderUser       user.User         `gorm:"foreignKey:OrderUserId"`
	OrderUserId     int               `gorm:"not null"`
	OrderGroup      group.Group       `gorm:"foreignKey:OrderGroupId"`
	OrderGroupId    int               `gorm:"not null"`
	OrderCategory   category.Category `gorm:"foreignKey:OrderCategoryId"`
	OrderCategoryId int               `gorm:"not null"`
	OrderAmount     float64           `gorm:"not null;check:order_amount>=0"`
	OrderCreatedAt  int64             `gorm:"autoCreateTime"`
	OrderUpdatedAt  int64             `gorm:"autoUpdateTime"`
}
