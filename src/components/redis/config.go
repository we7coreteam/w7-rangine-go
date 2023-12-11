package redis

type Config struct {
	Host     string         `mapstructure:"host" json:"host" yaml:"host" binding:"required"`
	Port     uint           `mapstructure:"port" json:"port" yaml:"port" binding:"required"`
	Password string         `mapstructure:"password" json:"password" yaml:"password"`
	Username string         `mapstructure:"user_name" json:"user_name" yaml:"user_name"`
	Db       uint           `mapstructure:"db" json:"db" yaml:"db"`
	PoolSize uint           `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
	Options  map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}
