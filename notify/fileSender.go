package notify

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/parnurzeal/gorequest"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// FileSend contains the information for the webhook and the files to be sent
type FileSend struct {
	logger     log.Logger
	workingDir string
	file       string
	webhook    string
}

// NewFileSender creates a new instance of FileSender
func NewFileSender(logger log.Logger, workingDir string, file string, webhook string) *FileSend {
	fs := &FileSend{logger, workingDir, file, webhook}
	return fs
}

// SendFileToWebhook handles sending the created results xml file to the defined webhook
func (fs *FileSend) SendFileToWebhook() error {

	if fs.webhook == "" {
		return nil
	}

	byteXML, errReadFile := fs.readFile()
	if errReadFile != nil {
		return errReadFile
	}

	level.Debug(fs.logger).Log("msg", "Going to sender")

	errWebhook := fs.sender(byteXML)
	if errWebhook != nil {
		return errWebhook
	}

	return nil
}

// Read the results file and return its content in a bytearray
func (fs *FileSend) readFile() ([]byte, error) {

	// Open our xmlFile
	filePath := filepath.Join(fs.workingDir, filepath.Clean(fs.file))
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, errors.Wrap(err, "Could not open file "+filePath)
	}

	level.Debug(fs.logger).Log("msg", "Successfully Opened "+fs.file)
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read file "+filePath)
	}
	return byteValue, nil
}

//Send bytearray to webhook in map format map[results: <xml oscap results>]
func (fs *FileSend) sender(byteXML []byte) error {

	m := make(map[string]string)
	m["results"] = string(byteXML)

	request := gorequest.New()
	_, _, errs := request.Post(fs.webhook).
		SendMap(m).
		Retry(3, 20*time.Second, http.StatusBadRequest, http.StatusRequestTimeout, http.StatusInternalServerError, http.StatusGatewayTimeout).
		End()
	if errs != nil {
		respErr := errors.New("Could not send to webhook " + fs.webhook)
		for _, err := range errs {
			respErr = errors.Wrap(respErr, fmt.Sprint(err))
		}
		return respErr
	}

	return nil
}
