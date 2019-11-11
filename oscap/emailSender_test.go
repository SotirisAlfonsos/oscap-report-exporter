package oscap

import (
	"fmt"
	"log"
	"testing"
)

var (
	smarthost = "dummy:25"
	from      = "from"
	to        = "to"
	password  = ""
)

func TestSendEmailReportDoesNotExist(t *testing.T) {
	nonExistentPath := "/path/to/report/that/should/not/exist"
	emailConf := EmailConf{smarthost, from, to, ""}
	err := emailConf.SendFileViaEmail(nonExistentPath)
	if err == nil {
		t.Errorf("File should not exist in the path " + nonExistentPath)
	}

}

func TestSendEmailCoundNotConctactSmarthost(t *testing.T) {
	reportPath := "../test-files/report.html"
	emailConf := EmailConf{smarthost, from, to, password}
	err := emailConf.SendFileViaEmail(reportPath)
	if err == nil {
		t.Errorf("Should be able to send email to " + to + ". Smarthost " + smarthost + " does not exist.")
	}

}
