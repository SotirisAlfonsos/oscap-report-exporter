package oscap

import (
	"log"
    "os"
    "fmt"
    "bytes"
    "os/exec"
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
		NetworkRetry: 3,
		Webhook: "http://localhost:8080",
	}

	DefaultVulnerabilityReportConf = VulnerabilityReport {
		GlobalVulnerabilityReportHttpsLocation: "security/data/metrics/ds/com.redhat.rhsa-all.ds.xml",
		BaseVulnerabilityReportUrl: "https://www.redhat.com/",
	}

	resultsFile = "results.xml"

)

type config struct {
	ScanDate string `yaml:"scan_date"`
	ScanTime string `yaml:"scan_time"`
	WorkingFolder string `yaml:"working_folder"`
	FileName string `yaml:"global_vulnerability_file_name"`
	VulnerabilityReportConf VulnerabilityReport `yaml:"artifactory"`
	NetworkRetry int `yaml:"network_retry"`
	Webhook string `yaml:"webhook"`
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
	        log.Printf("yamlFile.Get err   #%v ", err)
	    }

	    err = yaml.Unmarshal(yamlFile, conf)
	    if err != nil {
	        log.Fatalf("Unmarshal: %v", err)
	    }
	}
}

func (conf *config) OscapVulnerabilityScan() {
	vulnerabilityReport := conf.VulnerabilityReportConf
	vulnerabilityReport.DownloadFile(conf.WorkingFolder + conf.FileName, conf.NetworkRetry)
	
	err := conf.runOscapScan()
	if err != nil {
		log.Fatal("Error during oscap scan " + fmt.Sprint(err))
	}

	if !fileExists(conf.WorkingFolder+ resultsFile) {
		log.Fatalf("File " + conf.WorkingFolder + resultsFile + " does not exist (Or is a directory)")
	}

	fileSender := FileSender{conf.WorkingFolder+ resultsFile, conf.Webhook}
	errSendFile := fileSender.SendFileToWebhook()
	if errSendFile != nil {
		log.Fatalf("Error: Could not send data to webhook " + fmt.Sprint(errSendFile))
	}

	// filesToClean := []string{resultsFile, conf.FileName}
	// conf.cleanFiles(filesToClean)
	
}

func (conf *config) runOscapScan() error{
    	
	oscapCommand := "oscap xccdf eval --results " + resultsFile + " " + conf.FileName + " > log.out"
	cmd := exec.Command("bash", "-c", oscapCommand)
	cmd.Dir = conf.WorkingFolder
	log.Printf("Running Oscap command " + cmd.String())

	var stderr bytes.Buffer
    cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error during oscap scan " + string(stderr.Bytes()) + " " + fmt.Sprint(err) + ". Ignorring")
		// return err
	}

	log.Printf("Results for oscap scan created in folder " + conf.WorkingFolder)
	return nil
}

func fileExists(fileName string) bool{
	info, err := os.Stat(fileName)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func (conf *config) cleanFiles(filesToClean []string) {
	for _, fileName := range filesToClean {
		err := os.Remove(conf.WorkingFolder + fileName)
		if err != nil {
			log.Fatal("Unable to remove " + fileName + " with error " + fmt.Sprint(err))
		}
		log.Printf("Removed file " + fileName)
	}
}
