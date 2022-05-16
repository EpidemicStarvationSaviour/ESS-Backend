package cache

import (
	"ess/utils/stringu"
	"fmt"
)

type CacheType uint8

const (
	UserInfo CacheType = iota
	PhoneToId
	NameToId
	SubScription
)

func hash(t CacheType, id interface{}) string {
	return fmt.Sprintf("%d#!%s", t, stringu.Tostring(id))
}

// return the key of cache
func GetKey(t CacheType, id interface{}) string {
	return hash(t, id)
}

// return the key of get SubScription list
func GetKeySubScribeList() string {
	return hash(SubScription, -1)
}
