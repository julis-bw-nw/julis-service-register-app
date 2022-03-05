package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	_ "embed"

	"gopkg.in/yaml.v2"
)

//go:embed config.default.yml
var defaultConfig []byte

type Config struct {
	API      API      `yaml:"api"`
	Database Database `yaml:"database"`
	LLDAP    LLDAP    `yaml:"lldap"`
}

type API struct {
	Bind string `yaml:"bind"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     uint16 `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

func (cfg Database) DSN() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		cfg.Host, cfg.Port, cfg.Database, cfg.Username, cfg.Password)
}

type LLDAP struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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
