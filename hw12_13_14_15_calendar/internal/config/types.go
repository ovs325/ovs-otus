package config

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

type DbConf struct {
	Host       string `mapstructure:"host"`
	Port       string `mapstructure:"port"`
	User       string `mapstructure:"user"`
	Password   string `mapstructure:"password"`
	Database   string `mapstructure:"database"`
	IsPostgres bool   `mapstructure:"isPostgres"`
}

type ReindexerCnf struct {
	Namespace string `mapstructure:"namespace"`
}

type ServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
