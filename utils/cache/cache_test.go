package cache

import (
	"ess/utils/setting"
	"testing"
)

func TestCache(t *testing.T) {
	setting.ServerSetting.CacheSize = 1
	Setup()
	Set("test", "test")
	if GetOrCreate("test", nil) != "test" {
		t.Error("cache test fail")
	}
}
