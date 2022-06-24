package algorithm

import (
	"ess/utils/setting"
	"testing"
)

func TestAlgorithm(t *testing.T) {
	setting.GRPCSetting.Enable = false
	Setup()
	_ = Schedule(1, 1)
}
