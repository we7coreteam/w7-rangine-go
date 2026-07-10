package session

import (
	"net/http"
	"strconv"
	"strings"
)

type Cookie struct {
	MaxAge   int            `mapstructure:"expires" json:"expires" yaml:"expires"`
	Domain   string         `mapstructure:"domain" json:"domain" yaml:"domain"`
	HttpOnly bool           `mapstructure:"http_only" json:"http_only" yaml:"http_only"`
	Path     string         `mapstructure:"path" json:"path" yaml:"path"`
	Secure   bool           `mapstructure:"secure" json:"secure" yaml:"secure"`
	SameSite string         `mapstructure:"same_site" json:"same_site" yaml:"same_site"`
	Options  map[string]any `mapstructure:"options" json:"options" yaml:"options"`
}

func ParseSameSite(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "", "default":
		return http.SameSiteDefaultMode
	case "lax":
		return http.SameSiteLaxMode
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		parsed, err := strconv.Atoi(value)
		if err != nil {
			return http.SameSiteDefaultMode
		}

		return http.SameSite(parsed)
	}
}
