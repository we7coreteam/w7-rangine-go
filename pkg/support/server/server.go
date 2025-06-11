package server

type ServerInterface interface {
	GetServerName() string
	Start()
	Stop()
	GetOptions() map[string]string
}
