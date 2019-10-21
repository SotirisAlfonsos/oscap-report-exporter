package oscap

import (
	"log"
	"os"
	"io/ioutil"
	"fmt"
	"bytes"
	"net/http"
)

type FileSender struct {
	File string
	Webhook string
}

func (sf *FileSender) SendFileToWebhook() error {
	byteXml, err := sf.readFile()
	if err != nil {
		log.Printf("Error: reading " + sf.File + " : " + fmt.Sprint(err))
	}

	sf.sender(byteXml)

	return nil
}

func (sf *FileSender) readFile() ([]byte, error) {

	// Open our xmlFile
    xmlFile, errOpen := os.Open(sf.File)
    if errOpen != nil {
        log.Printf("Error: Could not open file")
        return nil, errOpen
    }
    
    log.Printf("Successfully Opened " + sf.File)
    defer xmlFile.Close()

    // read our opened xmlFile as a byte array.
    byteValue, errRead := ioutil.ReadAll(xmlFile)
    if errRead != nil {
        log.Printf("Error: Could not read file")
        return nil, errRead
    }
	return byteValue, nil
}

func (sf *FileSender) sender(byteXml []byte) error{
	client := &http.Client{}
	// build a new request, but not doing the POST yet
	req, errMakeReq := http.NewRequest("POST", sf.Webhook, bytes.NewBuffer(byteXml))
	if errMakeReq != nil {
		log.Printf("Could not create new request containing the xml bytearray. ")
		return errMakeReq
	}

    log.Printf("Sending results to webhook " + sf.Webhook)

	req.Header.Add("Content-Type", "application/xml; charset=utf-8")
	// now POST it
	resp, errHttp := client.Do(req)
	if errHttp != nil {
		fmt.Println("Posting the file to the webhook failed. ")
		return errHttp
	}
	fmt.Println(resp.Status)
	return nil
}