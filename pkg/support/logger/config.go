package logger

type Config struct {
	Driver     string         `mapstructure:"driver" json:"driver" yaml:"driver" binding:"required"`
	Path       string         `mapstructure:"path" json:"path" yaml:"path"`
	Level      string         `mapstructure:"level" json:"level" yaml:"level" binding:"required,oneof=debug info error warn dpanic panic fatal"`
	MaxDays    uint           `mapstructure:"days" json:"days" yaml:"days"`
	MaxSize    uint           `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups uint           `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	Options    map[string]any `mapstructure:"options" json:"options" yaml:"options"`
	Channels   []string       `mapstructure:"channels" json:"channels" yaml:"channels"`
}
