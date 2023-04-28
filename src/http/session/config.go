package session

type Cookie struct {
	Name     string         `mapstructure:"name" json:"name" yaml:"name"`
	MaxAge   int            `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	Domain   string         `mapstructure:"domain" json:"domain" yaml:"domain"`
	HttpOnly bool           `mapstructure:"http_only" json:"http_only" yaml:"http_only"`
	Path     string         `mapstructure:"path" json:"path" yaml:"path"`
	Secure   bool           `mapstructure:"secure" json:"secure" yaml:"secure"`
	SameSite int            `mapstructure:"same_site" json:"same_site" yaml:"same_site"`
	Options  map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}
