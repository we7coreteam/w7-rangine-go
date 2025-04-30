package server

type Config struct {
	Host        string         `mapstructure:"host" json:"host" yaml:"host"`
	Port        string         `mapstructure:"port" json:"port" yaml:"port" binding:"required"`
	MaxBodySize int64          `mapstructure:"max_body_size" json:"max_body_size" yaml:"max_body_size"`
	Options     map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}
