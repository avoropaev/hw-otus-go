package config

import (
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v3"
)

type CalendarConfig struct {
	Logger LoggerConf `yaml:"logger"`
	DB     DBConf     `yaml:"db"`
	HTTP   HTTPConf   `yaml:"http"`
	GRPC   GRPCConf   `yaml:"grpc"`
}

type HTTPConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type GRPCConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

func NewCalendarConfig() *CalendarConfig {
	return &CalendarConfig{
		Logger: LoggerConf{
			Level: "info",
		},
	}
}

func ParseCalendarConfig(filePath string) (*CalendarConfig, error) {
	c := NewCalendarConfig()

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
