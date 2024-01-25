package server

type Server interface {
	GetServerName() string
	Start()
	Stop()
	GetOptions() map[string]string
}
