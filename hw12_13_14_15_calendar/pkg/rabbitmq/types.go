package rabbitmq

import "fmt"

type RabbitConf struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Queue    string `mapstructure:"queue"`
}

func (r *RabbitConf) GetRabbitDSN() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s:%v/",
		r.User,
		r.Password,
		r.Host,
		r.Port,
	)
}
