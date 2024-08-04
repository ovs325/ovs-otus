package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer  HTTPServerConf  `mapstructure:"http_server"`
	GrpcServer  GrpcServerConf  `mapstructure:"grpc_server"`
	SwaggServer SwaggServerConf `mapstructure:"swagger_server"`
	Logger      LoggerConf      `mapstructure:"logger"`
	DB          DBConf          `mapstructure:"database"`
	PgCnf       PostgresCnf     `mapstructure:"postgres"`
}

func NewConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetDefault("logger.level", "Error")

	viper.SetDefault("postgres.host", "localhost")
	viper.SetDefault("postgres.port", 7777)
	viper.SetDefault("postgres.user", "user")
	viper.SetDefault("postgres.password", "pass")

	viper.SetDefault("database.isPostgres", true)
	viper.SetDefault("database.database", "calendar")

	viper.SetDefault("http_server.host", "localhost")
	viper.SetDefault("http_server.port", 7779)

	viper.SetDefault("grpc_server.host", "localhost")
	viper.SetDefault("grpc_server.port", 7780)

	viper.SetDefault("swagger_server.host", "localhost")
	viper.SetDefault("swagger_server.port", 7781)

	err := viper.ReadInConfig()
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config: %w", err)
	}

	var config Config
	err = viper.Unmarshal(&config)
	if err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}
