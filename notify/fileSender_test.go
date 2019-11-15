package notify

import (
	"github.com/stretchr/testify/assert"
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

	fs := FileSend{logger, getPwd() + "/../test-files/", fileFound, webhook}
	err := fs.SendFileToWebhook()

	assert.Error(t, err)
}

func TestReadFile(t *testing.T) {
	fs := FileSend{logger, getPwd() + "/../test-files/", fileFound, webhook}
	byteArray, err := fs.readFile()

	assert.NoError(t, err)
	assert.Greater(t, len(byteArray), 1)
}

func TestReadFileCouldNotFindFile(t *testing.T) {
	fs := FileSend{logger, getPwd() + "/../test-files/", fileNotFound, webhook}
	err := fs.SendFileToWebhook()
	//_, err := fs.readFile()
	assert.Error(t, err)
}

func TestSendDataToEndpoint(t *testing.T) {

	fs := NewFileSender(logger, getPwd()+"/../test-files/", fileFound, "")
	err := fs.SendFileToWebhook()

	assert.NoError(t, err)
}

func getPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
