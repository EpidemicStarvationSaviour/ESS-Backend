package authUtils

import (
	"errors"
	"ess/utils/setting"

	"github.com/dgrijalva/jwt-go"
)

func ParseToken(token string) (Policy, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Payload{}, func(token *jwt.Token) (i interface{}, err error) {
		return []byte(setting.SecretSetting.JwtKey), nil
	})

	if err != nil {
		return nil, err
	}
	if payload, ok := tokenClaims.Claims.(*Payload); ok && tokenClaims.Valid {
		return payload, nil
	}
	return nil, errors.New("invalid token")
}
