package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConf   `mapstructure:"server"`
	Logger LoggerConf   `mapstructure:"logger"`
	DB     DBConf       `mapstructure:"database"`
	RxCnf  ReindexerCnf `mapstructure:"reindexer"`
	PgCnf  PostgresCnf  `mapstructure:"postgres"`
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

	viper.SetDefault("reindexer.host", "localhost")
	viper.SetDefault("reindexer.port", 7778)
	viper.SetDefault("reindexer.namespace", "events")

	viper.SetDefault("database.isPostgres", true)
	viper.SetDefault("database.database", "calendar")

	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.port", 7779)

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
