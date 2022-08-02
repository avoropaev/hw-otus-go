package config

import (
	"github.com/spf13/viper"

	"github.com/avoropaev/hw-otus-go/hw12_13_14_15_calendar/pkg/viperenvreplacer"
)

type SchedulerConfig struct {
	Logger   LoggerConf   `mapstructure:"logger"`
	DB       DBConf       `mapstructure:"db"`
	Producer ProducerConf `mapstructure:"producer"`
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
