package oscap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepareOscapCommandWithProfile(t *testing.T) {
	oscan := &OScan{logger, "", "results.xml", "report.html", "fileName.xml", "profile.rhel", "xccdf"}
	result := "oscap xccdf eval --profile profile.rhel --results results.xml --report report.html fileName.xml > /dev/null"
	oscapCmd, _ := oscan.prepareOscapCommand()
	assert.Equal(t, oscapCmd, result)
}

func TestPrepareOscapCommandNoProfile(t *testing.T) {
	oscan := &OScan{logger, "", "results.xml", "report.html", "fileName.xml", "", "oval"}
	result := "oscap oval eval --results results.xml --report report.html fileName.xml > /dev/null"
	oscapCmd, _ := oscan.prepareOscapCommand()
	assert.Equal(t, oscapCmd, result)
}

func TestPrepareOscapCommandInvalidModule(t *testing.T) {
	oscan := &OScan{logger, "", "results.xml", "report.html", "fileName.xml", "", "invalid"}
	result := ""
	oscapCmd, err := oscan.prepareOscapCommand()
	assert.Equal(t, oscapCmd, result)
	assert.Error(t, err)
}
