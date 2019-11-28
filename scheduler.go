package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jasonlvhit/gocron"
	"github.com/pkg/errors"
	"os"
	"oscap-report-exporter/oscap"
	"oscap-report-exporter/oscaplogger"
)

func main() {

	configFile := flag.String("config.file", "", "the file that contains the configuration for oscap scan")
	debugLevel := flag.String("debug.level", "info", "the debug level for the exporter. Could be debug, info, warn, error.")
	flag.Parse()

	startScheduler(*configFile, *debugLevel)
}

func startScheduler(configFile string, debugLevel string) {
	logger := createLogger(debugLevel)

	config := oscap.GetConfig(configFile, logger)

	level.Info(logger).Log("msg", "Starting Scheduler for "+config.ScanDate+" "+config.ScanTime)
	level.Info(logger).Log("msg", "Working folder "+config.WorkingFolder)
	level.Info(logger).Log("msg", "Global vulnerability report url "+config.VulnerabilityReportConf.GlobalVulnerabilityReportHTTPSLocation)

	if config.VulnerabilityReportConf.UserName != "" && config.VulnerabilityReportConf.Password != "" {
		level.Debug(logger).Log("msg", "Username "+config.VulnerabilityReportConf.UserName)
	}

	job, err := createJob(config.ScanDate)
	if err != nil {
		level.Error(logger).Log("msg", "Could not schedule job ", "err", err)
		os.Exit(1)
	}

	job.At(config.ScanTime).Do(config.OscapVulnerabilityScan, logger)
	<-gocron.Start()
}

func createJob(date string) (*gocron.Job, error) {
	var job *gocron.Job
	switch date {
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
		return nil, errors.New("Scheduling option not supported")
	}
	return job, nil
}

func createLogger(debugLevel string) log.Logger {
	allowLevel := &oscaplogger.AllowedLevel{}
	if err := allowLevel.Set(debugLevel); err != nil {
		fmt.Printf("%v", err)
	}
	return oscaplogger.New(allowLevel)

}
