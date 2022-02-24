package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"os"
	"text/template"

	_ "embed"

	"gopkg.in/yaml.v2"
)

//go:embed config.default.yml
var defaultConfig string
var tmpl = template.New("defaultConfig")

func init() {
	tmpl = template.Must(tmpl.Parse(defaultConfig))
}

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

type ConfigData struct {
	Secret string
}

func generateDefaultData() (ConfigData, error) {
	secret := make([]byte, 16)
	if _, err := rand.Read(secret); err != nil {
		return ConfigData{}, err
	}

	return ConfigData{
		Secret: hex.EncodeToString(secret),
	}, nil
}

func createConfigIfNotExist(path string) error {
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := generateDefaultData()
	if err != nil {
		return err
	}

	return tmpl.Execute(f, d)
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
