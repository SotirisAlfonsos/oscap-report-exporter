package oscap

import (
	"bytes"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"os/exec"
	"oscap-report-exporter/common"
)

// OScan contains some configuration for the execution of the scan
type OScan struct {
	logger        log.Logger
	workingFolder string
	resultsFile   string
	reportFile    string
	fileName      string
}

// RunOscapScan runs the scan on the host machine
func (oscan *OScan) RunOscapScan() error {

	if err := oscan.runScan(); err != nil {
		return err
	}

	if err := common.FileExists(oscan.workingFolder + oscan.resultsFile); err != nil {
		return err
	}

	return nil
}

// Run oscap scan and store the results in the working folder
func (oscan *OScan) runScan() error {

	oscapCommand := "oscap xccdf eval --results " + oscan.resultsFile + " --report " + oscan.reportFile + " " + oscan.fileName + " > log.out"
	cmd := exec.Command("bash", "-c", oscapCommand)
	cmd.Dir = oscan.workingFolder
	level.Info(oscan.logger).Log("msg", "Running Oscap command "+cmd.String())

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
				return errors.Wrap(err, "could not complete oscap evaluation successfully. Exit code was "+string(exitError.ExitCode()))
			case 2:
				level.Info(oscan.logger).Log("msg", "at least one of the rules failed")
			}
		} else {
			return errors.Wrap(err, "unexpected error during oscap evaluation")
		}
	}
	level.Info(oscan.logger).Log("msg", "Results for oscap scan created in folder "+oscan.workingFolder)
	return nil
}
