package devereux

import (
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

var CONFIG_PATH = path.Join(BASE_PATH, "config.yaml")

type Config struct {
	CurrentRepository string `yaml:"current_repository"`
}

func (c *Config) Save(configPath string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	f.Write(data)
	return nil
}

func loadConfig(configPath string) (*Config, error) {
	if exists, _ := fileExists(configPath); !exists {
		return &Config{}, nil
	}

	configData, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config *Config
	yaml.Unmarshal(configData, &config)
	return config, nil
}
