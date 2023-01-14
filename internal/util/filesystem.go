package util

import (
	"os"
)

func IsFileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}
