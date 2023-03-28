package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/koooyooo/soi-go/pkg/common/file"
	"github.com/mitchellh/go-homedir"
)

// 設定情報
type Config struct {
	Server        string `json:"server"`
	DefaultBucket string `json:"default_bucket"`
}

func Load() (*Config, error) {
	path, err := confPath()
	if err != nil {
		return nil, err
	}
	exists, err := exists(path)
	if err != nil {
		return nil, fmt.Errorf("fail in checking conf file existence: %v", err)
	}
	if !exists {
		if err := initialize(path); err != nil {
			return nil, fmt.Errorf("fail in initializing conf file: %v", err)
		}
	}
	cfg, err := doLoad(path)
	if err != nil {
		fmt.Errorf("fail in loading path: %v", err)
	}
	return cfg, err
}

func exists(path string) (bool, error) {
	return file.Exists(path), nil
}

func initialize(path string) error {
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

func doLoad(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := json.Unmarshal(b, &conf); err != nil {
		return nil, err
	}
	if conf.DefaultBucket == "" {
		conf.DefaultBucket = "default"
	}
	return &conf, nil
}

func confPath() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	confPath := filepath.Join(dir, ".soi", "config.json")
	if !file.Exists(confPath) {
		return "", nil
	}
	return confPath, nil
}
