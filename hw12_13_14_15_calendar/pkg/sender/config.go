package sender

import (
	"fmt"

	"github.com/spf13/viper"
)

type SenderConfig struct {
	RabbitConfigPath string `mapstructure:"rabbitmq_path"`
}

func NewSenderConfig(path string) (SenderConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config_sender")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetDefault("config_sender", "../../")

	err := viper.ReadInConfig()
	if err != nil {
		return SenderConfig{}, fmt.Errorf("failed to read config_sender: %w", err)
	}

	var config SenderConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		return SenderConfig{}, fmt.Errorf("failed to unmarshal config_sender: %w", err)
	}

	return config, nil
}
