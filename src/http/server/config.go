package server

type Config struct {
	Host    string         `mapstructure:"host" json:"host" yaml:"host" binding:"required"`
	Port    string         `mapstructure:"port" json:"port" yaml:"port" binding:"required"`
	Options map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}
