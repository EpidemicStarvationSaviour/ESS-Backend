package amap_base

import "testing"

func TestAmapBase(t *testing.T) {
	_, err := DistanceByCoordination(0, 0, 0, 0)
	if err != nil {
		t.Error(err)
	}
}
