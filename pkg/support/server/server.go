package server

type ServerInterface interface {
	GetServerName() string
	Start()
	GetOptions() map[string]string
}
