package oscap

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPrepareOscapCommandWithProfile(t *testing.T) {
	oscan := &OScan{logger, "", "results.xml", "report.html", "fileName.xml", "profile.rhel"}
	result := "oscap xccdf eval --profile profile.rhel --results results.xml --report report.html fileName.xml > /dev/null"
	oscapCmd := oscan.prepareOscapCommand()
	assert.Equal(t, oscapCmd, result)
}

func TestPrepareOscapCommandNoProfile(t *testing.T) {
	oscan := &OScan{logger, "", "results.xml", "report.html", "fileName.xml", ""}
	result := "oscap xccdf eval --results results.xml --report report.html fileName.xml > /dev/null"
	oscapCmd := oscan.prepareOscapCommand()
	assert.Equal(t, oscapCmd, result)
}
