package user_service

import "ess/model/user"

// return a AuthResp model with a use model and token
func NewLoginResp(us user.User, _token string, _type string) *user.AuthResp {
	return &user.AuthResp{
		UserPhone: us.UserPhone,
		UserName:  us.UserName,
		UserRole:  us.UserRole,
		UserToken: _token,
		LoginType: _type,
	}
}
