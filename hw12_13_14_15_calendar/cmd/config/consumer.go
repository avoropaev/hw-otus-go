package config

type ConsumerConf struct {
	ConsumerTag  string `mapstructure:"consumer_tag"`
	URI          string `mapstructure:"uri"`
	ExchangeName string `mapstructure:"exchange_name"`
	ExchangeType string `mapstructure:"exchange_type"`
	Queue        string `mapstructure:"queue"`
	BindingKey   string `mapstructure:"binding_key"`
	Threads      int    `mapstructure:"threads"`
}
