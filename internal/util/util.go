package util

import (
	"os"
	"path/filepath"
)

func CurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	_, cd := filepath.Split(dir)
	return cd
}
