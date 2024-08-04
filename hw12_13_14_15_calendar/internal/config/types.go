package config

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

type DBConf struct {
	Database   string `mapstructure:"database"`
	IsPostgres bool   `mapstructure:"isPostgres"`
}

type PostgresCnf struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type HttpServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type GrpcServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}

type SwaggServerConf struct {
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
}
