package server

type Config struct {
	Host        string         `mapstructure:"host" json:"host" yaml:"host"`
	Port        string         `mapstructure:"port" json:"port" yaml:"port"`
	MaxBodySize int64          `mapstructure:"max_body_size" json:"max_body_size" yaml:"max_body_size"`
	TLS         TLSConfig      `mapstructure:"tls" json:"tls" yaml:"tls"`
	Options     map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}

type TLSConfig struct {
	Enable   bool   `mapstructure:"enable" json:"enable" yaml:"enable"`
	Port     string `mapstructure:"port" json:"port" yaml:"port"`
	CertFile string `mapstructure:"cert_file" json:"cert_file" yaml:"cert_file"`
	KeyFile  string `mapstructure:"key_file" json:"key_file" yaml:"key_file"`
}
