package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Version string `yaml:"version"`
	Port    int    `yaml:"port"`
}

func GetConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}
