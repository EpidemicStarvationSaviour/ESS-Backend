package crypto

import (
	"crypto/sha1"
	"encoding/hex"
	"ess/utils/setting"
	"strings"
)

// WARNING: it will transform string to upper
func SHA1(content string) string {
	h := sha1.New()
	h.Write([]byte(content))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// convert the password(plaintext) to salted, use double SHA1 and saltA saltB
func Password2Secret(password string) string {
	saltA := setting.SecretSetting.SaltA
	saltB := setting.SecretSetting.SaltB
	return SHA1(saltB + SHA1(saltA+password))
}
