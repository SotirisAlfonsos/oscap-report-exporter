package oscap

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
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
	CleanFiles              bool                `yaml:"clean_files"`
	EmailConfiguration      *EmailConf          `yaml:"email_config,omitempty"`
}

// GetConfig unmarshars the received conf file to the config struct
func GetConfig(configFile string) Config {
	var conf Config
	conf.unmarshalConfFromFile(configFile)

	return conf
}

func (conf *Config) unmarshalConfFromFile(file string) {
	*conf = DefaultConfig
	if file != "" {
		yamlFile, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Error: yamlFile.Get err %v ", err)
		}

		err = yaml.Unmarshal(yamlFile, conf)
		if err != nil {
			log.Fatalf("Error: Unmarshal: %v", err)
		}
	}
}

// OscapVulnerabilityScan is the main function that handles the scan and forwarding of all reports
func (conf *Config) OscapVulnerabilityScan() {

	createDir(conf.WorkingFolder, defaultPermission)

	vulnerabilityReport := conf.VulnerabilityReportConf
	if errDownload := vulnerabilityReport.DownloadFile(conf.WorkingFolder + conf.FileName); errDownload != nil {
		log.Fatalf("Error: File download failed : %v", errDownload)
	}

	if errOscapScan := RunOscapScan(conf.WorkingFolder, resultsFile, reportFile, conf.FileName); errOscapScan != nil {
		log.Fatalf("Error: Cound not run oscap scan in working folder " + conf.WorkingFolder + " : " + fmt.Sprint(errOscapScan))
	}

	if conf.Webhook != "" {
		if errWebhook := SendFileToWebhook(conf.WorkingFolder, resultsFile, conf.Webhook); errWebhook != nil {
			log.Fatalf("Error: sending xml to webhook " + conf.Webhook + " : " + fmt.Sprint(errWebhook))
		}
	}

	if conf.EmailConfiguration != nil {
		if errEmail := conf.EmailConfiguration.SendFileViaEmail(conf.WorkingFolder + reportFile); errEmail != nil {
			log.Fatalf("Error: Could not send report file via Email " + fmt.Sprint(errEmail))
		}
	}

	if conf.CleanFiles {
		filesToClean := []string{resultsFile, reportFile, conf.FileName}
		conf.cleanFiles(filesToClean)
	}
}

func createDir(dir string, permission os.FileMode) {
	err := os.MkdirAll(dir, permission)
	if err != nil {
		log.Fatalf("Eror: Could not create Dir "+dir+" : %v ", err)
	}
}

func (conf *Config) cleanFiles(filesToClean []string) {
	for _, fileName := range filesToClean {
		err := os.Remove(conf.WorkingFolder + fileName)
		if err != nil {
			log.Fatal("Error: Unable to remove " + fileName + " with error " + fmt.Sprint(err))
		}
		log.Printf("Removed file " + fileName)
	}
}

// Verify that the results file does exist
func fileExists(fileName string) error {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		log.Printf("Error: File " + fileName + " does not exist")
		return err
	} else if info.IsDir() {
		log.Printf("Error: " + fileName + " is a directory")
		return err
	} else {
		return nil
	}

}
