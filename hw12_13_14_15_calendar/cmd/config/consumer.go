package config

type ConsumerConf struct {
	ConsumerTag  string `yaml:"consumer_tag"`
	URI          string `yaml:"uri"`
	ExchangeName string `yaml:"exchange_name"`
	ExchangeType string `yaml:"exchange_type"`
	Queue        string `yaml:"queue"`
	BindingKey   string `yaml:"binding_key"`
	Threads      int    `yaml:"threads"`
}
