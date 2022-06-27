package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type SchedulerConfig struct {
	Logger   LoggerConf   `yaml:"logger"`
	DB       DBConf       `yaml:"db"`
	Producer ProducerConf `yaml:"producer"`
}

func NewSchedulerConfig() *SchedulerConfig {
	return &SchedulerConfig{
		Logger: LoggerConf{
			Level: "info",
		},
	}
}

func ParseSchedulerConfig(filePath string) (*SchedulerConfig, error) {
	c := NewSchedulerConfig()

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
