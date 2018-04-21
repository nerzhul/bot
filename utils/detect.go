package utils

import "os"

// IsInDocker detect if we are running in a docker container
func IsInDocker() bool {
	if _, err := os.Stat("/.dockerenv"); os.IsNotExist(err) {
		return false
	}

	return true
}
