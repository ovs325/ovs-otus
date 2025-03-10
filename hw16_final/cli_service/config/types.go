package config

type GroupParams struct {
	LenQueue    int64   `mapstructure:"lenQueue"`
	Capacity    int64   `mapstructure:"capacity"`
	LeakageRate float64 `mapstructure:"leakageRate"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

type HTTPServerConf struct {
	Scheme string `mapstructure:"scheme"`
	Host   string `mapstructure:"host"`
	Port   string `mapstructure:"port"`
}
