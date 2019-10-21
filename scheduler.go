package main

import (
	"log"
	"flag"
	"github.com/jasonlvhit/gocron"
	"oscap-report-exporter/oscap"
)

func main() {

	configFile := flag.String("config.file", "", "the file that contains the configuration for oscap scan")
    flag.Parse()

	startScheduler(*configFile)
	// config := oscap.GetConfig(*configFile)

	// log.Printf("Starting Scheduler for " + config.ScanDate + " " + config.ScanTime)
	// log.Printf("Working folder " + config.WorkingFolder)
	// log.Printf("Global vulnerability report url " + config.VulnerabilityReportConf.BaseVulnerabilityReportUrl + config.VulnerabilityReportConf.GlobalVulnerabilityReportHttpsLocation)
	
	// if config.VulnerabilityReportConf.UserName != "" && config.VulnerabilityReportConf.Password != ""{
	// 	log.Printf(config.VulnerabilityReportConf.UserName)
	// 	log.Printf(config.VulnerabilityReportConf.Password)
	// }	
	// scheduler.Every().Day().At("14:34").Run(job)
	// gocron.Every(1).Monday().At("15:17").Do(job)
	// gocron.Every(1).Tuesday().At("14:40").Do(job)
	// gocron.Every(1).Wednesday().At("14:40").Do(job)
	// gocron.Every(1).Thursday().At("14:40").Do(job)
	// gocron.Every(1).Friday().At("14:40").Do(job)
	// gocron.Every(1).Saturday().At("14:40").Do(job)
	// gocron.Every(1).Sunday().At("14:40").Do(job)

	// <- gocron.Start()
	//var job *scheduler.Job
	// switch config.ScanDate {   
	// 	case "Mon":
	// 		scheduler.Every().Monday().At("11:04").Run(config.OscapVulnerabilityScan)
	// 	// case "Tue":
	// 	// 	job := scheduler.Every().Tuesday()
	// 	// case "Wed":
	// 	// 	job := scheduler.Every().Wednesday()
	// 	// case "Thu":
	// 	// 	job := scheduler.Every().Thursday()
	// 	// case "Fri":
	// 	// 	job := scheduler.Every().Friday()
	// 	// case "Sat":
	// 	// 	job := scheduler.Every().Saturday()
	// 	// case "Sun":
	// 	// 	log.Printf("Chose Sunday")
	// 	// 	job := scheduler.Every().Sunday()
	// 	// case "Daily":
	// 	// 	job := scheduler.Every().Day()
	// 	default:
	// 		log.Fatalf("Scheduling option not supported")
	// }
	//Prevent job from exiting. Is that the right approach?
	//runtime.Goexit()
}

func startScheduler(configFile string) {

	config := oscap.GetConfig(configFile)

	log.Printf("Starting Scheduler for " + config.ScanDate + " " + config.ScanTime)
	log.Printf("Working folder " + config.WorkingFolder)
	log.Printf("Global vulnerability report url " + config.VulnerabilityReportConf.BaseVulnerabilityReportUrl + config.VulnerabilityReportConf.GlobalVulnerabilityReportHttpsLocation)
	
	if config.VulnerabilityReportConf.UserName != "" && config.VulnerabilityReportConf.Password != ""{
		log.Printf(config.VulnerabilityReportConf.UserName)
		log.Printf(config.VulnerabilityReportConf.Password)
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
	<- gocron.Start()
}

