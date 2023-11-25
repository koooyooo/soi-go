package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/caarlos0/env/v10"

	"soi-go/pkg/common/file"

	"github.com/mitchellh/go-homedir"
)

// 設定情報
type Config struct {
	Theme             string `env:"SOI_THEME" json:"theme"`
	DefaultBrowser    string `env:"SOI_DEFAULT_BROWSER" json:"default_browser"`
	DefaultBucket     string `env:"SOI_DEFAULT_BUCKET" json:"default_bucket"`
	Server            string `env:"SOI_SERVER" json:"server"`
	DefaultRepository string `env:"SOI_DEFAULT_REPOSITORY" json:"default_repository"`
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
	//sc := bufio.NewScanner(os.Stdin)
	//fmt.Println(`server url? (ex. https://server:80")`)
	//fmt.Print("> ")
	//sc.Scan()
	//txt := sc.Text()
	//if txt == "" {
	//	return errors.New("no servername specified")
	//}
	txt := ""
	cfg := Config{
		Server: txt,
	}
	b, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, b, 0644)
}

func doLoad(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := json.Unmarshal(b, &conf); err != nil {
		return nil, err
	}
	if err := env.Parse(&conf); err != nil {
		return nil, err
	}
	if conf.DefaultBucket == "" {
		conf.DefaultBucket = "default"
	}
	if conf.DefaultRepository == "" {
		conf.DefaultRepository = "file"
	}
	return &conf, nil
}

func confPath() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	confDir := filepath.Join(dir, ".soi")
	if !file.Exists(confDir) {
		if err := os.Mkdir(confDir, 0755); err != nil {
			return "", err
		}
	}
	return filepath.Join(confDir, "config.json"), nil
}
