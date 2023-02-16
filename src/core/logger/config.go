package logger

type Config struct {
	Driver     string                 `mapstructure:"driver" json:"driver" yaml:"driver"`
	Path       string                 `mapstructure:"path" json:"path" yaml:"path"`
	Level      string                 `mapstructure:"level" json:"level" yaml:"level"`
	MaxDays    int                    `mapstructure:"max_days" json:"max_days" yaml:"max_days"`
	MaxSize    int                    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	Options    map[string]interface{} `mapstructure:"options" json:"options" yaml:"options"`
	MaxBackups int                    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
}
