package main

import (
	"runtime"
	"log"
	"flag"
	"testing"
	"github.com/carlescere/scheduler"
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

	scheduler.Every(600).Seconds().Run(config.OscapVulnerabilityScan)

	//Prevent job from exiting. Is that the right approvach?
	runtime.Goexit()
}

