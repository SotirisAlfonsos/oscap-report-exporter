package oscap

import (
	"fmt"
	"log"
	"testing"
)

var (
	smarthost = "dummy:25"
	from = ""
	to = ""
	password = ""
)

func TestSendEmailReportDoesNotExist(t *testing.T) {
	nonExistentPath := "/path/to/report/that/should/not/exist"
	emailConf := EmailConf{smarthost, from, to, ""}
	err := emailConf.SendFileViaEmail(nonExistentPath)
	if err == nil {
		t.Errorf("File should not exist in the path " + nonExistentPath)
	}


	if expectedErr := "stat " + nonExistentPath + ": no such file or directory"; fmt.Sprint(err) != expectedErr {
		t.Errorf("Error received was different from error expected. \nExpected Err: " + expectedErr + "\nGot: " + fmt.Sprint(err))
	}
}

func TestSendEmailCoundNotConctactSmarthost(t *testing.T) {
	reportPath := "../example/report.html"
	emailConf := EmailConf{smarthost, from, to, password}
	err := emailConf.SendFileViaEmail(reportPath)
	log.Printf(fmt.Sprint(err))
	if err == nil {
		t.Errorf("Should be able to send email to " + to + ". Smarthost " + smarthost + " does not exist.")
	} 

	// if expectedErr := "dial tcp: lookup " + strings.Split(smarthost, ":")[0] + ": no such host"; fmt.Sprint(err) != expectedErr {
	// 	t.Errorf("Error received was different from error expected. \nExpected Err: " + expectedErr + "\nGot: " + fmt.Sprint(err))
	// }
	
}
