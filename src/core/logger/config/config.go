package config

type Channel struct {
	Drivers []Driver `mapstructure:"drivers" json:"drivers" yaml:"drivers" binding:"required"`
}

type Driver struct {
	Driver     string         `mapstructure:"driver" json:"driver" yaml:"driver" binding:"required"`
	Path       string         `mapstructure:"path" json:"path" yaml:"path" binding:"required"`
	Level      string         `mapstructure:"level" json:"level" yaml:"level" binding:"required"`
	MaxDays    uint           `mapstructure:"days" json:"days" yaml:"days"`
	MaxSize    uint           `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups uint           `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	Options    map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}