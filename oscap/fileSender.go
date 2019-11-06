package oscap

import (
	"log"
	"os"
	"io/ioutil"
	"bytes"
	"net/http"
)

func SendFileToWebhook(file string, webhook string) error {
	byteXml, errReadFile := readFile(file)
	if errReadFile != nil {
		return errReadFile
	}

	errWebhook := sender(byteXml, webhook)
	if errWebhook != nil {
		return errWebhook
	}

	return nil

}

// Read the results file and return its content in a bytearray
func readFile(file string) ([]byte, error) {

	// Open our xmlFile
    xmlFile, errOpen := os.Open(file)
    if errOpen != nil {
        log.Printf("Error: Could not open file " + file)
        return nil, errOpen
    }
    
    log.Printf("Successfully Opened " + file)
    defer xmlFile.Close()

    // read our opened xmlFile as a byte array.
    byteValue, errRead := ioutil.ReadAll(xmlFile)
    if errRead != nil {
        log.Printf("Error: Could not read file " + file)
        return nil, errRead
    }
	return byteValue, nil
}

//Send bytearray to webhook in xml format
func sender(byteXml []byte, webhook string) error {
	client := &http.Client{}
	// build a new request, but not doing the POST yet
	req, errMakeReq := http.NewRequest("POST", webhook, bytes.NewBuffer(byteXml))
	if errMakeReq != nil {
		log.Printf("Error: Could not create new request containing the xml bytearray. ")
		return errMakeReq
	}

    log.Printf("Sending results to webhook " + webhook)

	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	// now POST it
	resp, errHttp := client.Do(req)
	if errHttp != nil {
		log.Printf("Error: Posting the file to the webhook failed. ")
		return errHttp
	}
	log.Printf( "Webhook response " + string(resp.Status) )
	return nil
}