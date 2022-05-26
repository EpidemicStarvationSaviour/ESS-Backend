package admin

type UserAddress struct {
	AddressId        int     `json:"id"`
	AddressLat       float64 `json:"lat"`
	AddressLng       float64 `json:"lng"`
	AddressProvince  string  `json:"province"`
	AddressCity      string  `json:"city"`
	AddressArea      string  `json:"area"`
	AddressDetail    string  `json:"detail"`
	AddressIsDefault bool    `json:"is_default"`
}

type UserData struct {
	UserId        int           `json:"user_id"`
	UserPhone     string        `json:"user_phone"`
	UserName      string        `json:"user_name"`
	UserRole      int           `json:"user_role"`
	UserAddresses []UserAddress `json:"user_address"`
	UserCreatedAt int64         `json:"user_created_time"`
	UserUpdatedAt int64         `json:"user_updated_time"`
}

type AllUserResp struct {
	UserCount int        `json:"count"`
	UserList  []UserData `json:"data"`
}

type AdminDeleteUser struct {
	UserId int `json:"user_id"`
}
