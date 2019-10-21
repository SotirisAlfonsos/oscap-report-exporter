package main

import (
	"runtime"
	"log"
	"flag"
	"testing"
	"github.com/jasonlvhit/gocron"
	"oscap-report-exporter/oscap"
)

func TestMain(m *testing.M) {

	configFile := flag.String("config.file", "example/oscap-config.yaml", "the file that contains the configuration for oscap scan")
    flag.Parse()

	config := oscap.GetConfig(*configFile)
	log.Printf(config.ScanDate)
	log.Printf(config.ScanTime)
	log.Printf(config.WorkingFolder)
	log.Printf(config.VulnerabilityReportConf.GlobalVulnerabilityReportHttpsLocation)
	log.Printf(config.VulnerabilityReportConf.UserName)
	log.Printf(config.VulnerabilityReportConf.Password)
	log.Printf(config.VulnerabilityReportConf.BaseVulnerabilityReportUrl)

	gocron.Every(600).Seconds().Do(config.OscapVulnerabilityScan)

	<- gocron.Start()
}

