package scheduler

import (
	"fmt"

	"github.com/spf13/viper"
)

type ShedullerConfig struct {
	Interval         string `mapstructure:"interval"`
	RabbitConfigPath string `mapstructure:"rabbitmq_path"`
}

func NewShedullerConfig(path string) (ShedullerConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config_scheduler")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetDefault("interval", "00:01:00.000000")
	viper.SetDefault("config_scheduler", "../../")

	err := viper.ReadInConfig()
	if err != nil {
		return ShedullerConfig{}, fmt.Errorf("failed to read config_scheduler: %w", err)
	}

	var config ShedullerConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		return ShedullerConfig{}, fmt.Errorf("failed to unmarshal config_scheduler: %w", err)
	}

	return config, nil
}
