package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	Server struct {
		Addr string `yaml:"addr"`
		Port uint   `yaml:"port"`
	} `yaml:"server"`
	Database struct {
		Dsn string `yaml:"dsn"`
	} `yaml:"database"`
	Rabbit struct {
		Dsn      string `yaml:"dsn"`
		Exchange string `yaml:"exchange"`
		Queue    string `yaml:"queue"`
		Tag      string `yaml:"tag"`
	} `yaml:"rabbit"`
}

var ErrUnsupportedType = errors.New("unsupported type")

func NewConfig(configPath string) *Config {
	if filepath.Ext(configPath) != ".yaml" {
		return nil
	}

	file, _ := os.Open(configPath)
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	decoder := yaml.NewDecoder(file)

	config := &Config{}
	if err := decoder.Decode(&config); err != nil {
		return nil
	}

	return config
}
