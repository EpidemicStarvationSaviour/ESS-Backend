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

type GroupInfoResp struct {
	Count int `json:"count"`
}
