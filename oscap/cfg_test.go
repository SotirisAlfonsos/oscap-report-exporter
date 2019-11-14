package oscap

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"oscap-report-exporter/notify"
	"oscap-report-exporter/oscapLogger"
	"testing"
)

var logger = getLogger()

func TestSendResultsToChannel(t *testing.T) {
	conf := GetConfig("", logger)
	conf.Webhook = "http://localhost:8000"
	conf.EmailConfiguration = &notify.EmailConf{
		Smarthost: "smarthost",
		From:      "from",
		To:        "to",
		Password:  "",
	}
	conf.sendResultsToChannels(logger)
}

func getLogger() log.Logger {
	allowLevel := &oscapLogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}
	return oscapLogger.New(allowLevel)
}
