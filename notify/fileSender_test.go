package notify

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
	errChan := make(chan error)

	fs := FileSend{logger, getPwd() + "/../test-files/", fileFound, webhook}
	go fs.SendFileToWebhook(errChan)
	if <-errChan == nil {
		t.Errorf("Send file to webhook should fail")
	}
}

func TestReadFile(t *testing.T) {
	fs := FileSend{logger, getPwd() + "/../test-files/", fileFound, webhook}
	byteArray, err := fs.readFile()
	if err != nil || len(byteArray) <= 0 {
		t.Errorf("Unexpected error, received %v", err)
	}
}

func TestReadFileCouldNotFindFile(t *testing.T) {
	fs := FileSend{logger, getPwd() + "/../test-files/", fileNotFound, webhook}
	_, err := fs.readFile()
	if err == nil {
		t.Errorf("We should have received an error. File does not exist")
	}
}

func TestSendDataToEndpoint(t *testing.T) {
	errChan := make(chan error)

	fs := NewFileSender(logger, getPwd()+"/../test-files/", fileFound, "")
	go fs.SendFileToWebhook(errChan)
	if <-errChan != nil {
		t.Errorf("Send file to webhook should not happen. The call should return nil since the webhook is an empty string")
	}
}

func getPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
