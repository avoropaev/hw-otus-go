package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type SenderConfig struct {
	Logger   LoggerConf   `yaml:"logger"`
	DB       DBConf       `yaml:"db"`
	Consumer ConsumerConf `yaml:"consumer"`
}

func NewSenderConfig() *SenderConfig {
	return &SenderConfig{
		Logger: LoggerConf{
			Level: "info",
		},
	}
}

func ParseSenderConfig(filePath string) (*SenderConfig, error) {
	c := NewSenderConfig()

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0o644)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}

	return c, nil
}
