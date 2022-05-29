package group

import (
	"ess/model/address"
	"ess/model/category"
	"ess/model/user"

	"gorm.io/gorm"
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
	GroupDeleted     gorm.DeletedAt
}

type GroupInfoReq struct {
	Type     int `form:"type" json:"type"`
	PageNum  int `form:"page_num" json:"page_num"`
	PageSize int `form:"page_size" json:"page_size"`
}

type GroupInfoAddress struct {
	AddressId       int     `json:"id"`
	AddressLat      float64 `json:"lat"`
	AddressLng      float64 `json:"lng"`
	AddressProvince string  `json:"province"`
	AddressCity     string  `json:"city"`
	AddressArea     string  `json:"area"`
	AddressDetail   string  `json:"detail"`
}

type GroupInfoCommodity struct {
	CategoryId       int     `json:"type_id"`
	ParentId         int     `json:"id"`
	CategoryName     string  `json:"name"`
	CategoryImageUrl string  `json:"avatar"`
	CategoryPrice    float64 `json:"price"`
	OrderAmount      float64 `json:"number"`
	TotalAmount      float64 `json:"total_number"`
}

type GroupInfoData struct {
	GroupId          int                  `json:"id"`
	GroupName        string               `json:"name"`
	GroupStatus      Status               `json:"type"`
	UserId           int                  `json:"creator_id"`
	UserName         string               `json:"creator_name"`
	UserPhone        string               `json:"creator_phone"`
	CreatorAddr      GroupInfoAddress     `json:"creator_address"`
	UserNumber       int                  `json:"user_number"`
	TotalPrice       float64              `json:"total_price"`
	TotalMyPrice     float64              `json:"total_my_price"`
	Commodities      []GroupInfoCommodity `json:"commodity_detail"`
	GroupDescription string               `json:"description"`
	GroupRemark      string               `json:"remark"`
	GroupCreatedAt   int64                `json:"created_time"`
}

type GroupInfoResp struct {
	Count int             `json:"count"`
	Data  []GroupInfoData `json:"data"`
}

type GroupCreateReq struct {
	GroupName        string `json:"name"`
	GroupDescription string `json:"description"`
	GroupRemark      string `json:"remark"`
	GroupAddressId   int    `json:"address_id"`
	GroupUserGroupId int    `json:"user_group_id"`
	GroupCommodities []int  `json:"commodities"`
}

type GroupCreateResp struct {
	GroupId int `json:"id"`
}

type GroupSearchReq struct {
	PageSIze    int    `form:"page_size" json:"page_size"`
	PageNum     int    `form:"page_num" json:"page_num"`
	SearchType  int    `form:"search_type" json:"search_type"`
	GroupType   int    `form:"group_type" json:"group_type"`
	SearchValue string `form:"value" json:"value"`
}

type GroupJoinData struct {
	OrderCategoryId int     `json:"commodity_id"`
	OrderAmount     float64 `json:"number"`
}

type GroupJoinReq struct {
	GroupId   int             `json:"id"`
	OrderData []GroupJoinData `json:"data"`
}

type GroupEditReq struct {
	GroupName         string `json:"name"`
	GroupStatus       Status `json:"type"`
	GroupDescription  string `json:"description"`
	GroupRemark       string `json:"remark"`
	GroupAddressId    int    `json:"address_id"`
	GroupDeteledUsers []int  `json:"deleted_users"`
	GroupCommodityIds []int  `json:"commodities"`
}

type GroupInfoSupplierCommodity struct {
	CategoryId       int     `json:"type_id"`
	ParentId         int     `json:"id"`
	CategoryName     string  `json:"name"`
	CategoryImageUrl string  `json:"avatar"`
	CategoryPrice    float64 `json:"price"`
	RouteAmount      float64 `json:"number"`
}

type GroupRiderAddress struct {
	AddressLat         float64 `json:"lat"`
	AddressLng         float64 `json:"lng"`
	AddressUpdatedAt   int64   `json:"update_time"`
	RouteEstimatedTime int64   `json:"eta"`
}

type GroupInfoSupplierData struct {
	GroupId           int                          `json:"id"`
	GroupName         string                       `json:"name"`
	GroupStatus       Status                       `json:"type"`
	GroupCreatorId    int                          `json:"creator_id"`
	GroupCreatorName  string                       `json:"creator_name"`
	GroupCreatorPhone string                       `json:"creator_phone"`
	GroupTotalPrice   float64                      `json:"total_price"`
	GroupRemark       string                       `json:"remark"`
	GroupCommodity    []GroupInfoSupplierCommodity `json:"commodity_detail"`
	GroupRiderPhone   string                       `json:"rider_phone"`
	GroupRiderName    string                       `json:"rider_name"`
	GroupRiderPos     GroupRiderAddress            `json:"rider_pos"`
}

type GroupInfoSupplierResp struct {
	Count     int                     `json:"count"`
	GroupData []GroupInfoSupplierData `json:"data"`
}

type GroupCreatorAddress struct {
	AddressId       int     `json:"id"`
	AddressProvince string  `json:"province"`
	AddressCity     string  `json:"city"`
	AddressArea     string  `json:"area"`
	AddressDetail   string  `json:"detail"`
	AddressLat      float64 `json:"lat"`
	AddressLng      float64 `json:"lng"`
}

type GroupInfoRiderData struct {
	GroupId             int                 `json:"id"`
	GroupName           string              `json:"name"`
	GroupStatus         Status              `json:"type"`
	GroupCreatorId      int                 `json:"creator_id"`
	GroupCreatorName    string              `json:"creator_name"`
	GroupCreatorPhone   string              `json:"creator_phone"`
	GroupCreatorAddress GroupCreatorAddress `json:"creator_address"`
	GroupRemark         string              `json:"remark"`
	GroupReward         float64             `json:"reward"`
}

type GroupInfoRiderResp struct {
	Count     int                  `json:"count"`
	GroupData []GroupInfoRiderData `json:"data"`
}

type GroupAgentCommodityUser struct {
	UserId     int     `json:"user_id"`
	UserName   string  `json:"user_name"`
	UserPhone  string  `json:"user_phone"`
	UserAmount float64 `json:"number"`
}

type GroupAgentCommodity struct {
	CategoryId       int                       `json:"type_id"`
	Id               int                       `json:"id"`
	CategoryName     string                    `json:"name"`
	CategoryImageUrl string                    `json:"avatar"`
	CategoryPrice    float64                   `json:"price"`
	TotalAmount      float64                   `json:"total_number"`
	CategoryUser     []GroupAgentCommodityUser `json:"users"`
}

type GroupAgentDetail struct {
	GroupId             int                   `json:"id"`
	GroupName           string                `json:"name"`
	GroupStatus         Status                `json:"type"`
	GroupCreatorId      int                   `json:"creator_id"`
	GroupCreatorName    string                `json:"creator_name"`
	GroupCreatorPhone   string                `json:"creator_phone"`
	GroupUserNumber     int                   `json:"user_number"`
	GroupCreatorAddress GroupCreatorAddress   `json:"creator_address"`
	GroupRemark         string                `json:"remark"`
	GroupDescription    string                `json:"description"`
	GroupTotalPrice     float64               `json:"total_price"`
	GroupCreatedAt      int64                 `json:"created_time"`
	GroupUpdatedAt      int64                 `json:"updated_time"`
	GroupRiderPhone     string                `json:"rider_phone"`
	GroupRiderName      string                `json:"rider_name"`
	GroupRiderPos       GroupRiderAddress     `json:"rider_pos"`
	GroupCommodities    []GroupAgentCommodity `json:"commodity_detail"`
}
