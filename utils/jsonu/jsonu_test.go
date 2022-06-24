package jsonu

import "testing"

type test struct {
	Status string `json:"status"`
	Info   string `json:"info"`
}

func TestJSONU(t *testing.T) {
	object := test{
		Status: "success",
		Info:   "test",
	}
	json := Marshal(object)
	var new_object test
	Unmarshal(json, &new_object)
	if object != new_object {
		t.Error("jsonu test failed")
	}
}
