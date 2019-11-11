package oscap

import (
	"log"
	"os"
	"testing"
)

var (
	fileFound    = "example.xml"
	fileNotFound = "example-not-exists.xml"
	webhook      = "http://localhost:8000"
)

func TestSendFileToWebhook(t *testing.T) {
	err := SendFileToWebhook(getPwd()+"/../test-files/", fileFound, webhook)
	if err == nil {
		t.Errorf("Send file to webhook should fail")
	}
}

func TestReadFile(t *testing.T) {

	byteArray, err := readFile(getPwd()+"/../test-files/", fileFound)
	if err != nil || len(byteArray) <= 0 {
		t.Errorf("Unexpected error, received %v", err)
	}
}

func TestReadFileCouldNotFindFile(t *testing.T) {
	_, err := readFile(getPwd()+"/../test-files/", fileNotFound)
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
