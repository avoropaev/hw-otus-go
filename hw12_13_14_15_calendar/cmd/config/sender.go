package config

import (
	"github.com/spf13/viper"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/pkg/viperenvreplacer"
)

type SenderConfig struct {
	Logger   LoggerConf   `mapstructure:"logger"`
	DB       DBConf       `mapstructure:"db"`
	Consumer ConsumerConf `mapstructure:"consumer"`
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
