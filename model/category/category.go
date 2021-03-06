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

type CategoryAllResp struct {
	CategoryList []CategoryInfoResp `json:"data"`
}
type CategoryInfoResp struct {
	CategoryLevel    int                `json:"type_id"`
	CategoryName     string             `json:"type_name"`
	CategoryNumber   int                `json:"type_number"`
	CategoryImageUrl string             `json:"type_avatar"`
	Categorychild    []CategoryChildren `json:"children"`
}

type CategoryChildren struct {
	CategoryId     int     `json:"id"`
	CategoryName   string  `json:"name"`
	CategoryAvatar string  `json:"avatar"`
	CategoryTotal  float64 `json:"total"`
	CategoryPrice  float64 `json:"price"`
}

type CategoryCreateReq struct {
	CategoryFatherId int     `json:"type_id"`
	CategoryName     string  `json:"name"`
	CategoryPrice    float64 `json:"price"`
	CategoryAvatar   string  `json:"avatar"`
}

type CategoryCreateResp struct {
	CategoryId int `json:"id"`
}

type CategoryDeleted struct {
	CategoryId int `json:"id"`
}

type CategoryCertainInfoRep struct {
	CategoryId int `uri:"id" binding:"required"`
}

type CategoryCertainInfoResp struct {
	CategoryFatherId int                   `json:"type_id"`
	CategoryId       int                   `json:"id"`
	CategoryAvatar   string                `json:"avatar"`
	CategoryName     string                `json:"name"`
	CategoryPrice    float64               `json:"price"`
	CategoryTotal    float64               `json:"total"`
	CategoryDetails  []CategoryDetailsInfo `json:"details"`
}

type CategoryDetailsInfo struct {
	StoreId            int     `json:"store_id"`
	CategoryLat        float64 `json:"lat"`
	CategoryLng        float64 `json:"lng"`
	CategoryAddress    string  `json:"address"`
	CategoryStorephone string  `json:"phone"`
	Categorynumber     float64 `json:"number"`
}

type CategoryAddress struct {
	AddressProvince string `json:"province" form:"province" binding:"required" example:"?????????"`
	AddressCity     string `json:"city" form:"city" binding:"required" example:"?????????"`
	AddressArea     string `json:"area" form:"area" binding:"required" example:"?????????"`
	AddressDetail   string `json:"detail" form:"detail" binding:"required" example:"???????????????????????????"`
}

type CategoryMyInfoResp struct {
	CategoryLevel  int                  `json:"type_id"`
	CategoryName   string               `json:"type_name"`
	CategoryNumber int                  `json:"type_number"`
	CategoryAvatar string               `json:"type_avatar"`
	Categorychild  []CategoryMyChildren `json:"children"`
}

type CategoryMyChildren struct {
	CategoryId     int     `json:"id"`
	CategoryName   string  `json:"name"`
	CategoryAvatar string  `json:"avatar"`
	CategoryTotal  float64 `json:"total"`
	CategoryPrice  float64 `json:"price"`
}

type CategoryModifyRep struct {
	CategoryId     int     `json:"id"`
	Categorynumber float64 `json:"number"`
}
