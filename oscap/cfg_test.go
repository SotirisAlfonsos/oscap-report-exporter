package oscap

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"oscap-report-exporter/notify"
	"oscap-report-exporter/oscaplogger"
	"testing"
)

var logger = getLogger()

func TestSendResultsToChannel(t *testing.T) {
	conf := GetConfig("", logger)
	conf.Webhook = "http://localhost:8000"
	conf.EmailConfiguration = &notify.EmailConf{
		Smarthost: "smarthost:25",
		From:      "from",
		To:        "to",
		Password:  "",
	}
	err := conf.sendResultsToChannels(logger)
	assert.EqualError(t, err, "Could not send results to all available channels")
}

func TestSendResultsToChannelNoWebhook(t *testing.T) {
	conf := GetConfig("", logger)
	conf.Webhook = ""
	conf.EmailConfiguration = &notify.EmailConf{
		Smarthost: "smarthost",
		From:      "from",
		To:        "to",
		Password:  "",
	}
	err := conf.sendResultsToChannels(logger)
	assert.EqualError(t, err, "Could not send results to all available channels")
}

func TestSendResultsToChannelNoWebhookNoSmarthost(t *testing.T) {
	conf := GetConfig("", logger)
	conf.Webhook = ""
	conf.EmailConfiguration = &notify.EmailConf{
		Smarthost: "",
		From:      "from",
		To:        "to",
		Password:  "",
	}
	err := conf.sendResultsToChannels(logger)
	assert.NoError(t, err)
}

func TestSendResultsToChannelNoWebhookNoMailConf(t *testing.T) {
	conf := GetConfig("", logger)
	conf.Webhook = ""
	err := conf.sendResultsToChannels(logger)
	assert.NoError(t, err)
}

func TestPrepareAndRunScanFailDownload(t *testing.T) {
	conf := GetConfig("", logger)
	conf.VulnerabilityReportConf.GlobalVulnerabilityReportHTTPSLocation = ""
	code := conf.prepareAndRunScan(logger)
	assert.Equal(t, 1, code)
}

func TestGetConfigDefaults(t *testing.T) {

	logger := getLogger()

	configFileDefault := ""
	configDefault := GetConfig(configFileDefault, logger)

	assert.Equal(t, configDefault.ScanDate, DefaultConfig.ScanDate, "Date")
	assert.Equal(t, configDefault.ScanTime, DefaultConfig.ScanTime, "Time")
	assert.Equal(t, configDefault.WorkingFolder, DefaultConfig.WorkingFolder, "Working folder")
	assert.Equal(t, configDefault.Webhook, "", "Webhook")
	assert.Equal(t, configDefault.Profile, "", "Profile")
	assert.Equal(t, configDefault.VulnerabilityReportConf.GlobalVulnerabilityReportHTTPSLocation,
		DefaultConfig.VulnerabilityReportConf.GlobalVulnerabilityReportHTTPSLocation, "Vulnerability report https location")
	assert.Nil(t, configDefault.EmailConfiguration)
}

func TestGetConfigFromTestFullFile(t *testing.T) {

	logger := getLogger()

	configFile := "../test-files/oscap-full-config.yaml"
	config := GetConfig(configFile, logger)

	assert.Equal(t, config.ScanDate, "Mon", "Date")
	assert.Equal(t, config.ScanTime, "23:00", "Time")
	assert.Equal(t, config.WorkingFolder, "/tmp/downloads/", "Working folder")
	assert.Equal(t, config.Webhook, "http://localhost:8080", "Webhook")
	assert.Equal(t, config.Profile, "xccdf_org.ssgproject.content_profile_C2S", "Profile")
	assert.Equal(t, config.VulnerabilityReportConf.GlobalVulnerabilityReportHTTPSLocation,
		"https://www.redhat.com/security/data/metrics/ds/com.redhat.rhsa-all.ds.xml", "Vulnerability report https location")
	assert.Equal(t, config.EmailConfiguration.Smarthost, "", "Smarthost")
	assert.Equal(t, config.EmailConfiguration.To, "", "To")
	assert.Equal(t, config.EmailConfiguration.Password, "", "Password")

}

func TestGetConfigFromTestOmitedFile(t *testing.T) {
	logger := getLogger()

	configFile := "../test-files/oscap-omited-config.yaml"
	config := GetConfig(configFile, logger)

	if config.Webhook != "" {
		t.Errorf("The webhook as it was parsed by the exaple oscap config is wrong " + config.Webhook +
			". Should be empty string")
	}

	if config.EmailConfiguration != nil {
		t.Errorf("The email configuration from the exaple oscap config should be nil.")
	}
}

func getLogger() log.Logger {
	allowLevel := &oscaplogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}
	return oscaplogger.New(allowLevel)
}
