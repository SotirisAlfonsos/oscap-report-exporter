package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/jasonlvhit/gocron"
	"os"
	"oscap-report-exporter/oscap"
	"oscap-report-exporter/oscapLogger"
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
		level.Error(logger).Log("msg", "could not set up scheduling", "err", "scheduling option not supported")
		os.Exit(1)
	}

	job.At(config.ScanTime).Do(config.OscapVulnerabilityScan, logger)
	<-gocron.Start()
}

func createLogger(debugLevel string) log.Logger {
	allowLevel := &oscapLogger.AllowedLevel{}
	if err := allowLevel.Set(debugLevel); err != nil {
		fmt.Printf("%v", err)
	}
	return oscapLogger.New(allowLevel)

}
