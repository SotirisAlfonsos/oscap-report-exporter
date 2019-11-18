package main

import (
	"flag"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/stretchr/testify/assert"
	"os"
	"oscap-report-exporter/oscaplogger"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	if testing.Short() {
		os.Exit(m.Run())
	}

	// configFile := "example/oscap-config.yaml"
	// config := oscap.GetConfig(configFile)
	// config.CleanFiles = false

	// config.OscapVulnerabilityScan()

	// log.Printf("Verify that report and downloaded files exist")
	// if !fileExists(config.WorkingFolder+"results.xml") || !fileExists(config.WorkingFolder+config.FileName) {
	// 	log.Fatalf("One of the files we expected does not exist. Fail the tests")
	// }

	// errRemoveDownload := os.Remove(config.WorkingFolder + config.FileName)
	// if errRemoveDownload != nil {
	// 	log.Fatal("Unable to remove " + config.FileName + " with error " + fmt.Sprint(errRemoveDownload))
	// }
	// errRemoveResults := os.Remove(config.WorkingFolder + "results.xml")
	// if errRemoveResults != nil {
	// 	log.Fatal("Unable to remove results.xml with error " + fmt.Sprint(errRemoveResults))
	// }

	exitCode := m.Run()
	os.Exit(exitCode)

}

func TestStartSchedulerWrongSchedulingOption(t *testing.T) {
	_, err := createJob("Monday")
	assert.Error(t, err)
}

func TestStartSchedulerCorrectSchedulingOption(t *testing.T) {
	_, err := createJob("Mon")
	assert.NoError(t, err)
}

func getLogger() log.Logger {
	allowLevel := &oscaplogger.AllowedLevel{}
	if err := allowLevel.Set("debug"); err != nil {
		fmt.Printf("%v", err)
	}
	return oscaplogger.New(allowLevel)
}
