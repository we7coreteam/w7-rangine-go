package database

type Config struct {
	Driver        string `mapstructure:"driver" json:"driver" yaml:"driver"`
	Host          string `mapstructure:"host" json:"host" yaml:"host"`
	Port          int    `mapstructure:"port" json:"port" yaml:"port"`
	User          string `mapstructure:"user" json:"user" yaml:"user"`
	Password      string `mapstructure:"password" json:"password" yaml:"password"`
	DbName        string `mapstructure:"db_name" json:"db_name" yaml:"db_name"`
	Charset       string `mapstructure:"charset" json:"charset" yaml:"charset"`
	Prefix        string `mapstructure:"prefix" json:"prefix" yaml:"prefix"`
	MaxIdleConn   int    `mapstructure:"max_idel_conn" json:"max_idel_conn" yaml:"max_idel_conn"`
	MaxConn       int    `mapstructure:"max_conn" json:"max_conn" yaml:"max_conn"`
	SlowThreshold int64  `mapstructure:"slow_threshold" json:"slow_threshold" yaml:"slow_threshold"`
}
