package file

import (
	"errors"
	"os"
)

func Check(name string) bool {
	if _, err := os.Stat(name); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
