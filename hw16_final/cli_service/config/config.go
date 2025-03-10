package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	HTTPServer HTTPServerConf `mapstructure:"http_server"`
	Logger     LoggerConf     `mapstructure:"logger"`
	IsTest     bool           `mapstructure:"is_test"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.AddConfigPath("./")
	viper.AddConfigPath("../../")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetDefault("is_test", false)
	viper.SetDefault("http_server.scheme", "http")
	viper.SetDefault("http_server.host", "127.0.0.1")
	viper.SetDefault("http_server.port", "3009")

	if err = viper.ReadInConfig(); err != nil {
		return
	}
	err = viper.Unmarshal(&config)
	return
}
