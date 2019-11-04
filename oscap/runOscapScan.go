package oscap

import (
	"log"
	"fmt"
	"os"
	"bytes"
    "os/exec"
)

func RunOscapScan(workingFolder string, resultsFile string, fileName string) {

	err := runScan(workingFolder, resultsFile, fileName)
	if err != nil {
		log.Fatal("Error during oscap scan " + fmt.Sprint(err))
	}

	if !fileExists(workingFolder + resultsFile) {
		log.Fatalf("File " + workingFolder + resultsFile + " does not exist (Or is a directory)")
	}
}

// Run oscap scan and store the results in the working folder
func runScan(workingFolder string, resultsFile string, fileName string) error{
    	
	oscapCommand := "oscap xccdf eval --results " + resultsFile + " " + fileName + " > log.out"
	cmd := exec.Command("bash", "-c", oscapCommand)
	cmd.Dir = workingFolder
	log.Printf("Running Oscap command " + cmd.String())

	var stderr bytes.Buffer
    cmd.Stderr = &stderr

	/** oscap returns 0 if all rules pass. 
	 If there is an error during evaluation, the return code is 1. 
	 If there is at least one rule with either fail or unknown result, oscap-scan finishes with return code 2. **/
    err := cmd.Run()
    if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
        	switch exitError.ExitCode() {
        		case 1:
        			log.Fatalf("Error: During oscap evaluation %v", err)
        		case 2:
        			log.Printf("At least one of the rules failed.")
        	}
    	}else {
    		log.Fatalf("Unexpected error %v", err)
    	}
	}

	log.Printf("Results for oscap scan created in folder " + workingFolder)
	return nil
}

// Verify that the results file does exist
func fileExists(fileName string) bool{
	info, err := os.Stat(fileName)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}