package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Server ServerConf   `mapstructure:"server"`
	Logger LoggerConf   `mapstructure:"logger"`
	Db     DbConf       `mapstructure:"database"`
	RxCnf  ReindexerCnf `mapstructure:"reindexer"`
}

func NewConfig(path string) (Config, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetDefault("logger.level", "Error")
	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", 5432)
	viper.SetDefault("database.isPostgres", true)

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
