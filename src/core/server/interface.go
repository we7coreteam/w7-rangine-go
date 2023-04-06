package server

type Interface interface {
	GetServerName() string
	Start()
	GetOptions() map[string]string
}
