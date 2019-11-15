package common

import (
	"log"
	"os"
	"testing"
)

func TestFileOpsFileExistsFileDoesNotExist(t *testing.T) {
	dummyFile := "/dummy/path/should/not/exist"
	if err := FileExists(dummyFile); err == nil {
		t.Errorf("This should lead to an error because " + dummyFile + " should not exist")
	}
}

func TestFileOpsFileExistsFileIsADir(t *testing.T) {
	folder := "/bin"
	if err := FileExists(folder); err == nil {
		t.Errorf("This should lead to an error because " + folder + " is a folder")
	}
}

func TestFileOpsFileExists(t *testing.T) {
	file := getPwd() + "/../test-files/example.xml"
	if err := FileExists(file); err != nil {
		t.Error("The file " + file + " should exist. Something went wrong")
	}
}

func getPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
