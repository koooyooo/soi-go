package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Server string `json:"server"`
}

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
