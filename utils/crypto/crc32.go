package crypto

import "hash/crc32"

func EncodeCrc32(str string) uint32 {
	return crc32.ChecksumIEEE([]byte(str))
}
