package stringu

import "testing"

func TestStringu(t *testing.T) {
	if Tostring("string") != "string" {
		t.Error("Tostring(\"string\") != \"string\"")
	}
	if Tostring(1) != "1" {
		t.Error("Tostring(1) != \"1\"")
	}
	if Tostring(1.1) != "1.100000" {
		t.Error("Tostring(1.1) != \"1.1\"")
	}
}
