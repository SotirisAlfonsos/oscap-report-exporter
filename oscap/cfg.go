package oscap

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"oscap-report-exporter/notify"
)

var (
	// DefaultConfig provides a default configuration for oscap exporter
	DefaultConfig = Config{
		ScanDate:                "Sun",
		ScanTime:                "23:00",
		WorkingFolder:           "/tmp/downloads/",
		FileName:                "com.redhat.rhsa-all.ds.xml",
		VulnerabilityReportConf: DefaultVulnerabilityReportConf,
		CleanFiles:              true,
	}

	// DefaultVulnerabilityReportConf provides a default configuration for
	DefaultVulnerabilityReportConf = VulnerabilityReport{
		GlobalVulnerabilityReportHTTPSLocation: "https://www.redhat.com/security/data/metrics/ds/com.redhat.rhsa-all.ds.xml",
	}

	resultsFile                   = "results.xml"
	reportFile                    = "report.html"
	defaultPermission os.FileMode = 0744
)

// Config contains the configuration from the oscap config file
type Config struct {
	ScanDate                string              `yaml:"scan_date"`
	ScanTime                string              `yaml:"scan_time"`
	WorkingFolder           string              `yaml:"working_folder"`
	FileName                string              `yaml:"global_vulnerability_file_name"`
	VulnerabilityReportConf VulnerabilityReport `yaml:"vulnerability_report"`
	Webhook                 string              `yaml:"webhook,omitempty"`
	Profile                 string              `yaml:"profile,omitempty"`
	CleanFiles              bool                `yaml:"clean_files"`
	EmailConfiguration      *notify.EmailConf   `yaml:"email_config,omitempty"`
}

// GetConfig unmarshars the received conf file to the config struct
func GetConfig(configFile string, logger log.Logger) Config {
	var conf Config
	conf.unmarshalConfFromFile(configFile, logger)

	return conf
}

func (conf *Config) unmarshalConfFromFile(file string, logger log.Logger) {
	*conf = DefaultConfig
	if file != "" {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			level.Error(logger).Log("msg", "could not read yml", "err", err)
			os.Exit(1)
		}

		if err = yaml.Unmarshal(yamlFile, conf); err != nil {
			level.Error(logger).Log("msg", "could not unmarshal yml", "err", err)
			os.Exit(1)
		}
	}
}

// OscapVulnerabilityScan is the main function that handles the scan and forwarding of all reports
func (conf *Config) OscapVulnerabilityScan(logger log.Logger) {

	if code := createDir(conf.WorkingFolder, defaultPermission, logger); code != 0 {
		os.Exit(code)
	}

	if code := conf.prepareAndRunScan(logger); code != 0 {
		os.Exit(code)
	}

	if err := conf.sendResultsToChannels(logger); err != nil {
		level.Error(logger).Log("err", err)
	}

	if conf.CleanFiles {
		filesToClean := []string{resultsFile, reportFile, conf.FileName}
		conf.cleanFiles(filesToClean, logger)
	}
}

func (conf *Config) prepareAndRunScan(logger log.Logger) int {

	vulnerabilityReport := conf.VulnerabilityReportConf
	if errDownload := vulnerabilityReport.DownloadFile(conf.WorkingFolder+conf.FileName, logger); errDownload != nil {
		level.Error(logger).Log("msg", "file download failed", "err", errDownload)
		return 1
	}

	level.Info(logger).Log("msg", "starting scan")

	oscan := &OScan{logger, conf.WorkingFolder, resultsFile, reportFile, conf.FileName, conf.Profile}
	if errOscapScan := oscan.RunOscapScan(); errOscapScan != nil {
		level.Error(logger).Log("msg", "cound not run oscap scan in working folder "+conf.WorkingFolder, "err", errOscapScan)
		return 1
	}

	level.Info(logger).Log("msg", "scan completed")
	return 0
}

func (conf *Config) sendResultsToChannels(logger log.Logger) error {

	errWebhook := make(chan error)
	errEmail := make(chan error)

	level.Info(logger).Log("msg", "sending results to channels")

	go func() {
		if conf.Webhook != "" {
			fs := notify.NewFileSender(logger, conf.WorkingFolder, resultsFile, conf.Webhook)
			err := fs.SendFileToWebhook()
			errWebhook <- err
		} else {
			level.Debug(logger).Log("msg", "no webhook configuration")
			errWebhook <- nil
		}
	}()

	go func() {
		if conf.EmailConfiguration != nil {
			err := conf.EmailConfiguration.SendFileViaEmail(conf.WorkingFolder+reportFile, logger)
			errEmail <- err
		} else {
			level.Debug(logger).Log("msg", "no email configuration")
			errEmail <- nil
		}
	}()

	errW := <-errWebhook
	if errW != nil {
		level.Warn(logger).Log("msg", "could not send report file via webhook", "err", errW)
	}
	errE := <-errEmail
	if errE != nil {
		level.Warn(logger).Log("msg", "could not send report file via e-mail", "err", errE)
	}

	if errW != nil || errE != nil {
		return errors.New("Could not send results to all available channels")
	}
	level.Info(logger).Log("msg", "results send to available channels")
	return nil
}

func (conf *Config) cleanFiles(filesToClean []string, logger log.Logger) {
	for _, fileName := range filesToClean {
		err := os.Remove(conf.WorkingFolder + fileName)
		if err != nil {
			level.Error(logger).Log("msg", "unable to remove "+fileName, "err", err)
		}
		level.Debug(logger).Log("msg", "Removed file "+fileName)
	}
}

func createDir(dir string, permission os.FileMode, logger log.Logger) int {
	err := os.MkdirAll(dir, permission)
	if err != nil {
		level.Error(logger).Log("msg", "could not create dir "+dir, "err", err)
		return 1
	}
	return 0
}
