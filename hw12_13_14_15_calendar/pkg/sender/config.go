package sender

import (
	"fmt"

	"github.com/spf13/viper"
)

type SndConfig struct {
	RabbitConfigPath string `mapstructure:"rabbitmq_path"`
}

func NewSenderConfig(path string) (SndConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config_sender")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()

	viper.SetDefault("config_sender", "../../")

	err := viper.ReadInConfig()
	if err != nil {
		return SndConfig{}, fmt.Errorf("failed to read config_sender: %w", err)
	}

	var config SndConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		return SndConfig{}, fmt.Errorf("failed to unmarshal config_sender: %w", err)
	}

	return config, nil
}
