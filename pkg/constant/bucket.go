package constant

import (
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// SoisDir はSoiのルートディレクトリを取得する
func SoisDir() (string, error) {
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".soi"), nil
}
