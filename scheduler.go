package main

import (
	"flag"
	"github.com/jasonlvhit/gocron"
	"log"
	"oscap-report-exporter/oscap"
)

func main() {

	configFile := flag.String("config.file", "", "the file that contains the configuration for oscap scan")
	flag.Parse()

	startScheduler(*configFile)
}

func startScheduler(configFile string) {

	config := oscap.GetConfig(configFile)

	log.Printf("Starting Scheduler for " + config.ScanDate + " " + config.ScanTime)
	log.Printf("Working folder " + config.WorkingFolder)
	log.Printf("Global vulnerability report url " + config.VulnerabilityReportConf.GlobalVulnerabilityReportHTTPSLocation)

	if config.VulnerabilityReportConf.UserName != "" && config.VulnerabilityReportConf.Password != "" {
		log.Printf("Username " + config.VulnerabilityReportConf.UserName)
	}

	job := &gocron.Job{}
	switch config.ScanDate {
	case "Mon":
		job = gocron.Every(1).Monday()
	case "Tue":
		job = gocron.Every(1).Tuesday()
	case "Wed":
		job = gocron.Every(1).Wednesday()
	case "Thu":
		job = gocron.Every(1).Thursday()
	case "Fri":
		job = gocron.Every(1).Friday()
	case "Sat":
		job = gocron.Every(1).Saturday()
	case "Sun":
		job = gocron.Every(1).Sunday()
	case "Daily":
		job = gocron.Every(1).Day()
	default:
		log.Fatalf("Scheduling option not supported")
	}

	job.At(config.ScanTime).Do(config.OscapVulnerabilityScan)
	<-gocron.Start()
}
