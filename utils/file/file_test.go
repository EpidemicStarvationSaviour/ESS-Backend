package file

import "testing"

func TestFile(t *testing.T) {
	file := "file_test.go"
	if !CheckExist(file) {
		t.Error("file not exist")
	}
}
