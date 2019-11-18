package oscap

import (
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"oscap-report-exporter/notify"
	"oscap-report-exporter/oscapLogger"
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
	assert.Equal(t, code, 1)
}

func TestGetConfigDefaults(t *testing.T) {

	logger := getLogger()

	configFileDefault := ""
	configDefault := GetConfig(configFileDefault, logger)

	if configDefault.ScanDate != DefaultConfig.ScanDate {
		t.Errorf("The date as it was parsed by the exaple oscap config is wrong " + configDefault.ScanDate +
			". Should be " + DefaultConfig.ScanDate)
	}
	if configDefault.ScanTime != DefaultConfig.ScanTime {
		t.Errorf("The default time for the scan is wrong " + configDefault.ScanDate +
			". Should be " + DefaultConfig.ScanTime)
	}
	if configDefault.WorkingFolder != DefaultConfig.WorkingFolder {
		t.Errorf("The default working folder for the scan is wrong " + configDefault.WorkingFolder +
			". Should be " + DefaultConfig.WorkingFolder)
	}
	if configDefault.FileName != DefaultConfig.FileName {
		t.Errorf("The default working folder for the scan is wrong " + configDefault.FileName +
			". Should be " + DefaultConfig.FileName)
	}
	if configDefault.Webhook != "" {
		t.Errorf("The default webhook configuration is wrong %v . Should be nil", configDefault.Webhook)
	}
	expectedVulRepURL := DefaultConfig.VulnerabilityReportConf.GlobalVulnerabilityReportHTTPSLocation
	gotVulRepURL := configDefault.VulnerabilityReportConf.GlobalVulnerabilityReportHTTPSLocation
	if gotVulRepURL != expectedVulRepURL {
		t.Errorf("The default global vulnerability report url for the scan is wrong " + gotVulRepURL +
			". Should be " + expectedVulRepURL)
	}

	if configDefault.EmailConfiguration != nil {
		t.Errorf("The default email configuration is wrong %v . Should be nil", configDefault.EmailConfiguration)
	}
}

func TestGetConfigFromTestFullFile(t *testing.T) {

	logger := getLogger()

	configFile := "../test-files/oscap-full-config.yaml"
	config := GetConfig(configFile, logger)

	dateExpected := "Mon"
	if config.ScanDate != dateExpected {
		t.Errorf("The date as it was parsed by the exaple oscap config is wrong " + config.ScanDate +
			". Should be " + dateExpected)
	}
	timeExpected := "23:00"
	if config.ScanTime != timeExpected {
		t.Errorf("The time as it was parsed by the exaple oscap config is wrong " + config.ScanTime +
			". Should be " + timeExpected)
	}
	workFolderExpected := "/tmp/downloads/"
	if config.WorkingFolder != workFolderExpected {
		t.Errorf("The working folder as it was parsed by the exaple oscap config is wrong " + config.WorkingFolder +
			". Should be " + workFolderExpected)
	}
	globVulFileName := "com.redhat.rhsa-all.ds.xml"
	if config.FileName != globVulFileName {
		t.Errorf("The vulnerability file name as it was parsed by the exaple oscap config is wrong " + config.FileName +
			". Should be " + globVulFileName)
	}
	webhook := "http://localhost:8080"
	if config.Webhook != webhook {
		t.Errorf("The webhook as it was parsed by the exaple oscap config is wrong " + config.Webhook +
			". Should be " + webhook)
	}
	expectedVulRepURL := "https://www.redhat.com/security/data/metrics/ds/com.redhat.rhsa-all.ds.xml"
	gotVulRepURL := config.VulnerabilityReportConf.GlobalVulnerabilityReportHTTPSLocation
	if gotVulRepURL != expectedVulRepURL {
		t.Errorf("The vulnerability report url as it was parsed by the exaple oscap config is wrong " +
			gotVulRepURL + ". Should be " + expectedVulRepURL)
	}

	expectedEmailSmarthost := ""
	if config.EmailConfiguration.Smarthost != expectedEmailSmarthost {
		t.Errorf("The smarthost as it was parsed by the exaple oscap config is wrong " +
			config.EmailConfiguration.Smarthost + ". Should be " + expectedEmailSmarthost)
	}

	expectedEmailTo := ""
	if config.EmailConfiguration.To != expectedEmailTo {
		t.Errorf("The To as it was parsed by the exaple oscap config is wrong " +
			config.EmailConfiguration.To + ". Should be " + expectedEmailTo)
	}

	expectedEmailPassword := ""
	if config.EmailConfiguration.Password != expectedEmailPassword {
		t.Errorf("The password as it was parsed by the exaple oscap config is wrong " +
			config.EmailConfiguration.Password + ". Should be " + expectedEmailPassword + ".")
	}
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
	allowLevel := &oscapLogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}
	return oscapLogger.New(allowLevel)
}
