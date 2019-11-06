package main

import (
	"os"
	"flag"
	"testing"
	"oscap-report-exporter/oscap"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	// configFile := "example/oscap-config.yaml"
	// config := oscap.GetConfig(configFile)
	// config.CleanFiles = false

	// config.OscapVulnerabilityScan()

	// log.Printf("Verify that report and downloaded files exist")
	// if !fileExists(config.WorkingFolder+"results.xml") || !fileExists(config.WorkingFolder+config.FileName) {
	// 	log.Fatalf("One of the files we expected does not exist. Fail the tests")
	// }

	// errRemoveDownload := os.Remove(config.WorkingFolder + config.FileName)
	// if errRemoveDownload != nil {
	// 	log.Fatal("Unable to remove " + config.FileName + " with error " + fmt.Sprint(errRemoveDownload))
	// }
	// errRemoveResults := os.Remove(config.WorkingFolder + "results.xml")
	// if errRemoveResults != nil {
	// 	log.Fatal("Unable to remove results.xml with error " + fmt.Sprint(errRemoveResults))
	// }

	exitCode := m.Run()
	os.Exit(exitCode)
	
}

func TestConfigDefaults(t *testing.T) {

	configFileDefault := ""
	configDefault := oscap.GetConfig(configFileDefault)
	
	if configDefault.ScanDate != oscap.DefaultConfig.ScanDate {
		t.Errorf("The date as it was parsed by the exaple oscap config is wrong " + configDefault.ScanDate +
			". Should be " + oscap.DefaultConfig.ScanDate)
	}
	if configDefault.ScanTime != oscap.DefaultConfig.ScanTime {
		t.Errorf("The default time for the scan is wrong " + configDefault.ScanDate + 
			". Should be " + oscap.DefaultConfig.ScanTime)
	}
	if configDefault.WorkingFolder != oscap.DefaultConfig.WorkingFolder {
		t.Errorf("The default working folder for the scan is wrong " + configDefault.WorkingFolder + 
			". Should be " + oscap.DefaultConfig.WorkingFolder)
	}
	if configDefault.FileName != oscap.DefaultConfig.FileName {
		t.Errorf("The default working folder for the scan is wrong " + configDefault.FileName + 
			". Should be " + oscap.DefaultConfig.FileName)
	}
	expectedVulRepUrl := oscap.DefaultConfig.VulnerabilityReportConf.BaseVulnerabilityReportUrl + 
		oscap.DefaultConfig.VulnerabilityReportConf.GlobalVulnerabilityReportHttpsLocation
	gotVulRepUrl := configDefault.VulnerabilityReportConf.BaseVulnerabilityReportUrl + 
		configDefault.VulnerabilityReportConf.GlobalVulnerabilityReportHttpsLocation
	if gotVulRepUrl != expectedVulRepUrl {
		t.Errorf("The default global vulnerability report url for the scan is wrong " + gotVulRepUrl + 
			". Should be " + expectedVulRepUrl)
	}

	if configDefault.EmailConfiguration != nil {
		t.Errorf("The default email configuration is wrong %v . Should be nil", configDefault.EmailConfiguration )
	}
}

func TestConfigFromExampleFile(t *testing.T) {
	configFile := flag.String("config.file", "example/oscap-config.yaml", "the file that contains the configuration for oscap scan")
    flag.Parse()

    dateExpected := "Mon"
	config := oscap.GetConfig(*configFile)
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
	expectedVulRepUrl := "https://www.redhat.com/" + "security/data/metrics/ds/com.redhat.rhsa-all.ds.xml"
	gotVulRepUrl := config.VulnerabilityReportConf.BaseVulnerabilityReportUrl + config.VulnerabilityReportConf.GlobalVulnerabilityReportHttpsLocation
	if gotVulRepUrl != expectedVulRepUrl {
		t.Errorf("The vulnerability report url as it was parsed by the exaple oscap config is wrong " + 
			gotVulRepUrl + ". Should be " + expectedVulRepUrl)
	}

	expectedEmailSmarthost := "smarthost"
	if config.EmailConfiguration.Smarthost != expectedEmailSmarthost {
		t.Errorf("The smarthost as it was parsed by the exaple oscap config is wrong " + 
			config.EmailConfiguration.Smarthost + ". Should be " + expectedEmailSmarthost)
	}

	expectedEmailTo := "to"
	if config.EmailConfiguration.To != expectedEmailTo {
		t.Errorf("The To as it was parsed by the exaple oscap config is wrong " + 
			config.EmailConfiguration.To + ". Should be " + expectedEmailTo)
	}

	expectedEmailPassword := ""
	if config.EmailConfiguration.Password != expectedEmailPassword {
		t.Errorf("The username as it was parsed by the exaple oscap config is wrong " + 
			config.EmailConfiguration.Password + ". Should be " + expectedEmailPassword +".")
	}
}

