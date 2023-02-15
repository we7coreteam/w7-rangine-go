package session

import (
	"time"
)

type Cookie struct {
	Name     string `mapstructure:"name" json:"name" yaml:"name"`
	Domain   string `mapstructure:"domain" json:"domain" yaml:"domain"`
	HttpOnly bool   `mapstructure:"http_only" json:"http_only" yaml:"http_only"`
	Path     string `mapstructure:"path" json:"path" yaml:"path"`
	Persist  bool   `mapstructure:"persist" json:"persist" yaml:"persist"`
	Secure   bool   `mapstructure:"secure" json:"secure" yaml:"secure"`
	SameSite string `mapstructure:"same_site" json:"same_site" yaml:"same_site"`
}

type SessionConf struct {
	Lifetime    time.Duration `mapstructure:"life_time" json:"life_time" yaml:"life_time"`
	IdleTimeout time.Duration `mapstructure:"idle_time" json:"idle_time" yaml:"idle_time"`
}

type Config struct {
	Session SessionConf `mapstructure:"session" json:"session" yaml:"session"`
	Cookie  Cookie      `mapstructure:"cookie" json:"cookie" yaml:"cookie"`
}
