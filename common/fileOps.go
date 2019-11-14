package common

import (
	"github.com/pkg/errors"
	"os"
)

// FileExists verifies that the results file does exist
func FileExists(fileName string) error {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return errors.Wrap(err, "File "+fileName+" does not exist")
	} else if info.IsDir() {
		return errors.Wrap(err, fileName+" is a directory")
	} else {
		return nil
	}
}
