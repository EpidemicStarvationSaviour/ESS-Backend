package authUtils

import (
	"ess/model/user"
	"ess/utils/setting"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Policy interface {
	AdminOnly() bool
	LoginOnly() bool
	CheckExpired() bool
	SysAdminOnly() bool
	GetId() int
	GetName() string
	GetPhone() string
	ConvertToUser() user.User
}

type Payload struct {
	Name   string    `json:"name"`
	UserId int       `json:"user_id"`
	Phone  string    `json:"phone"`
	Role   user.Role `json:"role"`
	jwt.StandardClaims
}

func (p *Payload) LoginOnly() bool {
	return p.Role != user.NoLogin
}

func (p *Payload) AdminOnly() bool {
	return p.Role == user.Admin || p.Role == user.SysAdmin
}

func (p *Payload) CheckExpired() bool {
	return time.Now().Unix() < p.ExpiresAt
}

func (p *Payload) SysAdminOnly() bool {
	return p.Role == user.SysAdmin
}

func (p *Payload) GetId() int {
	return p.UserId
}

func (p *Payload) GetName() string {
	return p.Name
}

func (p *Payload) GetPhone() string {
	return p.Phone
}

func (p *Payload) ConvertToUser() user.User {
	return user.User{
		UserRole: p.Role,
		UserName: p.Name,
		UserId:   p.UserId,
	}
}

// convert token.LoginReq to Payload (already init the jwt.StandardClaims)
func GetClaimFromUser(user user.User) *Payload {
	nowTime := time.Now()
	expireTime := nowTime.Add(setting.ServerSetting.JwtExpireTime)

	return &Payload{
		Name:   user.UserName,
		UserId: user.UserId,
		Role:   user.UserRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    setting.SecretSetting.JwtIssuer,
		},
	}
}

// get a payload for admin
func GetClaimFromSysAdmin() *Payload {
	nowTime := time.Now()
	expireTime := nowTime.Add(setting.ServerSetting.JwtExpireTime)

	return &Payload{
		Name:   setting.AdminSetting.Name,
		UserId: setting.AdminSetting.UserId,
		Phone:  setting.AdminSetting.Phone,
		Role:   user.SysAdmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    setting.SecretSetting.JwtIssuer,
		},
	}
}
