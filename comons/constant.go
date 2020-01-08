package commons

import "github.com/mitchellh/go-homedir"

const (
//SoisFilePath = "sois.json"
)

func SoisFilePath() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return dir + "/sois.json", nil
}
