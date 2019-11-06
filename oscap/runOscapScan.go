package oscap

import (
	"log"
	"bytes"
    "os/exec"
)

func RunOscapScan(workingFolder string, resultsFile string, fileName string) error {

	errScan := runScan(workingFolder, resultsFile, fileName)
	if errScan != nil {
		return errScan
	}
	
	errFileExists := fileExists(workingFolder + resultsFile)
	if errFileExists != nil {
		return errFileExists
	}

	return nil
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
        			log.Printf("Error: During oscap evaluation")
        			return err
        		case 2:
        			log.Printf("At least one of the rules failed.")
        	}
    	}else {
    		log.Printf("Error: Unexpected error during oscap evaluation")
    		return err
    	}
	}

	log.Printf("Results for oscap scan created in folder " + workingFolder)
	return nil
}