package config

type ProducerConf struct {
	URI   string `yaml:"uri"`
	Queue string `yaml:"queue"`
}
