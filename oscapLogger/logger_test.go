package oscapLogger

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoggerSetLogLevelVerifyTrue(t *testing.T) {
	loglevel := "info"
	allowedlevel := &AllowedLevel{}
	err := allowedlevel.Set(loglevel)

	assert.NoError(t, err)
	assert.Equal(t, allowedlevel.s, loglevel, "The loglevel should be the same")
}

func TestLoggerSetLogLevelVerifyFalse(t *testing.T) {
	loglevel := "test"
	allowedlevel := &AllowedLevel{}
	err := allowedlevel.Set(loglevel)

	assert.Error(t, err)
	assert.EqualError(t, err, "unrecognized log level "+loglevel)
	assert.Equal(t, allowedlevel.s, "", "The loglevel should be the same")
}
