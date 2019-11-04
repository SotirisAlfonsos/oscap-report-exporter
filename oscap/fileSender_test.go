package oscap

import (
	"testing"
)

var (
	fileLoc = "../example/example.xml"
	fileLocFalse = "../example/example-not-exists.xml"
	webhook = "http://localhost:8080"
)

func TestSendFileToWebhook(t *testing.T) {
	SendFileToWebhook(fileLoc, webhook)
}

func TestReadFile(t *testing.T) {
	byteArray, err := readFile(fileLoc)
	if err != nil || len(byteArray) <= 0 {
		t.Errorf("Unexpected error, received %v", err)
	}
}

func TestReadFileCouldNotFindFile(t *testing.T) {
	_, err := readFile(fileLocFalse)
	if err == nil  {
		t.Errorf("We should have received an error. File does not exist")
	}
}

