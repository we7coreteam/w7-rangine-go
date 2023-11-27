package database

type Config struct {
	Driver        string         `mapstructure:"driver" json:"driver" yaml:"driver" binding:"required"`
	Host          string         `mapstructure:"host" json:"host" yaml:"host" binding:"required_if=Driver mysql|required_without=DSN"`
	Port          int            `mapstructure:"port" json:"port" yaml:"port" binding:"required_if=Driver mysql|required_without=DSN"`
	Username      string         `mapstructure:"user_name" json:"user_name" yaml:"user_name" binding:"required_without=DSN"`
	Password      string         `mapstructure:"password" json:"password" yaml:"password" binding:"required_without=DSN"`
	DbName        string         `mapstructure:"db_name" json:"db_name" yaml:"db_name" binding:"required_without=DSN"`
	Charset       string         `mapstructure:"charset" json:"charset" yaml:"charset"`
	Prefix        string         `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	MaxIdleConn   int            `mapstructure:"max_idle_conn" json:"max_idle_conn" yaml:"max_idle_conn"`
	MaxConn       int            `mapstructure:"max_conn" json:"max_conn" yaml:"max_conn"`
	SlowThreshold int64          `mapstructure:"slow_threshold" json:"slow_threshold" yaml:"slow_threshold"`
	DSN           string         `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
	Options       map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}
