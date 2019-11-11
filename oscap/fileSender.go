package oscap

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// SendFileToWebhook handles sending the created results xml file to the defined webhook
func SendFileToWebhook(workingDir string, file string, webhook string) error {

	if webhook == "" {
		return nil
	}

	byteXML, errReadFile := readFile(workingDir, file)
	if errReadFile != nil {
		return errReadFile
	}

	errWebhook := sender(byteXML, webhook)
	if errWebhook != nil {
		return errWebhook
	}

	return nil

}

// Read the results file and return its content in a bytearray
func readFile(workingDir string, file string) ([]byte, error) {

	// Open our xmlFile
	filePath := filepath.Join(workingDir, filepath.Clean(file))
	xmlFile, errOpen := os.Open(filePath)
	if errOpen != nil {
		log.Printf("Error: Could not open file " + filePath)
		return nil, errOpen
	}

	log.Printf("Successfully Opened " + file)
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, errRead := ioutil.ReadAll(xmlFile)
	if errRead != nil {
		log.Printf("Error: Could not read file " + filePath)
		return nil, errRead
	}
	return byteValue, nil
}

//Send bytearray to webhook in xml format
func sender(byteXML []byte, webhook string) error {
	client := &http.Client{}
	// build a new request, but not doing the POST yet
	req, errMakeReq := http.NewRequest("POST", webhook, bytes.NewBuffer(byteXML))
	if errMakeReq != nil {
		log.Printf("Error: Could not create new request containing the xml bytearray. ")
		return errMakeReq
	}

	log.Printf("Sending results to webhook " + webhook)

	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	// now POST it
	resp, errHTTP := client.Do(req)
	if errHTTP != nil {
		log.Printf("Error: Posting the file to the webhook failed. ")
		return errHTTP
	}
	log.Printf("Webhook response " + string(resp.Status))
	return nil
}
