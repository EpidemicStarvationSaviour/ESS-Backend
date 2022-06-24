package authUtils

import (
	"ess/model/user"
	"testing"
)

func TestAuth(t *testing.T) {
	userID := 1
	jwt, err := GetUserToken(user.User{
		UserId: userID,
	})
	if err != nil {
		t.Error(err)
	}
	policy, err := ParseToken(jwt)
	if err != nil {
		t.Error(err)
	}
	if userID != policy.GetId() {
		t.Error("userID != policy.GetId()")
	}
}
