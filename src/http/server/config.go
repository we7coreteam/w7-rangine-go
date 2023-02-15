package server

type Config struct {
	Host string `mapstructure:"host" json:"host" yaml:"host"`
	Port int    `mapstructure:"port" json:"port" yaml:"port"`
}
