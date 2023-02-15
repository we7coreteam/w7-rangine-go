package app

type Config struct {
	Env  string `mapstructure:"env" json:"env" yaml:"env"`
	Lang string `mapstructure:"lang" json:"lang" yaml:"lang"`
}
