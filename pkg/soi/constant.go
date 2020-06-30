package soi

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

func SoisDirPath() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".soi"), nil
}
