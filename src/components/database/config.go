package database

type Config struct {
	Driver        string         `mapstructure:"driver" json:"driver" yaml:"driver" binding:"required"`
	Host          string         `mapstructure:"host" json:"host" yaml:"host" binding:"required"`
	Port          uint           `mapstructure:"port" json:"port" yaml:"port"`
	Username      string         `mapstructure:"user_name" json:"user_name" yaml:"user_name" binding:"required"`
	Password      string         `mapstructure:"password" json:"password" yaml:"password" binding:"required"`
	DbName        string         `mapstructure:"db_name" json:"db_name" yaml:"db_name" binding:"required"`
	Charset       string         `mapstructure:"charset" json:"charset" yaml:"charset"`
	Prefix        string         `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	MaxIdleConn   uint           `mapstructure:"max_idle_conn" json:"max_idle_conn" yaml:"max_idle_conn"`
	MaxConn       uint           `mapstructure:"max_conn" json:"max_conn" yaml:"max_conn"`
	SlowThreshold uint64         `mapstructure:"slow_threshold" json:"slow_threshold" yaml:"slow_threshold"`
	Options       map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}
