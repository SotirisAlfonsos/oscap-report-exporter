package oscap

import (
	"log"
    "os"
    "fmt"
	"gopkg.in/yaml.v2"
    "io/ioutil"
)

var (
	DefaultConfig = config {
		ScanDate: "Sun",
		ScanTime: "23:00",
		WorkingFolder: "/tmp/downloads/",
		FileName: "com.redhat.rhsa-all.ds.xml",
		VulnerabilityReportConf: DefaultVulnerabilityReportConf,
		//NetworkRetry: 3,
		Webhook: "http://localhost:8080",
		CleanFiles: true,
	}

	DefaultVulnerabilityReportConf = VulnerabilityReport {
		GlobalVulnerabilityReportHttpsLocation: "security/data/metrics/ds/com.redhat.rhsa-all.ds.xml",
		BaseVulnerabilityReportUrl: "https://www.redhat.com/",
	}

	resultsFile = "results.xml"
	reportFile = "report.html"
	defaultPermission os.FileMode = 0744

)

type config struct {
	ScanDate string `yaml:"scan_date"`
	ScanTime string `yaml:"scan_time"`
	WorkingFolder string `yaml:"working_folder"`
	FileName string `yaml:"global_vulnerability_file_name"`
	VulnerabilityReportConf VulnerabilityReport `yaml:"vulnerability_report"`
	//NetworkRetry int `yaml:"network_retry"`
	Webhook string `yaml:"webhook"`
	CleanFiles bool `yaml:"clean_files"`
	EmailConfiguration *EmailConf `yaml:"email_config,omitempty"`
}

func GetConfig(configFile string) config {
	var conf config
	conf.unmarshalConfFromFile(configFile)

    return conf
}

func (conf *config) unmarshalConfFromFile(file string) {
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

func (conf *config) OscapVulnerabilityScan() {

	createDir(conf.WorkingFolder, defaultPermission)

	vulnerabilityReport := conf.VulnerabilityReportConf
	errDownload := vulnerabilityReport.DownloadFile(conf.WorkingFolder + conf.FileName)
	if errDownload != nil {
		log.Fatalf("Error: File download failed : %v", errDownload)
	}

	
	errOscapScan := RunOscapScan(conf.WorkingFolder, resultsFile, conf.FileName)
	if errOscapScan != nil {
		log.Fatalf("Error: Cound not run oscap scan in working folder " + conf.WorkingFolder + " : " + fmt.Sprint(errOscapScan))
	}

	errWebhook := SendFileToWebhook(conf.WorkingFolder + resultsFile, conf.Webhook)
	if errWebhook != nil {
		log.Fatalf("Error: sending xml to webhook " + conf.Webhook + " : " + fmt.Sprint(errWebhook))
	}

	if conf.EmailConfiguration != nil {
		errEmail := conf.EmailConfiguration.SendFileViaEmail(conf.WorkingFolder + reportFile)
		if errEmail != nil {
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
		log.Fatalf("Eror: Could not create Dir " + dir + " : %v ", err)
	}
}

func (conf *config) cleanFiles(filesToClean []string) {
	for _, fileName := range filesToClean {
		err := os.Remove(conf.WorkingFolder + fileName)
		if err != nil {
			log.Fatal("Error: Unable to remove " + fileName + " with error " + fmt.Sprint(err))
		}
		log.Printf("Removed file " + fileName)
	}
}

// Verify that the results file does exist
func fileExists(fileName string) error{
	info, err := os.Stat(fileName)
    if os.IsNotExist(err) {
        log.Printf("Error: File " + fileName + " does not exist")
        return err
    }else if info.IsDir() {
    	log.Printf("Error: " + fileName + " is a directory")
    	return err
    }else {
    	return nil
    }

}
