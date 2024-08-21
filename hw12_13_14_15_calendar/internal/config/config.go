package config

import (
	"fmt"

	rb "github.com/ovs325/ovs-otus/hw12_13_14_15_calendar/pkg/rabbitmq"
	"github.com/spf13/viper"
)

type CalendarConfig struct {
	HTTPServer  HTTPServerConf  `mapstructure:"http_server"`
	GrpcServer  GrpcServerConf  `mapstructure:"grpc_server"`
	SwaggServer SwaggServerConf `mapstructure:"swagger_server"`
	Logger      LoggerConf      `mapstructure:"logger"`
	DB          DBConf          `mapstructure:"database"`
	PgCnf       PostgresCnf     `mapstructure:"postgres"`
	Rabbit      rb.RabbitConf   `mapstructure:"rabbitmq"`
}

func NewCalendarConfig(path string) (CalendarConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config_calendar")
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

	viper.SetDefault("rabbitmq.host", "localhost")
	viper.SetDefault("rabbitmq.port", 5672)
	viper.SetDefault("rabbitmq.user", "guest")
	viper.SetDefault("rabbitmq.password", "guest")
	viper.SetDefault("rabbitmq.queue", "notifications")

	err := viper.ReadInConfig()
	if err != nil {
		return CalendarConfig{}, fmt.Errorf("failed to read calendar_config: %w", err)
	}

	var config CalendarConfig
	err = viper.Unmarshal(&config)
	if err != nil {
		return CalendarConfig{}, fmt.Errorf("failed to unmarshal calendar_config: %w", err)
	}

	return config, nil
}
