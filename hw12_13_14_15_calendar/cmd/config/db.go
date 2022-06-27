package config

type DBConf struct {
	Type string   `yaml:"type"`
	PSQL PSQLConf `yaml:"psql"`
}

type PSQLConf struct {
	URL string `yaml:"url"`
}
