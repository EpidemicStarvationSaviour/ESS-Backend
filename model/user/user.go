package user

import (
	"ess/model/address"

	"gorm.io/gorm"
)

type User struct {
	UserId               int             `gorm:"primaryKey"`
	UserName             string          `gorm:"uniqueIndex;size:30"`
	UserEmail            string          `gorm:"uniqueIndex;size:30"` // TODO(TO/GA): delete it
	UserType             Role            `gorm:"not null"`
	UserPhone            string          `gorm:"size:20;not null"`
	UserSecret           string          `gorm:"size:63;not null"`
	UserDefaultAddress   address.Address `gorm:"foreignKey:UserDefaultAddressId"`
	UserDefaultAddressId int             `gorm:"not null"`
	UserAvailable        bool            `gorm:"default:true;not null"`
	UserCreatedAt        int64           `gorm:"autoCreateTime"`
	UserUpdatedAt        int64           `gorm:"autoUpdateTime"`
	UserDeleted          gorm.DeletedAt
}

type UserInfoResp struct {
	ID    int    `json:"userId"`
	Name  string `json:"userName"`
	Email string `json:"userEmail"`
	Type  Role   `json:"userType"`
	Phone string `json:"userPhone"`
}

type UserModifyReq struct {
	UserName  string `json:"userName" binding:"required"`
	UserEmail string `json:"userEmail" binding:"required"`
	UserPhone string `json:"userPhone" binding:"required"`
}

type UserCreateReq struct {
	UserName   string `json:"userName" form:"userName" binding:"required,max=30"`
	UserEmail  string `json:"userEmail" form:"userEmail" binding:"required,max=30"`
	UserSecret string `json:"userSecret" form:"userSecret" binding:"required,max=20"`
	UserPhone  string `json:"userPhone" form:"userPhone" binding:"required,max=20"`
	NoCookie   bool   `json:"noCookie" form:"noCookie"`
}

type UserDeleteReq struct {
	ID int `json:"id" binding:"required"`
}

type UserChangeRoleReq struct {
	ID   int  `json:"userId" binding:"required"`
	Type Role `json:"userType" binding:"required"`
}

type AuthReq struct {
	Type    string `json:"type" example:"email" binding:"required,oneof=account email"`
	Account string `json:"account" example:"" `
	Email   string `json:"email" example:"admin@ess.org" form:"email"`
	Secret  string `json:"secret" example:"essess" form:"secret" binding:"required"`
}

type AuthResp struct {
	UserName  string `json:"userName"`
	UserEmail string `json:"userEmail"`
	UserType  Role   `json:"userType"`
	UserToken string `json:"userToken"`
	LoginType string `json:"loginType"`
}

type TokenAuth struct {
	Token string `json:"token" form:"token" binding:"required"`
}
