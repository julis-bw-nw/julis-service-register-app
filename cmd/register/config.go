package main

import (
	"errors"
	"io/ioutil"
	"os"

	_ "embed"

	"gopkg.in/yaml.v2"
)

//go:embed config.default.yml
var defaultConfig []byte

type Config struct {
	API struct {
		Bind string `yaml:"bind"`
	} `yaml:"api"`
	Database struct {
		Host             string `yaml:"host"`
		Database         string `yaml:"database"`
		Username         string `yaml:"username"`
		Password         string `yaml:"password"`
		EncryptionSecret string `yaml:"encryption_secret"`
	} `yaml:"database"`
}

func createConfigIfNotExist(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return ioutil.WriteFile(configPath, []byte(defaultConfig), 0644)
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
