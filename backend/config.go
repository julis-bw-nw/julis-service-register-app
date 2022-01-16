package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

//go:embed config.default.yaml
var defaultConfig []byte

type Config struct {
	API      APIConfig      `yaml:"api"`
	Database DatabaseConfig `yaml:"database"`
}

type APIConfig struct {
	Bind string `yaml:"bind"`
}
type DatabaseConfig struct {
	Host             string `yaml:"host"`
	Database         string `yaml:"database"`
	Username         string `yaml:"username"`
	Password         string `yaml:"password"`
	EncryptionSecret string `yaml:"encryption_secret"`
}

func createConfigIfNotExist(path string) error {
	if _, err := os.Stat(path); err == nil {
		return err
	}

	return ioutil.WriteFile(path, defaultConfig, 0644)
}

func loadConfig(path string) (Config, error) {
	bb, err := ioutil.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(bb, &cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
