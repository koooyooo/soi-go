package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	fileio2 "github.com/koooyooo/soi-go/pkg/common/file"

	"github.com/mitchellh/go-homedir"
)

// 設定情報
type Config struct {
	Server string `json:"server"`
}

// LoadConfig は設定をロードします
func LoadConfig() (*Config, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	confPath := filepath.Join(dir, ".soi", "config.json")
	if !fileio2.Exists(confPath) {
		initConfig(confPath)
	}
	b, err := ioutil.ReadFile(confPath)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := json.Unmarshal(b, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}

func initConfig(path string) error {
	sc := bufio.NewScanner(os.Stdin)
	fmt.Println(`server url? (ex. https://server:80")`)
	fmt.Print("> ")
	sc.Scan()
	txt := sc.Text()
	if txt == "" {
		return errors.New("no servername specified")
	}
	cfg := Config{
		Server: txt,
	}
	b, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, b, 0600)
}
