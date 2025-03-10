package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Params     map[string]GroupParams `mapstructure:"params"`
	Deault     GroupParams            `mapstructure:"default"`
	HTTPServer HTTPServerConf         `mapstructure:"http_server"`
	Logger     LoggerConf             `mapstructure:"logger"`
	IsTest     bool                   `mapstructure:"is_test"`
}

func LoadConfig(path string) (config Config, err error) {
	envPath := os.Getenv("DOCKER_CONF_PATH")
	if envPath != "" {
		path = envPath
	}
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetDefault("default.lenQueue", 10000)
	viper.SetDefault("default.capacity", 10)
	viper.SetDefault("default.leakageRate", 0.5)
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
