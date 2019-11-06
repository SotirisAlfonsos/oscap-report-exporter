package oscap

import (
	"log"
	"os"
	"testing"
)

var (
	fileLoc      = "example.xml"
	fileLocFalse = "example-not-exists.xml"
	webhook      = "http://localhost:8080"
)

func TestSendFileToWebhook(t *testing.T) {
	SendFileToWebhook(getPwd()+"/../example/", fileLoc, webhook)
}

func TestReadFile(t *testing.T) {

	byteArray, err := readFile(getPwd()+"/../example/", fileLoc)
	if err != nil || len(byteArray) <= 0 {
		t.Errorf("Unexpected error, received %v", err)
	}
}

func TestReadFileCouldNotFindFile(t *testing.T) {
	_, err := readFile(getPwd()+"/../example/", fileLocFalse)
	if err == nil {
		t.Errorf("We should have received an error. File does not exist")
	}
}

func getPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
