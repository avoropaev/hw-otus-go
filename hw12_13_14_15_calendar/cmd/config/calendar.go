package config

import (
	"github.com/spf13/viper"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/pkg/viperenvreplacer"
)

type CalendarConfig struct {
	Logger LoggerConf `mapstructure:"logger"`
	DB     DBConf     `mapstructure:"db"`
	HTTP   HTTPConf   `mapstructure:"http"`
	GRPC   GRPCConf   `mapstructure:"grpc"`
}

type HTTPConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type GRPCConf struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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

	viper.SetConfigFile(filePath)

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	viperenvreplacer.ViperReplaceEnvs()

	err = viper.Unmarshal(&c)
	if err != nil {
		return nil, err
	}

	return c, nil
}
