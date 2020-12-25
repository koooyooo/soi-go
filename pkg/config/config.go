package config

import (
	"encoding/json"
	"io/ioutil"
)

// 設定情報
type Config struct {
	Server string `json:"server"`
}

// 設定をロードします
func LoadConfig() (*Config, error) {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var conf Config
	if err := json.Unmarshal(b, &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
