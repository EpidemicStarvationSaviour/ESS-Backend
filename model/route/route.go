package route

import (
	"ess/model/group"
	"ess/model/user"
)

type Route struct {
	RouteId        int         `gorm:"primaryKey"`
	RouteGroup     group.Group `gorm:"foreignKey:RouteGroupId"`
	RouteGroupId   int         `gorm:"not null"`
	RouteIndex     uint        `gorm:"not null"`
	RouteUser      user.User   `gorm:"foreignKey:RouteUserId"`
	RouteUserId    int         `gorm:"not null"`
	RouteItems     string      `gorm:"size:255"`
	RouteTime      float64     `gorm:"not null;check:route_time>=0"`
	RouteDist      float64     `gorm:"not null;check:route_dist>=0"`
	RouteDone      bool        `gorm:"not null;default:false"`
	RouteCreatedAt int64       `gorm:"autoCreateTime"`
	RouteUpdatedAt int64       `gorm:"autoUpdateTime"`
}
