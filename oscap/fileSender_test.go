package oscap

import (
	"testing"
)

var (
	fileSenderCorrect = FileSender{"../example/example.xml", "http://localhost:8080"}
	fileSenderNoFile = FileSender{"../example/example-not-exists.xml", "http://localhost:8080"}
)

func TestSendFileToWebhook(t *testing.T) {
	err := fileSenderCorrect.SendFileToWebhook()
	if err != nil {
		t.Errorf("Unexpected error, received %v", err)
	}
}

func TestReadFile(t *testing.T) {
	byteArray, err := fileSenderCorrect.readFile()
	if err != nil || len(byteArray) <= 0 {
		t.Errorf("Unexpected error, received %v", err)
	}
}

func TestReadFileCouldNotFindFile(t *testing.T) {
	_, err := fileSenderNoFile.readFile()
	if err == nil  {
		t.Errorf("We should have received an error. File does not exist")
	}
}

