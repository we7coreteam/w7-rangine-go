package server

type Config struct {
	Host    string         `mapstructure:"host" json:"host" yaml:"host"`
	Port    string         `mapstructure:"port" json:"port" yaml:"port"`
	Options map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}
