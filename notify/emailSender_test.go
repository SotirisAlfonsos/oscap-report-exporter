package notify

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"oscap-report-exporter/oscapLogger"
	"testing"
)

var (
	smarthost       = "dummy:25"
	from            = "from"
	to              = "to"
	password        = ""
	reportPath      = "../test-files/report.html"
	nonExistentPath = "/path/to/report/that/should/not/exist"
	logger          = getLogger()
)

func getLogger() log.Logger {
	allowLevel := &oscapLogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}
	return oscapLogger.New(allowLevel)
}

func TestEmailSenderReportDoesNotExist(t *testing.T) {
	errChan := make(chan error)

	emailConf := EmailConf{smarthost, from, to, ""}
	go emailConf.SendFileViaEmail(nonExistentPath, logger, errChan)
	if <-errChan == nil {
		t.Errorf("File should not exist in the path " + nonExistentPath)
	}

}

func TestEmailSenderCoundNotConctactSmarthost(t *testing.T) {
	errChan := make(chan error)

	emailConf := EmailConf{smarthost, from, to, password}
	go emailConf.SendFileViaEmail(reportPath, logger, errChan)
	if <-errChan == nil {
		t.Errorf("Should not be able to send email to " + to + ". Smarthost " + smarthost + " does not exist.")
	}

}

func TestEmailSenderNoAuth(t *testing.T) {
	emailConf := EmailConf{smarthost, from, to, password}
	auth := emailConf.configureAuth()
	if auth != nil {
		t.Errorf("Auth should be nil since no pwd is provided. Instead it was %v", auth)
	}
}

func TestEmailSenderNoSmarthostDetails(t *testing.T) {
	errChan := make(chan error)

	emailConf := EmailConf{"", from, to, password}
	go emailConf.SendFileViaEmail(reportPath, logger, errChan)
	if <-errChan != nil {
		t.Errorf("Error should be nil and the email should not be sent")
	}
}
