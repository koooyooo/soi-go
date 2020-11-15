package fileio

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const (
	soisJSONName = "sois.json"
)

// SoisFilePath return default JSON path
func SoisFilePath() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, soisJSONName), nil
}

func SoisDirPath(bucket string) (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".soi", bucket), nil
}
