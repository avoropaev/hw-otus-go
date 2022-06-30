package config

type DBConf struct {
	Type string   `mapstructure:"type"`
	PSQL PSQLConf `mapstructure:"psql"`
}

type PSQLConf struct {
	URL string `mapstructure:"url"`
}
