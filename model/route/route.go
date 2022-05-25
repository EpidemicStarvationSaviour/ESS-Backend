package route

import (
	"ess/model/category"
	"ess/model/group"
	"ess/model/user"
	"time"
)

type Route struct {
	RouteId            int         `gorm:"primaryKey"`
	RouteGroup         group.Group `gorm:"foreignKey:RouteGroupId"`
	RouteGroupId       int         `gorm:"not null"`
	RouteIndex         uint        `gorm:"not null"`
	RouteUser          user.User   `gorm:"foreignKey:RouteUserId"`
	RouteUserId        int         `gorm:"not null"`
	RouteItems         []RouteItem `gorm:"not null;foreignKey:RouteId"`
	RouteEstimatedTime int64       `gorm:"not null;check:route_estimated_time>=0"` // seconds
	RouteDone          bool        `gorm:"not null;default:0"`
	RouteFinishedAt    time.Time   `gorm:""`
	RouteCreatedAt     int64       `gorm:"autoCreateTime"`
	RouteUpdatedAt     int64       `gorm:"autoUpdateTime"`
}

type RouteItem struct {
	RouteItemId         int               `gorm:"primaryKey"`
	RouteId             int               `gorm:"not null"`
	RouteItemCategory   category.Category `gorm:"foreignKey:RouteItemCategoryId"`
	RouteItemCategoryId int               `gorm:"not null"`
	RouteItemAmount     float64           `gorm:"not null;check:route_item_amount>=0"`
	RouteItemCreatedAt  int64             `gorm:"autoCreateTime"`
	RouteItemUpdatedAt  int64             `gorm:"autoUpdateTime"`
}
