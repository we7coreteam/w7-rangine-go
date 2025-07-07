package server

type Server interface {
	GetServerName() string
	Start()
	GetOptions() map[string]string
}
