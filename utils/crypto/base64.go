package crypto

import "encoding/base64"

func Base64Encode(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}

func Base64Decode(content string) (string, error) {
	s, err := base64.StdEncoding.DecodeString(content)
	return string(s), err
}
