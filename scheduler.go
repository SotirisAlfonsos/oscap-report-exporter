package main

import (
	"runtime"
	"log"
	"flag"
	"github.com/carlescere/scheduler"
	"oscap-report-exporter/oscap"
)

func main() {

	configFile := flag.String("config.file", "", "the file that contains the configuration for oscap scan")
    flag.Parse()

	startScheduler(*configFile)

	//Prevent job from exiting. Is that the right approvach?
	runtime.Goexit()
}

func startScheduler(configFile string) {

	config := oscap.GetConfig(configFile)

	log.Printf("Starting Scheduler for " + config.ScanDate + " " + config.ScanTime)
	log.Printf(config.WorkingFolder)
	log.Printf(config.VulnerabilityReportConf.GlobalVulnerabilityReportHttpsLocation)
	log.Printf(config.VulnerabilityReportConf.UserName)
	log.Printf(config.VulnerabilityReportConf.Password)
	log.Printf(config.VulnerabilityReportConf.BaseVulnerabilityReportUrl)

	var job *scheduler.Job
	switch config.ScanDate {   
		case "Mon":
			job = scheduler.Every().Monday()
		case "Tue":
			job = scheduler.Every().Tuesday()
		case "Wed":
			job = scheduler.Every().Wednesday()
		case "Thu":
			job = scheduler.Every().Thursday()
		case "Fri":
			job = scheduler.Every().Friday()
		case "Sat":
			job = scheduler.Every().Saturday()
		case "Sun":
			job = scheduler.Every().Sunday()
		case "Daily":
			job = scheduler.Every().Day()
		default:
			log.Fatalf("Scheduling option not supported")
	}
	job.At(config.ScanTime).Run(config.OscapVulnerabilityScan)

}

