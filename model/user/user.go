package user

import (
	"ess/model/address"

	"gorm.io/gorm"
)

type User struct {
	UserId               int             `gorm:"primaryKey"`
	UserName             string          `gorm:"uniqueIndex;size:30"`
	UserRole             Role            `gorm:"not null"`
	UserPhone            string          `gorm:"uniqueIndex;size:20;not null"`
	UserSecret           string          `gorm:"size:63;not null"`
	UserDefaultAddress   address.Address `gorm:"foreignKey:UserDefaultAddressId"`
	UserDefaultAddressId int             `gorm:"not null"`
	UserAvailable        bool            `gorm:"default:true;not null"`
	UserCreatedAt        int64           `gorm:"autoCreateTime"`
	UserUpdatedAt        int64           `gorm:"autoUpdateTime"`
	UserDeleted          gorm.DeletedAt
}

type UserCreateReq struct {
	UserPhone   string               `json:"user_phone" form:"user_phone" binding:"required,max=20" example:"13800138000"`
	UserSecret  string               `json:"userSecret" form:"userSecret" binding:"required,max=20"`
	UserRole    Role                 `json:"user_role" form:"user_role" binding:"required" example:"1"`
	UserName    string               `json:"user_secret" form:"user_secret" binding:"required,max=30"`
	UserAddress UserCreateReqAddress `json:"user_address" form:"user_address" binding:"required"`
}

type UserCreateReqAddress struct {
	AddressProvince string `json:"province" form:"province" binding:"required"`
	AddressCity     string `json:"city" form:"city" binding:"required"`
	AddressArea     string `json:"area" form:"area" binding:"required"`
	AddressDetail   string `json:"detail" form:"detail" binding:"required"`
}

type UserCreateResp struct {
	UserId int `json:"id" form:"id"`
}

type UserInfoResp struct {
	ID    int    `json:"userId"`
	Name  string `json:"userName"`
	Type  Role   `json:"userType"`
	Phone string `json:"userPhone"`
}

type UserModifyReq struct {
	UserName  string `json:"userName" binding:"required"`
	UserPhone string `json:"userPhone" binding:"required"`
}

type UserDeleteReq struct {
	ID int `json:"id" binding:"required"`
}

type UserChangeRoleReq struct {
	ID   int  `json:"userId" binding:"required"`
	Type Role `json:"userType" binding:"required"`
}

type AuthReq struct {
	Type    string `json:"type" form:"type" example:"name" binding:"required,oneof=name phone"`
	Account string `json:"account" form:"account" binding:"required"`
	Secret  string `json:"password" form:"password" binding:"required"`
}

type AuthResp struct {
	UserPhone string `json:"user_phone" form:"user_phone"`
	UserName  string `json:"user_name" form:"user_name"`
	UserRole  Role   `json:"user_role" form:"user_role"`
	UserToken string `json:"user_token" form:"user_token"`
	LoginType string `json:"login_type" form:"login_type"`
}

type TokenAuth struct {
	Token string `json:"token" form:"token" binding:"required"`
}
