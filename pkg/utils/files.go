package utils

import (
	"os"
)

// GetFileSize returns the size of the file at the given path in bytes.
// If the file does not exist or an error occurs, it returns 0.
func GetFileSize(f string) uint64 {
	info, err := os.Stat(f)
	if err != nil {
		return 0
	}
	return uint64(info.Size())
}
