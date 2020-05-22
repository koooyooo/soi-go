package common

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

const (
	SoisJSONName = "sois.json"
)

func SoisFilePath() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, SoisJSONName), nil
}
